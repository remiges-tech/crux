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
	"github.com/remiges-tech/crux/types"
	"github.com/stretchr/testify/require"
)

const (
	TestSchemaGet_1 = "ERROR_1- slice validation"
	TestSchemaGet_2 = "SUCCESS_2- get schema by valid req"
)

func TestBRESchemaGet(t *testing.T) {
	testCases := schemaGetTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(server.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/breschemaget", payload)
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

func schemaGetTestcase() []testutils.TestCasesStruct {

	var slice int32 = -1
	schemaNewTestcase := []testutils.TestCasesStruct{
		// 1st test case
		{
			Name: TestSchemaGet_1,
			RequestPayload: wscutils.Request{
				Data: breschema.BRESchemaGetReq{
					Slice: slice,
				},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			ExpectedResult: &wscutils.Response{
				Status: wscutils.ErrorStatus,
				Data:   nil,
				Messages: []wscutils.ErrorMessage{
					{
						MsgID:   102,
						ErrCode: "greater",
						Field:   &types.SLICE,
					}, {
						MsgID:   101,
						ErrCode: "required",
						Field:   &types.APP,
					}, {
						MsgID:   101,
						ErrCode: "required",
						Field:   &types.CLASS,
					},
				},
			},
		},

		// 2nd test case
		{
			Name: TestSchemaGet_2,
			RequestPayload: wscutils.Request{
				Data: breschema.BRESchemaGetReq{
					Slice: 11,
					App:   "amazon",
					Class: "inventoryitems",
				},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_get_response.json",
		},
	}
	return schemaNewTestcase
}
