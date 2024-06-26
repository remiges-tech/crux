package workflow_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/server/workflow"
	"github.com/remiges-tech/crux/types"
	"github.com/stretchr/testify/require"
)

const (
	TestWorkflowGet_test_1 = "ERROR_1- slice validation"
	TestWorkflowGet_test_2 = "SUCCESS_2- get workflow by valid req"
	TestWorkflowGet_test_3 = "Failed_3- get workflow by invalid req"
)

type TestCasesStruct struct {
	name             string
	requestPayload   wscutils.Request
	expectedHttpCode int
	testJsonFile     string
	expectedResult   *wscutils.Response
}

func TestWorkflowGet(t *testing.T) {
	testCases := workflowGetTestCase()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(server.MarshalJson(tc.requestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/workflowget", payload)
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

func workflowGetTestCase() []TestCasesStruct {
	var sliceStr int32 = 13
	app := "uccapp"
	class := "ucc"
	tname := "ucc_user_cr"
	var slice int32 = -1
	schemaNewTestcase := []TestCasesStruct{
		// 1st test case
		{
			name: TestWorkflowGet_test_1,
			requestPayload: wscutils.Request{
				Data: workflow.WorkflowGetReq{
					Slice: slice,
				},
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedResult: &wscutils.Response{
				Status: wscutils.ErrorStatus,
				Data:   nil,
				Messages: []wscutils.ErrorMessage{
					{
						MsgID:   102,
						ErrCode: "greater",
						Field:   &types.SLICE,
					}, {
						MsgID:   101,
						ErrCode: "required",
						Field:   &types.APP,
					}, {
						MsgID:   101,
						ErrCode: "required",
						Field:   &types.CLASS,
					},
					{
						MsgID:   101,
						ErrCode: "required",
						Field:   &types.NAME,
					},
				},
			},
		},

		// 2nd test case
		{
			name: TestWorkflowGet_test_2,
			requestPayload: wscutils.Request{
				Data: workflow.WorkflowGetReq{
					Slice: sliceStr,
					App:   app,
					Class: class,
					Name:  tname,
				},
			},

			expectedHttpCode: http.StatusOK,
			testJsonFile:     "./data/workflow_get_response.json",
		},
		// 3nd test case
		{
			name: TestWorkflowGet_test_3,
			requestPayload: wscutils.Request{
				Data: workflow.WorkflowGetReq{
					Slice: sliceStr,
					App:   "xyz",
					Class: class,
					Name:  tname,
				},
			},
			expectedHttpCode: http.StatusBadRequest,
			testJsonFile:     "./data/workflow_get_failed_response.json",
		},
	}
	return schemaNewTestcase
}
