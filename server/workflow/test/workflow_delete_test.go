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
	TestWorkflowDelete_test_1 = "ERROR_1- slice validation"
	TestWorkflowDelete_test_2 = "ERROR_2- if record not exists"
	TestWorkflowDelete_test_3 = "SUCCESS_3- delete workflow by valid req"
)

func TestWorkflowDelete(t *testing.T) {
	testCases := workflowDeleteTestCase()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(server.MarshalJson(tc.requestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/workflowdelete", payload)
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

func workflowDeleteTestCase() []TestCasesStruct {
	var sliceStr int32 = 13
	app := "uccapp"
	class := "ucc"
	tname := "ucc_user_cr_inactive"
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
						ErrCode: "no_record_found",
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
