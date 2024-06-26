package realmslice_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestRealmSliceApps(t *testing.T) {
	testCases := RealmSliceAppsTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, tc.Url, nil)
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

func RealmSliceAppsTestcase() []testutils.TestCasesStruct {
	realmSliceNeTestcase := []testutils.TestCasesStruct{
		{
			Name:             "Success- app list by id exist",
			Url:              "/realmsliceapps/11",
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status: wscutils.SuccessStatus,
				Data: []sqlc.RealmSliceAppsListRow{
					{
						Shortname: "Amazon",
						Longname:  "American multinational technology company, engaged in e-commerce",
					},
					{

						Shortname: "Myntra",
						Longname:  "American multinational technology company, engaged in e-commerce",
					},
				},
				Messages: nil,
			},
		},
		{
			Name:             "Success- app list by id not exist ",
			Url:              "/realmsliceapps/111",
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     []sqlc.RealmSliceAppsListRow{},
				Messages: nil,
			},
		},
	}
	return realmSliceNeTestcase
}
