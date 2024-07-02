package workflow_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	crux "github.com/remiges-tech/crux/matching-engine"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/server/workflow"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestWorkflowUpdate(t *testing.T) {
	testCases := WorkflowUpdateTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(server.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPut, "/workflowupdate", payload)
			require.NoError(t, err)

			r.ServeHTTP(res, req)

			require.Equal(t, tc.ExpectedHttpCode, res.Code)
			if tc.ExpectedResult != nil {
				jsonData := server.MarshalJson(tc.ExpectedResult)
				require.JSONEq(t, string(jsonData), res.Body.String())
			} else {
				jsonData, err := server.ReadJsonFromFile(tc.TestJsonFile)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonData), res.Body.String())
			}
		})
	}
}

func WorkflowUpdateTestcase() []testutils.TestCasesStruct {
	workflowUpdateTestcase := []testutils.TestCasesStruct{

		{
			Name: "err- standard validation",
			RequestPayload: wscutils.Request{
				Data: workflow.WorkflowUpdate{},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "./data/workflow_update_validation_error.json",
		},
		{
			Name: "Success - valid response",
			RequestPayload: wscutils.Request{
				Data: workflow.WorkflowUpdate{
					Slice:      12,
					App:        "fundify",
					Class:      "ucc",
					Name:       "temps",
					Flowrules: []crux.Rule_t{
						{
							RuleActions: crux.RuleActionBlock_t{
								Task:[]string{"pan_verification", "pan_aadhaar_linking"},
								Properties: map[string]string{
									"nextstep": "kyc_done",
								},
							},
							RulePatterns: []crux.RulePatternBlock_t{
								{
									Op:   "eq",
									Val:  "start",
									Attr: "step",
								},
								{
									Op:   "eq",
									Val:  "broker",
									Attr: "member_type",
								},
								{
									Op:   "eq",
									Val:  "physical",
									Attr: "ucc_type",
								},
								{
									Op:   "eq",
									Val:  "individual",
									Attr: "tax_status_type",
								},
							},
							NFailed:  0,
							NMatched: 0,
						},
					},
				},
			},
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     nil,
				Messages: nil,
			},
		},
	}
	return workflowUpdateTestcase
}
