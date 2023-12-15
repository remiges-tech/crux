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
			properties: []Property{{nextStep, "getcustdetails"}},
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
			properties: []Property{{nextStep, "aof"}},
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
			properties: []Property{{nextStep, "aof"}},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "getcustdetails"},
			{stepFailed, opEQ, true},
		},
		RuleActions{
			tasks:      []string{},
			properties: []Property{{done, trueStr}},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "aof"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"sendauthlinktoclient"},
			properties: []Property{{nextStep, "sendauthlinktoclient"}},
		},
	}
	rule6 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "aof"},
			{stepFailed, opEQ, true},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
		},
	}
	rule7 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "sendauthlinktoclient"},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
		},
	}
	ruleSets["ucccreation"] = RuleSet{1, uccCreationClass, "ucccreation",
		[]Rule{rule1, rule2, rule3, rule4, rule5, rule6, rule7},
	}
}

func testUCCStart(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, start},
		{"mode", "demat"},
	}}
	want := ActionSet{
		tasks:      []string{"getcustdetails"},
		properties: []Property{{nextStep, "getcustdetails"}},
	}
	*tests = append(*tests, doMatchTest{"ucc start", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCGetCustDetailsDemat(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "getcustdetails"},
		{stepFailed, falseStr},
		{"mode", "demat"},
	}}
	want := ActionSet{
		tasks:      []string{"aof", "kycvalid", "nomauth", "dpandbankaccvalid"},
		properties: []Property{{nextStep, "aof"}},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails demat", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCGetCustDetailsDematFail(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "getcustdetails"},
		{stepFailed, trueStr},
		{"mode", "demat"},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails demat fail", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCGetCustDetailsPhysical(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "getcustdetails"},
		{stepFailed, falseStr},
		{"mode", "physical"},
	}}
	want := ActionSet{
		tasks:      []string{"aof", "kycvalid", "nomauth", "bankaccvalid"},
		properties: []Property{{nextStep, "aof"}},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails physical", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCGetCustDetailsPhysicalFail(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "getcustdetails"},
		{stepFailed, trueStr},
		{"mode", "physical"},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"ucc getcustdetails physical fail", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCAOF(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "aof"},
		{stepFailed, falseStr},
		{"mode", "demat"},
	}}
	want := ActionSet{
		tasks:      []string{"sendauthlinktoclient"},
		properties: []Property{{nextStep, "sendauthlinktoclient"}},
	}
	*tests = append(*tests, doMatchTest{"ucc aof", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCAOFFail(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "aof"},
		{stepFailed, trueStr},
		{"mode", "demat"},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"ucc aof fail", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCEndSuccess(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "sendauthlinktoclient"},
		{stepFailed, falseStr},
		{"mode", "demat"},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"ucc end-success", entity, ruleSets["ucccreation"], ActionSet{}, want})
}

