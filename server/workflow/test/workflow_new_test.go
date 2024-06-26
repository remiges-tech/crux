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

func TestWorkflowNew(t *testing.T) {
	testCases := WorkflowNewTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(server.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/workflownew", payload)
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

func WorkflowNewTestcase() []testutils.TestCasesStruct {
	workflowNewTestcase := []testutils.TestCasesStruct{

		{
			Name: "err- standard validation",
			RequestPayload: wscutils.Request{
				Data: workflow.WorkflowNewRequest{},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "./data/workflow_new_validation_error.json",
		},
		{
			Name: "Success - valid response",
			RequestPayload: wscutils.Request{
				Data: workflow.WorkflowNewRequest{
					Slice:      11,
					App:        "myntra",
					Class:      "inventoryitems",
					Name:       "myntraruleset",
					IsInternal: true,
					Flowrules: []crux.Rule_t{
						{
							RuleActions: crux.RuleActionBlock_t{
								Task: []string{"cat"},
								Properties: map[string]string{
									"nextstep": "done",
								},
							},
							RulePatterns: []crux.RulePatternBlock_t{
								{
									Op:   "eq",
									Val:  "textbook",
									Attr: "cat",
								},
								{
									Op:   "ge",
									Val:  5000,
									Attr: "mrp",
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
	return workflowNewTestcase
}
