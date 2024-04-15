package realmslice

import (
	
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

func RealmSliceApps(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("Starting execution of RealmSliceApps()")
	userID, err := server.ExtractUserNameFromJwt(c)
	if err != nil {
		l.Info().Log("unable to extract userID from token")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
		return
	}

	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: rootCaps,
	}, false)

	if !isCapable {
		l.Info().LogActivity("Unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		l.Error(err).Log("Error while parsing string param to int")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Invalid_Request, server.ErrCode_InvalidRequest))
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Info().Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
	}

	realmSliceAppsListRow, err := query.RealmSliceAppsList(c, int32(id))
	if err != nil {
		l.Info().Error(err).Log("Error while creating new realmslice")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	if realmSliceAppsListRow != nil {
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: realmSliceAppsListRow, Messages: nil})
		l.Debug0().Log("Finished execution of RealmSliceApps()")
		return
	} else {
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: []sqlc.RealmSliceAppsListRow{}, Messages: nil})
		l.Debug0().Log("Finished execution of RealmSliceApps()")
		return
	}
}
