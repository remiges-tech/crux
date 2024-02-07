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
	ruleSchemas = append(ruleSchemas, RuleSchema{
		class: uccCreationClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeEnum},
			{name: stepFailed, valType: typeBool},
			{name: "mode", valType: typeEnum},
		},
		actionSchema: ActionSchema{
			tasks: []string{"getcustdetails", "aof", "kycvalid", "nomauth", "bankaccvalid",
				"dpandbankaccvalid", "sendauthlinktoclient"},
			properties: []string{nextStep, done},
		},
	})
}

func setupUCCCreationRuleSet() {
	rule1 := Rule{
		[]RulePatternTerm{
			{step, opEQ, start},
		},
		RuleActions{
			tasks:      []string{"getcustdetails"},
			properties: map[string]string{nextStep: "getcustdetails"},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "getcustdetails"},
			{stepFailed, opEQ, false},
			{"mode", opEQ, "physical"},
		},
		RuleActions{
			tasks:      []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
			properties: map[string]string{nextStep: "aof"},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "getcustdetails"},
			{stepFailed, opEQ, false},
			{"mode", opEQ, "demat"},
		},
		RuleActions{
			tasks:      []string{"aof", "kycvalid", "nomauth", "dpandbankaccvalid"},
			properties: map[string]string{nextStep: "aof"},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "getcustdetails"},
			{stepFailed, opEQ, true},
		},
		RuleActions{
			tasks:      []string{},
			properties: map[string]string{done: trueStr},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "aof"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"sendauthlinktoclient"},
			properties: map[string]string{nextStep: "sendauthlinktoclient"},
		},
	}
	rule6 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "aof"},
			{stepFailed, opEQ, true},
		},
		RuleActions{
			properties: map[string]string{done: trueStr},
		},
	}
	rule7 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "sendauthlinktoclient"},
		},
		RuleActions{
			properties: map[string]string{done: trueStr},
		},
	}
	ruleSets["ucccreation"] = RuleSet{1, uccCreationClass, "ucccreation",
		[]Rule{rule1, rule2, rule3, rule4, rule5, rule6, rule7},
	}
}

