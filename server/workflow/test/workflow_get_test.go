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
			payload := bytes.NewBuffer(types.MarshalJson(tc.requestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/workflowget", payload)
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

func workflowGetTestCase() []TestCasesStruct {
	var sliceStr int32 = 2
	app := "retailbank"
	class := "members"
	tname := "goldstatus"
	var slice int32 = -1
	schemaNewTestcase := []TestCasesStruct{
		// 1st test case
		{
			name: "ERROR- slice validation",
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
			name: "SUCCESS- get workflow by valid req ",
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
			name: "Failed- get workflow by invalid req ",
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
