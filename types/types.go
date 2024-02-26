package types

const (
	DevEnv           Environment = "dev_env"
	ProdEnv          Environment = "prod_env"
	UATEnv           Environment = "uat_env"
	RECORD_NOT_EXIST             = "record_does_not_exist"
	OPERATION_FAILED             = "operation_failed"
)

var APP, SLICE, CLASS, NAME string = "App", "Slice", "Class", "Name"

type Environment string

func (env Environment) IsValid() bool {
	switch env {
	case DevEnv, ProdEnv, UATEnv:
		return true
	}
	return false
}

type OpReq struct {
	User      string   `json:"user"`
	CapNeeded []string `json:"capNeeded"`
	Scope     Scope    `json:"scope"`
	Limit     Limit    `json:"limit"`
}

type Scope map[string]interface{}
type Limit map[string]interface{}

type QualifiedCap struct {
	Id    string `json:"id"`
	Cap   string `json:"cap"`
	Scope Scope  `json:"scope"`
	Limit Limit  `json:"limit"`
}

type Capabilities struct {
	Name          string         `json:"name"` //either user name or group name
	QualifiedCaps []QualifiedCap `json:"qualifiedCaps"`
}

type Attribute struct {
	Name      string   `json:"name" validate:"required"`
	ShortDesc string   `json:"shortdesc" validate:"required"`
	LongDesc  string   `json:"longdesc" validate:"required"`
	ValType   string   `json:"valtype" validate:"required"`
	Vals      []string `json:"vals,omitempty"`
	Enumdesc  []string `json:"enumdesc,omitempty"`
	ValMax    *int32   `json:"valmax,omitempty"`
	ValMin    *int32   `json:"valmin,omitempty"`
	LenMax    *int32   `json:"lenmax,omitempty"`
	LenMin    *int32   `json:"lenmin,omitempty"`
}

type PatternSchema struct {
	Class string      `json:"class" validate:"required,lowercase"`
	Attr  []Attribute `json:"attr" validate:"required,dive"`
}

type ActionSchema struct {
	Class      string   `json:"class" validate:"required,lowercase"`
	Tasks      []string `json:"tasks" validate:"required"`
	Properties []string `json:"properties" validate:"required"`
}
