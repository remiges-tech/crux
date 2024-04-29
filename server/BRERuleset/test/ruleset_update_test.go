package breruleset_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestRuleSetUpdate(t *testing.T) {
	var payload *bytes.Buffer
	testCases := RuleSetUpdateTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.PayloadFile != "" {
				jsonData, err := testutils.ReadJsonFromFile(tc.PayloadFile)
				require.NoError(t, err)
				payload = bytes.NewBuffer(jsonData)
			} else {
				payload = bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))
			}

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/BRErulesetUpdate", payload)
			require.NoError(t, err)

			r.ServeHTTP(res, req)

			require.Equal(t, tc.ExpectedHttpCode, res.Code)
			if tc.ExpectedResult != nil {
				jsonData := testutils.MarshalJson(tc.ExpectedResult)
				require.JSONEq(t, string(jsonData), res.Body.String())
			} else {
				jsonData, err := testutils.ReadJsonFromFile(tc.TestJsonFile)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonData), res.Body.String())
			}
		})
	}

}

func RuleSetUpdateTestcase() []testutils.TestCasesStruct {
	schemaNewTestcase := []testutils.TestCasesStruct{
		{
			Name: "err- binding_json_error",
			RequestPayload: wscutils.Request{
				Data: nil,
			},

			ExpectedHttpCode: http.StatusBadRequest,
			ExpectedResult: &wscutils.Response{
				Status: wscutils.ErrorStatus,
				Data:   nil,
				Messages: []wscutils.ErrorMessage{
					{
						MsgID:   0,
						ErrCode: "",
					},
				},
			},
		},
		{
			Name:             "err- standard validation",
			PayloadFile:      "./data/ruleset_new_validation_payload.json",
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "./data/ruleset_new_validation_error.json",
		},
		{
			Name:             "err- custom validation",
			PayloadFile:      "./data/ruleset_new_custom_validation_payload.json",
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "./data/ruleset_new_custom_validation_error.json",
		},
		{
			Name:             "Success- update ruleset",
			PayloadFile:      "./data/ruleset_new_success_payload.json",
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     nil,
				Messages: nil,
			},
		},
	}
	return schemaNewTestcase
}
