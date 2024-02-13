package wfinstance

/*
This file contains doMatch() and a helper function called by doMatch().

It also contains ruleSchemas and ruleSets, two global variables that currently store rule-schemas
and rulesets respectively for the purpose of testing doMatch().
*/

import (
	"fmt"
)

type RuleSchema struct {
	class         string
	patternSchema []AttrSchema
	actionSchema  ActionSchema
}
type ActionSchema struct {
	tasks      []string
	properties []string
}
type AttrSchema struct {
	name    string
	valType string
	vals    map[string]bool
	valMin  float64
	valMax  float64
	lenMin  int
	lenMax  int
}

type RuleSet struct {
	ver     int
	class   string
	setName string
	rules   []Rule
}

type Rule struct {
	rulePattern []RulePatternTerm
	ruleActions RuleActions
}

type RulePatternTerm struct {
	attrName string
	op       string
	attrVal  any
}

type RuleActions struct {
	tasks      []string
	properties []Property
	thenCall   string
	elseCall   string
	willReturn bool
	willExit   bool
}
type ActionSet struct {
	tasks      []string
	properties map[string]string
}

type Property struct {
	name string
	val  string
}

var ruleSchemas = []RuleSchema{}
var ruleSets = map[string]RuleSet{}

func doMatch(entity Entity, ruleSet RuleSet, actionSet ActionSet, seenRuleSets map[string]bool) (ActionSet, bool, error) {

	// Initializing the ActionSet struct

	//if task has only one task
	actionSet = ActionSet{
		tasks:      []string{"discount"},
		properties: map[string]string{"nextstep": "coupondistribution"},
	}

	//if task has multiple tasks
	// actionSet = ActionSet{
	// 	tasks:      []string{"discount", "yearendsale"},
	// 	properties: map[string]string{"nextstep": "coupondistribution"},
	// }

	// if task has only one task but no proper attri
	// actionSet = ActionSet{
	// 	tasks:      []string{"discount"}, //, "yearendsale"},
	// 	properties: map[string]string{"shipby": "fedex"},
	// }

	//if done attr  present
	// actionSet = ActionSet{
	// 	tasks:      []string{"discount", "yearendsale"},
	// 	properties: map[string]string{"done": "true"},
	// }

	return actionSet, true, nil

}

func inconsistentRuleSet(calledSetName string, currSetName string) (ActionSet, bool, error) {
	return ActionSet{}, false, fmt.Errorf("cannot call ruleset %v of class %v from ruleset %v of class %v",
		calledSetName, ruleSets[calledSetName].class, currSetName, ruleSets[currSetName].class,
	)
}
