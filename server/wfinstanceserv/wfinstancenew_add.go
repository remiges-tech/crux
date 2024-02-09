package wfinstanceserv

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
)

type addRecordRequest struct {
	Step     string
	Nextstep string
	Request  WFInstanceNewRequest
}

// To add single record in wfinstance table
func addSingleTask(r addRecordRequest, s *service.Service, c *gin.Context) (WFInstanceNewResponse, error) {
	lh := s.LogHarbour
	var response WFInstanceNewResponse
	var step = r.Step

	lh.Log("Inside addSingleTask()")

	// Add record in wfinstance table
	id, workflow, loggDate, error := addRecord(r, s, c)
	if error != nil {
		lh.LogActivity("error while adding single steps in wfinstance table :", error.Error())
		return response, error
	}

	// To form task array for respose with step name and id
	var tasks []map[string]int32
	task := map[string]int32{
		step: id,
	}
	tasks = append(tasks, task)

	//To form subflow for respose if present step is another workflow
	subflow := make(map[string]string)
	if workflow != "" {
		subflow[step] = workflow
	}

	// response
	response = WFInstanceNewResponse{
		Tasks:    tasks,
		Loggedat: loggDate,
		Subflows: &subflow,
	}

	return response, nil
}

// To add multiple records in wfinstance table
func addMultipleTasks(r []addRecordRequest, s *service.Service, c *gin.Context) (WFInstanceNewResponse, error) {
	var response WFInstanceNewResponse
	var tasks []map[string]int32
	subflow := make(map[string]string)
	var loggdates []pgtype.Timestamp
	var nextStep string
	lh := s.LogHarbour
	lh.Log("Inside addMultipleTasks()")

	// Add record in wfinstance table
	for _, req := range r {
		id, workflow, loggDate, error := addRecord(req, s, c)
		if error != nil {
			lh.LogActivity("error while adding multiple steps in wfinstance table :", error.Error())
			return response, error
		}
		// adding tasks
		tasks = append(tasks, map[string]int32{req.Step: id})

		// workflow if step require another workflow
		if workflow != "" {
			subflow[req.Step] = workflow
		}
		nextStep = req.Nextstep

		//loggingdates
		loggdates = append(loggdates, loggDate)

	}
	// response
	response = WFInstanceNewResponse{
		Tasks:    tasks,
		Nextstep: nextStep,
		Loggedat: loggdates[len(loggdates)-1],
		Subflows: &subflow,
	}

	return response, nil
}

// To add new instance in wfinstance table
func addRecord(r addRecordRequest, s *service.Service, c *gin.Context) (int32, string, pgtype.Timestamp, error) {
	//var errors []wscutils.ErrorMessage
	lh := s.LogHarbour
	lh.Log("Inside addRecord()")
	var loggDate pgtype.Timestamp

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		//errors := wscutils.BuildErrorMessage(wscutils.ErrcodeDatabaseError, nil)
		return 0, "", loggDate, errors.New(INVALID_DATABASE_DEPENDENCY)
	}

	parent := ConvertToPGType(r.Request.Parent)

	recordID, err := query.AddWFNewInstace(c, sqlc.AddWFNewInstaceParams{
		Entityid: *r.Request.EntityID,
		Slice:    *r.Request.Slice,
		App:      *r.Request.App,
		Class:    r.Request.Entity["class"],
		Workflow: *r.Request.Workflow,
		Step:     r.Step,
		Nextstep: r.Nextstep,
		Parent:   parent,
	})
	if err != nil {
		lh.LogActivity("error while inserting a record in wfinstance table :", err.Error())
		return 0, "", loggDate, errors.New(INSERT_OPERATION_FAILED)

	}

	// To get workflow if step is present in stepworkflow table
	workflow, _ := query.GetWorkflow(c, r.Step)

	//To get Loggedat from  wfinstance table
	loggDate, err = query.GetLoggedate(c, recordID)
	if err != nil {
		lh.LogActivity("error while getting loggedat from wfinstance table :", err.Error())
		return 0, "", loggDate, fmt.Errorf("error while fetching loggedat from wfinstance table")

	}

	return recordID, workflow, loggDate, nil

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
