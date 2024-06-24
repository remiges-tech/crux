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

func TestBRERuleSetList(t *testing.T) {
	testCases := RuleSetListTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(server.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/brerulesetlist", payload)
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

func RuleSetListTestcase() []testutils.TestCasesStruct {
	rulesetListTestcase := []testutils.TestCasesStruct{

		{
			Name: "success - no parameter present",
			RequestPayload: wscutils.Request{
				Data: breruleset.RuleSetListReq{},
			},
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/rulesetlist_resp.json",
		},
		{
			Name: "Success - valid response - root capability",
			RequestPayload: wscutils.Request{
				Data: breruleset.RuleSetListReq{
					Slice: (int32(11)),
					App:   "amazon",
				},
			},
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/rulesetlist_resp.json",
		},
	}
	return rulesetListTestcase
}
