package realmslice_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestRealmSlicePurge(t *testing.T) {
	testCases := RealmSlicePurgeTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/realmslicepurge", payload)
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

func RealmSlicePurgeTestcase() []testutils.TestCasesStruct {
	realmSliceNeTestcase := []testutils.TestCasesStruct{
		{
			Name:             "success all realmSlice purged",
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     nil,
				Messages: nil,
			},
		},
		{
			Name:             "error no realmSlice for  purged",
			ExpectedHttpCode: http.StatusBadRequest,
			ExpectedResult:   wscutils.NewErrorResponse(server.MsgId_NotFound, server.ErrCode_No_record_For_Purge),
		},
	}
	return realmSliceNeTestcase
}
