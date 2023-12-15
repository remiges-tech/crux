/* This file contains tests for collectActions() */

package main

import (
	"reflect"
	"testing"
)

func TestCollectActionsBasic(t *testing.T) {
	actionSet := ActionSet{
		tasks:      []string{"dodiscount", "yearendsale"},
		properties: []Property{{"discount", "7"}, {"shipby", "fedex"}},
	}

	ruleActions := RuleActions{
		tasks:      []string{"yearendsale", "summersale"},
		properties: []Property{{"cashback", "10"}, {"discount", "9"}},
		thenCall:   "domesticpo",
		willReturn: false,
		willExit:   true,
	}

	want := ActionSet{
		tasks:      []string{"dodiscount", "yearendsale", "summersale"},
		properties: []Property{{"discount", "9"}, {"shipby", "fedex"}, {"cashback", "10"}},
	}

	res := collectActions(actionSet, ruleActions)
	if !reflect.DeepEqual(want, res) {
		t.Errorf("\n\ncollectActions() = %v, \n\nwant %v\n\n", res, want)
	}
}

func TestCollectActionsWithEmptyRuleActions(t *testing.T) {
	actionSet := ActionSet{
		tasks:      []string{"dodiscount", "yearendsale"},
		properties: []Property{{"discount", "7"}, {"shipby", "fedex"}},
	}

	ruleActions := RuleActions{}

	want := ActionSet{
		tasks:      []string{"dodiscount", "yearendsale"},
		properties: []Property{{"discount", "7"}, {"shipby", "fedex"}},
	}

	res := collectActions(actionSet, ruleActions)
	if !reflect.DeepEqual(want, res) {
		t.Errorf("\n\ncollectActions() = %v, \n\nwant %v\n\n", res, want)
	}
}

func TestCollectActionsWithEmptyActionSet(t *testing.T) {
	actionSet := ActionSet{}

	ruleActions := RuleActions{
		tasks:      []string{"dodiscount", "yearendsale"},
		properties: []Property{{"discount", "7"}, {"shipby", "fedex"}},
		thenCall:   "overseaspo",
		willReturn: true,
		willExit:   false,
	}

	want := ActionSet{
		tasks:      []string{"dodiscount", "yearendsale"},
		properties: []Property{{"discount", "7"}, {"shipby", "fedex"}},
	}

	res := collectActions(actionSet, ruleActions)
	if !reflect.DeepEqual(want, res) {
		t.Errorf("\n\ncollectActions() = %v, \n\nwant %v\n\n", res, want)
	}
}
