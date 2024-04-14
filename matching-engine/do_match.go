/*
This file contains doMatch() and a helper function called by doMatch().

It also contains ruleSchemasCache and ruleSetsCache, two  variables that currently store rule-schemas
and rulesets respectively for the purpose of testing doMatch().
*/

package crux

import (
	"fmt"
	"time"
)

// var (
// 	trace trace_t
// )

func DoMatch(entity Entity, ruleset *Ruleset_t, ruleSchemasCache *Schema_t, actionSet ActionSet, seenRuleSets map[string]struct{}, trace Trace_t) (ActionSet, bool, error, Trace_t) {

	if trace.TraceData == nil {
		// time_s := time.Now()
		// this `time` is used for testing purposee revert with above in dev
		time_s := time.Date(2024, time.April, 10, 23, 0, 0, 0, time.UTC)
		trace = create_trace_level_one(time_s, entity.Realm, entity.App, ruleset.SetName, int(ruleset.Id), []TraceData_t{})
	}

	traceData := TraceData_t{
		SetID:   int(ruleset.Id),
		SetName: ruleset.SetName,
	}

	for ruleNumber, rule := range ruleset.Rules {
		var (
			eachRule TraceDataRule_t
			l2data   TraceDataRuleL2_t
		)
		eachRule.RuleNo = ruleNumber
		l2data.Tasks = rule.RuleActions.Task
		l2data.Properties = rule.RuleActions.Properties

		if len(l2data.ActionSet.Tasks) == 0 {
			l2data.ActionSet.Tasks = append(l2data.ActionSet.Tasks, rule.RuleActions.Task...)
		}
		if len(l2data.ActionSet.Properties) == 0 {
			l2data.ActionSet.Properties = rule.RuleActions.Properties
		} else {
			for k, v := range rule.RuleActions.Properties {
				l2data.ActionSet.Properties[k] = v
			}
		}

		//*********************************************************************************************
		DoExit := false

		matched, err := matchPattern(entity, rule.RulePatterns, actionSet, ruleSchemasCache)

		if err != nil {
			return ActionSet{
				Tasks:      []string{},
				Properties: make(map[string]string),
			}, false, err, trace
		}
		//*********************************************************************************************
		// if not match then set mismatch
		eachRule.Res = Mismatch
		//*********************************************************************************************
		if matched {
			//*********************************************************************************************
			// if match then set Match
			eachRule.Res = Match
			//*********************************************************************************************
			actionSet = collectActions(actionSet, rule.RuleActions)

			if len(rule.RuleActions.ThenCall) > 0 {

				// setToCall, exists := findRefRuleSetByName(ruleset, rule.RuleActions.ThenCall)

				// if !exists {
				// 	return ActionSet{}, false, errors.New("set not found")
				// }

				if ruleset.Class != entity.Class {
					return inconsistentRuleSet(ruleset.SetName, ruleset.SetName, ruleset, trace)
				}

				var err error
				actionSet, DoExit, err, trace = DoMatch(entity, ruleset, ruleSchemasCache, actionSet, seenRuleSets, trace)
				if err != nil {
					return ActionSet{
						Tasks:      []string{},
						Properties: make(map[string]string),
					}, false, err, trace
				}

			} else if DoExit || rule.RuleActions.DoExit {

				// eachRule.add_L2data(actionSet)
				// traceData.Rules = append(traceData.Rules, eachRule)
				// trace.add_tracedata(&traceData)

				return actionSet, true, nil, trace
			} else if rule.RuleActions.DoReturn {

				delete(seenRuleSets, ruleset.SetName)
				return actionSet, false, nil, trace
			} else if len(rule.RuleActions.ElseCall) > 0 {

				// setToCall, exists := findRefRuleSetByName(ruleSetsCache, rule.RuleActions.ElseCall)
				// if !exists {
				// 	return ActionSet{}, false, errors.New("set not found")
				// }

				if ruleset.Class != entity.Class {
					return inconsistentRuleSet(ruleset.SetName, ruleset.SetName, ruleset, trace)
				}

				var err error
				actionSet, DoExit, err, trace = DoMatch(entity, ruleset, ruleSchemasCache, actionSet, seenRuleSets, trace)
				if err != nil {
					return ActionSet{
						Tasks:      []string{},
						Properties: make(map[string]string),
					}, false, err, trace
				} else if DoExit {
					return actionSet, true, nil, trace
				}
			}
			//return actionSet, false, nil
		}
		//*********************************************************************************************
		eachRule.L2Data = l2data
		traceData.Rules = append(traceData.Rules, eachRule)
		//*********************************************************************************************
	}
	trace.add_tracedata(&traceData)

	delete(seenRuleSets, ruleset.SetName)

	//*********************************************************************************************
	// trace.End = time.Now()
	// this `time` is used for testing purposee revert with above in dev
	trace.End = time.Date(2024, time.April, 10, 23, 0, 0, 0, time.UTC)
	//*********************************************************************************************

	return actionSet, true, nil, trace
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

func inconsistentRuleSet(calledSetName string, currSetName string, ruleSets *Ruleset_t, trace Trace_t) (ActionSet, bool, error, Trace_t) {
	return ActionSet{}, false, fmt.Errorf("cannot call ruleset %v from ruleset %v",
		calledSetName, currSetName,
	), trace
}

func create_trace_level_one(start_t time.Time, realm, app, entryRulesetName string, entryRulesetID int, tracedata []TraceData_t) Trace_t {
	traced_data := Trace_t{
		Start:            start_t,
		End:              time.Time{},
		Realm:            realm,
		App:              app,
		EntryRulesetID:   entryRulesetID,
		EntryRulesetName: entryRulesetName,
		TraceData:        &tracedata,
	}
	return traced_data
}

func (traced_data *Trace_t) add_tracedata(record_to_add *TraceData_t) {
	*traced_data.TraceData = append(*traced_data.TraceData, *record_to_add)
}

// func (traced_data_rules *TraceData_t) add_rule(rule_to_add *TraceDataRule_t) {
// 	traced_data_rules.Rules = append(traced_data_rules.Rules, *rule_to_add)
// }
