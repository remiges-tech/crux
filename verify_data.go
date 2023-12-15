/*
This file contains verifyRuleSchema(), verifyRuleSet(), doReferentialChecks() and verifyEntity(),
and their helper functions
*/

package main

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
func verifyRuleSchema(rs RuleSchema, isWF bool) (bool, error) {
	if len(rs.class) == 0 {
		return false, fmt.Errorf("schema class is empty string")
	}
	if _, err := verifyPatternSchema(rs, isWF); err != nil {
		return false, err
	}
	if _, err := verifyActionSchema(rs, isWF); err != nil {
		return false, err
	}
	return true, nil
}

func verifyPatternSchema(rs RuleSchema, isWF bool) (bool, error) {
	if len(rs.patternSchema) == 0 {
		return false, fmt.Errorf("pattern-schema for %v is empty", rs.class)
	}
	re := regexp.MustCompile(cruxIDRegExp)
	// Bools needed for workflows only
	stepFound, stepFailedFound := false, false

	for _, attrSchema := range rs.patternSchema {
		if !re.MatchString(attrSchema.name) {
			return false, fmt.Errorf("attribute name %v is not a valid CruxID", attrSchema.name)
		} else if !validTypes[attrSchema.valType] {
			return false, fmt.Errorf("%v is not a valid value-type", attrSchema.valType)
		} else if attrSchema.valType == typeEnum && len(attrSchema.vals) == 0 {
			return false, fmt.Errorf("no valid values for enum %v", attrSchema.name)
		}
		for val := range attrSchema.vals {
			if !re.MatchString(val) && val != start {
				return false, fmt.Errorf("enum value %v is not a valid CruxID", val)
			}
		}

		// Workflows only
		if attrSchema.name == step && attrSchema.valType == typeEnum {
			stepFound = true
		}
		if isWF && attrSchema.name == step && !attrSchema.vals[start] {
			return false, fmt.Errorf("workflow schema for %v doesn't allow step=START", rs.class)
		}
		if attrSchema.name == stepFailed && attrSchema.valType == typeBool {
			stepFailedFound = true
		}
	}

	// Workflows only
	if isWF && (!stepFound || !stepFailedFound) {
		return false, fmt.Errorf("necessary workflow attributes absent in schema for class %v", rs.class)
	}

	return true, nil
}

func verifyActionSchema(rs RuleSchema, isWF bool) (bool, error) {
	re := regexp.MustCompile(cruxIDRegExp)
	if len(rs.actionSchema.tasks) == 0 && len(rs.actionSchema.properties) == 0 {
		return false, fmt.Errorf("both tasks and properties are empty in schema for class %v", rs.class)
	}
	for _, task := range rs.actionSchema.tasks {
		if !re.MatchString(task) {
			return false, fmt.Errorf("task %v is not a valid CruxID", task)
		}
	}

	// Workflows only
	if isWF && len(rs.actionSchema.properties) != 2 {
		return false, fmt.Errorf("action-schema for %v does not contain exactly two properties", rs.class)
	}
	nextStepFound, doneFound := false, false

	for _, propName := range rs.actionSchema.properties {
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
		return false, fmt.Errorf("action-schema for %v does not contain both the properties 'nextstep' and 'done'", rs.class)
	}
	if isWF && !reflect.DeepEqual(getTasksMapForWF(rs.actionSchema.tasks), getStepAttrVals(rs)) {
		return false, fmt.Errorf("action-schema tasks for %v are not the same as valid values for 'step' in pattern-schema", rs.class)
	}
	return true, nil
}

func getTasksMapForWF(tasks []string) map[string]bool {
	tm := map[string]bool{}
	for _, t := range tasks {
		tm[t] = true
	}
	// To allow comparison with the set of valid values for the 'step' attribute, which includes "START"
	tm[start] = true
	return tm
}

func getStepAttrVals(rs RuleSchema) map[string]bool {
	for _, ps := range rs.patternSchema {
		if ps.name == step {
			return ps.vals
		}
	}
	return nil
}

