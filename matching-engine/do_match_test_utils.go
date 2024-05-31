package crux

import (
	"fmt"
	"io"
	"os"
)

type TracingTestCasesStruct struct {
	TestCaseName string
	Entity
	*Ruleset_t
	EntityFilePath     string
	RulesetFilePath    string
	Trace_lev          int
	ExpectedResultFile string
	ExpectedTrace      Trace_t
	Url                string
	RuleSchemasCache   string
}

//Entity

var (
	Entity_return = Entity{
		Realm: "Remiges",
		App:   "tnt",
		Slice: 12,
		Class: "finance",
		Attrs: map[string]string{
			"step": "step2",
			"mode": "demat",
		},
	}
	Entity_match = Entity{
		Realm: "Remiges",
		App:   "tnt",
		Slice: 12,
		Class: "finance",
		Attrs: map[string]string{
			"step": "start",
			"mode": "demat",
		},
	}

	Entity_elsecall = Entity{
		Realm: "Remiges",
		App:   "tnt",
		Slice: 12,
		Class: "sale",
		Attrs: map[string]string{
			"step": "step3",
			"mode": "demat",
		},
	}

	Entity_thencall = Entity{
		Realm: "Remiges",
		App:   "tnt",
		Slice: 12,
		Class: "sale",
		Attrs: map[string]string{
			"cat": "textbook",
			"mrp": "6000",
		},
	}
)

// Ruleset_t
var (
	EmptyTrace      Trace_t
	ThenCallRuleset = Ruleset_t{
		Id:      3,
		Class:   "sale",
		SetName: "third",
		Rules: []Rule_t{
			{
				RulePatterns: []RulePatternBlock_t{
					{Attr: "cat", Op: "eq", Val: "textbook"},
					{Attr: "mrp", Op: "ge", Val: "5000"},
				},
				RuleActions: RuleActionBlock_t{
					Task:       []string{"yearendsale"},
					Properties: map[string]string{"discount": "15"},
					ThenCall:   "second",
				},
			}, {
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step1",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{
						"step2",
					},
					Properties: map[string]string{
						"nextstep": "step2",
					},
				},
				NMatched: 0,
				NFailed:  0,
			},
			{
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step2",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{},
					Properties: map[string]string{
						"nextstep": "step3", "done": "false",
					},
					DoReturn: true,
				},
				NMatched: 0,
				NFailed:  0,
			},
		},
	}
	ElseCallRuleset = Ruleset_t{
		Id:      3,
		Class:   "sale",
		SetName: "third",
		Rules: []Rule_t{
			{
				RulePatterns: []RulePatternBlock_t{
					{Attr: "cat", Op: "eq", Val: "textbook"},
					{Attr: "mrp", Op: "ge", Val: "5000"},
				},
				RuleActions: RuleActionBlock_t{
					Task:       []string{"yearendsale", "summersale", "wintersale"},
					Properties: map[string]string{"discount": "15", "freegift": "mug"},
					ElseCall:   "second",
				},
			}, {
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step2",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{},
					Properties: map[string]string{
						"nextstep": "step3", "done": "false",
					},
					DoReturn: true,
				},
				NMatched: 0,
				NFailed:  0,
			},
			{
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step3",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{
						"step4",
					},
					Properties: map[string]string{
						"nextstep": "step4",
					},
				},
				NMatched: 0,
				NFailed:  0,
			}, {
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step4",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{
						"step5",
					},
					Properties: map[string]string{
						"nextstep": "step5",
					},
				},
				NMatched: 0,
				NFailed:  0,
			},
			{
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step5",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{
						"step6",
					},
					Properties: map[string]string{
						"nextstep": "step6",
					},
				},
				NMatched: 0,
				NFailed:  0,
			},
		},
	}

	ThenCall2Ruleset = Ruleset_t{
		Id:      2,
		Class:   "sale",
		SetName: "second",
		Rules: []Rule_t{
			{
				RulePatterns: []RulePatternBlock_t{
					{Attr: "cat", Op: "eq", Val: "textbook"},
					{Attr: "mrp", Op: "ge", Val: "5000"},
				},
				RuleActions: RuleActionBlock_t{
					Task:       []string{"summersale", "wintersale"},
					Properties: map[string]string{"freegift": "mug"},
				},
			}, {
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step2",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{},
					Properties: map[string]string{
						"nextstep": "step3", "done": "false",
					},
					DoReturn: true,
				},
				NMatched: 0,
				NFailed:  0,
			},
			{
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step3",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{
						"step4",
					},
					Properties: map[string]string{
						"nextstep": "step4",
					},
				},
				NMatched: 0,
				NFailed:  0,
			}, {
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step4",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{
						"step5",
					},
					Properties: map[string]string{
						"nextstep": "step5",
					},
				},
				NMatched: 0,
				NFailed:  0,
			},
			{
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step5",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{
						"step6",
					},
					Properties: map[string]string{
						"nextstep": "step6",
					},
				},
				NMatched: 0,
				NFailed:  0,
			},
		},
	}

	SampleRuleset = Ruleset_t{
		Id:      1,
		Class:   "finance",
		SetName: "finance",
		Rules: []Rule_t{
			{
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "start",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{
						"step1",
					},
					Properties: map[string]string{
						"nextstep": "step1",
					},
				},
				NMatched: 0,
				NFailed:  0,
			},
			{
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step1",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{
						"step2",
					},
					Properties: map[string]string{
						"nextstep": "step2",
					},
				},
				NMatched: 0,
				NFailed:  0,
			},
			{
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step2",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{},
					Properties: map[string]string{
						"nextstep": "step3", "done": "false",
					},
					// DoReturn: true,
				},
				NMatched: 0,
				NFailed:  0,
			},
			{
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step3",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{
						"step4",
					},
					Properties: map[string]string{
						"nextstep": "step4",
					},
				},
				NMatched: 0,
				NFailed:  0,
			}, {
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step4",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{
						"step5",
					},
					Properties: map[string]string{
						"nextstep": "step5",
					},
					DoReturn: true,
				},
				NMatched: 0,
				NFailed:  0,
			},
			{
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step5",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{
						"step6",
					},
					Properties: map[string]string{
						"nextstep": "step6",
					},
				},
				NMatched: 0,
				NFailed:  0,
			},
			{
				RulePatterns: []RulePatternBlock_t{
					{
						Attr: "step",
						Op:   "eq",
						Val:  "step6",
					},
					{
						Attr: "mode",
						Op:   "eq",
						Val:  "demat",
					},
				},
				RuleActions: RuleActionBlock_t{
					Task: []string{},
					Properties: map[string]string{
						"done": "true",
					},
				},
				NMatched: 0,
				NFailed:  0,
			},
		},
		NCalled:       0,
		ReferenceType: "",
	}

	Set_rules = []*Ruleset_t{&ThenCallRuleset, &ThenCall2Ruleset}

	CruxCacheTest = &Cache{
		RulesetCache: RulesetCache_t{"Remiges": PerRealm_t{"tnt": PerApp_t{12: PerSlice_t{Workflows: map[ClassName_t][]*Ruleset_t{"sale": Set_rules}}}}},
	}
)

func readJsonFromFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("testFile path is not exist")
	}
	defer file.Close()
	jsonData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
