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

func TestSchemaList(t *testing.T) {
	testCases := schemaListTestcase()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(types.MarshalJson(tc.requestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/WFschemaList", payload)
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

func schemaListTestcase() []TestCasesStruct {
	feild := "Slice"
	slice := int32(-1)
	slice1 := int32(1)
	schemaNewTestcase := []TestCasesStruct{
		{
			name: "err- slice validation",
			requestPayload: wscutils.Request{
				Data: SchemaListStruct{
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
						Field:   &feild,
					},
				},
			},
		},
		{
			name: "suc- get schema by slice ",
			requestPayload: wscutils.Request{
				Data: SchemaListStruct{
					Slice: &slice1,
				},
			},

			expectedHttpCode: http.StatusOK,
			testJsonFile:     "./testData/schemaListByslice.json",
		},
	}
	return schemaNewTestcase
}
