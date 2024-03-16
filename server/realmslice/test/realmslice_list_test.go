package realmslice_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
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
				require.JSONEq(t, string(jsonData), res.Body.String())
			} else {
				jsonData, err := testutils.ReadJsonFromFile(tc.TestJsonFile)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonData), res.Body.String())
			}
		})
	}

}

func RealmSliceListTestcase() []testutils.TestCasesStruct {
	time1, _ := time.Parse("2006-01-02T15:04:05Z", "2021-12-01T14:30:15Z")
	realmSliceNeTestcase := []testutils.TestCasesStruct{
		// test 1
		{
			Name: TestRealmSliceActivate_1,
			RequestPayload: wscutils.Request{
				Data: realmslice.RealmSliceActivateReq{},
			},

			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status: wscutils.SuccessStatus,
				Data: map[string][]sqlc.GetRealmSliceListByRealmRow{"slices": []sqlc.GetRealmSliceListByRealmRow{
					sqlc.GetRealmSliceListByRealmRow{
						ID:           12,
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
	return realmSliceNeTestcase
}
