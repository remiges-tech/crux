package wfinstance

import (
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

// Incoming request format
type WFInstanceTryRequest struct {
	Slice      int32             `json:"slice" validate:"required,gt=0"`
	App        string            `json:"app" validate:"required,lt=50"`
	EntityID   string            `json:"entityid" validate:"required,gt=0,lt=40"`
	Entity     map[string]string `json:"entity" validate:"required"`
	Workflow   string            `json:"workflow" validate:"required,gt=0,lt=50"`
	Step       string            `json:"step" validate:"required,gt=0,lt=30"`
	StepFailed bool              `json:"stepfailed" validate:"required"`
	Trace      *int              `json:"trace,omitempty"`
}

type WFInstanceTryResponse struct {
	Tasks     []string          `json:"tasks,omitempty"`
	Nextstep  string            `json:"nextstep,omitempty"`
	Loggedat  pgtype.Timestamp  `json:"loggedat,omitempty"`
	Subflows  map[string]string `json:"subflows,omitempty"`
	Tracedata map[string]string `json:"tracedata,omitempty"`
}

// GetWFInstanceTry will be responsible for processing the /wfinstancetry request that comes through as a POST
func GetWFInstanceTry(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("GetWFInstanceTry request received")

	var (
		wfinstanceTryreq WFInstanceTryRequest
		actionSet        crux.ActionSet
		seenRuleSets     = make(map[string]struct{})
		attribute        map[string]string
		done, nextStep   string
		steps            []string
		subflow          = make(map[string]string)
		// ruleSet          *crux.Ruleset_t
	)
	// userID, err := server.ExtractUserNameFromJwt(c)
	// if err != nil {
	// 	lh.Info().Log("unable to extract userID from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	// realm, err := server.ExtractRealmFromJwt(c)
	// if err != nil {
	// 	lh.Info().Log("unable to extract realm from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	isCapable, _ := server.Authz_check(types.OpReq{
		User: userID,
	}, false)

	if !isCapable {
		lh.Info().LogActivity("unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	// Bind request
	err := wscutils.BindJSON(c, &wfinstanceTryreq)
	if err != nil {
		lh.Error(err).Log("GetWFInstanceTry||error while binding json request error")
		return
	}
	// Standard validation of Incoming Request
	validationErrors := wscutils.WscValidate(wfinstanceTryreq, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		lh.Debug0().LogActivity("GetWFInstanceTry||validation error:", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}
	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("WFInstanceNew||addTasks()||error while getting query instance from service dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	// Validate request
	existingEntity := wfinstanceTryreq.Entity

	wfInstancereq := WFInstanceNewRequest{
		Slice:    wfinstanceTryreq.Slice,
		App:      wfinstanceTryreq.App,
		EntityID: wfinstanceTryreq.EntityID,
		Entity:   existingEntity,
		Workflow: wfinstanceTryreq.Workflow,
		Trace:    wfinstanceTryreq.Trace,
	}
	isValidReq, errStr := validateWFInstanceNewReq(wfInstancereq, realm, wfinstanceTryreq.Step, s, c)
	if len(errStr) > 0 || !isValidReq {
		lh.Debug0().LogActivity("GetWFInstanceTry||request validation error:", errStr)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, errStr))
		return

	}

	//  doMatch() Processing
	// To get Entity
	entity := getEntityStructure(wfInstancereq, realm)

	cruxCache, ok := s.Dependencies["cruxCache"].(*crux.Cache)
	if !ok {
		lh.Debug0().Debug1().Log("Error while getting cruxCache instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	class := wfinstanceTryreq.Entity["class"]
	schema, ruleset, err := cruxCache.RetriveRuleSchemasAndRuleSetsFromCache(WFE, wfinstanceTryreq.App, realm, class, wfinstanceTryreq.Workflow, wfinstanceTryreq.Slice)
	if err != nil {
		lh.Debug0().Error(err).Log("error while Retrieve RuleSchemas and RuleSets FromCache")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Invalid, err.Error()))
		return
	} else if schema == nil || ruleset == nil {
		lh.Debug0().Log("didn't find any data in RuleSchemas or RuleSets cache")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NotFound, server.Errode_RuleSchemas_Or_Rulests_Not_Found))
		return
	}

	// call DoMatch()
	actionSet, _, err, _ = crux.DoMatch(entity, ruleset, schema, actionSet, seenRuleSets, crux.Trace_t{})
	if err != nil {
		lh.Error(err).Log("GetWFInstanceTry||error while calling doMatch Method")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Invalid, err.Error()))
		return
	}

	//To verify actionSet Properties and get their values
	attribute, error := getValidPropertyAttr(actionSet)
	if error != nil {
		lh.Debug0().LogActivity("GetWFInstanceTry||error while verifying actionset properties :", error.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Invalid, server.ErrCode_Invalid_property_attributes))
		return
	}
	if attribute[DONE] == TRUE {
		done = attribute[DONE]
	} else {
		nextStep = attribute[NEXTSTEP]
	}

	if done == TRUE {
		response := make(map[string]bool)
		response[DONE] = true
		lh.Log("Finished execution of GetWFInstanceTry")
		wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(response))
	}
	// To add records in table
	steps = actionSet.Tasks

	// To get workflow if step is present in stepworkflow table
	lh.Debug0().Log("WFInstanceTry||verifying whether step is workflow if it is, then append it to subflow")
	for _, step := range steps {
		workflow, err := query.GetWorkflowNameForStep(c, step)

		if err != nil {
			if err.Error() == "no rows in result set" {
				continue // If no workflow is found, continue to the next step
			}
		}

		// Only proceed if err is nils
		if err == nil {
			subflow[workflow.Step] = workflow.Workflow
		}
	}

	// if tasks of actionset contains only one task
	if len(actionSet.Tasks) == 1 && done == "" {
		response := WFInstanceTryResponse{
			Tasks:    []string{actionSet.Tasks[0]},
			Subflows: subflow,
			//Tracedata: map[string]string{},
		}

		lh.Log("Finished execution of GetWFInstanceTry")
		wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(response))
		return
	}
	// if tasks of actionset contains multiple tasks
	if len(actionSet.Tasks) > 1 && done == "" {
		response := WFInstanceTryResponse{
			Tasks:    steps,
			Subflows: subflow,
			Nextstep: nextStep,
			//Tracedata: map[string]string{},
		}
		lh.Log("Finished execution of GetWFInstanceTry")
		wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(response))
	}

}
