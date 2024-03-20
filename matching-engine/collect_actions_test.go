/* This file contains tests for collectActions() */

package crux

import (
	"reflect"
	"testing"
)

func TestCollectActionsBasic(t *testing.T) {
	actionSet := ActionSet{
		Tasks:      []string{"dodiscount", "yearendsale"},
		Properties: map[string]string{"discount": "7", "shipby": "fedex"},
	}

	ruleActions := RuleActionBlock_t{
		Task:       []string{"yearendsale", "summersale"},
		Properties: map[string]string{"cashback": "10", "discount": "9"},
		ThenCall:   "domesticpo",
		DoReturn:   false,
		DoExit:     true,
	}

	want := ActionSet{
		Tasks:      []string{"dodiscount", "yearendsale", "summersale"},
		Properties: map[string]string{"discount": "9", "shipby": "fedex", "cashback": "10"},
	}

	res := collectActions(actionSet, ruleActions)
	if !reflect.DeepEqual(want, res) {
		t.Errorf("\n\ncollectActions() = %v, \n\nwant %v\n\n", res, want)
	}
}

func TestCollectActionsWithEmptyRuleActions(t *testing.T) {
	actionSet := ActionSet{
		Tasks:      []string{"dodiscount", "yearendsale"},
		Properties: map[string]string{"discount": "7", "shipby": "fedex"},
	}

	ruleActions := RuleActionBlock_t{}

	want := ActionSet{
		Tasks:      []string{"dodiscount", "yearendsale"},
		Properties: map[string]string{"discount": "7", "shipby": "fedex"},
	}

	res := collectActions(actionSet, ruleActions)
	if !reflect.DeepEqual(want, res) {
		t.Errorf("\n\ncollectActions() = %v, \n\nwant %v\n\n", res, want)
	}
}

func TestCollectActionsWithEmptyActionSet(t *testing.T) {
	actionSet := ActionSet{
		Tasks:      []string{},
		Properties: make(map[string]string),
	}

	ruleActions := RuleActionBlock_t{
		Task:       []string{"dodiscount", "yearendsale"},
		Properties: map[string]string{"discount": "7", "shipby": "fedex"},
		ThenCall:   "overseaspo",
		DoReturn:   true,
		DoExit:     false,
	}

	want := ActionSet{
		Tasks:      []string{"dodiscount", "yearendsale"},
		Properties: map[string]string{"discount": "7", "shipby": "fedex"},
	}

	res := collectActions(actionSet, ruleActions)
	if !reflect.DeepEqual(want, res) {
		t.Errorf("\n\ncollectActions() = %v, \n\nwant %v\n\n", res, want)
	}
}
