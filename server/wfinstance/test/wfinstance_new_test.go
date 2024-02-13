package wfinstance_test

// import (
// 	"bytes"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/remiges-tech/alya/wscutils"
// 	"github.com/remiges-tech/crux/server/wfinstance"
// 	"github.com/remiges-tech/crux/testutils"
// 	"github.com/stretchr/testify/require"
// )

// func TestWFInstanceNew(t *testing.T) {
// 	testCases := wfInstanceNewTestcase()
// 	for _, tc := range testCases {
// 		t.Run(tc.Name, func(t *testing.T) {

// 			// Setting up buffer
// 			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))

// 			res := httptest.NewRecorder()
// 			req, err := http.NewRequest(http.MethodPost, "/wfinstancenew", payload)
// 			require.NoError(t, err)

// 			r.ServeHTTP(res, req)

// 			require.Equal(t, tc.ExpectedHttpCode, res.Code)
// 			if tc.ExpectedResult != nil {
// 				jsonData := testutils.MarshalJson(tc.ExpectedResult)
// 				require.JSONEq(t, string(jsonData), res.Body.String())
// 			} else {
// 				jsonData, err := testutils.ReadJsonFromFile(tc.TestJsonFile)
// 				require.NoError(t, err)
// 				// require.JSONEq(t, string(jsonData), res.Body.String())
// 				// response = unmarshal(res.Body.String(), &response)
// 				// unmarshal(strnig(jsonData), &expected)
// 				// response.tasks == expected.tasks
// 				// response.workflow
// 			}
// 		})
// 	}

// }

// func wfInstanceNewTestcase() []testutils.TestCasesStruct {
// 	var slice int32 = 2
// 	entityID := "0eb8da50-aece-11ee-b168-3b192f7cd2b6"
// 	workflow := "discountcheck"
// 	//workflow1 := "temp_set"
// 	app := "retailBANK"
// 	wfInstanceNewTestcase := []testutils.TestCasesStruct{
// 		// 1st test case
// 		// {
// 		// 	Name: "SUCCESS- Multi Step wfinstancenew by valid req ",
// 		// 	RequestPayload: wscutils.Request{
// 		// 		Data: wfinstance.WFInstanceNewRequest{
// 		// 			Slice:    &slice,
// 		// 			App:      &app,
// 		// 			EntityID: &entityID,
// 		// 			Entity: map[string]string{
// 		// 				"class":        "inventoryitems",
// 		// 				"mrp":          "200.00",
// 		// 				"fullname":     "belampally",
// 		// 				"ageinstock":   "2",
// 		// 				"inventoryqty": "2",
// 		// 			},
// 		// 			Workflow: &workflow,
// 		// 			Trace:    0,
// 		// 			Parent:   952,
// 		// 		},
// 		// 	},
// 		// 	ExpectedHttpCode: http.StatusOK,
// 		// 	TestJsonFile:     "./data/multiple_steps_success_response.json",
// 		// },

// 		//2nd  test case
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
// 					Workflow: &workflow,
// 					Trace:    0,
// 					Parent:   917,
// 				},
// 			},
// 			ExpectedHttpCode: http.StatusOK,
// 			TestJsonFile:     "./data/single_step_success_response.json",
// 			// ExpectedResult: &wscutils.Response{
// 			// 	Status: wscutils.SuccessStatus,
// 			// 	Data: wfinstance.WFInstanceNewResponse{
// 			// 		Tasks:     []map[string]int32{map[string]int32{"diwalisale": 3}},
// 			// 		Nextstep:  "",
// 			// 		Loggedat:  pgtype.Timestamp{Time: time.Now().In(time.UTC), Valid: true},
// 			// 		Subflows:  &map[string]string{},
// 			// 		Tracedata: nil,
// 			// 	},
// 			// 	Messages: nil,
// 			// },
// 		},

// 		// 3rd test case
// 		{
// 			Name: "ERROR- Single step",
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
// 					Trace:    0,
// 					Parent:   917,
// 				},
// 			},
// 			ExpectedHttpCode: http.StatusBadRequest,
// 			TestJsonFile:     "./data/record_already_exists_res.json",
// 		},
// 	}
// 	return wfInstanceNewTestcase
// }
