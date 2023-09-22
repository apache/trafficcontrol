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

package v4

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	totest "github.com/apache/trafficcontrol/v8/lib/go-tc/totestv4"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestDeliveryServicesRegexes(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, DeliveryServicesRegexes}, func() {

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.DeliveryServiceRegexPost]{
			"GET": {
				"OK when VALID request": {
					EndpointID:    totest.GetDeliveryServiceId(t, TOSession, "ds1"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(3)),
				},
				"OK when VALID ID parameter": {
					EndpointID:    totest.GetDeliveryServiceId(t, TOSession, "ds1"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"id": {strconv.Itoa(getDSRegexID(t, "ds1"))}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1)),
				},
			},
			"POST": {
				"BAD REQUEST when MISSING REGEX PATTERN": {
					EndpointID: totest.GetDeliveryServiceId(t, TOSession, "ds1"), ClientSession: TOSession,
					RequestBody: tc.DeliveryServiceRegexPost{
						Type:      totest.GetTypeId(t, TOSession, "HOST_REGEXP"),
						SetNumber: 3,
						Pattern:   "",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetDeliveryServiceRegexesByDSID(testCase.EndpointID(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.PostDeliveryServiceRegexesByDSID(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
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

func getDSRegexID(t *testing.T, dsName string) int {
	resp, _, err := TOSession.GetDeliveryServiceRegexesByDSID(totest.GetDeliveryServiceId(t, TOSession, dsName)(), client.RequestOptions{})
	assert.RequireNoError(t, err, "Get Delivery Service Regex failed with error: %v", err)
	assert.RequireGreaterOrEqual(t, len(resp.Response), 1, "Expected delivery service regex response object length 1, but got %d", len(resp.Response))
	assert.RequireNotNil(t, resp.Response[0].ID, "Expected id to not be nil")

	return resp.Response[0].ID
}
