package v4

/*

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	totest "github.com/apache/trafficcontrol/v8/lib/go-tc/totestv4"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestSteeringTargets(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, Users, SteeringTargets}, func() {

		steeringUserSession := utils.CreateV4Session(t, Config.TrafficOps.URL, "steering", "pa$$word", Config.Default.Session.TimeoutInSecs)

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.SteeringTargetNullable]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					EndpointID:    totest.GetDeliveryServiceId(t, TOSession, "ds1"),
					ClientSession: steeringUserSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					EndpointID:    totest.GetDeliveryServiceId(t, TOSession, "ds1"),
					ClientSession: steeringUserSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateSteeringTargetFields(map[string]interface{}{"DeliveryService": "ds1", "DeliveryServiceID": uint64(totest.GetDeliveryServiceId(t, TOSession, "ds1")()),
							"Target": "ds2", "TargetID": uint64(totest.GetDeliveryServiceId(t, TOSession, "ds2")()), "Type": "STEERING_WEIGHT", "TypeID": GetTypeID(t, "STEERING_WEIGHT")(), "Value": util.JSONIntStr(42)})),
				},
				"OK when CHANGES made": {
					EndpointID:    totest.GetDeliveryServiceId(t, TOSession, "ds1"),
					ClientSession: steeringUserSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					ClientSession: steeringUserSession,
					RequestBody: tc.SteeringTargetNullable{
						DeliveryServiceID: util.Ptr(uint64(totest.GetDeliveryServiceId(t, TOSession, "ds3")())),
						TargetID:          util.Ptr(uint64(totest.GetDeliveryServiceId(t, TOSession, "ds4")())),
						Value:             util.Ptr(util.JSONIntStr(-12345)),
						TypeID:            util.Ptr(GetTypeID(t, "STEERING_WEIGHT")()),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateSteeringTargetUpdateCreateFields(totest.GetDeliveryServiceId(t, TOSession, "ds3")(),
							map[string]interface{}{"DeliveryService": "ds3", "DeliveryServiceID": uint64(totest.GetDeliveryServiceId(t, TOSession, "ds3")()),
								"Target": "ds4", "TargetID": uint64(totest.GetDeliveryServiceId(t, TOSession, "ds4")()), "Type": "STEERING_WEIGHT",
								"TypeID": GetTypeID(t, "STEERING_WEIGHT")(), "Value": util.JSONIntStr(-12345)})),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					ClientSession: steeringUserSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: tc.SteeringTargetNullable{
						DeliveryServiceID: util.Ptr(uint64(totest.GetDeliveryServiceId(t, TOSession, "ds3")())),
						TargetID:          util.Ptr(uint64(totest.GetDeliveryServiceId(t, TOSession, "ds4")())),
						Value:             util.Ptr(util.JSONIntStr(-12345)),
						Type:              util.Ptr("STEERING_WEIGHT"),
						TypeID:            util.Ptr(GetTypeID(t, "STEERING_WEIGHT")()),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					ClientSession: steeringUserSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					RequestBody: tc.SteeringTargetNullable{
						DeliveryServiceID: util.Ptr(uint64(totest.GetDeliveryServiceId(t, TOSession, "ds3")())),
						TargetID:          util.Ptr(uint64(totest.GetDeliveryServiceId(t, TOSession, "ds4")())),
						Value:             util.Ptr(util.JSONIntStr(-12345)),
						Type:              util.Ptr("STEERING_WEIGHT"),
						TypeID:            util.Ptr(GetTypeID(t, "STEERING_WEIGHT")()),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetSteeringTargets(testCase.EndpointID(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateSteeringTarget(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateSteeringTarget(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							if _, ok := testCase.RequestOpts.QueryParameters["targetID"]; !ok {
								t.Fatalf("Query Parameter: \"name\" is required for PUT method tests.")
							}
							targetID, err := strconv.Atoi(testCase.RequestOpts.QueryParameters["targetID"][0])
							assert.RequireNoError(t, err, "Expected no error converting string to int for target ID: %v", err)
							alerts, reqInf, err := testCase.ClientSession.DeleteSteeringTarget(testCase.EndpointID(), targetID, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					}
				}
			})
		}
	})

}

func validateSteeringTargetFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Steering Targets response to not be nil.")
		steeringTargetsResp := resp.([]tc.SteeringTargetNullable)
		for field, expected := range expectedResp {
			for _, steeringTarget := range steeringTargetsResp {
				switch field {
				case "DeliveryService":
					assert.RequireNotNil(t, steeringTarget.DeliveryService, "Expected DeliveryService to not be nil.")
					assert.Equal(t, expected, string(*steeringTarget.DeliveryService), "Expected DeliveryService to be %v, but got %s", expected, *steeringTarget.DeliveryService)
				case "DeliveryServiceID":
					assert.RequireNotNil(t, steeringTarget.DeliveryServiceID, "Expected DeliveryServiceID to not be nil.")
					assert.Equal(t, expected, *steeringTarget.DeliveryServiceID, "Expected DeliveryServiceID to be %v, but got %s", expected, *steeringTarget.DeliveryServiceID)
				case "Target":
					assert.RequireNotNil(t, steeringTarget.Target, "Expected Target to not be nil.")
					assert.Equal(t, expected, string(*steeringTarget.Target), "Expected Target to be %v, but got %s", expected, *steeringTarget.Target)
				case "TargetID":
					assert.RequireNotNil(t, steeringTarget.TargetID, "Expected TargetID to not be nil.")
					assert.Equal(t, expected, *steeringTarget.TargetID, "Expected TargetID to be %v, but got %s", expected, *steeringTarget.TargetID)
				case "Type":
					assert.RequireNotNil(t, steeringTarget.Type, "Expected Type to not be nil.")
					assert.Equal(t, expected, *steeringTarget.Type, "Expected Type to be %v, but got %s", expected, *steeringTarget.Type)
				case "TypeID":
					assert.RequireNotNil(t, steeringTarget.Type, "Expected TypeID to not be nil.")
					assert.Equal(t, expected, *steeringTarget.TypeID, "Expected TypeID to be %v, but got %s", expected, *steeringTarget.TypeID)
				case "Value":
					assert.RequireNotNil(t, steeringTarget.Value, "Expected Value to not be nil.")
					assert.Equal(t, expected, *steeringTarget.Value, "Expected Value to be %v, but got %s", expected, *steeringTarget.Value)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateSteeringTargetUpdateCreateFields(dsId int, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		steeringTargets, _, err := TOSession.GetSteeringTargets(dsId, client.RequestOptions{})
		assert.RequireNoError(t, err, "Error getting Steering Targets: %v - alerts: %+v", err, steeringTargets.Alerts)
		assert.RequireEqual(t, 1, len(steeringTargets.Response), "Expected one Steering Target returned Got: %d", len(steeringTargets.Response))
		validateSteeringTargetFields(expectedResp)(t, toclientlib.ReqInf{}, steeringTargets.Response, tc.Alerts{}, nil)
	}
}
