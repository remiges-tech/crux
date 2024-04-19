package breruleset_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	breruleset "github.com/remiges-tech/crux/server/BRERuleset"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestBRERuleSetNew(t *testing.T) {
	var payload *bytes.Buffer
	testCases := RuleSetNewTestcase()
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
			req, err := http.NewRequest(http.MethodPost, "/brerulesetnew", payload)
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

func RuleSetNewTestcase() []testutils.TestCasesStruct {
	rulesetNewTestcase := []testutils.TestCasesStruct{

		{
			Name: "err- standard validation",
			RequestPayload: wscutils.Request{
				Data: breruleset.RuleSetNew{},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "./data/rulesetnew_standard_err.json",
		},
		{
			Name:             "Success - valid response",
			PayloadFile:      "./data/rulesetnew_payload.json",
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     nil,
				Messages: nil,
			},
		},
	}
	return rulesetNewTestcase
}
