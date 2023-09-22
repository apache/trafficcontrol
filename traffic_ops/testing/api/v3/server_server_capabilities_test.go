package v3

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
	"net/url"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestServerServerCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments, ServerCapabilities, ServerServerCapabilities, DeliveryServicesRequiredCapabilities}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.ServerServerCapability]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerServerCapabilitiesSort()),
				},
				"OK when VALID SERVERID parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"serverId": {strconv.Itoa(GetServerID(t, "dtrc-edge-01")())}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerServerCapabilitiesFields(map[string]interface{}{"ServerID": GetServerID(t, "dtrc-edge-01")()})),
				},
				"OK when VALID SERVERHOSTNAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"serverHostName": {"atlanta-edge-16"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerServerCapabilitiesFields(map[string]interface{}{"Server": "atlanta-edge-16"})),
				},
				"OK when VALID SERVERCAPABILITY parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"serverCapability": {"asdf"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerServerCapabilitiesFields(map[string]interface{}{"ServerCapability": "asdf"})),
				},
			},
			"POST": {
				"BAD REQUEST when ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody: tc.ServerServerCapability{
						ServerID:         util.Ptr(GetServerID(t, "dtrc-mid-01")()),
						ServerCapability: util.Ptr("disk"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING SERVER ID": {
					ClientSession: TOSession,
					RequestBody: tc.ServerServerCapability{
						ServerCapability: util.Ptr("disk"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING SERVER CAPABILITY": {
					ClientSession: TOSession,
					RequestBody: tc.ServerServerCapability{
						ServerID: util.Ptr(GetServerID(t, "dtrc-mid-01")()),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"NOT FOUND when SERVER CAPABILITY DOESNT EXIST": {
					ClientSession: TOSession,
					RequestBody: tc.ServerServerCapability{
						ServerID:         util.Ptr(GetServerID(t, "dtrc-mid-01")()),
						ServerCapability: util.Ptr("bogus"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"NOT FOUND when SERVER DOESNT EXIST": {
					ClientSession: TOSession,
					RequestBody: tc.ServerServerCapability{
						ServerID:         util.Ptr(99999999),
						ServerCapability: util.Ptr("bogus"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when SERVER TYPE NOT EDGE or MID": {
					ClientSession: TOSession,
					RequestBody: tc.ServerServerCapability{
						ServerID:         util.Ptr(GetServerID(t, "trafficvault")()),
						ServerCapability: util.Ptr("bogus"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"DELETE": {
				"OK when NOT the LAST SERVER of CACHE GROUP of TOPOLOGY DS which has REQUIRED CAPABILITIES": {
					ClientSession: TOSession,
					RequestParams: url.Values{"serverId": {strconv.Itoa(GetServerID(t, "dtrc-edge-01")())}, "serverCapability": {"ram"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when LAST SERVER of CACHE GROUP of TOPOLOGY DS which has REQUIRED CAPABILITIES": {
					ClientSession: TOSession,
					RequestParams: url.Values{"serverId": {strconv.Itoa(GetServerID(t, "edge-in-cdn1-only")())}, "serverCapability": {"ram"}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SERVER ASSIGNED TO DS with REQUIRED CAPABILITIES": {
					ClientSession: TOSession,
					RequestParams: url.Values{"serverId": {strconv.Itoa(GetServerID(t, "atlanta-org-2")())}, "serverCapability": {"bar"}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"NOT FOUND when SERVER SERVER CAPABILITY DOESNT EXIST": {
					ClientSession: TOSession,
					RequestParams: url.Values{"serverId": {strconv.Itoa(GetServerID(t, "atlanta-org-1")())}, "serverCapability": {"doesntexist"}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when MISSING SERVER CAPABILITY": {
					ClientSession: TOSession,
					RequestParams: url.Values{"serverId": {strconv.Itoa(GetServerID(t, "atlanta-org-1")())}, "serverCapability": {""}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var serverId *int
					var serverHostName *string
					var serverCapability *string

					if val, ok := testCase.RequestParams["serverId"]; ok {
						id, _ := strconv.Atoi(val[0])
						serverId = &id
					}
					if val, ok := testCase.RequestParams["serverCapability"]; ok {
						serverCapability = &val[0]
					}
					if val, ok := testCase.RequestParams["serverHostName"]; ok {
						serverHostName = &val[0]
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetServerServerCapabilitiesWithHdr(serverId, serverHostName, serverCapability, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateServerServerCapability(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteServerServerCapability(*serverId, *serverCapability)
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

func validateServerServerCapabilitiesFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Server Server Capabilities response to not be nil.")
		serverServerCapabilityResponse := resp.([]tc.ServerServerCapability)
		for field, expected := range expectedResp {
			for _, serverServerCapability := range serverServerCapabilityResponse {
				switch field {
				case "Server":
					assert.RequireNotNil(t, serverServerCapability.Server, "Expected Server to not be nil.")
					assert.Equal(t, expected, *serverServerCapability.Server, "Expected Server to be %v, but got %s", expected, *serverServerCapability.Server)
				case "ServerCapability":
					assert.RequireNotNil(t, serverServerCapability.ServerCapability, "Expected Server Capability to not be nil.")
					assert.Equal(t, expected, *serverServerCapability.ServerCapability, "Expected ServerCapability to be %v, but got %s", expected, *serverServerCapability.ServerCapability)
				case "ServerID":
					assert.RequireNotNil(t, serverServerCapability.ServerID, "Expected Server ID to not be nil.")
					assert.Equal(t, expected, *serverServerCapability.ServerID, "Expected ServerID to be %v, but got %d", expected, *serverServerCapability.ServerID)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateServerServerCapabilitiesSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Server Server Capabilities response to not be nil.")
		var serverNames []string
		serverServerCapabilityResponse := resp.([]tc.ServerServerCapability)
		for _, serverServerCapability := range serverServerCapabilityResponse {
			assert.RequireNotNil(t, serverServerCapability.Server, "Expected Server to not be nil.")
			serverNames = append(serverNames, *serverServerCapability.Server)
		}
		assert.Equal(t, true, sort.StringsAreSorted(serverNames), "List is not sorted by server names: %v", serverNames)
	}
}

func CreateTestServerServerCapabilities(t *testing.T) {
	for _, ssc := range testData.ServerServerCapabilities {
		assert.RequireNotNil(t, ssc.Server, "Expected Server to not be nil.")
		assert.RequireNotNil(t, ssc.ServerCapability, "Expected Server Capability to not be nil.")
		serverID := GetServerID(t, *ssc.Server)()
		ssc.ServerID = &serverID
		resp, _, err := TOSession.CreateServerServerCapability(ssc)
		assert.RequireNoError(t, err, "Could not associate Capability '%s' with server '%s': %v - alerts: %+v", *ssc.ServerCapability, *ssc.Server, err, resp.Alerts)
	}
}

func DeleteTestServerServerCapabilities(t *testing.T) {
	sscs, _, err := TOSession.GetServerServerCapabilitiesWithHdr(nil, nil, nil, nil)
	assert.RequireNoError(t, err, "Cannot get server server capabilities: %v", err)
	for _, ssc := range sscs {
		assert.RequireNotNil(t, ssc.Server, "Expected Server to not be nil.")
		assert.RequireNotNil(t, ssc.ServerCapability, "Expected Server Capability to not be nil.")
		alerts, _, err := TOSession.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability)
		assert.NoError(t, err, "Could not remove Capability '%s' from server '%s' (#%d): %v - alerts: %+v", *ssc.ServerCapability, *ssc.Server, *ssc.ServerID, err, alerts.Alerts)
	}
}
