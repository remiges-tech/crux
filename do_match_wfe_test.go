/*
This file contains the functions that represent WFE tests for doMatch(). These functions are called
inside TestDoMatch() in do_match_test.go.

Some of the definitions of rulesets below deliberately use a lot of whitespace to keep the code consistent
and to make it easier to understand, add to, and edit these tests
*/

package main

const (
	uccCreationClass = "ucccreation"
	prepareAOFClass  = "prepareaof"
	validateAOFClass = "validateaof"
)

func testUCCCreation(tests *[]doMatchTest) {
	setupUCCCreationSchema()
	setupUCCCreationRuleSet()

	// Each test below involves calling doMatch() with a different entity
	testUCCStart(tests)
	testUCCGetCustDetailsDemat(tests)
	testUCCGetCustDetailsDematFail(tests)
	testUCCGetCustDetailsPhysical(tests)
	testUCCGetCustDetailsPhysicalFail(tests)
	testUCCAOF(tests)
	testUCCAOFFail(tests)
	testUCCEndSuccess(tests)
	testUCCEndFailure(tests)
}

func setupUCCCreationSchema() {

	ruleSchemasTest = append(ruleSchemasTest, &schema_t{
		Class: uccCreationClass,
		PatternSchema: patternSchema_t{
			Attr: []attr_t{
				{Name: step, ValType: typeEnum},
				{Name: stepFailed, ValType: typeBool},
				{Name: "mode", ValType: typeEnum},
			},
		},
		ActionSchema: actionSchema_t{
			Tasks:      []string{"getcustdetails", "aof", "kycvalid", "nomauth", "bankaccvalid", "dpandbankaccvalid", "sendauthlinktoclient"},
			Properties: []string{nextStep, done},
		},
	})

}

func setupUCCCreationRuleSet() *Ruleset_t {
	rule1 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, start},
		},
		RuleActions: ruleActionBlock_t{
			Task:       []string{"getcustdetails"},
			Properties: map[string]string{nextStep: "getcustdetails"},
		},
	}
	rule2 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "getcustdetails"},
			{stepFailed, opEQ, false},
			{"mode", opEQ, "physical"},
		},
		RuleActions: ruleActionBlock_t{
			Task:       []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
			Properties: map[string]string{nextStep: "aof"},
		},
	}
	rule3 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "getcustdetails"},
			{stepFailed, opEQ, false},
			{"mode", opEQ, "demat"},
		},
		RuleActions: ruleActionBlock_t{
			Task:       []string{"aof", "kycvalid", "nomauth", "dpandbankaccvalid"},
			Properties: map[string]string{nextStep: "aof"},
		},
	}
	rule4 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "getcustdetails"},
			{stepFailed, opEQ, true},
		},
		RuleActions: ruleActionBlock_t{
			Task:       []string{},
			Properties: map[string]string{done: trueStr},
		},
	}
	rule5 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "aof"},
			{stepFailed, opEQ, false},
		},
		RuleActions: ruleActionBlock_t{
			Task:       []string{"sendauthlinktoclient"},
			Properties: map[string]string{nextStep: "sendauthlinktoclient"},
		},
	}
	rule6 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "aof"},
			{stepFailed, opEQ, true},
		},
		RuleActions: ruleActionBlock_t{
			Properties: map[string]string{done: trueStr},
		},
	}
	rule7 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "sendauthlinktoclient"},
		},
		RuleActions: ruleActionBlock_t{
			Properties: map[string]string{done: trueStr},
		},
	}

	rs := Ruleset_t{
		Id:      1,
		Class:   uccCreationClass,
		SetName: mainRS,
		Rules:   []rule_t{rule1, rule2, rule3, rule4, rule5, rule6, rule7, rule7},
		NCalled: 0,
	}

	return &rs
}

