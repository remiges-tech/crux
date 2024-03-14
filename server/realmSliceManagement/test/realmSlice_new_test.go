package realmSliceManagement_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server/realmSliceManagement"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestRealmSliceNew(t *testing.T) {
	testCases := RealmSliceNewTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/realmSliceNew", payload)
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

func RealmSliceNewTestcase() []testutils.TestCasesStruct {
	feild := "CopyOf"
	realmSliceNeTestcase := []testutils.TestCasesStruct{
		{
			Name: "err- standard validation",
			RequestPayload: wscutils.Request{
				Data: realmSliceManagement.RealmSliceNewRequest{
					CopyOf: -9,
				},
			},

			ExpectedHttpCode: http.StatusBadRequest,
			ExpectedResult:   wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{{MsgID: 102, ErrCode: "greater", Field: &feild}}),
		},
		{
			Name: "err- custom validation realmSlice not exist",
			RequestPayload: wscutils.Request{
				Data: realmSliceManagement.RealmSliceNewRequest{
					CopyOf: 99,
					App:    []string{},
				},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			ExpectedResult:   wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{{MsgID: 1006, ErrCode: "not_exist", Field: &feild}}),
		},
		{
			Name: "Success- create new copy of realmSlice by old realmSlice id",
			RequestPayload: wscutils.Request{
				Data: realmSliceManagement.RealmSliceNewRequest{
					CopyOf: 11,
				},
			},

			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     1,
				Messages: nil,
			},
		},
		{
			Name: "Success- create new copy of realmSlice by old realmSlice id with only the listed apps",
			RequestPayload: wscutils.Request{
				Data: realmSliceManagement.RealmSliceNewRequest{
					CopyOf: 11,
				},
			},

			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     2,
				Messages: nil,
			},
		},
		{
			Name: "Success- create new of realmSlice with description",
			RequestPayload: wscutils.Request{
				Data: realmSliceManagement.RealmSliceNewRequest{
					CopyOf: 11,
					Descr:  "description for new app",
				},
			},

			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     3,
				Messages: nil,
			},
		},
	}
	return realmSliceNeTestcase
}
