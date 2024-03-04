package ruleset

type RuleSetNew struct {
	Slice      int32   `json:"slice" validate:"required,gt=0"`
	App        string  `json:"app" validate:"required,alpha"`
	Class      string  `json:"class" validate:"required,lowercase"`
	Name       string  `json:"name" validate:"required,lowercase"`
	IsInternal bool    `json:"is_internal" validate:"required"`
	Flowrules  []Rules `json:"flowrules" validate:"required,dive"`
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
