package v5

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
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

func TestServersIDQueueUpdate(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers}, func() {

		methodTests := utils.V5TestCase{
			"POST": {
				"OK when VALID QUEUE request": {
					EndpointId:    GetServerID(t, "atlanta-edge-01"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"action": "queue",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServerQueueUpdateFields(map[string]interface{}{"Action": "queue", "ServerID": GetServerID(t, "atlanta-edge-01")()}),
						validateUpdPendingSpecificServers(map[string]bool{"atlanta-edge-01": true})),
				},
				"OK when VALID DEQUEUE request": {
					EndpointId:    GetServerID(t, "atlanta-edge-01"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"action": "dequeue",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServerQueueUpdateFields(map[string]interface{}{"Action": "dequeue", "ServerID": GetServerID(t, "atlanta-edge-01")()}),
						validateUpdPendingSpecificServers(map[string]bool{"atlanta-edge-01": false})),
				},
				/* COMMENTED UNTIL ISSUE IS FIXED:
					https://github.com/apache/trafficcontrol/issues/6691
					https://github.com/apache/trafficcontrol/issues/6801
				"NOT FOUND when NON-EXISTENT SERVER": {
					EndpointId:    func() int { return 999999 },
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"action": "queue",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				}, */
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var queueUpdate bool
					if val, ok := testCase.RequestBody["action"]; ok {
						if val == "queue" {
							queueUpdate = true
						} else if val == "dequeue" {
							queueUpdate = false
						}
					}

					switch method {
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.SetServerQueueUpdate(testCase.EndpointId(), queueUpdate, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func validateServerQueueUpdateFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Server Queue Update response to not be nil.")
		serverQueueUpdate := resp.(tc.ServerQueueUpdate)
		for field, expected := range expectedResp {
			switch field {
			case "Action":
				assert.Equal(t, expected, serverQueueUpdate.Action, "Expected Action to be %v, but got %s", expected, serverQueueUpdate.Action)
			case "ServerID":
				assert.Equal(t, util.JSONIntStr(expected.(int)), serverQueueUpdate.ServerID, "Expected ServerID to be %v, but got %d", expected, serverQueueUpdate.ServerID)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}
