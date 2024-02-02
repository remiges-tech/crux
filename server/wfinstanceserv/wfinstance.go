package wfinstanceserv

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
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
	lh := s.LogHarbour.WithWhatClass("wfinstance").WithWhatInstanceId("GetWFinstanceNew(")
	lh.Log("GetWFinstanceNew request received")

	// Bind request
	var wfinstanceNewreq WFInstanceNewRequest
	err := wscutils.BindJSON(c, &wfinstanceNewreq)
	if err != nil {
		lh.Debug0().LogActivity("error while binding json request error:", err.Error())
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

	lh.Log(fmt.Sprintf("Response : %v", map[string]any{"response": existingEntity}))
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(existingEntity))
}
