package main

import (
	"sync"
	"time"
)

type realm_t string
type app_t string
type slice_t int
type className_t string
type BrwfEnum string

var ConnectionString = "host=localhost port=5432 user=postgres password=postgres dbname=crux sslmode=disable"

type statsSchema_t struct {
	NChecked int
}

type statsRuleset_t struct {
	NCalled    int
	RulesStats []map[realm_t]map[app_t]map[slice_t][]statsSchema_t
}

type statsPerSlice_t struct {
	LoadedAt   time.Time
	BRSchema   map[className_t]statsSchema_t
	BRRulesets map[className_t][]statsRuleset_t
	WFSchema   map[className_t]statsSchema_t
	Workflows  map[className_t][]statsRuleset_t
}

type statsPerApp_t map[slice_t]statsPerSlice_t

type statsPerRealm_t map[app_t]statsPerApp_t

type rulesetStats_t map[realm_t]statsPerRealm_t

type schema_t struct {
	Class         string          `json:"class"`
	PatternSchema patternSchema_t `json:"patternschema"`
	ActionSchema  actionSchema_t  `json:"actionschema"`
	NChecked      int32           `json:"n_checked"`
}
type patternSchema_t struct {
	Attr []attr_t `json:"attr"`
}
type attr_t struct {
	Name     string   `json:"name"`
	ValType  string   `json:"valtype"`
	EnumVals []string `json:"vals,omitempty"`
	ValMin   float64  `json:"valmin,omitempty"`
	ValMax   float64  `json:"valmax,omitempty"`
	LenMin   int      `json:"lenmin,omitempty"`
	LenMax   int      `json:"lenmax,omitempty"`
}
type actionSchema_t struct {
	Tasks      []string          `json:"tasks"`
	Properties map[string]string `json:"properties"`
}

type rulePatternBlock_t struct {
	Attr    string    `json:"attr"`
	Op      string    `json:"op"`
	Val     string    `json:"val"`
	ValType valType_t `json:"valtype,omitempty"`
}

type ruleActionBlock_t struct {
	Task          []string          `json:"tasks"`
	Properties    map[string]string `json:"properties"`
	ThenCall      string            `json:"thencall,omitempty"`
	ElseCall      string            `json:"elsecall,omitempty"`
	DoReturn      bool              `json:"doreturn,omitempty"`
	DoExit        bool              `json:"doexit,omitempty"`
	References    []*Ruleset_t      `json:"references,omitempty"`
	ReferenceType string            `json:"referencetype,omitempty"`
}

type valType_t int

const (
	valInt_t valType_t = iota
	valFloat_t
	valString_t
	valBool_t
	valTimestamp_t
	valEnum_t
)

type Ruleset_t struct {
	Class         string               `json:"class,omitempty""`
	SetName       string               `json:"setname,omitempty""`
	RulePatterns  []rulePatternBlock_t `json:"rulepattern"`
	RuleActions   ruleActionBlock_t    `json:"ruleactions"`
	NMatched      int                  `json:"nMatched"`
	NFailed       int                  `json:"nFailed"`
	ReferenceType string               `json:"referenceType,omitempty""`
	NextRuleset   *Ruleset_t
}

type perSlice_t struct {
	LoadedAt   time.Time
	BRSchema   map[className_t][]*schema_t
	BRRulesets map[className_t][]*Ruleset_t
	WFSchema   map[className_t][]*schema_t
	Workflows  map[className_t][]*Ruleset_t
}

type perApp_t map[slice_t]perSlice_t

type perRealm_t map[app_t]perApp_t

type rulesetCache_t map[realm_t]perRealm_t

type schemaCache_t map[realm_t]perRealm_t

var (
	rulesetCache    rulesetCache_t
	schemaCache     schemaCache_t
	cacheLock       sync.RWMutex
	nCheckedcounter int32
)
