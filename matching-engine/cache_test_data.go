package crux

import (
	"time"

	"github.com/remiges-tech/crux/db/sqlc-gen"

	"github.com/jackc/pgx/v5/pgtype"
)

var mockSchemasets = []sqlc.Schema{

	{
		Realm: "1",
		App:   "Test1",
		Slice: 1,
		Class: "inventoryitem2",

		Brwf: sqlc.BrwfEnum("W"),
		Patternschema: []byte(`[
		
				{
					"attr": "cat",
					"valtype": "str"
					
				},
				{
					"attr": "mrp",
					"valtype": "float"
				},
				{
					"attr": "fullname",
					"valtype": "str"
				},
				{
					"attr": "ageinstock",
					"valtype": "int"
				},
				{
					"attr": "inventoryqty",
					"valtype": "int"
				}
			
		]`),
		Actionschema: []byte(`{
			
		}`),
		Createdat: pgtype.Timestamp{Time: time.Now()},
		Createdby: "user4",
		Editedat:  pgtype.Timestamp{Time: time.Now()},
		Editedby:  pgtype.Text{String: "user4"},
	},
	{
		Realm: "1",
		App:   "Test6",
		Slice: 6,
		Class: "inventoryitem",

		Brwf: sqlc.BrwfEnum("W"),
		Patternschema: []byte(`[
			
				{
					"attr": "cat",
					"valtype": "enum",
					"vals": {"textbook":{}, "book":{}, "stationery":{}, "refbooks":{}}
				},
				{
					"attr": "mrp",
					"valtype": "float"
				},
				{
					"attr": "fullname",
					"valtype": "str"
				},
				{
					"attr": "ageinstock",
					"valtype": "int"
				},
				{
					"attr": "inventoryqty",
					"valtype": "int"
				},
				{
					"attr": "received",
					"valtype": "ts"
				}
			
		]`),
		Actionschema: []byte(`{
			
		}`),
		Createdat: pgtype.Timestamp{Time: time.Now()},
		Createdby: "user4",
		Editedat:  pgtype.Timestamp{Time: time.Now()},
		Editedby:  pgtype.Text{String: "user4"},
	},

	{
		Realm: "1",
		App:   "Test7",
		Slice: 7,
		Class: "transaction",
		Brwf:  sqlc.BrwfEnum("W"),
		Patternschema: []byte(`[
			
				{
					"attr": "productname",
					"valtype": "str"
					
				},
				{
					"attr": "price",
					"valtype": "int"
				},
				{
					"attr": "inwintersale",
					"valtype": "bool"
				},
				{
					"attr": "paymenttype",
					"valtype": "enum"
				},
				{
					"attr": "ismember",
					"valtype": "bool"
				},
				{  
					"attr": "received",
					"valtype": "ts"
				}
			
		]`),
		Actionschema: []byte(`{
			
		}`),
		Createdat: pgtype.Timestamp{Time: time.Now()},
		Createdby: "user4",
		Editedat:  pgtype.Timestamp{Time: time.Now()},
		Editedby:  pgtype.Text{String: "user4"},
	},
	{
		Realm: "1",
		App:   "Test8",
		Slice: 8,
		Class: "purchase",
		Brwf:  sqlc.BrwfEnum("W"),
		Patternschema: []byte(`[
			
				{
					"attr": "product",
					"valtype": "str"
					
				},
				{
					"attr": "price",
					"valtype": "float"
				},
				{
					"attr": "ismember",
					"valtype": "bool"
				}
				
			
		]`),
		Actionschema: []byte(`{
			"tasks": ["freepen", "freebottle", "freepencil", "freemug", "freejar", "freeplant","freebag", "freenotebook"],
			"properties": ["discount", "pointsmult"]
		}`),
		Createdat: pgtype.Timestamp{Time: time.Now()},
		Createdby: "user4",
		Editedat:  pgtype.Timestamp{Time: time.Now()},
		Editedby:  pgtype.Text{String: "user4"},
	},
	{

		Realm: "1",
		App:   "Test9",
		Slice: 9,
		Class: "order",
		Brwf:  sqlc.BrwfEnum("W"),
		Patternschema: []byte(`[
			
				{
					"attr": "ordertype",
					"valtype": "enum"
					
				},
				{
					"attr": "mode",
					"valtype": "enum"
				},
				{
					"attr": "liquidscheme",
					"valtype": "bool"
				},{
					"attr": "overnightscheme",
					"valtype": "bool"
				},{
					"attr": "extendedhours",
					"valtype": "bool"
				}
				
			
		]`),
		Actionschema: []byte(`{
			"tasks": ["unitstoamc", "unitstorta"],
			"properties": ["amfiordercutoff", "bseordercutoff", "fundscutoff", "unitscutoff"]
		}`),
		Createdat: pgtype.Timestamp{Time: time.Now()},
		Createdby: "user4",
		Editedat:  pgtype.Timestamp{Time: time.Now()},
		Editedby:  pgtype.Text{String: "user4"},
	},

	{
		Realm: "1",
		App:   "Test10",
		Slice: 10,
		Class: "ucccreation",
		Brwf:  sqlc.BrwfEnum("W"),
		Patternschema: []byte(`[
		
				{
					"attr": "step",
					"valtype": "enum"
					
				},
				{
					"attr": "stepfailed",
					"valtype": "bool"
					
				},
				{
					"attr": "mode",
					"valtype": "enum"
				}
				
			
		]`),
		Actionschema: []byte(`{
			"tasks": ["getcustdetails", "aof", "kycvalid", "nomauth", "bankaccvalid", "dpandbankaccvalid", "sendauthlinktoclient"],
			"properties": ["nextstep", "done"]
		}`),
		Createdat: pgtype.Timestamp{Time: time.Now()},
		Createdby: "user4",
		Editedat:  pgtype.Timestamp{Time: time.Now()},
		Editedby:  pgtype.Text{String: "user4"},
	},

	{
		Realm: "1",
		App:   "Test11",
		Slice: 11,
		Class: "prepareaof",
		Brwf:  sqlc.BrwfEnum("W"),
		Patternschema: []byte(`[
			
				{
					"attr": "step",
					"valtype": "enum"
					
				},
				{
					"attr": "stepfailed",
					"valtype": "bool"
					
				}
		]`),
		Actionschema: []byte(`{
			
		}`),
		Createdat: pgtype.Timestamp{Time: time.Now()},
		Createdby: "user4",
		Editedat:  pgtype.Timestamp{Time: time.Now()},
		Editedby:  pgtype.Text{String: "user4"},
	},
	{
		Realm: "1",
		App:   "Test12",
		Slice: 12,
		Class: "validateaof",
		Brwf:  sqlc.BrwfEnum("W"),
		Patternschema: []byte(`[

				{
					"attr": "step",
					"valtype": "enum"
					
				},
				{
					"attr": "stepfailed",
					"valtype": "bool"
					
				},
				{
					"attr": "aofexists",
					"valtype": "bool"
					
				}

			
		]`),
		Actionschema: []byte(`{
			
		}`),
		Createdat: pgtype.Timestamp{Time: time.Now()},
		Createdby: "user4",
		Editedat:  pgtype.Timestamp{Time: time.Now()},
		Editedby:  pgtype.Text{String: "user4"},
	},
}

