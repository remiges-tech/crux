package crux

import "time"

type ruleOp_t int

const (
	ruleOpEQ ruleOp_t = iota // Equal to
	ruleOpNE                 // Not equal to
	ruleOpGT                 // Greater than
	ruleOpGE                 // Greater than or equal to
	ruleOpLT                 // Less than
	ruleOpLE                 // Less than or equal to

	Match    = "+"  // "+"
	Mismatch = "-"  // "-"
	ThenCall = "+c" // "+c"
	ElseCall = "-c" // "-c"
	Return   = "<"  // "<"
	Exit     = "<<" // "<<"
)

type ActionSet_t struct {
	Tasks      []string          `json:"tasks,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type TraceDataRuleL2_t struct {
	Pattern    []map[string][]string `json:"pattern,omitempty"`
	Tasks      []string              `json:"tasks,omitempty"`
	Properties map[string]string     `json:"properties,omitempty"`
	NextSet    int                   `json:"nextset,omitempty"`
	ActionSet  *ActionSet_t          `json:"actionset,omitempty"`
}
type TraceDataRule_t struct {
	RuleNo  int    `json:"r"`
	Res     string `json:"res"`
	NextSet string `json:"nextset,omitempty"`
	TraceDataRuleL2_t
}
type TraceData_t struct {
	SetID   int               `json:"id"`
	SetName string            `json:"setName"`
	Rules   []TraceDataRule_t `json:"rules"`
}

type Trace_t struct {
	Start            *time.Time    `json:"start"`
	End              time.Time     `json:"end"`
	Realm            string        `json:"realm"`
	App              string        `json:"app"`
	EntryRulesetID   int           `json:"entryRulesetID"`
	EntryRulesetName string        `json:"entryRulesetName"`
	TraceData        []TraceData_t `json:"tracedata"`
}
