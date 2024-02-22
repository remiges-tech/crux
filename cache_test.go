/*
This file contains the functions that represent Cache tests for Load()/Purge()/Reload(). These functions are called
inside TestCache()) in do_matchest.go.

Some of the definitions of rulesets below deliberately use a lot of whitespace to keep the code consistent
and to make it easier to understand, add to, and edit these tests
*/

package main

import (
	"encoding/json"
	"testing"
)

func testCache(tests *[]doMatchTest, t *testing.T) {

	//testLoadDB(tests)
	testLoad(tests, t)
	testPurge(tests, t)
	testReload(tests, t)
}

func testLoadDB(tests *[]doMatchTest, t *testing.T) {

	err := Load()
	if err != nil {
		t.Errorf("Error:%+v", err)
	}
}

func setSchemaRulesetCacheBuffer(t *testing.T) {

	err := loadInternal(mockSchemasets, mockRulesets)
	if err != nil {
		t.Errorf(" %v", err)
		return
	}

}

func testLoad(tests *[]doMatchTest, t *testing.T) {

	var RealmId = 1
	var App = "test3"
	var Class = "inventoryMaterial"
	var Slice = 3
	var Brwf = "W"

	setSchemaRulesetCacheBuffer(t)

	p, a, err := retrieveSchemasFromCache(RealmId, App, Class, Slice, Brwf)

	if err != "success" {
		t.Errorf("%v", err)
		return
	}
	//UT 1 check the valid fields

	if !containsField(p, "material", t) {
		t.Errorf("Expected fieldnot found in the actualpattern")
	}

	if !containsField(a, "schedule", t) {
		t.Errorf("Expected field not found in the  actualAction")
	}

	//UT 2 rulesets check the  valid fields
	RealmId = 1
	App = "Test2"
	Slice = 2
	Class = "inventoryNewyear"
	Brwf = "B"

	rp, ra, rval, refrenceRuleset := retrieveRulesetFromCache(RealmId, App, Class, Slice, Brwf)

	if rval != "success" {
		t.Errorf("%v", rval)
		return
	}

	if !containsField(rp, "notebook", t) {
		t.Errorf("Expected field not found in the actualpattern")
	}

	if !containsField(ra, "newyearsale", t) {
		t.Errorf("Expected field not found in the  actualAction")
	}

	//UT 3 rulesets check refrecnce to other rulesets ThenCall
	for _, reference := range refrenceRuleset {

		if reference.ReferenceType == "thencall" {
			jsonData, err := json.Marshal(reference.RulePatterns)
			if err != nil {
				t.Errorf("Error:%+v", err)
				return
			}

			if !containsField(jsonData, "textbook", t) {
				t.Errorf("Expected field not found in the  RulePatterns")
			}
			jsonData, err = json.Marshal(reference.RuleActions)
			if err != nil {
				t.Errorf("Error:%+v", err)
				return
			}
			if !containsField(jsonData, "fedex", t) {
				t.Errorf("Expected field not found in the RuleActions")
			}
		}
	}
	//UT 4 rulesets check refrecnce to other rulesets ElseCall
	RealmId = 1
	App = "Test3"
	Slice = 3
	Class = "inventoryClearance"
	Brwf = "B"
	_, _, rval, refrenceRuleset = retrieveRulesetFromCache(RealmId, App, Class, Slice, Brwf)

	if rval != "success" {
		t.Errorf("%v", rval)
		return
	}
	for _, reference := range refrenceRuleset {

		if reference.ReferenceType == "elsecall" {
			jsonData, err := json.Marshal(reference.RulePatterns)
			if err != nil {
				t.Errorf("Error:%+v", err)
				return
			}
			if !containsField(jsonData, "refbooks", t) {
				t.Errorf("Expected field not found in the  RulePatterns")
			}
			jsonData, err = json.Marshal(reference.RuleActions)
			if err != nil {
				t.Errorf("Error:%+v", err)
				return
			}
			if !containsField(jsonData, "usps", t) {
				t.Errorf("Expected field not found in the RuleActions")
			}
		}
	}

}

func testPurge(tests *[]doMatchTest, t *testing.T) {

	err := Purge()
	if err != nil {
		t.Errorf("ERROR Purge %+v", err)
	}
}

func testReload(tests *[]doMatchTest, t *testing.T) {

	/*err := Reload()
	if err != nil {
		t.Errorf("ERROR Reload %+v", err)
	}*/
	// Not needed its a combination of purge and load func
}
