package wfinstance_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/server/wfinstance"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestGetWFInstanceAbort(t *testing.T) {
	testCases := wfInstanceAbortTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(server.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/wfinstanceabort", payload)
			require.NoError(t, err)

			r.ServeHTTP(res, req)

			require.Equal(t, tc.ExpectedHttpCode, res.Code)
			if tc.ExpectedResult != nil {
				jsonData := server.MarshalJson(tc.ExpectedResult)
				expectedJSON := string(jsonData)
				actualJSON := res.Body.String()
				require.JSONEq(t, expectedJSON, actualJSON)
			} else {
				jsonData, err := server.ReadJsonFromFile(tc.TestJsonFile)
				require.NoError(t, err)
				expectedJSON := string(jsonData)
				actualJSON := res.Body.String()
				require.JSONEq(t, expectedJSON, actualJSON)

			}
		})
	}

}

func wfInstanceAbortTestcase() []testutils.TestCasesStruct {
	var (
		ID              int32  = 777777
		entityID        string = "0eb8da50-aece-11ee-b168-3b192f7cd2b6"
		invalidID       int32  = 233
		invalidEntityID string = "0eb8da50-aece-11ee-b168"
	)
	wfInstanceAbortTestcase := []testutils.TestCasesStruct{

		// 1st test case
		{
			Name: "ERROR- Invalid request with two parameters",
			RequestPayload: wscutils.Request{
				Data: wfinstance.WFInstanceAbortRquest{
					ID:       &ID,
					EntityID: &entityID,
				},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "./data/wfinstanceabort_invalid_request_res.json",
		},
		// 2nd test case
		{
			Name: "ERROR- ID does not exist",
			RequestPayload: wscutils.Request{
				Data: wfinstance.WFInstanceAbortRquest{
					ID: &invalidID,
				},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "./data/wfinstaceabort_record_not_found_res.json",
		},

		// 3rd test case
		{
			Name: "ERROR- EntityID does not exist",
			RequestPayload: wscutils.Request{
				Data: wfinstance.WFInstanceAbortRquest{
					EntityID: &invalidEntityID,
				},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "./data/wfinstaceabort_record_not_found_res.json",
		},
		{
			Name: "SUCCESS- valid request ID",
			RequestPayload: wscutils.Request{
				Data: wfinstance.WFInstanceAbortRquest{
					ID: &ID,
				},
			},
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/wfinstanceabort_success_res.json",
		},
	}
	return wfInstanceAbortTestcase
}