func testUCCStart(tests *[]doMatchTest) {

	rc := setupUCCCreationRuleSet()

	entity := Entity{
		realm: "1",
		app:   "Test10",
		slice: "10",
		class: uccCreationClass,
		attrs: map[string]string{
			step:   start,
			"mode": "demat",
		},
	}

	want := ActionSet{
		tasks:      []string{"getcustdetails"},
		properties: map[string]string{nextStep: "getcustdetails"},
	}
	*tests = append(*tests, doMatchTest{"ucc start", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCGetCustDetailsDemat(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test10",
		slice: "10",
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "getcustdetails",
			stepFailed: falseStr,
			"mode":     "demat",
		},
	}
	rc := setupUCCCreationRuleSet()
	want := ActionSet{
		tasks:      []string{"aof", "kycvalid", "nomauth", "dpandbankaccvalid"},
		properties: map[string]string{nextStep: "aof"},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails demat", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCGetCustDetailsDematFail(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test10",
		slice: "10",
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "getcustdetails",
			stepFailed: trueStr,
			"mode":     "demat",
		},
	}
	rc := setupUCCCreationRuleSet()
	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails demat fail", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCGetCustDetailsPhysical(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test10",
		slice: "10",
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "getcustdetails",
			stepFailed: falseStr,
			"mode":     "physical",
		},
	}
	rc := setupUCCCreationRuleSet()
	want := ActionSet{
		tasks:      []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
		properties: map[string]string{nextStep: "aof"},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails physical", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCGetCustDetailsPhysicalFail(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test10",
		slice: "10",
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "getcustdetails",
			stepFailed: trueStr,
			"mode":     "physical",
		},
	}
	rc := setupUCCCreationRuleSet()
	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails physical fail", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCAOF(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test10",
		slice: "10",
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "aof",
			stepFailed: falseStr,
			"mode":     "demat",
		},
	}
	rc := setupUCCCreationRuleSet()
	want := ActionSet{
		tasks:      []string{"sendauthlinktoclient"},
		properties: map[string]string{nextStep: "sendauthlinktoclient"},
	}
	*tests = append(*tests, doMatchTest{"ucc aof", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCAOFFail(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test10",
		slice: "10",
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "aof",
			stepFailed: trueStr,
			"mode":     "demat",
		},
	}
	rc := setupUCCCreationRuleSet()
	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"ucc aof fail", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCEndSuccess(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test10",
		slice: "10",
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "sendauthlinktoclient",
			stepFailed: falseStr,
			"mode":     "demat",
		},
	}
	rc := setupUCCCreationRuleSet()
	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"ucc end-success", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCEndFailure(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test10",
		slice: "10",
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "sendauthlinktoclient",
			stepFailed: trueStr,
			"mode":     "demat",
		},
	}
	rc := setupUCCCreationRuleSet()
	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"ucc end-failure", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testPrepareAOF(tests *[]doMatchTest) {

	ruleSchemasTest = append(ruleSchemasTest, &schema_t{
		Class: prepareAOFClass,
		PatternSchema: patternSchema_t{
			Attr: []attr_t{
				{Name: step, ValType: typeEnum},
				{Name: stepFailed, ValType: typeBool},
			},
		},
	})

	setupRuleSetForPrepareAOF()

	// Each test below involves calling doMatch() with a different entity
	testDownloadAOF(tests)
	testDownloadAOFFail(tests)
	testPrintAOF(tests)
	testSignAOF(tests)
	testSignAOFFail(tests)
	testReceiveSignedAOF(tests)
	testUploadAOF(tests)
	testPrepareAOFEnd(tests)
}

