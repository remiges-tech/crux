/*
Contains the TestDoMatch() function that sets up and runs all tests for doMatch()



All rulesets must be associated with a schema, so that matchPattern() (and hence doMatch()) works.
matchPattern() needs the schema to know what is meant to be the type of a particular attribute.

However, matchPattern() does not need the action-schema. In some of the tests below, the pattern-schema
is present in the schema, but the action-schema is omitted. This is permissible for the purpose of testing
doMatch() since the action-schema is not used by matchPattern() or doMatch().
*/

package main

import (
	"reflect"
	"testing"
)

type doMatchTest struct {
	name      string
	entity    Entity
	ruleSet   RuleSet
	actionSet ActionSet
	want      ActionSet
}

func deepEqualMap(actualResult, expectedResult ActionSet) bool {
	// Check if tasks slices are both nil or empty
	tasksEqual := (actualResult.tasks == nil && expectedResult.tasks == nil) ||
		(len(actualResult.tasks) == 0 && len(expectedResult.tasks) == 0) ||
		reflect.DeepEqual(actualResult.tasks, expectedResult.tasks)

	// Check if properties maps are both nil or empty
	propertiesEqual := (actualResult.properties == nil && expectedResult.properties == nil) ||
		(len(actualResult.properties) == 0 && len(expectedResult.properties) == 0) ||
		reflect.DeepEqual(actualResult.properties, expectedResult.properties)

	// Return true if both tasks and properties are equal
	return tasksEqual && propertiesEqual
}
func TestDoMatch(t *testing.T) {
	tests := []doMatchTest{}

	/****************
	    BRE tests
	*****************/
	// Simple tests involving entities of class "inventoryitem"
	setupInventoryItemSchema()
	testCache(&tests, t)
	testBasic(&tests)
	testExit(&tests)
	testReturn(&tests)
	// More complex BRE tests
	testTransactions(&tests)
	testPurchases(&tests)
	testOrders(&tests)

	/****************
	    WFE tests
	*****************/
	testUCCCreation(&tests)
	testPrepareAOF(&tests)
	testValidateAOF(&tests)

	// Run all the tests above
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, _ := doMatch(tt.entity, tt.ruleSet, tt.actionSet, map[string]bool{})

			if !deepEqualMap(got, tt.want) {

				t.Errorf("\n\n  doMatch() = %v, \n\nwant        %v\n\n", got, tt.want)

			}
		})
	}

	/******************
	    Error tests
	*******************/
	// Test for cyclical rulesets that could lead to an infinite loop
	testCycleError(t)
	// Test for a THENCALL to a ruleset that's for a different class
	testThenCallWrongClass(t)
}

func testCycleError(t *testing.T) {
	t.Log("Running cycle test")
	setupRuleSetsForCycleError()
	_, _, err := doMatch(sampleEntity, ruleSets[mainRS], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, map[string]bool{})
	if err == nil {
		t.Errorf("test cycle: expected but did not get error")
	}
}

func setupRuleSetsForCycleError() {
	// main ruleset that contains a ThenCall to ruleset "second"
	rule1 := Rule{
		[]RulePatternTerm{
			{"cat", opEQ, "textbook"},
		},
		RuleActions{
			thenCall: "second",
		},
	}
	ruleSets[mainRS] = RuleSet{1, inventoryItemClass, mainRS,
		[]Rule{rule1},
	}

	// "second" ruleset that contains a ThenCall to ruleset "third"
	rule1 = Rule{
		[]RulePatternTerm{
			{"cat", opEQ, "textbook"},
		},
		RuleActions{
			thenCall: "third",
		},
	}
	ruleSets["second"] = RuleSet{1, inventoryItemClass, "second",
		[]Rule{rule1},
	}

	// "third" ruleset that contains a ThenCall back to ruleset "second"
	rule1 = Rule{
		[]RulePatternTerm{
			{"cat", opEQ, "textbook"},
		},
		RuleActions{
			tasks: []string{"testtask"},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"cat", opEQ, "textbook"},
		},
		RuleActions{
			thenCall: "second",
		},
	}
	ruleSets["third"] = RuleSet{1, inventoryItemClass, "third",
		[]Rule{rule1, rule2},
	}
}

func testThenCallWrongClass(t *testing.T) {
	ruleSets["wrongclassrs"] = RuleSet{1, inventoryItemClass, "wrongclassrs",
		[]Rule{Rule{
			[]RulePatternTerm{
				{"cat", opEQ, "textbook"},
			},
			RuleActions{
				thenCall: "winterdisc",
			},
		}},
	}
	entity := Entity{
		class: inventoryItemClass,
		attrs: map[string]string{
			"cat": "textbook",
		},
	}

	_, _, err := doMatch(entity, ruleSets["wrongclassrs"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, map[string]bool{})
	if err == nil {
		t.Errorf("unexpected output when erroneously 'calling' ruleset of different class")
	}
}
