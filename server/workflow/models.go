package workflow

const (
	setBy    = "admin"
	realmID  = 1
	brwf     = "W"
	isActive = true
)

type RulePattern struct {
	AttrName string      `json:"attrname" validate:"required"`
	Op       string      `json:"op" validate:"required"`
	AttrVal  interface{} `json:"attrval" validate:"required"`
}

type RuleActions struct {
	Tasks      []string            `json:"tasks" validate:"required"`
	Properties []map[string]string `json:"properties" validate:"required"`
}
type Rule struct {
	RulePattern RulePattern `json:"rulepattern" validate:"required"`
	RuleActions RuleActions `json:"ruleactions" validate:"required"`
}

type workflowNew struct {
	Slice      int32  `json:"slice" validate:"required,gt=0"`
	App        string `json:"App" validate:"required,alpha"`
	Class      string `json:"class" validate:"required,lowercase"`
	Name       string `json:"name" validate:"required,lowercase"`
	IsInternal bool   `json:"isInternal" validate:"required"`
	Flowrules  []Rule `json:"flowRule" validate:"required"`
}
