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

var (
	RESERVED_APPNAMES = []string{"ALL"}
)

// Incoming request format
type GetAppNewRequest struct {
	Name        string `json:"name" validate:"required,lt=25"`
	Description string `json:"descr" validate:"required,gt=0,lt=40"`
}

// AppNew will be responsible for processing the /appnew request that comes through as a POST
func AppNew(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour.WithClass("app")
	lh.Log("AppNew request received")

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

	var (
		request GetAppNewRequest
	)

	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: rootCapability,
	}, false)

	if !isCapable {
		lh.Info().LogActivity("unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	// Bind request
	err = wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Error(err).Log("AppNew() || error while binding json request ")
		return
	}
	// Standard validation of Incoming Request
	validationErrors := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		lh.Debug0().LogActivity("AppNew()||validation error:", validationErrors)
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
		lh.Debug0().LogActivity("AppNew()||error occurred while validating app name", customError)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, customError))
		return
	}

	// Insert record
	applc := strings.ToLower(request.Name)
	appData, err := query.AppNew(c, sqlc.AppNewParams{
		Realm:       realmName,
		Shortname:   request.Name,
		Shortnamelc: applc,
		Longname:    request.Description,
		Setby:       userID,
	})

	if err != nil {
		lh.Info().Error(err).Log("error while inserting app")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	// data change log
	for _, val := range appData {
		dclog := lh.WithClass("app").WithInstanceId(strconv.Itoa(int(val.ID)))
		dclog.LogDataChange("created app ", logharbour.ChangeInfo{
			Entity: "app",
			Op:     "create",
			Changes: []logharbour.ChangeDetail{
				{
					Field:  "row",
					OldVal: nil,
					NewVal: val},
			},
		})
	}
	lh.Debug0().Log("finished execution of AppNew()")
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(nil))

}
