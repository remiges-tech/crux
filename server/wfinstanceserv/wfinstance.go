package wfinstanceserv

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/types"
)

// Incoming request format
type WFinstanceNewRequest struct {
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
	var wfinstanceNewreq WFinstanceNewRequest
	err := wscutils.BindJSON(c, &wfinstanceNewreq)
	if err != nil {
		lh.Debug0().LogActivity("error while binding json request error:", err.Error)
		return
	}

	// Standard validation of Incoming Request
	valError := wscutils.WscValidate(wfinstanceNewreq, getValsForGetWFinstanceNewReqError)
	if len(valError) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, valError))
		lh.Debug0().LogActivity("validation error:", valError)
		return
	}

	// 1.To verify whether app,slice ,class present in schema and get patternschema against it
	class := wfinstanceNewreq.Entity["class"]
	patternSchema, err := s.Database.(*sqlc.Queries).WfPatternSchemaGet(c, sqlc.WfPatternSchemaGetParams{
		App:   *wfinstanceNewreq.App,
		Slice: *wfinstanceNewreq.Slice,
		Class: class,
	})
	if err != nil {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(RECORD_NOT_EXIST, nil)}))
		lh.Debug0().LogActivity("failed to get data from DB:", err.Error)
		return
	}
	// 2.validate name and type of entity to patternschema

	lh.Log(fmt.Sprintf("Record found: %v", map[string]any{"response": string(patternSchema)}))
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(string(patternSchema)))

}

func getValsForGetWFinstanceNewReqError(err validator.FieldError) []string {
	return types.CommonValidation(err)
}
