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
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
)

func TestFederationsDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations, FederationDeliveryServices}, func() {

		methodTests := utils.V3TestCaseT[tc.FederationDSPost]{
			"GET": {
				"OK when VALID request": {
					EndpointID:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
			},
			"DELETE": {
				"OK when VALID request": {
					EndpointID:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestParams: url.Values{"dsID": {strconv.Itoa(GetDeliveryServiceId(t, "ds1")())}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when LAST DELIVERY SERVICE": {
					EndpointID:    GetFederationID(t, "google.com."),
					ClientSession: TOSession,
					RequestParams: url.Values{"dsID": {strconv.Itoa(GetDeliveryServiceId(t, "ds2")())}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetFederationDeliveryServicesWithHdr(testCase.EndpointID(), testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							reqInf, err := testCase.ClientSession.CreateFederationDeliveryServices(testCase.EndpointID(), testCase.RequestBody.DSIDs, *testCase.RequestBody.Replace)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, tc.Alerts{}, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							var dsID int
							if val, ok := testCase.RequestParams["dsID"]; ok {
								id, err := strconv.Atoi(val[0])
								assert.RequireNoError(t, err, "Failed to convert dsID to an integer.")
								dsID = id
							}
							alerts, reqInf, err := testCase.ClientSession.DeleteFederationDeliveryService(testCase.EndpointID(), dsID)
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

func CreateTestFederationDeliveryServices(t *testing.T) {
	// Prerequisite Federation Delivery Services
	federationDS := map[string]tc.FederationDSPost{
		"the.cname.com.": {
			DSIDs:   []int{GetDeliveryServiceId(t, "ds1")(), GetDeliveryServiceId(t, "ds2")(), GetDeliveryServiceId(t, "ds3")(), GetDeliveryServiceId(t, "ds4")()},
			Replace: util.BoolPtr(true),
		},
		"google.com.": {
			DSIDs:   []int{GetDeliveryServiceId(t, "ds1")()},
			Replace: util.BoolPtr(true),
		},
	}

	for federation, fedDS := range federationDS {
		_, err := TOSession.CreateFederationDeliveryServices(GetFederationID(t, federation)(), fedDS.DSIDs, *fedDS.Replace)
		assert.RequireNoError(t, err, "Creating federations delivery services: %v", err)
	}
}
