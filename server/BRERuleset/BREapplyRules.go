package breruleset

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	crux "github.com/remiges-tech/crux/matching-engine"
	"github.com/remiges-tech/crux/server"
)

type applyRules struct {
	Entity  crux.Entity `json:"entity" validate:"required"`
	Ruleset string      `json:"ruleset" validate:"required"`
}

func BREapplyRules(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("starting execution of BREapplyRules()")

	cruxCache, ok := s.Dependencies["cruxCache"].(*crux.Cache)
	if !ok {
		l.Debug0().Debug1().Log("Error while getting cruxCache instance from service Dependencies")
		// return wfinstance.WFInstanceNewResponse{}, fmt.Errorf("error while getting cruxCache instance from service Dependencies")
	}

	var req applyRules

	err := wscutils.BindJSON(c, &req)
	if err != nil {
		l.Error(err).Debug0().Log("Error Unmarshalling Query parameters to struct:")
		return
	}

	req.Entity.Realm = realmName
	// Validate request
	validationErrors := wscutils.WscValidate(req, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		l.Debug0().LogDebug("validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	schema, ruleset, err := cruxCache.RetriveRuleSchemasAndRuleSetsFromCache(brwf, req.Entity.App, realmName, req.Entity.Class, req.Ruleset, req.Entity.Slice)
	if err != nil {
		l.Debug0().Error(err).Log("error while Retrieve RuleSchemas and RuleSets FromCache")
		// return wfinstance.WFInstanceNewResponse{}, fmt.Errorf("error while Retrieve RuleSchemas and RuleSets FromCache: %v", err)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.Errcode_cache))
		return
	} else if schema == nil || ruleset == nil {
		l.Debug0().Error(err).Log("didn't find any data in RuleSchemas or RuleSets cache")
		// return wfinstance.WFInstanceNewResponse{}, fmt.Errorf("didn't find any data in RuleSchemas or RuleSets cache: ")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.Errcode_cache))
		return
	}
	actionSet := crux.ActionSet{}
	seenRuleSets := make(map[string]struct{})

	// Call the doMatch function passing the entity.entity, ruleset, and the empty actionSet and seenRuleSets
	actionset, _, err, _ := crux.DoMatch(req.Entity, ruleset, schema, actionSet, seenRuleSets, crux.Trace_t{})

	if err != nil {
		l.Debug0().Error(err).Log("error while performing DoMatch")
		// return wfinstance.WFInstanceNewResponse{}, err
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.Errcode_DoMatch_Failed))
		return
	}

	wscutils.SendSuccessResponse(c, wscutils.NewResponse(wscutils.SuccessStatus, actionset, nil))

}
