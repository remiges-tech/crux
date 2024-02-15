package schema

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"

	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

var capForList = []string{"report", "ruleset", "schema"}

func SchemaList(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("Starting execution of SchemaList()")

	isCapable, _ := types.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: capForList,
	}, false)

	if !isCapable {
		l.Info().LogActivity("Unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	var sh SchemaListStruct

	err := wscutils.BindJSON(c, &sh)
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

	switch true {
	case sh.App != nil && sh.Class == nil && sh.Slice == nil:
		schemaList, err := query.SchemaListByApp(c, *sh.App)
		if err != nil || len(schemaList) == 0 {
			l.LogActivity("Error while retrieving schema list by app", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NoSchemaFound, server.ErrCode_Invalid))
			return
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	case sh.Class != nil && sh.App == nil && sh.Slice == nil:
		schemaList, err := query.SchemaListByClass(c, *sh.Class)
		if err != nil || len(schemaList) == 0 {
			l.LogActivity("Error while retrieving schema list by class", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NoSchemaFound, server.ErrCode_Invalid))
			return
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	case sh.Slice != nil && sh.Class == nil && sh.App == nil:
		schemaList, err := query.SchemaListBySlice(c, *sh.Slice)
		if err != nil || len(schemaList) == 0 {
			l.LogActivity("Error while retrieving schema list by slice", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NoSchemaFound, server.ErrCode_Invalid))
			return
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	case sh.App != nil && sh.Class != nil && sh.Slice == nil:
		schemaList, err := query.SchemaListByAppAndClass(c, sqlc.SchemaListByAppAndClassParams{App: *sh.App, Class: *sh.Class})
		if err != nil || len(schemaList) == 0 {
			l.LogActivity("Error while retrieving schema list by app & class", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NoSchemaFound, server.ErrCode_Invalid))
			return
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	case sh.App != nil && sh.Slice != nil && sh.Class == nil:
		schemaList, err := query.SchemaListByAppAndSlice(c, sqlc.SchemaListByAppAndSliceParams{App: *sh.App, Slice: *sh.Slice})
		if err != nil || len(schemaList) == 0 {
			l.LogActivity("Error while retrieving schema list by app & slice", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NoSchemaFound, server.ErrCode_Invalid))
			return
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	case sh.Class != nil && sh.Slice != nil && sh.App == nil:
		schemaList, err := query.SchemaListByClassAndSlice(c, sqlc.SchemaListByClassAndSliceParams{Class: *sh.Class, Slice: *sh.Slice})
		if err != nil || len(schemaList) == 0 {
			l.LogActivity("Error while retrieving schema list by class & slice", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NoSchemaFound, server.ErrCode_Invalid))
			return
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	case sh.App != nil && sh.Class != nil && sh.Slice != nil:
		Schema, err := query.SchemaGet(c, sqlc.SchemaGetParams{App: *sh.App, Class: *sh.Class, Slice: *sh.Slice})
		if err != nil {
			l.LogActivity("Error while retrieving schema list by app, class & slice", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NoSchemaFound, server.ErrCode_Invalid))
			return
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: Schema, Messages: nil})

	default:
		schemaList, err := query.SchemaList(c)
		if err != nil || len(schemaList) == 0 {
			l.LogActivity("Error while retrieving schema list by app, class & slice", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NoSchemaFound, server.ErrCode_Invalid))
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	}
	l.Debug0().Log("Finished execution of SchemaNew()")
}
