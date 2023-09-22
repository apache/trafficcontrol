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
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

var fedIDs = make(map[string]int)

// All prerequisite Federations are associated to this cdn and this xmlID
var cdnName = "cdn1"
var xmlId = "ds1"

func TestCDNFederations(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Tenants, CacheGroups, Statuses, Divisions, Regions, PhysLocations, Servers, Topologies, ServiceCategories, ServerCapabilities, DeliveryServices, CDNFederations}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.CDNFederationV5]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID ID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"id": {strconv.Itoa(GetFederationID(t, "the.cname.com.")())}}},
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
					EndpointID:    GetFederationID(t, "google.com."),
					ClientSession: TOSession,
					RequestBody: tc.CDNFederationV5{
						CName:       "new.cname.",
						TTL:         64,
						Description: util.Ptr("updated"),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCDNFederationUpdateFields(map[string]interface{}{"CName": "new.cname."})),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:    GetFederationID(t, "booya.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: tc.CDNFederationV5{
						CName:       "booya.com.",
						TTL:         64,
						Description: util.Ptr("fooya"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetFederationID(t, "booya.com."),
					ClientSession: TOSession,
					RequestBody: tc.CDNFederationV5{
						CName:       "new.cname.",
						TTL:         64,
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
							resp, reqInf, err := testCase.ClientSession.GetCDNFederations(cdnName, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateCDNFederation(testCase.RequestBody, cdnName, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateCDNFederation(testCase.RequestBody, cdnName, testCase.EndpointID(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteCDNFederation(cdnName, testCase.EndpointID(), testCase.RequestOpts)
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
		CDNFederationResp := resp.(tc.CDNFederationV5)
		for field, expected := range expectedResp {
			switch field {
			case "CName":
				assert.Equal(t, expected, CDNFederationResp.CName, "Expected CName to be %v, but got %s", expected, CDNFederationResp.CName)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}

func validateCDNFederationPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.CDNFederationV5)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetCDNFederations(cdnName, opts)
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
		CDNFederationResp := resp.([]tc.CDNFederationV5)
		for _, CDNFederationV5 := range CDNFederationResp {
			federationCNames = append(federationCNames, CDNFederationV5.CName)
		}
		assert.Equal(t, true, sort.StringsAreSorted(federationCNames), "List is not sorted by their names: %v", federationCNames)
	}
}

func validateCDNFederationIDDescSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected CDN Federation response to not be nil.")
		var CDNFederationIDs []int
		CDNFederationResp := resp.([]tc.CDNFederationV5)
		for _, federation := range CDNFederationResp {
			CDNFederationIDs = append([]int{federation.ID}, CDNFederationIDs...)
		}
		assert.Equal(t, true, sort.IntsAreSorted(CDNFederationIDs), "List is not sorted by their ids: %v", CDNFederationIDs)
	}
}

func GetFederationID(t *testing.T, cname string) func() int {
	return func() int {
		ID, ok := fedIDs[cname]
		assert.RequireEqual(t, true, ok, "Expected to find Federation CName: %s to have associated ID", cname)
		return ID
	}
}

func setFederationID(t *testing.T, cdnFederation tc.CDNFederationV5) {
	fedIDs[cdnFederation.CName] = cdnFederation.ID
}

func CreateTestCDNFederations(t *testing.T) {
	for _, federation := range testData.Federations {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", federation.DeliveryService.XMLID)
		dsResp, _, err := TOSession.GetDeliveryServices(opts)
		assert.RequireNoError(t, err, "Could not get Delivery Service by XML ID: %v", err)
		assert.RequireEqual(t, 1, len(dsResp.Response), "Expected one Delivery Service, but got %d", len(dsResp.Response))
		assert.RequireNotNil(t, dsResp.Response[0].CDNName, "Expected Delivery Service CDN Name to not be nil.")

		resp, _, err := TOSession.CreateCDNFederation(federation, *dsResp.Response[0].CDNName, client.RequestOptions{})
		assert.NoError(t, err, "Could not create CDN Federations: %v - alerts: %+v", err, resp.Alerts)

		// Need to save the ids, otherwise the other tests won't be able to reference the federations
		setFederationID(t, resp.Response)
		assert.RequireNotNil(t, resp.Response.ID, "Federation ID was nil after posting.")
		assert.RequireNotNil(t, dsResp.Response[0].ID, "Delivery Service ID was nil.")
		_, _, err = TOSession.CreateFederationDeliveryServices(resp.Response.ID, []int{*dsResp.Response[0].ID}, false, client.NewRequestOptions())
		assert.NoError(t, err, "Could not create Federation Delivery Service: %v", err)
	}
}

func DeleteTestCDNFederations(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, id := range fedIDs {
		resp, _, err := TOSession.DeleteCDNFederation(cdnName, id, opts)
		assert.NoError(t, err, "Cannot delete federation #%d: %v - alerts: %+v", id, err, resp.Alerts)

		opts.QueryParameters.Set("id", strconv.Itoa(id))
		data, _, err := TOSession.GetCDNFederations(cdnName, opts)
		assert.Equal(t, 0, len(data.Response), "expected federation to be deleted")
	}
	fedIDs = make(map[string]int) // reset the global variable for the next test
}
