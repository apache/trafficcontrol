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
	"net/url"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestCDNQueueUpdates(t *testing.T) {
	WithObjs(t, []TCObj{Types, CDNs, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers}, func() {

		methodTests := utils.TestCase[client.Session, client.RequestOptions, bool]{
			"POST": {
				"OK when VALID TYPE parameter": {
					EndpointID:    GetCDNID(t, "cdn1"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"type": {"EDGE"}}},
					RequestBody:   true,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServersUpdatePending(GetCDNID(t, "cdn1")(), map[string]string{"type": "EDGE"})),
				},
				"OK when VALID PROFILE parameter": {
					EndpointID:    GetCDNID(t, "cdn1"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"profile": {"EDGE1"}}},
					RequestBody:   true,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServersUpdatePending(GetCDNID(t, "cdn1")(), map[string]string{"profileName": "EDGE1"})),
				},
			},
		}
		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "POST":
						t.Run(name, func(t *testing.T) {
							// Clear updates on all associated cdn servers to begin with
							_, _, err := TOSession.QueueUpdatesForCDN(testCase.EndpointID(), false, client.RequestOptions{})
							assert.RequireNoError(t, err, "Failed to clear updates for cdn %d", testCase.EndpointID())
							resp, reqInf, err := testCase.ClientSession.QueueUpdatesForCDN(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					}
				}
			})
		}
	})
}

func validateServersUpdatePending(cdnID int, params map[string]string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, _ tc.Alerts, _ error) {
		// Get all the servers for the same CDN and type as that of the first server
		serverIDMap := make(map[int]bool, 0)
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("cdn", strconv.Itoa(cdnID))
		for k, v := range params {
			opts.QueryParameters.Set(k, v)
		}

		servers, _, err := TOSession.GetServers(opts)
		assert.RequireNoError(t, err, "Couldn't get servers by cdn and parameters: %v", err)
		assert.RequireGreaterOrEqual(t, len(servers.Response), 1, "expected atleast one server in response, got %d", len(servers.Response))

		for _, server := range servers.Response {
			assert.Equal(t, true, server.UpdatePending(), "Expected updates to be queued on all the servers filtered by CDN and parameter, but %s didn't queue updates", server.HostName)
			serverIDMap[server.ID] = true
		}

		// Make sure that the servers that are not filtered by the above criteria do not have updates queued
		allServersResp, _, err := TOSession.GetServers(client.NewRequestOptions())
		assert.RequireNoError(t, err, "Couldn't get all servers: %v", err)

		for _, server := range allServersResp.Response {
			if _, ok := serverIDMap[server.ID]; !ok {
				assert.Equal(t, false, server.UpdatePending(), "Did not expect server %s to have queued updates", server.HostName)
			}
		}
	}
}
