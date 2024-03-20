/*
Contains the TestDoMatch() function that sets up and runs all tests for DoMatch()



All rulesets must be associated with a schema, so that matchPattern() (and hence DoMatch()) works.
matchPattern() needs the schema to know what is meant to be the type of a particular attribute.

However, matchPattern() does not need the action-schema. In some of the tests below, the pattern-schema
is present in the schema, but the action-schema is omitted. This is permissible for the purpose of testing
DoMatch() since the action-schema is not used by matchPattern() or DoMatch().
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
	// Check if Tasks slices are both nil or empty
	tasksEqual := (actualResult.Tasks == nil && expectedResult.Tasks == nil) ||
		(len(actualResult.Tasks) == 0 && len(expectedResult.Tasks) == 0) ||
		reflect.DeepEqual(actualResult.Tasks, expectedResult.Tasks)

	// Check if Properties maps are both nil or empty
	propertiesEqual := (actualResult.Properties == nil && expectedResult.Properties == nil) ||
		(len(actualResult.Properties) == 0 && len(expectedResult.Properties) == 0) ||
		reflect.DeepEqual(actualResult.Properties, expectedResult.Properties)

	// Return true if both Tasks and Properties are equal
	return tasksEqual && propertiesEqual
}
func TestDoMatch(t *testing.T) {
	tests := []doMatchTest{}

	/****************
	    BRE tests
	*****************/
	// Simple tests involving entities of Class "inventoryitem"
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
			got, _, _ := DoMatch(tt.entity, tt.ruleSet, tt.actionSet, map[string]struct{}{})

			if !deepEqualMap(got, tt.want) {
				t.Errorf("\n\n  DoMatch() = %v, \n\nwant        %v\n\n", got, tt.want)
			}
		})
	}

	/******************
	    Error tests
	*******************/
	// Test for cyclical rulesets that could lead to an infinite loop
	testCycleError(t)
	// Test for a THENCALL to a ruleset that's for a different Class
	testThenCallWrongClass(t)

	testGetStats(t)
}

func testCycleError(t *testing.T) {

	var sampleEntity = Entity{
		Realm: "1",
		App:   "Test1",
		Slice: "1",
		Class: "inventoryitem2",
		Attrs: map[string]string{
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
	_, _, err := DoMatch(sampleEntity, rs, ActionSet{
		Tasks:      []string{},
		Properties: make(map[string]string),
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
		Realm: "1",
		App:   "Test1",
		Slice: "1",
		Class: "inventoryitem2",
		Attrs: map[string]string{
			"cat": "textbook",
		},
	}

	_, _, err := DoMatch(entity, &rs, ActionSet{
		Tasks:      []string{},
		Properties: make(map[string]string),
	}, map[string]struct{}{})

	if err == nil {
		t.Errorf("unexpected output when erroneously 'calling' ruleset of different Class")
	}
}

func testGetStats(t *testing.T) {
	entity := Entity{
		Realm: "1",
		App:   "Test1",
		Slice: "1",
		Class: "inventoryitem2",
		Attrs: map[string]string{
			"cat": "textbook",
		},
	}
	Realm := Realm_t(entity.Realm)
	App := App_t(entity.App)
	Slice, err := strconv.Atoi(entity.Slice)
	if err != nil {
		t.Fatalf("Failed to convert Slice to int: %v", err)
	}

	_, _, err1 := getStats(Realm, App, Slice_t(Slice))
	//fmt.Printf("GetStats  time %v \n", timestamp)
	//printStats(stats)
	if err1 != nil {
		t.Errorf("unexpected output when erroneously 'calling' ruleset of different Class")
	}

}
