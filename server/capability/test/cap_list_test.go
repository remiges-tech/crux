package auth_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server/capability"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

const (
	test_cap_list_1 = "ERROR_1: Invalid cap_list"
	test_cap_list_2 = "SUCCESS_2: app & cap both"
	test_cap_list_3 = "SUCCESS_3: all list"
	test_cap_list_4 = "SUCCESS_4: app only"
	test_cap_list_5 = "SUCCESS_5: cap only"
)

func TestCapList(t *testing.T) {
	testCases := testCapList()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/caplist", payload)
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

func testCapList() []testutils.TestCasesStruct {
	app := []string{"hdfcbank", "nedbank"}
	appOnly := []string{"amazon"}
	cap := []string{"root"}
	testUsrAct := []testutils.TestCasesStruct{
		// test 1 : bad req
		{
			Name:             test_cap_list_1,
			RequestPayload:   wscutils.Request{},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "../test/data/bad_req.json",
		},
		// test 2 : app & cap both
		{
			Name: test_cap_list_2,
			RequestPayload: wscutils.Request{
				Data: capability.CapListReq{
					App: app,
					Cap: cap,
				},
			},
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "../test/data/cap_list_both_resp.json",
		},
		// test 3 : all list
		{
			Name: test_cap_list_3,
			RequestPayload: wscutils.Request{
				Data: capability.CapListReq{},
			},
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "../test/data/cap_list_all_resp.json",
		},
		// test 4 : app only
		{
			Name: test_cap_list_4,
			RequestPayload: wscutils.Request{
				Data: capability.CapListReq{
					App: appOnly,
				},
			},
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "../test/data/cap_list_app_only_resp.json",
		},
		// test 5 : cap only
		{
			Name: test_cap_list_5,
			RequestPayload: wscutils.Request{
				Data: capability.CapListReq{
					Cap: cap,
				},
			},
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "../test/data/cap_list_cap_only_resp.json",
		},
	}
	return testUsrAct
}
