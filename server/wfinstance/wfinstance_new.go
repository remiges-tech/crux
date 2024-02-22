package wfinstance

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

// Incoming request format
type WFInstanceNewRequest struct {
	Slice    int32             `json:"slice" validate:"required"`
	App      string            `json:"app" validate:"required,alpha"`
	EntityID string            `json:"entityid" validate:"required"`
	Entity   map[string]string `json:"entity" validate:"required"`
	Workflow string            `json:"workflow" validate:"required"`
	Trace    *int              `json:"trace,omitempty"`
	Parent   *int32            `json:"parent,omitempty"`
}

// WFInstanceNew response format
type WFInstanceNewResponse struct {
	Tasks     []map[string]int32 `json:"tasks"  validate:"required"`
	Nextstep  string             `json:"nextstep"`
	Loggedat  pgtype.Timestamp   `json:"loggedat"`
	Subflows  *map[string]string `json:"subflows"`
	Tracedata *map[string]string `json:"tracedata"`
}

// GetWFinstanceNew will be responsible for processing the /wfinstanceNew request that comes through as a POST
func GetWFinstanceNew(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour.WithWhatClass("wfinstance")
	lh.Log("GetWFinstanceNew request received")
	var wfinstanceNewreq WFInstanceNewRequest
	var actionSet ActionSet
	var ruleSet RuleSet
	var entity = getEntity(wfinstanceNewreq.Entity)
	var seenRuleSets map[string]bool
	var response WFInstanceNewResponse
	var attribute map[string]string
	var done, nextStep string
	var steps []string

	isCapable, _ := types.Authz_check(types.OpReq{
		User: USERID,
	}, false)

	if !isCapable {
		lh.Info().LogActivity("Unauthorized user:", USERID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	// Bind request
	err := wscutils.BindJSON(c, &wfinstanceNewreq)
	if err != nil {
		lh.Debug0().LogActivity("error while binding json request error:", err)
		return
	}
	// Standard validation of Incoming Request
	validationErrors := wscutils.WscValidate(wfinstanceNewreq, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		lh.Debug0().LogActivity("validation error:", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}
	// Validate request
	existingEntity := wfinstanceNewreq.Entity
	isValidReq, errStr := validateWFInstanceNewReq(wfinstanceNewreq, s, c)
	if len(errStr) > 0 || !isValidReq {
		lh.Debug0().LogActivity("request validation error:", errStr)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, errStr))
		return

	} else {
		// Additional attributes to append
		existingEntity[STEP] = START
		existingEntity[STEPFALED] = FALSE
	}
	lh.Debug0().LogActivity("wfinstanceNewRequest after adding additional attributes :", wfinstanceNewreq)

	// call doMatch()
	actionSet, _, err = doMatch(entity, ruleSet, actionSet, seenRuleSets)
	if err != nil {
		lh.LogActivity("error while calling doMatch Method :", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Invalid, server.ErrCode_Invalid))
		return
	}

	// To verify actionSet Properties and get their values
	attribute, error := getValidPropertyAttr(actionSet)
	if error != nil {
		lh.Debug0().LogActivity("error while verifying actionset properties :", error.Error())
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
		lh.Log(fmt.Sprintf("Response : %v", map[string]any{"response": response}))
		wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(response))
	}

	// To add records in table
	steps = actionSet.Tasks

	// if tasks of actionset contains only one task
	if len(actionSet.Tasks) == 1 && done == "" {
		addTaskRequest := AddTaskRequest{
			Steps:    steps,
			Nextstep: steps[0],
			Request:  wfinstanceNewreq,
		}
		response, err = addTasks(addTaskRequest, s, c)
		if err != nil {
			lh.LogActivity("error while adding single step in wfinstance table :", err.Error())
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
			return
		}
		wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(response))

	}
	// if tasks of actionset contains multiple tasks
	if len(actionSet.Tasks) > 1 && done == "" {
		addTaskRequest := AddTaskRequest{
			Steps:    steps,
			Nextstep: nextStep,
			Request:  wfinstanceNewreq,
		}
		response, err = addTasks(addTaskRequest, s, c)
		if err != nil {
			lh.LogActivity("error while adding multiple steps in wfinstance table :", error.Error())
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
			return
		}
		wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(response))
	}

}

// To verify whether actionset.properties attributes valid and get their values
func getValidPropertyAttr(a ActionSet) (map[string]string, error) {
	attribute := make(map[string]string)
	attributes := a.Properties
	for attr, val := range attributes {
		if attr == DONE {
			attribute[attr] = val
		} else if attr == NEXTSTEP {
			attribute[attr] = val
		} else {
			return nil, fmt.Errorf("property attributes does not contain either done or nextstep")
		}

	}
	return attribute, nil
}
