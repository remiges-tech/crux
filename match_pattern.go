/* This file contains matchPattern(), and helper functions called by matchPattern() */

package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	typeBool  = "bool"
	typeInt   = "int"
	typeFloat = "float"
	typeStr   = "str"
	typeEnum  = "enum"
	typeTS    = "ts"

	timeLayout = "2006-01-02T15:04:05Z"

	opEQ = "eq"
	opNE = "ne"
	opLT = "lt"
	opLE = "le"
	opGT = "gt"
	opGE = "ge"

	trueStr  = "true"
	falseStr = "false"
)

func matchPattern(entity Entity, rulePattern []RulePatternTerm, actionSet ActionSet) (bool, error) {
	for _, term := range rulePattern {
		valType := ""
		entityAttrVal := ""

		// Check whether the attribute name in the pattern term exists in the entity attrs map
		if val, ok := entity.attrs[term.attrName]; ok {
			entityAttrVal = val
			valType = getTypeFromSchema(entity.class, term.attrName)
		}

		// If the attribute value is still empty, check whether it matches any of the tasks in the action-set
		if entityAttrVal == "" {
			for _, task := range actionSet.tasks {
				if task == term.attrName {
					entityAttrVal = trueStr
					valType = typeBool
				}
			}
		}

		// If the attribute value is still empty, default to false
		if entityAttrVal == "" {
			entityAttrVal = falseStr
			valType = typeBool
		}

		matched, err := makeComparison(entityAttrVal, term.attrVal, valType, term.op)
		if err != nil {
			return false, fmt.Errorf("error making comparison %w", err)
		}

		if !matched {
			return false, nil
		}
	}

	return true, nil
}

func getTypeFromSchema(class string, attrName string) string {
	for _, ruleSchema := range ruleSchemas {

		if ruleSchema.class == class {
			for _, attrSchema := range ruleSchema.patternSchema {

				if attrSchema.name == attrName {
					return attrSchema.valType
				}
			}
		}
	}
	return ""
}

// Returns whether or not the comparison represented by {entityAttrVal, op, termAttrVal} is true
// For example, {7, gt (greater than), 5} is true but {3, gt, 5} is false
func makeComparison(entityAttrVal string, termAttrVal any, valType string, op string) (bool, error) {
	entityAttrValConv, err := convertEntityAttrVal(entityAttrVal, valType)
	if err != nil {
		return false, fmt.Errorf("error converting value: %w", err)
	}
	switch op {
	case opEQ:
		return entityAttrValConv == termAttrVal, nil
	case opNE:
		return entityAttrValConv != termAttrVal, nil
	}
	orderedTypes := map[string]bool{typeInt: true, typeFloat: true, typeTS: true, typeStr: true}
	if !orderedTypes[valType] {
		return false, errors.New("not an ordered type")
	}
	var result int8
	var match bool
	switch op {
	case opLT:
		result, err = compare(entityAttrValConv, termAttrVal)
		match = (result == -1)
	case opLE:
		result, err = compare(entityAttrValConv, termAttrVal)
		match = (result == -1) || (result == 0)
	case opGT:
		result, err = compare(entityAttrValConv, termAttrVal)
		match = (result == 1)
	case opGE:
		result, err = compare(entityAttrValConv, termAttrVal)
		match = (result == 1) || (result == 0)
	}
	if err != nil {
		return false, fmt.Errorf("error making comparison %w", err)
	}
	return match, nil
}

// Converts the string entityAttrVal to its schema-provided type
func convertEntityAttrVal(entityAttrVal string, valType string) (any, error) {
	var entityAttrValConv any
	var err error
	switch valType {
	case typeBool:
		entityAttrValConv, err = strconv.ParseBool(entityAttrVal)
	case typeInt:
		entityAttrValConv, err = strconv.Atoi(entityAttrVal)
	case typeFloat:
		entityAttrValConv, err = strconv.ParseFloat(entityAttrVal, 64)
	case typeStr, typeEnum:
		entityAttrValConv = entityAttrVal
	case typeTS:
		entityAttrValConv, err = time.Parse(timeLayout, entityAttrVal)
	}
	if err != nil {
		return nil, err
	}
	return entityAttrValConv, nil
}

// The compare function returns:
// 0 if a == b,
// -1 if a < b, or
// 1 if a > b
func compare(a any, b any) (int8, error) {
	if a == b {
		return 0, nil
	}
	var lessThan bool
	switch a.(type) {
	case int:
		if a.(int) < b.(int) {
			lessThan = true
		}
	case float64:
		if a.(float64) < b.(float64) {
			lessThan = true
		}
	case string:
		if a.(string) < b.(string) {
			lessThan = true
		}
	case time.Time:
		if a.(time.Time).Before(b.(time.Time)) {
			lessThan = true
		}
	default:
		return -2, errors.New("invalid type")
	}
	if lessThan {
		return -1, nil
	} else {
		return 1, nil
	}
}
