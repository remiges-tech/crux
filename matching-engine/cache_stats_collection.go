package crux

import (
	"errors"
	"time"
)

type rulesetStats_t map[Realm_t]statsPerRealm_t

type statsPerRealm_t map[App_t]statsPerApp_t

type statsPerApp_t map[Slice_t]statsPerSlice_t

type statsPerSlice_t struct {
	loadedAt   time.Time
	BRSchema   map[ClassName_t]statsSchema_t
	BRRulesets map[ClassName_t][]statsRuleset_t
	WFSchema   map[ClassName_t]statsSchema_t
	Workflows  map[ClassName_t][]statsRuleset_t
}

type statsSchema_t struct {
	nChecked int32
}

type statsRuleset_t struct {
	nCalled    int32
	rulesStats []statsRule_t
}

type statsRule_t struct {
	nMatched int32
	nFailed  int32
}

func getStats(realm Realm_t, app App_t, slice Slice_t) (rulesetStats_t, time.Time, error) {
	var timestamp time.Time
	var err error

	statsData := make(rulesetStats_t)

	perRealm, realmExists := RulesetCache[realm]
	if !realmExists {
		return nil, timestamp, errors.New("Realm not found")
	}
	perApp, appExists := perRealm[app]
	if !appExists {
		return nil, timestamp, errors.New("App not found")
	}
	perSlice, sliceExists := perApp[slice]
	if !sliceExists {
		return nil, timestamp, errors.New("Slice not found")
	}
	timestamp = perSlice.LoadedAt
	statsPerSlice := statsPerSlice_t{
		loadedAt:   perSlice.LoadedAt,
		BRSchema:   make(map[ClassName_t]statsSchema_t),
		BRRulesets: make(map[ClassName_t][]statsRuleset_t),
		WFSchema:   make(map[ClassName_t]statsSchema_t),
		Workflows:  make(map[ClassName_t][]statsRuleset_t),
	}

	if schemaData, exists := SchemaCache[realm][app][slice]; exists {
		for className, schemas := range schemaData.BRSchema {
			for _, schema := range schemas {
				statsPerSlice.BRSchema[className] = statsSchema_t{nChecked: schema.NChecked}

			}
		}
		for className, schemas := range schemaData.WFSchema {
			for _, schema := range schemas {
				statsPerSlice.WFSchema[className] = statsSchema_t{nChecked: schema.NChecked}

			}
		}
	}
	if rulesetData, exists := RulesetCache[realm][app][slice]; exists {
		for className, rulesets := range rulesetData.BRRulesets {
			statsPerSlice.BRRulesets[className] = make([]statsRuleset_t, len(rulesets))
			for i, ruleset := range rulesets {
				statsPerSlice.BRRulesets[className][i] = statsRuleset_t{
					nCalled:    ruleset.NCalled,
					rulesStats: make([]statsRule_t, len(ruleset.Rules)),
				}
				for j, rule := range ruleset.Rules {
					statsPerSlice.BRRulesets[className][i].rulesStats[j] = statsRule_t{
						nMatched: rule.NMatched,
						nFailed:  rule.NFailed,
					}
				}
			}
		}
		for className, rulesets := range rulesetData.Workflows {
			statsPerSlice.Workflows[className] = make([]statsRuleset_t, len(rulesets))
			for i, ruleset := range rulesets {
				statsPerSlice.Workflows[className][i] = statsRuleset_t{
					nCalled:    ruleset.NCalled,
					rulesStats: make([]statsRule_t, len(ruleset.Rules)),
				}
				for j, rule := range ruleset.Rules {
					statsPerSlice.Workflows[className][i].rulesStats[j] = statsRule_t{
						nMatched: rule.NMatched,
						nFailed:  rule.NFailed,
					}
				}
			}
		}
	}

	statsPerApp := make(statsPerApp_t)
	statsPerApp[slice] = statsPerSlice

	statsPerRealm := make(statsPerRealm_t)
	statsPerRealm[app] = statsPerApp

	statsData[realm] = statsPerRealm

	return statsData, timestamp, err
}
