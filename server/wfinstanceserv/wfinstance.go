package wfinstanceserv

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/types"
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

type Entity struct {
	Class      string             `json:"class"`
	Attributes []EntityAttributes `json:"attr"`
}

type EntityAttributes struct {
	Name string `json:"name"`
	Val  string `json:"val"`
}

// GetWFinstanceNew will be responsible for processing the /wfinstanceNew request that comes through as a POST
func GetWFinstanceNew(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("GetWFinstanceNew request received")

	// Bind request
	var wfinstanceNewreq WFInstanceNewRequest
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

	// Validate Entity
	entity := wfinstanceNewreq.Entity
	// The entity must have valid structure
	_, isKeyExist := wfinstanceNewreq.Entity["class"]
	if !(entity != nil && isKeyExist) {
		lh.Debug0().LogActivity("Entity does not match with standard structure:", err.Error)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(INVALID_ENTITY, nil, err.Error())}))

	}

	query, ok := s.Database.(*sqlc.Queries)
	if !ok {
		lh.Log("Error while getting query instance from service")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse("unable_to_get_database_query"))
		return
	}

	// To verify whether app,slice ,class present in schema and get patternschema against it
	class := wfinstanceNewreq.Entity["class"]
	pattern, err := query.WfPatternSchemaGet(c, sqlc.WfPatternSchemaGetParams{
		App:   *wfinstanceNewreq.App,
		Slice: *wfinstanceNewreq.Slice,
		Class: class,
	})
	if err != nil {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(wscutils.ERRCODE_INVALID_REQUEST, nil, err.Error())}))
		lh.Debug0().LogActivity("failed to get data from DB:", err.Error)
		return
	}

	// unmarshalling byte data
	schmapattern, err := byteToPatternSchema(pattern)
	if err != nil {
		lh.Debug0().LogActivity(" error while converting byte patternschema to struct:", err.Error)

		return
	}

	// segregating class and attributes from entity
	EntityStruct := Entity{
		Class: entity["class"],
	}

	for key, val := range entity {
		if key != "class" {
			attribute := EntityAttributes{Name: key, Val: val}
			EntityStruct.Attributes = append(EntityStruct.Attributes, attribute)

		}
	}

	// match entity against patternschema
	isValidEntity, err := validateEntity(EntityStruct, schmapattern, s)
	if !isValidEntity || err != nil {
		lh.Debug0().LogActivity(" error while validating entity against patternschema:", err.Error)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(wscutils.ERRCODE_INVALID_REQUEST, nil, err.Error())}))
		return

	}

	lh.Log(fmt.Sprintf("Record found: %v", map[string]any{"response": schmapattern}))
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(schmapattern))

}

func getValsForGetWFinstanceNewReqError(err validator.FieldError) []string {
	return types.CommonValidation(err)
}

// to convert byte data to patternschema struct
func byteToPatternSchema(byteData []byte) (*types.Patternschema, error) {
	var response *types.Patternschema
	err := json.Unmarshal(byteData, &response)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	return response, nil
}

// validate entity
func validateEntity(e Entity, ps *types.Patternschema, s *service.Service) (bool, error) {
	// all attributes of entity match the names in the schema

	// Check if the entity class matches the expected class from the schema
	if e.Class != ps.Class {
		return false, errors.New("Entity class does not match the expected class in the schema")
	}
	// Validate attributes
	for _, a := range e.Attributes {
		t := getType(*ps, a.Name)
		if t == "" {
			return false, fmt.Errorf("schema does not contain attribute %v", a.Name)
		}
		_, err := convertEntityAttrVal(a.Val, t)
		if err != nil {
			return false, fmt.Errorf("attribute %v in entity has value of wrong type", a.Name)
		}
	}
	// if len(e.Attributes) != len(ps.Attr) {
	// 	return false, fmt.Errorf("entity does not contain all the attributes in its pattern-schema")
	// }

	// all values of the fields in entity match types of the attributes as specified in the schema

	return true, nil
}

func getType(ps types.Patternschema, name string) string {
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
