package workflow

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
)

// WorkflowGet will be responsible for processing the /workflowget request that comes through as a POST
func WorkflowList(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("WorkflowList request received")
	var (
		request    WorkflowListReq
		slice      int32
		app        string
		class      string
		name       string
		active     bool
		intrnl     bool
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
	if request.Slice != nil {
		slice = *request.Slice
	}

	if request.App != nil {
		app = *request.App
	}
	if request.Class != nil {
		class = *request.Class
	}
	if request.Name != nil {
		class = *request.Name
	}
	if request.IsActive != nil {
		active = *request.IsActive
	}
	if request.IsInternal != nil {
		intrnl = *request.IsInternal
	}

	// Check if the caller has root capabilities
	hasRootCapabilities := false // HasCapabilities()

	if !hasRootCapabilities {
		if request.App != nil {
			// Check if the named app matches the name of one of the apps for which the user has "ruleset" rights
			if HasRulesetRights(*request.App) {
				dbResponse, err = query.WorkflowList(c, sqlc.WorkflowListParams{
					// Slice:      pgtype.Int4{Int32: slice, Valid: request.Slice != nil},
					App: pgtype.Text{String: app, Valid: !isStringEmpty(request.App)},
					// Class:      pgtype.Text{String: class, Valid: !isStringEmpty(request.Class)},
					// Setname:    pgtype.Text{String: name, Valid: !isStringEmpty(request.Name)},
					// IsActive:   pgtype.Bool{Bool: active, Valid: request.IsActive != nil},
					// IsInternal: pgtype.Bool{Bool: intrnl, Valid: request.IsInternal != nil},
				})

				if err != nil {
					wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_NotFound, server.ErrCode_NotFound, nil)}))
					lh.LogActivity(DB_ERROR, err.Error)
					return
				}
			} else {
				// Generate "auth" error
				wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Unauthorized, server.ErrCode_Unauthorized, nil)}))
				return
			}
		} else {
			// The user does not have root capabilities, and app parameter is not present
			dbResponse, err = query.WorkflowList(c, sqlc.WorkflowListParams{
				Slice:      pgtype.Int4{Int32: slice, Valid: request.Slice != nil},
				App:        pgtype.Text{String: app, Valid: !isStringEmpty(request.App)},
				Class:      pgtype.Text{String: class, Valid: !isStringEmpty(request.Class)},
				Setname:    pgtype.Text{String: name, Valid: !isStringEmpty(request.Name)},
				IsActive:   pgtype.Bool{Bool: active, Valid: request.IsActive != nil},
				IsInternal: pgtype.Bool{Bool: intrnl, Valid: request.IsInternal != nil},
			})

			if err != nil {
				wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_NotFound, server.ErrCode_NotFound, nil)}))
				lh.LogActivity(DB_ERROR, err.Error)
				return
			}
		}
	} else {
		if request.App != nil {
			// The user does not have root capabilities, and app parameter is not present
			dbResponse, err = query.WorkflowList(c, sqlc.WorkflowListParams{
				Slice:      pgtype.Int4{Int32: slice, Valid: request.Slice != nil},
				App:        pgtype.Text{String: app, Valid: !isStringEmpty(request.App)},
				Class:      pgtype.Text{String: class, Valid: !isStringEmpty(request.Class)},
				Setname:    pgtype.Text{String: name, Valid: !isStringEmpty(request.Name)},
				IsActive:   pgtype.Bool{Bool: active, Valid: request.IsActive != nil},
				IsInternal: pgtype.Bool{Bool: intrnl, Valid: request.IsInternal != nil},
			})

			if err != nil {
				wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_NotFound, server.ErrCode_NotFound, nil)}))
				lh.LogActivity(DB_ERROR, err.Error)
				return
			}
		} else {
			// Check if the named app matches the name of one of the apps for which the user has "ruleset" rights
			if HasRulesetRights(*request.App) {
				dbResponse, err = query.WorkflowList(c, sqlc.WorkflowListParams{
					// Slice:      pgtype.Int4{Int32: slice, Valid: request.Slice != nil},
					App: pgtype.Text{String: app, Valid: !isStringEmpty(request.App)},
					// Class:      pgtype.Text{String: class, Valid: !isStringEmpty(request.Class)},
					// Setname:    pgtype.Text{String: name, Valid: !isStringEmpty(request.Name)},
					// IsActive:   pgtype.Bool{Bool: active, Valid: request.IsActive != nil},
					// IsInternal: pgtype.Bool{Bool: intrnl, Valid: request.IsInternal != nil},
				})

				if err != nil {
					wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_NotFound, server.ErrCode_NotFound, nil)}))
					lh.LogActivity(DB_ERROR, err.Error)
					return
				}
			} else {
				// Generate "auth" error
				wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Unauthorized, server.ErrCode_Unauthorized, nil)}))
				return
			}
		}
	}

	// lh.Log(fmt.Sprintf("Record found: %v", map[string]any{"response": tempData}))

	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(map[string][]sqlc.WorkflowListRow{"workflows": dbResponse}))
}

func isStringEmpty(s *string) bool {
	return s == nil || strings.TrimSpace(*s) == ""
}

func HasRulesetRights(app string) bool {
	// Implement logic to check if the user has "ruleset" rights for the given app
	return true
}
