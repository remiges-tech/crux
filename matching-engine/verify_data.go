/*
This file contains verifyRuleSchema(), verifyRuleSet(), doReferentialChecks() and verifyEntity(),
and their helper functions
*/

package crux

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"time"
)

const (
	step       = "step"
	stepFailed = "stepfailed"
	start      = "START"
	nextStep   = "nextstep"
	done       = "done"

	cruxIDRegExp = `^[a-z][a-z0-9_]*$`
)

var validTypes = map[string]bool{
	typeBool: true, typeInt: true, typeFloat: true, typeStr: true, typeEnum: true, typeTS: true,
}

var validOps = map[string]bool{
	opEQ: true, opNE: true, opLT: true, opLE: true, opGT: true, opGE: true,
}

type CruxError struct {
	Keyword   string `json:"keyword"`
	FieldName string `json:"fieldName"`
	Messages  string `json:"messages"`
	Vals      string `json:"vals,omitempty"`
}

// Error returns the error message string
func (e CruxError) Error() string {
	return fmt.Sprintf("%s: %s: %s: %s", e.Keyword, e.FieldName, e.Vals, e.Messages)
}

// Parameters
// rs RuleSchema: the RuleSchema to be verified
// isWF bool: true if the RuleSchema applies to a workflow, otherwise false
func VerifyRuleSchema(rschema []*Schema_t, isWF bool) []error {
	var errs []error
	for _, rs := range rschema {

		if len(rs.Class) == 0 {
			err := CruxError{Keyword: "Empty", FieldName: "class", Messages: "schema class is an empty string"} //fmt.Errorf("schema class is an empty string")
			errs = append(errs, err)
		}
		if err := VerifyPatternSchema(*rs, isWF); err != nil {
			errs = append(errs, err...)
		}
		if err := VerifyActionSchema(*rs, isWF); err != nil {
			errs = append(errs, err...)
		}
	}
	return errs
}

func VerifyPatternSchema(rs Schema_t, isWF bool) []error {
	var errs []error
	if len(rs.PatternSchema) == 0 {
		err := CruxError{Keyword: "Empty", FieldName: "patternSchema", Messages: "pattern-schema is empty"} //fmt.Errorf("pattern-schema for %v is empty", rs.Class)
		errs = append(errs, err)
	}
	re := regexp.MustCompile(cruxIDRegExp)
	// Bools needed for workflows only
	stepFound, stepFailedFound := false, false

	for i, attrSchema := range rs.PatternSchema {
		i++
		if _, exists := validTypes[attrSchema.ValType]; exists {

		} else {
			fieldName := fmt.Sprintf("patternSchema[%d].ValType", i)
			err := CruxError{Keyword: "Invalid", FieldName: fieldName, Vals: attrSchema.ValType, Messages: "not a valid value-type"} //fmt.Errorf("%v is not a valid value-type", attrSchema.ValType)
			errs = append(errs, err)
		}

		if !re.MatchString(attrSchema.Attr) {
			fieldName := fmt.Sprintf("patternSchema[%d].attr", i)
			err := CruxError{Keyword: "Invalid", FieldName: fieldName, Vals: attrSchema.Attr, Messages: "attribute name is not a valid CruxID"} //fmt.Errorf("attribute name %v is not a valid CruxID", attrSchema.Attr)
			errs = append(errs, err)
		} else if attrSchema.ValType == typeEnum && len(attrSchema.EnumVals) == 0 {
			fieldName := fmt.Sprintf("patternSchema[%d].Vals", i)
			err := CruxError{Keyword: "Empty", FieldName: fieldName, Messages: "no valid values for enum"} //fmt.Errorf("no valid values for enum %v", attrSchema.Attr)
			errs = append(errs, err)
		}
		for val := range attrSchema.EnumVals {
			if !re.MatchString(val) && val != start {
				fieldName := fmt.Sprintf("patternSchema[%d].Vals", i)
				err := CruxError{Keyword: "Invalid", FieldName: fieldName, Vals: val, Messages: "enum value is not a valid CruxID"} //fmt.Errorf("enum value %v is not a valid CruxID", val)
				errs = append(errs, err)
			}
		}

		// Workflows only
		if attrSchema.Attr == step && attrSchema.ValType == typeEnum {
			stepFound = true
		}

		if isWF && attrSchema.Attr == step && attrSchema.EnumVals[start] != struct{}{} {
			fieldName := fmt.Sprintf("patternSchema[%d].Vals", i)
			err := CruxError{Keyword: "NotAllowed", FieldName: fieldName, Vals: "Start", Messages: "workflow schema doesn't allow step=START"} //fmt.Errorf("workflow schema for %v doesn't allow step=START", rs.Class)
			errs = append(errs, err)
		}
		if attrSchema.Attr == stepFailed && attrSchema.ValType == typeBool {
			stepFailedFound = true
		}
	}

	// Workflows only
	if isWF && (!stepFound || !stepFailedFound) {
		err := CruxError{Keyword: "Required", FieldName: "attr", Vals: "step/stepfailed", Messages: "necessary attributes aNovant in schema"} //fmt.Errorf("necessary workflow attributes aNovant in schema for class %v", rs.Class)
		errs = append(errs, err)
	}

	return errs

}

