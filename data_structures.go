/*
This file contains the data structures used by the matching engine
*/

package main

type RuleSchema struct {
	class         string
	patternSchema []AttrSchema
	actionSchema  ActionSchema
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

type ActionSchema struct {
	tasks      []string
	properties []string
}

type Entity struct {
	realm string
	app   string
	slice string
	class string
	attrs map[string]string
}

type ActionSet struct {
	tasks      []string
	properties map[string]string
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
	properties map[string]string
	thenCall   string
	elseCall   string
	willReturn bool
	willExit   bool
}
