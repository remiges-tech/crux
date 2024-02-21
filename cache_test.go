/*
This file contains the functions that represent Cache tests for Load()/Purge()/Reload(). These functions are called
inside TestCache()) in do_matchest.go.

Some of the definitions of rulesets below deliberately use a lot of whitespace to keep the code consistent
and to make it easier to understand, add to, and edit these tests
*/

package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func testCache(tests *[]doMatchTest, t *testing.T) {

	testLoad(tests, t)
	testPurge(tests)
	testReload(tests)
}

func testLoad(tests *[]doMatchTest, t *testing.T) {

	// UT Test Data passed to loadInternal()

	var RealmId = 2
	var App = "test3"
	var Class = "inventoryMaterial"
	var Slice = 1
	var Brwf = "W"
	expectedSchemasPattern := []byte(`{
		"attr": [
			{"name": "material", "valtype": "enum", "vals": ["cotton", "leather", "metal"]},
			{"name": "quantity", "valtype": "int"}
		]
	}`)

	expectedSchemasAction := []byte(`{
		"tasks": ["notify", "cancel", "schedule"],
		"properties": {"message": "timestamp"}
	}`)

	err := loadInternal(mockSchemasets, mockRulesets)
	if err != nil {
		t.Errorf(" %v", err)
		return
	}

	p, a, val := retrieveSchemasFromCache(RealmId, App, Class, Slice, Brwf)

	if val == "success" {
		expectedPattern := make(map[string]interface{})
		if err := json.Unmarshal(expectedSchemasPattern, &expectedPattern); err != nil {
			t.Errorf("Error unmarshalling expectedSchemasPattern: %v", err)
			return
		}

		expectedAction := make(map[string]interface{})
		if err := json.Unmarshal(expectedSchemasAction, &expectedAction); err != nil {
			t.Errorf("Error unmarshalling expectedSchemasAction: %v", err)
			return
		}
		actualPattern := make(map[string]interface{})
		if err := json.Unmarshal(p, &actualPattern); err != nil {
			t.Errorf("Error unmarshalling actual pattern: %v", err)
			return
		}

		actualAction := make(map[string]interface{})
		if err := json.Unmarshal(a, &actualAction); err != nil {
			t.Errorf("Error unmarshalling actual action: %v", err)
			return
		}

		fmt.Printf("ExpectedSchemaPattern: %+v \n ", expectedPattern)
		fmt.Printf(" \n ")
		fmt.Printf("ActualSchemaPattern: %+v \n ", actualPattern)
		fmt.Printf(" \n ")
		fmt.Printf("ExpectedSchemaAction: %+v \n", expectedAction)
		fmt.Printf(" \n ")
		fmt.Printf("ActualSchemaActual: %+v \n", actualAction)
	}

	RealmId = 1
	App = "Test2"
	Slice = 2
	Class = "inventoryNewyear"
	Brwf = "B"
	expectedRulePattern := []byte(`[
		{"attr": "cat", "op": "eq", "val": "notebook"},
		{"attr": "mrp", "op": "ge", "val": "3000"}
	  ]`)

	expectedRuleAction := []byte(`{
		"tasks": ["newyearsale"],
		"properties": {"shipby": "dhl"},
		"thencall": "inventoryChristmas"
	  }`)
	rp, ra, rval, refrenceRuleset := retrieveRulesetFromCache(RealmId, App, Class, Slice, Brwf)

	if rval == "success" {

		var expectedPattern []rulePatternBlock_t
		if err := json.Unmarshal(expectedRulePattern, &expectedPattern); err != nil {
			t.Errorf("Error unmarshalling expectedRulePattern: %v", err)
			return
		}

		var expectedAction ruleActionBlock_t
		if err := json.Unmarshal(expectedRuleAction, &expectedAction); err != nil {
			t.Errorf("Error unmarshalling expectedRuleAction: %v", err)
			return
		}
		var actualPattern []rulePatternBlock_t
		if err := json.Unmarshal(rp, &actualPattern); err != nil {
			t.Errorf("Error unmarshalling actualPattern: %v", err)
			return
		}

		var actualAction ruleActionBlock_t

		if err := json.Unmarshal(ra, &actualAction); err != nil {
			t.Errorf("Error unmarshalling actual action: %v", err)
			return
		}

		fmt.Printf("ExpectedRulePattern: %+v \n ", expectedPattern)
		fmt.Printf(" \n ")
		fmt.Printf("ActualRulePattern: %+v \n ", actualPattern)
		fmt.Printf(" \n ")
		fmt.Printf("ExpectedRuleAction: %+v \n", expectedAction)
		fmt.Printf(" \n ")
		fmt.Printf("ActualRuleActual: %+v \n", actualAction)
		fmt.Printf(" \n ")
		fmt.Printf("refrenceRuleset: %+v \n", refrenceRuleset)

	}
}

func testPurge(tests *[]doMatchTest) {

	/*err := Purge()
	if err != nil {
		log.Println("ERROR Purge", err)
	}*/
}

func testReload(tests *[]doMatchTest) {

	/*err := Reload()
	if err != nil {
		log.Println("ERROR Reload", err)
	}*/
}