// Parameters
// rs RuleSet: the RuleSet to be verified
// isWF bool: true if the RuleSet is a workflow, otherwise false
func verifyRuleSet(rs RuleSet, isWF bool) (bool, error) {
	schema, err := getSchema(rs.class)
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

func verifyRulePatterns(ruleSet RuleSet, schema RuleSchema, isWF bool) (bool, error) {
	for _, rule := range ruleSet.rules {
		for _, term := range rule.rulePattern {
			valType := getType(schema, term.attrName)
			if valType == "" {
				// If the attribute name is not in the pattern-schema, we check if it's a task "tag"
				// by checking for its presence in the action-schema
				if !isStringInArray(term.attrName, schema.actionSchema.tasks) {
					return false, fmt.Errorf("attribute does not exist in schema: %v", term.attrName)
				}
				// If it is a tag, the value type is set to bool
				valType = typeBool
			}
			if !verifyType(term.attrVal, valType) {
				return false, fmt.Errorf("value of this attribute does not match schema type: %v", term.attrName)
			}
			if !validOps[term.op] {
				return false, fmt.Errorf("invalid operation in rule: %v", term.op)
			}
		}
		// Workflows only
		if isWF {
			stepFound := false
			for _, term := range rule.rulePattern {
				if term.attrName == step {
					stepFound = true
					break
				}
			}
			if !stepFound {
				return false, fmt.Errorf("no 'step' attribute found in a rule in workflow %v", ruleSet.setName)
			}
		}
	}
	return true, nil
}

func getSchema(class string) (RuleSchema, error) {
	for _, s := range ruleSchemas {
		if class == s.class {
			return s, nil
		}
	}
	return RuleSchema{}, fmt.Errorf("no schema found for class %v", class)
}

func getType(rs RuleSchema, name string) string {
	for _, as := range rs.patternSchema {
		if as.name == name {
			return as.valType
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
	case typeStr, typeEnum:
		_, ok = val.(string)
	case typeTS:
		s, _ := val.(string)
		_, err := time.Parse(timeLayout, s)
		ok = (err == nil)
	}
	return ok
}

func verifyRuleActions(ruleSet RuleSet, schema RuleSchema, isWF bool) (bool, error) {
	for _, rule := range ruleSet.rules {
		for _, t := range rule.ruleActions.tasks {
			if !isStringInArray(t, schema.actionSchema.tasks) {
				return false, fmt.Errorf("task %v not found in action-schema", t)
			}
		}
		for _, p := range rule.ruleActions.properties {
			if !isStringInArray(p.name, schema.actionSchema.properties) {
				return false, fmt.Errorf("property name %v not found in action-schema", p.name)
			}
		}
		if rule.ruleActions.willReturn && rule.ruleActions.willExit {
			return false, fmt.Errorf("there is a rule with both the RETURN and EXIT instructions in ruleset %v", ruleSet.setName)
		}
		// Workflows only
		if isWF {
			nsFound, doneFound := areNextStepAndDoneInProps(rule.ruleActions.properties)
			if !nsFound && !doneFound {
				return false, fmt.Errorf("rule found with neither 'nextstep' nor 'done' in ruleset %v", ruleSet.setName)
			}
			if !doneFound && len(rule.ruleActions.tasks) == 0 {
				return false, fmt.Errorf("no tasks and no 'done=true' in a rule in ruleset %v", ruleSet.setName)
			}
			currNS := getNextStep(rule.ruleActions.properties)
			if len(currNS) > 0 && !isStringInArray(currNS, rule.ruleActions.tasks) {
				return false, fmt.Errorf("`nextstep` value not found in `tasks` in a rule in ruleset %v", ruleSet.setName)
			}
		}
	}
	return true, nil
}

func areNextStepAndDoneInProps(props []Property) (bool, bool) {
	nsFound, doneFound := false, false
	for _, p := range props {
		if p.name == nextStep {
			nsFound = true
		}
		if p.name == done && p.val == trueStr {
			doneFound = true
		}
	}
	return nsFound, doneFound
}

func getNextStep(props []Property) string {
	for _, p := range props {
		if p.name == nextStep {
			return p.val
		}
	}
	return ""
}

func doReferentialChecks() (bool, error) {
	for _, rs := range ruleSets {
		for _, r := range rs.rules {
			ok, err := doesRuleSetExist(r.ruleActions.thenCall)
			if !ok {
				return false, err
			}
			ok, err = doesRuleSetExist(r.ruleActions.elseCall)
			if !ok {
				return false, err
			}
		}
	}
	return true, nil
}

func doesRuleSetExist(name string) (bool, error) {
	// a zero-length THENCALL/ELSECALL means that we don't need to check for the existence of
	// any ruleset and can simply return "true"
	if len(name) == 0 {
		return true, nil
	}
	_, found := ruleSets[name]
	if !found {
		return false, fmt.Errorf("ruleset %v does not exist", name)
	}
	return true, nil
}

func verifyEntity(e Entity) (bool, error) {
	rs, err := getSchema(e.class)
	if err != nil {
		return false, err
	}
	for _, a := range e.attrs {
		t := getType(rs, a.name)
		if t == "" {
			return false, fmt.Errorf("schema does not contain attribute %v", a.name)
		}
		_, err := convertEntityAttrVal(a.val, t)
		if err != nil {
			return false, fmt.Errorf("attribute %v in entity has value of wrong type", a.name)
		}
	}
	if len(e.attrs) != len(rs.patternSchema) {
		return false, fmt.Errorf("entity does not contain all the attributes in its pattern-schema")
	}
	return true, nil
}
