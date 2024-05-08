package auth_test

// import (
// 	"bytes"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/remiges-tech/alya/wscutils"
// 	"github.com/remiges-tech/crux/server/capability"
// 	"github.com/remiges-tech/crux/testutils"
// 	"github.com/stretchr/testify/require"
// )

// func TestCapRevoke(t *testing.T) {
// 	testCases := capRevokeTestcase()
// 	for _, tc := range testCases {
// 		t.Run(tc.Name, func(t *testing.T) {

// 			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))

// 			res := httptest.NewRecorder()
// 			req, err := http.NewRequest(http.MethodPost, "/caprevoke", payload)
// 			require.NoError(t, err)

// 			r.ServeHTTP(res, req)

// 			require.Equal(t, tc.ExpectedHttpCode, res.Code)
// 			if tc.ExpectedResult != nil {
// 				jsonData := testutils.MarshalJson(tc.ExpectedResult)
// 				require.JSONEq(t, string(jsonData), res.Body.String())
// 			} else {
// 				jsonData, err := testutils.ReadJsonFromFile(tc.TestJsonFile)
// 				require.NoError(t, err)
// 				require.JSONEq(t, string(jsonData), res.Body.String())
// 			}
// 		})
// 	}

// }

// func capRevokeTestcase() []testutils.TestCasesStruct {

// 	return []testutils.TestCasesStruct{
// 		{
// 			Name:             "ERROR: Invalid User",
// 			RequestPayload:   wscutils.Request{Data: capability.CapRevokeReq{User: "kan@gmail.com", App: []string{"nedbank", "retailbank", "retailbank1"}, Cap: []string{"schema", "root"}}},
// 			ExpectedHttpCode: http.StatusBadRequest,
// 			ExpectedResult: wscutils.NewErrorResponse(),
// 		},
// 		{
// 			Name: " SUCCESS : Granting Only Realm Level Capabilities",
// 			RequestPayload: wscutils.Request{
// 				Data: capability.CapGrantRequest{
// 					User: "kanchan@gmail.com",
// 					App:  &[]string{"nedbank"},
// 					Cap:  []string{"root", "config"},
// 					From: fromTSPtr,
// 					To:   toTSPtr,
// 				},
// 			},
// 			ExpectedHttpCode: http.StatusOK,
// 		},
// 		{
// 			Name: " SUCCESS : Granting Only App Level Capabilities",
// 			RequestPayload: wscutils.Request{
// 				Data: capability.CapGrantRequest{
// 					User: "kanchan@gmail.com",
// 					App:  &[]string{"nedbank", "retailbank", "retailbank1"},
// 					Cap:  []string{"schema", "rules"},
// 					From: fromTSPtr,
// 					To:   toTSPtr,
// 				},
// 			},
// 			ExpectedHttpCode: http.StatusOK,
// 		},
// 		{
// 			Name: " SUCCESS : Granting both App and Realm Level Capabilities",
// 			RequestPayload: wscutils.Request{
// 				Data: capability.CapGrantRequest{
// 					User: "kanchan@gmail.com",
// 					App:  &[]string{"nedbank", "retailbank", "retailbank1"},
// 					Cap:  []string{"schema", "root"},
// 					From: fromTSPtr,
// 					To:   toTSPtr,
// 				},
// 			},
// 			ExpectedHttpCode: http.StatusOK,
// 		},

// 		{
// 			Name: " ERROR : Invalid timestamp",
// 			RequestPayload: wscutils.Request{
// 				Data: capability.CapGrantRequest{
// 					User: "kanchan@gmail.com",
// 					App:  &[]string{"nedbank", "retailbank", "retailbank1"},
// 					Cap:  []string{"schema", "root"},
// 					From: invalidfromTSPtr,
// 					To:   toTSPtr,
// 				},
// 			},
// 			ExpectedHttpCode: http.StatusBadRequest,
// 			TestJsonFile:     "../test/data/cap_grant_invalid_timestamp.json",
// 		},
// 	}
// 	return capGrantTestcase
// }
