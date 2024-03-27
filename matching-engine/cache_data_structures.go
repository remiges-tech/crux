package crux

import (
	"sync"
	"time"
)

type Realm_t string
type App_t string
type Slice_t int
type ClassName_t string
type BrwfEnum string

const (
	ValInt_t ValType_t = iota
	ValFloat_t
	ValString_t
	ValBool_t
	ValTimestamp_t
	ValEnum_t
)

type PerSlice_t struct {
	LoadedAt   time.Time
	BRSchema   map[ClassName_t]Schema_t
	BRRulesets map[ClassName_t][]*Ruleset_t
	WFSchema   map[ClassName_t]Schema_t
	Workflows  map[ClassName_t][]*Ruleset_t
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

type ValType_t int

const (
	ValInt ValType_t = iota
	ValFloat
	ValString
	ValBool
	ValTimestamp
	ValEnum
)

type ActionSchema_t struct {
	Tasks      []string `json:"tasks" validate:"required"`
	Properties []string `json:"properties" validate:"required"`
}

type RuleOp_t int

const (
	RuleOpEQ RuleOp_t = iota // Equal to
	RuleOpNE                 // Not equal to
	RuleOpGT                 // Greater than
	RuleOpGE                 // Greater than or equal to
	RuleOpLT                 // Less than
	RuleOpLE                 // Less than or equal to
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

type PerApp_t map[Slice_t]PerSlice_t

type PerRealm_t map[App_t]PerApp_t

type RulesetCache_t map[Realm_t]PerRealm_t

type SchemaCache_t map[Realm_t]PerRealm_t

var RulesetCache RulesetCache_t
var SchemaCache SchemaCache_t

var (
	rulesetCache RulesetCache_t
	schemaCache  SchemaCache_t
	cacheLock    sync.RWMutex
)

func init() {
	RulesetCache = make(RulesetCache_t)
	SchemaCache = make(SchemaCache_t)
}
