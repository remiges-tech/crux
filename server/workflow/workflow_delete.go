package workflow

import (
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
	"github.com/remiges-tech/logharbour/logharbour"
)

// WorkflowDelete will be responsible for processing the /workflowdelete request that comes through as a DELETE
func WorkflowDelete(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("WorkflowDelete request received")

	var (
		request WorkflowGetReq
	)

	// implement the user realm and all here
	var capForList = []string{"workflow"}

	userID, err := server.ExtractUserNameFromJwt(c)
	if err != nil {
		lh.Info().Log("unable to extract userID from token")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
		return
	}

	// realmName, err := server.ExtractRealmFromJwt(c)
	// if err != nil {
	// 	lh.Info().Log("unable to extract realm from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	realmName, ok := s.Dependencies["realmName"].(string)
	if !ok {
		lh.Debug0().Log("error while getting realmName instance from service dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
		return
	}

	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: capForList,
	}, false)

	if !isCapable {
		lh.Info().LogActivity(server.ErrCode_Unauthorized, userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	err = wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Debug0().Error(err).Log("error while binding json request error")
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
		lh.Log("error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	dbRecord, err := query.Workflowget(c, sqlc.WorkflowgetParams{
		Slice:   request.Slice,
		App:     request.App,
		Class:   request.Class,
		Setname: request.Name,
		Realm:   realmName,
		Brwf:    sqlc.BrwfEnumW,
	})
	if err != nil {
		lh.Debug0().Error(err).Log("error while retriving record")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	tag, err := query.WorkflowDelete(c, sqlc.WorkflowDeleteParams{
		Slice:   request.Slice,
		App:     request.App,
		Class:   request.Class,
		Setname: request.Name,
		Realm:   realmName,
		Brwf:    sqlc.BrwfEnumW,
	})
	if err != nil {
		lh.Debug0().Error(err).Log("failed to delete data from db")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	// if '1' contains means db transaction done
	if strings.Contains(tag.String(), "1") {
		// do data change log
		dataChangeLog(lh, dbRecord)
		lh.Debug0().Log("record found finished execution of WorkflowDelete()")
		wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(nil))
		return
	}
	lh.Debug0().LogActivity("failed to delete data from db:", tag.String())
	wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_InvalidRequest, nil)}))
}

func dataChangeLog(lh *logharbour.Logger, dbRecord sqlc.WorkflowgetRow) {
	dclog := lh.WithClass("ruleset").WithInstanceId(string(dbRecord.ID))
	dclog.LogDataChange("delete ruleset", logharbour.ChangeInfo{
		Entity: "ruleset",
		Op:     "delete",
		Changes: []logharbour.ChangeDetail{
			{
				Field:  "row",
				OldVal: dbRecord,
				NewVal: nil},
		},
	})
}

// to check if the user has "schema" capability for the app this workflow belongs to
func HasSchemaCap(app string) bool {
	userRights := server.GetWorkflowsByRulesetRights()
	return slices.Contains(userRights, app)
}