func testDownloadAOF(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test11",
		slice: "11",
		class: prepareAOFClass,
		attrs: map[string]string{
			step: start,
		},
	}
	rc := setupRuleSetForPrepareAOF()
	want := ActionSet{
		tasks:      []string{"downloadform"},
		properties: map[string]string{nextStep: "downloadform"},
	}
	*tests = append(*tests, doMatchTest{"download aof", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testDownloadAOFFail(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test11",
		slice: "11",
		class: prepareAOFClass,
		attrs: map[string]string{
			step:       "downloadform",
			stepFailed: trueStr,
		},
	}
	rc := setupRuleSetForPrepareAOF()
	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"download aof fail", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testPrintAOF(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test11",
		slice: "11",
		class: prepareAOFClass,
		attrs: map[string]string{
			step:       "downloadform",
			stepFailed: falseStr,
		},
	}
	rc := setupRuleSetForPrepareAOF()
	want := ActionSet{
		tasks:      []string{"printprefilledform"},
		properties: map[string]string{nextStep: "printprefilledform"},
	}
	*tests = append(*tests, doMatchTest{"print prefilled aof", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testSignAOF(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test11",
		slice: "11",
		class: prepareAOFClass,
		attrs: map[string]string{
			step:       "printprefilledform",
			stepFailed: falseStr,
		},
	}
	rc := setupRuleSetForPrepareAOF()

	want := ActionSet{
		tasks:      []string{"signform"},
		properties: map[string]string{nextStep: "signform"},
	}
	*tests = append(*tests, doMatchTest{"sign aof", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testSignAOFFail(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test11",
		slice: "11",
		class: prepareAOFClass,
		attrs: map[string]string{
			step:       "signform",
			stepFailed: trueStr,
		},
	}
	rc := setupRuleSetForPrepareAOF()
	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"sign aof fail", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testReceiveSignedAOF(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test11",
		slice: "11",
		class: prepareAOFClass,
		attrs: map[string]string{
			step:       "signform",
			stepFailed: falseStr,
		},
	}
	rc := setupRuleSetForPrepareAOF()
	want := ActionSet{
		tasks:      []string{"receivesignedform"},
		properties: map[string]string{nextStep: "receivesignedform"},
	}
	*tests = append(*tests, doMatchTest{"receive signed aof", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUploadAOF(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test11",
		slice: "11",
		class: prepareAOFClass,
		attrs: map[string]string{
			step:       "receivesignedform",
			stepFailed: falseStr,
		},
	}
	rc := setupRuleSetForPrepareAOF()
	want := ActionSet{
		tasks:      []string{"uploadsignedform"},
		properties: map[string]string{nextStep: "uploadsignedform"},
	}
	*tests = append(*tests, doMatchTest{"upload signed aof", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testPrepareAOFEnd(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test11",
		slice: "11",
		class: prepareAOFClass,
		attrs: map[string]string{
			step:       "uploadsignedform",
			stepFailed: falseStr,
		},
	}
	rc := setupRuleSetForPrepareAOF()
	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"prepare aof end", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func setupRuleSetForPrepareAOF() *Ruleset_t {
	rule1 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, start},
		},
		RuleActions: ruleActionBlock_t{
			Task:       []string{"downloadform"},
			Properties: map[string]string{nextStep: "downloadform"},
		},
	}
	rule2 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "downloadform"},
			{stepFailed, opEQ, false},
		},
		RuleActions: ruleActionBlock_t{
			Task:       []string{"printprefilledform"},
			Properties: map[string]string{nextStep: "printprefilledform"},
		},
	}
	rule2F := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "downloadform"},
			{stepFailed, opEQ, true},
		},
		RuleActions: ruleActionBlock_t{
			Properties: map[string]string{done: trueStr},
		},
	}
	rule3 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "printprefilledform"},
			{stepFailed, opEQ, false},
		},
		RuleActions: ruleActionBlock_t{
			Task:       []string{"signform"},
			Properties: map[string]string{nextStep: "signform"},
		},
	}
	rule3F := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "printprefilledform"},
			{stepFailed, opEQ, true},
		},
		RuleActions: ruleActionBlock_t{
			Properties: map[string]string{done: trueStr},
		},
	}
	rule4 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "signform"},
			{stepFailed, opEQ, false},
		},
		RuleActions: ruleActionBlock_t{
			Task:       []string{"receivesignedform"},
			Properties: map[string]string{nextStep: "receivesignedform"},
		},
	}
	rule4F := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "signform"},
			{stepFailed, opEQ, true},
		},
		RuleActions: ruleActionBlock_t{
			Properties: map[string]string{done: trueStr},
		},
	}
	rule5 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "receivesignedform"},
			{stepFailed, opEQ, false},
		},
		RuleActions: ruleActionBlock_t{
			Task:       []string{"uploadsignedform"},
			Properties: map[string]string{nextStep: "uploadsignedform"},
		},
	}
	rule5F := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "receivesignedform"},
			{stepFailed, opEQ, true},
		},
		RuleActions: ruleActionBlock_t{
			Properties: map[string]string{done: trueStr},
		},
	}
	rule6 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "uploadsignedform"},
		},
		RuleActions: ruleActionBlock_t{
			Task:       []string{},
			Properties: map[string]string{done: trueStr},
		},
	}

	rs := Ruleset_t{
		Id:      1,
		Class:   prepareAOFClass,
		SetName: mainRS,
		Rules:   []rule_t{rule1, rule2, rule2F, rule3, rule3F, rule4, rule4F, rule5, rule5F, rule6},
		NCalled: 0,
	}

	return &rs
}

