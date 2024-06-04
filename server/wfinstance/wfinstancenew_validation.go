package wfinstance

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	crux "github.com/remiges-tech/crux/matching-engine"
	"github.com/remiges-tech/crux/server"
)

func validateWFInstanceNewReq(r WFInstanceNewRequest, realm string, s *service.Service, c *gin.Context) (bool, []wscutils.ErrorMessage) {
	lh := s.LogHarbour.WithClass("wfinstance")
	entity := r.Entity
	var errRes []wscutils.ErrorMessage

	lh.Debug0().Log("Inside ValidateWFInstaceNewReq()")
	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("GetWFinstanceNew||validateWFInstanceNewReq()||error while getting query instance from service Dependencies")
		errRes := append(errRes, wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_DatabaseError, nil))
		return false, errRes
	}
	// Validate request
	isValidReq, errAry := validateWorkflow(r, s, c, realm)
	if len(errAry) > 0 || !isValidReq {
		lh.Debug0().LogActivity("GetWFinstanceNew||validateWFInstanceNewReq()||invalid request:", errAry)
		errRes = errAry
		return false, errRes
	}

	// To verify whether entity has valid structure
	_, isKeyExist := entity[CLASS]

	if !(entity != nil && isKeyExist) {
		lh.Debug0().Log("GetWFinstanceNew||validateWFInstanceNewReq()||entity does not match with standard structure")
		errRes = append(errRes, wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid_Entity, &ENTITY))
		return false, errRes
	}

	// To verify whether app,slice,class present in schema and get patternschema against it
	class := entity[CLASS]
	pattern, err := query.WfPatternSchemaGet(c, sqlc.WfPatternSchemaGetParams{
		Slice: r.Slice,
		Class: class,
		App:   r.App,
		Realm: realm,
	})
	if err != nil {
		lh.Error(err).Log("GetWFinstanceNew||validateWFInstanceNewReq()||failed to get schema pattern from DB")
		errRes = append(errRes, wscutils.BuildErrorMessage(server.MsgId_NoSchemaFound, server.ErrCode_NotFound, nil))
		return false, errRes
	}

	// Unmarshalling byte data to schemapatten struct
	schemaPattern, err := byteToPatternSchema(pattern)
	lh.Debug0().LogActivity("GetWFinstanceNew||validateWFInstanceNewReq()||patternschema :", schemaPattern)
	if err != nil {
		lh.Debug0().LogActivity("GetWFinstanceNew||validateWFInstanceNewReq()|| error while converting byte patternschema to struct:", err)
		errRes = append(errRes, wscutils.BuildErrorMessage(server.MsgId_NoSchemaFound, server.ErrCode_Invalid_pattern_schema, nil))
		return false, errRes
	}
	schema := crux.Schema_t{
		Class:         class,
		PatternSchema: *schemaPattern,
	}

	// Forming  requested entity  into proper Entity struct
	EntityStruct := getEntityStructure(r, realm)
	lh.Debug1().LogActivity("GetWFinstanceNew||validateWFInstanceNewReq()||entity stucture:", EntityStruct)

	//  To match entity against patternschema
	isValidEntity, err := ValidateEntity(EntityStruct, &schema, s)
	if !isValidEntity || err != nil {
		lh.Debug0().LogActivity(" GetWFinstanceNew||validateWFInstanceNewReq()||error while validating entity against patternschema:", err)
		errRes = append(errRes, wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid_Entity, &ENTITY))
		return false, errRes

	}
	return true, nil
}

