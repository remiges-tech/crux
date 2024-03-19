/*
Contains the TestDoMatch() function that sets up and runs all tests for doMatch()



All rulesets must be associated with a schema, so that matchPattern() (and hence doMatch()) works.
matchPattern() needs the schema to know what is meant to be the type of a particular attribute.

However, matchPattern() does not need the action-schema. In some of the tests below, the pattern-schema
is present in the schema, but the action-schema is omitted. This is permissible for the purpose of testing
doMatch() since the action-schema is not used by matchPattern() or doMatch().
*/

package crux

import (
	"reflect"
	"strconv"
	"testing"
)

var ruleSetsTests = []*Ruleset_t{}
var ruleSchemasTest = []*Schema_t{}

type doMatchTest struct {
	name      string
	entity    Entity
	ruleSet   *Ruleset_t
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
			// Assuming you have an index, adjust the index below accordingly
			got, _, _ := doMatch(tt.entity, tt.ruleSet, tt.actionSet, map[string]struct{}{})

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

	testGetStats(t)
}

func testCycleError(t *testing.T) {

	var sampleEntity = Entity{
		realm: "1",
		app:   "Test1",
		slice: "1",
		class: "inventoryitem2",
		attrs: map[string]string{
			"cat":        "textbook",
			"fullname":   "Advanced Physics",
			"ageinstock": "5",
			"mrp":        "50.80",
			"received":   "2018-06-01T15:04:05Z",
			"bulkorder":  trueStr,
		},
	}

	t.Log("Running cycle test")

	rs := setupRuleSetsForCycleError()
	_, _, err := doMatch(sampleEntity, rs, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, map[string]struct{}{})
	if err == nil {
		t.Errorf("test cycle: expected but did not get error")
	}
}

func setupRuleSetsForCycleError() *Ruleset_t {

	// main ruleset that contains a ThenCall to ruleset "second"
	rule1 := Rule_t{
		RulePatterns: []RulePatternBlock_t{
			{"cat", opEQ, "textbook"},
		},
		RuleActions: RuleActionBlock_t{
			ThenCall: "second",
		},
	}

	// "second" ruleset that contains a ThenCall to ruleset "third"
	rule2 := Rule_t{
		RulePatterns: []RulePatternBlock_t{
			{"cat", opEQ, "textbook"},
		},
		RuleActions: RuleActionBlock_t{
			ThenCall: "third",
		},
	}

	// "third" ruleset that contains a ThenCall back to ruleset "second"
	rule3 := Rule_t{
		RulePatterns: []RulePatternBlock_t{
			{"cat", opEQ, "textbook"},
		},
		RuleActions: RuleActionBlock_t{
			Task: []string{"testtask"},
		},
	}
	rule4 := Rule_t{
		RulePatterns: []RulePatternBlock_t{
			{"cat", opEQ, "textbook"},
		},
		RuleActions: RuleActionBlock_t{
			ThenCall: "main",
		},
	}

	rs := Ruleset_t{
		Id:      1,
		Class:   "inventoryitem2",
		SetName: "main",
		Rules:   []Rule_t{rule1, rule2, rule3, rule4},
		NCalled: 0,
	}
	return &rs

}

func testThenCallWrongClass(t *testing.T) {

	rule2 := Rule_t{
		RulePatterns: []RulePatternBlock_t{
			{"cat", opEQ, "textbook"},
		},
		RuleActions: RuleActionBlock_t{
			ThenCall: "winterdisc",
		},
	}
	rs := Ruleset_t{
		Id:      1,
		Class:   inventoryItemClass,
		SetName: "wrongclassrs",
		Rules:   []Rule_t{rule2},
		NCalled: 0,
	}

	entity := Entity{
		realm: "1",
		app:   "Test1",
		slice: "1",
		class: "inventoryitem2",
		attrs: map[string]string{
			"cat": "textbook",
		},
	}

	_, _, err := doMatch(entity, &rs, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, map[string]struct{}{})

	if err == nil {
		t.Errorf("unexpected output when erroneously 'calling' ruleset of different class")
	}
}

func testGetStats(t *testing.T) {
	entity := Entity{
		realm: "1",
		app:   "Test1",
		slice: "1",
		class: "inventoryitem2",
		attrs: map[string]string{
			"cat": "textbook",
		},
	}
	realm := realm_t(entity.realm)
	app := app_t(entity.app)
	slice, err := strconv.Atoi(entity.slice)
	if err != nil {
		t.Fatalf("Failed to convert slice to int: %v", err)
	}

	_, _, err1 := getStats(realm, app, slice_t(slice))
	//fmt.Printf("GetStats  time %v \n", timestamp)
	//printStats(stats)
	if err1 != nil {
		t.Errorf("unexpected output when erroneously 'calling' ruleset of different class")
	}

}
