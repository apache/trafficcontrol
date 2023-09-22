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
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestFederationFederationResolvers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Tenants, CacheGroups, Statuses, Divisions, Regions, PhysLocations, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations, FederationResolvers, FederationFederationResolvers}, func() {

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.AssignFederationResolversRequest]{
			"GET": {
				"OK when VALID request AND RESOLVERS ASSIGNED": {
					EndpointID:    totest.GetFederationID(t, "booya.com."),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID request AND NO RESOLVERS ASSIGNED": {
					EndpointID:    totest.GetFederationID(t, "google.com."),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
			},
			"POST": {
				"OK when ASSIGNING ONE FEDERATION RESOLVER": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestBody: tc.AssignFederationResolversRequest{
						FedResolverIDs: []int{GetFederationResolverID(t, "1.2.3.4")()},
						Replace:        false,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when ASSIGNING MULTIPLE FEDERATION RESOLVERS": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestBody: tc.AssignFederationResolversRequest{
						FedResolverIDs: []int{
							GetFederationResolverID(t, "1.2.3.4")(),
							GetFederationResolverID(t, "0.0.0.0/12")(),
							GetFederationResolverID(t, "::f1d0:f00d/123")(),
						},
						Replace: false,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when REPLACING ALL FEDERATION RESOLVERS": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestBody: tc.AssignFederationResolversRequest{
						FedResolverIDs: []int{GetFederationResolverID(t, "dead::babe")()},
						Replace:        true,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when FEDERATION DOESNT EXIST": {
					EndpointID:    func() int { return -1 },
					ClientSession: TOSession,
					RequestBody: tc.AssignFederationResolversRequest{
						FedResolverIDs: []int{GetFederationResolverID(t, "1.2.3.4")()},
						Replace:        false,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetFederationFederationResolvers(testCase.EndpointID(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							frAssignment := testCase.RequestBody
							resp, reqInf, err := testCase.ClientSession.AssignFederationFederationResolver(testCase.EndpointID(), frAssignment.FedResolverIDs, frAssignment.Replace, testCase.RequestOpts)
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
		fedID := totest.GetFederationID(t, cname)()
		resp, _, err := TOSession.AssignFederationFederationResolver(fedID, federationFederationResolver.FedResolverIDs, federationFederationResolver.Replace, client.RequestOptions{})
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
		fedID := totest.GetFederationID(t, cname)()
		resp, _, err := TOSession.AssignFederationFederationResolver(fedID, federationFederationResolver.FedResolverIDs, federationFederationResolver.Replace, client.RequestOptions{})
		assert.RequireNoError(t, err, "Assigning resolvers %v to federation %d: %v - alerts: %+v", federationFederationResolver.FedResolverIDs, fedID, err, resp.Alerts)
	}
}
