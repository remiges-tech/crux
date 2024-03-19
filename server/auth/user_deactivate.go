package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

// UserActivate: will handle "/userdeactivate/:userid" POST
func UserDeactivate(c *gin.Context, s *service.Service) {
	from_t := time.Now()
	// uncomment below time while running test case
	// from_t, _ := time.Parse("2006-01-02T15:04:05Z", "2021-12-01T14:30:15Z")
	l := s.LogHarbour
	l.Debug0().Log("starting execution of UserDeactivate()")

	userID, err := server.ExtractUserNameFromJwt(c)
	if err != nil {
		l.Info().Log("unable to extract userID from token")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ERRCode_Token_Data_Missing))
		return
	}

	realmName, err := server.ExtractRealmFromJwt(c)
	if err != nil {
		l.Info().Log("unable to extract realm from token")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ERRCode_Token_Data_Missing))
		return
	}

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
	newSliceID, err := query.UserDeactivate(c, sqlc.UserDeactivateParams{
		Userid:       opUserId,
		Deactivateat: pgtype.Timestamp{Time: from_t, Valid: true},
		Realm:        realmName,
	})
	if err != nil {
		l.Info().Error(err).Log("error while changing active status in db with func UserDeactivate")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	err = query.DeactivateRecord(c, sqlc.DeactivateRecordParams{
		Realm:   newSliceID.Realm,
		Userid:  pgtype.Text{String: newSliceID.User, Valid: true},
		Deactby: newSliceID.Setby,
		Deactat: newSliceID.To,
	})
	if err != nil {
		l.Info().Error(err).Log("error while adding record in deactivated table")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	// Step:4 - send response
	l.Debug0().LogActivity("exiting from UserDeactivate()", newSliceID)
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(nil))
	// uncomment below while running test cases
	// wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(newSliceID))
}
