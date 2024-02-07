/*
This file contains the functions that represent BRE tests for doMatch(). These functions are called
inside TestDoMatch() in do_match_test.go.

Some of the definitions of rulesets below deliberately use a lot of whitespace to keep the code consistent
and to make it easier to understand, add to, and edit these tests
*/

package main

const (
	// The "main" ruleset that may contain "thenCall"s/"elseCall"s to other rulesets
	mainRS = "main"

	inventoryItemClass = "inventoryitem"
	transactionClass   = "transaction"
	purchaseClass      = "purchase"
	orderClass         = "order"
)

var sampleEntity = Entity{
	class: inventoryItemClass,
	attrs: map[string]string{
		"cat":        "textbook",
		"fullname":   "Advanced Physics",
		"ageinstock": "5",
		"mrp":        "50.80",
		"received":   "2018-06-01T15:04:05Z",
		"bulkorder":  trueStr,
	},
}

func testBasic(tests *[]doMatchTest) {
	ruleSet := RuleSet{1, inventoryItemClass, mainRS,
		[]Rule{{
			[]RulePatternTerm{{"cat", opEQ, "textbook"}},
			RuleActions{
				tasks:      []string{"yearendsale", "summersale"},
				properties: map[string]string{"cashback": "10", "discount": "9"},
			},
		}},
	}
	*tests = append(*tests, doMatchTest{
		"basic test", sampleEntity, ruleSet, ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		ActionSet{
			tasks:      []string{"yearendsale", "summersale"},
			properties: map[string]string{"cashback": "10", "discount": "9"},
		},
	})
}

func testExit(tests *[]doMatchTest) {
	rA1 := RuleActions{
		tasks:      []string{"springsale"},
		properties: map[string]string{"cashback": "15"},
	}
	rA2 := RuleActions{
		tasks:      []string{"yearendsale", "summersale"},
		properties: map[string]string{"discount": "10", "freegift": "mug"},
	}
	rA3 := RuleActions{
		tasks:      []string{"wintersale"},
		properties: map[string]string{"discount": "15"},
		willExit:   true,
	}
	rA4 := RuleActions{
		tasks: []string{"autumnsale"},
	}
	ruleSet := RuleSet{1, inventoryItemClass, mainRS, []Rule{
		{[]RulePatternTerm{{"cat", opEQ, "refbook"}}, rA1},                           // no match
		{[]RulePatternTerm{{"ageinstock", opLT, 7}, {"cat", opEQ, "textbook"}}, rA2}, // match
		{[]RulePatternTerm{{"summersale", opEQ, true}}, rA3},                         // match then exit
		{[]RulePatternTerm{{"ageinstock", opLT, 7}}, rA4},                            // ignored
	}}
	want := ActionSet{
		tasks:      []string{"yearendsale", "summersale", "wintersale"},
		properties: map[string]string{"discount": "15", "freegift": "mug"},
	}
	*tests = append(*tests, doMatchTest{"exit", sampleEntity, ruleSet, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testReturn(tests *[]doMatchTest) {
	rA1 := RuleActions{
		tasks:      []string{"yearendsale", "summersale"},
		properties: map[string]string{"discount": "10", "freegift": "mug"},
	}
	rA2 := RuleActions{
		tasks:      []string{"springsale"},
		properties: map[string]string{"discount": "15"},
		willReturn: true,
	}
	rA3 := RuleActions{
		tasks: []string{"autumnsale"},
	}
	ruleSet := RuleSet{1, inventoryItemClass, mainRS, []Rule{
		{[]RulePatternTerm{{"ageinstock", opLT, 7}, {"cat", opEQ, "textbook"}}, rA1}, // match
		{[]RulePatternTerm{{"summersale", opEQ, true}}, rA2},                         // match then return
		{[]RulePatternTerm{{"ageinstock", opLT, 7}}, rA3},                            // ignored
	}}
	want := ActionSet{
		tasks:      []string{"yearendsale", "summersale", "springsale"},
		properties: map[string]string{"discount": "15", "freegift": "mug"},
	}
	*tests = append(*tests, doMatchTest{"return", sampleEntity, ruleSet, ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}, want})
}

func testTransactions(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		class: transactionClass,
		patternSchema: []AttrSchema{
			{name: "productname", valType: typeStr},
			{name: "price", valType: typeInt},
			{name: "inwintersale", valType: typeBool},
			{name: "paymenttype", valType: typeEnum},
			{name: "ismember", valType: typeBool},
		},
	})

	setupRuleSetsForTransaction()

	// Each test below involves calling doMatch() with a different entity
	testWinterDiscJacket60(tests)
	testWinterDiscJacket40(tests)
	testWinterDiscKettle110Cash(tests)
	testWinterDiscKettle110Card(tests)
	testMemberDiscLamp60(tests)
	testMemberDiscKettle60Card(tests)
	testMemberDiscKettle60Cash(tests)
	testMemberDiscKettle110Card(tests)
	testMemberDiscKettle110Cash(tests)
	testNonMemberDiscLamp30(tests)
	testNonMemberDiscKettle70(tests)
	testNonMemberDiscKettle110Cash(tests)
	testNonMemberDiscKettle110Card(tests)
}

