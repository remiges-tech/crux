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

	// userID, err := server.ExtractUserNameFromJwt(c)
	// if err != nil {
	// 	lh.Info().Log("unable to extract userID from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	// realmName, err := server.ExtractRealmFromJwt(c)
	// if err != nil {
	// 	lh.Info().Log("unable to extract realm from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	// implement the user realm and all here
	var (
		capForList = []string{"workflow"}
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
		dbResponse []sqlc.WorkflowListRow
	)

	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Debug0().Error(err).Log("error while binding json request")
		return
	}

	// Check for validation error
	validationErrors := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		lh.Debug0().LogDebug("standard validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Debug0().Log("error while getting query instance from service dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	// Check if the caller has root capabilities
	hasRootCapabilities := HasRootCapabilities()

	// Process the request based on the provided BRD
	dbResponse, err = processRequest(c, lh, hasRootCapabilities, query, &request)

	if err != nil {
		if err.Error() == AUTH_ERROR {
			// Generate "auth" error
			lh.Debug0().Error(err).Log(AUTH_ERROR)
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

// Function to process the request and get the workflows
func processRequest(c *gin.Context, lh *logharbour.Logger, hasRootCapabilities bool, query *sqlc.Queries, request *WorkflowListReq) ([]sqlc.WorkflowListRow, error) {
	lh.Debug0().Log("processRequest request received")
	var (
		isAct, isIntr bool
	)
	if request.IsActive != nil {
		isAct = *request.IsActive
	}
	if request.IsInternal != nil {
		isIntr = *request.IsInternal
	}
	if !hasRootCapabilities {
		lh.Debug0().Log("user not have root cap")
		//  if app parameter is present then
		if !server.IsStringEmpty(&request.App) {
			lh.Debug0().Log("user has app params present")
			// check if named app matches = user has "ruleset" rights
			if server.HasRulesetRights(request.App) {
				lh.Debug0().Log("user has app rights")
				return query.WorkflowList(c, sqlc.WorkflowListParams{
					Brwf:       sqlc.BrwfEnumW,
					Slice:      pgtype.Int4{Int32: request.Slice, Valid: request.Slice > 0},
					App:        []string{request.App},
					Realm:      realmName,
					Class:      pgtype.Text{String: request.Class, Valid: !server.IsStringEmpty(&request.Class)},
					Setname:    pgtype.Text{String: request.Name, Valid: !server.IsStringEmpty(&request.Name)},
					IsActive:   pgtype.Bool{Bool: isAct, Valid: (request.IsActive != nil)},
					IsInternal: pgtype.Bool{Bool: isIntr, Valid: (request.IsInternal != nil)},
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
			Realm: realmName,
			Brwf:       sqlc.BrwfEnumW,
		})
	}

	if !server.IsStringEmpty(&request.App) {
		lh.Debug0().Log("user have root cap with 'app'")
		// the workflows of that app
		return query.WorkflowList(c, sqlc.WorkflowListParams{
			Brwf:       sqlc.BrwfEnumW,
			Realm:      realmName,
			Slice:      pgtype.Int4{Int32: request.Slice, Valid: request.Slice > 0},
			App:        []string{request.App},
			Class:      pgtype.Text{String: request.Class, Valid: !server.IsStringEmpty(&request.Class)},
			Setname:    pgtype.Text{String: request.Name, Valid: !server.IsStringEmpty(&request.Name)},
			IsActive:   pgtype.Bool{Valid: request.IsActive != nil, Bool: isAct},
			IsInternal: pgtype.Bool{Valid: request.IsInternal != nil, Bool: isIntr},
		})
	}
	lh.Debug0().Log("user have root cap or 'app' is nil")
	// the workflows of all the apps in the realm
	return query.WorkflowList(c, sqlc.WorkflowListParams{Realm: realmName, Brwf: sqlc.BrwfEnumW})
}
