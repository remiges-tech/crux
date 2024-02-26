package workflow

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
	"github.com/remiges-tech/logharbour/logharbour"
)

type WorkflowListParams struct {
	Slice      int32
	App        string
	Class      string
	Name       string
	IsActive   bool
	IsInternal bool
}

// WorkflowList will be responsible for processing the /WorkflowList request that comes through as a POST
func WorkflowList(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("WorkflowList request received")

	// implement the user realm and all here
	var (
		userID     = "1234"
		capForList = []string{"workflow"}
		// userRealm  int32 = 1 //this is implemented in processRequest() below
	)
	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: capForList,
	}, false)

	if !isCapable {
		lh.Info().LogActivity(server.ErrCode_Unauthorized, userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	var (
		request    WorkflowListReq
		params     WorkflowListParams
		dbResponse []sqlc.WorkflowListRow
	)

	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Debug0().Error(err).Log("error while binding json request")
		return
	}

	// Check for validation error
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
	// Populate parameters
	populateParams(&request, &params)

	// Check if the caller has root capabilities
	hasRootCapabilities := HasRootCapabilities()

	// Process the request based on the provided BRD
	dbResponse, err = processRequest(c, lh, hasRootCapabilities, query, params, &request)

	if err != nil {
		if err.Error() == AUTH_ERROR {
			// Generate "auth" error
			lh.Debug0().Error(err).Log("error while getting query instance from service dependencies")
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Unauthorized, server.ErrCode_Unauthorized, nil)}))
			return
		}
		lh.Debug0().Error(err).Log(DB_ERROR)
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	lh.Debug0().Log("record found finished execution of WorkflowList()")
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(map[string][]sqlc.WorkflowListRow{"workflows": dbResponse}))
}

// Function to populate parameters from the request
func populateParams(request *WorkflowListReq, params *WorkflowListParams) {
	if request.Slice != nil {
		params.Slice = *request.Slice
	}

	if request.App != nil {
		params.App = *request.App
	}
	if request.Class != nil {
		params.Class = *request.Class
	}
	if request.Name != nil {
		params.Name = *request.Name
	}
	if request.IsActive != nil {
		params.IsActive = *request.IsActive
	}
	if request.IsInternal != nil {
		params.IsInternal = *request.IsInternal
	}
}

// Function to process the request and get the workflows
func processRequest(c *gin.Context, lh *logharbour.Logger, hasRootCapabilities bool, query *sqlc.Queries, params WorkflowListParams, request *WorkflowListReq) ([]sqlc.WorkflowListRow, error) {
	lh.Debug0().Log("processRequest request received")

	// implement the user realm here
	var userRealm int32 = 1

	if !hasRootCapabilities {
		lh.Debug0().Log("user not have root cap")
		//  if app parameter is present then
		if !server.IsStringEmpty(request.App) {
			lh.Debug0().Log("user has app params present")
			// check if named app matches = user has "ruleset" rights
			if server.HasRulesetRights(*request.App) {
				lh.Debug0().Log("user has app rights")
				return query.WorkflowList(c, sqlc.WorkflowListParams{
					Slice:      pgtype.Int4{Int32: *request.Slice, Valid: request.Slice != nil},
					App:        []string{*request.App},
					Realm:      userRealm,
					Class:      pgtype.Text{String: *request.Class, Valid: !server.IsStringEmpty(request.Class)},
					Setname:    pgtype.Text{String: *request.Name, Valid: !server.IsStringEmpty(request.Name)},
					IsActive:   pgtype.Bool{Bool: *request.IsActive, Valid: request.IsActive != nil},
					IsInternal: pgtype.Bool{Bool: *request.IsInternal, Valid: request.IsInternal != nil},
				})
			}
			lh.Debug0().Log("user not have 'ruleset' rights")
			// if user doesn't have "ruleset" rights then -> "auth" error
			return nil, fmt.Errorf(AUTH_ERROR)
		}
		// show the workflows of all the apps for which the user has "ruleset" rights
		app := server.GetWorkflowsByRulesetRights()
		lh.Debug0().LogActivity("app not present hence all user 'ruleset' rights:", app)
		return query.WorkflowList(c, sqlc.WorkflowListParams{
			App:   app,
			Realm: userRealm,
		})
	}

	if !server.IsStringEmpty(request.App) {
		lh.Debug0().Log("user have root cap with 'app'")
		// the workflows of that app
		return query.WorkflowList(c, sqlc.WorkflowListParams{
			Slice:      pgtype.Int4{Int32: *request.Slice, Valid: request.Slice != nil},
			App:        []string{*request.App},
			Realm:      userRealm,
			Class:      pgtype.Text{String: *request.Class, Valid: !server.IsStringEmpty(request.Class)},
			Setname:    pgtype.Text{String: *request.Name, Valid: !server.IsStringEmpty(request.Name)},
			IsActive:   pgtype.Bool{Bool: *request.IsActive, Valid: request.IsActive != nil},
			IsInternal: pgtype.Bool{Bool: *request.IsInternal, Valid: request.IsInternal != nil},
		})
	}
	lh.Debug0().Log("user have root cap or 'app' is nil")
	// the workflows of all the apps in the realm
	return query.WorkflowList(c, sqlc.WorkflowListParams{Realm: userRealm})
}
