package wfinstanceserv

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
)

// Incoming request format
type WFInstanceNewRequest struct {
	Slice    *int32            `json:"slice" validate:"required"`
	App      *string           `json:"app" validate:"required,alpha"`
	EntityID *string           `json:"entityid" validate:"required"`
	Entity   map[string]string `json:"entity" validate:"required"`
	Workflow *string           `json:"workflow" validate:"required"`
	Trace    int               `json:"trace,omitempty"`
	Parent   int               `json:"parent,omitempty"`
}

// GetWFinstanceNew will be responsible for processing the /wfinstanceNew request that comes through as a POST
func GetWFinstanceNew(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("GetWFinstanceNew request received")

	// Bind request
	var wfinstanceNewreq WFInstanceNewRequest
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
		// lh.Debug0().LogActivity("Invalid request:", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, errStr))
		return

	} else {
		// Additional attributes to append
		existingEntity["step"] = "START"
		existingEntity["stepfailed"] = "false"
	}

	// call doMatch()
	var actionSet ActionSet
	var ruleSet RuleSet
	var entity = getEntity(wfinstanceNewreq.Entity)
	var seenRuleSets map[string]bool
	actionSet, result, err := doMatch(entity, ruleSet, actionSet, seenRuleSets)
	if err != nil {
		lh.LogActivity("error while calling doMatch Method :", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(INVALID_PATTERN))
		return
	}
	if !result {
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(INVALID_PATTERN))
		return
	}

	// call insertRecord()
	errArray := insertRecord(actionSet, wfinstanceNewreq, s, c)
	if len(errArray) > 0 {
		lh.LogActivity("error while inserting data in wfinstance table:", errArray)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, errArray))
		return
	}

	lh.Log(fmt.Sprintf("Response : %v", map[string]any{"response": existingEntity}))
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(existingEntity))
}

func insertRecord(a ActionSet, rq WFInstanceNewRequest, s *service.Service, c *gin.Context) []wscutils.ErrorMessage {
	var errors []wscutils.ErrorMessage

	lh := s.LogHarbour

	lh.Log("Inside insertRecord()")

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		errors := append(errors, wscutils.BuildErrorMessage(wscutils.ErrcodeDatabaseError, nil))
		return errors
	}

	// The properties object of actionset must always contain either an attribute called done or an
	// attribute called nextstep.

	var propertyAttr []string
	var temp []string
	var done, nextstep string
	attributes := a.properties
	for attr, _ := range attributes {
		if attr == "done" {
			done = attr
			propertyAttr = append(propertyAttr, attr)
		} else if attr == "nextstep" {
			nextstep = attr
			propertyAttr = append(propertyAttr, attr)
		} else {
			temp = append(temp, attr)
		}
	}
	if len(propertyAttr) == 0 {
		errors = append(errors, wscutils.BuildErrorMessage(INVALID_PROPERTY_ATTRIBUTES, &ACTIONSET_PROPERTIES, temp...))
		return errors
	}
	fmt.Println(" nextstep ", nextstep)

	// If there is only one name in tasks and done attribute is not there, only one record will be
	// entered into the table
	if done == "" && len(a.tasks) == 1 {

		record, err := query.AddWFNewInstace(c, sqlc.AddWFNewInstaceParams{
			Entityid: *rq.EntityID,
			Slice:    *rq.Slice,
			App:      *rq.App,
			Class:    rq.Entity["class"],
			Workflow: *rq.Workflow,
			Step:     a.tasks[0],
			Nextstep: a.tasks[0],
		})
		if err != nil {
			lh.LogActivity("error while inserting a record in wfinstance table :", err.Error())
			errors = append(errors, wscutils.BuildErrorMessage(INSERT_OPERATION_FAILED, &TASK, a.tasks[0]))

		}
		fmt.Println(" record inserted  ", record)

	}

	return errors
}
