package wfinstance_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server/wfinstance"
	"github.com/remiges-tech/crux/testutils"
	"github.com/remiges-tech/crux/types"
	"github.com/stretchr/testify/require"
)

func TestWFinstanceNew(t *testing.T) {
	testCases := wfInstanceNewTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			switch tc.Name {
			case "SUCCESS- Single step":
				wfinstance.SWITCH = wfinstance.ActionSet{
					Tasks:      []string{"diwalisale"},
					Properties: map[string]string{"nextstep": "coupondistribution"},
				}
			case "Success- Multiple steps":
				wfinstance.SWITCH = wfinstance.ActionSet{
					Tasks:      []string{"diwalisale", "yearendsale"},
					Properties: map[string]string{"nextstep": "coupondistribution"},
				}
			default:
				wfinstance.SWITCH = wfinstance.ActionSet{
					Tasks:      []string{"diwalisale"},
					Properties: map[string]string{"nextstep": "coupondistribution"},
				}
			}

			// Setting up buffer
			payload := bytes.NewBuffer(types.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/wfinstancenew", payload)
			require.NoError(t, err)

			r.ServeHTTP(res, req)

			require.Equal(t, tc.ExpectedHttpCode, res.Code)
			if tc.ExpectedResult != nil {
				jsonData := types.MarshalJson(tc.ExpectedResult)
				expectedJSON := string(jsonData)
				actualJSON := res.Body.String()
				require.JSONEq(t, expectedJSON, actualJSON)
			} else {
				fmt.Println("inside else part")
				jsonData, err := types.ReadJsonFromFile(tc.TestJsonFile)
				require.NoError(t, err)
				expectedJSON := removeFieldFromJSON(string(jsonData), "loggedat")
				fmt.Println(">>>>>>>>>>>>>>>expected :", expectedJSON)
				actualJSON := removeFieldFromJSON(res.Body.String(), "loggedat")
				fmt.Println(">>>>>>>>>>>>>>>actual :", expectedJSON)
				require.JSONEq(t, expectedJSON, actualJSON)

			}
		})
	}

}

func wfInstanceNewTestcase() []testutils.TestCasesStruct {
	var slice int32 = 2
	entityID := "0eb8da50-aece-11ee-b168-3b192f7cd2b6"
	entityID1 := "0eb8da50-aece-11ee-b168-3b192f7cd2b667"
	workflow := "discountcheck"
	workflow1 := "temp"
	trace := 0
	parent := int32(917)
	app := "retailBANK"
	wfInstanceNewTestcase := []testutils.TestCasesStruct{

		//1st test case
		{
			Name: "SUCCESS- Single step",
			RequestPayload: wscutils.Request{
				Data: wfinstance.WFInstanceNewRequest{
					Slice:    slice,
					App:      app,
					EntityID: entityID,
					Entity: map[string]string{
						"class":        "inventoryitems",
						"mrp":          "200.00",
						"fullname":     "belampally",
						"ageinstock":   "2",
						"inventoryqty": "2",
					},
					Workflow: workflow,
					Trace:    &trace,
					Parent:   &parent,
				},
			},
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/single_step_success_response.json",
		},

		// 2nd test case
		{
			Name: "Success- Multiple steps",
			RequestPayload: wscutils.Request{
				Data: wfinstance.WFInstanceNewRequest{
					Slice:    slice,
					App:      app,
					EntityID: entityID1,
					Entity: map[string]string{
						"class":        "members",
						"mrp":          "200.00",
						"fullname":     "belampally",
						"ageinstock":   "2",
						"inventoryqty": "2",
					},
					Workflow: workflow1,
					Trace:    &trace,
					Parent:   &parent,
				},
			},
			ExpectedHttpCode: http.StatusOK,
			TestJsonFile:     "./data/multiple_steps_success_response.json",
		},
	}
	return wfInstanceNewTestcase
}

func removeFieldFromJSON(jsonStr string, field string) string {
	// Parse JSON
	var obj wscutils.Response
	obj.Data = wfinstance.WFInstanceNewRequest{}

	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		return jsonStr
	}

	// Convert the object to a map using reflection
	if obj.Data != nil {
		// Remove the specified field from 'Data'
		delete(obj.Data.(map[string]interface{}), field)
	}

	// Convert back to JSON
	modifiedJSON, err := json.Marshal(obj)
	if err != nil {
		return jsonStr
	}

	return string(modifiedJSON)
}
