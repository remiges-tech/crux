/*
This file contains TestVerifySchema(), TestVerifyRuleSet(), TestDoReferentialChecks() and TestVerifyEntity()
(and their helper functions)
*/

package main

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
	rs      RuleSchema
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
			got, err := verifyRuleSchema(tt.rs, tt.isWF)
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyRuleSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("verifyRuleSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testCorrectBRSchema(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "freemug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "correct business-rules schema",
		rs:      rs,
		isWF:    false,
		want:    true,
		wantErr: false,
	})
}

func testSchemaEmptyClass(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: "",
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "freemug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "schema with empty class",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testEmptyPatternSchema(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "freemug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "empty pattern schema",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testAttrNameIsNotCruxID(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			// 1productname is not a CruxID
			{name: "1productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "freemug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "attr name is not CruxID",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testInvalidValType(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			// "abc" is not a valid valType
			{name: "inwintersale", valType: "abc"},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "freemug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "invalid value type",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testNoValsForEnum(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			// The "vals" "hash-set" below, which is the set of valid values for the
			// enum "paymenttype", shold not be empty
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "freemug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "no vals for enum",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testEnumValIsNotCruxID(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			// 1cash is not a CruxID
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"1cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"freepen", "freemug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "enum val is not CruxID",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testBothTasksAndPropsEmpty(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		// Both tasks and properties should not be empty
		actionSchema: ActionSchema{
			tasks:      []string{},
			properties: []string{},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "both tasks and properties empty",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testTaskIsNotCruxID(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			// free*mug is not a CruxID
			tasks:      []string{"freepen", "free*mug", "freebag"},
			properties: []string{"discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "task is not CruxID",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testPropNameNotCruxID(tests *[]verifySchemaTest) {
	rs := RuleSchema{class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum, vals: map[string]bool{"cash": true, "card": true}},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks: []string{"freepen", "freemug", "freebag"},
			// Discount is not a CruxID
			properties: []string{"Discount", "pointsmult"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "property name is not CruxID",
		rs:      rs,
		isWF:    false,
		want:    false,
		wantErr: true,
	})
}

func testCorrectWFSchema(tests *[]verifySchemaTest) {
	rs := RuleSchema{
		class: uccCreationClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeEnum,
				vals: map[string]bool{start: true, "getcustdetails": true, "aof": true, "sendauthlinktoclient": true},
			},
			{name: stepFailed, valType: typeBool},
			{name: "mode", valType: typeEnum, vals: map[string]bool{"physical": true, "demat": true}},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			properties: []string{nextStep, done},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "correct workflow schema",
		rs:      rs,
		isWF:    true,
		want:    true,
		wantErr: false,
	})
}

func testMissingStart(tests *[]verifySchemaTest) {
	rs := RuleSchema{
		class: uccCreationClass,
		patternSchema: []AttrSchema{
			// vals below should also contain '"START": true'
			{name: step, valType: typeEnum,
				vals: map[string]bool{"getcustdetails": true, "aof": true, "sendauthlinktoclient": true},
			},
			{name: stepFailed, valType: typeBool},
			{name: "mode", valType: typeEnum, vals: map[string]bool{"physical": true, "demat": true}},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			properties: []string{nextStep, done},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "missing START",
		rs:      rs,
		isWF:    true,
		want:    false,
		wantErr: true,
	})
}

func testMissingStep(tests *[]verifySchemaTest) {
	rs := RuleSchema{
		class: uccCreationClass,
		patternSchema: []AttrSchema{
			// there should be a "step" attribute-schema here
			{name: stepFailed, valType: typeBool},
			{name: "mode", valType: typeEnum, vals: map[string]bool{"physical": true, "demat": true}},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			properties: []string{nextStep, done},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "missing step",
		rs:      rs,
		isWF:    true,
		want:    false,
		wantErr: true,
	})
}

func testAdditionalProps(tests *[]verifySchemaTest) {
	rs := RuleSchema{
		class: uccCreationClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeEnum,
				vals: map[string]bool{start: true, "getcustdetails": true, "aof": true, "sendauthlinktoclient": true},
			},
			{name: stepFailed, valType: typeBool},
			{name: "mode", valType: typeEnum, vals: map[string]bool{"physical": true, "demat": true}},
		},
		actionSchema: ActionSchema{
			tasks: []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			// "abcd" should not be in properties
			properties: []string{nextStep, done, "abcd"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "additional property other than nextstep and done",
		rs:      rs,
		isWF:    true,
		want:    false,
		wantErr: true,
	})
}

func testMissingNextStep(tests *[]verifySchemaTest) {
	rs := RuleSchema{
		class: uccCreationClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeEnum,
				vals: map[string]bool{start: true, "getcustdetails": true, "aof": true, "sendauthlinktoclient": true},
			},
			{name: stepFailed, valType: typeBool},
			{name: "mode", valType: typeEnum, vals: map[string]bool{"physical": true, "demat": true}},
		},
		actionSchema: ActionSchema{
			tasks: []string{"getcustdetails", "aof", "sendauthlinktoclient"},
			// properties should contain "nextstep" (and should not contain "abcd")
			properties: []string{done, "abcd"},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "missing nextstep",
		rs:      rs,
		isWF:    true,
		want:    false,
		wantErr: true,
	})
}

func testTasksAndStepDiscrepancy(tests *[]verifySchemaTest) {
	rs := RuleSchema{
		class: uccCreationClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeEnum,
				// "vals" should have exactly the same strings as "tasks" below, except "start" which is only in "vals"
				vals: map[string]bool{start: true, "getcustdetails": true, "aof": true, "sendauthlinktoclient": true},
			},
			{name: stepFailed, valType: typeBool},
			{name: "mode", valType: typeEnum, vals: map[string]bool{"physical": true, "demat": true}},
		},
		actionSchema: ActionSchema{
			// "tasks" should have exactly the same strings as "vals" above, except for "start"
			tasks:      []string{"getcustinfo", "aof", "sendauthlinktoclient"},
			properties: []string{nextStep, done},
		},
	}
	*tests = append(*tests, verifySchemaTest{
		name:    "tasks and steps discrepancy",
		rs:      rs,
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

func testCorrectRS(t *testing.T) {
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if !ok || err != nil {
		t.Errorf(incorrectOutputRS + "no issues")
	}
}

// In each of the rule-pattern tests below, a rule-pattern is modified temporarily.
// After each test, we must reset the rule-pattern to the correct one below before
// moving on to the next test.
var correctRP = []RulePatternTerm{
	{"product", opEQ, "jacket"},
	{"price", opGT, 50.0},
}

func testInvalidAttrName(t *testing.T) {
	ruleSets[mainRS].rules[1].rulePattern = []RulePatternTerm{
		{"product", opEQ, "jacket"},
		// priceabc is not in the schema
		{"priceabc", opGT, 50.0},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRS + "invalid attr name")
	}
	ruleSets[mainRS].rules[1].rulePattern = correctRP
}

func testTaskAsAttrName(t *testing.T) {
	ruleSets[mainRS].rules[1].rulePattern = []RulePatternTerm{
		{"product", opEQ, "jacket"},
		// freejar is not in the pattern-schema, but it is a task in the action-schema
		{"freejar", opEQ, true},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if !ok || err != nil {
		t.Errorf(incorrectOutputRS + "a task 'tag' as an attribute name")
	}
	ruleSets[mainRS].rules[1].rulePattern = correctRP
}

func testWrongAttrValType(t *testing.T) {
	ruleSets[mainRS].rules[1].rulePattern = []RulePatternTerm{
		{"product", opEQ, "jacket"},
		// price should be a float, not a string
		{"price", opGT, "abc"},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRS + "wrong attribute value type")
	}
	ruleSets[mainRS].rules[1].rulePattern = correctRP
}

func testInvalidOp(t *testing.T) {
	ruleSets[mainRS].rules[1].rulePattern = []RulePatternTerm{
		{"product", opEQ, "jacket"},
		// it should be "gt" (opGT), not "greater than"
		{"price", "greater than", 50.0},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRS + "invalid operation")
	}
	ruleSets[mainRS].rules[1].rulePattern = correctRP
}

