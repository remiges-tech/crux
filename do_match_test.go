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

func TestDoMatch(t *testing.T) {
	tests := []doMatchTest{}

	/****************
	    BRE tests
	*****************/
	// Simple tests involving entities of class "inventoryitem"
	setupInventoryItemSchema()
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\n\ndoMatch() = %v, \n\nwant        %v\n\n", got, tt.want)
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
	_, _, err := doMatch(sampleEntity, ruleSets[mainRS], ActionSet{}, map[string]bool{})
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
	entity := Entity{inventoryItemClass, []Attr{{"cat", "textbook"}}}

	_, _, err := doMatch(entity, ruleSets["wrongclassrs"], ActionSet{}, map[string]bool{})
	if err == nil {
		t.Errorf("unexpected output when erroneously 'calling' ruleset of different class")
	}
}
