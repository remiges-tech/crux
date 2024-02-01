package types

import "github.com/go-playground/validator/v10"

const (
	DevEnv           Environment = "dev_env"
	ProdEnv          Environment = "prod_env"
	UATEnv           Environment = "uat_env"
	RECORD_NOT_EXIST             = "record_does_not_exist"
	OPERATION_FAILED             = "operation_failed"
)

type Environment string

type OpReq struct {
	User      string   `json:"user"`
	CapNeeded []string `json:"capNeeded"`
	Scope     Scope    `json:"scope"`
	Limit     Limit    `json:"limit"`
}

type Scope map[string]interface{}
type Limit map[string]interface{}

type QualifiedCap struct {
	Id    string `json:"id"`
	Cap   string `json:"cap"`
	Scope Scope  `json:"scope"`
	Limit Limit  `json:"limit"`
}

type Capabilities struct {
	Name          string         `json:"name"` //either user name or group name
	QualifiedCaps []QualifiedCap `json:"qualifiedcaps"`
}

type Attribute struct {
	Name      string   `json:"name" validate:"required"`
	ShortName string   `json:"shortname" validate:"required"`
	LongDesc  string   `json:"longdesc" validate:"required"`
	ValType   string   `json:"valtype" validate:"required"`
	Vals      []string `json:"vals,omitempty"`
	Enumdesc  []string `json:"enumdesc,omitempty"`
	ValMax    int32    `json:"valmax,omitempty"`
	ValMin    int32    `json:"valmin,omitempty"`
	LenMax    int32    `json:"lemmax,omitempty"`
	LenMin    int32    `json:"lenmin,omitempty"`
}
type Patternschema struct {
	Class string      `json:"class" validate:"required"`
	Attr  []Attribute `json:"attr"`
}

type Actionschema struct {
	Class      string   `json:"class" validate:"required"`
	Task       []string `json:"tasks" validate:"required"`
	Properties []string `json:"properties" validate:"required"`
}

func (env Environment) IsValid() bool {
	switch env {
	case DevEnv, ProdEnv, UATEnv:
		return true
	}
	return false
}

// CommonValidation is a generic function which setup standard validation utilizing
// validator package and Maps the errorVals based on the map parameter and
// return []errorVals
func CommonValidation(err validator.FieldError) []string {
	var vals []string
	switch err.Tag() {
	case "required":
		vals = append(vals, "not_provided")
	case "alpha":
		vals = append(vals, "only_alphabets_are_allowed")
	default:
		vals = append(vals, "not_valid_input")
	}
	return vals
}
