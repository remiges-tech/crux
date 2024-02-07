package workflow

const (
	setBy     = "admin"
	realmID   = 1
	brwf      = "W"
	isActive  = true
	typeBool  = "bool"
	typeInt   = "int"
	typeFloat = "float"
	typeStr   = "str"
	typeEnum  = "enum"
	typeTS    = "ts"

	timeLayout = "2006-01-02T15:04:05Z"

	opEQ = "eq"
	opNE = "ne"
	opLT = "lt"
	opLE = "le"
	opGT = "gt"
	opGE = "ge"

	trueStr  = "true"
	falseStr = "false"

	step       = "step"
	stepFailed = "stepfailed"
	start      = "START"
	nextStep   = "nextstep"
	done       = "done"
)


var validOps = map[string]bool{
	opEQ: true, opNE: true, opLT: true, opLE: true, opGT: true, opGE: true,
}

type workflowNew struct {
	Slice      int32   `json:"slice" validate:"required,gt=0"`
	App        string  `json:"app" validate:"required,alpha"`
	Class      string  `json:"class" validate:"required,lowercase"`
	Name       string  `json:"name" validate:"required,lowercase"`
	IsInternal bool    `json:"is_internal" validate:"required"`
	Flowrules  []Rules `json:"flowrules" validate:"required"`
}

type Rules struct {
	RulePattern []RulePattern `json:"rulepattern" validate:"required"`
	RuleActions RuleActions   `json:"ruleactions" validate:"required"`
}

type RulePattern struct {
	AttrName string `json:"attr" validate:"required"`
	Op       string `json:"op" validate:"required"`
	AttrVal  any    `json:"val" validate:"required"`
}
type RuleActions struct {
	Tasks      []string   `json:"tasks" validate:"required"`
	Properties []Property `json:"properties" validate:"required"`
	ThenCall   string     `json:"thenCall,omitempty"`
	ElseCall   string     `json:"elseCall,omitempty"`
	WillReturn bool       `json:"willReturn,omitempty"`
	WillExit   bool       `json:"willExit,omitempty"`
}

type Property struct {
	Name string `json:"name" validate:"required"`
	Val  string `json:"val" validate:"required"`
}
