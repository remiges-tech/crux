package crux

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
	BRSchema   map[className_t][]*Schema_t
	BRRulesets map[className_t][]*Ruleset_t
	WFSchema   map[className_t][]*Schema_t
	Workflows  map[className_t][]*Ruleset_t
}

type Schema_t struct {
	Class         string            `json:"class"`
	PatternSchema []PatternSchema_t `json:"patternschema"`
	ActionSchema  ActionSchema_t    `json:"actionschema"`
	NChecked      int32             `json:"n_checked"`
}
type PatternSchema_t struct {
	Attr      string              `json:"attr" validate:"required"`
	ShortDesc string              `json:"shortdesc" validate:"required"`
	LongDesc  string              `json:"longdesc" validate:"required"`
	ValType   string              `json:"valtype" validate:"required"`
	EnumVals  map[string]struct{} `json:"vals,omitempty"`
	ValMin    float64             `json:"valmin,omitempty"`
	ValMax    float64             `json:"valmax,omitempty"`
	LenMin    int                 `json:"lenmin,omitempty"`
	LenMax    int                 `json:"lenmax,omitempty"`
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

type ActionSchema_t struct {
	Tasks      []string `json:"tasks" validate:"required"`
	Properties []string `json:"properties" validate:"required"`
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

type RulePatternBlock_t struct {
	Attr string `json:"attr" validate:"required"`
	Op   string `json:"op" validate:"required"`
	Val  any    `json:"val" validate:"required"`
}

type RuleActionBlock_t struct {
	Task          []string          `json:"tasks" validate:"required"`
	Properties    map[string]string `json:"properties" validate:"required"`
	ThenCall      string            `json:"thencall,omitempty"`
	ElseCall      string            `json:"elsecall,omitempty"`
	DoReturn      bool              `json:"doreturn,omitempty"`
	DoExit        bool              `json:"doexit,omitempty"`
	References    []*Ruleset_t      `json:"-"`
	ReferenceType string            `json:"referencetype,omitempty"`
}

type Rule_t struct {
	RulePatterns []RulePatternBlock_t `json:"rulepattern" validate:"required,dive"`
	RuleActions  RuleActionBlock_t    `json:"ruleactions" validate:"required"`
	NMatched     int32
	NFailed      int32
}

type Ruleset_t struct {
	Id            int32    `json:"id"`
	Class         string   `json:"class"`
	SetName       string   `json:"setname"`
	Rules         []Rule_t `json:"rule"`
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
