package crux

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	schemaCachePath = "./data/schema_sample.json"
	TestTracing_0   = "1. L0: with trace_level 0"
	TestTracing_1   = "2. L1: Match rule with return"
	TestTracing_2   = "3. L1: Match 1st rule & other mismatch"
	TestTracing_3   = "4. L1: Entity attributes are more than ruleset attrs"
)

var emptyTrace = Trace_t{}

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

			err = json.Unmarshal(ruleSchemasCache, &schema)
			require.NoError(t, err)
		}

		t.Run(tc.TestCaseName, func(t *testing.T) {
			var expected []byte
			_, _, err, trace := DoMatch(tc.entity, tc.ruleset, schema, ActionSet{}, map[string]struct{}{}, emptyTrace, tc.trace_lev)
			if tc.TestCaseName == TestTracing_3 {
				require.EqualError(t, err, "size of entity attrs is not equal to rulePattern") // if the entity attrs arre more
				return
			} else {
				require.NoError(t, err)
			}

			// traceByt, _ := json.Marshal(trace)
			// fmt.Println("<<<<<<<<<<<<<<<<<<< traceByt:", string(traceByt))

			actual, err := json.Marshal(trace)
			require.NoError(t, err)

			if tc.ExpectedResultFile != "" {
				expected, err = readJsonFromFile(tc.ExpectedResultFile)
				require.NoError(t, err)
			} else {
				expected, err = json.Marshal(tc.ExpectedTrace)
				require.NoError(t, err)
			}
			require.JSONEq(t, string(expected), string(actual))

		})
	}

}

func testcase() []tracingTestCasesStruct {
	tracingTestcases := []tracingTestCasesStruct{
		// 1st test case
		{
			TestCaseName:     TestTracing_0,
			entity:           entity_return,
			ruleset:          &sampleRuleset,
			trace_lev:        0,
			ruleSchemasCache: schemaCachePath,
			ExpectedTrace:    emptyTrace,
		},

		// 2nd test case
		{
			TestCaseName:       TestTracing_1,
			entity:             entity_return,
			ruleset:            &sampleRuleset,
			trace_lev:          1,
			ruleSchemasCache:   schemaCachePath,
			ExpectedResultFile: "./data/test_2_expected.json",
		},

		// 3rd test case
		{
			TestCaseName:       TestTracing_2,
			entity:             entity_match,
			ruleset:            &sampleRuleset,
			trace_lev:          1,
			ruleSchemasCache:   schemaCachePath,
			ExpectedResultFile: "./data/test_3_expected.json",
		},

		// 4th test case
		{
			TestCaseName: TestTracing_3,
			entity:       entity_attrs,
			ruleset:      &sampleRuleset,
			trace_lev:    1,
		},
	}
	return tracingTestcases
}
