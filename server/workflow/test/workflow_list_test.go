package workflow_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server/workflow"
	"github.com/remiges-tech/crux/types"
	"github.com/stretchr/testify/require"
)

func TestWorkflowList(t *testing.T) {
	testCases := workflowListTestCase()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(types.MarshalJson(tc.requestPayload))

			res := httptest.NewRecorder()
			if tc.name == "SUCCESS- with app name but HasRulesetRights = false & HasRootCapabilities = true" {
				workflow.TRIGGER = true
			}
			req, err := http.NewRequest(http.MethodPost, "/workflowlist", payload)
			require.NoError(t, err)
			r.ServeHTTP(res, req)
			require.Equal(t, tc.expectedHttpCode, res.Code)
			if tc.expectedResult != nil {
				jsonData := types.MarshalJson(tc.expectedResult)
				require.JSONEq(t, string(jsonData), res.Body.String())
			} else {
				jsonData, err := types.ReadJsonFromFile(tc.testJsonFile)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonData), res.Body.String())
			}
		})
	}

}

func workflowListTestCase() []TestCasesStruct {
	var sliceStr int32 = 2
	app := "retailbank"
	class := "members"
	tname := "goldstatus"
	isActive := true
	schemaNewTestcase := []TestCasesStruct{
		// 1st test case
		{
			name: "ERROR- App + HasRootCapabilities()= false & HasRulesetRights()= false",
			requestPayload: wscutils.Request{
				Data: workflow.WorkflowListReq{
					Slice:      &sliceStr,
					App:        &app,
					Class:      &class,
					Name:       &tname,
					IsActive:   &isActive,
					IsInternal: &isActive,
				},
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedResult: &wscutils.Response{
				Status: wscutils.ErrorStatus,
				Data:   nil,
				Messages: []wscutils.ErrorMessage{
					{
						MsgID:   1010,
						ErrCode: "Unauthorized",
					},
				},
			},
		},

		// 2nd test case
		{
			name: "SUCCESS- No app name but HasRulesetRights = false & HasRootCapabilities = false",
			requestPayload: wscutils.Request{
				Data: workflow.WorkflowListReq{
					Slice:      &sliceStr,
					Class:      &class,
					Name:       &tname,
					IsActive:   &isActive,
					IsInternal: &isActive,
				},
			},
			expectedHttpCode: http.StatusOK,
			testJsonFile:     "./data/workflow_list_no_app_response.json",
		},

		// 3rd test case
		{
			name: "SUCCESS- with app name but HasRulesetRights = false & HasRootCapabilities = true",
			requestPayload: wscutils.Request{
				Data: workflow.WorkflowListReq{},
			},
			expectedHttpCode: http.StatusOK,
			testJsonFile:     "./data/workflow_list_with_app_response.json",
		},
	}
	return schemaNewTestcase
}
