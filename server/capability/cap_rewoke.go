package capability

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
)

type CapRevokeReq struct {
	User string   `json:"user" validate:"required"`
	App  []string `json:"app" validate:"omitempty,gt=0"`
	Cap  []string `json:"cap" validate:"required,gt=0"`
}

func CapRevoke(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Log("started execution of capRevoke()")

	// userID, err := server.ExtractUserNameFromJwt(c)
	// if err != nil {
	// 	l.Info().Log("unable to extract userID from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ERRCode_Token_Data_Missing))
	// 	return
	// }

	// realmName, err := server.ExtractRealmFromJwt(c)
	// if err != nil {
	// 	l.Info().Log("unable to extract realm from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ERRCode_Token_Data_Missing))
	// 	return
	// }

	// realmName, ok := s.Dependencies["realmName"].(string)
	// if !ok {
	// 	l.Debug0().Log("error while getting realmName instance from service dependencies")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
	// 	return
	// }

	reqCaps := []string{"auth"}
	isCapable, userCapabilities := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: reqCaps,
	}, false)

	if !isCapable {
		l.Info().LogActivity("Unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}
	var req CapRevokeReq
	err := wscutils.BindJSON(c, &req)
	if err != nil {
		l.Debug0().Error(err).Log("error while binding json request error:")
		return
	}

	valError := wscutils.WscValidate(req, func(err validator.FieldError) []string { return []string{} })
	if len(valError) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, valError))
		l.Debug0().LogActivity("validation error:", valError)
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Log("error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	valError = customValidation(c, query, req)
	if len(valError) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, valError))
		l.Debug0().LogActivity("validation error:", valError)
		return
	}
	var app []string
	var cap []string

	// check is there 'ALL' present in app []string
	if slices.Contains(req.App, capALL) {
		l.LogActivity("contains special capability name, ALL in app", req)
		app = nil
	} else {
		// Convert each app to lowercase
		for _, str := range req.App {
			app = append(app, strings.ToLower(str))
		}
		// check provided app exist or not
		appCount, err := query.AppExists(c, app)
		if err != nil {
			errMsg := db.HandleDatabaseError(err)
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errMsg}))
			return
			// no app found
		} else if appCount == 0 {
			feildName := "app"
			errMsg := wscutils.BuildErrorMessage(server.MsgId_NotFound, server.ErrCode_Invalid_APP, &feildName, req.App...)
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errMsg}))
			return
		}
	}

	// check is there 'ALL' present in cap []string
	if slices.Contains(req.Cap, capALL) {
		l.LogActivity("contains special capability name, ALL in cap", req)
		count, err := query.CountOfRootCapUser(c)
		if err != nil {
			l.LogActivity("Error while counting user with root cap", err.Error())
			errmsg := db.HandleDatabaseError(err)
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
			return
		}
		if count == 1 {
			l.Err().LogActivity("cap array includes root/all and system has only one user with root cap", req)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
			return
		}
		cap = nil
	} else {
		// Convert each cap to lowercase
		for _, str := range req.Cap {
			cap = append(cap, strings.ToLower(str))
		}

		capCount, err := query.CapExists(c, req.Cap)
		if err != nil {
			errMsg := db.HandleDatabaseError(err)
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errMsg}))
			return
		}
		if capCount == 0 {
			feildName := "cap"
			errMsg := wscutils.BuildErrorMessage(server.MsgId_NotFound, server.ErrCode_Invalid_Cap, &feildName, req.App...)
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errMsg}))
			return
		}
	}

	if slices.Contains(req.Cap, capRoot) && !slices.Contains(userCapabilities, capRoot) {
		l.Log("cap array includes root but the calling user has not root capability")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	} else if slices.Contains(req.Cap, capRoot) && slices.Contains(userCapabilities, capRoot) {
		count, err := query.CountOfRootCapUser(c)
		if err != nil {
			l.LogActivity("Error while counting user with root cap", err.Error())
			errmsg := db.HandleDatabaseError(err)
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
			return
		}
		if count == 1 {
			l.Err().LogActivity("cap array includes root and system has only one user with root cap", req)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
			return
		}
	}

	tag, err := query.CapRevoke(c, sqlc.CapRevokeParams{
		User: req.User,
		App:  app,
		Cap:  cap,
	})
	if err != nil {
		l.LogActivity("Error while deleting data in cap", err.Error())
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	if strings.Contains(tag.String(), "0") {
		l.Log("no record found to delete")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NotFound, server.ErrCode_NotFound))
		return
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: nil, Messages: nil})
	l.Debug0().Log("Finished execution of CapRevoke()")

}

func customValidation(c *gin.Context, queries *sqlc.Queries, req CapRevokeReq) []wscutils.ErrorMessage {
	var valError []wscutils.ErrorMessage
	userCount, err := queries.UserExists(c, req.User)
	if err != nil {
		errMsg := db.HandleDatabaseError(err)
		valError = append(valError, errMsg)
	}
	if userCount == 0 {
		feildName := "user"
		errMsg := wscutils.BuildErrorMessage(server.MsgId_NotFound, server.ErrCode_Invalid_User, &feildName, req.User)
		valError = append(valError, errMsg)
	}

	appCount, err := queries.AppExists(c, req.App)
	if err != nil {
		errMsg := db.HandleDatabaseError(err)
		valError = append(valError, errMsg)
	}
	if appCount == 0 {
		feildName := "app"
		errMsg := wscutils.BuildErrorMessage(server.MsgId_NotFound, server.ErrCode_Invalid_APP, &feildName, req.App...)
		valError = append(valError, errMsg)
	}

	capCount, err := queries.CapExists(c, req.Cap)
	if err != nil {
		errMsg := db.HandleDatabaseError(err)
		valError = append(valError, errMsg)
	}
	if capCount == 0 {
		feildName := "cap"
		errMsg := wscutils.BuildErrorMessage(server.MsgId_NotFound, server.ErrCode_Invalid_Cap, &feildName, req.App...)
		valError = append(valError, errMsg)
	}

	return valError

}
