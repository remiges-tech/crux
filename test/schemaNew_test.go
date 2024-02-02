package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/test/testutils"
	"github.com/stretchr/testify/require"
)

type TestCasesStruct struct {
	name             string
	requestPayload   wscutils.Request
	expectedHttpCode int
	testJsonFile     string
	expectedResult   *wscutils.Response
}

func TestSchemaNew(t *testing.T) {
	testCases := schemaNewTestcase()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(testutils.MarshalJson(tc.requestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/WFschemaNew", payload)
			require.NoError(t, err)

			r.ServeHTTP(res, req)

			require.Equal(t, tc.expectedHttpCode, res.Code)
			if tc.expectedResult != nil {
				jsonData := testutils.MarshalJson(tc.expectedResult)
				require.JSONEq(t, string(jsonData), res.Body.String())
			} else {
				jsonData, err := testutils.ReadJsonFromFile(tc.testJsonFile)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonData), res.Body.String())
			}
		})
	}

}

func schemaNewTestcase() []TestCasesStruct {
	schemaNewTestcase := []TestCasesStruct{
		{
			name: "err- binding_json_error",
			requestPayload: wscutils.Request{
				Data: nil,
			},

			expectedHttpCode: http.StatusBadRequest,
			expectedResult: &wscutils.Response{
				Status: wscutils.ErrorStatus,
				Data:   nil,
				Messages: []wscutils.ErrorMessage{
					{
						MsgID:   0,
						ErrCode: wscutils.ErrcodeInvalidJson,
					},
				},
			},
		},
	}
	return schemaNewTestcase
}
