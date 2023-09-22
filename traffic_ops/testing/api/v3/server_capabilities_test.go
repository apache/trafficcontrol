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
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestServerCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{ServerCapabilities}, func() {

		methodTests := utils.V3TestCaseT[tc.ServerCapability]{
			"GET": {
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerCapabilitiesSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"ram"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"BAD REQUEST when INVALID NAME": {
					ClientSession: TOSession,
					RequestBody:   tc.ServerCapability{Name: "b@dname"},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						if name == "OK when VALID NAME parameter" {
							resp, reqInf, err := testCase.ClientSession.GetServerCapabilityWithHdr(testCase.RequestParams["name"][0], testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						} else {
							resp, reqInf, err := testCase.ClientSession.GetServerCapabilitiesWithHdr(testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						}
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateServerCapability(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								if resp != nil {
									check(t, reqInf, resp.Response, resp.Alerts, err)
								}
							}
						})
					}
				}
			})
		}
	})
}

func validateServerCapabilitiesSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Server Capabilities response to not be nil.")
		var serverCapabilityNames []string
		serverCapabilitiesResp := resp.([]tc.ServerCapability)
		for _, serverCapability := range serverCapabilitiesResp {
			serverCapabilityNames = append(serverCapabilityNames, serverCapability.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(serverCapabilityNames), "List is not sorted by their names: %v", serverCapabilityNames)
	}
}

func CreateTestServerCapabilities(t *testing.T) {
	for _, sc := range testData.ServerCapabilities {
		resp, _, err := TOSession.CreateServerCapability(sc)
		assert.RequireNoError(t, err, "Unexpected error creating Server Capability '%s': %v - alerts: %+v", sc.Name, err, resp.Alerts)
	}
}

func DeleteTestServerCapabilities(t *testing.T) {
	serverCapabilities, _, err := TOSession.GetServerCapabilitiesWithHdr(nil)
	assert.NoError(t, err, "Cannot get Server Capabilities: %v", err)

	for _, serverCapability := range serverCapabilities {
		alerts, _, err := TOSession.DeleteServerCapability(serverCapability.Name)
		assert.NoError(t, err, "Unexpected error deleting Server Capability '%s': %v - alerts: %+v", serverCapability.Name, err, alerts.Alerts)
		// Retrieve the Server Capability to see if it got deleted
		getServerCapability, _, err := TOSession.GetServerCapabilityWithHdr(serverCapability.Name, nil)
		assert.Error(t, err, "Expected error getting Server Capability '%s' after deletion: %v", serverCapability.Name, err)
		assert.Equal(t, (*tc.ServerCapability)(nil), getServerCapability, "Expected Server Capability '%s' to be deleted, but it was found in Traffic Ops", serverCapability.Name)
	}
}
