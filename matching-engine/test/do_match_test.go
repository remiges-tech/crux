package crux_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/remiges-tech/crux/db/sqlc-gen"
	crux "github.com/remiges-tech/crux/matching-engine"
	"github.com/remiges-tech/crux/server"
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
		var schema *crux.Schema_t
		if tc.RulesetFilePath != "" {
			var tm crux.Ruleset_t
			rulesetFile, err := server.ReadJsonFromFile(tc.RulesetFilePath)
			require.NoError(t, err)
			err = json.Unmarshal(rulesetFile, &tm)
			require.NoError(t, err)
			tc.Ruleset_t = &tm
		}
		if tc.RuleSchemasCache != "" {
			ruleSchemasCache, err := server.ReadJsonFromFile(tc.RuleSchemasCache)
			require.NoError(t, err)
			err = json.Unmarshal(ruleSchemasCache, &schema)
			require.NoError(t, err)
		}
		// setting wfschema
		crux.CruxCacheTest.Ctx = context.Background()
		crux.CruxCacheTest.Query = &sqlc.Queries{}
		crux.CruxCacheTest.RulesetCache = crux.RulesetCache_t{"Remiges": crux.PerRealm_t{"tnt": crux.PerApp_t{12: crux.PerSlice_t{Workflows: map[crux.ClassName_t][]*crux.Ruleset_t{"sale": crux.Set_rules}}}}}
		crux.CruxCacheTest.SchemaCache = crux.SchemaCache_t{"Remiges": crux.PerRealm_t{"tnt": crux.PerApp_t{12: crux.PerSlice_t{WFSchema: map[crux.ClassName_t]crux.Schema_t{"sale": *schema}}}}}
		t.Run(tc.TestCaseName, func(t *testing.T) {
			var expected []byte
			_, _, err, trace := crux.DoMatch(tc.Entity, tc.Ruleset_t, schema, crux.ActionSet{}, map[string]struct{}{}, crux.EmptyTrace, tc.Trace_lev, crux.CruxCacheTest)
			require.NoError(t, err)
			actual, err := json.Marshal(trace)
			require.NoError(t, err)
			if tc.ExpectedResultFile != "" {
				expected, err = server.ReadJsonFromFile(tc.ExpectedResultFile)
				require.NoError(t, err)
			} else {
				expected, err = json.Marshal(tc.ExpectedTrace)
				require.NoError(t, err)
			}
			require.JSONEq(t, string(expected), string(actual))
		})
	}

}

func testcase() []crux.TracingTestCasesStruct {
	tracingTestcases := []crux.TracingTestCasesStruct{
		// 1st test case
		{
			TestCaseName:     TestTracing_1,
			Entity:           crux.Entity_return,
			Ruleset_t:        &crux.SampleRuleset,
			Trace_lev:        0,
			ExpectedTrace:    crux.EmptyTrace,
			RuleSchemasCache: schemaCachePath,
		},

		// 2nd test case
		{
			TestCaseName:       TestTracing_2,
			Entity:             crux.Entity_return,
			Ruleset_t:          &crux.SampleRuleset,
			Trace_lev:          1,
			RuleSchemasCache:   schemaCachePath,
			ExpectedResultFile: "./data/test_2_expected.json",
		},

		// 3rd test case
		{
			TestCaseName:       TestTracing_3,
			Entity:             crux.Entity_match,
			Ruleset_t:          &crux.SampleRuleset,
			Trace_lev:          1,
			RuleSchemasCache:   schemaCachePath,
			ExpectedResultFile: "./data/test_3_expected.json",
		},

		// 4th test case
		{
			TestCaseName:       TestTracing_4,
			Entity:             crux.Entity_thencall,
			Ruleset_t:          &crux.ThenCallRuleset,
			Trace_lev:          2,
			RuleSchemasCache:   schemaCacheThencall,
			ExpectedResultFile: "./data/test_4_expected.json",
		},

		// 5th test case : Else_Call
		{
			TestCaseName:       TestTracing_5,
			Entity:             crux.Entity_elsecall,
			Ruleset_t:          &crux.ElseCallRuleset,
			Trace_lev:          1,
			RuleSchemasCache:   schemaCachePath,
			ExpectedResultFile: "./data/test_5_expected.json",
		},
	}
	return tracingTestcases
}
