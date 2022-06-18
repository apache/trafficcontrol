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
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestFederationsDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations, FederationDeliveryServices}, func() {

		methodTests := utils.V4TestCase{
			"GET": {
				"OK when VALID request": {
					EndpointId:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"SORTED when ORDERBY=DSID parameter": {
					EndpointId:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"dsID"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationDeliveryServicesSort(false)),
				},
				"SORTED when ORDERBY=DSID and SORTORDER=DESC parameter": {
					EndpointId:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"dsID"}, "sortOrder": {"desc"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationDeliveryServicesSort(true)),
				},
				"FIRST RESULT when LIMIT=1": {
					EndpointId:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"dsID"}, "limit": {"1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateFederationDeliveryServicesPagination(GetFederationID(t, "the.cname.com.")(), "limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					EndpointId:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"dsID"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateFederationDeliveryServicesPagination(GetFederationID(t, "the.cname.com.")(), "offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					EndpointId:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"dsID"}, "limit": {"1"}, "page": {"2"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateFederationDeliveryServicesPagination(GetFederationID(t, "the.cname.com.")(), "page")),
				},
				"BAD REQUEST when INVALID LIMIT parameter": {
					EndpointId:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"-2"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID OFFSET parameter": {
					EndpointId:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "offset": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PAGE parameter": {
					EndpointId:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "page": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"DELETE": {
				"OK when VALID request": {
					EndpointId:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"dsID": {strconv.Itoa(GetDeliveryServiceId(t, "ds1")())}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when LAST DELIVERY SERVICE": {
					EndpointId:    GetFederationID(t, "google.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"dsID": {strconv.Itoa(GetDeliveryServiceId(t, "ds2")())}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var dsID int
					fedDS := tc.FederationDSPost{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &fedDS)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetFederationDeliveryServices(testCase.EndpointId(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateFederationDeliveryServices(testCase.EndpointId(), fedDS.DSIDs, *fedDS.Replace, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							if val, ok := testCase.RequestOpts.QueryParameters["dsID"]; ok {
								id, err := strconv.Atoi(val[0])
								assert.RequireNoError(t, err, "Failed to convert dsID to an integer.")
								dsID = id
							}
							alerts, reqInf, err := testCase.ClientSession.DeleteFederationDeliveryService(testCase.EndpointId(), dsID, testCase.RequestOpts)
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
			DSIDs:   []int{GetDeliveryServiceId(t, "ds1")(), GetDeliveryServiceId(t, "ds2")(), GetDeliveryServiceId(t, "ds3")(), GetDeliveryServiceId(t, "ds4")()},
			Replace: util.BoolPtr(true),
		},
		"google.com.": {
			DSIDs:   []int{GetDeliveryServiceId(t, "ds1")()},
			Replace: util.BoolPtr(true),
		},
	}

	for federation, fedDS := range federationDS {
		alerts, _, err := TOSession.CreateFederationDeliveryServices(GetFederationID(t, federation)(), fedDS.DSIDs, *fedDS.Replace, client.RequestOptions{})
		assert.RequireNoError(t, err, "Creating federations delivery services: %v - alerts: %+v", err, alerts.Alerts)
	}
}

func DeleteTestFederationDeliveryServices(t *testing.T) {
	// Prerequisite Federation Delivery Services
	federation := []string{"the.cname.com.", "google.com."}
	for _, cname := range federation {
		resp, _, err := TOSession.GetFederationDeliveryServices(GetFederationID(t, cname)(), client.RequestOptions{})
		assert.RequireNoError(t, err, "Error when getting Federation Delivery Services.")
		for _, fedDS := range resp.Response {
			assert.RequireNotNil(t, fedDS.ID, "Expected Federation Delivery Service ID to not be nil.")
			alerts, _, err := TOSession.DeleteFederationDeliveryService(GetFederationID(t, cname)(), *fedDS.ID, client.RequestOptions{})
			assert.NoError(t, err, "Unexpected error deleting Federation Delivery Service: %v - alerts: %+v", err, alerts.Alerts)
		}
		resp, _, err = TOSession.GetFederationDeliveryServices(GetFederationID(t, cname)(), client.RequestOptions{})
		assert.RequireNoError(t, err, "Error when getting Federation Delivery Services.")
		assert.Equal(t, 0, len(resp.Response), "Expected Federation Delivery Services length to be 0. Got: %d", len(resp.Response))
	}
}
