package capability

import (
	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

// CapGet: will handle "/capget/:userid" GET
func CapGet(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("starting execution of CapGet()")

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

	// Step:1 - get params / json binding
	opUserId := c.Param("userid")
	if opUserId == "" {
		l.Log("no operational user id found")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NotFound, server.ERRCode_User_Id_Not_Exist))
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
	dbResponse, err := query.CapGet(c, sqlc.CapGetParams{
		Userid: opUserId,
		Realm:  realmName,
	})
	if err != nil {
		l.Info().Error(err).Log("error while getting data from db with func CapGet")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	// Step:4 - send response
	l.Debug0().Log("exiting from CapGet()")
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(map[string][]sqlc.CapGetRow{"capabilities": dbResponse}))
}
