package schema

import (
	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
)

type SchemaListStruct struct {
	Slice *int32  `form:"slice"`
	App   *string `form:"app"`
	Class *string `form:"class"`
}

func SchemaList(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Log("Starting execution of SchemaList()")
	var sh SchemaListStruct

	//check the deactivated table to check whether first, the realm and then, the user has been
	//  deactivated
	// isDeactivated()

	// check the capgrant table to see if the calling user has the capability to perform the
	// operation
	// isCapable, _ := utils.Authz_check(types.OpReq{
	// 	User:      username,
	// 	CapNeeded: []string{"report","ruleset","schema"},
	// }, false)

	// if !isCapable {
	// 	l.Log("Unauthorized user:")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(utils.ErrUnauthorized))
	// 	return
	// }

	err := c.ShouldBindQuery(&sh)
	if err != nil {
		l.LogActivity("Error Unmarshalling Query paramaeters to struct:", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeInvalidJson))
		return
	}
	query, ok := s.Database.(*sqlc.Queries)
	if !ok {
		l.Log("Error while getting query instance from service")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse("unable_to_get_database_query"))
		return
	}

	switch true {
	case sh.App != nil && sh.Class == nil && sh.Slice == nil:
		schemaList, err := query.SchemaListByApp(c, *sh.App)
		if err != nil || len(schemaList) == 0 {
			l.LogActivity("Error while retrieving schema list by app", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	case sh.Class != nil && sh.App == nil && sh.Slice == nil:
		schemaList, err := query.SchemaListByClass(c, *sh.Class)
		if err != nil || len(schemaList) == 0 {
			l.LogActivity("Error while retrieving schema list by class", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	case sh.Slice != nil && sh.Class == nil && sh.App == nil:
		schemaList, err := query.SchemaListBySlice(c, *sh.Slice)
		if err != nil || len(schemaList) == 0 {
			l.LogActivity("Error while retrieving schema list by slice", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	case sh.App != nil && sh.Class != nil:
		schemaList, err := query.SchemaListByAppAndClass(c, sqlc.SchemaListByAppAndClassParams{App: *sh.App, Class: *sh.Class})
		if err != nil || len(schemaList) == 0 {
			l.LogActivity("Error while retrieving schema list by app & class", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	case sh.App != nil && sh.Slice != nil:
		schemaList, err := query.SchemaListByAppAndSlice(c, sqlc.SchemaListByAppAndSliceParams{App: *sh.App, Slice: *sh.Slice})
		if err != nil || len(schemaList) == 0 {
			l.LogActivity("Error while retrieving schema list by app & slice", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	case sh.Class != nil && sh.Slice != nil:
		schemaList, err := query.SchemaListByClassAndSlice(c, sqlc.SchemaListByClassAndSliceParams{Class: *sh.Class, Slice: *sh.Slice})
		if err != nil || len(schemaList) == 0 {
			l.LogActivity("Error while retrieving schema list by class & slice", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	case sh.App != nil && sh.Class != nil && sh.Slice != nil:
		Schema, err := query.SchemaGet(c, sqlc.SchemaGetParams{App: *sh.App, Class: *sh.Class, Slice: *sh.Slice})
		if err != nil || Schema != nil {
			l.LogActivity("Error while retrieving schema list by app, class & slice", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: Schema, Messages: nil})

	default:
		schemaList, err := query.SchemaList(c)
		if err != nil || len(schemaList) == 0 {
			l.LogActivity("Error while retrieving schema list by app, class & slice", err)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	}

}
