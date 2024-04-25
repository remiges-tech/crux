package breruleset

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	crux "github.com/remiges-tech/crux/matching-engine"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

type BRERuleSetActivateReq struct {
	Slice int32  `json:"slice" validate:"required,gt=0"`
	App   string `json:"app" validate:"required,alpha,lt=15"`
	Class string `json:"class" validate:"required,alpha,lt=15"`
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

	realmName := "Ecommerce"
	userID := "Raj"
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

	isCapable = true

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
	cruxCache, ok := s.Dependencies["cruxCache"].(*crux.Cache)
	if !ok {
		lh.Debug0().Debug1().Log("Error while getting cruxCache instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	ntraversed, nactivated, err := cruxCache.RetrieveAndCheckIsActiveRuleSet(B, applc, realmName, request.Class, request.Name, request.Slice)
	if err != nil {
		lh.Debug0().Error(err).Log("error while retriving and verifying whether all rulesetmust be active ")
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
func hasRulesetCapability(app string, query *sqlc.Queries, c *gin.Context) (bool, error) {
	count, err := query.GetRuleSetCapabilityForApp(c, sqlc.GetRuleSetCapabilityForAppParams{
		Userid: userID,
		Realm:  realmName,
		App:    pgtype.Text{String: app, Valid: true},
		Cap:    RULESET,
	})
	if err != nil {
		return false, err
	}
fmt.Println(">>>>>>>>>>>>>>.count ",count ," userid",userID," App: ",app," cap: ",RULESET," realm",realmName)
	if count > 0 {
		return true, nil
	} else {
		return false, fmt.Errorf("caller does not have ruleset capability for the specified app : %v", app)
	}
}

//================================================================================================================================================
// func retrieveAndCheckIsActive(cruxCache *crux.Cache, request BRERuleSetActivateReq) {
// 	applc := strings.ToLower(request.App)
// 	currentRuleset, res := cruxCache.GetRulesetsFromCacheWithName(B, applc, realmName, request.Class, request.Name, request.Slice)
// 	actionBlocks := currentRuleset.Rules

// 	for _, rule := range actionBlocks {
// 		if rule.RuleActions.ThenCall != "" {
// 			thenRuleset, res := cruxCache.GetRulesetsFromCacheWithName(B, applc, realmName, request.Class, rule.RuleActions.ThenCall, request.Slice)
// 			// for _, rule := range thenRuleset.Rules {
// 			// 	if rule.IsActive == true
// 		}
// 	}
//==========================================================================================================================================================================

// 	currentRulesetName = ruleset in the request

//	retrieveAndCheckIsActive(rulesetName string) {
//	    currentRuleset = retrieve from cache (currentRulesetName)
//	    actionBlocks = currentRuleset.RuleActions
//	    for each action in actionBlocks
//	        if action.thencall != nil
//	            thenRuleset = retrieve the target ruleset
//	            if thenRuleset.active == false
//	                return error
//	            endif
//	            retrieveAndCheck(thenRuleset)
//	        endif
//	    endfo
//	}
//==================================================================================================================================================================

//}
