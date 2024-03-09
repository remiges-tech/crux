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

func doMatch(entity Entity, ruleset []*Ruleset_t, actionSet ActionSet, seenRuleSets map[string]struct{}) (ActionSet, bool, error) {
	ruleSchemasCache, ruleSetsCache := retriveRuleSchemasAndRuleSetsFromCache(entity.realm, entity.app, entity.class, entity.slice)

	for _, ruleSet := range ruleset {
		if _, seen := seenRuleSets[ruleSet.SetName]; seen {
			return ActionSet{
				tasks:      []string{},
				properties: make(map[string]string),
			}, false, errors.New("ruleset has already been traversed")
		}

		seenRuleSets[ruleSet.SetName] = struct{}{}

		for _, rule := range ruleSet.Rules {

			DoExit := false

			matched, err := matchPattern(entity, rule.RulePatterns, actionSet, ruleSchemasCache)

			if err != nil {
				return ActionSet{
					tasks:      []string{},
					properties: make(map[string]string),
				}, false, err
			}

			if matched {

				actionSet = collectActions(actionSet, rule.RuleActions)

				if len(rule.RuleActions.ThenCall) > 0 {

					setToCall, exists := findRuleSetByName(ruleSetsCache, rule.RuleActions.ThenCall)

					if !exists {
						return ActionSet{}, false, errors.New("set not found")
					}

					if setToCall.Class != entity.class {
						return inconsistentRuleSet(setToCall.SetName, ruleSet.SetName, ruleset)
					}

					var err error
					actionSet, DoExit, err = doMatch(entity, []*Ruleset_t{setToCall}, actionSet, seenRuleSets)
					if err != nil {
						return ActionSet{
							tasks:      []string{},
							properties: make(map[string]string),
						}, false, err
					}
				}

				if DoExit || rule.RuleActions.DoExit {

					return actionSet, true, nil
				}

				if rule.RuleActions.DoReturn {

					delete(seenRuleSets, ruleSet.SetName)
					return actionSet, false, nil
				}
			} else if len(rule.RuleActions.ElseCall) > 0 {

				setToCall, exists := findRuleSetByName(ruleset, rule.RuleActions.ElseCall)
				if !exists {
					return ActionSet{}, false, errors.New("set not found")
				}

				if setToCall.Class != entity.class {
					return inconsistentRuleSet(setToCall.SetName, ruleSet.SetName, ruleset)
				}

				var err error
				actionSet, DoExit, err = doMatch(entity, []*Ruleset_t{setToCall}, actionSet, seenRuleSets)
				if err != nil {
					return ActionSet{
						tasks:      []string{},
						properties: make(map[string]string),
					}, false, err
				} else if DoExit {
					return actionSet, true, nil
				}
			}
		}

		delete(seenRuleSets, ruleSet.SetName)
	}

	return actionSet, false, nil
}

func findRuleSetByName(ruleSets []*Ruleset_t, setName string) (*Ruleset_t, bool) {
	for _, ruleset := range ruleSets {
		if ruleset.SetName == setName {

			ruleset.NCalled++
			for _, rule := range ruleset.Rules {

				rule.NMatched++

			}
			return ruleset, true
		} else {

			for _, rule := range ruleset.Rules {

				rule.NFailed++

			}
		}
	}
	return nil, false
}

func inconsistentRuleSet(calledSetName string, currSetName string, ruleSets []*Ruleset_t) (ActionSet, bool, error) {
	return ActionSet{}, false, fmt.Errorf("cannot call ruleset %v from ruleset %v",
		calledSetName, currSetName,
	)
}
