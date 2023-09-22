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
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	totest "github.com/apache/trafficcontrol/v8/lib/go-tc/totestv4"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestCacheGroupsDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, CacheGroupsDeliveryServices}, func() {

		methodTests := utils.TestCase[client.Session, client.RequestOptions, []int]{
			"POST": {
				"BAD REQUEST assigning TOPOLOGY-BASED DS to CACHEGROUP": {
					EndpointID:    totest.GetCacheGroupId(t, TOSession, "cachegroup3"),
					ClientSession: TOSession,
					RequestBody:   []int{totest.GetDeliveryServiceId(t, TOSession, "top-ds-in-cdn1")()},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when valid request": {
					EndpointID:    totest.GetCacheGroupId(t, TOSession, "cachegroup3"),
					ClientSession: TOSession,
					RequestBody: []int{
						totest.GetDeliveryServiceId(t, TOSession, "ds1")(),
						totest.GetDeliveryServiceId(t, TOSession, "ds2")(),
						totest.GetDeliveryServiceId(t, TOSession, "ds3")(),
						totest.GetDeliveryServiceId(t, TOSession, "ds3")(),
						totest.GetDeliveryServiceId(t, TOSession, "DS5")(),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCGDSServerAssignments()),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.SetCacheGroupDeliveryServices(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
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

func validateCGDSServerAssignments() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		cgDsResp := resp.(tc.CacheGroupPostDSResp)
		opts := client.NewRequestOptions()
		for _, serverName := range cgDsResp.ServerNames {
			opts.QueryParameters.Set("hostName", string(serverName))
			resp, _, err := TOSession.GetServers(opts)
			assert.NoError(t, err, "Error: Getting server: %v - alerts: %+v", err, resp.Alerts)
			assert.Equal(t, len(resp.Response), 1, "Error: Getting servers: expected 1 got %v", len(resp.Response))

			serverDSes, _, err := TOSession.GetDeliveryServicesByServer(*resp.Response[0].ID, client.RequestOptions{})
			assert.NoError(t, err, "Error: Getting Delivery Service Servers #%d: %v - alerts: %+v", *resp.Response[0].ID, err, serverDSes.Alerts)
			for _, dsID := range cgDsResp.DeliveryServices {
				found := false
				for _, serverDS := range serverDSes.Response {
					if *serverDS.ID == dsID {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("POST succeeded, but didn't assign delivery service %v to server", dsID)
				}
			}
		}
	}
}
