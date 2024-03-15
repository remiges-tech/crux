package app

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
	"github.com/remiges-tech/logharbour/logharbour"
)

// AppDelete will be responsible for processing the /appdelete request that comes through as a POST
func AppDelete(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour.WithClass("app")
	lh.Log("AppDelete request received")

	var (
		appName = c.Param("name")
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
	err := c.ShouldBindQuery(appName)
	if err != nil {
		lh.Error(err).Log("AppDelete() || error while binding  app name")
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Debug0().Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
		return
	}

	// validate app name
	customError := validateAppName(appName)
	if len(customError) > 0 {
		lh.Debug0().LogActivity("AppDelete()||error occurred while validating app name", customError)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, customError))
		return
	}

	// To verify whether app name already exist in database
	applc := strings.ToLower(appName)
	appData, err := query.GetAppName(c, sqlc.GetAppNameParams{
		Shortnamelc: applc,
		Realm:       REALM,
	})
	if err != nil {
		lh.Info().Error(err).Log("error while getting app details if present")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	if len(appData) == 0 {
		lh.Debug0().LogActivity("app name does not exist in db", appName)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NotFound, server.ERRCode_Name_Not_Exist))
		return

	}
	// To verify whether app name is present in schema and ruleset table
	count, err := query.AppExist(c, applc)
	if err != nil {
		lh.Info().Error(err).Log("error while verifying app name present in schema or ruleset table")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	if count == 1 {
		// wscutils.SendErrorResponse(c,wscutils.NewResponse( wscutils.BuildErrorMessage(server.MsgId__NonEmpty, server.ErrCode_NonEmpty, &APP))
		wscutils.SendErrorResponse(c, &wscutils.Response{Status: wscutils.ErrorStatus, Data: nil, Messages: []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId__NonEmpty, server.ErrCode_NonEmpty, &APP)}})
		return
	}
	// To get cap grants for app
	capGrantData, err := query.GetCapGrantForApp(c, sqlc.GetCapGrantForAppParams{
		App:   pgtype.Text{String: applc, Valid: true},
		Realm: REALMID,
	})
	if err != nil {
		lh.Info().Error(err).Log("error while getting app capablity grants from db")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	if len(capGrantData) > 0 {
		// If the app is to be deleted, and there are cap grants to users for this app in the capgrant
		// table, then those capability grants are deleted too.
		error := query.DeleteCapGranForApp(c, sqlc.DeleteCapGranForAppParams{
			App:   pgtype.Text{String: applc, Valid: true},
			Realm: REALMID,
		})
		if error != nil {
			lh.Info().Error(err).Log("error while deleting app capablity grants from db")
			errmsg := db.HandleDatabaseError(err)
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
			return
		}
		// data change log
		for _, val := range capGrantData {
			dclog := lh.WithClass("cap grants for app").WithInstanceId(strconv.Itoa(int(val.ID)))
			dclog.LogDataChange("delete cap grants for app", logharbour.ChangeInfo{
				Entity: "capgrant",
				Op:     "delete",
				Changes: []logharbour.ChangeDetail{
					{
						Field:  "row",
						OldVal: val,
						NewVal: nil},
				},
			})
		}

	}

	// delete app
	err = query.AppDelete(c, sqlc.AppDeleteParams{
		Shortnamelc: applc,
		Realm:       REALM,
	})
	if err != nil {
		lh.Info().Error(err).Log("error while deleting app from db")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	// data change log
	for _, val := range appData {
		dclog := lh.WithClass("app").WithInstanceId(strconv.Itoa(int(val.ID)))
		dclog.LogDataChange("delete app", logharbour.ChangeInfo{
			Entity: "app",
			Op:     "delete",
			Changes: []logharbour.ChangeDetail{
				{
					Field:  "row",
					OldVal: val,
					NewVal: nil},
			},
		})
	}

	lh.Debug0().Log("finished execution of AppDelete()")
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(nil))

}
