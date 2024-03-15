package schema

import "github.com/remiges-tech/crux/types"

type Schema struct {
	Slice         int32               `json:"slice" validate:"required,gt=0"`
	App           string              `json:"App" validate:"required,alpha"`
	Class         string              `json:"class" validate:"required,lowercase"`
	PatternSchema types.PatternSchema `json:"patternSchema" validate:"required"`
	ActionSchema  types.ActionSchema  `json:"actionSchema" validate:"required"`
}

type updateSchema struct {
	Slice         int32          `json:"slice" validate:"required,gt=0"`
	App           string         `json:"App" validate:"required,alpha"`
	Class         string         `json:"class" validate:"required,lowercase"`
	PatternSchema *patternSchema `json:"patternSchema,omitempty"`
	ActionSchema  *actionSchema  `json:"actionSchema,omitempty"`
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

type attribute struct {
	Name      string   `json:"name" validate:"required"`
	ShortDesc string   `json:"shortdesc" validate:"required"`
	LongDesc  string   `json:"longdesc" validate:"required"`
	ValType   string   `json:"valtype" validate:"required"`
	Vals      []string `json:"vals,omitempty"`
	Enumdesc  []string `json:"enumdesc,omitempty"`
	ValMax    *int32   `json:"valmax,omitempty"`
	ValMin    *int32   `json:"valmin,omitempty"`
	LenMax    *int32   `json:"lenmax,omitempty"`
	LenMin    *int32   `json:"lenmin,omitempty"`
}

type SchemaGetReq struct {
	Slice int32  `json:"slice" validate:"required,gt=0"`
	App   string `json:"app" validate:"required,alpha"`
	Class string `json:"class" validate:"required,alpha"`
}

type SchemaListStruct struct {
	Slice *int32  `form:"slice" validate:"omitempty,gt=0"`
	App   *string `form:"app" validate:"omitempty,alpha"`
	Class *string `form:"class" validate:"omitempty,lowercase"`
}
