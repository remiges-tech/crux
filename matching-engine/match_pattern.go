/* This file contains matchPattern(), and helper functions called by matchPattern() */

package crux

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

func matchPattern(entity Entity, rulePattern []RulePatternBlock_t, actionSet ActionSet, rSchema []*Schema_t) (bool, error) {

	for _, term := range rulePattern {
		valType := ""
		entityAttrVal := ""

		// Check whether the attribute name in the pattern term exists in the entity Attrs map
		if val, ok := entity.Attrs[term.Attr]; ok {

			entityAttrVal = val
			valType = getTypeFromSchema(entity.Class, term.Attr, rSchema)
			incrementStatsSchemaCounterNChecked(entity.Class, rSchema)
		}

		// If the attribute value is still empty, check whether it matches any of the Tasks in the action-set
		if entityAttrVal == "" {
			for _, task := range actionSet.Tasks {
				if task == term.Attr {
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

		matched, err := makeComparison(entityAttrVal, term.Val, valType, term.Op)

		if err != nil {
			return false, fmt.Errorf("error making comparison %w", err)
		}

		if !matched {
			return false, nil
		}

	}

	return true, nil
}

func getTypeFromSchema(Class string, attrName string, ruleSchemas []*Schema_t) string {

	for _, ruleSchema := range ruleSchemas {

		if ruleSchema.Class == Class {
			for _, attrSchema := range ruleSchema.PatternSchema {

				if attrSchema.Attr == attrName {

					return attrSchema.ValType
				}
			}
		}
	}
	return ""
}

// Returns whether or not the comparison represented by {entityAttrVal, op, termAttrVal} is true
// For example, {7, gt (greater than), 5} is true but {3, gt, 5} is false
func makeComparison(entityAttrVal string, termAttrVal any, valType string, op string) (bool, error) {
	entityAttrValConv, err := ConvertEntityAttrVal(entityAttrVal, valType)

	if err != nil {
		return false, fmt.Errorf("error converting value: %w", err)
	}
	// switch op {
	// case opEQ:
	// 	return entityAttrValConv == termAttrVal, nil
	// case opNE:
	// 	return entityAttrValConv != termAttrVal, nil
	// }
	orderedTypes := map[string]bool{typeInt: true, typeFloat: true, typeTS: true, typeStr: true}
	if !orderedTypes[valType] {
		return false, errors.New("not an ordered type")
	}
	var result int8
	var match bool
	switch op {
	case opEQ:
		result, err = compare(entityAttrValConv, termAttrVal)
		match = (result == 0)
	case opNE:
		result, err = compare(entityAttrValConv, termAttrVal)
		match = (result == -1)
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
func ConvertEntityAttrVal(entityAttrVal string, valType string) (any, error) {

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
	var lessThan bool
	switch a.(type) {
	case int:
		bInt, _ := strconv.Atoi(b.(string))
		if a.(int) < bInt {
			lessThan = true
		} else if a.(int) == bInt {
			return 0, nil
		}
	case float64:
		bFloat, _ := strconv.ParseFloat(b.(string), 64)
		if a.(float64) < bFloat {
			lessThan = true
		} else if a.(float64) == bFloat {
			return 0, nil
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
func incrementStatsSchemaCounterNChecked(Class string, schema []*Schema_t) {

	for _, ruleSchema := range schema {

		if ruleSchema.Class == Class {

			ruleSchema.NChecked++
		}
	}

}
