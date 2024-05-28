package crux

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/stretchr/testify/require"
)

const (
	schemaCachePath     = "./data/schema_sample.json"
	schemaCacheThencall = "./data/schema_thencall.json"
	TestTracing_1       = "1. L0: with trace_level 0"
	TestTracing_2       = "2. L1: Match rule with return"
	TestTracing_3       = "3. L1: Match 1st rule & other mismatch"
	TestTracing_4       = "4. L2: Then_Call"
	TestTracing_5       = "5. L1: Else_Call & Return call"
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
			err = json.Unmarshal(ruleSchemasCache, &schema)
			require.NoError(t, err)
		}
		// setting wfschema
		cruxCache.Ctx = context.Background()
		cruxCache.Query = &sqlc.Queries{}
		cruxCache.RulesetCache = RulesetCache_t{"Remiges": PerRealm_t{"tnt": PerApp_t{12: PerSlice_t{Workflows: map[ClassName_t][]*Ruleset_t{"sale": set_rules}}}}}
		cruxCache.SchemaCache = SchemaCache_t{"Remiges": PerRealm_t{"tnt": PerApp_t{12: PerSlice_t{WFSchema: map[ClassName_t]Schema_t{"sale": *schema}}}}}
		t.Run(tc.TestCaseName, func(t *testing.T) {
			var expected []byte
			_, _, err, trace := DoMatch(tc.entity, tc.ruleset, schema, ActionSet{}, map[string]struct{}{}, emptyTrace, tc.trace_lev, cruxCache)
			traceByt, _ := json.Marshal(trace)
			fmt.Println("<<<<<<<<<<<<<<<<<<< traceByt:", string(traceByt))
			require.NoError(t, err)
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
			TestCaseName:     TestTracing_1,
			entity:           entity_return,
			ruleset:          &sampleRuleset,
			trace_lev:        0,
			ruleSchemasCache: schemaCachePath,
			ExpectedTrace:    emptyTrace,
		},

		// 2nd test case
		{
			TestCaseName:       TestTracing_2,
			entity:             entity_return,
			ruleset:            &sampleRuleset,
			trace_lev:          1,
			ruleSchemasCache:   schemaCachePath,
			ExpectedResultFile: "./data/test_2_expected.json",
		},

		// 3rd test case
		{
			TestCaseName:       TestTracing_3,
			entity:             entity_match,
			ruleset:            &sampleRuleset,
			trace_lev:          1,
			ruleSchemasCache:   schemaCachePath,
			ExpectedResultFile: "./data/test_3_expected.json",
		},

		// 4th test case
		{
			TestCaseName:       TestTracing_4,
			entity:             entity_thencall,
			ruleset:            &thenCallRuleset,
			trace_lev:          2,
			ruleSchemasCache:   schemaCacheThencall,
			ExpectedResultFile: "./data/test_4_expected.json",
		},

		// 5th test case : Else_Call
		{
			TestCaseName:       TestTracing_5,
			entity:             entity_elsecall,
			ruleset:            &elseCallRuleset,
			trace_lev:          1,
			ruleSchemasCache:   schemaCachePath,
			ExpectedResultFile: "./data/test_5_expected.json",
		},
	}
	return tracingTestcases
}
