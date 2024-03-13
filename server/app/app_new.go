package app

import (
	"fmt"
	"regexp"
	"slices"
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

var RESERVED_APPNAMES = []string{"ALL"}

// Incoming request format
type GetAppNewRequest struct {
	Name        string `json:"name" validate:"required,lt=25"`
	Description string `json:"descr" validate:"required,gt=0,lt=40"`
}

// AppNew will be responsible for processing the /appnew request that comes through as a POST
func AppNew(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour.WithClass("app")
	lh.Log("AppNew request received")

	var (
		request   GetAppNewRequest
		capForNew = []string{"root"}
	)

	isCapable, _ := server.Authz_check(types.OpReq{
		User:      USERID,
		CapNeeded: capForNew,
	}, false)

	if !isCapable {
		lh.Info().LogActivity("unauthorized user:", USERID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	// Bind request
	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Error(err).Log("AppNew() || error while binding json request error")
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
	err = validateAppName(request.Name)
	if err != nil {
		lh.Error(err).Log("AppNew()|| error occurred while validating app name")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Invalid, server.ErrCode_Invalid_NAME))
		return
	}

	// Insert record
	applc := strings.ToLower(request.Name)
	appData, err := query.AppNew(c, sqlc.AppNewParams{
		Realm:       REALM,
		Shortname:   request.Name,
		Shortnamelc: applc,
		Longname:    request.Description,
		Setby:       ADMIN,
	})

	if err != nil {
		lh.Info().Error(err).Log("Error while inserting app")
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
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(nil))

}

// validate app name
func validateAppName(app string) error {

	// Check if the app name is one-word and follows identifier syntax rules
	pattern := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	if !pattern.MatchString(app) || strings.Contains(app, " ") {
		return fmt.Errorf("%v must be a one-word and follows identifier syntax rules", app)
	}

	// Check if the app name is reserved
	App := strings.ToUpper(app)

	isReservedName := slices.Contains(RESERVED_APPNAMES, App)
	if isReservedName {
		return fmt.Errorf("reserved app name")
	}
	return nil

}
