package breruleset

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

type RuleSetNew struct {
	Slice      int32         `json:"slice" validate:"required,gt=0"`
	App        string        `json:"app" validate:"required,alpha"`
	Class      string        `json:"class" validate:"required,lowercase"`
	Name       string        `json:"name" validate:"required,lowercase"`
	IsInternal bool          `json:"is_internal" validate:"required"`
	RulseSet   []crux.Rule_t `json:"ruleset" validate:"required,dive"`
}

func BRERuleSetNew(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Log("Starting execution of BRERuleSetNew()")

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
	var request RuleSetNew

	err := wscutils.BindJSON(c, &request)
	if err != nil {
		l.Error(err).Log("Error Unmarshalling Query parameters to struct:")
		return
	}

	validationErrors := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
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
		l.Info().Error(err).Log("Error while Begin tx")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	defer tx.Rollback(c)
	qtx := query.WithTx(tx)

	schema, err := qtx.GetSchemaWithLock(c, sqlc.GetSchemaWithLockParams{
		RealmName: realmName,
		Slice:     request.Slice,
		App:       strings.ToLower(request.App),
		Class:     request.Class,
		Brwf:      sqlc.BrwfEnumB,
	})
	if err != nil {
		l.Info().Error(err).Log("failed to get schema from DB:")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	schema_t := crux.Schema_t{
		Class: request.Class,
	}
	err = json.Unmarshal([]byte(schema.Patternschema), &schema_t.PatternSchema)
	if err != nil {
		l.Debug1().Error(err).Log("Error while Unmarshalling PatternSchema")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	err = json.Unmarshal(schema.Actionschema, &schema_t.ActionSchema)
	if err != nil {
		l.Debug1().LogDebug("Error while Unmarshaling ActionSchema", err)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	// custom Validation
	customValidationErrors := customValidationErrorsForRulesetNew(schema_t, request)
	validationErrors = append(validationErrors, customValidationErrors...)
	if len(validationErrors) > 0 {
		l.Debug0().LogDebug("custom validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	ruleSetByte, err := json.Marshal(request.RulseSet)
	if err != nil {
		patternSchema := "RuleSet"
		l.Debug1().Error(err).Log("Error while marshaling ruleset")
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_InvalidRequest, &patternSchema)}))
		return
	}

	id, err := qtx.WorkFlowNew(c, sqlc.WorkFlowNewParams{
		RealmName:  realmName,
		Slice:      request.Slice,
		App:        strings.ToLower(request.App),
		Brwf:       sqlc.BrwfEnumB,
		Class:      request.Class,
		Setname:    request.Name,
		Schemaid:   schema.ID,
		IsInternal: request.IsInternal,
		Ruleset:    ruleSetByte,
		Createdby:  userID,
	})
	if err != nil {
		l.Info().LogActivity("Error while Inserting data in request", err.Error())
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	if err := tx.Commit(c); err != nil {
		l.Info().Error(err).Log("Error while commits the transaction")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	dclog := l.WithClass("BRERuleSet").WithInstanceId(string(id))
	dclog.LogDataChange("insert request", logharbour.ChangeInfo{
		Entity: "BRERuleSet",
		Op:     "insert",
		Changes: []logharbour.ChangeDetail{
			{
				Field:  "realmName",
				OldVal: nil,
				NewVal: realmName,
			},
			{
				Field:  "Slice",
				OldVal: nil,
				NewVal: request.Slice,
			},
			{
				Field:  "App",
				OldVal: nil,
				NewVal: strings.ToLower(request.App),
			},
			{
				Field:  "brwf",
				OldVal: nil,
				NewVal: sqlc.BrwfEnumB,
			},
			{
				Field:  "Class",
				OldVal: nil,
				NewVal: request.Class,
			},
			{
				Field:  "setname",
				OldVal: nil,
				NewVal: request.Name,
			},
			{
				Field:  "Schemaid",
				OldVal: nil,
				NewVal: schema.ID,
			},
			{
				Field:  "IsInternal",
				OldVal: nil,
				NewVal: request.IsInternal,
			},
			{
				Field:  "request",
				OldVal: nil,
				NewVal: request.RulseSet,
			},
			{
				Field:  "Createdby",
				OldVal: nil,
				NewVal: userID,
			},
		},
	})

	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: nil, Messages: nil})
	l.Debug0().Log("Finished execution of BRERuleSetNew()")
}

func customValidationErrorsForRulesetNew(schema_t crux.Schema_t, r RuleSetNew) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	if len(r.RulseSet) == 0 {
		fieldName := "RuleSet"
		vErr := wscutils.BuildErrorMessage(server.MsgId_Empty, server.ErrCode_Empty, &fieldName)
		validationErrors = append(validationErrors, vErr)
		return validationErrors
	}
	ruleSet := crux.Ruleset_t{
		Class:   r.Class,
		SetName: r.Name,
		Rules:   r.RulseSet,
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