func VerifyActionSchema(rs Schema_t, isWF bool) []error {
	var errs []error
	re := regexp.MustCompile(cruxIDRegExp)
	if len(rs.ActionSchema.Tasks) == 0 && len(rs.ActionSchema.Properties) == 0 {
		err := CruxError{Keyword: "Empty", FieldName: "Tasks/Properties", Messages: "both tasks and properties are empty in schema"} //fmt.Errorf("both tasks and properties are empty in schema for class %v", rs.Class)
		errs = append(errs, err)
	}
	for i, task := range rs.ActionSchema.Tasks {
		i++
		if !re.MatchString(task) {
			fieldName := fmt.Sprintf("actionSchema.Tasks[%d]", i)
			err := CruxError{Keyword: "Invalid", FieldName: fieldName, Vals: task, Messages: "task is not a valid CruxID"} // fmt.Errorf("task %v is not a valid CruxID", task)
			errs = append(errs, err)
		}
	}

	// Workflows only

	if isWF && len(rs.ActionSchema.Properties) != 2 {
		err := CruxError{Keyword: "NotAllowed", FieldName: "properties", Messages: "contain exactly two properties"} // //fmt.Errorf("action-schema for %v does not contain exactly two properties", rs.Class)
		errs = append(errs, err)
	}
	nextStepFound, doneFound := false, false

	for i, propName := range rs.ActionSchema.Properties {
		i++
		if !re.MatchString(propName) {
			fieldName := fmt.Sprintf("actionSchema.Properties[%d]", i)
			err := CruxError{Keyword: "Invalid", FieldName: fieldName, Vals: propName, Messages: "property name is not a valid CruxID"} //fmt.Errorf("property name %v is not a valid CruxID", propName)
			errs = append(errs, err)
		} else if propName == nextStep {
			nextStepFound = true
		} else if propName == done {
			doneFound = true
		}
	}

	// Workflows only
	if isWF && (!nextStepFound || !doneFound) {
		err := CruxError{Keyword: "NotAllowed", FieldName: "properties", Messages: "does not contain both the properties 'nextstep' and 'done'"} // fmt.Errorf("action-schema for %v does not contain both the properties 'nextstep' and 'done'", rs.Class)
		errs = append(errs, err)
	}

	if isWF && !reflect.DeepEqual(getTasksMapForWF(rs.ActionSchema.Tasks), getStepAttrVals(rs)) {
		err := CruxError{Keyword: "NotAllowed", FieldName: "task", Messages: "action-schema tasks are not the same as valid values for 'step' in pattern-schema"} // fmt.Errorf("action-schema tasks for %v are not the same as valid values for 'step' in pattern-schema", rs.Class)
		errs = append(errs, err)
	}
	return errs
}

func getTasksMapForWF(tasks []string) map[string]struct{} {
	tm := map[string]struct{}{}
	for _, t := range tasks {
		tm[t] = struct{}{}
	}
	// To allow comparison with the set of valid values for the 'step' attribute, which includes "START"
	// tm[start] = struct{}{}

	return tm
}

func getStepAttrVals(rs Schema_t) map[string]struct{} {

	for _, ps := range rs.PatternSchema {
		if ps.Attr == step {

			return ps.EnumVals
		}
	}
	return nil

}

// Parameters
// rs RuleSet: the RuleSet to be verified
// isWF bool: true if the RuleSet is a workflow, otherwise false
// func verifyRuleSet(entiry Entity, rs *Ruleset_t, isWF bool) []error {
// 	var errs []error
// 	schema, err := getSchema(entiry, entiry.Class)

