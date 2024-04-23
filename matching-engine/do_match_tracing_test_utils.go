package crux

import (
	"fmt"
	"io"
	"os"
)

type tracingTestCasesStruct struct {
	TestCaseName       string
	entity             Entity
	ruleset            *Ruleset_t
	entityFilePath     string
	rulesetFilePath    string
	trace_lev          int
	ExpectedResultFile string
	Url                string
	ruleSchemasCache   string
}

//***********************************  Entity  ****************************************

var (
	entity_test_1 = Entity{
		Realm: "Remiges",
		App:   "tnt",
		Slice: 12,
		Class: "finance",
		Attrs: map[string]string{
			"step": "step1",
			"mode": "demat",
			// 		"stepfailed": "false",

		},
	}

	// sampleEntityWfinstancenew = Entity{
	// 	Realm: "BSE",
	// 	App:   "uccapp",
	// 	Slice: 12,
	// 	Class: "ucc",
	// 	Attrs: map[string]string{
	// 		"mode":       "demat",
	// 		"step":       "start",
	// 		"stepfailed": "false",
	// 	},
	// }
)

// *********************************  Ruleset_t  ****************************************
var (
	sampleRuleset = Ruleset_t{
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
						"nextstep": "step3",
						// "done":     "true",
					},
					// ThenCall: "finance_2",
					// References: []*Ruleset_t{
					// 	&Ruleset_t{
					// 		Id:      28,
					// 		Class:   "finance_2",
					// 		SetName: "finance_2",
					// 		Rules: []Rule_t{{
					// 			RulePatterns: []RulePatternBlock_t{
					// 				{
					// 					Attr: "step",
					// 					Op:   "eq",
					// 					Val:  "step3",
					// 				},
					// 				{
					// 					Attr: "mode",
					// 					Op:   "eq",
					// 					Val:  "demat",
					// 				},
					// 			},
					// 			RuleActions: RuleActionBlock_t{
					// 				Task: []string{},
					// 				Properties: map[string]string{
					// 					"done": "true",
					// 				},
					// 			},
					// 			NMatched: 0,
					// 			NFailed:  0,
					// 		}},
					// 		NCalled:       0,
					// 		ReferenceType: "",
					// 	},
					// },
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
					ElseCall: "ruleset_name-will be here", // elsecall
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
					ThenCall: "true", // thencall
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

	// rulSet = &Ruleset_t{
	// 	Id:      28,
	// 	Class:   "finance_2",
	// 	SetName: "finance_2",
	// 	Rules: []Rule_t{{
	// 		RulePatterns: []RulePatternBlock_t{
	// 			{
	// 				Attr: "step",
	// 				Op:   "eq",
	// 				Val:  "step3",
	// 			},
	// 			{
	// 				Attr: "mode",
	// 				Op:   "eq",
	// 				Val:  "demat",
	// 			},
	// 		},
	// 		RuleActions: RuleActionBlock_t{
	// 			Task: []string{},
	// 			Properties: map[string]string{
	// 				"done": "true",
	// 			},
	// 		},
	// 		NMatched: 0,
	// 		NFailed:  0,
	// 	}},
	// 	NCalled:       0,
	// 	ReferenceType: "",
	// }
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
