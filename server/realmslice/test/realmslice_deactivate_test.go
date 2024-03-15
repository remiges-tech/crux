package realmslice_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server/realmslice"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

const (
	TestRealmSliceDeactivate_1 = "success: realm_slice_deactivate"
	TestRealmSliceDeactivate_2 = "error: field validation"
)

func TestRealmSliceDeactivate(t *testing.T) {
	testCases := RealmSliceDectivateTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))
			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/realmslicedeactivate", payload)
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

func RealmSliceDectivateTestcase() []testutils.TestCasesStruct {
	id := "Id"
	realmSliceNeTestcase := []testutils.TestCasesStruct{
		// test 1
		{
			Name: TestRealmSliceActivate_1,
			RequestPayload: wscutils.Request{
				Data: realmslice.RealmSliceActivateReq{
					Id: 11,
					// From: ,
				},
			},

			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     nil,
				Messages: nil,
			},
		},
		// test 2
		{
			Name: TestRealmSliceActivate_2,
			RequestPayload: wscutils.Request{
				Data: realmslice.RealmSliceActivateReq{},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			ExpectedResult: &wscutils.Response{
				Status: wscutils.ErrorStatus,
				Data:   nil,
				Messages: []wscutils.ErrorMessage{
					wscutils.ErrorMessage{
						MsgID:   101,
						ErrCode: "required",
						Field:   &id,
					},
				},
			},
		},
	}
	return realmSliceNeTestcase
}
