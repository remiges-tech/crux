package wfinstance

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

type Entity struct {
	Class      string             `json:"class"`
	Attributes []EntityAttributes `json:"attr"`
}

type EntityAttributes struct {
	Name string `json:"name"`
	Val  string `json:"val"`
}

func validateWFInstanceNewReq(r WFInstanceNewRequest, s *service.Service, c *gin.Context) (bool, []wscutils.ErrorMessage) {
	lh := s.LogHarbour
	entity := r.Entity
	var errRes []wscutils.ErrorMessage

	lh.Debug0().Log("Inside ValidateWFInstaceNewReq()")
	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		errRes := append(errRes, wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_DatabaseError, nil))
		return false, errRes
	}
	// Validate request
	isValidReq, errAry := validateWorkflow(r, s, c)
	if len(errAry) > 0 || !isValidReq {
		lh.Debug0().LogActivity("Invalid request:", errAry)
		errRes = errAry
		return false, errRes
	}

	// To verify whether entity has valid structure
	_, isKeyExist := entity[CLASS]

	if !(entity != nil && isKeyExist) {
		lh.Debug0().Log("Entity does not match with standard structure")
		errRes = append(errRes, wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid_Entity, &ENTITY))
		return false, errRes
	}

	// To verify whether app,slice,class present in schema and get patternschema against it
	class := entity[CLASS]
	pattern, err := query.WfPatternSchemaGet(c, sqlc.WfPatternSchemaGetParams{
		Slice: r.Slice,
		Class: class,
		App:   r.App,
		Realm: REALM,
	})
	if err != nil {
		lh.LogActivity("failed to get schema pattern from DB:", err)
		errRes = append(errRes, wscutils.BuildErrorMessage(server.MsgId_NoSchemaFound, server.ErrCode_NotFound, nil))
		return false, errRes
	}

	// Unmarshalling byte data to schemapatten struct
	schmapattern, err := byteToPatternSchema(pattern)
	lh.Debug0().LogActivity("patternschema :", schmapattern)
	if err != nil {
		lh.Debug0().LogActivity(" error while converting byte patternschema to struct:", err)
		errRes = append(errRes, wscutils.BuildErrorMessage(server.MsgId_NoSchemaFound, server.ErrCode_Invalid_pattern_schema, nil))
	}

	// Forming  requested entity  into proper Entity struct
	EntityStruct := getEntity(r.Entity)
	lh.Debug1().LogActivity("Entity stucture:", EntityStruct)

	//  To match entity against patternschema
	isValidEntity, err := ValidateEntity(EntityStruct, schmapattern, s)
	if !isValidEntity || err != nil {
		lh.Debug0().LogActivity(" error while validating entity against patternschema:", err)
		errRes = append(errRes, wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid_Entity, &ENTITY))
		return false, errRes

	}
	return true, nil
}

