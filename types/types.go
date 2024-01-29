package types

type AppConfig struct {
	DBHost        string `json:"db_host"`
	DBPort        int    `json:"db_port"`
	DBUser        string `json:"db_user"`
	DBPassword    string `json:"db_password"`
	DBName        string `json:"db_name"`
	DriverName    string `json:"driver_name"`
	AppServerPort int    `json:"driver_name"`
	ErrorTypeFile string `json:"error_type_file"`
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
	QualifiedCaps []QualifiedCap `json:"qualifiedcaps"`
}
type Attribute struct {
	Name      string   `json:"name"`
	ShortName string   `json:"shortname"`
	LongDesc  string   `json:"longdesc"`
	ValType   string   `json:"valtype"`
	Vals      []string `json:"vals,omitempty"`
	Enumdesc  []string `json:"enumdesc,omitempty"`
	ValMax    int32    `json:"valmax,omitempty"`
	ValMin    int32    `json:"valmin,omitempty"`
	LenMax    int32    `json:"lemmax,omitempty"`
	LenMin    int32    `json:"lenmin,omitempty"`
}
type Patternschema struct {
	Class string      `json:"class"`
	Attr  []Attribute `json:"attr"`
}

type Actionschema struct {
	Class      string   `json:"class"`
	Task       []string `json:"task"`
	Properties []string `json:"properties"`
}
