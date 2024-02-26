package schema

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
)

func BRESchemaList(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Log("Starting execution of BRESchemaList()")
	var sh SchemaListStruct

	err := wscutils.BindJSON(c, &sh)
	if err != nil {
		l.LogActivity("Error Unmarshalling Query parameter to struct:", err.Error())
		return
	}

	validationErrors := wscutils.WscValidate(sh, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}
	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Log("Error while getting query instance from service Dependencies")
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

}