func testUCCEndFailure(tests *[]doMatchTest) {
	entity := Entity{uccCreationClass, []Attr{
		{step, "sendauthlinktoclient"},
		{stepFailed, trueStr},
		{"mode", "demat"},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"ucc end-failure", entity, ruleSets["ucccreation"], ActionSet{}, want})
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
	entity := Entity{prepareAOFClass, []Attr{
		{step, start},
	}}
	want := ActionSet{
		tasks:      []string{"downloadform"},
		properties: []Property{{nextStep, "downloadform"}},
	}
	*tests = append(*tests, doMatchTest{"download aof", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func testDownloadAOFFail(tests *[]doMatchTest) {
	entity := Entity{prepareAOFClass, []Attr{
		{step, "downloadform"},
		{stepFailed, trueStr},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"download aof fail", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func testPrintAOF(tests *[]doMatchTest) {
	entity := Entity{prepareAOFClass, []Attr{
		{step, "downloadform"},
		{stepFailed, falseStr},
	}}
	want := ActionSet{
		tasks:      []string{"printprefilledform"},
		properties: []Property{{nextStep, "printprefilledform"}},
	}
	*tests = append(*tests, doMatchTest{"print prefilled aof", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func testSignAOF(tests *[]doMatchTest) {
	entity := Entity{prepareAOFClass, []Attr{
		{step, "printprefilledform"},
		{stepFailed, falseStr},
	}}
	want := ActionSet{
		tasks:      []string{"signform"},
		properties: []Property{{nextStep, "signform"}},
	}
	*tests = append(*tests, doMatchTest{"sign aof", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func testSignAOFFail(tests *[]doMatchTest) {
	entity := Entity{prepareAOFClass, []Attr{
		{step, "signform"},
		{stepFailed, trueStr},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"sign aof fail", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func testReceiveSignedAOF(tests *[]doMatchTest) {
	entity := Entity{prepareAOFClass, []Attr{
		{step, "signform"},
		{stepFailed, falseStr},
	}}
	want := ActionSet{
		tasks:      []string{"receivesignedform"},
		properties: []Property{{nextStep, "receivesignedform"}},
	}
	*tests = append(*tests, doMatchTest{"receive signed aof", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func testUploadAOF(tests *[]doMatchTest) {
	entity := Entity{prepareAOFClass, []Attr{
		{step, "receivesignedform"},
		{stepFailed, falseStr},
	}}
	want := ActionSet{
		tasks:      []string{"uploadsignedform"},
		properties: []Property{{nextStep, "uploadsignedform"}},
	}
	*tests = append(*tests, doMatchTest{"upload signed aof", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func testPrepareAOFEnd(tests *[]doMatchTest) {
	entity := Entity{prepareAOFClass, []Attr{
		{step, "uploadsignedform"},
		{stepFailed, falseStr},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"prepare aof end", entity, ruleSets["prepareaof"], ActionSet{}, want})
}

func setupRuleSetForPrepareAOF() {
	rule1 := Rule{
		[]RulePatternTerm{
			{step, opEQ, start},
		},
		RuleActions{
			tasks:      []string{"downloadform"},
			properties: []Property{{nextStep, "downloadform"}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "downloadform"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"printprefilledform"},
			properties: []Property{{nextStep, "printprefilledform"}},
		},
	}
	rule2F := Rule{
		[]RulePatternTerm{
			{step, opEQ, "downloadform"},
			{stepFailed, opEQ, true},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "printprefilledform"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"signform"},
			properties: []Property{{nextStep, "signform"}},
		},
	}
	rule3F := Rule{
		[]RulePatternTerm{
			{step, opEQ, "printprefilledform"},
			{stepFailed, opEQ, true},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "signform"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"receivesignedform"},
			properties: []Property{{nextStep, "receivesignedform"}},
		},
	}
	rule4F := Rule{
		[]RulePatternTerm{
			{step, opEQ, "signform"},
			{stepFailed, opEQ, true},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "receivesignedform"},
			{stepFailed, opEQ, false},
		},
		RuleActions{
			tasks:      []string{"uploadsignedform"},
			properties: []Property{{nextStep, "uploadsignedform"}},
		},
	}
	rule5F := Rule{
		[]RulePatternTerm{
			{step, opEQ, "receivesignedform"},
			{stepFailed, opEQ, true},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
		},
	}
	rule6 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "uploadsignedform"},
		},
		RuleActions{
			tasks:      []string{},
			properties: []Property{{done, trueStr}},
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
	entity := Entity{validateAOFClass, []Attr{
		{step, start},
		{"aofexists", trueStr},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"validate existing aof", entity, ruleSets["validateaof"], ActionSet{}, want})
}

func testValidateAOFStart(tests *[]doMatchTest) {
	entity := Entity{validateAOFClass, []Attr{
		{step, start},
		{"aofexists", falseStr},
	}}
	want := ActionSet{
		tasks:      []string{"sendaoftorta"},
		properties: []Property{{nextStep, "sendaoftorta"}},
	}
	*tests = append(*tests, doMatchTest{"send aof to rta", entity, ruleSets["validateaof"], ActionSet{}, want})
}

func testSendAOFToRTAFail(tests *[]doMatchTest) {
	entity := Entity{validateAOFClass, []Attr{
		{step, "sendaoftorta"},
		{stepFailed, trueStr},
		{"aofexists", falseStr},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"send aof to rta fail", entity, ruleSets["validateaof"], ActionSet{}, want})
}

func testAOFGetResponseFromRTA(tests *[]doMatchTest) {
	entity := Entity{validateAOFClass, []Attr{
		{step, "sendaoftorta"},
		{stepFailed, falseStr},
		{"aofexists", falseStr},
	}}
	want := ActionSet{
		tasks:      []string{"getresponsefromrta"},
		properties: []Property{{nextStep, "getresponsefromrta"}},
	}
	*tests = append(*tests, doMatchTest{"aof - get response from rta", entity, ruleSets["validateaof"], ActionSet{}, want})
}

func testValidateAOFEnd(tests *[]doMatchTest) {
	entity := Entity{validateAOFClass, []Attr{
		{step, "getresponsefromrta"},
		{stepFailed, falseStr},
		{"aofexists", falseStr},
	}}
	want := ActionSet{
		properties: []Property{{done, trueStr}},
	}
	*tests = append(*tests, doMatchTest{"validate aof end", entity, ruleSets["validateaof"], ActionSet{}, want})
}

func setupRuleSetForValidateAOF() {
	rule1 := Rule{
		[]RulePatternTerm{
			{step, opEQ, start},
			{"aofexists", opEQ, true},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{step, opEQ, start},
			{"aofexists", opEQ, false},
		},
		RuleActions{
			tasks:      []string{"sendaoftorta"},
			properties: []Property{{nextStep, "sendaoftorta"}},
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
			properties: []Property{{nextStep, "getresponsefromrta"}},
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
			properties: []Property{{done, trueStr}},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{step, opEQ, "getresponsefromrta"},
			{"aofexists", opEQ, false},
		},
		RuleActions{
			properties: []Property{{done, trueStr}},
		},
	}
	ruleSets["validateaof"] = RuleSet{1, validateAOFClass, "validateaof",
		[]Rule{rule1, rule2, rule3, rule3F, rule4},
	}
}
