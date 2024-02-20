/*
This file contains the functions that represent Cache tests for Load()/Purge()/Reload(). These functions are called
inside TestCache()) in do_matchest.go.

Some of the definitions of rulesets below deliberately use a lot of whitespace to keep the code consistent
and to make it easier to understand, add to, and edit these tests
*/

package main

import (
	sqlc "crux/db/sqlc-gen"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

var mockSchemasets = []sqlc.Schema{

	{
		Realm: 1,
		App:   "test1",
		Slice: 1,
		Class: "inventoryitems",
		Brwf:  sqlc.BrwfEnum("B"),
		Patternschema: []byte(`{
			"attr": [{
				"name": "cat",
				"valtype": "enum",
				"vals": ["textbook", "notebook", "stationery", "refbooks"]
			},{
				"name": "mrp",
				"valtype": "float"
			},{
				"name": "fullname",
				"valtype": "str"
			},{
				"name": "ageinstock",
				"valtype": "int"
			},{
				"name": "inventoryqty",
				"valtype": "int"
			}]
		}`),
		Actionschema: []byte(`{
			"tasks": ["invitefordiwali", "allowretailsale", "assigntotrash"],
			"properties": ["discount", "shipby"]
		}`),
		Createdat: pgtype.Timestamp{Time: time.Now()},
		Createdby: "user1",
		Editedat:  pgtype.Timestamp{Time: time.Now()},
		Editedby:  pgtype.Text{String: "user1"},
	},
	{
		Realm: 1,
		App:   "test2",
		Slice: 2,
		Class: "inventorySize",
		Brwf:  sqlc.BrwfEnum("W"),
		Patternschema: []byte(`{
			"attr": [
				{"name": "size", "valtype": "enum", "vals": ["small", "medium", "large"]},
				{"name": "price", "valtype": "float"}
			]
		}`),

		Actionschema: []byte(`{
			"tasks": ["approve", "dispatch", "verify"],
			"properties": ["status", "destination"]
		}`),
		Createdat: pgtype.Timestamp{Time: time.Now()},
		Createdby: "user2",
		Editedat:  pgtype.Timestamp{Time: time.Now()},
		Editedby:  pgtype.Text{String: "user2"},
	},
	{
		Realm: 2,
		App:   "test3",
		Slice: 1,
		Class: "inventoryMaterial",
		Brwf:  sqlc.BrwfEnum("W"),
		Patternschema: []byte(`{
			"attr": [
				{"name": "material", "valtype": "enum", "vals": ["cotton", "leather", "metal"]},
				{"name": "quantity", "valtype": "int"}
			]
		}`),

		Actionschema: []byte(`{
			"tasks": ["notify", "cancel", "schedule"],
			"properties": ["message", "timestamp"]
		}`),
		Createdat: pgtype.Timestamp{Time: time.Now()},
		Createdby: "user3",
		Editedat:  pgtype.Timestamp{Time: time.Now()},
		Editedby:  pgtype.Text{String: "user3"},
	},
	{
		Realm: 2,
		App:   "test4",
		Slice: 3,
		Class: "inventoryColor",
		Brwf:  sqlc.BrwfEnum("W"),
		Patternschema: []byte(`{
			"attr": [
				{"name": "color", "valtype": "enum", "vals": ["red", "blue", "green"]},
				{"name": "weight", "valtype": "float"}
			]
		}`),

		Actionschema: []byte(`{
			"tasks": ["ship", "receive", "track"],
			"properties": ["carrier", "trackingNumber"]
		}`),
		Createdat: pgtype.Timestamp{Time: time.Now()},
		Createdby: "user4",
		Editedat:  pgtype.Timestamp{Time: time.Now()},
		Editedby:  pgtype.Text{String: "user4"},
	},
	// Add more mock rulesets as needed
}

func testCache(tests *[]doMatchTest, t *testing.T) {

	testLoad(tests, t)
	testPurge(tests)
	testReload(tests)
}

var mockRulesets = []sqlc.Ruleset{
	{
		Realm: 1,
		App:   "Test1",
		Slice: 1,
		Class: "inventoryChristmas",
		Brwf:  "B",
		Ruleset: []byte(`{
			"rulepattern": [
				{"attrname": "cat", "op": "eq", "attrval": "textbook"},
				{"attrname": "mrp", "op": "ge", "attrval": 5000}
			],
			"ruleactions": {
				"tasks": ["christmassale"],
				"properties": {"shipby": "fedex"}
			}
		}`),
	},
	{
		Realm: 1,
		App:   "Test1",
		Slice: 2,
		Class: "inventoryNewyear",
		Brwf:  "B",
		Ruleset: []byte(`{
			"rulepattern": [
				{"attr": "cat", "op": "eq", "val": "notebook"},
				{"attr": "mrp", "op": "ge", "val": 3000}
			],
			"ruleactions": {
				"tasks": ["newyearsale"],
				"properties": {"shipby": "dhl"}
				"thencall": "inventoryChristmas",
			}
		}`),
	},
	{
		Realm: 2,
		App:   "Test1",
		Slice: 1,
		Class: "inventoryClearance",
		Brwf:  "B",
		Ruleset: []byte(`{
			"rulepattern": [
				{"attr": "cat", "op": "eq", "val": "stationery"},
				{"attr": "mrp", "op": "ge", "val": 1000}
			],
			"ruleactions": {
				"tasks": ["clearancesale"],
				"properties": {"shipby": "ups"}
			}
		}`),
	},
	{
		Realm: 1,
		App:   "Test1",
		Slice: 1,
		Class: "inventorySummer",
		Brwf:  "B",
		Ruleset: []byte(`{
			"rulepattern": [
				{"attr": "cat", "op": "eq", "val": "refbooks"},
				{"attr": "mrp", "op": "ge", "val": 200}
			],
			"ruleactions": {
				"tasks": ["summersale"],
				"properties": [{"shipby": "usps"}]
			}
		}`),
	},
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
		"properties": ["message", "timestamp"]
	}`)

	loadInternal(mockSchemasets, []sqlc.Ruleset{})

	PrintAllRuleSetCache()

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

		fmt.Printf("ExpectedPattern: %+v \n ", expectedPattern)
		fmt.Printf(" \n ")
		fmt.Printf("ActualPattern: %+v \n ", actualPattern)
		fmt.Printf(" \n ")
		fmt.Printf("ExpectedAction: %+v \n", expectedAction)
		fmt.Printf(" \n ")
		fmt.Printf("ActualActual: %+v \n", actualAction)
	}

	//p, a, val := retrieveRulesetFromCache(RealmId, App, Class, Slice, Brwf)

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
