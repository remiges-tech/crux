package schema

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

type updateSchema struct {
	Slice         int32               `json:"slice" validate:"required,gt=0,lt=50"`
	App           string              `json:"App" validate:"required,alpha,lt=50"`
	Class         string              `json:"class" validate:"required,lowercase,lt=50"`
	PatternSchema []PatternSchema     `json:"patternSchema,omitempty"`
	ActionSchema  crux.ActionSchema_t `json:"actionSchema,omitempty"`
}

func SchemaUpdate(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("Starting execution of SchemaUpdate()")

	// userID, err := server.ExtractUserNameFromJwt(c)
	// if err != nil {
	// 	l.Info().Log("unable to extract userID from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	realmName, ok := s.Dependencies["realmName"].(string)
	if !ok {
		l.Debug0().Log("error while getting realmName instance from service dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
		return
	}

	capForUpdate := []string{"schema"}
	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: capForUpdate,
	}, false)

	if !isCapable {
		l.Info().LogActivity("Unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	var req updateSchema

	err = wscutils.BindJSON(c, &req)
	if err != nil {
		l.Error(err).Log("Error Unmarshalling Query paramaeters to struct:")
		return
	}

	newPatternSchema := convertPatternSchema(req.PatternSchema)
	schema := crux.Schema_t{
		Class:         req.Class,
		PatternSchema: newPatternSchema,
		ActionSchema:  req.ActionSchema,
		NChecked:      0,
	}

	// Validate request
	validationErrors := wscutils.WscValidate(req, func(err validator.FieldError) []string { return []string{} })
	customValidationErrors := customValidationErrors(schema)
	validationErrors = append(validationErrors, customValidationErrors...)
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
		l.Debug0().Log("Error while getting query instance from service Database")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	patternSchema, err := json.Marshal(newPatternSchema)
	if err != nil {
		patternSchema := "patternSchema"
		l.Debug1().LogDebug("Error while marshaling patternSchema", err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_InvalidJson, &patternSchema)}))
		return
	}

	actionSchema, err := json.Marshal(req.ActionSchema)
	if err != nil {
		actionSchema := "actionSchema"
		l.Debug1().LogDebug("Error while marshaling actionSchema", err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_InvalidJson, &actionSchema)}))
		return
	}

	tx, err := connpool.Begin(c)
	if err != nil {
		l.Info().Error(err).Log("Error while beginning transaction")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	defer tx.Rollback(c)
	qtx := query.WithTx(tx)

	tag, err := qtx.IsWorkflowReferringSchema(c, sqlc.IsWorkflowReferringSchemaParams{
		Realm: realmName,
		Slice: req.Slice,
		App:   req.App,
		Class: req.Class,
	})
	if err != nil {
		tx.Rollback(c)
		l.Info().Error(err).Log("Error while checking does any workflow refers this schema")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	if tag != 0 {
		l.Info().Log("cannot update schema because workflows referring it")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId__NotAllowed, server.ErrCode_NotAllowed))
		return
	}

	getSchema, err := qtx.GetSchemaWithLock(c, sqlc.GetSchemaWithLockParams{
		RealmName: realmName,
		Slice:     req.Slice,
		Class:     req.Class,
		App:       strings.ToLower(req.App),
		Brwf:      sqlc.BrwfEnumW,
	})
	if err != nil {
		tx.Rollback(c)
		l.Info().Error(err).Log("Error while locking schema to get old value")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
		RealmName:     realmName,
		Slice:         req.Slice,
		Class:         req.Class,
		App:           strings.ToLower(req.App),
		Brwf:          sqlc.BrwfEnumW,
		Patternschema: patternSchema,
		Actionschema:  actionSchema,
		Editedby:      pgtype.Text{String: userID, Valid: true},
	})
	if err != nil {
		tx.Rollback(c)
		l.Info().Error(err).Log("Error while updating schema")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	if err := tx.Commit(c); err != nil {
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	dclog := l.WithClass("schema").WithInstanceId(string(getSchema.ID))
	dclog.LogDataChange("Updated schema", logharbour.ChangeInfo{
		Entity: "schema",
		Op:     "Update",
		Changes: []logharbour.ChangeDetail{
			{
				Field:  "patternSchema",
				OldVal: string(getSchema.Patternschema),
				NewVal: newPatternSchema},
			{
				Field:  "actionSchema",
				OldVal: string(getSchema.Actionschema),
				NewVal: req.ActionSchema},
		},
	})
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: nil, Messages: nil})
	l.Debug0().Log("Finished execution of SchemaUpdate()")
}
