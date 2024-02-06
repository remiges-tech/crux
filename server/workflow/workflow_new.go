package workflow

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
)



func WorkFlowNew(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Log("Starting execution of WorkFlowNew()")
	
	var wf workflowNew

	err := wscutils.BindJSON(c, &wf)
	if err != nil {
		l.LogActivity("Error Unmarshalling Query parameters to struct:", err.Error())
		return
	}

	// Validate request
	validationErrors := wscutils.WscValidate(wf, func(err validator.FieldError) []string { return []string{} })
	customValidationErrors := customValidationErrors(wf)
	validationErrors = append(validationErrors, customValidationErrors...)
	if len(validationErrors) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	ruleset, err := json.Marshal(wf.Flowrules)
	if err != nil {
		patternSchema := "flowrules"
		l.LogDebug("Error while marshaling Flowrules", err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(wscutils.ErrcodeInvalidJson, &patternSchema)}))
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}

	_, err = query.WorkFlowNew(c, sqlc.WorkFlowNewParams{
		Realm:      realmID,
		Slice:      wf.Slice,
		App:        wf.App,
		Brwf:       brwf,
		Class:      wf.Class,
		Setname:    setBy,
		IsActive:   pgtype.Bool{Bool: isActive},
		IsInternal: wf.IsInternal,
		Ruleset:    ruleset,
		Createdby: setBy,
	})
	if err != nil {
		l.LogActivity("Error while creating schema", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: "Created successfully", Messages: nil})
	l.Log("Finished execution of SchemaNew()")
}

func customValidationErrors(wf workflowNew) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage

	return validationErrors
}
