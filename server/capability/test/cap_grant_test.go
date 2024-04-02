package auth_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server/capability"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

// {
// 	"data":{
// 	  "user": "kanchan@gmail.com",
// 	  "app": [ "nedbank" ],
// 	  "cap": [ "root", "config" ],
// 	  "from": "2024-05-20T10:05:00Z",
// 	  "to": "2024-05-31T10:05:00Z"
// 	}
//   }    only realm level

// {
// 	"data":{
// 	  "user": "kanchan@gmail.com",
// 	  "app": [ "nedbank","retailbank","retailbank1"],
// 	  "cap": [ "schema", "rules" ],
// 	  "from": "2024-05-20T10:05:00Z",
// 	  "to": "2024-05-31T10:05:00Z"
// 	}
//   }   only app level

// {
// 	"data":{
// 	  "user": "kanchan@gmail.com",
// 	  "app": [ "nedbank","retailbank","retailbank1"],
// 	  "cap": [ "schema", "root" ],
// 	  "from": "2024-05-20T10:05:00Z",
// 	  "to": "2024-05-31T10:05:00Z"
// 	}
//   }   mix

// {
// 	"data":{
// 	  "user": "kanchan@gmail.com",
// 	  "app": [ "nedbank","retailbank","retailbank1"],
// 	  "cap": [ "schema", "root" ],
// 	  "from": "2024-05-20T10:05:00Z",
// 	  "to": "2024-05-31T10:05:00Z"
// 	}
//   }   // invalid user

// {
// 	"data":{
// 	  "user": "kanchan@gmail.com",
// 	  "app": [ "nedbank","retailbank","retailbank1"],
// 	  "cap": [ "schema", "root" ],
// 	  "from": "2022-05-20T10:05:00Z",
// 	  "to": "2024-05-31T10:05:00Z"
// 	}
//   }   // invalid time

func TestCapGrant(t *testing.T) {
	testCases := capGrantTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/capgrant", payload)
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

func capGrantTestcase() []testutils.TestCasesStruct {
	fromTS := time.Date(2024, 5, 20, 10, 5, 0, 0, time.UTC)
	toTS := time.Date(2024, 5, 31, 10, 5, 0, 0, time.UTC)
	InvalidfromTS := time.Date(2022, 5, 20, 10, 5, 0, 0, time.UTC)

	fromTSPtr := &fromTS
	toTSPtr := &toTS
	invalidfromTSPtr := &InvalidfromTS

	capGrantTestcase := []testutils.TestCasesStruct{
		{
			Name:             "ERROR: Invalid User",
			RequestPayload:   wscutils.Request{Data: capability.CapGrantRequest{User: "kan@gmail.com", App: &[]string{"nedbank", "retailbank", "retailbank1"}, Cap: []string{"schema", "root"}, From: fromTSPtr, To: toTSPtr}},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "../test/data/cap_grant_invalid_user.json",
		},
		{
			Name: " SUCCESS : Granting Only Realm Level Capabilities",
			RequestPayload: wscutils.Request{
				Data: capability.CapGrantRequest{
					User: "kanchan@gmail.com",
					App:  &[]string{"nedbank"},
					Cap:  []string{"root", "config"},
					From: fromTSPtr,
					To:   toTSPtr,
				},
			},
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     nil,
				Messages: nil,
			},
		},
		{
			Name: " SUCCESS : Granting Only App Level Capabilities",
			RequestPayload: wscutils.Request{
				Data: capability.CapGrantRequest{
					User: "kanchan@gmail.com",
					App:  &[]string{"nedbank", "retailbank", "retailbank1"},
					Cap:  []string{"schema", "rules"},
					From: fromTSPtr,
					To:   toTSPtr,
				},
			},
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     nil,
				Messages: nil,
			},
		},
		{
			Name: " SUCCESS : Granting both App and Realm Level Capabilities",
			RequestPayload: wscutils.Request{
				Data: capability.CapGrantRequest{
					User: "kanchan@gmail.com",
					App:  &[]string{"nedbank", "retailbank", "retailbank1"},
					Cap:  []string{"schema", "root"},
					From: fromTSPtr,
					To:   toTSPtr,
				},
			},
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     nil,
				Messages: nil,
			},
		},

		{
			Name: " ERROR : Invalid timestamp",
			RequestPayload: wscutils.Request{
				Data: capability.CapGrantRequest{
					User: "kanchan@gmail.com",
					App:  &[]string{"nedbank", "retailbank", "retailbank1"},
					Cap:  []string{"schema", "root"},
					From: invalidfromTSPtr,
					To:   toTSPtr,
				},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "../test/data/cap_grant_invalid_timestamp.json",
		},
	}
	return capGrantTestcase
}
