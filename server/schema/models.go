package schema

import "github.com/remiges-tech/crux/types"

type Schema struct {
	Slice         int32               `json:"slice" validate:"required,gt=0,lt=15"`
	App           string              `json:"App" validate:"required,alpha,lt=15"`
	Class         string              `json:"class" validate:"required,lowercase,lt=15"`
	PatternSchema types.PatternSchema `json:"patternSchema"`
	ActionSchema  types.ActionSchema  `json:"actionSchema"`
}

type SchemaGetReq struct {
	Slice int32  `json:"slice" validate:"required,gt=0,lt=15"`
	App   string `json:"app" validate:"required,alpha,lt=15"`
	Class string `json:"class" validate:"required,alpha,lt=15"`
}

// used in test cases no need for validation
type SchemaListStruct struct {
	Slice *int32  `form:"slice"`
	App   *string `form:"app"`
	Class *string `form:"class"`
}

type attribute struct {
	Name      string   `json:"name,omitempty"`
	ShortDesc string   `json:"shortdesc,omitempty"`
	LongDesc  string   `json:"longdesc,omitempty"`
	ValType   string   `json:"valtype,omitempty"`
	Vals      []string `json:"vals,omitempty"`
	Enumdesc  []string `json:"enumdesc,omitempty"`
	Vallt     int32    `json:"vallt,omitempty"`
	ValMin    int32    `json:"valmin,omitempty"`
	Lenlt     int32    `json:"lenlt,omitempty"`
	LenMin    int32    `json:"lenmin,omitempty"`
}
type patternSchema struct {
	Class string      `json:"class,omitempty"`
	Attr  []attribute `json:"attr,omitempty"`
}
type actionSchema struct {
	Class      string   `json:"class,omitempty"`
	Tasks      []string `json:"tasks,omitempty"`
	Properties []string `json:"properties,omitempty"`
}
