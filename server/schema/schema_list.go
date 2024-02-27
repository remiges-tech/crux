package schema

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"

	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

type SchemaStruct struct {
	Slice int32  `json:"slice" validate:"lt=15"`
	App   string `json:"app" validate:"lt=15"`
	Class string `json:"class" validate:"lt=15"`
}

var CapForList = []string{"ruleset", "schema", "root", "report"}
var schemaList []sqlc.WfSchemaListRow
var err error

func SchemaList(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("starting execution of SchemaList()")

	isCapable, capList := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: CapForList,
	}, false)

	if !isCapable {
		l.Info().LogActivity(server.ErrCode_Unauthorized, userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	var sh SchemaStruct

	err = wscutils.BindJSON(c, &sh)
	if err != nil {
		l.Debug0().Error(err).Log("error unmarshalling request payload to struct")
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
		l.Debug0().Log("error while getting query instance from service dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
		return
	}

	if slices.Contains(capList, "root") || slices.Contains(capList, "report") {
		schemaList, err = getSchemaList(c, sh, query)
	} else if slices.Contains(capList, "ruleset") || slices.Contains(capList, "schema") {
		if sh.App != "" {
			schemaList, err = getSchemaList(c, sh, query)
		} else {
			l.Info().LogActivity(server.ErrCode_Unauthorized, userID)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
			return
		}
	} else {
		l.Info().LogActivity(server.ErrCode_Unauthorized, userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}
	if err != nil || len(schemaList) == 0 {
		l.Debug0().Error(err).Log("error while retrieving from db")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	l.Debug0().Log("finished execution of SchemaList()")
}

func getSchemaList(c *gin.Context, sh SchemaStruct, query *sqlc.Queries) ([]sqlc.WfSchemaListRow, error) {
	return query.WfSchemaList(c, sqlc.WfSchemaListParams{
		Relam: realmID,
		Slice: pgtype.Int4{Int32: sh.Slice, Valid: sh.Slice > 0},
		App:   pgtype.Text{String: sh.App, Valid: !server.IsStringEmpty(&sh.App)},
		Class: pgtype.Text{String: sh.Class, Valid: !server.IsStringEmpty(&sh.Class)},
	})
}
