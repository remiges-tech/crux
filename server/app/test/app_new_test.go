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

func TestAppNew(t *testing.T) {
	testCases := appNewTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/appnew", payload)
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

func appNewTestcase() []testutils.TestCasesStruct {

	appNewTestcase := []testutils.TestCasesStruct{
		{
			Name: "ERROR: empty request",
			RequestPayload: wscutils.Request{
				Data: app.GetAppNewRequest{
					Name:        "",
					Description: "",
				},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "../test/data/app_new_empty_request_res.json",
		},
		{
			Name: "ERROR: Invalid App Name",
			RequestPayload: wscutils.Request{Data: app.GetAppNewRequest{
				Name:        "_SBI Bank",
				Description: "SBI bank services",
			}},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "../test/data/app_new_invalidapp_name_res.json",
		},
		{
			Name: "SUCCESS: Valid request",
			RequestPayload: wscutils.Request{
				Data: app.GetAppNewRequest{
					Name:        "SBI",
					Description: "SBI bank services",
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
			Name: "ERROR: App already exist",
			RequestPayload: wscutils.Request{
				Data: app.GetAppNewRequest{
					Name:        "SBI",
					Description: "SBI bank services",
				},
			},

			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "../test/data/app_new_app_exists_res.json"},
	}
	return appNewTestcase
}
