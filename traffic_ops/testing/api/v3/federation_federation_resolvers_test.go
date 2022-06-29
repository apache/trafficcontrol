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
	"encoding/json"
	"net/http"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
)

func TestFederationFederationResolvers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Tenants, CacheGroups, Statuses, Divisions, Regions, PhysLocations, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations, FederationResolvers, FederationFederationResolvers}, func() {

		methodTests := utils.V3TestCase{
			"GET": {
				"OK when VALID request AND RESOLVERS ASSIGNED": {
					EndpointId:    GetFederationID(t, "booya.com."),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID request AND NO RESOLVERS ASSIGNED": {
					EndpointId:    GetFederationID(t, "google.com."),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
			},
			"POST": {
				"OK when ASSIGNING ONE FEDERATION RESOLVER": {
					EndpointId:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"resolverIDs": []int{GetFederationResolverID(t, "1.2.3.4")()},
						"replace":     false,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when ASSIGNING MULTIPLE FEDERATION RESOLVERS": {
					EndpointId:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"resolverIDs": []int{
							GetFederationResolverID(t, "1.2.3.4")(),
							GetFederationResolverID(t, "0.0.0.0/12")(),
							GetFederationResolverID(t, "::f1d0:f00d/123")(),
						},
						"replace": false,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when REPLACING ALL FEDERATION RESOLVERS": {
					EndpointId:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"resolverIDs": []int{GetFederationResolverID(t, "dead::babe")()},
						"replace":     true,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when FEDERATION DOESNT EXIST": {
					EndpointId:    func() int { return -1 },
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"resolverIDs": []int{GetFederationResolverID(t, "1.2.3.4")()},
						"replace":     false,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					frAssignment := tc.AssignFederationResolversRequest{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &frAssignment)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetFederationFederationResolversByID(testCase.EndpointId())
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.AssignFederationFederationResolver(testCase.EndpointId(), frAssignment.FedResolverIDs, frAssignment.Replace)
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

func CreateTestFederationFederationResolvers(t *testing.T) {
	// Prerequisite Federation Federation Resolvers
	federationFederationResolvers := map[string]tc.AssignFederationResolversRequest{
		"booya.com.": {
			FedResolverIDs: []int{
				GetFederationResolverID(t, "1.2.3.4")(),
				GetFederationResolverID(t, "0.0.0.0/12")(),
				GetFederationResolverID(t, "::f1d0:f00d/123")(),
				GetFederationResolverID(t, "dead::babe")(),
			},
			Replace: false,
		},
	}

	for cname, federationFederationResolver := range federationFederationResolvers {
		fedID := GetFederationID(t, cname)()
		resp, _, err := TOSession.AssignFederationFederationResolver(fedID, federationFederationResolver.FedResolverIDs, federationFederationResolver.Replace)
		assert.RequireNoError(t, err, "Assigning resolvers %v to federation %d: %v - alerts: %+v", federationFederationResolver.FedResolverIDs, fedID, err, resp.Alerts)
	}
}

func DeleteTestFederationFederationResolvers(t *testing.T) {
	// Prerequisite Federation Federation Resolvers
	federationFederationResolvers := map[string]tc.AssignFederationResolversRequest{
		"booya.com.": {
			FedResolverIDs: []int{},
			Replace:        true,
		},
	}

	for cname, federationFederationResolver := range federationFederationResolvers {
		fedID := GetFederationID(t, cname)()
		resp, _, err := TOSession.AssignFederationFederationResolver(fedID, federationFederationResolver.FedResolverIDs, federationFederationResolver.Replace)
		assert.RequireNoError(t, err, "Assigning resolvers %v to federation %d: %v - alerts: %+v", federationFederationResolver.FedResolverIDs, fedID, err, resp.Alerts)
	}
}
