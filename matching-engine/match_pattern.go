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

func matchPattern(entity Entity, rulePattern []RulePatternBlock_t, actionSet ActionSet, rSchema *Schema_t) (bool, error) {
	var (
		// inPattern bool
		termval interface{}
		err     error
	)
	for _, term := range rulePattern {
		valType := ""
		entityAttrVal := ""

		// Check whether the attribute name in the pattern term exists in the entity attrs map
		if val, ok := entity.Attrs[term.Attr]; ok {
			// inPattern = true
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

		// if inPattern {
		termval, err = convertTermAttrVal(term.Val, valType)
		if err != nil {
			return false, err
		}
		// }

		matched, err := makeComparison(entityAttrVal, termval, valType, term.Op)
		if err != nil {
			return false, fmt.Errorf("error making comparison %w", err)
		}

		if !matched {
			return false, nil
		}

	}

	return true, nil
}

// Function to convert termAttrVal to the expected type based on valType
func convertTermAttrVal(termAttrVal any, valType string) (interface{}, error) {
	switch valType {
	case typeInt:
		if val, ok := termAttrVal.(string); ok {
			intValue, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("error converting termAttrVal to int: %w", err)
			}
			return intValue, nil
		}
		return termAttrVal, nil
	case typeFloat:
		if val, ok := termAttrVal.(string); ok {
			floatValue, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, fmt.Errorf("error converting termAttrVal to float64: %w", err)
			}
			return floatValue, nil
		}
		return termAttrVal, nil
	case typeBool:
		if val, ok := termAttrVal.(string); ok {
			boolValue, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("error converting termAttrVal to bool: %w", err)
			}
			return boolValue, nil
		}
		return termAttrVal, nil
	case typeTS:
		// Assuming termAttrVal is already in the correct type
		return termAttrVal, nil
	case typeStr, typeEnum:
		// Assuming termAttrVal is already in the correct type
		return termAttrVal, nil
	default:
		return nil, errors.New("unsupported valType")
	}
}

func getTypeFromSchema(class string, attrName string, ruleSchema *Schema_t) string {

	if ruleSchema == nil {
		return ""
	}
	if ruleSchema.Class == class {
		for _, attrSchema := range ruleSchema.PatternSchema {

			if attrSchema.Attr == attrName {

				return attrSchema.ValType
			}
		}
	}
	return ""
}

// Returns whether or not the comparison represented by {entityAttrVal, op, termAttrVal} is true
// For example, {7, gt (greater than), 5} is true but {3, gt, 5} is false
// Returns whether or not the comparison represented by {entityAttrVal, op, termAttrVal} is true
// For example, {7, gt (greater than), 5} is true but {3, gt, 5} is false
func makeComparison(entityAttrVal string, termAttrVal any, valType string, op string) (bool, error) {
	entityAttrValConv, err := ConvertEntityAttrVal(entityAttrVal, valType)

	if err != nil {
		return false, fmt.Errorf("error converting value: %w", err)
	}
	switch op {
	case opEQ:

		if entityAttrValConv != termAttrVal {
			return false, nil
		}
		return true, nil
	case opNE:
		if entityAttrValConv == termAttrVal {
			return false, nil
		}
		return true, nil
	}
	orderedTypes := map[string]bool{typeInt: true, typeFloat: true, typeTS: true, typeStr: true, typeBool: true}
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
	case bool:
		if a.(bool) == b.(bool) {
			return 0, nil
		}
	case int:
		if a.(int) < b.(int) {
			lessThan = true
		}
		if a.(int) == b.(int) {
			return 0, nil
		}
	case float64:
		if a.(float64) < b.(float64) {
			lessThan = true
		}
		if a.(float64) == b.(float64) {
			return 0, nil
		}
	case string:
		if a.(string) < b.(string) {
			lessThan = true
		}
		if a.(string) == b.(string) {
			return 0, nil
		}
	case time.Time:
		if a.(time.Time).Before(b.(time.Time)) {
			lessThan = true
		}
		if a.(time.Time) == b.(time.Time) {
			return 0, nil
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

func incrementStatsSchemaCounterNChecked(class string, schema *Schema_t) {

	if schema != nil {

		if schema.Class == class {

			schema.NChecked++
		}
	}

}
