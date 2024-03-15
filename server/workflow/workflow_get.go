package workflow

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

// WorkflowGet will be responsible for processing the /workflowget request that comes through as a POST
func WorkflowGet(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Debug0().Log("WorkflowGet request received")
	var (
		request WorkflowGetReq
	)
	userID, err := server.ExtractUserNameFromJwt(c)
	if err != nil {
		lh.Info().Log("unable to extract userID from token")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ERRCode_Token_Data_Missing))
		return
	}

	realmName, err := server.ExtractRealmFromJwt(c)
	if err != nil {
		lh.Info().Log("unable to extract realm from token")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ERRCode_Token_Data_Missing))
		return
	}
	// implement the user realm and all here
	var (
		capForList = []string{"workflow"}
	)
	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: capForList,
	}, false)

	if !isCapable {
		lh.Debug0().LogActivity(server.ErrCode_Unauthorized, userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	err = wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Debug0().Error(err).Log("error while binding json request")
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
		lh.Debug0().Log("error while getting query instance from service dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	dbResponse, err := query.Workflowget(c, sqlc.WorkflowgetParams{
		Slice:   request.Slice,
		App:     request.App,
		Class:   request.Class,
		Setname: request.Name,
		Realm:   realmName,
	})
	if err != nil {
		lh.Debug0().Error(err).Log("failed to get data from db")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	actualResponse := responseBinding(dbResponse)

	err = json.Unmarshal(dbResponse.Flowrules, &actualResponse.Flowrules)
	if err != nil {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, nil)}))
		lh.Debug0().Error(err).Log("failed to unmarshal data")
		return
	}
	lh.Debug0().Log("record found finished execution of WorkflowGet()")
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(actualResponse))
}

func responseBinding(dbResponse sqlc.WorkflowgetRow) WorkflowgetRow {
	actualResponse := WorkflowgetRow{
		ID:         dbResponse.ID,
		Slice:      dbResponse.Slice,
		App:        dbResponse.App,
		Class:      dbResponse.Class,
		Name:       dbResponse.Name,
		IsActive:   dbResponse.IsActive.Bool,
		IsInternal: dbResponse.IsInternal,
		Createdat:  dbResponse.Createdat,
		Createdby:  dbResponse.Createdby,
		Editedat:   dbResponse.Editedat,
		Editedby:   dbResponse.Editedby,
	}
	return actualResponse
}
