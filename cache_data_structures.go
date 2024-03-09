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

const (
	valInt_t valType_t = iota
	valFloat_t
	valString_t
	valBool_t
	valTimestamp_t
	valEnum_t
)

type perSlice_t struct {
	LoadedAt   time.Time
	BRSchema   map[className_t][]*schema_t
	BRRulesets map[className_t][]*Ruleset_t
	WFSchema   map[className_t][]*schema_t
	Workflows  map[className_t][]*Ruleset_t
}

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
	Name     string              `json:"name"`
	ValType  string              `json:"valtype"`
	EnumVals map[string]struct{} `json:"vals,omitempty"`
	//EnumVals map[string]bool `json:"vals,omitempty"`
	ValMin float64 `json:"valmin,omitempty"`
	ValMax float64 `json:"valmax,omitempty"`
	LenMin int     `json:"lenmin,omitempty"`
	LenMax int     `json:"lenmax,omitempty"`
}

type valType_t int

const (
	valInt valType_t = iota
	valFloat
	valString
	valBool
	valTimestamp
	valEnum
)

type actionSchema_t struct {
	Tasks      []string
	Properties []string
}

type ruleOp_t int

const (
	ruleOpEQ ruleOp_t = iota // Equal to
	ruleOpNE                 // Not equal to
	ruleOpGT                 // Greater than
	ruleOpGE                 // Greater than or equal to
	ruleOpLT                 // Less than
	ruleOpLE                 // Less than or equal to
)

type rulePatternBlock_t struct {
	Attr string `json:"attr"`
	Op   string `json:"op"`
	Val  any    `json:"val"`
}

type ruleActionBlock_t struct {
	Task          []string          `json:"tasks"`
	Properties    map[string]string `json:"properties"`
	ThenCall      string            `json:"thencall,omitempty"`
	ElseCall      string            `json:"elsecall,omitempty"`
	DoReturn      bool              `json:"doreturn,omitempty"`
	DoExit        bool              `json:"doexit,omitempty"`
	References    []*Ruleset_t      `json:"-"`
	ReferenceType string            `json:"referencetype,omitempty"`
}

type rule_t struct {
	RulePatterns []rulePatternBlock_t `json:"rulepattern"`
	RuleActions  ruleActionBlock_t    `json:"ruleactions"`
	NMatched     int32
	NFailed      int32
}

type Ruleset_t struct {
	Id            int32    `json:"id"`
	Class         string   `json:"class"`
	SetName       string   `json:"setname"`
	Rules         []rule_t `json:"rule"`
	NCalled       int32
	ReferenceType string
}

type perApp_t map[slice_t]perSlice_t

type perRealm_t map[app_t]perApp_t

type rulesetCache_t map[realm_t]perRealm_t

type schemaCache_t map[realm_t]perRealm_t

var (
	rulesetCache rulesetCache_t
	schemaCache  schemaCache_t
	cacheLock    sync.RWMutex
)
