package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/types"
)

type WorkflowGetReq struct {
	Slice int32  `json:"slice" validate:"required,gt=0"`
	App   string `json:"app" validate:"required,alpha"`
	Class string `json:"class" validate:"required,alpha"`
	Name  string `json:"name" validate:"required,alpha"`
}

type WorkflowgetRow struct {
	ID         int32            `json:"id"`
	Slice      int32            `json:"slice"`
	App        string           `json:"app"`
	Class      string           `json:"class"`
	Name       string           `json:"name"`
	IsActive   bool             `json:"is_active"`
	IsInternal bool             `json:"is_internal"`
	Flowrules  interface{}      `json:"flowrules"`
	Createdat  pgtype.Timestamp `json:"createdat"`
	Createdby  string           `json:"createdby"`
	Editedat   pgtype.Timestamp `json:"editedat"`
	Editedby   pgtype.Text      `json:"editedby"`
}

// WorkflowGet will be responsible for processing the /workflowget request that comes through as a POST
func WorkflowGet(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("WorkflowGet request received")

	// var response schemaGetResp
	var request WorkflowGetReq
	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.LogActivity("error while binding json request error:", err.Error)
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
		lh.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}

	dbResponse, err := query.Workflowget(c, sqlc.WorkflowgetParams{
		Slice:   request.Slice,
		App:     request.App,
		Class:   request.Class,
		Setname: request.Name,
	})
	if err != nil {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(types.RECORD_NOT_EXIST, nil)}))
		lh.LogActivity("failed to get data from DB:", err.Error)
		return
	}

	tempData := responseBinding(dbResponse)

	err = json.Unmarshal(dbResponse.Flowrules, &tempData.Flowrules)
	if err != nil {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(types.OPERATION_FAILED, nil)}))
		lh.LogActivity("failed to unmarshal data:", err.Error)
		return
	}
	lh.Log(fmt.Sprintf("Record found: %v", map[string]any{"response": tempData}))
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(tempData))
}

func responseBinding(dbResponse sqlc.WorkflowgetRow) WorkflowgetRow {
	tempData := WorkflowgetRow{
		ID:         dbResponse.ID,
		Slice:      dbResponse.Slice,
		App:        dbResponse.App,
		Class:      dbResponse.Class,
		Name:       dbResponse.Name,
		IsActive:   dbResponse.IsActive.Bool,
		IsInternal: dbResponse.IsInternal,
		Createdat:  dbResponse.Createdat,
		Createdby:  dbResponse.Createdby,
		Editedat:   dbResponse.Editedat,
		Editedby:   dbResponse.Editedby,
	}
	return tempData
}