// validate workflow
func validateWorkflow(r WFInstanceNewRequest, s *service.Service, c *gin.Context) (bool, []wscutils.ErrorMessage) {
	var errors []wscutils.ErrorMessage
	lh := s.LogHarbour
	entityClass := r.Entity[CLASS]

	lh.Debug0().Log("Inside validateWorkflow()")
	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Debug0().Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		errors := append(errors, wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_DatabaseError, nil))
		return false, errors
	}

	//1.The value of app specified in the request matches the app ID with which this workflow is associated
	lh.Debug0().Log("verifying whether app present in request is valid")
	app, err := query.GetApp(c, sqlc.GetAppParams{
		Slice: r.Slice,
		App:   r.App,
		Class: entityClass,
		Realm: REALM,
	})

	if err != nil {
		lh.LogActivity("failed to get app from DB:", err.Error())
		errors = append(errors, wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid_APP, &APP, r.App))
	}

	//2. The class of the workflow must match that of entity
	lh.Debug0().Log("verifying whether class present in request is valid")
	class, err := query.GetClass(c, sqlc.GetClassParams{
		Slice: r.Slice,
		App:   app,
		Class: entityClass,
		Realm: REALM,
	})

	if err != nil {
		lh.LogActivity("failed to get class from DB:", err)
		errors = append(errors, wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid_Class, &CLASS, entityClass))
	}

	// 3.The workflow named has is_active == true and internal == false

	//To get worflow active status
	lh.Debug0().Log("verifying whether workflow active status is valid")
	wfActiveStatus, err := query.GetWFActiveStatus(c, sqlc.GetWFActiveStatusParams{
		Slice:   r.Slice,
		App:     app,
		Class:   class,
		Realm:   REALM,
		Setname: r.Workflow,
	})

	if err != nil || !wfActiveStatus.Bool {
		lh.LogActivity("Invalid workflow is_active status :", err)
		errors = append(errors, wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid_workflow_active_status, &WORKFLOW, fmt.Sprintf("%v", wfActiveStatus.Bool)))
	}

	//To get worflow Internal status
	lh.Debug0().Log("verifying whether workflow internal status is valid")
	wfInternalStatus, err := query.GetWFInternalStatus(c, sqlc.GetWFInternalStatusParams{
		Slice:   r.Slice,
		App:     app,
		Class:   class,
		Realm:   REALM,
		Setname: r.Workflow,
	})

	if err != nil || wfInternalStatus {
		lh.LogActivity("Invalid workflow is_internal status:", err)
		errors = append(errors, wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid_workflow_Internal_status, &WORKFLOW, fmt.Sprintf("%v", wfInternalStatus)))

	}

	//4.There is no record in the wfinstance table with the same values for the tuple (slice,app,workflow,entityid)
	lh.Log("verifying whether record is already exist in wfinstance table")
	wfinstanceRecordCount, err := query.GetWFINstance(c, sqlc.GetWFINstanceParams{
		Slice:    r.Slice,
		App:      app,
		Workflow: r.Workflow,
		Entityid: r.EntityID,
	})

	if err != nil || wfinstanceRecordCount > 0 {
		lh.LogActivity("Record already exist in wfinstance table :", err)
		errors = append(errors, wscutils.BuildErrorMessage(server.MsgId_AlreadyExist, server.ErrCode_AlreadyExist, &ENTITYID, r.EntityID))
	}

	if len(errors) > 0 {
		return false, errors
	}

	return true, nil
}

// validate entity
func ValidateEntity(e Entity, ps *types.PatternSchema, s *service.Service) (bool, error) {

	lh := s.LogHarbour
	lh.Debug0().Log("Inside validateEntity()")
	// Check if the entity class matches the expected class from the schema
	if e.Class != ps.Class {
		lh.Debug0().Log("Entity class does not match the expected class in the schema")
		return false, errors.New("Entity class does not match the expected class in the schema")
	}
	// Validate attributes

	lh.Log("validating entity attributes")
	for _, a := range e.Attributes {
		t := getType(*ps, a.Name)
		if t == "" {
			lh.Debug0().LogActivity("schema does not contain attribute %v", a.Name)
			return false, fmt.Errorf("schema does not contain attribute %v", a.Name)

		}
		_, err := convertEntityAttrVal(a.Val, t)
		if err != nil {
			lh.Debug0().LogActivity("attribute %v in entity has value of wrong type", a.Name)
			return false, fmt.Errorf("attribute %v in entity has value of wrong type", a.Name)
		}
	}

	return true, nil
}

// To get type of request entity attributes
func getType(ps types.PatternSchema, name string) string {
	for _, as := range ps.Attr {
		if as.Name == name {
			return as.ValType
		}
	}
	return ""
}

// Converts the string entityAttrVal to its schema-provided type
func convertEntityAttrVal(entityAttrVal string, valType string) (any, error) {
	var entityAttrValConv any
	var err error
	switch valType {
	case typeBool:
		entityAttrValConv, err = strconv.ParseBool(entityAttrVal)
	case typeInt:
		entityAttrValConv, err = strconv.Atoi(entityAttrVal)
	case typeFloat:
		entityAttrValConv, err = strconv.ParseFloat(entityAttrVal, 64)
	case typeStr, typeEnum:
		entityAttrValConv = entityAttrVal
	case typeTS:
		entityAttrValConv, err = time.Parse(timeLayout, entityAttrVal)
	}

	if err != nil {
		return nil, err
	}
	return entityAttrValConv, nil
}

// To convert byte data to patternschema struct
func byteToPatternSchema(byteData []byte) (*types.PatternSchema, error) {
	var response *types.PatternSchema
	err := json.Unmarshal(byteData, &response)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}
	return response, nil
}

// To convert request entity into proper Entity structure
func getEntity(en map[string]string) Entity {
	entity := en

	EntityStruct := Entity{
		Class: entity[CLASS],
	}

	for key, val := range entity {
		if key != CLASS {
			attribute := EntityAttributes{Name: key, Val: val}
			EntityStruct.Attributes = append(EntityStruct.Attributes, attribute)

		}
	}

	return EntityStruct
}
