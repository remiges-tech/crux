package schema

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"

	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

type SchemaStruct struct {
	Slice int32  `json:"slice"`
	App   string `json:"app"`
	Class string `json:"class"`
}

var capForList = []string{"ruleset", "schema", "root", "report"}
var schemaList []sqlc.WfSchemaListRow
var err error

func SchemaList(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("Starting execution of SchemaList()")

	isCapable, capList := types.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: capForList,
	}, false)

	if !isCapable {
		l.Info().LogActivity("Unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	var sh SchemaStruct

	err = wscutils.BindJSON(c, &sh)
	if err != nil {
		l.LogActivity("Error Unmarshalling request payload to struct:", err.Error())
		return
	}

	// Validate request
	validationErrors := wscutils.WscValidate(sh, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		l.Debug0().LogDebug("standard validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}
	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Debug0().Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
		return
	}

	if slices.Contains(capList, "root") || slices.Contains(capList, "report") {
		schemaList, err = getSchemaList(c, sh, query)
	} else if slices.Contains(capList, "ruleset") || slices.Contains(capList, "schema") {
		if sh.App != "" {
			schemaList, err = getSchemaList(c, sh, query)
		} else {
			l.Info().LogActivity("Unauthorized user:", userID)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
			return
		}
	} else {
		l.Info().LogActivity("Unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}
	if err != nil || len(schemaList) == 0 {
		l.LogActivity("Error while retrieving schema list", err)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NoSchemaFound, server.ErrCode_Invalid))
		return
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	l.Debug0().Log("Finished execution of SchemaNew()")
}

func getSchemaList(c *gin.Context, sh SchemaStruct, query *sqlc.Queries) ([]sqlc.WfSchemaListRow, error) {
	return query.WfSchemaList(c, sqlc.WfSchemaListParams{
		Relam: realmID,
		Slice: pgtype.Int4{Int32: sh.Slice, Valid: sh.Slice > 0},
		App:   pgtype.Text{String: sh.App, Valid: !types.IsStringEmpty(&sh.App)},
		Class: pgtype.Text{String: sh.Class, Valid: !types.IsStringEmpty(&sh.Class)},
	})
}
