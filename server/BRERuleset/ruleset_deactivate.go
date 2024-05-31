package breruleset

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	crux "github.com/remiges-tech/crux/matching-engine"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

type BRERuleSetDeActivateReq struct {
	Slice int32  `json:"slice" validate:"required,gt=0"`
	App   string `json:"app" validate:"required,alpha,lt=15"`
	Class string `json:"class" validate:"required,alpha,lt=15"`
	Name  string `json:"name" validate:"required,lt=20"`
}

func BRERuleSetDeActivate(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("BRERuleSetDeActivate request received")

	var (
		request BRERuleSetDeActivateReq
	)

	// implement the user realm and all here
	var capForList = []string{"ruleset"}

	// userID := "Raj"
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

	realmName, ok := s.Dependencies["realmName"].(string)
	if !ok {
		lh.Debug0().Log("error while getting realmName instance from service dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
		return
	}


	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: capForList,
	}, false)

	if !isCapable {
		lh.Info().LogActivity(server.ErrCode_Unauthorized, userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Debug0().Error(err).Log("error while binding json request error")
		return
	}

	valError := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
	if len(valError) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, valError))
		lh.LogActivity("validation error:", valError)
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	cruxCache, ok := s.Dependencies["cruxCache"].(*crux.Cache)
	if !ok {
		lh.Debug0().Debug1().Log("error while getting cruxCache instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	// verifying whether user has ruleset capability for app
	applc := strings.ToLower(request.App)
	isValidUser, err := hasRulesetCapability(applc, query, c, realmName)
	if err != nil {
		lh.Error(err).Log("error while verifying whether caller has ruleset cap for app")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, err.Error()))
		return
	}
	if !isValidUser {
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_app_does_not_have_ruleset_cap))
		return
	}

	// check if a given ruleset exist in db
	count, err := query.GetBRERuleSetCount(c, sqlc.GetBRERuleSetCountParams{
		Slice:   request.Slice,
		App:     applc,
		Class:   request.Class,
		Setname: request.Name,
		Realm:   realmName,
		Brwf:    sqlc.BrwfEnumB,
	})
	if err != nil {
		lh.Error(err).Log("error while getting a ruleset")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return

	}

	if count == 0 {
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NotFound, server.ErrCode_NotFound))
		return
	}

	// deactive a ruleset in db
	err = query.DeActivateBRERuleSet(c, sqlc.DeActivateBRERuleSetParams{
		Realm:   realmName,
		Slice:   request.Slice,
		App:     applc,
		Class:   request.Class,
		Setname: request.Name,
		Brwf:    sqlc.BrwfEnumB,
	})
	if err != nil {
		lh.Error(err).Log("error while deactivating ruleset")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return

	}

	// deactivating ruleset in cache
	cruxCache.Purge(B, applc, realmName, request.Class, "ruleset", request.Name, request.Slice)

	lh.Debug0().Log(" finished execution of BRERuleSetDeActivate()")
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: nil, Messages: nil})

}