func testValidateAOF(tests *[]doMatchTest) {
	ruleSchemasTest = append(ruleSchemasTest, &schema_t{
		Class: validateAOFClass,
		PatternSchema: patternSchema_t{
			Attr: []attr_t{
				{Name: step, ValType: typeEnum},
				{Name: stepFailed, ValType: typeBool},
				{Name: "aofexists", ValType: typeBool},
			},
		},
	})

	setupRuleSetForValidateAOF()

	// Each test below involves calling doMatch() with a different entity
	testValidateExistingAOF(tests)
	testValidateAOFStart(tests)
	testSendAOFToRTAFail(tests)
	testAOFGetResponseFromRTA(tests)
	testValidateAOFEnd(tests)
}

func testValidateExistingAOF(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test12",
		slice: "12",
		class: validateAOFClass,
		attrs: map[string]string{
			step:        start,
			"aofexists": trueStr,
		},
	}
	rc := setupRuleSetForValidateAOF()
	want := ActionSet{

		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"validate existing aof", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testValidateAOFStart(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test12",
		slice: "12",
		class: validateAOFClass,
		attrs: map[string]string{
			step:        start,
			"aofexists": falseStr,
		},
	}
	rc := setupRuleSetForValidateAOF()

	want := ActionSet{
		tasks:      []string{"sendaoftorta"},
		properties: map[string]string{nextStep: "sendaoftorta"},
	}
	*tests = append(*tests, doMatchTest{"send aof to rta", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testSendAOFToRTAFail(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test12",
		slice: "12",
		class: validateAOFClass,
		attrs: map[string]string{
			step:        "sendaoftorta",
			stepFailed:  trueStr,
			"aofexists": falseStr,
		},
	}
	rc := setupRuleSetForValidateAOF()
	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"send aof to rta fail", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testAOFGetResponseFromRTA(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test12",
		slice: "12",
		class: validateAOFClass,
		attrs: map[string]string{
			step:        "sendaoftorta",
			stepFailed:  falseStr,
			"aofexists": falseStr,
		},
	}
	rc := setupRuleSetForValidateAOF()
	want := ActionSet{
		tasks:      []string{"getresponsefromrta"},
		properties: map[string]string{nextStep: "getresponsefromrta"},
	}
	*tests = append(*tests, doMatchTest{"aof - get response from rta", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testValidateAOFEnd(tests *[]doMatchTest) {
	entity := Entity{
		realm: "1",
		app:   "Test12",
		slice: "12",
		class: validateAOFClass,
		attrs: map[string]string{
			step:        "getresponsefromrta",
			stepFailed:  falseStr,
			"aofexists": falseStr,
		},
	}
	rc := setupRuleSetForValidateAOF()
	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"validate aof end", entity, rc, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func setupRuleSetForValidateAOF() *Ruleset_t {
	rule1 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, start},
			{"aofexists", opEQ, true},
		},
		RuleActions: ruleActionBlock_t{
			Properties: map[string]string{done: trueStr},
		},
	}
	rule2 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, start},
			{"aofexists", opEQ, false},
		},
		RuleActions: ruleActionBlock_t{
			Task:       []string{"sendaoftorta"},
			Properties: map[string]string{nextStep: "sendaoftorta"},
		},
	}
	rule3 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "sendaoftorta"},
			{stepFailed, opEQ, false},
			{"aofexists", opEQ, false},
		},
		RuleActions: ruleActionBlock_t{
			Task:       []string{"getresponsefromrta"},
			Properties: map[string]string{nextStep: "getresponsefromrta"},
		},
	}
	rule4 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "sendaoftorta"},
			{stepFailed, opEQ, true},
			{"aofexists", opEQ, false},
		},
		RuleActions: ruleActionBlock_t{
			Task:       []string{},
			Properties: map[string]string{done: trueStr},
		},
	}
	rule5 := rule_t{
		RulePatterns: []rulePatternBlock_t{
			{step, opEQ, "getresponsefromrta"},
			{"aofexists", opEQ, false},
		},
		RuleActions: ruleActionBlock_t{
			Properties: map[string]string{done: trueStr},
		},
	}

	rs := Ruleset_t{
		Id:      1,
		Class:   validateAOFClass,
		SetName: mainRS,
		Rules:   []rule_t{rule1, rule2, rule3, rule4, rule5},
		NCalled: 0,
	}
	return &rs

}
