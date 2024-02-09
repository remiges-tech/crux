package workflow

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
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

// WorkflowGet will be responsible for processing the /workflowget request that comes through as a POST
func WorkflowList(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("WorkflowList request received")
	var (
		request WorkflowListReq
		params  WorkflowListParams
		// slice      int32
		// app        string
		// class      string
		// name       string
		// active     bool
		// intrnl     bool
		dbResponse []sqlc.WorkflowListRow
	)

	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.LogActivity("error while binding json request error:", err.Error)
		return
	}

	// Check for validation error
	valError := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
	if len(valError) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, valError))
		lh.LogActivity("validation error:", valError)
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	// Populate parameters
	populateParams(&request, &params)

	// Check if the caller has root capabilities
	hasRootCapabilities := HasRootCapabilities()

	// Process the request based on the provided BRD
	dbResponse, err = processRequest(c, lh, hasRootCapabilities, query, params, &request)
	// dbResponse, err = query.WorkflowList(c, sqlc.WorkflowListParams{
	// 	Slice:      pgtype.Int4{Int32: *request.Slice, Valid: request.Slice != nil},
	// 	App:        pgtype.Text{String: *request.App, Valid: !isStringEmpty(request.App)},
	// 	Class:      pgtype.Text{String: *request.Class, Valid: !isStringEmpty(request.Class)},
	// 	Setname:    pgtype.Text{String: *request.Name, Valid: !isStringEmpty(request.Name)},
	// 	IsActive:   pgtype.Bool{Bool: *request.IsActive, Valid: request.IsActive != nil},
	// 	IsInternal: pgtype.Bool{Bool: *request.IsInternal, Valid: request.IsInternal != nil},
	// })

	if err != nil {
		if err.Error() == AUTH_ERROR {
			// Generate "auth" error
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Unauthorized, server.ErrCode_Unauthorized, nil)}))
			return
		}
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		// wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(types.RECORD_NOT_EXIST, nil)}))
		lh.LogActivity(DB_ERROR, err.Error)
		return
	}

	lh.Log(fmt.Sprintf("Record found: %v", map[string]any{"workflows": dbResponse}))

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
	lh.Log("processRequest request received")
	if !hasRootCapabilities {
		lh.Log("User not have root cap")
		//  if app parameter is present then
		if !isStringEmpty(request.App) {
			lh.Log("User has app params present")
			// check if named app matches = user has "ruleset" rights
			if HasRulesetRights(*request.App) {
				lh.Log("User has app rights")

				return query.WorkflowList(c, sqlc.WorkflowListParams{
					Slice: pgtype.Int4{Int32: *request.Slice, Valid: request.Slice != nil},
					// App:        pgtype.Text{String: *request.App, Valid: !isStringEmpty(request.App)},
					App:        []string{*request.App},
					Class:      pgtype.Text{String: *request.Class, Valid: !isStringEmpty(request.Class)},
					Setname:    pgtype.Text{String: *request.Name, Valid: !isStringEmpty(request.Name)},
					IsActive:   pgtype.Bool{Bool: *request.IsActive, Valid: request.IsActive != nil},
					IsInternal: pgtype.Bool{Bool: *request.IsInternal, Valid: request.IsInternal != nil},
				})
			}
			lh.Log("User not have 'ruleset' rights")
			// if user doesn't have "ruleset" rights then -> "auth" error
			return nil, fmt.Errorf(AUTH_ERROR)
		}
		// show the workflows of all the apps for which the user has "ruleset" rights
		app := GetWorkflowsByRulesetRights()
		// appStr := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprint(app), "[", "\""), "]", "\""), " ", "\",\"")
		lh.LogActivity("app not present hence all user 'ruleset' rights:", app)
		return query.WorkflowList(c, sqlc.WorkflowListParams{
			// App: pgtype.Text{String: app, Valid: isStringEmpty(request.App)},
			App: app,
		})
	}

	if !isStringEmpty(request.App) {
		lh.Log("User have root cap with 'app'")
		// the workflows of that app
		return query.WorkflowList(c, sqlc.WorkflowListParams{
			Slice:      pgtype.Int4{Int32: *request.Slice, Valid: request.Slice != nil},
			App:        []string{*request.App},
			Class:      pgtype.Text{String: *request.Class, Valid: !isStringEmpty(request.Class)},
			Setname:    pgtype.Text{String: *request.Name, Valid: !isStringEmpty(request.Name)},
			IsActive:   pgtype.Bool{Bool: *request.IsActive, Valid: request.IsActive != nil},
			IsInternal: pgtype.Bool{Bool: *request.IsInternal, Valid: request.IsInternal != nil},
		})
	}
	lh.Log("User have root cap and 'app' is nil")
	// the workflows of all the apps in the realm
	return query.WorkflowList(c, sqlc.WorkflowListParams{})
}

// to check given string is nil or not
func isStringEmpty(s *string) bool {
	return s == nil || strings.TrimSpace(*s) == ""
}

// to check if the user has "ruleset" rights for the given app
func HasRulesetRights(app string) bool {
	return true
}

// to check if the caller has root capabilities
func HasRootCapabilities() bool {
	return false
}

// to get workflows for all apps for which the user has "ruleset" rights
func GetWorkflowsByRulesetRights() []string {
	return []string{"retailBANK", "nedbank"}
	// return "nedbank"
}
