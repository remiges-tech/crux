package workflow

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
)

type WorkflowNew struct {
	Slice      int32         `json:"slice" validate:"required,gt=0,lt=15"`
	App        string        `json:"app" validate:"required,alpha,lt=15"`
	Class      string        `json:"class" validate:"required,lowercase,lt=15"`
	Name       string        `json:"name" validate:"required,lowercase,lt=15"`
	IsInternal bool          `json:"is_internal" validate:"required"`
	Flowrules  []crux.Rule_t `json:"flowrules" validate:"required,dive"`
}

func WorkFlowNew(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("Starting execution of WorkFlowNew()")

	// userID, err := server.ExtractUserNameFromJwt(c)
	// if err != nil {
	// 	l.Info().Log("unable to extract userID from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ERRCode_Token_Data_Missing))
	// 	return
	// }

	// realmName, err := server.ExtractRealmFromJwt(c)
	// if err != nil {
	// 	l.Info().Log("unable to extract realm from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ERRCode_Token_Data_Missing))
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
	var wf WorkflowNew

	err := wscutils.BindJSON(c, &wf)
	if err != nil {
		l.Error(err).Log("Error Unmarshalling Query parameters to struct:")
		return
	}

	validationErrors := wscutils.WscValidate(wf, func(err validator.FieldError) []string { return []string{} })
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
		Slice:     wf.Slice,
		App:       strings.ToLower(wf.App),
		Class:     wf.Class,
	})
	if err != nil {
		l.Info().Error(err).Log("failed to get schema from DB:")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	schema_t := crux.Schema_t{
		Class: wf.Class,
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
	customValidationErrors := customValidationErrors(schema_t, wf)
	validationErrors = append(validationErrors, customValidationErrors...)
	if len(validationErrors) > 0 {
		l.Debug0().LogDebug("custom validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	ruleset, err := json.Marshal(wf.Flowrules)
	if err != nil {
		patternSchema := "flowrules"
		l.Debug1().Error(err).Log("Error while marshaling Flowrules")
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_InvalidRequest, &patternSchema)}))
		return
	}

	err = qtx.WorkFlowNew(c, sqlc.WorkFlowNewParams{
		RealmName:  realmName,
		Slice:      wf.Slice,
		App:        strings.ToLower(wf.App),
		Brwf:       brwf,
		Class:      wf.Class,
		Setname:    wf.Name,
		Schemaid:   schema.ID,
		IsActive:   pgtype.Bool{Bool: false, Valid: true},
		IsInternal: wf.IsInternal,
		Ruleset:    ruleset,
		Createdby:  userID,
	})
	if err != nil {
		l.Info().LogActivity("Error while Inserting data in ruleset", err.Error())
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	if err := tx.Commit(c); err != nil {
		l.Info().Error(err).Log("Error while commits the transaction")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: nil, Messages: nil})
	l.Debug0().Log("Finished execution of WorkFlowNew()")
}

func customValidationErrors(schema_t crux.Schema_t, wf WorkflowNew) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	if len(wf.Flowrules) == 0 {
		fieldName := "flowrules"
		vErr := wscutils.BuildErrorMessage(server.MsgId_Empty, server.ErrCode_Empty, &fieldName)
		validationErrors = append(validationErrors, vErr)
		return validationErrors
	}
	ruleSet := crux.Ruleset_t{
		// Id:            0,
		Class:   wf.Class,
		SetName: wf.Name, //
		Rules:   wf.Flowrules,
		// NCalled:       0,
		// ReferenceType: "",
	}
	err := crux.VerifyRulePatterns(&ruleSet, &schema_t, true)
	if err != nil {
		rulePatternError := server.HandleCruxError(err)
		validationErrors = append(validationErrors, rulePatternError...)
	}

	err = crux.VerifyRuleActions(&ruleSet, &schema_t, true)
	if err != nil {
		ruleActionError := server.HandleCruxError(err)
		validationErrors = append(validationErrors, ruleActionError...)
	}

	return validationErrors
}
