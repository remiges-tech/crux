package schema_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server/schema"
	"github.com/remiges-tech/crux/testutils"
	"github.com/remiges-tech/crux/types"
	"github.com/stretchr/testify/require"
)

func TestSchemaDelete(t *testing.T) {
	testCases := schemaDeleteTestCase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(types.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/wfschemadelete", payload)
			require.NoError(t, err)

			r.ServeHTTP(res, req)

			require.Equal(t, tc.ExpectedHttpCode, res.Code)
			if tc.ExpectedResult != nil {
				jsonData := types.MarshalJson(tc.ExpectedResult)
				require.JSONEq(t, string(jsonData), res.Body.String())
			} else {
				jsonData, err := types.ReadJsonFromFile(tc.TestJsonFile)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonData), res.Body.String())
			}
		})
	}

}

func schemaDeleteTestCase() []testutils.TestCasesStruct {
	var sliceStr int32 = 1
	app := "retailbank"
	class := "tempclass"
	var slice int32 = -1
	schemaNewTestcase := []testutils.TestCasesStruct{
		// 1st test case
		{
			Name: "ERROR- slice validation",
			RequestPayload: wscutils.Request{
				Data: schema.SchemaGetReq{
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
			Name: "SUCCESS- delete schema by valid req ",
			RequestPayload: wscutils.Request{
				Data: schema.SchemaListStruct{
					Slice: &sliceStr,
					App:   &app,
					Class: &class,
				},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_delete_response.json",
		},
	}
	return schemaNewTestcase
}
