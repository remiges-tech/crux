package types

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
	QualifiedCaps []QualifiedCap `json:"qualifiedcaps"`
}

type Attribute struct {
	Name      string   `json:"name" validate:"required"`
	ShortName string   `json:"shortname" validate:"required"`
	LongDesc  string   `json:"longdesc" validate:"required"`
	ValType   string   `json:"valtype" validate:"required"`
	Vals      []string `json:"vals,omitempty"`
	Enumdesc  []string `json:"enumdesc,omitempty"`
	ValMax    int32    `json:"valmax,omitempty"`
	ValMin    int32    `json:"valmin,omitempty"`
	LenMax    int32    `json:"lemmax,omitempty"`
	LenMin    int32    `json:"lenmin,omitempty"`
}
type Patternschema struct {
	Class string      `json:"class" validate:"required"`
	Attr  []Attribute `json:"attr"`
}

type Actionschema struct {
	Class      string   `json:"class" validate:"required"`
	Task       []string `json:"tasks" validate:"required"`
	Properties []string `json:"properties" validate:"required"`
}
