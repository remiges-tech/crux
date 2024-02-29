package serverBRE

var ValidTypes = map[string]bool{
	"int": true, "float": true, "str": true, "enum": true, "bool": true, "timestamps": true,
}

const (
	CreatedBy    = "admin"
	EditedBy     = "admin"
	RealmID      = 1
	CruxIDRegExp = `^[a-z][a-z0-9_]*$`
	BRWF         = "B"
	Queries      = "queries"
	SetBy        = "admin"
	TimeLayout   = "2006-01-02T15:04:05Z"

	TypeBool  = "bool"
	TypeInt   = "int"
	TypeFloat = "float"
	TypeStr   = "str"
	TypeEnum  = "enum"
	TypeTS    = "ts"

	opEQ = "eq"
	opNE = "ne"
	opLT = "lt"
	opLE = "le"
	opGT = "gt"
	opGE = "ge"

	TrueStr  = "true"
	FalseStr = "false"

	Step       = "step"
	StepFailed = "stepfailed"
	Start      = "START"
	NextStep   = "nextstep"
	Done       = "done"
)

var ValidOps = map[string]bool{
	opEQ: true, opNE: true, opLT: true, opLE: true, opGT: true, opGE: true,
}
