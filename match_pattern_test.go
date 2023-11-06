/* This file contains TestMatchPattern() and one helper function */

package main

import (
	"testing"
	"time"
)

func TestMatchPattern(t *testing.T) {
	var testNames []string
	var entities []Entity
	var rulePatterns []([]RulePatternTerm)
	var resultsExpected []any

	setupInventoryItemSchema()
	actionSet := ActionSet{tasks: []string{"dodiscount", "yearendsale"}}

	// Test: many terms, everything matches
	testNames = append(testNames, "everything matches")
	entities = append(entities, sampleEntity)
	receivedTime, _ := time.Parse(timeLayout, "2018-05-15T12:00:00Z")
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"cat", opEQ, "textbook"},
		{"fullname", opEQ, "Advanced Physics"},
		{"ageinstock", opLE, 7},
		{"mrp", opLT, 51.2},
		{"received", opGT, receivedTime},
		{"bulkorder", opNE, false},
		{"dodiscount", opEQ, true},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: many terms, mrp doesn't match
	testNames = append(testNames, "mrp doesn't match")
	entities = append(entities, sampleEntity)
	receivedTime, _ = time.Parse(timeLayout, "2018-05-15T12:00:00Z")
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"cat", opEQ, "textbook"},
		{"fullname", opEQ, "Advanced Physics"},
		{"ageinstock", opLE, 7},
		{"mrp", opGE, 51.2},
		{"received", opGT, receivedTime},
		{"bulkorder", opNE, false},
		{"dodiscount", opEQ, true},
	})
	resultsExpected = append(resultsExpected, false)

	// Test: bool "ne"
	testNames = append(testNames, "bool ne")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"bulkorder", opNE, true},
	})
	resultsExpected = append(resultsExpected, false)

	// Test: enum "ne"
	testNames = append(testNames, "enum ne")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"cat", opNE, "refbook"},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: float "eq"
	testNames = append(testNames, "float eq")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"mrp", opEQ, 50.8},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: float "ge"
	testNames = append(testNames, "float ge")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"mrp", opGE, 50.8},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: timestamp "lt"
	testNames = append(testNames, "timestamp lt")
	entities = append(entities, sampleEntity)
	receivedTime, _ = time.Parse(timeLayout, "2018-06-10T15:04:05Z")
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"received", opLT, receivedTime},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: timestamp "le"
	testNames = append(testNames, "timestamp le")
	entities = append(entities, sampleEntity)
	receivedTime, _ = time.Parse(timeLayout, "2018-05-01T15:04:05Z")
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"received", opLE, receivedTime},
	})
	resultsExpected = append(resultsExpected, false)

	// Test: string "lt"
	testNames = append(testNames, "string lt")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"fullname", opLT, "Advanced Science"},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: string "gt"
	testNames = append(testNames, "string gt")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"fullname", opGT, "Accelerated Physics"},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: task wanted in pattern, found in action set
	testNames = append(testNames, "tasks found in action set")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"dodiscount", opEQ, true},
		{"yearendsale", opNE, false},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: task wanted in pattern, but not found in action set
	testNames = append(testNames, "task not in action set")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"dodiscount", opEQ, true},
		{"summersale", opEQ, true},
	})
	resultsExpected = append(resultsExpected, false)

	// Test: task not wanted in pattern, and not found in action set
	testNames = append(testNames, "task 'eq false' in pattern, and not in action set")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"summersale", opEQ, false},
	})
	resultsExpected = append(resultsExpected, true)

	// Test: edge case - no rule pattern
	testNames = append(testNames, "no rule pattern")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{})
	resultsExpected = append(resultsExpected, true)

	// Test: error converting value
	testNames = append(testNames, "deliberate error converting value")
	entities = append(entities, Entity{inventoryItemClass, []Attr{
		{"ageinstock", "abc"},
	}})
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"ageinstock", opGT, 5},
	})
	resultsExpected = append(resultsExpected, nil)

	// Test: error - not an ordered type
	testNames = append(testNames, "deliberate error: not an ordered type")
	entities = append(entities, sampleEntity)
	rulePatterns = append(rulePatterns, []RulePatternTerm{
		{"bulkorder", opGT, true},
	})
	resultsExpected = append(resultsExpected, nil)

	// Run the tests
	t.Log("==Running", len(rulePatterns), "matchPattern tests==")

	for i, rulePattern := range rulePatterns {
		t.Logf("Test: %s", testNames[i])
		res, err := matchPattern(entities[i], rulePattern, actionSet)
		if resultsExpected[i] == nil && err == nil {
			t.Errorf("Expected but did not get error")
			continue
		} else if resultsExpected[i] == nil && err != nil {
			continue
		} else if res != resultsExpected[i] {
			t.Errorf("FAIL output=%v, expected=%v", res, resultsExpected[i])
		}
	}
}

func setupInventoryItemSchema() {
	ruleSchemas = append(ruleSchemas, RuleSchema{class: inventoryItemClass,
		patternSchema: []AttrSchema{
			{name: "cat", valType: typeEnum},
			{name: "fullname", valType: typeStr},
			{name: "ageinstock", valType: typeInt},
			{name: "mrp", valType: typeFloat},
			{name: "received", valType: typeTS},
			{name: "bulkorder", valType: typeBool},
		},
	})
}
