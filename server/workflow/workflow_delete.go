package workflow

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

// WorkflowDelete will be responsible for processing the /workflowdelete request that comes through as a DELETE
func WorkflowDelete(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("WorkflowDelete request received")

	var (
		request WorkflowGetReq
		// implement the user realm here
		userRealm int32 = 1
	)

	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.LogActivity("error while binding json request error:", err.Error)
		return
	}

	valError := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
	if len(valError) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, valError))
		lh.LogActivity("validation error:", valError)
		return
	}
	if !HasSchemaCap(request.App) {
		// Generate "auth" error if no
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Unauthorized, server.ErrCode_Unauthorized, nil)}))
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	tag, err := query.WorkflowDelete(c, sqlc.WorkflowDeleteParams{
		Slice:   request.Slice,
		App:     request.App,
		Class:   request.Class,
		Setname: request.Name,
		Realm:   userRealm,
	})
	if err != nil {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, nil)}))
		lh.LogActivity("failed to delete data from DB:", err.Error)
		return
	}
	if strings.Contains(tag.String(), "1") {
		lh.Log(fmt.Sprintf("Record found: %v", map[string]any{"response": tag.String()}))
		wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse("record_deleted_success"))
		return
	}
	wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_InvalidRequest, nil)}))
	lh.LogActivity("failed to delete data from DB:", tag.String())
}

// to check if the user has "schema" capability for the app this workflow belongs to
func HasSchemaCap(app string) bool {
	userRights := types.GetWorkflowsByRulesetRights()
	return slices.Contains(userRights, app)
}
