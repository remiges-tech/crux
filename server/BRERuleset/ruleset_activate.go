package breruleset

import (
	"encoding/json"
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
	"github.com/remiges-tech/logharbour/logharbour"
)

type BRERuleSetActivateReq struct {
	Slice int32  `json:"slice" validate:"required,gt=0"`
	App   string `json:"app" validate:"required,alpha,lt=15"`
	Class string `json:"class" validate:"required,lt=15"`
	Name  string `json:"name" validate:"required,lt=20"`
}

type BRERuleSetActivateRes struct {
	Ntraversed int32 `json:"ntraversed"`
	Nactivated int32 `json:"nactivated"`
}

func BRERuleSetActivate(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("BRERuleSetActivate request received")

	var (
		request BRERuleSetActivateReq
	)

	// implement the user realm and all here
	var capForList = []string{"ruleset"}
	realmName = "Nova"
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

	// verifying whether user has ruleset capability for app
	applc := strings.ToLower(request.App)
	isValidUser, err := hasRulesetCapability(applc, query, c)
	if err != nil {
		lh.Error(err).Log("error while verifying whether caller has ruleset cap for app")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, err.Error()))
		return
	}
	if !isValidUser {
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_app_does_not_have_ruleset_cap))
		return
	}

	//  activate  rulesets in db
	ntraversed, nactivated, err := activateAllRulesets(query, lh, request.Slice, applc, request.Class, request.Name, c)
	if err != nil {
		lh.Debug0().Error(err).Log("error while retriving and verifying whether all ruleset must be active ")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, err.Error()))
		return
	}

	response := BRERuleSetActivateRes{
		Ntraversed: int32(ntraversed),
		Nactivated: int32(nactivated),
	}
	lh.Debug0().Log(" finished execution of BRERuleSetActivate()")
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(response))
}

func activateAllRulesets(query *sqlc.Queries, lh *logharbour.Logger, slice int32, app, class, rulesetName string, c *gin.Context) (ntraversed, nactivated int, err error) {

	var (
		rules           []crux.Rule_t
		thenRulesetName string
		thenRulesetArr  []string
	)

	// First check if ruleset exist in db
	count, err := query.GetBRERuleSetCount(c, sqlc.GetBRERuleSetCountParams{
		Realm:   realmName,
		Slice:   slice,
		App:     app,
		Class:   class,
		Setname: rulesetName,
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
	// To get rulesetdata and active status of ruleset
	rulesetData, err := query.GetBRERuleSetActiveStatus(c, sqlc.GetBRERuleSetActiveStatusParams{
		Realm:   realmName,
		Slice:   slice,
		App:     app,
		Class:   class,
		Setname: rulesetName,
		Brwf:    sqlc.BrwfEnumB,
	})
	if err != nil {
		lh.Error(err).Log("error while getting a ruleset data and it's active state")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return ntraversed, nactivated, err
	}

	// Unmarshalling into []Rule_t
	err = json.Unmarshal(rulesetData.Ruleset, &rules)
	if err != nil {
		lh.Error(err).Log("error while unmarshalling ruleset into crux.Rule_t")
		//wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Invalid, server.ErrCode_InvalidJson))
		return ntraversed, nactivated, err
	}

	if !rulesetData.IsActive.Bool {
		//  activate given ruleset in db
		err = query.ActivateBRERuleSet(c, sqlc.ActivateBRERuleSetParams{
			Realm:   realmName,
			Slice:   slice,
			App:     app,
			Class:   class,
			Setname: rulesetName,
			Brwf:    sqlc.BrwfEnumB,
		})
		if err != nil {
			lh.Error(err).Log("error while activating a ruleset ")
			errmsg := db.HandleDatabaseError(err)
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
			return
		}
		nactivated++

	}

	for _, rule := range rules {
		if rule.RuleActions.ThenCall != "" {
			ntraversed++
			thenRulesetName = rule.RuleActions.ThenCall
			thenRulesetArr = append(thenRulesetArr, thenRulesetName)
		}
	}

	// Recursive call
	if len(thenRulesetArr) > 0 {
		for _, thenRulesetName := range thenRulesetArr {
			ntrav, nact, err := activateAllRulesets(query, lh, slice, app, class, thenRulesetName, c)
			if err != nil {
				return 0, 0, err
			}
			ntraversed += ntrav
			nactivated += nact
		}
	}

	return ntraversed, nactivated, nil
}