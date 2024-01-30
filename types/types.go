package types

import "github.com/go-playground/validator/v10"

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
