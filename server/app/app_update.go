package app

import (
	"strconv"
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

// Incoming request format
type GetAppUpdateRequest struct {
	Name        string `json:"name" validate:"required,lt=25"`
	Description string `json:"descr" validate:"required,gt=0,lt=40"`
}

// AppUpdate will be responsible for processing the /appUpdate request that comes through as a POST
func AppUpdate(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour.WithClass("app")
	lh.Log("AppUpdate request received")

	var (
		request GetAppUpdateRequest
	)

	isCapable, _ := server.Authz_check(types.OpReq{
		User:      USERID,
		CapNeeded: rootCapability,
	}, false)

	if !isCapable {
		lh.Info().LogActivity("unauthorized user:", USERID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	// Bind request
	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Error(err).Log("AppUpdate() || error while binding json request")
		return
	}
	// Standard validation of Incoming Request
	validationErrors := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		lh.Debug0().LogActivity("AppUpdate()||validation error:", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}
	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Debug0().Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
		return
	}
	// validate app name
	customError := validateAppName(request.Name)
	if len(customError) > 0 {
		lh.Debug0().LogActivity("AppUpdate()||error occurred while validating app name", customError)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, customError))
		return
	}

	// To check whether app name already exist in database
	appName := strings.ToLower(request.Name)
	appData, err := query.GetAppName(c, sqlc.GetAppNameParams{
		Shortnamelc: appName,
		Realm:       REALM,
	})
	if err != nil {
		lh.Info().Error(err).Log("Error while getting app details if present")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	if len(appData) == 0 {
		lh.Debug0().LogActivity("app name does not exist in db", request.Name)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NotFound, server.ERRCode_Name_Not_Exist))
		return

	}

	// Update record
	err = query.AppUpdate(c, sqlc.AppUpdateParams{
		Longname:    request.Description,
		Setby:       USERID,
		Shortnamelc: appName,
		Realm:       REALM,
	})
	if err != nil {
		lh.Info().Error(err).Log("error while updating app")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	// data change log
	for _, val := range appData {
		dclog := lh.WithClass("app").WithInstanceId(strconv.Itoa(int(val.ID)))
		dclog.LogDataChange("update app ", logharbour.ChangeInfo{
			Entity: "app",
			Op:     "update",
			Changes: []logharbour.ChangeDetail{
				{
					Field:  "longname",
					OldVal: val.Longname,
					NewVal: request.Description,
				}, {
					Field:  "setby",
					OldVal: val.Setby,
					NewVal: USERID,
				},
			},
		})
	}
	lh.Debug0().Log("finished execution of AppUpdate()")
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(nil))

}
