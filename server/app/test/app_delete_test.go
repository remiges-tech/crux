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

func TestAppDelete(t *testing.T) {
	testCases := appDeleteTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, tc.Url, payload)
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

func appDeleteTestcase() []testutils.TestCasesStruct {

	appDeleteTestcase := []testutils.TestCasesStruct{
		{
			Name:             "ERROR: Invalid App Name",
			Url:              "/appdelete/SBI Bank",
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "../test/data/invalid_app_name.json",
		},
		{
			Name:             "SUCCESS: Valid request",
			Url:              "/appdelete/nedBank1",
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     nil,
				Messages: nil,
			},
		},
		{
			Name:             "ERROR: Name does not exist",
			Url:              "/appdelete/SBI",
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "../test/data/name_not_exist.json",
		},
		{
			Name:             "ERROR: Name has assets attached to it",
			Url:              "/appdelete/retailbank",
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "../test/data/non_empty.json",
		},
	}
	return appDeleteTestcase
}
