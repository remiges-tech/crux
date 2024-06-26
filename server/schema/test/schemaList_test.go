package schema_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/remiges-tech/alya/wscutils"

	"github.com/remiges-tech/crux/server/schema"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestSchemaList(t *testing.T) {
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
			req, err := http.NewRequest(http.MethodPost, "/wfschemalist", payload)
			require.NoError(t, err)

			r.ServeHTTP(res, req)

			require.Equal(t, tc.ExpectedHttpCode, res.Code)
			if tc.ExpectedResult != nil {
				jsonData := testutils.MarshalJson(tc.ExpectedResult)
				compareJSON(t, jsonData, res.Body.Bytes())
			} else {
				jsonData, err := testutils.ReadJsonFromFile(tc.TestJsonFile)
				require.NoError(t, err)
				compareJSON(t, jsonData, res.Body.Bytes())
			}
		})
	}
}

func schemaListTestcase() []testutils.TestCasesStruct {
	// slice := int32(-1)
	// app := "retailBANK1"
	// class := "Inventoryitemsded"
	slice1 := int32(13)
	app1 := "retailbank"
	class1 := "members"
	schemaListTestcase := []testutils.TestCasesStruct{
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
		{
			Name: "root cap",
			RequestPayload: wscutils.Request{
				Data: schema.SchemaListStruct{},
			},

			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/schema_list.json",
		},
		{
			Name: "schema cap",
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
	}
	return schemaListTestcase
}
func removeField(data map[string]interface{}, field string) {
	for _, v := range data {
		switch vv := v.(type) {
		case map[string]interface{}:
			removeField(vv, field)
		case []interface{}:
			for _, u := range vv {
				if umap, ok := u.(map[string]interface{}); ok {
					removeField(umap, field)
				}
			}
		}
	}
	delete(data, field)
}

func compareJSON(t *testing.T, expected, actual []byte) {
	var expectedMap, actualMap map[string]interface{}

	if err := json.Unmarshal(expected, &expectedMap); err != nil {
		t.Fatalf("Error unmarshaling expected JSON: %v", err)
	}
	if err := json.Unmarshal(actual, &actualMap); err != nil {
		t.Fatalf("Error unmarshaling actual JSON: %v", err)
	}

	fieldsToRemove := []string{"createdat", "editedat", "deactivateat"}

	for _, field := range fieldsToRemove {
		removeField(expectedMap, field)
		removeField(actualMap, field)
	}

	if !reflect.DeepEqual(expectedMap, actualMap) {
		t.Errorf("Expected JSON does not match actual JSON")
	}
}
