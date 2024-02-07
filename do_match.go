/*
This file contains doMatch() and a helper function called by doMatch().

It also contains ruleSchemas and ruleSets, two global variables that currently store rule-schemas
and rulesets respectively for the purpose of testing doMatch().
*/

package main

import (
	"errors"
	"fmt"
)

var ruleSchemas = []RuleSchema{}
var ruleSets = map[string]RuleSet{}

func doMatch(entity Entity, ruleSet RuleSet, actionSet ActionSet, seenRuleSets map[string]bool) (ActionSet, bool, error) {
	if seenRuleSets[ruleSet.setName] {
		return ActionSet{
			tasks:      []string{},
			properties: make(map[string]string),
		}, false, errors.New("ruleset has already been traversed")
	}
	seenRuleSets[ruleSet.setName] = true
	for _, rule := range ruleSet.rules {
		willExit := false
		matched, err := matchPattern(entity, rule.rulePattern, actionSet)
		if err != nil {
			return ActionSet{
				tasks:      []string{},
				properties: make(map[string]string),
			}, false, err
		}
		if matched {
			actionSet = collectActions(actionSet, rule.ruleActions)
			if len(rule.ruleActions.thenCall) > 0 {
				setToCall := ruleSets[rule.ruleActions.thenCall]
				if setToCall.class != entity.class {
					return inconsistentRuleSet(setToCall.setName, ruleSet.setName)
				}
				var err error
				actionSet, willExit, err = doMatch(entity, setToCall, actionSet, seenRuleSets)
				if err != nil {
					return ActionSet{
						tasks:      []string{},
						properties: make(map[string]string),
					}, false, err
				}
			}
			if willExit || rule.ruleActions.willExit {
				return actionSet, true, nil
			}
			if rule.ruleActions.willReturn {
				delete(seenRuleSets, ruleSet.setName)
				return actionSet, false, nil
			}
		} else if len(rule.ruleActions.elseCall) > 0 {
			setToCall := ruleSets[rule.ruleActions.elseCall]
			if setToCall.class != entity.class {
				return inconsistentRuleSet(setToCall.setName, ruleSet.setName)
			}
			var err error
			actionSet, willExit, err = doMatch(entity, setToCall, actionSet, seenRuleSets)
			if err != nil {
				return ActionSet{
					tasks:      []string{},
					properties: make(map[string]string),
				}, false, err
			} else if willExit {
				return actionSet, true, nil
			}
		}
	}
	delete(seenRuleSets, ruleSet.setName)
	return actionSet, false, nil
}

func inconsistentRuleSet(calledSetName string, currSetName string) (ActionSet, bool, error) {
	return ActionSet{}, false, fmt.Errorf("cannot call ruleset %v of class %v from ruleset %v of class %v",
		calledSetName, ruleSets[calledSetName].class, currSetName, ruleSets[currSetName].class,
	)
}
