package breschema_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server"
	breschema "github.com/remiges-tech/crux/server/BRESchema"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestBRESchemaDelete(t *testing.T) {
	testCases := breschemaDeleteTestCase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(server.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/breschemadelete", payload)
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

func breschemaDeleteTestCase() []testutils.TestCasesStruct {
	var validSlice int32 = 12
	app := "fundify"
	class := "custonboarding"
	var slice int32 = -1
	schemaNewTestcase := []testutils.TestCasesStruct{
		// 1st test case
		{
			Name: "ERROR- slice validation",
			RequestPayload: wscutils.Request{
				Data: breschema.BRESchemaGetReq{
					Slice: slice,
				},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "./data/scema_delete_error.json",
		},

		// 2nd test case
		{
			Name: "SUCCESS- delete schema by valid req ",
			RequestPayload: wscutils.Request{
				Data: breschema.BRESchemaListStruct{
					Slice: &validSlice,
					App:   &app,
					Class: &class,
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
	return schemaNewTestcase
}
