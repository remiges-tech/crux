/*
This file contains TestVerifySchema(), TestVerifyRuleSet(), TestDoReferentialChecks() and TestVerifyEntity()
(and their helper functions)
*/

package crux

import (
	"testing"
)

const (
	incorrectOutputRS     = "incorrect output when verifying ruleset with "
	incorrectOutputWF     = "incorrect output when verifying workflow with "
	incorrectOutputEntity = "incorrect output when verifying entity with "

	uccCreation = "ucccreation"
)

type verifySchemaTest struct {
	name    string
	rs      []*Schema_t
	isWF    bool
	want    bool
	wantErr bool
}

func TestVerifySchema(t *testing.T) {
	tests := []verifySchemaTest{}

	/* Business rules schema tests */
	// the only test that involves no error, because the schema is correct
	testCorrectBRSchema(&tests)
	// in the rest of these tests, verifyRuleSchema() should return an error
	testSchemaEmptyClass(&tests)
	testEmptyPatternSchema(&tests)
	testAttrNameIsNotCruxID(&tests)
	testInvalidValType(&tests)
	testNoValsForEnum(&tests)
	testEnumValIsNotCruxID(&tests)
	testBothTasksAndPropsEmpty(&tests)
	testTaskIsNotCruxID(&tests)
	testPropNameNotCruxID(&tests)

	/* Workflow schema tests */
	// the only test that involves no error, because the workflow schema is correct
	testCorrectWFSchema(&tests)
	// in the rest of these tests, verifyRuleSchema() should return an error
	testMissingStart(&tests)
	testMissingStep(&tests)
	testAdditionalProps(&tests)
	testMissingNextStep(&tests)
	testTasksAndStepDiscrepancy(&tests)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := VerifyRuleSchema(tt.rs, tt.isWF)

			if (err != nil) != tt.wantErr {
				t.Errorf("tt.name %v error = %v, wantErr %v tt.want %v \n", tt.name, err, tt.wantErr, tt.want)
				return
			}
			if got != tt.want {
				t.Errorf("got = %v, want %v \n", got, tt.want)

			}
		})
	}
}

func testCorrectBRSchema(tests *[]verifySchemaTest) {
	ruleSchemasTest = make([]*Schema_t, 0)
	rs := &Schema_t{
		Class: transactionClass,
		PatternSchema: []PatternSchema_t{

			{Attr: "productname", ValType: typeStr},
			{Attr: "price", ValType: typeInt},
			{Attr: "inwintersale", ValType: typeBool},
			{Attr: "paymenttype", ValType: typeEnum, EnumVals: map[string]struct{}{"cash": {}, "card": {}}},
			{Attr: "ismember", ValType: typeBool},
		},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{"freepen", "freemug", "freebag"},
			Properties: []string{"discount", "pointsmult"},
		},
	}

	// Append the newly created instance to the ruleSchemasTest slice
	ruleSchemasTest = append(ruleSchemasTest, rs)
	*tests = append(*tests, verifySchemaTest{
		name:    "correct business-rules schema",
		rs:      ruleSchemasTest,
		isWF:    false,
		want:    true,
		wantErr: false,
	})
}

