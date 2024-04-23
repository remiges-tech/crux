package crux

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	TestTracing_1 = "L1: Match rule ( 0 , 1) Return ( 3 )"
	TestTracing_2 = "wfinstancenew matching"
)

func TestDoMatchTracee(t *testing.T) {
	testCases := testcase()
	for _, tc := range testCases {
		var schema *Schema_t
		if tc.rulesetFilePath != "" {
			var tm Ruleset_t
			rulesetFile, err := readJsonFromFile(tc.rulesetFilePath)
			require.NoError(t, err)

			err = json.Unmarshal(rulesetFile, &tm)
			require.NoError(t, err)
			tc.ruleset = &tm
		}

		if tc.ruleSchemasCache != "" {
			ruleSchemasCache, err := readJsonFromFile(tc.ruleSchemasCache)
			require.NoError(t, err)

			var tmp interface{}
			err = json.Unmarshal(ruleSchemasCache, &tmp)
			require.NoError(t, err)
		}

		t.Run(tc.TestCaseName, func(t *testing.T) {
			_, _, err, trace := DoMatch(tc.entity, tc.ruleset, schema, ActionSet{}, map[string]struct{}{}, Trace_t{}, tc.trace_lev)
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
			TestCaseName:       TestTracing_1,
			entity:             entity_test_1,
			ruleset:            &sampleRuleset,
			trace_lev:          1,
			ruleSchemasCache:   "./data/schema_sample.json",
			ExpectedResultFile: "./data/expected_trace_test_1.json",
		},

		// 2nd test case
		// {
		// 	TestCaseName: TestTracing_2,
		// 	entity:       sampleEntityWfinstancenew,
		// 	// ruleset:          &sampleRuleset,
		// 	rulesetFilePath:    "./data/ruleset_wfinstancenew.json",
		// 	ruleSchemasCache:   "./data/schema_t.json",
		// 	ExpectedResultFile: "./data/expected_trace2.json",
		// },
	}
	return tracingTestcases
}
