package breruleset

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	crux "github.com/remiges-tech/crux/matching-engine"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
	"github.com/remiges-tech/logharbour/logharbour"
)

type RulesSetUpdate struct {
	Slice   int32         `json:"slice" validate:"required,gt=0,lt=50"`
	App     string        `json:"app" validate:"required,alpha,lt=50"`
	Class   string        `json:"class" validate:"required,lowercase,lt=50"`
	Name    string        `json:"name" validate:"required,lowercase,lt=50"`
	RuleSet []crux.Rule_t `json:"ruleset" validate:"required,dive"`
}

func RuleSetUpdate(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("Starting execution of RuleSetUpdate()")

	// userID, err := server.ExtractUserNameFromJwt(c)
	// if err != nil {
	// 	l.Info().Log("unable to extract userID from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	// realmName, err := server.ExtractRealmFromJwt(c)
	// if err != nil {
	// 	l.Info().Log("unable to extract realm from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	realmName, ok := s.Dependencies["realmName"].(string)
	if !ok {
		l.Debug0().Log("error while getting realmName instance from service dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
		return
	}

	// delete below line whie actual implementation (reason: kept for testing while writting api)
	// realmName = "Ecommerce"

	capNeeded := []string{"ruleset"}
	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: capNeeded,
	}, false)

	if !isCapable {
		l.Info().LogActivity("Unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	var req RulesSetUpdate

	err := wscutils.BindJSON(c, &req)
	if err != nil {
		l.Error(err).Log("Error Unmarshalling Query parameters to struct:")
		return
	}

	validationErrors := wscutils.WscValidate(req, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		l.Debug0().LogDebug("standard validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Debug0().Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	connpool, ok := s.Database.(*pgxpool.Pool)
	if !ok {
		l.Debug0().Log("Error while getting connection pool instance from service Database")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	tx, err := connpool.Begin(c)
	if err != nil {
		l.LogActivity("Error while Begin tx", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	defer tx.Rollback(c)
	qtx := query.WithTx(tx)

	schema, err := qtx.GetSchemaWithLock(c, sqlc.GetSchemaWithLockParams{
		RealmName: realmName,
		Slice:     req.Slice,
		App:       strings.ToLower(req.App),
		Class:     req.Class,
		Brwf:      sqlc.BrwfEnumB,
	})
	if err != nil {
		l.LogActivity("failed to get schema from DB:", err.Error())
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	schema_t := crux.Schema_t{
		Class: req.Class,
	}

	err = json.Unmarshal(schema.Patternschema, &schema_t.PatternSchema)
	if err != nil {
		l.Debug1().Error(err).Log("Error while Unmarshalling PatternSchema")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	err = json.Unmarshal(schema.Actionschema, &schema_t.ActionSchema)
	if err != nil {
		l.Debug1().Error(err).Log("Error while Unmarshaling ActionSchema")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	// custom Validation
	customValidationErrors := customValidationErrors(schema_t, req)
	validationErrors = append(validationErrors, customValidationErrors...)
	if len(validationErrors) > 0 {
		l.Debug0().LogDebug("custom validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	ruleSet, err := json.Marshal(req.RuleSet)
	if err != nil {
		patternSchema := "ruleSet"
		l.Debug1().LogDebug("Error while marshaling ruleSet", err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_InvalidRequest, &patternSchema)}))
		return
	}

	ruleset, err := qtx.RulesetRowLock(c, sqlc.RulesetRowLockParams{
		RealmName: realmName,
		Slice:     req.Slice,
		App:       strings.ToLower(req.App),
		Class:     req.Class,
	})
	if err != nil {
		l.LogActivity("Error while locking row of ruleset", err.Error())
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	tag, err := qtx.WorkFlowUpdate(c, sqlc.WorkFlowUpdateParams{
		RealmName: realmName,
		Slice:     req.Slice,
		App:       strings.ToLower(req.App),
		Brwf:      sqlc.BrwfEnumB,
		Class:     req.Class,
		Setname:   req.Name,
		Ruleset:   ruleSet,
		Editedby:  pgtype.Text{String: userID, Valid: true},
	})
	if err != nil {
		l.LogActivity("Error while Updating data in ruleset", err.Error())
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	if strings.Contains(tag.String(), "0") {
		l.Log("no record found to update")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NotFound, server.ErrCode_No_record_Found))
		return
	}
	if err := tx.Commit(c); err != nil {
		l.LogActivity("Error while commits the transaction", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	dclog := l.WithClass("ruleset").WithInstanceId(string(ruleset.ID))
	dclog.LogDataChange("Updated ruleset", logharbour.ChangeInfo{
		Entity: "ruleset",
		Op:     "Update",
		Changes: []logharbour.ChangeDetail{
			{
				Field:  "brwf",
				OldVal: ruleset.Brwf,
				NewVal: brwf},
			{
				Field:  "setname",
				OldVal: ruleset.Setname,
				NewVal: req.Name},
			{
				Field:  "ruleset",
				OldVal: string(ruleset.Ruleset),
				NewVal: req.RuleSet},
		},
	})
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: nil, Messages: nil})
	l.Debug0().Log("Finished execution of RuleSetUpdate()")

}

func customValidationErrors(schema_t crux.Schema_t, rs RulesSetUpdate) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	if len(rs.RuleSet) == 0 {
		fieldName := "RuleSet"
		vErr := wscutils.BuildErrorMessage(server.MsgId_Empty, server.ErrCode_Empty, &fieldName)
		validationErrors = append(validationErrors, vErr)
		return validationErrors
	}
	ruleSet := crux.Ruleset_t{
		// Id:            0,
		Class:   rs.Class,
		SetName: rs.Name, //
		Rules:   rs.RuleSet,
		// NCalled:       0,
		// ReferenceType: "",
	}
	err := crux.VerifyRulePatterns(&ruleSet, &schema_t, false)
	if err != nil {
		rulePatternError := server.HandleCruxError(err)
		validationErrors = append(validationErrors, rulePatternError...)
	}

	err = crux.VerifyRuleActions(&ruleSet, &schema_t, false)
	if err != nil {
		ruleActionError := server.HandleCruxError(err)
		validationErrors = append(validationErrors, ruleActionError...)
	}

	return validationErrors
}