// In each of the rule-action tests below, a rule-action is modified temporarily.
// After each test, we must reset the rule-action to the correct one below before
// moving on to the next test.
var correctRA RuleActions = RuleActions{
	tasks:      []string{"freemug", "freejar", "freeplant"},
	properties: map[string]string{"discount": "20"},
}

func testTaskNotInSchema(t *testing.T) {
	ruleSets[mainRS].rules[3].ruleActions = RuleActions{
		// freeeraser is not in the schema
		tasks:      []string{"freemug", "freeeraser"},
		properties: map[string]string{"discount": "20"},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRS + "task not in schema")
	}
	ruleSets[mainRS].rules[3].ruleActions = correctRA
}

func testPropNameNotInSchema(t *testing.T) {
	ruleSets[mainRS].rules[3].ruleActions = RuleActions{
		tasks: []string{"freemug", "freejar", "freeplant"},
		// cashback is not a property in the action-schema
		properties: map[string]string{"cashback": "5"},
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRS + "property name not in schema")
	}
	ruleSets[mainRS].rules[3].ruleActions = correctRA
}

func testBothReturnAndExit(t *testing.T) {
	ruleSets[mainRS].rules[3].ruleActions = RuleActions{
		tasks:      []string{"freemug", "freejar", "freeplant"},
		properties: map[string]string{"discount": "20"},
		// both WillReturn and WillExit below should not be true
		willReturn: true,
		willExit:   true,
	}
	ok, err := verifyRuleSet(ruleSets[mainRS], false)
	if ok || err == nil {
		t.Errorf(incorrectOutputRS + "both RETURN and EXIT instructions")
	}
	ruleSets[mainRS].rules[3].ruleActions = correctRA
}

