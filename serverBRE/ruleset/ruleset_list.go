package ruleset

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
	"github.com/remiges-tech/logharbour/logharbour"
)

func BRERuleSetList(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("RuleSetList request received")
	var (
		request    RuleSetListReq
		dbResponse []sqlc.WorkflowListRow
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

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	hasRootCapabilities := HasRootCapabilities()

	dbResponse, err = processRequest(c, lh, hasRootCapabilities, query, &request)

	if err != nil {
		if err.Error() == AuthError {
			// Generate "auth" error
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Unauthorized, server.ErrCode_Unauthorized, nil)}))
			return
		}
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		lh.LogActivity(DbError, err.Error)
		return
	}

	lh.Log(fmt.Sprintf("Record found: %v", map[string]any{"RuleSets": dbResponse}))

	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(map[string][]sqlc.WorkflowListRow{"rulesets": dbResponse}))
}

// Function to process the request and get the RuleSet
func processRequest(c *gin.Context, lh *logharbour.Logger, hasRootCapabilities bool, query *sqlc.Queries, request *RuleSetListReq) ([]sqlc.WorkflowListRow, error) {
	lh.Log("processRequest request received")
	// implement the user realm here
	var userRealm int32 = 1

	if !hasRootCapabilities {
		lh.Log("User not have root cap")
		//  if app parameter is present then
		if !types.IsStringEmpty(request.App) {
			lh.Log("User has app params present")
			// check if named app matches = user has "ruleset" rights
			if types.HasRulesetRights(*request.App) {
				lh.Log("User has app rights")
				return query.WorkflowList(c, sqlc.WorkflowListParams{
					Slice:      pgtype.Int4{Int32: *request.Slice, Valid: request.Slice != nil},
					App:        []string{*request.App},
					Realm:      userRealm,
					Class:      pgtype.Text{String: *request.Class, Valid: !types.IsStringEmpty(request.Class)},
					Setname:    pgtype.Text{String: *request.Name, Valid: !types.IsStringEmpty(request.Name)},
					IsActive:   pgtype.Bool{Bool: *request.IsActive, Valid: request.IsActive != nil},
					IsInternal: pgtype.Bool{Bool: *request.IsInternal, Valid: request.IsInternal != nil},
				})
			}
			lh.Log("User not have 'rulesets' rights")
			// if user doesn't have "ruleset" rights then -> "auth" error
			return nil, fmt.Errorf(AuthError)
		}
		// show the RuleSet of all the apps for which the user has "ruleset" rights
		app := GeRuleSetsByRulesetRights()
		lh.LogActivity("app not present hence all user 'ruleset' rights:", app)
		return query.WorkflowList(c, sqlc.WorkflowListParams{
			App:   app,
			Realm: userRealm,
		})
	}

	if !types.IsStringEmpty(request.App) {
		lh.Log("User have root cap with 'app'")
		// the RuleSet of that app
		return query.WorkflowList(c, sqlc.WorkflowListParams{
			Slice:      pgtype.Int4{Int32: *request.Slice, Valid: request.Slice != nil},
			App:        []string{*request.App},
			Realm:      userRealm,
			Class:      pgtype.Text{String: *request.Class, Valid: !types.IsStringEmpty(request.Class)},
			Setname:    pgtype.Text{String: *request.Name, Valid: !types.IsStringEmpty(request.Name)},
			IsActive:   pgtype.Bool{Bool: *request.IsActive, Valid: request.IsActive != nil},
			IsInternal: pgtype.Bool{Bool: *request.IsInternal, Valid: request.IsInternal != nil},
		})
	}
	lh.Log("User have root cap and 'app' is nil")
	// the RuleSet of all the apps in the realm
	return query.WorkflowList(c, sqlc.WorkflowListParams{Realm: userRealm})
}
