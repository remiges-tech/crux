package markdone_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server/markdone"
	"github.com/remiges-tech/crux/server/wfinstance"
	"github.com/remiges-tech/crux/testutils"
	"github.com/stretchr/testify/require"
)

func TestWFInstanceMarkDone(t *testing.T) {
	testCases := WFInstanceMarkDoneTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))
			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/wFinstancemarkdone", payload)
			require.NoError(t, err)
			r.ServeHTTP(res, req)
			require.Equal(t, tc.ExpectedHttpCode, res.Code)
			if tc.ExpectedResult != nil {
				jsonData := testutils.MarshalJson(tc.ExpectedResult)
				if strings.HasPrefix(tc.Name, "err") {
					require.JSONEq(t, string(jsonData), res.Body.String())
				} else {
					compareJSON(t, jsonData, res.Body.Bytes())
				}
			} else {
				jsonData, err := testutils.ReadJsonFromFile(tc.TestJsonFile)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonData), res.Body.String())
			}
		})
	}

}

func WFInstanceMarkDoneTestcase() []testutils.TestCasesStruct {
	fieldID := "ID"
	fieldStep := "Step"
	fieldEntity := "Entity"
	realmSliceNeTestcase := []testutils.TestCasesStruct{
		{
			Name: "mark done getcustdetails",
			RequestPayload: wscutils.Request{Data: markdone.WFInstanceMarkDoneReq{
				ID:         1,
				Entity:     map[string]string{"mode": "demat"},
				Step:       "getcustdetails",
				Stepfailed: false,
			}},
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{Status: wscutils.SuccessStatus, Data: wfinstance.WFInstanceNewResponse{
				Tasks: []map[string]int32{
					{
						"aof": 2,
					},
					{
						"dpandbankaccvalid": 3,
					},
					{
						"kycvalid": 4,
					},
					{
						"nomauth": 5,
					}},
				Nextstep: "auth_done",
				Loggedat: pgtype.Timestamp{Time: time.Now(), Valid: true},
				Subflows: map[string]string{"aof": "aofworkflow"},
			}, Messages: nil},
		},
		{
			Name: "mark done aof",
			RequestPayload: wscutils.Request{Data: markdone.WFInstanceMarkDoneReq{
				ID:         2,
				Entity:     map[string]string{"mode": "demat"},
				Step:       "aof",
				Stepfailed: false,
			}},
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{Status: wscutils.SuccessStatus, Data: wfinstance.WFInstanceNewResponse{
				ID:       "2",
				Loggedat: pgtype.Timestamp{Time: time.Now(), Valid: true},
			}, Messages: nil},
		},
		{
			Name: "mark done dpandbankaccvalid",
			RequestPayload: wscutils.Request{Data: markdone.WFInstanceMarkDoneReq{
				ID:         3,
				Entity:     map[string]string{"mode": "demat"},
				Step:       "dpandbankaccvalid",
				Stepfailed: false,
			}},
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{Status: wscutils.SuccessStatus, Data: wfinstance.WFInstanceNewResponse{
				ID:       "3",
				Loggedat: pgtype.Timestamp{Time: time.Now(), Valid: true},
			}, Messages: nil},
		},
		{
			Name: "mark done kycvalid",
			RequestPayload: wscutils.Request{Data: markdone.WFInstanceMarkDoneReq{
				ID:         4,
				Entity:     map[string]string{"mode": "demat"},
				Step:       "kycvalid",
				Stepfailed: false,
			}},
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{Status: wscutils.SuccessStatus, Data: wfinstance.WFInstanceNewResponse{
				ID:       "4",
				Loggedat: pgtype.Timestamp{Time: time.Now(), Valid: true},
			}, Messages: nil},
		},
		{
			Name: "mark done nomauth",
			RequestPayload: wscutils.Request{Data: markdone.WFInstanceMarkDoneReq{
				ID:         5,
				Entity:     map[string]string{"mode": "demat"},
				Step:       "nomauth",
				Stepfailed: false,
			}},
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{Status: wscutils.SuccessStatus, Data: wfinstance.WFInstanceNewResponse{
				Tasks: []map[string]int32{{
					"sendauthlinktoclient": 6,
				}},
				Loggedat: pgtype.Timestamp{Time: time.Now(), Valid: true},
				Subflows: map[string]string{},
			}, Messages: nil},
		},
		{
			Name: "Stepfailed sendauthlinktoclient",
			RequestPayload: wscutils.Request{Data: markdone.WFInstanceMarkDoneReq{
				ID:         6,
				Entity:     map[string]string{"mode": "demat"},
				Step:       "sendauthlinktoclient",
				Stepfailed: true,
			}},
			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{Status: wscutils.SuccessStatus, Data: wfinstance.WFInstanceNewResponse{
				Loggedat: pgtype.Timestamp{Time: time.Now(), Valid: true},
				Done:     "true",
			}, Messages: nil},
		},
		{
			Name: "error reattempting markdone for last step of workflow",
			RequestPayload: wscutils.Request{Data: markdone.WFInstanceMarkDoneReq{
				ID:         6,
				Entity:     map[string]string{"mode": "demat"},
				Step:       "sendauthlinktoclient",
				Stepfailed: false,
			}},
			ExpectedHttpCode: http.StatusBadRequest,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.ErrorStatus,
				Data:     "error while GetWFInstanceFromId() in WFInstanceMarkDone",
				Messages: nil,
			},
		},
		{
			Name: "err standard validation",
			RequestPayload: wscutils.Request{Data: markdone.WFInstanceMarkDoneReq{
				ID:         0,
				Entity:     nil,
				Step:       "",
				Stepfailed: false,
			}},
			ExpectedHttpCode: http.StatusBadRequest,
			ExpectedResult: &wscutils.Response{Status: wscutils.ErrorStatus, Data: nil, Messages: []wscutils.ErrorMessage{
				{
					MsgID:   101,
					ErrCode: "required",
					Field:   &fieldID,
				},
				{
					MsgID:   101,
					ErrCode: "required",
					Field:   &fieldEntity,
				},
				{
					MsgID:   101,
					ErrCode: "required",
					Field:   &fieldStep,
				},
			}},
		},
	}
	return realmSliceNeTestcase
}

func compareJSON(t *testing.T, expected, actual []byte) {
	var expectedMap, actualMap map[string]interface{}

	if err := json.Unmarshal(expected, &expectedMap); err != nil {
		t.Fatalf("Error unmarshaling expected JSON: %v", err)
	}
	if err := json.Unmarshal(actual, &actualMap); err != nil {
		t.Fatalf("Error unmarshaling actual JSON: %v", err)
	}

	// Remove the "Loggedat" field
	delete(expectedMap["data"].(map[string]interface{}), "loggedat")

	delete(actualMap["data"].(map[string]interface{}), "loggedat")

	if !reflect.DeepEqual(expectedMap, actualMap) {
		t.Errorf("Expected JSON does not match actual JSON")
	}
}