func testSchemaEmptyClass(tests *[]verifySchemaTest) {
	ruleSchemasTest = make([]*Schema_t, 0)
	rs := &Schema_t{
		Class: "",
		PatternSchema: []PatternSchema_t{

			{Attr: "productname", ValType: typeStr},
			{Attr: "price", ValType: typeInt},
			{Attr: "inwintersale", ValType: typeBool},
			{Attr: "paymenttype", ValType: typeEnum, EnumVals: map[string]struct{}{"cash": {}, "card": {}}},
			{Attr: "ismember", ValType: typeBool},
		},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{"freepen", "freemug", "freebag"},
			Properties: []string{"discount", "pointsmult"},
		},
	}

	// Append the newly created instance to the ruleSchemasTest slice
	ruleSchemasTest = append(ruleSchemasTest, rs)
	*tests = append(*tests, verifySchemaTest{
		name:    "schema with empty class",
		rs:      ruleSchemasTest,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testEmptyPatternSchema(tests *[]verifySchemaTest) {
	ruleSchemasTest = make([]*Schema_t, 0)

	rs := &Schema_t{
		Class:         transactionClass,
		PatternSchema: []PatternSchema_t{},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{"freepen", "freemug", "freebag"},
			Properties: []string{"discount", "pointsmult"},
		},
	}
	ruleSchemasTest = append(ruleSchemasTest, rs)
	*tests = append(*tests, verifySchemaTest{
		name:    "empty pattern schema",
		rs:      ruleSchemasTest,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testAttrNameIsNotCruxID(tests *[]verifySchemaTest) {
	ruleSchemasTest = make([]*Schema_t, 0)
	rs := &Schema_t{
		Class: transactionClass,
		PatternSchema: []PatternSchema_t{

			{Attr: "1productname", ValType: typeStr},
			{Attr: "price", ValType: typeInt},
			{Attr: "inwintersale", ValType: typeBool},
			{Attr: "paymenttype", ValType: typeEnum, EnumVals: map[string]struct{}{"cash": {}, "card": {}}},
			{Attr: "ismember", ValType: typeBool},
		},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{"freepen", "freemug", "freebag"},
			Properties: []string{"discount", "pointsmult"},
		},
	}
	ruleSchemasTest = append(ruleSchemasTest, rs)

	*tests = append(*tests, verifySchemaTest{
		name:    "attr name is not CruxID",
		rs:      ruleSchemasTest,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testInvalidValType(tests *[]verifySchemaTest) {
	ruleSchemasTest = make([]*Schema_t, 0)
	rs := &Schema_t{
		Class: transactionClass,
		PatternSchema: []PatternSchema_t{

			{Attr: "productname", ValType: typeStr},
			{Attr: "price", ValType: typeInt},
			// "abc" is not a valid valType
			{Attr: "inwintersale", ValType: "abc"},
			{Attr: "paymenttype", ValType: typeEnum, EnumVals: map[string]struct{}{"cash": {}, "card": {}}},
			{Attr: "ismember", ValType: typeBool},
		},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{"freepen", "freemug", "freebag"},
			Properties: []string{"discount", "pointsmult"},
		},
	}
	ruleSchemasTest = append(ruleSchemasTest, rs)

	*tests = append(*tests, verifySchemaTest{
		name:    "invalid value type",
		rs:      ruleSchemasTest,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testNoValsForEnum(tests *[]verifySchemaTest) {
	ruleSchemasTest = make([]*Schema_t, 0)
	rs := &Schema_t{
		Class: transactionClass,
		PatternSchema: []PatternSchema_t{

			{Attr: "productname", ValType: typeStr},
			{Attr: "price", ValType: typeInt},
			{Attr: "inwintersale", ValType: typeBool},
			// The "EnumVals" "hash-set" below, which is the set of valid values for the
			// enum "paymenttype", should not be empty
			{Attr: "paymenttype", ValType: typeEnum, EnumVals: map[string]struct{}{}},
			{Attr: "ismember", ValType: typeBool},
		},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{"freepen", "freemug", "freebag"},
			Properties: []string{"discount", "pointsmult"},
		},
	}
	ruleSchemasTest = append(ruleSchemasTest, rs)

	*tests = append(*tests, verifySchemaTest{
		name:    "no vals for enum",
		rs:      ruleSchemasTest,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testEnumValIsNotCruxID(tests *[]verifySchemaTest) {

	ruleSchemasTest = make([]*Schema_t, 0)
	rs := &Schema_t{
		Class: transactionClass,
		PatternSchema: []PatternSchema_t{

			{Attr: "productname", ValType: typeStr},
			{Attr: "price", ValType: typeInt},
			{Attr: "inwintersale", ValType: typeBool},
			// "1cash" is not a CruxID
			{Attr: "paymenttype", ValType: typeEnum, EnumVals: map[string]struct{}{"1cash": {}, "card": {}}},
			{Attr: "ismember", ValType: typeBool},
		},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{"freepen", "freemug", "freebag"},
			Properties: []string{"discount", "pointsmult"},
		},
	}
	ruleSchemasTest = append(ruleSchemasTest, rs)

	*tests = append(*tests, verifySchemaTest{
		name:    "enum val is not CruxID",
		rs:      ruleSchemasTest,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testBothTasksAndPropsEmpty(tests *[]verifySchemaTest) {
	ruleSchemasTest = make([]*Schema_t, 0)
	rs := &Schema_t{
		Class: transactionClass,
		PatternSchema: []PatternSchema_t{

			{Attr: "productname", ValType: typeStr},
			{Attr: "price", ValType: typeInt},
			{Attr: "inwintersale", ValType: typeBool},
			{Attr: "paymenttype", ValType: typeEnum, EnumVals: map[string]struct{}{"cash": {}, "card": {}}},
			{Attr: "ismember", ValType: typeBool},
		},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{},
			Properties: []string{},
		},
	}
	ruleSchemasTest = append(ruleSchemasTest, rs)

	*tests = append(*tests, verifySchemaTest{
		name:    "both tasks and properties empty",
		rs:      ruleSchemasTest,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testTaskIsNotCruxID(tests *[]verifySchemaTest) {
	ruleSchemasTest = make([]*Schema_t, 0)
	rs := &Schema_t{
		Class: transactionClass,
		PatternSchema: []PatternSchema_t{

			{Attr: "productname", ValType: typeStr},
			{Attr: "price", ValType: typeInt},
			{Attr: "inwintersale", ValType: typeBool},
			{Attr: "paymenttype", ValType: typeEnum, EnumVals: map[string]struct{}{"cash": {}, "card": {}}},
			{Attr: "ismember", ValType: typeBool},
		},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{"freepen", "free*mug", "freebag"},
			Properties: []string{"discount", "pointsmult"},
		},
	}
	ruleSchemasTest = append(ruleSchemasTest, rs)

	*tests = append(*tests, verifySchemaTest{
		name:    "task is not CruxID",
		rs:      ruleSchemasTest,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testPropNameNotCruxID(tests *[]verifySchemaTest) {
	ruleSchemasTest = make([]*Schema_t, 0)
	rs := &Schema_t{
		Class: transactionClass,
		PatternSchema: []PatternSchema_t{

			{Attr: "productname", ValType: typeStr},
			{Attr: "price", ValType: typeInt},
			{Attr: "inwintersale", ValType: typeBool},
			{Attr: "paymenttype", ValType: typeEnum, EnumVals: map[string]struct{}{"cash": {}, "card": {}}},
			{Attr: "ismember", ValType: typeBool},
		},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{"freepen", "freemug", "freebag"},
			Properties: []string{"Discount", "pointsmult"},
		},
	}
	ruleSchemasTest = append(ruleSchemasTest, rs)

	*tests = append(*tests, verifySchemaTest{
		name:    "property name is not CruxID",
		rs:      ruleSchemasTest,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testCorrectWFSchema(tests *[]verifySchemaTest) {

	ruleSchemasTest = make([]*Schema_t, 0)
	rs := &Schema_t{
		Class: transactionClass,
		PatternSchema: []PatternSchema_t{

			{Attr: step, ValType: typeEnum, EnumVals: map[string]struct{}{start: {}, "getcustdetails": {}, "aof": {}, "sendauthlinktoclient": {}}},
			{Attr: stepFailed, ValType: typeBool},
			{Attr: "mode", ValType: typeEnum, EnumVals: map[string]struct{}{"physical": {}, "demat": {}}},
		},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			Properties: []string{nextStep, done},
		},
	}

	ruleSchemasTest = append(ruleSchemasTest, rs)
	*tests = append(*tests, verifySchemaTest{
		name:    "correct workflow schema",
		rs:      ruleSchemasTest,
		isWF:    true,
		want:    true,
		wantErr: false,
	})
}

func testMissingStart(tests *[]verifySchemaTest) {

	rs := &Schema_t{
		Class: uccCreationClass,
		PatternSchema: []PatternSchema_t{

			{Attr: step, ValType: typeEnum, EnumVals: map[string]struct{}{"getcustdetails": {}, "aof": {}, "sendauthlinktoclient": {}}},
			{Attr: stepFailed, ValType: typeBool},
			{Attr: "mode", ValType: typeEnum, EnumVals: map[string]struct{}{"physical": {}, "demat": {}}},
		},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			Properties: []string{nextStep, done},
		},
	}

	ruleSchemasTest = append(ruleSchemasTest, rs)

	*tests = append(*tests, verifySchemaTest{
		name:    "missing START",
		rs:      ruleSchemasTest,
		isWF:    true,
		want:    false,
		wantErr: true,
	})
}

func testMissingStep(tests *[]verifySchemaTest) {

	rs := &Schema_t{
		Class: uccCreationClass,
		PatternSchema: []PatternSchema_t{

			{Attr: stepFailed, ValType: typeBool},
			{Attr: "mode", ValType: typeEnum, EnumVals: map[string]struct{}{"physical": {}, "demat": {}}},
		},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			Properties: []string{nextStep, done},
		},
	}
	ruleSchemasTest = append(ruleSchemasTest, rs)
	*tests = append(*tests, verifySchemaTest{
		name:    "missing step",
		rs:      ruleSchemasTest,
		isWF:    true,
		want:    false,
		wantErr: true,
	})

}

func testAdditionalProps(tests *[]verifySchemaTest) {
	rs := &Schema_t{
		Class: uccCreationClass,
		PatternSchema: []PatternSchema_t{

			{Attr: step, ValType: typeEnum,
				EnumVals: map[string]struct{}{start: {}, "getcustdetails": {}, "aof": {}, "sendauthlinktoclient": {}},
			},
			{Attr: stepFailed, ValType: typeBool},
			{Attr: "mode", ValType: typeEnum, EnumVals: map[string]struct{}{"physical": {}, "demat": {}}},
		},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			Properties: []string{nextStep, done, "abcd"},
		},
	}
	ruleSchemasTest = append(ruleSchemasTest, rs)

	*tests = append(*tests, verifySchemaTest{
		name:    "additional property other than nextstep and done",
		rs:      ruleSchemasTest,
		isWF:    true,
		want:    false,
		wantErr: true,
	})
}

func testMissingNextStep(tests *[]verifySchemaTest) {

	rs := &Schema_t{
		Class: uccCreationClass,
		PatternSchema: []PatternSchema_t{

			{Attr: step, ValType: typeEnum,
				EnumVals: map[string]struct{}{start: {}, "getcustdetails": {}, "aof": {}, "sendauthlinktoclient": {}},
			},
			{Attr: stepFailed, ValType: typeBool},
			{Attr: "mode", ValType: typeEnum, EnumVals: map[string]struct{}{"physical": {}, "demat": {}}},
		},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			Properties: []string{done, "abcd"},
		},
	}
	ruleSchemasTest = append(ruleSchemasTest, rs)

	*tests = append(*tests, verifySchemaTest{
		name:    "missing nextstep",
		rs:      ruleSchemasTest,
		isWF:    true,
		want:    false,
		wantErr: true,
	})
}

func testTasksAndStepDiscrepancy(tests *[]verifySchemaTest) {

	rs := &Schema_t{
		Class: uccCreationClass,
		PatternSchema: []PatternSchema_t{

			{Attr: step, ValType: typeEnum,
				EnumVals: map[string]struct{}{start: {}, "getcustdetails": {}, "aof": {}, "sendauthlinktoclient": {}},
			},
			{Attr: stepFailed, ValType: typeBool},
			{Attr: "mode", ValType: typeEnum, EnumVals: map[string]struct{}{"physical": {}, "demat": {}}},
		},
		ActionSchema: ActionSchema_t{
			Tasks:      []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			Properties: []string{nextStep, done},
		},
	}
	ruleSchemasTest = append(ruleSchemasTest, rs)
	*tests = append(*tests, verifySchemaTest{
		name:    "tasks and steps discrepancy",
		rs:      ruleSchemasTest,
		isWF:    true,
		want:    false,
		wantErr: true,
	})
}

func TestVerifyRuleSet(t *testing.T) {

	/* Business rules tests */
	setupPurchaseRuleSchema()
	setupRuleSetForPurchases()
	// the only two tests that involve no error, because the ruleset is correct
	testCorrectRS(t)
	testTaskAsAttrName(t)
	// in the rest of these tests, verifyRuleSet() should return an error
	testInvalidAttrName(t)
	testWrongAttrValType(t)
	testInvalidOp(t)
	testTaskNotInSchema(t)
	testPropNameNotInSchema(t)
	testBothReturnAndExit(t)

	/* Workflow tests */
	setupUCCCreationSchema()
	setupUCCCreationRuleSet()
	// the only test that involves no error, because the ruleset is correct
	testCorrectWF(t)
	// in the rest of these tests, verifyRuleSet() should return an error
	testWFRuleMissingStep(t)
	testWFRuleMissingBothNSAndDone(t)
	testWFNoTasksAndNotDone(t)
	testWFNextStepValNotInTasks(t)
}

var sampleEntityVerify = Entity{
	realm: "1",
	app:   "Test8",
	slice: "8",
	class: "purchase",
	attrs: map[string]string{
		"cat":        "textbook",
		"fullname":   "Advanced Physics",
		"ageinstock": "5",
		"mrp":        "50.80",
		"received":   "2018-06-01T15:04:05Z",
		"bulkorder":  trueStr,
	},
}

func testCorrectRS(t *testing.T) {

	rc := Ruleset_t{
		Id:      1,
		Class:   purchaseClass,
		SetName: mainRS,

		NCalled: 0,
	}
	rc.Rules = []rule_t{{
		RulePatterns: []rulePatternBlock_t{},
	}}
	ok, err := verifyRuleSet(sampleEntityVerify, &rc, false)

	if !ok || err != nil {
		t.Errorf(incorrectOutputRS + "no issues")
	}

}

// In each of the rule-pattern tests below, a rule-pattern is modified temporarily.
// After each test, we must reset the rule-pattern to the correct one below before
// moving on to the next test.
var correctRP = []rulePatternBlock_t{
	{"product", opEQ, "jacket"},
	{"price", opGT, 50.0},
}

func testInvalidAttrName(t *testing.T) {

	rc := Ruleset_t{
		Id:    1,
		Class: mainRS,

		NCalled: 0,
	}
	rc.Rules = []rule_t{{
		RulePatterns: []rulePatternBlock_t{{"product", opEQ, "jacket"},
			// priceabc is not in the schema
			{"priceabc", opGT, "50.0"}},
	}}

	ok, err := verifyRuleSet(sampleEntityVerify, &rc, false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRS + "invalid attr name")
	}

}

func testTaskAsAttrName(t *testing.T) {

	rc := Ruleset_t{
		Id:    1,
		Class: mainRS,

		NCalled: 0,
	}
	rc.Rules = []rule_t{{
		RulePatterns: []rulePatternBlock_t{{"product", opEQ, "jacket"},
			// freejar is not in the pattern-schema, but it is a task in the action-schema
			{"freejar", opEQ, ""}},
	}}

	ok, err := verifyRuleSet(sampleEntityVerify, &rc, false)

	if !ok || err != nil {
		t.Errorf(incorrectOutputRS + "a task 'tag' as an attribute name")

	}
	ruleSetsTests = append(ruleSetsTests, &Ruleset_t{})
}

func testWrongAttrValType(t *testing.T) {

	rc := Ruleset_t{
		Id:    1,
		Class: mainRS,

		NCalled: 0,
	}
	rc.Rules = []rule_t{{
		RulePatterns: []rulePatternBlock_t{{"product", opEQ, "jacket"},
			// price should be a float, not a string
			{"price", opGT, "abc"}},
	}}

	ok, err := verifyRuleSet(sampleEntityVerify, &rc, false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRS + "wrong attribute value type")
	}
	ruleSetsTests = append(ruleSetsTests, &Ruleset_t{})
}

func testInvalidOp(t *testing.T) {

	rc := Ruleset_t{
		Id:    1,
		Class: mainRS,

		NCalled: 0,
	}
	rc.Rules = []rule_t{{
		RulePatterns: []rulePatternBlock_t{{"product", opEQ, "jacket"},
			// it should be "gt" (opGT), not "greater than"
			{"price", "greater than", "50.0"}},
	}}

	ok, err := verifyRuleSet(sampleEntityVerify, &rc, false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRS + "invalid operation")
	}
	ruleSetsTests = append(ruleSetsTests, &Ruleset_t{})
}

func testTaskNotInSchema(t *testing.T) {

	rc := Ruleset_t{
		Id:    1,
		Class: mainRS,

		NCalled: 0,
	}
	rc.Rules = []rule_t{{
		RulePatterns: []rulePatternBlock_t{},
		// freeeraser is not in the schema
		RuleActions: ruleActionBlock_t{Task: []string{"freemug", "freeeraser"},
			Properties: map[string]string{"discount": "20"}},
	}}

	ok, err := verifyRuleSet(sampleEntityVerify, &rc, false)

	if !ok || err != nil {

		t.Errorf(incorrectOutputRS + "task not in schema")
	}
	ruleSetsTests = append(ruleSetsTests, &Ruleset_t{})
}

func testPropNameNotInSchema(t *testing.T) {

	rc := Ruleset_t{
		Id:    1,
		Class: mainRS,

		NCalled: 0,
	}
	rc.Rules = []rule_t{{
		RulePatterns: []rulePatternBlock_t{},
		RuleActions: ruleActionBlock_t{Task: []string{"freemug", "freejar", "freeplant"},
			Properties: map[string]string{"cashback": "5"}},
	}}

	ok, err := verifyRuleSet(sampleEntityVerify, &rc, false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRS + "property name not in schema")
	}
	ruleSetsTests = append(ruleSetsTests, &Ruleset_t{})
}

func testBothReturnAndExit(t *testing.T) {

	rc := Ruleset_t{
		Id:    1,
		Class: mainRS,

		NCalled: 0,
	}
	rc.Rules = []rule_t{{
		RulePatterns: []rulePatternBlock_t{},
		RuleActions: ruleActionBlock_t{Task: []string{"freemug", "freejar", "freeplant"},
			Properties: map[string]string{"discount": "20"},
			DoReturn:   true,
			DoExit:     true},
	}}

	ok, err := verifyRuleSet(sampleEntityVerify, &rc, false)

	if ok || err == nil {

		t.Errorf(incorrectOutputRS + "both RETURN and EXIT instructions")

	}
	ruleSetsTests = append(ruleSetsTests, &Ruleset_t{})
}

var sampleEntityUCC = Entity{
	realm: "1",
	app:   "Test10",
	slice: "10",
	class: "ucccreation",

	attrs: map[string]string{
		"cat":        "textbook",
		"fullname":   "Advanced Physics",
		"ageinstock": "5",
		"mrp":        "50.80",
		"received":   "2018-06-01T15:04:05Z",
		"bulkorder":  trueStr,
	},
}

func testCorrectWF(t *testing.T) {

	rc := setupUCCCreationRuleSet()

	ok, err := verifyRuleSet(sampleEntityUCC, rc, true)

	if !ok || err != nil {
		t.Errorf(incorrectOutputWF + "no issues")

	}

}

func testWFRuleMissingStep(t *testing.T) {

	// Modify the specific rule within the ruleset

	rc := Ruleset_t{
		Id:    1,
		Class: uccCreation,

		NCalled: 0,
	}
	rc.Rules = []rule_t{{
		RulePatterns: []rulePatternBlock_t{{stepFailed, opEQ, true},
			{"mode", opEQ, "physical"}},
		RuleActions: ruleActionBlock_t{},
	}}

	ok, err := verifyRuleSet(sampleEntityUCC, &rc, true)
	if ok || err == nil {
		t.Errorf(incorrectOutputWF + "a rule missing 'step'")
	}
	// Reset to original correct rule-pattern
	ruleSetsTests = append(ruleSetsTests, &Ruleset_t{})

}

func testWFRuleMissingBothNSAndDone(t *testing.T) {
	// Assuming uccCreation is an index in the rs slice
	rc := Ruleset_t{
		Id:    1,
		Class: uccCreation,

		NCalled: 0,
	}
	rc.Rules = []rule_t{{
		RulePatterns: []rulePatternBlock_t{},
		RuleActions: ruleActionBlock_t{
			Task:       []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
			Properties: map[string]string{}, // Properties below should contain at least one of "nextstep" and "done"
		},
	}}

	// Call the verifyRuleSet function with the modified ruleset
	ok, err := verifyRuleSet(sampleEntityUCC, &rc, true)

	// Check if the verification fails as expected
	if ok || err == nil {
		t.Errorf(incorrectOutputWF + "a rule missing both 'nextstep' and 'done'")
	}

	// Reset the ruleActions to the original value after the test
	ruleSetsTests = append(ruleSetsTests, &Ruleset_t{})
}

func testWFNoTasksAndNotDone(t *testing.T) {

	rc := Ruleset_t{
		Id:    1,
		Class: uccCreation,

		NCalled: 0,
	}
	rc.Rules = []rule_t{{
		RulePatterns: []rulePatternBlock_t{},
		RuleActions: ruleActionBlock_t{
			Task:       []string{},
			Properties: map[string]string{nextStep: "abc"},
		},
	}}

	// Call the verifyRuleSet function with the modified ruleset
	ok, err := verifyRuleSet(sampleEntityUCC, &rc, true)

	// Check if the verification fails as expected
	if ok || err == nil {
		t.Errorf(incorrectOutputWF + "a rule missing both 'nextstep' and 'done'")
	}

	// Reset the ruleActions to the original value after the test
	ruleSetsTests = append(ruleSetsTests, &Ruleset_t{})
}

func testWFNextStepValNotInTasks(t *testing.T) {

	rc := Ruleset_t{
		Id:    1,
		Class: uccCreation,

		NCalled: 0,
	}
	rc.Rules = []rule_t{{
		RulePatterns: []rulePatternBlock_t{},
		RuleActions: ruleActionBlock_t{
			Task:       []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
			Properties: map[string]string{nextStep: "abcd"},
		},
	}}

	ok, err := verifyRuleSet(sampleEntityUCC, &rc, true)
	if ok || err == nil {
		t.Errorf(incorrectOutputWF + "a 'nextstep' value not in its rule's 'tasks'")
	}
	ruleSetsTests = append(ruleSetsTests, &Ruleset_t{})
}

func TestDoReferentialChecks(t *testing.T) {
	testNoReferentialIssues(t)
	testWrongThenCall(t)
	testWrongElseCall(t)
}

func testNoReferentialIssues(t *testing.T) {
	setupRuleSetsForTransaction()
	ok, err := doReferentialChecks(sampleEntityVerify)
	if !ok || err != nil {
		t.Errorf("unexpected output when there are no referential issues")
	}
}

func testWrongThenCall(t *testing.T) {

	ruleSetsTests[0].Class = mainRS
	ruleSetsTests[0].Rules = []rule_t{{
		RulePatterns: []rulePatternBlock_t{},
		RuleActions: ruleActionBlock_t{
			ThenCall: "summerdisc",
		},
	}}

	ok, err := doReferentialChecks(sampleEntityVerify)
	if !ok || err != nil {
		t.Errorf("unexpected output when there is an incorrect THENCALL")
	}

	ruleSetsTests = append(ruleSetsTests, &Ruleset_t{})
}

func testWrongElseCall(t *testing.T) {

	var sampleEntityVerify1 = Entity{
		realm: "1",
		app:   "Test1",
		slice: "1",
		class: "inventoryitem1",
		attrs: map[string]string{
			"cat":        "textbook",
			"fullname":   "Advanced Physics",
			"ageinstock": "5",
			"mrp":        "50.80",
			"received":   "2018-06-01T15:04:05Z",
			"bulkorder":  trueStr,
		},
	}
	// there is no ruleset "normaldisc"

	ruleSetsTests[0].Class = "inventoryitem1"
	ruleSetsTests[0].Rules = []rule_t{{
		RulePatterns: []rulePatternBlock_t{},
		RuleActions: ruleActionBlock_t{
			ElseCall: "normaldisc",
		},
	}}
	ok, err := doReferentialChecks(sampleEntityVerify1)
	if !ok || err != nil {
		t.Errorf("unexpected output when there is an incorrect ELSECALL")
	}
	// reset the ELSECALL back to "regulardisc"
	ruleSetsTests = append(ruleSetsTests, &Ruleset_t{})
}

func TestVerifyEntity(t *testing.T) {
	// The tests below use this schema
	setupPurchaseRuleSchema()

	testCorrectEntity(t)
	testEntityWithoutSchema(t)
	testEntityWrongAttr(t)
	testEntityWrongType(t)
	testEntityMissingAttr(t)

}

func testCorrectEntity(t *testing.T) {
	e := Entity{
		realm: "1",
		app:   "Test8",
		slice: "8",
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "jacket",
			"price":    "50",
			"ismember": trueStr,
		},
	}

	ok, err := verifyEntity(e)
	if !ok || err != nil {
		t.Errorf(incorrectOutputEntity + "no issues")
	}
}

func testEntityWithoutSchema(t *testing.T) {
	e := Entity{
		class: "wrongclass",
		attrs: map[string]string{
			"product": "jacket",
		},
	}

	ok, err := verifyEntity(e)
	if ok || err == nil {
		t.Errorf(incorrectOutputEntity + "no schema")
	}
}

func testEntityWrongAttr(t *testing.T) {
	e := Entity{
		realm: "1",
		app:   "Test8",
		slice: "8",
		class: purchaseClass,
		attrs: map[string]string{
			"product": "jacket",
			// discount is not in the schema
			"discount": "5",
			"ismember": trueStr,
		},
	}

	ok, err := verifyEntity(e)
	if ok || err == nil {
		t.Errorf(incorrectOutputEntity + "an attribute not in the schema")
	}
}

func testEntityWrongType(t *testing.T) {
	e := Entity{
		realm: "1",
		app:   "Test8",
		slice: "8",
		class: purchaseClass,
		attrs: map[string]string{
			"product": "jacket",
			// price should be a float, not the string "fifty"
			"price":    "fifty",
			"ismember": trueStr,
		},
	}

	ok, err := verifyEntity(e)
	if ok || err == nil {
		t.Errorf(incorrectOutputEntity + "a wrongly-typed attribute")
	}
}

func testEntityMissingAttr(t *testing.T) {
	e := Entity{
		realm: "1",
		app:   "Test8",
		slice: "8",
		class: purchaseClass,
		attrs: map[string]string{
			"product": "jacket",
			// price is missing
			"ismember": trueStr,
		},
	}

	ok, err := verifyEntity(e)
	if ok || err == nil {
		t.Errorf(incorrectOutputEntity + "a missing attribute")
	}
}
