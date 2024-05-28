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

func DoMatch(entity Entity, ruleset *Ruleset_t, ruleSchemasCache *Schema_t, actionSet ActionSet, seenRuleSets map[string]struct{}, trace Trace_t, trace_level int, cruxCache *Cache) (ActionSet, bool, error, Trace_t) {

	if !(-1 < trace_level && trace_level < 3) {
		return ActionSet{}, false, fmt.Errorf("trace_level out of bound"), trace
	}
	is_trace_enable := trace_level > 0

	if trace.Start == nil && trace_level > 0 {
		// this `time` is used for testing purposee revert with Now() in dev
		// time_s := time.Date(2024, time.April, 10, 23, 0, 0, 0, time.UTC)
		time_s := time.Now()
		trace = create_trace_level_one(&time_s, entity.Realm, entity.App, ruleset.SetName, int(ruleset.Id), []TraceData_t{})
	}

	traceData := TraceData_t{
		SetID:   int(ruleset.Id),
		SetName: ruleset.SetName,
	}
	for ruleNumber, rule := range ruleset.Rules {
		var (
			l2data = TraceDataRuleL2_t{
				Pattern:    []map[string][]string{},
				Tasks:      []string{},
				Properties: map[string]string{},
				NextSet:    0,
				ActionSet:  &ActionSet_t{},
			}
			eachRule TraceDataRule_t
		)
		DoExit := false
		matched, err := matchPattern(entity, rule.RulePatterns, actionSet, ruleSchemasCache)
		if err != nil {
			return ActionSet{}, false, err, trace
		}
		/*
			This part contains trace data for thencall and it is same for elsecall and return call,
			in which when there is a thencall/elsecall then all rules data will be saved into the `trace` then rules
			within running process will be set to empty coz we need to show rules after saved rules when it comes back
			from thencal ruleset.
				For e.g. in a DB ruleset if there is thencall for 2nd rule then we save that data object till 2nd rule in
			tracing and it will keep adding data for given 2nd ruleset within thencall consider there are 5 entries and then
			it will come back to main entry ruleset, now from here it will start to add rules data in entry Ruleset from 3rd onward.
			refer ( https://github.com/remiges-tech/crux/wiki/Traversal-trace-structure#trace-level-1 )

			Note: Please change time.End & ...Start with `time.Date(2024, time.April, 10, 23, 0, 0, 0, time.UTC)` while running test cases
		*/
		eachRule.RuleNo = ruleNumber
		eachRule.Res = Mismatch // set mismatch as `-`
		if matched {
			eachRule.Res = Match // if match then set Match as `+`
			if trace_level > 1 {
				l2data.Tasks = rule.RuleActions.Task
				l2data.Properties = rule.RuleActions.Properties
				l2data.ActionSet.Tasks = append(l2data.ActionSet.Tasks, rule.RuleActions.Task...)
				for k, v := range rule.RuleActions.Properties {
					l2data.ActionSet.Properties = map[string]string{k: v}
				}
				if trace.TraceData != nil {
					l2data.ActionSet.Tasks = append(l2data.ActionSet.Tasks, actionSet.Tasks...)
					for k, v := range actionSet.Properties {
						l2data.ActionSet.Properties[k] = v
					}
				}
			}
			actionSet = collectActions(actionSet, rule.RuleActions)
			if len(rule.RuleActions.ThenCall) > 0 {
				newRulesetSchema, thencallRuleset, err := cruxCache.RetriveRuleSchemasAndRuleSetsFromCache(WFE, entity.App, entity.Realm, entity.Class, rule.RuleActions.ThenCall, entity.Slice)
				if err != nil {
					return ActionSet{}, false, err, trace
				}
				if is_trace_enable {
					eachRule.Res = ThenCall // if ThenCall then set Match as `+c`
					eachRule.NextSet = rule.RuleActions.ThenCall
					traceData.add_rule(&eachRule, l2data, trace_level)
					trace.add_tracedata(&traceData)
					trace.End = time.Now()
					traceData.Rules = []TraceDataRule_t{}
				}
				if ruleset.Class != entity.Class {
					return inconsistentRuleSet(ruleset.SetName, ruleset.SetName, ruleset, trace)
				}
				actionSet, DoExit, err, trace = DoMatch(entity, thencallRuleset, newRulesetSchema, actionSet, seenRuleSets, trace, trace_level, cruxCache)
				if err != nil {
					return ActionSet{}, false, err, trace
				}
			} else if DoExit || rule.RuleActions.DoExit {
				if is_trace_enable {
					eachRule.Res = Exit // if Exit then set Match as `<<`
					traceData.add_rule(&eachRule, l2data, trace_level)
					trace.add_tracedata(&traceData)
					// trace.End = time.Date(2024, time.April, 10, 23, 0, 0, 0, time.UTC)
					trace.End = time.Now()
					// traceData.Rules = []TraceDataRule_t{}
				}
				return actionSet, true, nil, trace
			} else if rule.RuleActions.DoReturn {
				delete(seenRuleSets, ruleset.SetName)
				if is_trace_enable {
					eachRule.Res = Return // if Return then set Match as `<`
					traceData.add_rule(&eachRule, l2data, trace_level)
					trace.add_tracedata(&traceData)
					// trace.End = time.Date(2024, time.April, 10, 23, 0, 0, 0, time.UTC)
					trace.End = time.Now()

				}
				return actionSet, false, nil, trace
			} else {
				if is_trace_enable {
					traceData.add_rule(&eachRule, l2data, trace_level)
				}
			}
		} else if len(rule.RuleActions.ElseCall) > 0 {
			newRulesetSchema, elsecallRuleset, err := cruxCache.RetriveRuleSchemasAndRuleSetsFromCache(WFE, entity.App, entity.Realm, entity.Class, rule.RuleActions.ElseCall, entity.Slice)
			if err != nil {
				return ActionSet{}, false, err, trace
			}
			if is_trace_enable {
				eachRule.Res = ElseCall // if ElseCall then set Match as `-c`
				eachRule.NextSet = rule.RuleActions.ElseCall
				traceData.add_rule(&eachRule, l2data, trace_level)
				trace.add_tracedata(&traceData)
				// trace.End = time.Date(2024, time.April, 10, 23, 0, 0, 0, time.UTC)
				trace.End = time.Now()
				traceData.Rules = []TraceDataRule_t{}
			}
			if ruleset.Class != entity.Class {
				return inconsistentRuleSet(ruleset.SetName, ruleset.SetName, ruleset, trace)
			}
			actionSet, DoExit, err, trace = DoMatch(entity, elsecallRuleset, newRulesetSchema, actionSet, seenRuleSets, trace, trace_level, cruxCache)
			if err != nil {
				return ActionSet{}, false, err, trace
			} else if DoExit {
				if is_trace_enable {
					eachRule.Res = Exit
					traceData.add_rule(&eachRule, l2data, trace_level)
					trace.add_tracedata(&traceData)
					// trace.End = time.Date(2024, time.April, 10, 23, 0, 0, 0, time.UTC)
					trace.End = time.Now()
				}
				return actionSet, true, nil, trace
			}
		} else {
			if is_trace_enable {
				if trace_level > 1 {
					for _, val := range rule.RulePatterns {
						l2data.Pattern = append(l2data.Pattern, map[string][]string{val.Attr: {
							val.Val.(string), val.Op, entity.Attrs[val.Attr], returnRespSymbol(val.Val.(string), entity.Attrs[val.Attr]),
						}})
					}
				}
				traceData.add_rule(&eachRule, l2data, trace_level)
				trace.add_tracedata(&traceData) // append collected tracedata to trace
				// trace.End = time.Date(2024, time.April, 10, 23, 0, 0, 0, time.UTC)
				trace.End = time.Now()
			}
			return actionSet, false, nil, trace
		}
	}
	delete(seenRuleSets, ruleset.SetName)
	if is_trace_enable && (len(traceData.Rules) > 0) {
		trace.add_tracedata(&traceData) // append collected tracedata to trace
		// set end time of trace
		// this `time` is used for testing purposee revert with Now() in dev
		// trace.End = time.Date(2024, time.April, 10, 23, 0, 0, 0, time.UTC)
		trace.End = time.Now()
	}
	return actionSet, false, nil, trace
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

func create_trace_level_one(start_t *time.Time, realm, app, entryRulesetName string, entryRulesetID int, tracedata []TraceData_t) Trace_t {
	return Trace_t{
		Start:            start_t,
		End:              time.Time{},
		Realm:            realm,
		App:              app,
		EntryRulesetID:   entryRulesetID,
		EntryRulesetName: entryRulesetName,
		TraceData:        tracedata,
	}
}

func (traced_data *Trace_t) add_tracedata(record_to_add *TraceData_t) {
	traced_data.TraceData = append(traced_data.TraceData, *record_to_add)
}

func (traced_data_rules *TraceData_t) add_rule(each_rule *TraceDataRule_t, l2_data_to_add TraceDataRuleL2_t, trace_level int) {
	if trace_level == 2 {
		each_rule.TraceDataRuleL2_t = l2_data_to_add
	}
	if trace_level >= 1 {
		traced_data_rules.Rules = append(traced_data_rules.Rules, *each_rule)
	}
}

func returnRespSymbol(rulAttr, entyAttr string) string {
	res := "-"
	if rulAttr == entyAttr {
		res = "+"
	}
	return res
}
