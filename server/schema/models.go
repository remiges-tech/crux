package schema

import "github.com/remiges-tech/crux/types"



type Schema struct {
	Slice         int32               `json:"slice" validate:"required,gt=0"`
	App           string              `json:"App" validate:"required,alpha"`
	Class         string              `json:"class" validate:"required,lowercase"`
	PatternSchema types.PatternSchema `json:"patternSchema"`
	ActionSchema  types.ActionSchema  `json:"actionSchema"`
}

type SchemaListStruct struct {
	Slice *int32  `form:"slice" validate:"omitempty,gt=0"`
	App   *string `form:"app" validate:"omitempty,alpha"`
	Class *string `form:"class" validate:"omitempty,lowercase"`
}