func testCorrectWF(t *testing.T) {
	ok, err := verifyRuleSet(ruleSets[uccCreation], true)
	if !ok || err != nil {
		t.Errorf(incorrectOutputWF + "no issues")
	}
}

func testWFRuleMissingStep(t *testing.T) {
	ruleSets[uccCreation].rules[1].rulePattern = []RulePatternTerm{
		// there should be a "step" attribute here
		{stepFailed, opEQ, false},
		{"mode", opEQ, "physical"},
	}
	ok, err := verifyRuleSet(ruleSets[uccCreation], true)
	if ok || err == nil {
		t.Errorf(incorrectOutputWF + "a rule missing 'step'")
	}
	// Reset to original correct rule-pattern
	ruleSets[uccCreation].rules[1].rulePattern = []RulePatternTerm{
		{step, opEQ, "getcustdetails"},
		{stepFailed, opEQ, false},
		{"mode", opEQ, "physical"},
	}
}

// In each of the (workflow) rule-action tests below, a rule-action is modified temporarily.
// After each test, we must reset the rule-action to the correct one below before
// moving on to the next test.
var correctWorkflowRA RuleActions = RuleActions{
	tasks:      []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
	properties: map[string]string{nextStep: "aof"},
}

func testWFRuleMissingBothNSAndDone(t *testing.T) {
	ruleSets[uccCreation].rules[1].ruleActions = RuleActions{
		tasks: []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
		// Properties below should contain at least one of "nextstep" and "done"
		properties: map[string]string{},
	}
	ok, err := verifyRuleSet(ruleSets[uccCreation], true)
	if ok || err == nil {
		t.Errorf(incorrectOutputWF + "a rule missing both 'nextstep' and 'done'")
	}
	ruleSets[uccCreation].rules[1].ruleActions = correctWorkflowRA
}

func testWFNoTasksAndNotDone(t *testing.T) {
	ruleSets[uccCreation].rules[1].ruleActions = RuleActions{
		// Either Tasks below should not be empty, or Properties below should contain {"done", "true"}
		tasks:      []string{},
		properties: map[string]string{nextStep: "abc"},
	}
	ok, err := verifyRuleSet(ruleSets[uccCreation], true)
	if ok || err == nil {
		t.Errorf(incorrectOutputWF + "a rule with no tasks and no 'done=true'")
	}
	ruleSets[uccCreation].rules[1].ruleActions = correctWorkflowRA
}

func testWFNextStepValNotInTasks(t *testing.T) {
	ruleSets[uccCreation].rules[1].ruleActions = RuleActions{
		tasks: []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
		// "abcd" below is not in "Tasks" above
		properties: map[string]string{nextStep: "abcd"},
	}
	ok, err := verifyRuleSet(ruleSets[uccCreation], true)
	if ok || err == nil {
		t.Errorf(incorrectOutputWF + "a 'nextstep' value not in its rule's 'tasks'")
	}
	ruleSets[uccCreation].rules[1].ruleActions = correctWorkflowRA
}

func TestDoReferentialChecks(t *testing.T) {
	testNoReferentialIssues(t)
	testWrongThenCall(t)
	testWrongElseCall(t)
}

func testNoReferentialIssues(t *testing.T) {
	setupRuleSetsForTransaction()
	ok, err := doReferentialChecks()
	if !ok || err != nil {
		t.Errorf("unexpected output when there are no referential issues")
	}
}

func testWrongThenCall(t *testing.T) {
	// there is no ruleset "summerdisc"
	ruleSets[mainRS].rules[0].ruleActions.thenCall = "summerdisc"
	ok, err := doReferentialChecks()
	if ok || err == nil {
		t.Errorf("unexpected output when there is an incorrect THENCALL")
	}
	// reset the THENCALL back to "winterdisc"
	ruleSets[mainRS].rules[0].ruleActions.thenCall = "winterdisc"
}

func testWrongElseCall(t *testing.T) {
	// there is no ruleset "normaldisc"
	ruleSets[mainRS].rules[0].ruleActions.elseCall = "normaldisc"
	ok, err := doReferentialChecks()
	if ok || err == nil {
		t.Errorf("unexpected output when there is an incorrect ELSECALL")
	}
	// reset the ELSECALL back to "regulardisc"
	ruleSets[mainRS].rules[0].ruleActions.elseCall = "regulardisc"
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