var mockRulesets = []sqlc.Ruleset{

	{
		Realm:   "1",
		App:     "Test1",
		Slice:   1,
		Class:   "inventoryitem2",
		Setname: "main",
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
				"thenCall":"second"
			  }
			}
		  ]`),
	},

	{
		Realm:   "1",
		App:     "Test1",
		Slice:   2,
		Class:   "inventoryitem2",
		Setname: "second",
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
				"thenCall":"third"
			  }
			}
		  ]`),
	},
	{
		Realm:   "1",
		App:     "Test1",
		Slice:   3,
		Class:   "inventoryitem2",
		Setname: "third",
		Brwf:    "B",
		Ruleset: []byte(`[{
			"rulepattern": [
				{"attr": "cat", "op": "eq", "val": "textbook"},
				{"attr": "mrp", "op": "ge", "val": "5000"}
			],
			"ruleactions": {
				"tasks": ["yearendsale", "summersale", "wintersale"],
				"properties": {"discount": "15", "freegift": "mug"},
				"thencall": "second"
			
			}
		}]`),
	},
	{
		Realm:   "1",
		App:     "Test3",
		Slice:   3,
		Class:   "inventoryNewyear",
		Setname: "third",
		Brwf:    "B",
		Ruleset: []byte(`[{
			"rulepattern": [
				{"attr": "cat", "op": "eq", "val": "stationery"},
				{"attr": "mrp", "op": "ge", "val": "1000"}
			],
			"ruleactions": {
				"tasks": ["clearancesale"],
				"properties": {"shipby": "ups"}

			}
		}]`),
	},

	{
		Realm:   "1",
		App:     "Test4",
		Slice:   4,
		Class:   "inventoryNewyear",
		Brwf:    "B",
		Setname: "nosale",
		Ruleset: []byte(`[{
			"rulepattern": [
				{"attr": "cat", "op": "eq", "val": "refbooks"},
				{"attr": "mrp", "op": "ge", "val": "200"}
			],
			"ruleactions": {
				"tasks": ["summersale"],
				"properties": {"shipby": "usps"},
				"thencall": "yearendoffer"
			}
		}]`),
	},
	{

		Realm:   "1",
		App:     "Test5",
		Slice:   5,
		Class:   "inventoryClearance",
		Brwf:    "B",
		Setname: "second",
		Ruleset: []byte(`[{
			"rulepattern": [
				{"attr": "cat", "op": "eq", "val": "refbooks"},
				{"attr": "mrp", "op": "ge", "val": "200"}
			],
			"ruleactions": {
				"tasks": ["summersale"],
				"properties": {"shipby": "usps"},
				"thencall": "yearendoffer"
			}
		}]`),
	},
	{
		Realm:   "1",
		App:     "Test6",
		Slice:   6,
		Class:   "inventoryitem",
		Brwf:    "B",
		Setname: "second",
		Ruleset: []byte(`[{
			"rulepattern": [
				{"attr": "cat", "op": "eq", "val": "textbook"},
				{"attr": "cat", "op": "eq", "val": "refbook"},
				{"attr": "mrp", "op": "ge", "val": 69.50},
				{"attr": "ageinstock", "op": "lt", "val": 7},
				{"attr": "summersale", "op": "eq", "val": true},
				{"attr": "bulkorder", "op": "ne", "val": false} ,
				{"attr": "received", "op": "le", "val": "2018-06-10T15:04:05Z"} 
			],
			"ruleactions": {
				"tasks": ["yearendsale", "summersale", "springsale", "wintersale"],
				"properties": {
					"cashback": "15",
					"discount": "10",
					"freegift": "mug"
				},
				"DoExit": false,
				"DoReturn": true
			}
			}]`),
	},
	{
		Realm:   "1",
		App:     "Test7",
		Slice:   7,
		Class:   "transaction",
		Brwf:    "B",
		Setname: "main",
		Ruleset: []byte(`[{
			"rulepattern": [
				{"attr": "inwintersale", "op": "eq", "val": true}
			],
			"ruleActions": {
				"thenCall": "winterdisc",
				"elseCall": "regulardisc",
				"DoExit": false,
				"DoReturn": false
			}
		},
		{
			"rulepattern": [
				{"attr": "paymenttype", "op": "eq", "val": "cash"},
				{"attr": "price", "op": "gt", "val": 10}
			],
			"ruleActions": {
				"tasks": ["freepen"],
				"DoExit": false,
				"DoReturn": false
			}
		},
		{
			"rulepattern": [
				{"attr": "paymenttype", "op": "eq", "val": "card"},
				{"attr": "price", "op": "gt", "val": 10}
			],
			"ruleActions": {
				"tasks": ["freemug"],
				"DoExit": false,
				"DoReturn": false
			}
		},
		{
			"rulepattern": [
				{"attr": "freehat", "op": "eq", "val": true}
			],
			"ruleActions": {
				"tasks": ["freebag"],
				"DoExit": false,
				"DoReturn": false
			}
		}]`),
	},
	{
		Realm:   "1",
		App:     "Test8",
		Slice:   8,
		Class:   "transaction",
		Brwf:    "B",
		Setname: "memberdisc",
		Ruleset: []byte(`[{
			"rulepattern": [
				{"attr": "productname", "op": "eq", "val": "lamp"},
				{"attr": "price", "op": "gt", "val": 50}
			],
			"ruleActions": {
				"properties": {"discount": "35", "pointsmult": "2"},
				"DoExit": true,
				"DoReturn": false
			}
		},
		{
			"rulepattern": [
				{"attr": "price", "op": "lt", "val": 100}
			],
			"ruleActions": {
				"properties": {"discount": "20"},
				"DoExit": false,
				"DoReturn": false
			}
		},
		{
			"rulepattern": [
				{"attr": "price", "op": "ge", "val": 100}
			],
			"ruleActions": {
				"properties": {"discount": "25"},
				"DoExit": false,
				"DoReturn": false
			}
		}]`),
	},

	{
		Realm:   "1",
		App:     "Test7",
		Slice:   9,
		Class:   "transaction",
		Brwf:    "B",
		Setname: "nonmemberdisc",
		Ruleset: []byte(`[{
			"rulepattern": [
				{"attr": "price", "op": "lt", "val": 50}
			],
			"ruleActions": {
				"properties": {"discount": "5"},
				"DoExit": false,
				"DoReturn": false
			}
		},
		{
			"rulepattern": [
				{"attr": "price", "op": "ge", "val": 50}
			],
			"ruleActions": {
				"properties": {"discount": "10"},
				"DoExit": false,
				"DoReturn": false
			}
		},
		{
			"rulepattern": [
				{"attr": "price", "op": "ge", "val": 100}
			],
			"ruleActions": {
				"properties": {"discount": "15"},
				"DoExit": false,
				"DoReturn": false
			}
		}]`),
	},
}
