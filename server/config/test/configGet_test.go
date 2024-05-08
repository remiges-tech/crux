package config_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestConfigGet(t *testing.T) {
	testCases := ConfigGetTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/configGet", nil)
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

func ConfigGetTestcase() []testutils.TestCasesStruct {
	realmSliceNeTestcase := []testutils.TestCasesStruct{
		{
			Name:             "Success- create get config list",
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status: wscutils.SuccessStatus,
				Data: []sqlc.ConfigGetRow{
					{
						Attr: "CONFIG_A",
						Val:  pgtype.Text{String: "Value for CONFIG_A", Valid: true},
						Ver:  pgtype.Int4{Int32: 1, Valid: true},
						By:   "User1",
					},
				},
				Messages: nil,
			},
		},
	}
	return realmSliceNeTestcase
}
