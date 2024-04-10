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
	"github.com/remiges-tech/logharbour/logharbour"
)

type WorkflowUpdate struct {
	Slice     int32         `json:"slice" validate:"required,gt=0,lt=15"`
	App       string        `json:"app" validate:"required,alpha,lt=15"`
	Class     string        `json:"class" validate:"required,lowercase,lt=15"`
	Name      string        `json:"name" validate:"required,lowercase,lt=15"`
	Flowrules []crux.Rule_t `json:"flowrules" validate:"required,dive"`
}

func WorkFlowUpdate(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("Starting execution of WorkflowUpdate()")

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

	var wf WorkflowUpdate
	// var ruleSchema schema.SchemaNewReq

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
		l.LogActivity("Error while Begin tx", err.Error())
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
		l.LogActivity("failed to get schema from DB:", err.Error())
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
		l.Debug1().Error(err).Log("Error while Unmarshaling ActionSchema")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	// custom Validation
	customValidationErrors := customValidationErrors(schema_t, WorkflowNew{
		Slice:     wf.Slice,
		App:       strings.ToLower(wf.App),
		Class:     wf.Class,
		Name:      wf.Name,
		Flowrules: wf.Flowrules,
	})
	validationErrors = append(validationErrors, customValidationErrors...)
	if len(validationErrors) > 0 {
		l.Debug0().LogDebug("custom validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	flowrules, err := json.Marshal(wf.Flowrules)
	if err != nil {
		patternSchema := "flowrules"
		l.Debug1().LogDebug("Error while marshaling Flowrules", err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_InvalidRequest, &patternSchema)}))
		return
	}

	ruleset, err := qtx.RulesetRowLock(c, sqlc.RulesetRowLockParams{
		RealmName: realmName,
		Slice:     wf.Slice,
		App:       strings.ToLower(wf.App),
		Class:     wf.Class,
	})
	if err != nil {
		l.LogActivity("Error while locking row of ruleset", err.Error())
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	tag, err := qtx.WorkFlowUpdate(c, sqlc.WorkFlowUpdateParams{
		RealmName: realmName,
		Slice:     wf.Slice,
		App:       strings.ToLower(wf.App),
		Brwf:      brwf,
		Class:     wf.Class,
		Setname:   wf.Name,
		Ruleset:   flowrules,
		Editedby:  pgtype.Text{String: editedBy, Valid: true},
	})
	if err != nil {
		l.LogActivity("Error while Updating data in ruleset", err.Error())
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	if strings.Contains(tag.String(), "0") {
		l.Log("no record found to update")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NotFound, server.ErrCode_NotFound))
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
				NewVal: wf.Name},
			{
				Field:  "ruleset",
				OldVal: string(ruleset.Ruleset),
				NewVal: wf.Flowrules},
		},
	})
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: nil, Messages: nil})
	l.Debug0().Log("Finished execution of WorkflowUpdate()")

}
