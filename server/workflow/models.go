package workflow

import "github.com/jackc/pgx/v5/pgtype"

const (
	editedBy = "admin"

	brwf      = "W"
	isActive  = false
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

// error constant

const (
	ONLY_NUMBERS = "only_numbers_allowed"
	ONLY_BOOL    = "only_boolean_allowed"
	MARK         = "%"
	DB_ERROR     = "failed to get data from DB"
	AUTH_ERROR   = "authorization_error"
)

var validOps = map[string]bool{
	opEQ: true, opNE: true, opLT: true, opLE: true, opGT: true, opGE: true,
}

type Rules struct {
	RulePattern []RulePattern `json:"rulepattern" validate:"required,dive"`
	RuleActions RuleActions   `json:"ruleactions" validate:"required"`
}

type RulePattern struct {
	Attr string `json:"attr" validate:"required"`
	Op   string `json:"op" validate:"required"`
	Val  any    `json:"val" validate:"required"`
}
type RuleActions struct {
	Tasks      []string   `json:"tasks" validate:"required"`
	Properties []Property `json:"properties" validate:"required,dive"`
	ThenCall   string     `json:"thenCall,omitempty"`
	ElseCall   string     `json:"elseCall,omitempty"`
	WillReturn bool       `json:"willReturn,omitempty"`
	WillExit   bool       `json:"willExit,omitempty"`
}

type Property struct {
	Name string `json:"name" validate:"required"`
	Val  string `json:"val" validate:"required"`
}

type WorkflowGetReq struct {
	Slice int32  `json:"slice" validate:"required,gt=0,lt=15"`
	App   string `json:"app" validate:"required,alpha,lt=15"`
	Class string `json:"class" validate:"required,alpha,lt=15"`
	Name  string `json:"name" validate:"required,alpha,lt=15"`
}

type WorkflowgetRow struct {
	ID         int32            `json:"id"`
	Slice      int32            `json:"slice"`
	App        string           `json:"app"`
	Class      string           `json:"class"`
	Name       string           `json:"name"`
	IsActive   bool             `json:"is_active"`
	IsInternal bool             `json:"is_internal"`
	Flowrules  interface{}      `json:"flowrules"`
	Createdat  pgtype.Timestamp `json:"createdat"`
	Createdby  string           `json:"createdby"`
	Editedat   pgtype.Timestamp `json:"editedat"`
	Editedby   pgtype.Text      `json:"editedby"`
}

type WorkflowListReq struct {
	Slice      int32  `json:"slice,omitempty" validate:"lt=15"`
	App        string `json:"app,omitempty" validate:"lt=15"`
	Class      string `json:"class,omitempty" validate:"lt=15"`
	Name       string `json:"name,omitempty" validate:"lt=15"`
	IsActive   bool   `json:"is_active,omitempty"`
	IsInternal bool   `json:"is_internal,omitempty"`
}
