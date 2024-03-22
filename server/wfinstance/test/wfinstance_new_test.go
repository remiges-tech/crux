package wfinstance_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/remiges-tech/alya/wscutils"
// 	"github.com/remiges-tech/crux/server"
// 	"github.com/remiges-tech/crux/server/wfinstance"
// 	"github.com/remiges-tech/crux/testutils"
// 	"github.com/stretchr/testify/require"
// )

// func TestWFinstanceNew(t *testing.T) {
// 	testCases := wfInstanceNewTestcase()
// 	for _, tc := range testCases {
// 		t.Run(tc.Name, func(t *testing.T) {

// 			switch tc.Name {
// 			case "SUCCESS- Multiple steps":
// 				wfinstance.SWITCH = wfinstance.ActionSet{
// 					Tasks:      []string{"diwalisale", "yearendsale"},
// 					Properties: map[string]string{"nextstep": "coupondistribution"},
// 				}
// 			case "ERROR- Invalid property attributes":
// 				wfinstance.SWITCH = wfinstance.ActionSet{
// 					Tasks:      []string{"discount"},
// 					Properties: map[string]string{"shipby": "fedex"},
// 				}

// 			case "SUCCESS- done attribute present in domatch() response":
// 				wfinstance.SWITCH = wfinstance.ActionSet{
// 					Tasks:      []string{"discount", "yearendsale"},
// 					Properties: map[string]string{"done": "true"},
// 				}
// 			default:
// 				wfinstance.SWITCH = wfinstance.ActionSet{
// 					Tasks:      []string{"diwalisale"},
// 					Properties: map[string]string{"nextstep": "coupondistribution"},
// 				}
// 			}

// 			// Setting up buffer
// 			payload := bytes.NewBuffer(server.MarshalJson(tc.RequestPayload))

// 			res := httptest.NewRecorder()
// 			req, err := http.NewRequest(http.MethodPost, "/wfinstancenew", payload)
// 			require.NoError(t, err)

// 			r.ServeHTTP(res, req)

// 			require.Equal(t, tc.ExpectedHttpCode, res.Code)
// 			if tc.ExpectedResult != nil {
// 				jsonData := server.MarshalJson(tc.ExpectedResult)
// 				expectedJSON := string(jsonData)
// 				actualJSON := res.Body.String()
// 				require.JSONEq(t, expectedJSON, actualJSON)
// 			} else {
// 				jsonData, err := server.ReadJsonFromFile(tc.TestJsonFile)
// 				require.NoError(t, err)
// 				expectedJSON := removeFieldFromJSON(string(jsonData), "loggedat")
// 				actualJSON := removeFieldFromJSON(res.Body.String(), "loggedat")
// 				require.JSONEq(t, expectedJSON, actualJSON)

// 			}
// 		})
// 	}

// }

// func wfInstanceNewTestcase() []testutils.TestCasesStruct {
// 	var slice int32 = 2
// 	entityID := "0eb8da50-aece-11ee-b168-3b192f7cd2b6"
// 	entityID1 := "0eb8da50-aece-11ee-b168-3b192f7cd2b667"
// 	workflow := "discountcheck"
// 	workflow1 := "temp"
// 	trace := 0
// 	parent := int32(917)
// 	app := "retailBANK"
// 	wfInstanceNewTestcase := []testutils.TestCasesStruct{

// 		// 1st test case
// 		{
// 			Name: "ERROR- Invalid request",
// 			RequestPayload: wscutils.Request{
// 				Data: wfinstance.WFInstanceNewRequest{
// 					Slice:    slice,
// 					App:      app,
// 					EntityID: entityID,
// 					Entity: map[string]string{
// 						"class":        "inventoryitemss",
// 						"mrp":          "200.00",
// 						"fullname":     "belampally",
// 						"ageinstock":   "2",
// 						"inventoryqty": "2",
// 					},
// 					Workflow: workflow,
// 					Trace:    &trace,
// 					Parent:   &parent,
// 				},
// 			},
// 			ExpectedHttpCode: http.StatusBadRequest,
// 			TestJsonFile:     "./data/invalid_wfinstancenew_request_response.json",
// 		},

