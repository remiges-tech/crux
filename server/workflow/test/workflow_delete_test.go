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

const (
	TestWorkflowDelete_test_1 = "ERROR- slice validation"
	TestWorkflowDelete_test_2 = "ERROR- if record not exists"
	TestWorkflowDelete_test_3 = "SUCCESS- delete workflow by valid req"
)

func TestWorkflowDelete(t *testing.T) {
	testCases := workflowDeleteTestCase()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(types.MarshalJson(tc.requestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/workflowdelete", payload)
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

func workflowDeleteTestCase() []TestCasesStruct {
	var sliceStr int32 = 2
	app := "retailBANK"
	class := "members"
	tname := "tempset"
	wrongName := "tempse"
	var slice int32 = -1
	schemaNewTestcase := []TestCasesStruct{
		// 1st test case
		{
			name: TestWorkflowDelete_test_1,
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

		// 2nd error: if given payload record not exists
		{
			name: TestWorkflowDelete_test_2,
			requestPayload: wscutils.Request{
				Data: workflow.WorkflowGetReq{
					Slice: sliceStr,
					App:   app,
					Class: class,
					Name:  wrongName,
				},
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedResult: &wscutils.Response{
				Status: wscutils.ErrorStatus,
				Data:   nil,
				Messages: []wscutils.ErrorMessage{
					{
						MsgID:   1006,
						ErrCode: "invalid_request",
					},
				},
			},
		},

		// 3rd test case
		{
			name: TestWorkflowDelete_test_3,
			requestPayload: wscutils.Request{
				Data: workflow.WorkflowGetReq{
					Slice: sliceStr,
					App:   app,
					Class: class,
					Name:  tname,
				},
			},
			expectedHttpCode: http.StatusOK,
			expectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     nil,
				Messages: nil,
			},
		},
	}
	return schemaNewTestcase
}
