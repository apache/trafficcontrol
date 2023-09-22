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
	"net/url"
	"sort"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	totest "github.com/apache/trafficcontrol/v8/lib/go-tc/totestv4"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestFederationsDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations, FederationDeliveryServices}, func() {

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.FederationDSPost]{
			"GET": {
				"OK when VALID request": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"SORTED when ORDERBY=DSID parameter": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"dsID"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationDeliveryServicesSort(false)),
				},
				"SORTED when ORDERBY=DSID and SORTORDER=DESC parameter": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"dsID"}, "sortOrder": {"desc"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationDeliveryServicesSort(true)),
				},
				"FIRST RESULT when LIMIT=1": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"dsID"}, "limit": {"1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateFederationDeliveryServicesPagination(totest.GetFederationID(t, "the.cname.com.")(), "limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"dsID"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateFederationDeliveryServicesPagination(totest.GetFederationID(t, "the.cname.com.")(), "offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"dsID"}, "limit": {"1"}, "page": {"2"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateFederationDeliveryServicesPagination(totest.GetFederationID(t, "the.cname.com.")(), "page")),
				},
				"BAD REQUEST when INVALID LIMIT parameter": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"-2"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID OFFSET parameter": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "offset": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PAGE parameter": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "page": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"DELETE": {
				"OK when VALID request": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"dsID": {strconv.Itoa(totest.GetDeliveryServiceId(t, TOSession, "ds1")())}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when LAST DELIVERY SERVICE": {
					EndpointID:    totest.GetFederationID(t, "google.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"dsID": {strconv.Itoa(totest.GetDeliveryServiceId(t, TOSession, "ds2")())}}},
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
							resp, reqInf, err := testCase.ClientSession.GetFederationDeliveryServices(testCase.EndpointID(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							fedDS := testCase.RequestBody
							alerts, reqInf, err := testCase.ClientSession.CreateFederationDeliveryServices(testCase.EndpointID(), fedDS.DSIDs, *fedDS.Replace, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							var dsID int
							if val, ok := testCase.RequestOpts.QueryParameters["dsID"]; ok {
								id, err := strconv.Atoi(val[0])
								assert.RequireNoError(t, err, "Failed to convert dsID to an integer.")
								dsID = id
							}
							alerts, reqInf, err := testCase.ClientSession.DeleteFederationDeliveryService(testCase.EndpointID(), dsID, testCase.RequestOpts)
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

func validateFederationDeliveryServicesSort(desc bool) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Federation DeliveryServices response to not be nil.")
		var federationDSIDs []int
		federationDSResp := resp.([]tc.FederationDeliveryServiceNullable)
		for _, federationDS := range federationDSResp {
			if desc {
				federationDSIDs = append([]int{*federationDS.ID}, federationDSIDs...)
			} else {
				federationDSIDs = append(federationDSIDs, *federationDS.ID)
			}
		}
		assert.Equal(t, true, sort.IntsAreSorted(federationDSIDs), "List is not sorted by their ids: %v", federationDSIDs)
	}
}

func validateFederationDeliveryServicesPagination(fedID int, paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Federation DeliveryServices response to not be nil.")
		paginationResp := resp.([]tc.FederationDeliveryServiceNullable)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "dsID")
		respBase, _, err := TOSession.GetFederationDeliveryServices(fedID, opts)
		assert.RequireNoError(t, err, "Cannot get Federation DeliveryServices: %v - alerts: %+v", err, respBase.Alerts)

		federationDS := respBase.Response
		assert.RequireGreaterOrEqual(t, len(federationDS), 3, "Need at least 3 Federation DeliveryServices in Traffic Ops to test pagination support, found: %d", len(federationDS))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, federationDS[:1], paginationResp, "expected GET Federation DeliveryServices with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, federationDS[1:2], paginationResp, "expected GET Federation DeliveryServices with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, federationDS[1:2], paginationResp, "expected GET Federation DeliveryServices with limit = 1, page = 2 to return second result")
		}
	}
}

func CreateTestFederationDeliveryServices(t *testing.T) {
	// Prerequisite Federation Delivery Services
	federationDS := map[string]tc.FederationDSPost{
		"the.cname.com.": {
			DSIDs:   []int{totest.GetDeliveryServiceId(t, TOSession, "ds1")(), totest.GetDeliveryServiceId(t, TOSession, "ds2")(), totest.GetDeliveryServiceId(t, TOSession, "ds3")(), totest.GetDeliveryServiceId(t, TOSession, "ds4")()},
			Replace: util.BoolPtr(true),
		},
		"google.com.": {
			DSIDs:   []int{totest.GetDeliveryServiceId(t, TOSession, "ds1")()},
			Replace: util.BoolPtr(true),
		},
	}

	for federation, fedDS := range federationDS {
		alerts, _, err := TOSession.CreateFederationDeliveryServices(totest.GetFederationID(t, federation)(), fedDS.DSIDs, *fedDS.Replace, client.RequestOptions{})
		assert.RequireNoError(t, err, "Creating federations delivery services: %v - alerts: %+v", err, alerts.Alerts)
	}
}
