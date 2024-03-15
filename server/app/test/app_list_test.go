package app_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

const (
	AppListTest_1 = "SUCCESS: valid request"
)

func TestAppList(t *testing.T) {
	testCases := appListTestCases()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/applist", payload)
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

func appListTestCases() []testutils.TestCasesStruct {

	appListTestCases := []testutils.TestCasesStruct{
		// test 1
		{
			Name:             AppListTest_1,
			RequestPayload:   wscutils.Request{},
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "../test/data/get_app_list.json",
		},
	}
	return appListTestCases
}
