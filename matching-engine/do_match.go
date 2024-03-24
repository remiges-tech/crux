/*
This file contains doMatch() and a helper function called by doMatch().

It also contains ruleSchemasCache and ruleSetsCache, two  variables that currently store rule-schemas
and rulesets respectively for the purpose of testing doMatch().
*/

package crux

import (
	"errors"
	"fmt"
)

func DoMatch(entity Entity, ruleset *Ruleset_t, actionSet ActionSet, seenRuleSets map[string]struct{}) (ActionSet, bool, error) {

	ruleSchemasCache, ruleSetsCache := retriveRuleSchemasAndRuleSetsFromCache(entity.Realm, entity.App, entity.Class, entity.Slice)
	if _, seen := seenRuleSets[ruleset.SetName]; seen {
		return ActionSet{
			Tasks:      []string{},
			Properties: make(map[string]string),
		}, false, errors.New("ruleset has already been traversed")
	}

	seenRuleSets[ruleset.SetName] = struct{}{}

	for _, rule := range ruleset.Rules {

		DoExit := false

		matched, err := matchPattern(entity, rule.RulePatterns, actionSet, ruleSchemasCache)

		if err != nil {
			return ActionSet{
				Tasks:      []string{},
				Properties: make(map[string]string),
			}, false, err
		}
        
		if matched {

			actionSet = collectActions(actionSet, rule.RuleActions)
            
			if len(rule.RuleActions.ThenCall) > 0 {

				setToCall, exists := findRefRuleSetByName(ruleSetsCache, rule.RuleActions.ThenCall)

				if !exists {
					return ActionSet{}, false, errors.New("set not found")
				}

				if setToCall.Class != entity.Class {
					return inconsistentRuleSet(setToCall.SetName, ruleset.SetName, ruleset)
				}

				var err error
				actionSet, DoExit, err = DoMatch(entity, setToCall, actionSet, seenRuleSets)
				if err != nil {
					return ActionSet{
						Tasks:      []string{},
						Properties: make(map[string]string),
					}, false, err
				}
			}

			if DoExit || rule.RuleActions.DoExit {

				return actionSet, true, nil
			}

			if rule.RuleActions.DoReturn {

				delete(seenRuleSets, ruleset.SetName)
				return actionSet, false, nil
			}
		} else if len(rule.RuleActions.ElseCall) > 0 {

			setToCall, exists := findRefRuleSetByName(ruleSetsCache, rule.RuleActions.ElseCall)
			if !exists {
				return ActionSet{}, false, errors.New("set not found")
			}

			if setToCall.Class != entity.Class {
				return inconsistentRuleSet(setToCall.SetName, ruleset.SetName, ruleset)
			}

			var err error
			actionSet, DoExit, err = DoMatch(entity, setToCall, actionSet, seenRuleSets)
			if err != nil {
				return ActionSet{
					Tasks:      []string{},
					Properties: make(map[string]string),
				}, false, err
			} else if DoExit {
				return actionSet, true, nil
			}
		}
	}

	delete(seenRuleSets, ruleset.SetName)

	return actionSet, false, nil
}

func findRefRuleSetByName(ruleSets []*Ruleset_t, setName string) (*Ruleset_t, bool) {
	for _, ruleset := range ruleSets {
		for _, rule := range ruleset.Rules {
			found := false
			for _, referRuleset := range rule.RuleActions.References {
				if referRuleset.SetName == setName {
					rule.NMatched++
					return referRuleset, true
				}
			}
			if !found {
				rule.NFailed++
			}
		}
	}
	return nil, false
}

func inconsistentRuleSet(calledSetName string, currSetName string, ruleSets *Ruleset_t) (ActionSet, bool, error) {
	return ActionSet{}, false, fmt.Errorf("cannot call ruleset %v from ruleset %v",
		calledSetName, currSetName,
	)
}
