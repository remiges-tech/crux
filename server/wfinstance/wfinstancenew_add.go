package wfinstance

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/logharbour/logharbour"
)

// Task request
type AddTaskRequest struct {
	Steps    []string
	Nextstep string
	Request  WFInstanceNewRequest
}

// getResponse request
type ResponseRequest struct {
	Subflow      map[string]string
	NextStep     string
	ResponseData []sqlc.Wfinstance
	Service      *service.Service
}

// To add multiple records in wfinstance table
func AddTasks(req AddTaskRequest, s *service.Service, c *gin.Context) (WFInstanceNewResponse, error) {
	var response WFInstanceNewResponse
	var parent pgtype.Int4
	subflow := make(map[string]string)

	lh := s.LogHarbour.WithClass("wfinstance")
	lh.Debug0().Log("inside addTasks()")

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("WFInstanceNew||addTasks()||error while getting query instance from service dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return WFInstanceNewResponse{}, errors.New(INVALID_DATABASE_DEPENDENCY)
	}

	if req.Request.Parent != nil {
		// convert int32 tp pgtype.Int4
		parent = ConvertToPGType(*req.Request.Parent)
	}

	// Add record in wfinstance table
	responseData, error := query.AddWFNewInstances(c, sqlc.AddWFNewInstancesParams{
		Entityid: req.Request.EntityID,
		Slice:    req.Request.Slice,
		App:      req.Request.App,
		Class:    req.Request.Entity[CLASS],
		Workflow: req.Request.Workflow,
		Step:     req.Steps,
		Nextstep: req.Nextstep,
		Parent:   parent,
	})
	if error != nil {
		lh.Error(error).Log("WFInstanceNew||addTasks()||error while adding Task steps in wfinstance table")
		return response, error
	}

	// data change log
	for _, val := range responseData {
		dclog := lh.WithClass("wfinstance").WithInstanceId(strconv.Itoa(int(val.ID)))
		dclog.LogDataChange("created wfinstance ", logharbour.ChangeInfo{
			Entity: "wfinstance",
			Op:     "create",
			Changes: []logharbour.ChangeDetail{
				{
					Field:  "row",
					OldVal: nil,
					NewVal: val},
			},
		})
	}

	// To get workflow if step is present in stepworkflow table
	lh.Debug0().Log("WFInstanceNew||addTasks()||verifying whether step is workflow if it is, then append it to subflow")
	for _, step := range req.Steps {
		workflow, _ := query.GetWorkflowNameForStep(c, step)

		if !reflect.ValueOf(workflow).IsZero() {
			subflow[workflow.Step] = workflow.Workflow
		}

	}

	// To get response
	responseRequest := ResponseRequest{
		Subflow:      subflow,
		NextStep:     req.Nextstep,
		ResponseData: responseData,
		Service:      s,
	}

	response = getResponse(responseRequest)
	lh.Debug0().LogActivity("response of addTask() :", response)
	return response, nil
}

// response structure
func getResponse(r ResponseRequest) WFInstanceNewResponse {
	var tasks []map[string]int32
	var loggedDate pgtype.Timestamp
	var response WFInstanceNewResponse

	lh := r.Service.LogHarbour.WithClass("wfinstance")
	lh.Debug0().Log("Inside getResponse()")

	for _, val := range r.ResponseData {
		// adding tasks
		task := make(map[string]int32)
		task[val.Step] = val.ID
		tasks = append(tasks, task)
		//loggingdates
		loggedDate = val.Loggedat
	}
	//var response WFInstanceNewResponse
	if len(tasks) > 1 {
		// response for multiple task steps
		response = WFInstanceNewResponse{
			Tasks:    tasks,
			Nextstep: &r.NextStep,
			Loggedat: loggedDate,
			Subflows: &r.Subflow,
		}
	} else {
		//response for single task step
		response = WFInstanceNewResponse{
			Tasks:    tasks,
			Loggedat: loggedDate,
			Subflows: &r.Subflow,
		}
	}
	return response
}

// To convert int to pgtype.Int4
func ConvertToPGType(value int32) pgtype.Int4 {
	if value != 0 {
		return pgtype.Int4{
			Int32: value,
			Valid: true,
		}
	}
	return pgtype.Int4{Int32: value, Valid: false}
}