// 	if err != nil {
// 		errs = append(errs, err)
// 	}
// 	if err := VerifyRulePatterns(rs, schema, isWF); err != nil {
// 		errs = append(errs, err...)
// 	}
// 	if err := VerifyRuleActions(rs, schema, isWF); err != nil {
// 		errs = append(errs, err...)
// 	}

// 	return errs
// }

func VerifyRulePatterns(ruleset *Ruleset_t, schema *Schema_t, isWF bool) []error {
	var errs []error
	for i, rule := range ruleset.Rules {
		i++
		for j, term := range rule.RulePatterns {
			j++
			valType := GetType(schema, term.Attr)

			if !verifyType(term.Val, valType) {
				fieldName := fmt.Sprintf("rule[%d].term[%d]val", i, j)
				err := CruxError{Keyword: "NotMatch", FieldName: fieldName, Messages: "value of this attribute does not match schema"} // fmt.Errorf("value of this attribute does not match schema type: %v", term.Attr)
				errs = append(errs, err)
			}
			if !validOps[term.Op] {
				fieldName := fmt.Sprintf("rule[%d].term[%d]op", i, j)
				err := CruxError{Keyword: "Invalid", FieldName: fieldName, Vals: term.Op, Messages: "invalid operation in rule"} // fmt.Errorf("invalid operation in rule: %v", term.Op)
				errs = append(errs, err)
			}
		}
		// Workflows only
		if isWF {
			stepFound := false
			for j, term := range rule.RulePatterns {
				if term.Attr == step {
					str := term.Val.(string)
					if str == "start" {
						stepFound = true
						break
					} else if !slices.Contains(schema.ActionSchema.Tasks, str) {
						fieldName := fmt.Sprintf("rule[%d].term[%d]Attr", i, j)
						err := CruxError{Keyword: "NotAllowed", FieldName: fieldName, Vals: str, Messages: "step must be present in schema Pattern or action"}
						errs = append(errs, err)
					}
					stepFound = true
					break
				}
			}
			if !stepFound {
				fieldName := fmt.Sprintf("rule[%d].attr", i)
				err := CruxError{Keyword: "Required", FieldName: fieldName, Messages: "required one 'step' attribute in a rule in workflow"} // fmt.Errorf("no 'step' attribute found in a rule in workflow %v", ruleset.SetName)
				errs = append(errs, err)
			}
		}
	}
	return errs
}

// func getSchema(entity Entity, Class string) (*Schema_t, error) {
// 	ruleSchemas, _ := retriveRuleSchemasAndRuleSetsFromCache(entity.Realm, entity.App, entity.Class, entity.Slice)

// 	if ruleSchemas == nil {
// 		return nil, fmt.Errorf("no schema found for Class %v", Class)
// 	}
// 	if Class == ruleSchemas.Class {
// 		return ruleSchemas, nil
// 	}

// 	return nil, fmt.Errorf("no schema found for Class %v", Class)
// }

func GetType(rs *Schema_t, name string) string {
	for _, as := range rs.PatternSchema {
		if as.Attr == name {
			return as.ValType
		}
	}
	return ""
}

func isStringInArray(s string, arr []string) bool {
	for _, a := range arr {

		if a == s {

			return true
		}
	}
	return false
}

// Returns whether or not the type of "val" is the same as "valType"
func verifyType(val any, valType string) bool {

	var ok bool
	switch valType {
	case typeBool:
		_, ok = val.(bool)
	case typeInt:
		_, ok = val.(int)
	case typeFloat:
		_, ok = val.(float64)
	case typeStr:

		_, ok = val.(string)
	case typeEnum:

		_, ok = val.(string)
	case typeTS:
		s, _ := val.(string)
		_, err := time.Parse(timeLayout, s)
		ok = (err == nil)
	}

	return ok
}

