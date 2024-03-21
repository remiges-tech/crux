package auth_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestCapGet(t *testing.T) {
	testCases := testUserActivate()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, tc.Url, payload)
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

func testUserActivate() []testutils.TestCasesStruct {

	testUsrAct := []testutils.TestCasesStruct{
		{
			Name:             "ERROR: Invalid user_activate userId",
			Url:              "/capget/jhohns",
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "../test/data/invalid_request.json",
		},
		{
			Name:             "SUCCESS: Valid request",
			Url:              "/capget/Raj",
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "../test/data/cap_get_resp.json",
		},
	}
	return testUsrAct
}
