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

// This TesFile tests WFInstanceList Api for this you must need to run WfinstanceNew API
func TestGetWFInstanceList(t *testing.T) {

	testCases := wfInstanceListTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(server.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/wfinstancelist", payload)
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

func wfInstanceListTestcase() []testutils.TestCasesStruct {
	var (
		InvalidSlice int32  = 1
		Slice        int32  = 2
		entityID     string = "tempentityid"
		App          string = "retailBANK"
		Workflow     string = "temp"
		Parent       int32  = 78
	)

	wfInstanceListTestcase := []testutils.TestCasesStruct{

		// 1st test case
		{
			Name: "ERROR- Invalid request ",
			RequestPayload: wscutils.Request{
				Data: wfinstance.WFInstanceListRequest{
					Slice:    &InvalidSlice,
					EntityID: &entityID,
				},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "./data/wfinstancelist_invalid_request_res.json",
		},
		// 2nd test case
		{
			Name: "SUCCESS- valid request",
			RequestPayload: wscutils.Request{
				Data: wfinstance.WFInstanceListRequest{
					Slice:    &Slice,
					EntityID: &entityID,
					App:      &App,
					Workflow: &Workflow,
					Parent:   &Parent,
				},
			},
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/wfinstancelist_valid_request_res.json",
		},
	}
	return wfInstanceListTestcase
}
