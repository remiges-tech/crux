package breschema_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"

	breschema "github.com/remiges-tech/crux/server/BRESchema"
	"github.com/remiges-tech/crux/server/schema"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestBRESchemaList(t *testing.T) {
	testCases := schemaListTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			if tc.Name == "root cap" {
				schema.CapForList = []string{"root"}
			} else if tc.Name == "schema cap" {
				schema.CapForList = []string{"schema"}
			}
			req, err := http.NewRequest(http.MethodPost, "/breschemalist", payload)
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
	slice2 := int32(11)
	app1 := "amazon"
	class1 := "inventoryitems"
	schemaListTestcase := []testutils.TestCasesStruct{

		{
			Name: "Success- get schema by app slice class",
			RequestPayload: wscutils.Request{
				Data: breschema.BRESchemaListStruct{
					Slice: &slice2,
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
				Data:  breschema.BRESchemaListStruct{
					App: &app1,
				},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_list_by_app_class_slice.json",
		},
		{
			Name: "Success- get schema by slice",
			RequestPayload: wscutils.Request{
				Data: breschema.BRESchemaListStruct{
					Slice: &slice2,
				},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_list_by_slice.json",
		},
		{
			Name: "Success- get schema by class",
			RequestPayload: wscutils.Request{
				Data:  breschema.BRESchemaListStruct{
					Class: &class1,
				},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_list_by_class.json",
		},
		{
			Name: "Success- get schema list",
			RequestPayload: wscutils.Request{
				Data:  breschema.BRESchemaListStruct{},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_list.json",
		},
		{
			Name: "root cap",
			RequestPayload: wscutils.Request{
				Data:  breschema.BRESchemaListStruct{},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_list.json",
		},
	}
	return schemaListTestcase
}