func VerifyRuleActions(ruleset *Ruleset_t, schema *Schema_t, isWF bool) []error {
	var errs []error
	//for _, ruleset := range ruleSets {
	for i, rule := range ruleset.Rules {
		i++
		for j, t := range rule.RuleActions.Task {
			j++
			doRet := rule.RuleActions.DoReturn
			doExit := rule.RuleActions.DoExit

			found := false

			if doRet && doExit {
				fieldName := fmt.Sprintf("rule[%d].ruleActions[%d]", i, j)
				err := CruxError{Keyword: "Required", FieldName: fieldName, Messages: "required one of RETURN or EXIT"} // fmt.Errorf("there is a rule with both the RETURN and EXIT instructions in ruleset %v", ruleset.SetName)
				errs = append(errs, err)
			}
			if isStringInArray(t, schema.ActionSchema.Tasks) {
				found = true
				break
			}

			if !found {
				fieldName := fmt.Sprintf("rule[%d].ruleActions[%d].task", i, j)
				err := CruxError{Keyword: "NotExist", FieldName: fieldName, Messages: "task not found in any action-schema"} // fmt.Errorf("task %v not found in any action-schema", t)
				errs = append(errs, err)
			}
		}

		for propName := range rule.RuleActions.Properties {
			found := false

			if isStringInArray(propName, schema.ActionSchema.Properties) {
				found = true
				break
			}

			if !found {
				fieldName := fmt.Sprintf("rule[%d].property", i)
				err := CruxError{Keyword: "NotExist", FieldName: fieldName, Messages: "property name not found in any action-schema"} // fmt.Errorf("property name %v not found in any action-schema", propName)
				errs = append(errs, err)

			}
		}

		// Workflows only
		if isWF {

			nsFound, doneFound := areNextStepAndDoneInProps(rule.RuleActions.Properties)
			if !nsFound && !doneFound {
				fieldName := fmt.Sprintf("rule[%d]", i)
				err := CruxError{Keyword: "NotAllowed", FieldName: fieldName, Messages: "rule found with neither 'nextstep' nor 'done' in ruleset"} // fmt.Errorf("rule found with neither 'nextstep' nor 'done' in ruleset %v", ruleset.SetName)
				errs = append(errs, err)
			}
			if !doneFound && len(rule.RuleActions.Task) == 0 {
				fieldName := fmt.Sprintf("rule[%d]", i)
				err := CruxError{Keyword: "NotExist", FieldName: fieldName, Messages: "no tasks and no 'done=true' in a rule"} // fmt.Errorf("no tasks and no 'done=true' in a rule in ruleset %v", ruleset.SetName)
				errs = append(errs, err)
			}
			// currNS := getNextStep(rule.RuleActions.Properties)
			// fieldName := fmt.Sprintf("rule[%d].task", i)
			// if len(currNS) > 0 && !isStringInArray(currNS, rule.RuleActions.Task) {
			// 	err := CruxError{Keyword: "NotExist", FieldName: fieldName, Messages: "`nextstep` value not found in `tasks` in a rule "} // fmt.Errorf("`nextstep` value not found in `tasks` in a rule in ruleset %v", ruleset.SetName)
			// 	errs = append(errs, err)
			// }
		}
	}

	return errs
}

func areNextStepAndDoneInProps(props map[string]string) (bool, bool) {
	nsFound, doneFound := false, false
	for name, val := range props {

		if name == nextStep {
			nsFound = true
		}
		if name == done && val == trueStr {
			doneFound = true
		}
	}
	return nsFound, doneFound
}

// func getNextStep(props map[string]string) string {
// 	for name, val := range props {
// 		if name == nextStep {
// 			return val
// 		}
// 	}
// 	return ""
// }

// func doReferentialChecks(e Entity) (bool, error) {
// 	_, ruleSets := retriveRuleSchemasAndRuleSetsFromCache(e.Realm, e.App, e.Class, e.Slice)

// 	for _, ruleset := range ruleSets {
// 		for _, rule := range ruleset.Rules {

// 			if rule.RuleActions.ThenCall != "" || rule.RuleActions.ElseCall != "" {

// 				return true, nil
// 			}

// 		}
// 	}
// 	return true, nil
// }

func VerifyEntity(e Entity, rs *Schema_t) error {

	for attrName, attrVal := range e.Attrs {

		t := GetType(rs, attrName)
		if t == "" {
			return fmt.Errorf("schema does not contain attribute %v", attrName)
		}
		_, err := ConvertEntityAttrVal(attrVal, t)
		if err != nil {
			return fmt.Errorf("attribute %v in entity has value of wrong type", attrName)
		}
	}
	if len(e.Attrs) != len(rs.PatternSchema) {
		return fmt.Errorf("entity does not contain all the attributes in its pattern-schema")
	}
	return nil
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
