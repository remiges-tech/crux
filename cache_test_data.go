package main

import (
	sqlc "crux/db/sqlc-gen"
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
			"properties": {"discount": "shipby"}
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
			"properties": {"status": "destination"}
		}`),
		Createdat: pgtype.Timestamp{Time: time.Now()},
		Createdby: "user2",
		Editedat:  pgtype.Timestamp{Time: time.Now()},
		Editedby:  pgtype.Text{String: "user2"},
	},
	{
		Realm:         1,
		App:           "test3",
		Slice:         3,
		Class:         "inventoryMaterial",
		Brwf:          sqlc.BrwfEnum("W"),
		Patternschema: []byte(`{"attr": [{"name": "material", "valtype": "enum", "vals": ["cotton", "leather", "metal"]},{"name": "quantity", "valtype": "int"}]}`),
		Actionschema:  []byte(`{"tasks": ["notify", "cancel", "schedule"],"properties": {"message": "timestamp"}}`),
		Createdat:     pgtype.Timestamp{Time: time.Now()},
		Createdby:     "user3",
		Editedat:      pgtype.Timestamp{Time: time.Now()},
		Editedby:      pgtype.Text{String: "user3"},
	},
	{
		Realm: 1,
		App:   "test4",
		Slice: 4,
		Class: "inventoryColor",

		Brwf: sqlc.BrwfEnum("W"),
		Patternschema: []byte(`{
			"attr": [
				{"name": "color", "valtype": "enum", "vals": ["red", "blue", "green"]},
				{"name": "weight", "valtype": "float"}
			]
		}`),

		Actionschema: []byte(`{
			"tasks": ["ship", "receive", "track"],
			"properties": {"carrier": "trackingNumber"}
		}`),
		Createdat: pgtype.Timestamp{Time: time.Now()},
		Createdby: "user4",
		Editedat:  pgtype.Timestamp{Time: time.Now()},
		Editedby:  pgtype.Text{String: "user4"},
	},
}

var mockRulesets = []sqlc.Ruleset{
	{
		Realm:   1,
		App:     "Test1",
		Slice:   1,
		Class:   "inventoryChristmas",
		Setname: "yearendoffer",
		Brwf:    "B",
		Ruleset: []byte(`[{
			"rulepattern": [
				{"attr": "cat", "op": "eq", "val": "textbook"},
				{"attr": "mrp", "op": "ge", "val": "5000"}
			],
			"ruleactions": {
				"tasks": ["christmassale"],
				"properties": {"shipby": "fedex"}
			}
		}]`),
	},
	{
		Realm:   1,
		App:     "Test2",
		Slice:   2,
		Class:   "inventoryNewyear",
		Setname: "newyear",
		Brwf:    "B",
		Ruleset: []byte(`[
			{
			  "rulepattern": [
				{"attr": "cat", "op": "eq", "val": "notebook"},
				{"attr": "mrp", "op": "ge", "val": "3000"}
			  ],
			  "ruleactions": {
				"tasks": ["newyearsale"],
				"properties": {"shipby": "dhl"},
				"thencall": "yearendoffer"
			  }
			}
		  ]`),
	},
	{
		Realm: 1,
		App:   "Test3",
		Slice: 3,
		Class: "inventoryClearance",
		Brwf:  "B",
		Ruleset: []byte(`[{
			"rulepattern": [
				{"attr": "cat", "op": "eq", "val": "stationery"},
				{"attr": "mrp", "op": "ge", "val": "1000"}
			],
			"ruleactions": {
				"tasks": ["clearancesale"],
				"properties": {"shipby": "ups"},
				"thencall": "yearendoffer",
				"elsecall":"nosale"
			}
		}]`),
	},
	{
		Realm:   1,
		App:     "Test4",
		Slice:   4,
		Class:   "inventorySummer",
		Brwf:    "B",
		Setname: "nosale",
		Ruleset: []byte(`[{
			"rulepattern": [
				{"attr": "cat", "op": "eq", "val": "refbooks"},
				{"attr": "mrp", "op": "ge", "val": "200"}
			],
			"ruleactions": {
				"tasks": ["summersale"],
				"properties": {"shipby": "usps"}
			}
		}]`),
	},
}
