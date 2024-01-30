package types

import "github.com/go-playground/validator/v10"

type AppConfig struct {
	DBHost        string `json:"db_host"`
	DBPort        int    `json:"db_port"`
	DBUser        string `json:"db_user"`
	DBPassword    string `json:"db_password"`
	DBName        string `json:"db_name"`
	AppServerPort int    `json:"app_server_port"`
	ErrorTypeFile string `json:"error_type_file"`
}

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
	Name      string   `json:"name"`
	ShortName string   `json:"shortname"`
	LongDesc  string   `json:"longdesc"`
	ValType   string   `json:"valtype"`
	Vals      []string `json:"vals,omitempty"`
	Enumdesc  []string `json:"enumdesc,omitempty"`
	ValMax    int32    `json:"valmax,omitempty"`
	ValMin    int32    `json:"valmin,omitempty"`
	LenMax    int32    `json:"lemmax,omitempty"`
	LenMin    int32    `json:"lenmin,omitempty"`
}
type Patternschema struct {
	Class string      `json:"class"`
	Attr  []Attribute `json:"attr"`
}

type Actionschema struct {
	Class      string   `json:"class"`
	Task       []string `json:"task"`
	Properties []string `json:"properties"`
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
	case "lowercase":
		vals = append(vals, "only_lowercase_letter_are_allowed")
	default:
		vals = append(vals, "not_valid_input")
	}
	return vals
}
