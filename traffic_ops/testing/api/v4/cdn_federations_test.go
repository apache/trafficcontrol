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
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	totest "github.com/apache/trafficcontrol/v8/lib/go-tc/totestv4"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestCDNFederations(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Tenants, CacheGroups, Statuses, Divisions, Regions, PhysLocations, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.CDNFederation]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID ID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"id": {strconv.Itoa(totest.GetFederationID(t, "the.cname.com.")())}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1)),
				},
				"SORTED by CNAME when ORDERBY=CNAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"cname"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCDNFederationCNameSort()),
				},
				"SORTED when ORDERBY=ID and SORTORDER=DESC": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "sortOrder": {"desc"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCDNFederationIDDescSort()),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCDNFederationPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCDNFederationPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCDNFederationPagination("page")),
				},
				"BAD REQUEST when INVALID LIMIT parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"-2"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID OFFSET parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "offset": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PAGE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "page": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    totest.GetFederationID(t, "google.com."),
					ClientSession: TOSession,
					RequestBody: tc.CDNFederation{
						CName:       util.Ptr("new.cname."),
						TTL:         util.Ptr(34),
						Description: util.Ptr("updated"),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCDNFederationUpdateFields(map[string]interface{}{"CName": "new.cname."})),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:    totest.GetFederationID(t, "booya.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: tc.CDNFederation{
						CName:       util.Ptr("booya.com."),
						TTL:         util.Ptr(34),
						Description: util.Ptr("fooya"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    totest.GetFederationID(t, "booya.com."),
					ClientSession: TOSession,
					RequestBody: tc.CDNFederation{
						CName:       util.Ptr("new.cname."),
						TTL:         util.Ptr(34),
						Description: util.Ptr("updated"),
					},
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetCDNFederationsByName(totest.FederationCDNName, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateCDNFederation(testCase.RequestBody, totest.FederationCDNName, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateCDNFederation(testCase.RequestBody, totest.FederationCDNName, testCase.EndpointID(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteCDNFederation(totest.FederationCDNName, testCase.EndpointID(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts.Alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func validateCDNFederationUpdateFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected CDN Federation response to not be nil.")
		CDNFederationResp := resp.(tc.CDNFederation)
		for field, expected := range expectedResp {
			switch field {
			case "CName":
				assert.RequireNotNil(t, CDNFederationResp.CName, "Expected CName to not be nil.")
				assert.Equal(t, expected, *CDNFederationResp.CName, "Expected CName to be %v, but got %s", expected, *CDNFederationResp.CName)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}

func validateCDNFederationPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.CDNFederation)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetCDNFederationsByName(totest.FederationCDNName, opts)
		assert.RequireNoError(t, err, "Cannot get Federation Users: %v - alerts: %+v", err, respBase.Alerts)

		CDNfederations := respBase.Response
		assert.RequireGreaterOrEqual(t, len(CDNfederations), 3, "Need at least 3 CDN Federations in Traffic Ops to test pagination support, found: %d", len(CDNfederations))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, CDNfederations[:1], paginationResp, "expected GET CDN Federations with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, CDNfederations[1:2], paginationResp, "expected GET CDN Federations with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, CDNfederations[1:2], paginationResp, "expected GET CDN Federations with limit = 1, page = 2 to return second result")
		}
	}
}

func validateCDNFederationCNameSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected CDN Federation response to not be nil.")
		var federationCNames []string
		CDNFederationResp := resp.([]tc.CDNFederation)
		for _, CDNFederation := range CDNFederationResp {
			assert.RequireNotNil(t, CDNFederation.CName, "Expected CDN Federation CName to not be nil.")
			federationCNames = append(federationCNames, *CDNFederation.CName)
		}
		assert.Equal(t, true, sort.StringsAreSorted(federationCNames), "List is not sorted by their names: %v", federationCNames)
	}
}

func validateCDNFederationIDDescSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected CDN Federation response to not be nil.")
		var CDNFederationIDs []int
		CDNFederationResp := resp.([]tc.CDNFederation)
		for _, federation := range CDNFederationResp {
			CDNFederationIDs = append([]int{*federation.ID}, CDNFederationIDs...)
		}
		assert.Equal(t, true, sort.IntsAreSorted(CDNFederationIDs), "List is not sorted by their ids: %v", CDNFederationIDs)
	}
}