// 		// 2nd test case
// 		{
// 			Name: "SUCCESS- done attribute present in domatch() response",
// 			RequestPayload: wscutils.Request{
// 				Data: wfinstance.WFInstanceNewRequest{
// 					Slice:    slice,
// 					App:      app,
// 					EntityID: entityID,
// 					Entity: map[string]string{
// 						"class":        "inventoryitems",
// 						"mrp":          "200.00",
// 						"fullname":     "belampally",
// 						"ageinstock":   "2",
// 						"inventoryqty": "2",
// 					},
// 					Workflow: workflow,
// 					Trace:    &trace,
// 					Parent:   &parent,
// 				},
// 			},
// 			ExpectedHttpCode: http.StatusOK,
// 			TestJsonFile:     "./data/done_attribute_response.json",
// 		},

// 		// 3rd test case
// 		{
// 			Name: "ERROR- Invalid property attributes",
// 			RequestPayload: wscutils.Request{
// 				Data: wfinstance.WFInstanceNewRequest{
// 					Slice:    slice,
// 					App:      app,
// 					EntityID: entityID,
// 					Entity: map[string]string{
// 						"class":        "inventoryitems",
// 						"mrp":          "200.00",
// 						"fullname":     "belampally",
// 						"ageinstock":   "2",
// 						"inventoryqty": "2",
// 					},
// 					Workflow: workflow,
// 					Trace:    &trace,
// 					Parent:   &parent,
// 				},
// 			},
// 			ExpectedHttpCode: http.StatusBadRequest,
// 			TestJsonFile:     "./data/invalid_property_attributes_response.json",
// 		},
// 		// 4th test case
// 		{
// 			Name: "SUCCESS- Single step",
// 			RequestPayload: wscutils.Request{
// 				Data: wfinstance.WFInstanceNewRequest{
// 					Slice:    slice,
// 					App:      app,
// 					EntityID: entityID,
// 					Entity: map[string]string{
// 						"class":        "inventoryitems",
// 						"mrp":          "200.00",
// 						"fullname":     "belampally",
// 						"ageinstock":   "2",
// 						"inventoryqty": "2",
// 					},
// 					Workflow: workflow,
// 					Trace:    &trace,
// 					Parent:   &parent,
// 				},
// 			},
// 			ExpectedHttpCode: http.StatusOK,
// 			TestJsonFile:     "./data/single_step_success_response.json",
// 		},

// 		// 5th test case
// 		{
// 			Name: "SUCCESS- Multiple steps",
// 			RequestPayload: wscutils.Request{
// 				Data: wfinstance.WFInstanceNewRequest{
// 					Slice:    slice,
// 					App:      app,
// 					EntityID: entityID1,
// 					Entity: map[string]string{
// 						"class":        "members",
// 						"mrp":          "200.00",
// 						"fullname":     "belampally",
// 						"ageinstock":   "2",
// 						"inventoryqty": "2",
// 					},
// 					Workflow: workflow1,
// 					Trace:    &trace,
// 					Parent:   &parent,
// 				},
// 			},
// 			ExpectedHttpCode: http.StatusOK,
// 			TestJsonFile:     "./data/multiple_steps_success_response.json",
// 		},
// 		// 6th test case
// 		{
// 			Name: "ERROR- Instance already exist in database",
// 			RequestPayload: wscutils.Request{
// 				Data: wfinstance.WFInstanceNewRequest{
// 					Slice:    slice,
// 					App:      app,
// 					EntityID: entityID1,
// 					Entity: map[string]string{
// 						"class":        "members",
// 						"mrp":          "200.00",
// 						"fullname":     "belampally",
// 						"ageinstock":   "2",
// 						"inventoryqty": "2",
// 					},
// 					Workflow: workflow1,
// 					Trace:    &trace,
// 					Parent:   &parent,
// 				},
// 			},
// 			ExpectedHttpCode: http.StatusBadRequest,
// 			TestJsonFile:     "./data/instance_already_exist_response.json",
// 		},
// 	}
// 	return wfInstanceNewTestcase
// }

// func removeFieldFromJSON(jsonStr string, field string) string {
// 	var obj wscutils.Response
// 	obj.Data = wfinstance.WFInstanceNewRequest{}

// 	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
// 		return jsonStr
// 	}

// 	if obj.Data != nil {
// 		delete(obj.Data.(map[string]interface{}), field)
// 	}

// 	modifiedJSON, err := json.Marshal(obj)
// 	if err != nil {
// 		return jsonStr
// 	}

// 	return string(modifiedJSON)
// }
