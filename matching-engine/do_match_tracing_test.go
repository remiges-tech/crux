package crux

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	TestTracing_1 = "L1 trace: only `match` and `mismatch`"
	TestTracing_2 = "wfinstancenew matching"
)

type tracingTestCasesStruct struct {
	entity  Entity
	ruleset *Ruleset_t
	// TestJsonFile       string
	TestCaseName       string
	entityFilePath     string
	rulesetFilePath    string
	ruleSchemasCache   *Schema_t
	ExpectedResultFile string
	Url                string
}

var (
	sampleEntity = Entity{
		Realm: "BSE",
		App:   "starmf",
		Slice: 12,
		Class: "ucctest",
		Attrs: map[string]string{
			"mode":       "demat",
			"step":       "start",
			"stepfailed": "false",
		},
	}
	sampleEntityWfinstancenew = Entity{
		Realm: "BSE",
		App:   "uccapp",
		Slice: 12,
		Class: "ucc",
		Attrs: map[string]string{
			"mode":       "demat",
			"step":       "start",
			"stepfailed": "false",
		},
	}
	sampleRuleset = Ruleset_t{
		Id:      10,
		Class:   "ucctest",
		SetName: "ucctest",
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
)

func TestDoMatchTracee(t *testing.T) {
	testCases := testcase()
	for _, tc := range testCases {
		if tc.rulesetFilePath != "" {
			var temp_ruleset Ruleset_t
			rulesetFile, err := readJsonFromFile(tc.rulesetFilePath)
			require.NoError(t, err)

			err = json.Unmarshal(rulesetFile, &temp_ruleset)
			require.NoError(t, err)
			tc.ruleset = &temp_ruleset
		}

		t.Run(tc.TestCaseName, func(t *testing.T) {
			_, _, err, trace := DoMatch(tc.entity, tc.ruleset, tc.ruleSchemasCache, ActionSet{}, map[string]struct{}{}, Trace_t{})
			require.NoError(t, err)

			traceByt, _ := json.Marshal(trace)
			fmt.Println("<<<<<<<<<<<<<<<<<<< traceByt:", string(traceByt))

			expected, err := readJsonFromFile(tc.ExpectedResultFile)
			require.NoError(t, err)
			actual, err := json.Marshal(trace)
			// fmt.Println(">>>>>> actual:", string(actual))
			require.NoError(t, err)
			require.JSONEq(t, string(expected), string(actual))
		})
	}

}

func testcase() []tracingTestCasesStruct {
	tracingTestcases := []tracingTestCasesStruct{
		// 1st test case
		{
			TestCaseName:     TestTracing_1,
			entity:           sampleEntity,
			ruleset:          &sampleRuleset,
			ruleSchemasCache: &Schema_t{},
			// TestJsonFile:       "./data/action_.json",
			ExpectedResultFile: "./data/expected_trace.json",
		},
		// 2nd test case
		{
			TestCaseName: TestTracing_2,
			entity:       sampleEntityWfinstancenew,
			// ruleset:          &sampleRuleset,
			rulesetFilePath:  "./data/ruleset_wfinstancenew.json",
			ruleSchemasCache: &Schema_t{},
			// TestJsonFile:       "./data/action_.json",
			ExpectedResultFile: "./data/expected_trace2.json",
		},
	}
	return tracingTestcases
}

func readJsonFromFile(filepath string) ([]byte, error) {
	// var err error
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
