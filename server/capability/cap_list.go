package capability

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

type CapListReq struct {
	App []string `json:"app,omitempty"`
	Cap []string `json:"cap,omitempty"`
}

// CapList: will handle "/caplist" GET
func CapList(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("starting execution of CapList()")

	var (
		request CapListReq
	)

	// userID, err := server.ExtractUserNameFromJwt(c)
	// if err != nil {
	// 	l.Info().Log("unable to extract userID from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	// realmName, err := server.ExtractRealmFromJwt(c)
	// if err != nil {
	// 	l.Info().Log("unable to extract realm from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	// Step:0 - check user caps
	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: authCaps,
	}, false)

	if !isCapable {
		l.Info().LogActivity("unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	err := wscutils.BindJSON(c, &request)
	if err != nil {
		l.Debug0().Error(err).Log("error while binding json request")
		return
	}

	// Check for validation error
	validationErrors := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		l.Debug0().LogDebug("standard validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	// Step:2 - get required dependency
	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Info().Log("error while getting query instance from service dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
		return
	}

	// Step:3 - do the db transaction
	dbResponse, err := query.CapList(c, sqlc.CapListParams{
		Realm: realmName,
		App:   request.App,
		Cap:   request.Cap,
	})

	if err != nil {
		l.Info().Error(err).Log("error while getting data from db with func CapList")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	// Step:4 - send response
	l.Debug0().Log("exiting from CapList()")
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(map[string][]sqlc.CapListRow{"capabilities": dbResponse}))
}
