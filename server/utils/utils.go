package utils

import (
	"fmt"
	"regexp"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/types"
)

const (
	cruxIDRegExp = `^[a-z][a-z0-9_]*$`
)

var validTypes = map[string]bool{
	"int": true, "float": true, "str": true, "enum": true, "bool": true, "timestamps": true,
}

func VerifyPatternSchema(ps types.PatternSchema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	re := regexp.MustCompile(cruxIDRegExp)

	for i, attrSchema := range ps.Attr {
		i++
		if !re.MatchString(attrSchema.Name) {
			fieldName := fmt.Sprintf("attrSchema[%d].Name", i)
			vErr := wscutils.BuildErrorMessage("not_valid", &fieldName, attrSchema.Name)
			validationErrors = append(validationErrors, vErr)
		}
		if !validTypes[attrSchema.ValType] {
			fieldName := fmt.Sprintf("attrSchema[%d].ValType", i)
			vErr := wscutils.BuildErrorMessage("not_valid", &fieldName, attrSchema.ValType)
			validationErrors = append(validationErrors, vErr)
		}
		if attrSchema.ValType == "enum" && len(attrSchema.Vals) == 0 {
			fieldName := fmt.Sprintf("attrSchema[%d].Vals", i)
			vErr := wscutils.BuildErrorMessage("empty", &fieldName)
			validationErrors = append(validationErrors, vErr)
		}
	}
	return validationErrors
}

func VerifyActionSchema(as types.ActionSchema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	re := regexp.MustCompile(cruxIDRegExp)
	for i, task := range as.Tasks {
		if !re.MatchString(task) {
			fieldName := fmt.Sprintf("actionSchema.Tasks[%d]", i)
			vErr := wscutils.BuildErrorMessage("not_valid", &fieldName, task)
			validationErrors = append(validationErrors, vErr)
		}
	}
	for i, propName := range as.Properties {
		if !re.MatchString(propName) {
			fieldName := fmt.Sprintf("actionSchema.Properties[%d]", i)
			vErr := wscutils.BuildErrorMessage("not_valid", &fieldName, propName)
			validationErrors = append(validationErrors, vErr)
		}
	}
	return validationErrors
}
