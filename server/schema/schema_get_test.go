package schema

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/types"
	"github.com/stretchr/testify/require"
)

func TestSchemaGet(t *testing.T) {
	testCases := schemaGetTestcase()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(types.MarshalJson(tc.requestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/wfschemaget", payload)
			require.NoError(t, err)

			r.ServeHTTP(res, req)

			require.Equal(t, tc.expectedHttpCode, res.Code)
			if tc.expectedResult != nil {
				jsonData := types.MarshalJson(tc.expectedResult)
				require.JSONEq(t, string(jsonData), res.Body.String())
			} else {
				jsonData, err := types.ReadJsonFromFile(tc.testJsonFile)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonData), res.Body.String())
			}
		})
	}

}

func schemaGetTestcase() []TestCasesStruct {
	var sliceStr int32 = 2
	app := "nedbank"
	class := "custonboarding"
	slice := int32(-1)
	// slice1 := int32(1)
	schemaNewTestcase := []TestCasesStruct{
		// 1st test case
		{
			name: "err- slice validation",
			requestPayload: wscutils.Request{
				Data: schemaGetReq{
					Slice: &slice,
				},
			},

			expectedHttpCode: http.StatusBadRequest,
			expectedResult: &wscutils.Response{
				Status: wscutils.ErrorStatus,
				Data:   nil,
				Messages: []wscutils.ErrorMessage{

					{
						MsgID:   101,
						ErrCode: "gt",
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
			name: "success- get schema by valid req ",
			requestPayload: wscutils.Request{
				Data: SchemaListStruct{
					Slice: &sliceStr,
					App:   &app,
					Class: &class,
				},
			},

			expectedHttpCode: http.StatusOK,
			testJsonFile:     "./testData/schema_get_response.json",
		},
	}
	return schemaNewTestcase
}