func setupRuleSetsForTransaction() {
	setupRuleSetMainForTransaction()
	setupRuleSetWinterDisc()
	setupRuleSetRegularDisc()
	setupRuleSetMemberDisc()
	setupRuleSetNonMemberDisc()
}

func setupRuleSetMainForTransaction() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"inwintersale", opEQ, true},
		},
		RuleActions{
			thenCall: "winterdisc",
			elseCall: "regulardisc",
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"paymenttype", opEQ, "cash"},
			{"price", opGT, 10},
		},
		RuleActions{
			tasks: []string{"freepen"},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"paymenttype", opEQ, "card"},
			{"price", opGT, 10},
		},
		RuleActions{
			tasks: []string{"freemug"},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{"freehat", opEQ, true},
		},
		RuleActions{tasks: []string{"freebag"}},
	}
	ruleSets[mainRS] = RuleSet{1, transactionClass, mainRS,
		[]Rule{rule1, rule2, rule3, rule4},
	}
}

func setupRuleSetWinterDisc() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"productname", opEQ, "jacket"},
			{"price", opGT, 50},
		},
		RuleActions{
			tasks:      []string{"freehat"},
			properties: map[string]string{"discount": "50"},
			willReturn: true,
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"price", opLT, 100},
		},
		RuleActions{
			properties: map[string]string{"discount": "40", "pointsmult": "2"},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"price", opGE, 100},
		},
		RuleActions{
			properties: map[string]string{"discount": "45", "pointsmult": "3"},
		},
	}
	ruleSets["winterdisc"] = RuleSet{1, transactionClass, "winterdisc",
		[]Rule{rule1, rule2, rule3},
	}
}

func setupRuleSetRegularDisc() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"ismember", opEQ, true},
		},
		RuleActions{
			thenCall: "memberdisc",
			elseCall: "nonmemberdisc",
		},
	}
	ruleSets["regulardisc"] = RuleSet{1, transactionClass, "regulardisc",
		[]Rule{rule1},
	}
}

func setupRuleSetMemberDisc() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"productname", opEQ, "lamp"},
			{"price", opGT, 50},
		},
		RuleActions{
			properties: map[string]string{"discount": "35", "pointsmult": "2"},
			willExit:   true,
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"price", opLT, 100},
		},
		RuleActions{
			properties: map[string]string{"discount": "20"},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"price", opGE, 100},
		},
		RuleActions{
			properties: map[string]string{"discount": "25"},
		},
	}
	ruleSets["memberdisc"] = RuleSet{1, transactionClass, "memberdisc",
		[]Rule{rule1, rule2, rule3},
	}
}