func testUCCStart(tests *[]doMatchTest) {
	entity := Entity{
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
	*tests = append(*tests, doMatchTest{"ucc start", entity, ruleSets["ucccreation"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCGetCustDetailsDemat(tests *[]doMatchTest) {
	entity := Entity{
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "getcustdetails",
			stepFailed: falseStr,
			"mode":     "demat",
		},
	}

	want := ActionSet{
		tasks:      []string{"aof", "kycvalid", "nomauth", "dpandbankaccvalid"},
		properties: map[string]string{nextStep: "aof"},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails demat", entity, ruleSets["ucccreation"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCGetCustDetailsDematFail(tests *[]doMatchTest) {
	entity := Entity{
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "getcustdetails",
			stepFailed: trueStr,
			"mode":     "demat",
		},
	}

	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails demat fail", entity, ruleSets["ucccreation"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCGetCustDetailsPhysical(tests *[]doMatchTest) {
	entity := Entity{
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "getcustdetails",
			stepFailed: falseStr,
			"mode":     "physical",
		},
	}

	want := ActionSet{
		tasks:      []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
		properties: map[string]string{nextStep: "aof"},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails physical", entity, ruleSets["ucccreation"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCGetCustDetailsPhysicalFail(tests *[]doMatchTest) {
	entity := Entity{
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "getcustdetails",
			stepFailed: trueStr,
			"mode":     "physical",
		},
	}

	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails physical fail", entity, ruleSets["ucccreation"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCAOF(tests *[]doMatchTest) {
	entity := Entity{
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "aof",
			stepFailed: falseStr,
			"mode":     "demat",
		},
	}

	want := ActionSet{
		tasks:      []string{"sendauthlinktoclient"},
		properties: map[string]string{nextStep: "sendauthlinktoclient"},
	}
	*tests = append(*tests, doMatchTest{"ucc aof", entity, ruleSets["ucccreation"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCAOFFail(tests *[]doMatchTest) {
	entity := Entity{
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "aof",
			stepFailed: trueStr,
			"mode":     "demat",
		},
	}

	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"ucc aof fail", entity, ruleSets["ucccreation"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCEndSuccess(tests *[]doMatchTest) {
	entity := Entity{
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "sendauthlinktoclient",
			stepFailed: falseStr,
			"mode":     "demat",
		},
	}

	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"ucc end-success", entity, ruleSets["ucccreation"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUCCEndFailure(tests *[]doMatchTest) {
	entity := Entity{
		class: uccCreationClass,
		attrs: map[string]string{
			step:       "sendauthlinktoclient",
			stepFailed: trueStr,
			"mode":     "demat",
		},
	}

	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"ucc end-failure", entity, ruleSets["ucccreation"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testPrepareAOF(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		class: prepareAOFClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeEnum},
			{name: stepFailed, valType: typeBool},
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
		class: prepareAOFClass,
		attrs: map[string]string{
			step: start,
		},
	}

	want := ActionSet{
		tasks:      []string{"downloadform"},
		properties: map[string]string{nextStep: "downloadform"},
	}
	*tests = append(*tests, doMatchTest{"download aof", entity, ruleSets["prepareaof"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testDownloadAOFFail(tests *[]doMatchTest) {
	entity := Entity{
		class: prepareAOFClass,
		attrs: map[string]string{
			step:       "downloadform",
			stepFailed: trueStr,
		},
	}

	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"download aof fail", entity, ruleSets["prepareaof"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testPrintAOF(tests *[]doMatchTest) {
	entity := Entity{
		class: prepareAOFClass,
		attrs: map[string]string{
			step:       "downloadform",
			stepFailed: falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"printprefilledform"},
		properties: map[string]string{nextStep: "printprefilledform"},
	}
	*tests = append(*tests, doMatchTest{"print prefilled aof", entity, ruleSets["prepareaof"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testSignAOF(tests *[]doMatchTest) {
	entity := Entity{
		class: prepareAOFClass,
		attrs: map[string]string{
			step:       "printprefilledform",
			stepFailed: falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"signform"},
		properties: map[string]string{nextStep: "signform"},
	}
	*tests = append(*tests, doMatchTest{"sign aof", entity, ruleSets["prepareaof"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testSignAOFFail(tests *[]doMatchTest) {
	entity := Entity{
		class: prepareAOFClass,
		attrs: map[string]string{
			step:       "signform",
			stepFailed: trueStr,
		},
	}

	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"sign aof fail", entity, ruleSets["prepareaof"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testReceiveSignedAOF(tests *[]doMatchTest) {
	entity := Entity{
		class: prepareAOFClass,
		attrs: map[string]string{
			step:       "signform",
			stepFailed: falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"receivesignedform"},
		properties: map[string]string{nextStep: "receivesignedform"},
	}
	*tests = append(*tests, doMatchTest{"receive signed aof", entity, ruleSets["prepareaof"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testUploadAOF(tests *[]doMatchTest) {
	entity := Entity{
		class: prepareAOFClass,
		attrs: map[string]string{
			step:       "receivesignedform",
			stepFailed: falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"uploadsignedform"},
		properties: map[string]string{nextStep: "uploadsignedform"},
	}
	*tests = append(*tests, doMatchTest{"upload signed aof", entity, ruleSets["prepareaof"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testPrepareAOFEnd(tests *[]doMatchTest) {
	entity := Entity{
		class: prepareAOFClass,
		attrs: map[string]string{
			step:       "uploadsignedform",
			stepFailed: falseStr,
		},
	}
	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"prepare aof end", entity, ruleSets["prepareaof"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func setupRuleSetForPrepareAOF() {
	rule1 := Rule{
		[]RulePatternTerm{
			{step, opEQ, start},
		},
		RuleActions{
			tasks:      []string{"downloadform"},
			properties: map[string]string{nextStep: "downloadform"},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "downloadform"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"printprefilledform"},
			properties: map[string]string{nextStep: "printprefilledform"},
		},
	}
	rule2F := Rule{
		[]RulePatternTerm{
			{step, opEQ, "downloadform"},
			{stepFailed, opEQ, true},
		},
		RuleActions{
			properties: map[string]string{done: trueStr},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "printprefilledform"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"signform"},
			properties: map[string]string{nextStep: "signform"},
		},
	}
	rule3F := Rule{
		[]RulePatternTerm{
			{step, opEQ, "printprefilledform"},
			{stepFailed, opEQ, true},
		},
		RuleActions{
			properties: map[string]string{done: trueStr},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "signform"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"receivesignedform"},
			properties: map[string]string{nextStep: "receivesignedform"},
		},
	}
	rule4F := Rule{
		[]RulePatternTerm{
			{step, opEQ, "signform"},
			{stepFailed, opEQ, true},
		},
		RuleActions{
			properties: map[string]string{done: trueStr},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "receivesignedform"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"uploadsignedform"},
			properties: map[string]string{nextStep: "uploadsignedform"},
		},
	}
	rule5F := Rule{
		[]RulePatternTerm{
			{step, opEQ, "receivesignedform"},
			{stepFailed, opEQ, true},
		},
		RuleActions{
			properties: map[string]string{done: trueStr},
		},
	}
	rule6 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "uploadsignedform"},
		},
		RuleActions{
			tasks:      []string{},
			properties: map[string]string{done: trueStr},
		},
	}
	ruleSets["prepareaof"] = RuleSet{1, prepareAOFClass, "prepareaof",
		[]Rule{rule1, rule2, rule2F, rule3, rule3F, rule4, rule4F, rule5, rule5F, rule6},
	}
}

func testValidateAOF(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		class: validateAOFClass,
		patternSchema: []AttrSchema{
			{name: step, valType: typeEnum},
			{name: stepFailed, valType: typeBool},
			{name: "aofexists", valType: typeBool},
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
		class: validateAOFClass,
		attrs: map[string]string{
			step:        start,
			"aofexists": trueStr,
		},
	}

	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"validate existing aof", entity, ruleSets["validateaof"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testValidateAOFStart(tests *[]doMatchTest) {
	entity := Entity{
		class: validateAOFClass,
		attrs: map[string]string{
			step:        start,
			"aofexists": falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"sendaoftorta"},
		properties: map[string]string{nextStep: "sendaoftorta"},
	}
	*tests = append(*tests, doMatchTest{"send aof to rta", entity, ruleSets["validateaof"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testSendAOFToRTAFail(tests *[]doMatchTest) {
	entity := Entity{
		class: validateAOFClass,
		attrs: map[string]string{
			step:        "sendaoftorta",
			stepFailed:  trueStr,
			"aofexists": falseStr,
		},
	}

	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"send aof to rta fail", entity, ruleSets["validateaof"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testAOFGetResponseFromRTA(tests *[]doMatchTest) {
	entity := Entity{
		class: validateAOFClass,
		attrs: map[string]string{
			step:        "sendaoftorta",
			stepFailed:  falseStr,
			"aofexists": falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"getresponsefromrta"},
		properties: map[string]string{nextStep: "getresponsefromrta"},
	}
	*tests = append(*tests, doMatchTest{"aof - get response from rta", entity, ruleSets["validateaof"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testValidateAOFEnd(tests *[]doMatchTest) {
	entity := Entity{
		class: validateAOFClass,
		attrs: map[string]string{
			step:        "getresponsefromrta",
			stepFailed:  falseStr,
			"aofexists": falseStr,
		},
	}

	want := ActionSet{
		properties: map[string]string{done: trueStr},
	}
	*tests = append(*tests, doMatchTest{"validate aof end", entity, ruleSets["validateaof"], ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func setupRuleSetForValidateAOF() {
	rule1 := Rule{
		[]RulePatternTerm{
			{step, opEQ, start},
			{"aofexists", opEQ, true},
		},
		RuleActions{
			properties: map[string]string{done: trueStr},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{step, opEQ, start},
			{"aofexists", opEQ, false},
		},
		RuleActions{
			tasks:      []string{"sendaoftorta"},
			properties: map[string]string{nextStep: "sendaoftorta"},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "sendaoftorta"},
			{stepFailed, opEQ, false},
			{"aofexists", opEQ, false},
		},
		RuleActions{
			tasks:      []string{"getresponsefromrta"},
			properties: map[string]string{nextStep: "getresponsefromrta"},
		},
	}
	rule3F := Rule{
		[]RulePatternTerm{
			{step, opEQ, "sendaoftorta"},
			{stepFailed, opEQ, true},
			{"aofexists", opEQ, false},
		},
		RuleActions{
			tasks:      []string{},
			properties: map[string]string{done: trueStr},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "getresponsefromrta"},
			{"aofexists", opEQ, false},
		},
		RuleActions{
			properties: map[string]string{done: trueStr},
		},
	}
	ruleSets["validateaof"] = RuleSet{1, validateAOFClass, "validateaof",
		[]Rule{rule1, rule2, rule3, rule3F, rule4},
	}
}
