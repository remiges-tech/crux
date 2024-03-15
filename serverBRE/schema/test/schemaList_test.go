package test

import (
	"bytes"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/serverBRE/schema"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSchemaList(t *testing.T) {
	testCases := schemaListTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/BREschemaList", payload)
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

func schemaListTestcase() []testutils.TestCasesStruct {
	slice := int32(-1)
	app := "retailBANK1"
	class := "Inventoryitemsded"
	slice1 := int32(1)
	app1 := "retailBANK"
	class1 := "custonboarding"
	schemaListTestcase := []testutils.TestCasesStruct{
		{
			Name: "err- binding_json_error",
			RequestPayload: wscutils.Request{
				Data: nil,
			},

			ExpectedHttpCode: http.StatusBadRequest,
			ExpectedResult: &wscutils.Response{
				Status: wscutils.ErrorStatus,
				Data:   nil,
				Messages: []wscutils.ErrorMessage{
					{
						MsgID:   1001,
						ErrCode: wscutils.ErrcodeInvalidJson,
					},
				},
			},
		},
		{
			Name: "err- validation",
			RequestPayload: wscutils.Request{
				Data: schema.SchemaListStruct{
					Slice: &slice,
					App:   &app,
					Class: &class,
				},
			},

			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "./data/schema_list_validation_error.json",
		},
		{
			Name: "Success- get schema by app slice class",
			RequestPayload: wscutils.Request{
				Data: schema.SchemaListStruct{
					Slice: &slice1,
					App:   &app1,
					Class: &class1,
				},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_list_by_app_class_slice.json",
		},
		{
			Name: "Success- get schema by app",
			RequestPayload: wscutils.Request{
				Data: schema.SchemaListStruct{
					App: &app1,
				},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_list_by_app.json",
		},
		{
			Name: "Success- get schema by slice",
			RequestPayload: wscutils.Request{
				Data: schema.SchemaListStruct{
					Slice: &slice1,
				},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_list_by_slice.json",
		},
		{
			Name: "Success- get schema by class",
			RequestPayload: wscutils.Request{
				Data: schema.SchemaListStruct{
					Class: &class1,
				},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_list_by_class.json",
		},
		{
			Name: "Success- get schema by app slice",
			RequestPayload: wscutils.Request{
				Data: schema.SchemaListStruct{
					Slice: &slice1,
					App:   &app1,
				},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_list_by_app_slice.json",
		},
		{
			Name: "Success- get schema by slice class",
			RequestPayload: wscutils.Request{
				Data: schema.SchemaListStruct{
					Slice: &slice1,
					Class: &class1,
				},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_list_by_class_slice.json",
		},
		{
			Name: "Success- get schema by app class",
			RequestPayload: wscutils.Request{
				Data: schema.SchemaListStruct{
					App:   &app1,
					Class: &class1,
				},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_list_by_app_class.json",
		},
		{
			Name: "Success- get schema list",
			RequestPayload: wscutils.Request{
				Data: schema.SchemaListStruct{},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_list.json",
		},
	}
	return schemaListTestcase
}
