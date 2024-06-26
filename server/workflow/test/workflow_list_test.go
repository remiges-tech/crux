package workflow_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/server/workflow"
	"github.com/stretchr/testify/require"
)

const (
	TestWorkflowList_test_1 = "ERROR_1- App + HasRootCapabilities()= false & HasRulesetRights()= false"
	TestWorkflowList_test_2 = "SUCCESS_2- No app name but HasRulesetRights = false & HasRootCapabilities = false"
	TestWorkflowList_test_3 = "SUCCESS_3- empty req with HasRulesetRights = false & HasRootCapabilities = true"
)

func TestWorkflowList(t *testing.T) {
	testCases := workflowListTestCase()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(server.MarshalJson(tc.requestPayload))

			res := httptest.NewRecorder()
			if tc.name == TestWorkflowList_test_3 {
				workflow.TRIGGER = true
			}
			req, err := http.NewRequest(http.MethodPost, "/workflowlist", payload)
			require.NoError(t, err)
			r.ServeHTTP(res, req)
			require.Equal(t, tc.expectedHttpCode, res.Code)
			if tc.expectedResult != nil {
				jsonData := server.MarshalJson(tc.expectedResult)
				require.JSONEq(t, string(jsonData), res.Body.String())
			} else {
				jsonData, err := server.ReadJsonFromFile(tc.testJsonFile)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonData), res.Body.String())
			}
		})
	}

}

func workflowListTestCase() []TestCasesStruct {
	var sliceStr int32 = 14
	app := "retailbank"
	class := "members"
	tname := "temp"
	isActive := true
	schemaNewTestcase := []TestCasesStruct{
		// 1st test case
		{
			name: TestWorkflowList_test_1,
			requestPayload: wscutils.Request{
				Data: workflow.WorkflowListReq{
					Slice:      sliceStr,
					App:        app,
					Class:      class,
					Name:       tname,
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
			name: TestWorkflowList_test_2,
			requestPayload: wscutils.Request{
				Data: workflow.WorkflowListReq{
					Slice:      sliceStr,
					Class:      class,
					Name:       tname,
					IsActive:   &isActive,
					IsInternal: &isActive,
				},
			},
			expectedHttpCode: http.StatusOK,
			testJsonFile:     "./data/workflow_list_no_app_response.json",
		},

		// 3rd test case
		{
			name: TestWorkflowList_test_3,
			requestPayload: wscutils.Request{
				Data: workflow.WorkflowListReq{},
			},
			expectedHttpCode: http.StatusOK,
			testJsonFile:     "./data/workflow_list_empty_req_response.json",
		},
	}
	return schemaNewTestcase
}
