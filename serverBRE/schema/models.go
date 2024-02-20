package schema

import "github.com/remiges-tech/crux/types"

type Schema struct {
	Slice         int32               `json:"slice" validate:"required,gt=0"`
	App           string              `json:"App" validate:"required,alpha"`
	Class         string              `json:"class" validate:"required,lowercase"`
	PatternSchema types.PatternSchema `json:"patternSchema" validate:"required"`
	ActionSchema  types.ActionSchema  `json:"actionSchema" validate:"required"`
}
