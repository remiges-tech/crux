package test

import (
	"bytes"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/serverBRE/ruleset"
	"github.com/remiges-tech/crux/types"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRuleSetGet(t *testing.T) {
	testCases := ruleSetGetTestCase()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(types.MarshalJson(tc.requestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/ruleSetGet", payload)
			require.NoError(t, err)

			r.ServeHTTP(res, req)

			require.Equal(t, tc.expectedHttpCode, res.Code)
			if tc.expectedResult != nil {
				jsonData := types.MarshalJson(tc.expectedResult)
				require.JSONEq(t, string(jsonData), res.Body.String())
			} else {
				jsonData, err := types.ReadJsonFromFile(tc.testJsonFile)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonData), res.Body.String())
			}
		})
	}

}

func ruleSetGetTestCase() []TestCasesStruct {
	var sliceStr int32 = 2
	app := "retailbank"
	class := "members"
	tname := "goldstatus"
	var slice int32 = -1
	schemaNewTestcase := []TestCasesStruct{
		{
			name: "ERROR- slice validation",
			requestPayload: wscutils.Request{
				Data: ruleset.RuleSetGetReq{
					Slice: slice,
				},
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedResult: &wscutils.Response{
				Status: wscutils.ErrorStatus,
				Data:   nil,
				Messages: []wscutils.ErrorMessage{
					{
						MsgID:   102,
						ErrCode: "greater",
						Field:   &types.SLICE,
					}, {
						MsgID:   101,
						ErrCode: "required",
						Field:   &types.APP,
					}, {
						MsgID:   101,
						ErrCode: "required",
						Field:   &types.CLASS,
					},
					{
						MsgID:   101,
						ErrCode: "required",
						Field:   &types.NAME,
					},
				},
			},
		},
		{
			name: "SUCCESS- get ruleset by valid req ",
			requestPayload: wscutils.Request{
				Data: ruleset.RuleSetGetReq{
					Slice: sliceStr,
					App:   app,
					Class: class,
					Name:  tname,
				},
			},

			expectedHttpCode: http.StatusOK,
			testJsonFile:     "./data/ruleset_get_response.json",
		},
		{
			name: "Failed- get ruleset by invalid req ",
			requestPayload: wscutils.Request{
				Data: ruleset.RuleSetGetReq{
					Slice: sliceStr,
					App:   "xyz",
					Class: class,
					Name:  tname,
				},
			},

			expectedHttpCode: http.StatusBadRequest,
			testJsonFile:     "./data/ruleset_get_failed_response.json",
		},
	}
	return schemaNewTestcase
}
