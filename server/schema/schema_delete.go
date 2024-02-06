package schema

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/types"
)

// SchemaDelete will be responsible for processing the /WFschemaDelete request that comes through as a POST
func SchemaDelete(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("SchemaDelete request received")

	// var response schemaGetResp
	var request SchemaGetReq
	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Debug0().LogActivity("error while binding json request error:", err.Error)
		return
	}

	valError := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
	if len(valError) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, valError))
		lh.Debug0().LogActivity("validation error:", valError)
		return
	}
	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}
	err = query.Wfschemadelete(c, sqlc.WfschemadeleteParams{
		Slice: request.Slice,
		App:   request.App,
		Class: request.Class,
	})
	if err != nil {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(types.OPERATION_FAILED, nil)}))
		lh.Debug0().LogActivity("failed while deleting record:", err.Error)
		return
	}

	// lh.ChangeMinLogPriority(logharbour.LogPriority(logharbour.Change))
	lh.Log(fmt.Sprintf("Record delete: %v", map[string]any{"err": err}))
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(err))
}
