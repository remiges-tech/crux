package workflow_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/server/workflow"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestWorkflowList(t *testing.T) {
	testCases := WorkflowListTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(server.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/workflowlist", payload)
			require.NoError(t, err)

			r.ServeHTTP(res, req)

			require.Equal(t, tc.ExpectedHttpCode, res.Code)

			if tc.ExpectedResult != nil {
				expectedData := server.MarshalJson(tc.ExpectedResult)
				actualData := res.Body.Bytes()
				compareJSON(t, expectedData, actualData)
			} else {
				expectedData, err := server.ReadJsonFromFile(tc.TestJsonFile)
				require.NoError(t, err)
				actualData := res.Body.Bytes()
				fmt.Println("expected>>>>>>>>>>>>>>>>>>>", string(expectedData))
				fmt.Println("actual>>>>>>>>>>>>>>>>>>>", string(actualData))

				compareJSON(t, expectedData, actualData)
			}
		})
	}
}

func WorkflowListTestcase() []testutils.TestCasesStruct {
	rulesetListTestcase := []testutils.TestCasesStruct{
		{
			Name: "success - no parameter present",
			RequestPayload: wscutils.Request{
				Data: workflow.WorkflowListReq{},
			},
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/workflow_list_empty_req_response.json",
		},
		{
			Name: "Success - valid response - root capability",
			RequestPayload: wscutils.Request{
				Data: workflow.WorkflowListReq{
					Slice: int32(12),
					App:   "fundify",
				},
			},
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/workflow_list_success_response.json",
		},
	}
	return rulesetListTestcase
}

func removeField(data map[string]interface{}, field string) {
	for _, v := range data {
		switch vv := v.(type) {
		case map[string]interface{}:
			removeField(vv, field)
		case []interface{}:
			for _, u := range vv {
				if umap, ok := u.(map[string]interface{}); ok {
					removeField(umap, field)
				}
			}
		}
	}
	delete(data, field)
}

func canonicalizeJSON(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			v[key] = canonicalizeJSON(value)
		}
		return v
	case []interface{}:
		for i, value := range v {
			v[i] = canonicalizeJSON(value)
		}
		sort.Slice(v, func(i, j int) bool {
			return fmt.Sprintf("%v", v[i]) < fmt.Sprintf("%v", v[j])
		})
		return v
	default:
		return v
	}
}

func compareJSON(t *testing.T, expected, actual []byte) {
	var expectedMap, actualMap map[string]interface{}

	if err := json.Unmarshal(expected, &expectedMap); err != nil {
		t.Fatalf("Error unmarshaling expected JSON: %v", err)
	}
	if err := json.Unmarshal(actual, &actualMap); err != nil {
		t.Fatalf("Error unmarshaling actual JSON: %v", err)
	}

	fieldsToRemove := []string{"createdat", "editedat"}

	for _, field := range fieldsToRemove {
		removeField(expectedMap, field)
		removeField(actualMap, field)
	}

	expectedMap = canonicalizeJSON(expectedMap).(map[string]interface{})
	actualMap = canonicalizeJSON(actualMap).(map[string]interface{})

	if !reflect.DeepEqual(expectedMap, actualMap) {
		expectedStr, _ := json.MarshalIndent(expectedMap, "", "  ")
		actualStr, _ := json.MarshalIndent(actualMap, "", "  ")
		t.Errorf("Expected JSON does not match actual JSON\nExpected: %s\nActual: %s", string(expectedStr), string(actualStr))
	}
}
