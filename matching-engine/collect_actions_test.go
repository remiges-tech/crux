/* This file contains tests for collectActions() */

package crux

import (
	"reflect"
	"testing"
)

func TestCollectActionsBasic(t *testing.T) {
	actionSet := ActionSet{
		tasks:      []string{"dodiscount", "yearendsale"},
		properties: map[string]string{"discount": "7", "shipby": "fedex"},
	}

	ruleActions := ruleActionBlock_t{
		Task:       []string{"yearendsale", "summersale"},
		Properties: map[string]string{"cashback": "10", "discount": "9"},
		ThenCall:   "domesticpo",
		DoReturn:   false,
		DoExit:     true,
	}

	want := ActionSet{
		tasks:      []string{"dodiscount", "yearendsale", "summersale"},
		properties: map[string]string{"discount": "9", "shipby": "fedex", "cashback": "10"},
	}

	res := collectActions(actionSet, ruleActions)
	if !reflect.DeepEqual(want, res) {
		t.Errorf("\n\ncollectActions() = %v, \n\nwant %v\n\n", res, want)
	}
}

func TestCollectActionsWithEmptyRuleActions(t *testing.T) {
	actionSet := ActionSet{
		tasks:      []string{"dodiscount", "yearendsale"},
		properties: map[string]string{"discount": "7", "shipby": "fedex"},
	}

	ruleActions := ruleActionBlock_t{}

	want := ActionSet{
		tasks:      []string{"dodiscount", "yearendsale"},
		properties: map[string]string{"discount": "7", "shipby": "fedex"},
	}

	res := collectActions(actionSet, ruleActions)
	if !reflect.DeepEqual(want, res) {
		t.Errorf("\n\ncollectActions() = %v, \n\nwant %v\n\n", res, want)
	}
}

func TestCollectActionsWithEmptyActionSet(t *testing.T) {
	actionSet := ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}

	ruleActions := ruleActionBlock_t{
		Task:       []string{"dodiscount", "yearendsale"},
		Properties: map[string]string{"discount": "7", "shipby": "fedex"},
		ThenCall:   "overseaspo",
		DoReturn:   true,
		DoExit:     false,
	}

	want := ActionSet{
		tasks:      []string{"dodiscount", "yearendsale"},
		properties: map[string]string{"discount": "7", "shipby": "fedex"},
	}

	res := collectActions(actionSet, ruleActions)
	if !reflect.DeepEqual(want, res) {
		t.Errorf("\n\ncollectActions() = %v, \n\nwant %v\n\n", res, want)
	}
}