// validate workflow
func validateWorkflow(r WFInstanceNewRequest, s *service.Service, c *gin.Context, realm string) (bool, []wscutils.ErrorMessage) {
	var errors []wscutils.ErrorMessage
	lh := s.LogHarbour.WithClass("wfinstance")
	entityClass := r.Entity[CLASS]

	lh.Debug0().Log("Inside validateWorkflow()")
	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Debug0().Log("GetWFinstanceNew||validateWorkflow()||error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		errors := append(errors, wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_DatabaseError, nil))
		return false, errors
	}

	// The value of app specified in the request matches the app ID with which this workflow is associated

	lh.Debug0().Log("GetWFinstanceNew||validateWorkflow()||verifying whether app present in request is valid")
	applc := strings.ToLower(r.App)

	app, err := query.GetApp(c, sqlc.GetAppParams{
		Slice: r.Slice,
		App:   applc,
		Class: entityClass,
		Realm: realm,
	})

	if err != nil {
		lh.Error(err).Log("GetWFinstanceNew||validateWorkflow()||failed to get app from ruleset table")
		errors = append(errors, wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid_APP, &APP, r.App))
	}

	// The class of the workflow must match that of entity
	lh.Debug0().Log("GetWFinstanceNew||validateWorkflow()||verifying whether class present in request is valid")
	class, err := query.GetClass(c, sqlc.GetClassParams{
		Slice: r.Slice,
		App:   app,
		Class: entityClass,
		Realm: realm,
	})

	if err != nil {
		lh.Error(err).Log("GetWFinstanceNew||validateWorkflow()||failed to get class from ruleset table")
		errors = append(errors, wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid_Class, &CLASS, entityClass))
	}

	// The workflow named has is_active == true and internal == false

	// To get worflow active status
	lh.Debug0().Log("GetWFinstanceNew||validateWorkflow()||verifying whether workflow active status is valid")
	wfActiveStatus, err := query.GetWFActiveStatus(c, sqlc.GetWFActiveStatusParams{
		Slice:   r.Slice,
		App:     applc,
		Class:   class,
		Realm:   realm,
		Setname: r.Workflow,
	})

	if err != nil || !wfActiveStatus.Bool {
		lh.LogActivity("GetWFinstanceNew||validateWorkflow()||invalid workflow is_active status", err)
		errors = append(errors, wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid_workflow_active_status, &WORKFLOW, fmt.Sprintf("%v", wfActiveStatus.Bool)))
	}

	// To get worflow Internal status
	lh.Debug0().Log("GetWFinstanceNew||validateWorkflow()||verifying whether workflow internal status is valid")
	wfInternalStatus, err := query.GetWFInternalStatus(c, sqlc.GetWFInternalStatusParams{
		Slice:   r.Slice,
		App:     applc,
		Class:   class,
		Realm:   realm,
		Setname: r.Workflow,
	})

	if err != nil || wfInternalStatus {
		lh.LogActivity("GetWFinstanceNew||validateWorkflow()||invalid workflow is_internal status", err)
		errors = append(errors, wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid_workflow_Internal_status, &WORKFLOW, fmt.Sprintf("%v", wfInternalStatus)))

	}

	// There is no record in the wfinstance table with the same values for the tuple (slice,app,workflow,entityid)
	lh.Log("GetWFinstanceNew||validateWorkflow()||verifying whether record is already exist in wfinstance table")
	wfinstanceRecordCount, err := query.GetWFINstance(c, sqlc.GetWFINstanceParams{
		Slice:    r.Slice,
		App:      applc,
		Workflow: r.Workflow,
		Entityid: r.EntityID,
	})

	if err != nil || wfinstanceRecordCount > 0 {
		lh.LogActivity("GetWFinstanceNew||validateWorkflow()||record already exist in wfinstance table ", err)
		errors = append(errors, wscutils.BuildErrorMessage(server.MsgId_AlreadyExist, server.ErrCode_AlreadyExist, &ENTITYID, r.EntityID))
	}

	if len(errors) > 0 {
		return false, errors
	}

	return true, nil
}

// validate entity
func ValidateEntity(e crux.Entity, ps *crux.Schema_t, s *service.Service) (bool, error) {

	lh := s.LogHarbour
	lh.Debug0().Log("Inside validateEntity()")

	// Check if the entity class matches the expected class from the schema
	if e.Class != ps.Class {
		lh.Debug0().Log("GetWFinstanceNew||validateEntity()||entity class does not match the expected class in the schema")
		return false, errors.New("Entity class does not match the expected class in the schema")
	}

	// Validate attributes
	lh.Log("validating entity attributes")
	for name, val := range e.Attrs {
		t := getType(*ps, name)
		if t == "" {
			lh.Debug0().LogActivity("GetWFinstanceNew||validateEntity()||schema does not contain attribute %v", name)
			return false, fmt.Errorf("schema does not contain attribute %v", name)

		}
		_, err := crux.ConvertEntityAttrVal(val, t)
		if err != nil {
			lh.Debug0().LogActivity("GetWFinstanceNew||validateEntity()||attribute %v in entity has value of wrong type", name)
			return false, fmt.Errorf("attribute %v in entity has value of wrong type", name)
		}
	}

	return true, nil
}

// To get type of request entity attributes
func getType(rs crux.Schema_t, name string) string {
	for _, as := range rs.PatternSchema {
		if as.Attr == name {
			return as.ValType
		}
	}
	return ""
}

// To convert byte data to patternschema struct
func byteToPatternSchema(byteData []byte) (*[]crux.PatternSchema_t, error) {
	var ptternSchema *[]crux.PatternSchema_t
	err := json.Unmarshal(byteData, &ptternSchema)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}
	return ptternSchema, nil
}

// To convert request entity into proper Entity structure
func getEntityStructure(req WFInstanceNewRequest, realm string) crux.Entity {

	var attributes = make(map[string]string)
	for key, val := range req.Entity {
		if key != CLASS {
			attributes[key] = val
		}
	}
	entityStruct := crux.Entity{
		Realm: realm,
		App:   req.App,
		Slice: req.Slice,
		Class: req.Entity[CLASS],
		Attrs: attributes,
	}
	return entityStruct
}
