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

func TestWFinstanceTry(t *testing.T) {
	testCases := wfInstanceTryTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			// Setting up buffer
			payload := bytes.NewBuffer(server.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/wfinstancetry", payload)
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
				expectedJSON := removeFieldFromJSON(string(jsonData), "loggedat")
				actualJSON := removeFieldFromJSON(res.Body.String(), "loggedat")
				require.JSONEq(t, expectedJSON, actualJSON)

			}
		})
	}

}
func wfInstanceTryTestcase() []testutils.TestCasesStruct {
	entityID := "0eb8da50-aece-11ee-b168-3b192f7cd2b6"
	trace := 0
	stepfailed := false
	app := "uccapp"
	workflow := "aofworkflow"
	InvalidStep := "sendtortaaa"
	wfInstanceNewTestcase := []testutils.TestCasesStruct{

		// 1st test case
		{
			Name: "ERROR- Standard error",
			RequestPayload: wscutils.Request{
				Data: wfinstance.WFInstanceTryRequest{},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "./data/wfinstancetry_invalid_response.json",
		},

		// 2nd test case
		{
			Name: "SUCCESS- valid response",
			RequestPayload: wscutils.Request{
				Data: wfinstance.WFInstanceTryRequest{
					Slice:    int32(13),
					App:      app,
					EntityID: entityID,
					Entity: map[string]string{
						"class":     "ucc_aof",
						"step":      "getsigneddocument",
						"aofexists": "false",
					},
					Workflow:   workflow,
					Step:       "sendtorta",
					StepFailed: &stepfailed,
					Trace:      &trace,
				},
			},
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/wfinstancetry_valid_response.json",
		},

		// 3rd test case
		{
			Name: "ERROR : Invalid step",
			RequestPayload: wscutils.Request{
				Data: wfinstance.WFInstanceTryRequest{
					Slice:    int32(13),
					App:      app,
					EntityID: entityID,
					Entity: map[string]string{
						"class":     "ucc_aof",
						"step":      "getsigneddocument",
						"aofexists": "false",
					},
					Workflow:   workflow,
					Step:       InvalidStep,
					StepFailed: &stepfailed,
					Trace:      &trace,
				},
			},
			ExpectedHttpCode: http.StatusBadRequest,
			TestJsonFile:     "./data/invalid_step_response.json",
		},
	}
	return wfInstanceNewTestcase
}
