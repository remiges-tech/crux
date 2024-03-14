package app_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server/app"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestAppUpdate(t *testing.T) {
	testCases := appUpdateTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/appupdate", payload)
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

func appUpdateTestcase() []testutils.TestCasesStruct {

	appUpdateTestcase := []testutils.TestCasesStruct{
		{
			Name: "ERROR: empty request",
			RequestPayload: wscutils.Request{
				Data: app.GetAppUpdateRequest{
					Name:        "",
					Description: "",
				},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "../test/data/empty_request.json",
		},
		{
			Name: "ERROR: Invalid App Name",
			RequestPayload: wscutils.Request{Data: app.GetAppNewRequest{
				Name:        "SBI Bank",
				Description: "SBI bank services",
			}},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "../test/data/invalid_app_name.json",
		},
		{
			Name: "SUCCESS: Valid request",
			RequestPayload: wscutils.Request{
				Data: app.GetAppNewRequest{
					Name:        "retailBANK",
					Description: "Retail bank services",
				},
			},

			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     nil,
				Messages: nil,
			},
		},
		{
			Name: "ERROR: Name does not exist",
			RequestPayload: wscutils.Request{
				Data: app.GetAppNewRequest{
					Name:        "SBIs",
					Description: "SBI bank services",
				},
			},

			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "../test/data/name_not_exist.json"},
	}
	return appUpdateTestcase
}
