package schema

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
)

func SchemaNew(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Log("Starting execution of SchemaNew()")
	createdby := "admin"

	var sh Schema

	// check the capgrant table to see if the calling user has the capability to perform the
	// operation
	// isCapable, _ := utils.Authz_check(types.OpReq{
	// 	User:      username,
	// 	CapNeeded: []string{"schema"},
	// }, false)

	// if !isCapable {
	// 	l.Log("Unauthorized user:")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(utils.ErrUnauthorized))
	// 	return
	// }

	// The system will check whether there are any rulesets in the ruleset table whose
	// (slice,app,class) match this record. If this is true, then the call will fail.
	// In other words, updating a schema is not allowed once rulesets referring to it are defined.

	err := wscutils.BindJSON(c, &sh)
	if err != nil {
		l.LogActivity("Error Unmarshalling Query paramaeters to struct:", err.Error())
		return
	}

	// Validate request
	validationErrors := wscutils.WscValidate(sh, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	query, ok := s.Database.(*sqlc.Queries)
	if !ok {
		l.Log("Error while getting query instance from service")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}
	paternSchema, err := json.Marshal(sh.Patternschema)
	if err != nil {
		l.LogDebug("Error while marshaling paternSchema", err)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse("MarshalJSON"))
		return
	}

	actionschema, err := json.Marshal(sh.Actionschema)
	if err != nil {
		l.LogDebug("Error while marshaling actionschema", err)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse("MarshalJSON"))
		return
	}
	_, err = query.SchemaNew(c, sqlc.SchemaNewParams{Realm: 1, Slice: sh.Slice, Class: sh.Class, App: sh.App, Brwf: "W", Patternschema: paternSchema, Actionschema: actionschema, Createdby: createdby, Editedby: createdby})
	if err != nil {
		l.LogActivity("Error while creating schema", err)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: "Created successfully", Messages: nil})
	l.Log("Finished execution of SchemaNew()")
}
