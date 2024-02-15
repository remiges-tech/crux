package wfinstance

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
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
	ResponseData []sqlc.AddWFNewInstancesRow
	Service      *service.Service
}

// To add multiple records in wfinstance table
func addTasks(req AddTaskRequest, s *service.Service, c *gin.Context) (WFInstanceNewResponse, error) {
	var response WFInstanceNewResponse
	var parent pgtype.Int4
	subflow := make(map[string]string)

	lh := s.LogHarbour
	lh.Log("Inside addTasks()")

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		//errors := wscutils.BuildErrorMessage(wscutils.ErrcodeDatabaseError, nil)
		return WFInstanceNewResponse{}, errors.New(INVALID_DATABASE_DEPENDENCY)
	}

	// convert int32 tp pgtype.Int4
	parent = ConvertToPGType(*req.Request.Parent)

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
		lh.LogActivity("error while adding Task steps in wfinstance table :", error.Error())
		return response, error
	}

	// To get workflow if step is present in stepworkflow table
	for _, step := range req.Steps {
		workflowData, _ := query.GetWorkflow(c, step)

		if len(workflowData) > 0 && workflowData[0].Workflow != "" {
			subflow[workflowData[0].Step] = workflowData[0].Workflow
		}

	}

	// to get response
	responseRequest := ResponseRequest{
		Subflow:      subflow,
		NextStep:     req.Nextstep,
		ResponseData: responseData,
		Service:      s,
	}

	response = getResponse(responseRequest)

	return response, nil
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

// response structure
func getResponse(r ResponseRequest) WFInstanceNewResponse {
	var tasks []map[string]int32
	var loggdates []pgtype.Timestamp
	var response WFInstanceNewResponse

	lh := r.Service.LogHarbour
	lh.Log("Inside addTasks()")

	for _, val := range r.ResponseData {
		// adding tasks
		task := make(map[string]int32)
		task[val.Step] = val.ID
		tasks = append(tasks, task)
		//loggingdates
		loggdates = append(loggdates, val.Loggedat)

	}

	//var response WFInstanceNewResponse
	if len(tasks) > 1 {
		// response for multiple task steps
		response = WFInstanceNewResponse{
			Tasks:    tasks,
			Nextstep: r.NextStep,
			Loggedat: loggdates[len(loggdates)-1],
			Subflows: &r.Subflow,
		}
	} else {
		//response for single task step
		response = WFInstanceNewResponse{
			Tasks:    tasks,
			Loggedat: loggdates[len(loggdates)-1],
			Subflows: &r.Subflow,
		}
	}
	return response
}