func setupRuleSetNonMemberDisc() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"price", opLT, 50},
		},
		RuleActions{
			properties: map[string]string{"discount": "5"},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"price", opGE, 50},
		},
		RuleActions{
			properties: map[string]string{"discount": "10"},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"price", opGE, 100},
		},
		RuleActions{
			properties: map[string]string{"discount": "15"},
		},
	}
	ruleSets["nonmemberdisc"] = RuleSet{1, transactionClass, "nonmemberdisc",
		[]Rule{rule1, rule2, rule3},
	}
}

func testWinterDiscJacket60(tests *[]doMatchTest) {
	entity := Entity{
		class: transactionClass,
		attrs: map[string]string{
			"productname":  "jacket",
			"price":        "60",
			"inwintersale": trueStr,
			"paymenttype":  "card",
			"ismember":     trueStr,
		},
	}
	want := ActionSet{
		tasks:      []string{"freehat", "freemug", "freebag"},
		properties: map[string]string{"discount": "50"},
	}
	*tests = append(*tests, doMatchTest{
		"winterdisc jacket 60",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testWinterDiscJacket40(tests *[]doMatchTest) {
	entity := Entity{
		class: transactionClass,
		attrs: map[string]string{
			"productname":  "jacket",
			"price":        "40",
			"inwintersale": trueStr,
			"paymenttype":  "card",
			"ismember":     trueStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freemug"},
		properties: map[string]string{"discount": "40", "pointsmult": "2"},
	}
	*tests = append(*tests, doMatchTest{
		"winterdisc jacket 40",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testWinterDiscKettle110Cash(tests *[]doMatchTest) {
	entity := Entity{
		class: transactionClass,
		attrs: map[string]string{
			"productname":  "kettle",
			"price":        "110",
			"inwintersale": trueStr,
			"paymenttype":  "cash",
			"ismember":     trueStr,
		},
	}
	want := ActionSet{
		tasks:      []string{"freepen"},
		properties: map[string]string{"discount": "45", "pointsmult": "3"},
	}
	*tests = append(*tests, doMatchTest{
		"winterdisc kettle 110 cash",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testWinterDiscKettle110Card(tests *[]doMatchTest) {
	entity := Entity{
		class: transactionClass,
		attrs: map[string]string{
			"productname":  "kettle",
			"price":        "110",
			"inwintersale": trueStr,
			"paymenttype":  "card",
			"ismember":     trueStr,
		},
	}
	want := ActionSet{
		tasks:      []string{"freemug"},
		properties: map[string]string{"discount": "45", "pointsmult": "3"},
	}
	*tests = append(*tests, doMatchTest{
		"winterdisc kettle 110 card",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testMemberDiscLamp60(tests *[]doMatchTest) {
	entity := Entity{
		class: transactionClass,
		attrs: map[string]string{
			"productname":  "lamp",
			"price":        "60",
			"inwintersale": falseStr,
			"paymenttype":  "card",
			"ismember":     trueStr,
		},
	}
	want := ActionSet{
		properties: map[string]string{"discount": "35", "pointsmult": "2"},
	}
	*tests = append(*tests, doMatchTest{
		"memberdisc lamp 60",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
	//fmt.Println("tests", tests)
}

func testMemberDiscKettle60Card(tests *[]doMatchTest) {
	entity := Entity{
		class: transactionClass,
		attrs: map[string]string{
			"productname":  "kettle",
			"price":        "60",
			"inwintersale": falseStr,
			"paymenttype":  "card",
			"ismember":     trueStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freemug"},
		properties: map[string]string{"discount": "20"},
	}
	*tests = append(*tests, doMatchTest{
		"memberdisc kettle 60 card",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testMemberDiscKettle60Cash(tests *[]doMatchTest) {
	entity := Entity{
		class: transactionClass,
		attrs: map[string]string{
			"productname":  "kettle",
			"price":        "60",
			"inwintersale": falseStr,
			"paymenttype":  "cash",
			"ismember":     trueStr,
		},
	}
	want := ActionSet{
		tasks:      []string{"freepen"},
		properties: map[string]string{"discount": "20"},
	}
	*tests = append(*tests, doMatchTest{
		"memberdisc kettle 60 cash",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testMemberDiscKettle110Card(tests *[]doMatchTest) {
	entity := Entity{
		class: transactionClass,
		attrs: map[string]string{
			"productname":  "kettle",
			"price":        "110",
			"inwintersale": falseStr,
			"paymenttype":  "card",
			"ismember":     trueStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freemug"},
		properties: map[string]string{"discount": "25"},
	}
	*tests = append(*tests, doMatchTest{
		"memberdisc kettle 110 card",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testMemberDiscKettle110Cash(tests *[]doMatchTest) {
	entity := Entity{
		class: transactionClass,
		attrs: map[string]string{
			"productname":  "kettle",
			"price":        "110",
			"inwintersale": falseStr,
			"paymenttype":  "cash",
			"ismember":     trueStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freepen"},
		properties: map[string]string{"discount": "25"},
	}
	*tests = append(*tests, doMatchTest{
		"memberdisc kettle 110 cash",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testNonMemberDiscLamp30(tests *[]doMatchTest) {
	entity := Entity{
		class: transactionClass,
		attrs: map[string]string{
			"productname":  "lamp",
			"price":        "30",
			"inwintersale": falseStr,
			"paymenttype":  "cash",
			"ismember":     falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freepen"},
		properties: map[string]string{"discount": "5"},
	}
	*tests = append(*tests, doMatchTest{
		"nonmemberdisc lamp 30",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testNonMemberDiscKettle70(tests *[]doMatchTest) {
	entity := Entity{
		class: transactionClass,
		attrs: map[string]string{
			"productname":  "kettle",
			"price":        "70",
			"inwintersale": falseStr,
			"paymenttype":  "cash",
			"ismember":     falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freepen"},
		properties: map[string]string{"discount": "10"},
	}
	*tests = append(*tests, doMatchTest{
		"nonmemberdisc kettle 70",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testNonMemberDiscKettle110Cash(tests *[]doMatchTest) {
	entity := Entity{
		class: transactionClass,
		attrs: map[string]string{
			"productname":  "kettle",
			"price":        "110",
			"inwintersale": falseStr,
			"paymenttype":  "cash",
			"ismember":     falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freepen"},
		properties: map[string]string{"discount": "15"},
	}
	*tests = append(*tests, doMatchTest{
		"nonmemberdisc kettle 110 cash",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testNonMemberDiscKettle110Card(tests *[]doMatchTest) {
	entity := Entity{
		class: transactionClass,
		attrs: map[string]string{
			"productname":  "kettle",
			"price":        "110",
			"inwintersale": falseStr,
			"paymenttype":  "card",
			"ismember":     falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freemug"},
		properties: map[string]string{"discount": "15"},
	}
	*tests = append(*tests, doMatchTest{
		"nonmemberdisc kettle 110 card",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testPurchases(tests *[]doMatchTest) {
	setupPurchaseRuleSchema()
	setupRuleSetForPurchases()

	// Each test below involves calling doMatch() with a different entity
	testJacket35(tests)
	testJacket55ForMember(tests)
	testJacket55ForNonMember(tests)
	testJacket75ForMember(tests)
	testJacket75ForNonMember(tests)
	testLamp35(tests)
	testLamp55(tests)
	testLamp75ForMember(tests)
	testLamp75ForNonMember(tests)
	testKettle35(tests)
	testKettle55(tests)
	testKettle75ForMember(tests)
	testKettle75ForNonMember(tests)
	testOven35(tests)
	testOven55(tests)
}

func setupPurchaseRuleSchema() {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		class: purchaseClass,
		patternSchema: []AttrSchema{
			{name: "product", valType: typeStr},
			{name: "price", valType: typeFloat},
			{name: "ismember", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks: []string{"freepen", "freebottle", "freepencil", "freemug", "freejar", "freeplant",
				"freebag", "freenotebook"},
			properties: []string{"discount", "pointsmult"},
		},
	})
}

func testJacket35(tests *[]doMatchTest) {
	entity := Entity{
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "jacket",
			"price":    "35",
			"ismember": trueStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freepen", "freebottle", "freepencil"},
		properties: map[string]string{"discount": "5"},
	}
	*tests = append(*tests, doMatchTest{
		"jacket price 35",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testJacket55ForMember(tests *[]doMatchTest) {
	entity := Entity{
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "jacket",
			"price":    "55",
			"ismember": trueStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freepen", "freebottle", "freepencil", "freenotebook"},
		properties: map[string]string{"discount": "10"},
	}
	*tests = append(*tests, doMatchTest{
		"jacket price 55 for member",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testJacket55ForNonMember(tests *[]doMatchTest) {
	entity := Entity{
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "jacket",
			"price":    "55",
			"ismember": falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freepen", "freebottle", "freepencil", "freenotebook"},
		properties: map[string]string{"discount": "10"},
	}
	*tests = append(*tests, doMatchTest{
		"jacket price 55 for non-member",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testJacket75ForMember(tests *[]doMatchTest) {
	entity := Entity{
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "jacket",
			"price":    "75",
			"ismember": trueStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freepen", "freebottle", "freepencil", "freenotebook"},
		properties: map[string]string{"discount": "15", "pointsmult": "2"},
	}
	*tests = append(*tests, doMatchTest{
		"jacket price 75 for member",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testJacket75ForNonMember(tests *[]doMatchTest) {
	entity := Entity{
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "jacket",
			"price":    "75",
			"ismember": falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freepen", "freebottle", "freepencil", "freenotebook"},
		properties: map[string]string{"discount": "10"},
	}
	*tests = append(*tests, doMatchTest{
		"jacket price 75 for non-member",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testLamp35(tests *[]doMatchTest) {
	entity := Entity{
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "lamp",
			"price":    "35",
			"ismember": falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freemug", "freejar", "freeplant", "freebag"},
		properties: map[string]string{"discount": "20"},
	}
	*tests = append(*tests, doMatchTest{
		"lamp price 35",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testLamp55(tests *[]doMatchTest) {
	entity := Entity{
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "lamp",
			"price":    "55",
			"ismember": falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freemug", "freejar", "freeplant", "freebag", "freenotebook"},
		properties: map[string]string{"discount": "25"},
	}
	*tests = append(*tests, doMatchTest{
		"lamp price 55",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testLamp75ForMember(tests *[]doMatchTest) {
	entity := Entity{
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "lamp",
			"price":    "75",
			"ismember": trueStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freemug", "freejar", "freeplant"},
		properties: map[string]string{"discount": "30", "pointsmult": "3"},
	}
	*tests = append(*tests, doMatchTest{
		"lamp price 75 for member",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testLamp75ForNonMember(tests *[]doMatchTest) {
	entity := Entity{
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "lamp",
			"price":    "75",
			"ismember": falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freemug", "freejar", "freeplant", "freebag", "freenotebook"},
		properties: map[string]string{"discount": "25"},
	}
	*tests = append(*tests, doMatchTest{
		"lamp price 75 for non-member",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testKettle35(tests *[]doMatchTest) {
	entity := Entity{
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "kettle",
			"price":    "35",
			"ismember": trueStr,
		},
	}

	want := ActionSet{
		properties: map[string]string{"discount": "35"},
	}
	*tests = append(*tests, doMatchTest{
		"kettle price 35",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testKettle55(tests *[]doMatchTest) {
	entity := Entity{
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "kettle",
			"price":    "55",
			"ismember": falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freenotebook"},
		properties: map[string]string{"discount": "40"},
	}
	*tests = append(*tests, doMatchTest{
		"kettle price 55",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testKettle75ForMember(tests *[]doMatchTest) {
	entity := Entity{
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "kettle",
			"price":    "75",
			"ismember": trueStr,
		},
	}

	want := ActionSet{
		properties: map[string]string{"discount": "45", "pointsmult": "4"},
	}
	*tests = append(*tests, doMatchTest{
		"kettle price 75 for member",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testKettle75ForNonMember(tests *[]doMatchTest) {
	entity := Entity{
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "kettle",
			"price":    "75",
			"ismember": falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"freenotebook"},
		properties: map[string]string{"discount": "40"},
	}
	*tests = append(*tests, doMatchTest{
		"kettle price 75 for non-member",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testOven35(tests *[]doMatchTest) {
	entity := Entity{
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "oven",
			"price":    "35",
			"ismember": falseStr,
		},
	}

	want := ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}
	*tests = append(*tests, doMatchTest{
		"oven price 35",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testOven55(tests *[]doMatchTest) {
	entity := Entity{
		class: purchaseClass,
		attrs: map[string]string{
			"product":  "oven",
			"price":    "55",
			"ismember": trueStr,
		},
	}

	want := ActionSet{
		tasks: []string{"freenotebook"},
	}
	*tests = append(*tests, doMatchTest{
		"oven price 55",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func setupRuleSetForPurchases() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "jacket"},
			{"price", opGT, 30.0},
		},
		RuleActions{
			tasks:      []string{"freepen", "freebottle", "freepencil"},
			properties: map[string]string{"discount": "5"},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "jacket"},
			{"price", opGT, 50.0},
		},
		RuleActions{
			properties: map[string]string{"discount": "10"},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "jacket"},
			{"price", opGT, 70.0},
			{"ismember", opEQ, true},
		},
		RuleActions{
			properties: map[string]string{"discount": "15", "pointsmult": "2"},
		},
	}
	rule4 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "lamp"},
			{"price", opGT, 30.0},
		},
		RuleActions{
			tasks:      []string{"freemug", "freejar", "freeplant"},
			properties: map[string]string{"discount": "20"},
		},
	}
	rule5 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "lamp"},
			{"price", opGT, 50.0},
		},
		RuleActions{
			properties: map[string]string{"discount": "25"},
		},
	}
	rule6 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "lamp"},
			{"price", opGT, 70.0},
			{"ismember", opEQ, true},
		},
		RuleActions{
			properties: map[string]string{"discount": "30", "pointsmult": "3"},
			willExit:   true,
		},
	}
	rule7 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "kettle"},
			{"price", opGT, 30.0},
		},
		RuleActions{
			properties: map[string]string{"discount": "35"},
		},
	}
	rule8 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "kettle"},
			{"price", opGT, 50.0},
		},
		RuleActions{
			properties: map[string]string{"discount": "40"},
		},
	}
	rule9 := Rule{
		[]RulePatternTerm{
			{"product", opEQ, "kettle"},
			{"price", opGT, 70.0},
			{"ismember", opEQ, true},
		},
		RuleActions{
			properties: map[string]string{"discount": "45", "pointsmult": "4"},
			willReturn: true,
		},
	}
	rule10 := Rule{
		[]RulePatternTerm{
			{"freemug", opEQ, true},
		},
		RuleActions{
			tasks: []string{"freebag"},
		},
	}
	rule11 := Rule{
		[]RulePatternTerm{
			{"price", opGT, 50.0},
		},
		RuleActions{
			tasks: []string{"freenotebook"},
		},
	}
	ruleSets[mainRS] = RuleSet{1, purchaseClass, mainRS,
		[]Rule{rule1, rule2, rule3, rule4, rule5, rule6, rule7, rule8, rule9, rule10, rule11},
	}
}

func testOrders(tests *[]doMatchTest) {
	ruleSchemas = append(ruleSchemas, RuleSchema{
		class: orderClass,
		patternSchema: []AttrSchema{
			{name: "ordertype", valType: typeEnum},
			{name: "mode", valType: typeEnum},
			{name: "liquidscheme", valType: typeBool},
			{name: "overnightscheme", valType: typeBool},
			{name: "extendedhours", valType: typeBool},
		},
		actionSchema: ActionSchema{
			tasks:      []string{"unitstoamc", "unitstorta"},
			properties: []string{"amfiordercutoff", "bseordercutoff", "fundscutoff", "unitscutoff"},
		},
	})

	setupRuleSetMainForOrder()
	setupRuleSetPurchaseOrSIPForOrder()
	setupRuleSetOtherOrderTypesForOrder()

	// Each test below involves calling doMatch() with a different entity
	testSIPOrder(tests)
	testSwitchDematOrder(tests)
	testSwitchDematExtHours(tests)
	testRedemptionDematExtHours(tests)
	testPurchaseOvernightOrder(tests)
	testSIPLiquidOrder(tests)
	testSwitchPhysicalOrder(tests)
}

func setupRuleSetMainForOrder() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"ordertype", opEQ, "purchase"},
		},
		RuleActions{
			thenCall: "purchaseorsip",
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"ordertype", opEQ, "sip"},
		},
		RuleActions{
			thenCall: "purchaseorsip",
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"ordertype", opNE, "purchase"},
			{"ordertype", opNE, "sip"},
		},
		RuleActions{
			properties: map[string]string{"amfiordercutoff": "1500", "bseordercutoff": "1500"},
			thenCall:   "otherordertypes",
		},
	}
	ruleSets[mainRS] = RuleSet{1, orderClass, mainRS,
		[]Rule{rule1, rule2, rule3},
	}
}

func setupRuleSetPurchaseOrSIPForOrder() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"liquidscheme", opEQ, false},
			{"overnightscheme", opEQ, false},
		},
		RuleActions{
			properties: map[string]string{"amfiordercutoff": "1500", "bseordercutoff": "1430",
				"fundscutoff": "1430"},
			willReturn: true,
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{},
		RuleActions{
			properties: map[string]string{"amfiordercutoff": "1330", "bseordercutoff": "1300",
				"fundscutoff": "1230"},
		},
	}
	ruleSets["purchaseorsip"] = RuleSet{1, orderClass, "purchaseorsip",
		[]Rule{rule1, rule2},
	}
}

func setupRuleSetOtherOrderTypesForOrder() {
	rule1 := Rule{
		[]RulePatternTerm{
			{"mode", opEQ, "physical"},
		},
		RuleActions{
			tasks: []string{"unitstoamc", "unitstorta"},
		},
	}
	rule2 := Rule{
		[]RulePatternTerm{
			{"mode", opEQ, "demat"},
			{"extendedhours", opEQ, false},
		},
		RuleActions{
			properties: map[string]string{"unitscutoff": "1630"},
		},
	}
	rule3 := Rule{
		[]RulePatternTerm{
			{"mode", opEQ, "demat"},
			{"extendedhours", opEQ, true},
		},
		RuleActions{
			properties: map[string]string{"unitscutoff": "1730"},
		},
	}
	ruleSets["otherordertypes"] = RuleSet{1, orderClass, "otherordertypes",
		[]Rule{rule1, rule2, rule3},
	}
}

func testSIPOrder(tests *[]doMatchTest) {
	entity := Entity{
		class: orderClass,
		attrs: map[string]string{
			"ordertype":       "sip",
			"mode":            "demat",
			"liquidscheme":    falseStr,
			"overnightscheme": falseStr,
			"extendedhours":   falseStr,
		},
	}

	want := ActionSet{
		properties: map[string]string{"amfiordercutoff": "1500", "bseordercutoff": "1430",
			"fundscutoff": "1430"},
	}
	*tests = append(*tests, doMatchTest{
		"sip order",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testSwitchDematOrder(tests *[]doMatchTest) {
	entity := Entity{
		class: orderClass,
		attrs: map[string]string{
			"ordertype":       "switch",
			"mode":            "demat",
			"liquidscheme":    falseStr,
			"overnightscheme": falseStr,
			"extendedhours":   falseStr,
		},
	}

	want := ActionSet{
		properties: map[string]string{"amfiordercutoff": "1500", "bseordercutoff": "1500",
			"unitscutoff": "1630"},
	}
	*tests = append(*tests, doMatchTest{
		"switch demat order",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testSwitchDematExtHours(tests *[]doMatchTest) {
	entity := Entity{
		class: orderClass,
		attrs: map[string]string{
			"ordertype":       "switch",
			"mode":            "demat",
			"liquidscheme":    falseStr,
			"overnightscheme": falseStr,
			"extendedhours":   trueStr,
		},
	}

	want := ActionSet{
		properties: map[string]string{"amfiordercutoff": "1500", "bseordercutoff": "1500",
			"unitscutoff": "1730"},
	}
	*tests = append(*tests, doMatchTest{
		"switch demat ext-hours order",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testRedemptionDematExtHours(tests *[]doMatchTest) {
	entity := Entity{
		class: orderClass,
		attrs: map[string]string{
			"ordertype":       "redemption",
			"mode":            "demat",
			"liquidscheme":    falseStr,
			"overnightscheme": falseStr,
			"extendedhours":   trueStr,
		},
	}

	want := ActionSet{
		properties: map[string]string{"amfiordercutoff": "1500", "bseordercutoff": "1500",
			"unitscutoff": "1730"},
	}
	*tests = append(*tests, doMatchTest{
		"redemption demat ext-hours order",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testPurchaseOvernightOrder(tests *[]doMatchTest) {
	entity := Entity{
		class: orderClass,
		attrs: map[string]string{
			"ordertype":       "purchase",
			"mode":            "physical",
			"liquidscheme":    falseStr,
			"overnightscheme": trueStr,
			"extendedhours":   falseStr,
		},
	}

	want := ActionSet{
		properties: map[string]string{"amfiordercutoff": "1330", "bseordercutoff": "1300",
			"fundscutoff": "1230"},
	}
	*tests = append(*tests, doMatchTest{
		"purchase overnight order",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testSIPLiquidOrder(tests *[]doMatchTest) {
	entity := Entity{
		class: orderClass,
		attrs: map[string]string{
			"ordertype":       "sip",
			"mode":            "physical",
			"liquidscheme":    trueStr,
			"overnightscheme": falseStr,
			"extendedhours":   falseStr,
		},
	}

	want := ActionSet{
		properties: map[string]string{"amfiordercutoff": "1330", "bseordercutoff": "1300",
			"fundscutoff": "1230"},
	}
	*tests = append(*tests, doMatchTest{
		"sip liquid order",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}

func testSwitchPhysicalOrder(tests *[]doMatchTest) {
	entity := Entity{
		class: orderClass,
		attrs: map[string]string{
			"ordertype":       "switch",
			"mode":            "physical",
			"liquidscheme":    falseStr,
			"overnightscheme": trueStr,
			"extendedhours":   trueStr,
		},
	}

	want := ActionSet{
		tasks:      []string{"unitstoamc", "unitstorta"},
		properties: map[string]string{"amfiordercutoff": "1500", "bseordercutoff": "1500"},
	}
	*tests = append(*tests, doMatchTest{
		"switch physical order",
		entity,
		ruleSets[mainRS],
		ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		},
		want,
	})
}
