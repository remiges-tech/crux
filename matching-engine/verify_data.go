/*
This file contains verifyRuleSchema(), verifyRuleSet(), doReferentialChecks() and verifyEntity(),
and their helper functions
*/

package crux

import (
	"fmt"
	"reflect"
	"regexp"
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

// Parameters
// rs RuleSchema: the RuleSchema to be verified
// isWF bool: true if the RuleSchema applies to a workflow, otherwise false
func VerifyRuleSchema(rschema []*Schema_t, isWF bool) (bool, error) {
	for _, rs := range rschema {

		if len(rs.Class) == 0 {
			return false, fmt.Errorf("schema class is an empty string")
		}
		if _, err := VerifyPatternSchema(*rs, isWF); err != nil {
			return false, err
		}
		if _, err := VerifyActionSchema(*rs, isWF); err != nil {
			return false, err
		}
	}
	return true, nil
}

func VerifyPatternSchema(rs Schema_t, isWF bool) (bool, error) {
	if len(rs.PatternSchema) == 0 {
		return false, fmt.Errorf("pattern-schema for %v is empty", rs.Class)
	}
	re := regexp.MustCompile(cruxIDRegExp)
	// Bools needed for workflows only
	stepFound, stepFailedFound := false, false

	for _, attrSchema := range rs.PatternSchema {

		if _, exists := validTypes[attrSchema.ValType]; exists {

		} else {
			return false, fmt.Errorf("%v is not a valid value-type", attrSchema.ValType)
		}

		if !re.MatchString(attrSchema.Attr) {
			return false, fmt.Errorf("attribute name %v is not a valid CruxID", attrSchema.Attr)
		} else if attrSchema.ValType == typeEnum && len(attrSchema.EnumVals) == 0 {
			return false, fmt.Errorf("no valid values for enum %v", attrSchema.Attr)
		}
		for val := range attrSchema.EnumVals {
			if !re.MatchString(val) && val != start {
				return false, fmt.Errorf("enum value %v is not a valid CruxID", val)
			}
		}

		// Workflows only
		if attrSchema.Attr == step && attrSchema.ValType == typeEnum {
			stepFound = true
		}

		if isWF && attrSchema.Attr == step && attrSchema.EnumVals[start] != struct{}{} {
			return false, fmt.Errorf("workflow schema for %v doesn't allow step=START", rs.Class)
		}
		if attrSchema.Attr == stepFailed && attrSchema.ValType == typeBool {
			stepFailedFound = true
		}
	}

	// Workflows only
	if isWF && (!stepFound || !stepFailedFound) {
		return false, fmt.Errorf("necessary workflow attributes absent in schema for class %v", rs.Class)
	}

	return true, nil

}

func VerifyActionSchema(rs Schema_t, isWF bool) (bool, error) {
	re := regexp.MustCompile(cruxIDRegExp)
	if len(rs.ActionSchema.Tasks) == 0 && len(rs.ActionSchema.Properties) == 0 {
		return false, fmt.Errorf("both tasks and properties are empty in schema for class %v", rs.Class)
	}
	for _, task := range rs.ActionSchema.Tasks {
		if !re.MatchString(task) {
			return false, fmt.Errorf("task %v is not a valid CruxID", task)
		}
	}

	// Workflows only

	if isWF && len(rs.ActionSchema.Properties) != 2 {
		return false, fmt.Errorf("action-schema for %v does not contain exactly two properties", rs.Class)
	}
	nextStepFound, doneFound := false, false

	for _, propName := range rs.ActionSchema.Properties {

		if !re.MatchString(propName) {
			return false, fmt.Errorf("property name %v is not a valid CruxID", propName)
		} else if propName == nextStep {
			nextStepFound = true
		} else if propName == done {
			doneFound = true
		}
	}

	// Workflows only
	if isWF && (!nextStepFound || !doneFound) {
		return false, fmt.Errorf("action-schema for %v does not contain both the properties 'nextstep' and 'done'", rs.Class)
	}

	if isWF && !reflect.DeepEqual(getTasksMapForWF(rs.ActionSchema.Tasks), getStepAttrVals(rs)) {
		return false, fmt.Errorf("action-schema tasks for %v are not the same as valid values for 'step' in pattern-schema", rs.Class)
	}
	return true, nil
}

func getTasksMapForWF(tasks []string) map[string]struct{} {
	tm := map[string]struct{}{}
	for _, t := range tasks {
		tm[t] = struct{}{}
	}
	// To allow comparison with the set of valid values for the 'step' attribute, which includes "START"
	tm[start] = struct{}{}

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
func verifyRuleSet(entiry Entity, rs *Ruleset_t, isWF bool) (bool, error) {
	schema, err := getSchema(entiry, entiry.class)

	if err != nil {
		return false, err
	}
	if _, err = verifyRulePatterns(rs, schema, isWF); err != nil {

		return false, err
	}
	if _, err = verifyRuleActions(rs, schema, isWF); err != nil {

		return false, err
	}

	return true, nil
}

func verifyRulePatterns(ruleset *Ruleset_t, schema *Schema_t, isWF bool) (bool, error) {
	//for _, ruleset := range ruleSets {
	for _, rule := range ruleset.Rules {

		for _, term := range rule.RulePatterns {

			valType := getType(schema, term.Attr)

			if valType == "" {
				// If the attribute name is not in the pattern-schema, we check if it's a task "tag"
				// by checking for its presence in the action-schema

				if !isStringInArray(term.Attr, schema.ActionSchema.Tasks) {
					return false, fmt.Errorf("attribute does not exist in schema: %v", term.Attr)
				}
				// If it is a tag, the value type is set to bool
				valType = typeStr
			}
			if !verifyType(term.Val, valType) {

				return false, fmt.Errorf("value of this attribute does not match schema type: %v", term.Attr)
			}
			if !validOps[term.Op] {

				return false, fmt.Errorf("invalid operation in rule: %v", term.Op)
			}
		}
		// Workflows only
		if isWF {
			stepFound := false
			for _, term := range rule.RulePatterns {
				if term.Attr == step {

					stepFound = true
					break
				}
			}
			if !stepFound {

				return false, fmt.Errorf("no 'step' attribute found in a rule in workflow %v", ruleset.SetName)
			}
		}
	}
	//}
	return true, nil
}

func getSchema(entity Entity, class string) (*Schema_t, error) {
	ruleSchemas, _ := retriveRuleSchemasAndRuleSetsFromCache(entity.realm, entity.app, entity.class, entity.slice)

	if len(ruleSchemas) > 0 {

		for _, s := range ruleSchemas {

			if class == s.Class {
				return s, nil
			}
		}
	}
	return nil, fmt.Errorf("no schema found for class %v", class)
}

func getType(rs *Schema_t, name string) string {
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

func verifyRuleActions(ruleset *Ruleset_t, schema *Schema_t, isWF bool) (bool, error) {
	//for _, ruleset := range ruleSets {
	for _, rule := range ruleset.Rules {
		for _, t := range rule.RuleActions.Task {

			doRet := rule.RuleActions.DoReturn
			doExit := rule.RuleActions.DoExit

			found := false

			if doRet && doExit {
				return false, fmt.Errorf("there is a rule with both the RETURN and EXIT instructions in ruleset %v", ruleset.SetName)
			}
			if isStringInArray(t, schema.ActionSchema.Tasks) {
				found = true
				break
			}

			if !found {
				return false, fmt.Errorf("task %v not found in any action-schema", t)
			}
		}

		for propName := range rule.RuleActions.Properties {
			found := false

			if isStringInArray(propName, schema.ActionSchema.Properties) {
				found = true
				break
			}

			if !found {
				return false, fmt.Errorf("property name %v not found in any action-schema", propName)
			}
		}

		// Workflows only
		if isWF {

			nsFound, doneFound := areNextStepAndDoneInProps(rule.RuleActions.Properties)
			if !nsFound && !doneFound {
				return false, fmt.Errorf("rule found with neither 'nextstep' nor 'done' in ruleset %v", ruleset.SetName)
			}
			if !doneFound && len(rule.RuleActions.Task) == 0 {
				return false, fmt.Errorf("no tasks and no 'done=true' in a rule in ruleset %v", ruleset.SetName)
			}
			currNS := getNextStep(rule.RuleActions.Properties)
			if len(currNS) > 0 && !isStringInArray(currNS, rule.RuleActions.Task) {
				return false, fmt.Errorf("`nextstep` value not found in `tasks` in a rule in ruleset %v", ruleset.SetName)
			}
		}
	}
	//}

	return true, nil
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

func getNextStep(props map[string]string) string {
	for name, val := range props {
		if name == nextStep {
			return val
		}
	}
	return ""
}

func doReferentialChecks(e Entity) (bool, error) {
	_, ruleSets := retriveRuleSchemasAndRuleSetsFromCache(e.realm, e.app, e.class, e.slice)

	for _, ruleset := range ruleSets {
		for _, rule := range ruleset.Rules {

			if rule.RuleActions.ThenCall != "" || rule.RuleActions.ElseCall != "" {

				return true, nil
			}

		}
	}
	return true, nil
}

func verifyEntity(e Entity) (bool, error) {
	rs, err := getSchema(e, e.class)
	if err != nil {
		return false, err
	}
	for attrName, attrVal := range e.attrs {

		t := getType(rs, attrName)
		if t == "" {
			return false, fmt.Errorf("schema does not contain attribute %v", attrName)
		}
		_, err := convertEntityAttrVal(attrVal, t)
		if err != nil {
			return false, fmt.Errorf("attribute %v in entity has value of wrong type", attrName)
		}
	}
	if len(e.attrs) != len(rs.PatternSchema) {
		return false, fmt.Errorf("entity does not contain all the attributes in its pattern-schema")
	}
	return true, nil
}
