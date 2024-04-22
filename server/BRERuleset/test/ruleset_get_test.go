package breruleset_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server"
	breruleset "github.com/remiges-tech/crux/server/BRERuleset"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestBRERuleSetGet(t *testing.T) {
	testCases := RuleSetGetTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(server.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/brerulesetget", payload)
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

func RuleSetGetTestcase() []testutils.TestCasesStruct {
	rulesetNewTestcase := []testutils.TestCasesStruct{

		{
			Name: "err- standard validation",
			RequestPayload: wscutils.Request{
				Data: breruleset.RuleSetGetReq{},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "./data/rulesetget_standard_err.json",
		},
		{
			Name: "Success - valid response",
			RequestPayload: wscutils.Request{
				Data: breruleset.RuleSetGetReq{
					Slice: (int32(11)),
					App:   "amazon",
					Class: "inventoryitems",
					Name:  "amazonruleset",
				},
			},
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/rulesetget_valid_response.json",
		},
	}
	return rulesetNewTestcase
}
