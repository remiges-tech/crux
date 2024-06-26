package realmslice_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server/realmslice"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

const (
	TestRealmSliceList_1 = "success: TestRealmSliceList"
)

var resp wscutils.Response

func TestRealmSliceList(t *testing.T) {
	testCases := RealmSliceListTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))
			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/realmslicelist", payload)
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

func RealmSliceListTestcase() []testutils.TestCasesStruct {
	time1, _ := time.Parse("2006-01-02T15:04:05Z", "2021-12-01T14:30:15Z")
	realmSliceTestCase := []testutils.TestCasesStruct{
		{
			Name: TestRealmSliceList_1,
			RequestPayload: wscutils.Request{
				Data: realmslice.RealmSliceActivateReq{},
			},
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status: wscutils.SuccessStatus,
				Data: map[string][]sqlc.GetRealmSliceListByRealmRow{"slices": []sqlc.GetRealmSliceListByRealmRow{
					{
						ID:           12,
						Descr:        "Stock Market",
						Active:       true,
						Deactivateat: pgtype.Timestamp{},
						Createdat:    pgtype.Timestamp{Time: time1, Valid: true},
						Createdby:    "aniket",
						Editedat:     pgtype.Timestamp{},
						Editedby:     pgtype.Text{},
					},
					{
						ID:           13,
						Descr:        "Stock Market",
						Active:       true,
						Deactivateat: pgtype.Timestamp{},
						Createdat:    pgtype.Timestamp{Time: time1, Valid: true},
						Createdby:    "aniket",
						Editedat:     pgtype.Timestamp{},
						Editedby:     pgtype.Text{},
					},
					{
						ID:           1,
						Descr:        "Stock Market",
						Active:       true,
						Deactivateat: pgtype.Timestamp{},
						Createdat:    pgtype.Timestamp{Time: time1, Valid: true},
						Createdby:    "aniket",
						Editedat:     pgtype.Timestamp{},
						Editedby:     pgtype.Text{},
					},
					{
						ID:           2,
						Descr:        "Stock Market",
						Active:       true,
						Deactivateat: pgtype.Timestamp{},
						Createdat:    pgtype.Timestamp{Time: time1, Valid: true},
						Createdby:    "aniket",
						Editedat:     pgtype.Timestamp{},
						Editedby:     pgtype.Text{},
					},
					{
						ID:           3,
						Descr:        "Stock Market",
						Active:       true,
						Deactivateat: pgtype.Timestamp{},
						Createdat:    pgtype.Timestamp{Time: time1, Valid: true},
						Createdby:    "aniket",
						Editedat:     pgtype.Timestamp{},
						Editedby:     pgtype.Text{},
					},
				}},
				Messages: nil,
			},
		},
	}
	return realmSliceTestCase
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
