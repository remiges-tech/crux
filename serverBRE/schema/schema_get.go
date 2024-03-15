package schema

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	constants "github.com/remiges-tech/crux/serverBRE"
)

func BRESchemaGet(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("SchemaGet request received")

	var request SchemaGetReq
	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Debug0().LogActivity("error while binding json request error:", err.Error())
		return
	}

	valError := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
	if len(valError) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, valError))
		lh.Debug0().LogActivity("validation error:", valError)
		return
	}
	query, ok := s.Dependencies[constants.Queries].(*sqlc.Queries)
	if !ok {
		lh.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	dbResponse, err := query.Wfschemaget(c, sqlc.WfschemagetParams{
		Slice: request.Slice,
		App:   request.App,
		Class: request.Class,
	})
	if err != nil {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_NoSchemaFound, server.ErrCode_Invalid, nil)}))
		lh.Debug0().LogActivity("failed to get data from DB:", err.Error())
		return
	}

	lh.Log(fmt.Sprintf("Record found: %v", map[string]any{"response": dbResponse}))
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(dbResponse))
}
