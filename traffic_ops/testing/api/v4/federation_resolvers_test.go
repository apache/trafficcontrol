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

func TestFederationResolvers(t *testing.T) {
	WithObjs(t, []TCObj{Types, FederationResolvers}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.FederationResolver]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID ID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"id": {strconv.Itoa(GetFederationResolverID(t, "0.0.0.0/12")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateFederationResolverFields(map[string]interface{}{"ID": uint(GetFederationResolverID(t, "0.0.0.0/12")())})),
				},
				"OK when VALID IPADDRESS parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"ipAddress": {"1.2.3.4"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateFederationResolverFields(map[string]interface{}{"IPAddress": "1.2.3.4"})),
				},
				"OK when VALID TYPE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"type": {"RESOLVE4"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateFederationResolverFields(map[string]interface{}{"Type": "RESOLVE4"})),
				},
				"SORTED by ID when ORDERBY=ID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationResolverIDSort()),
				},
				"VALID when SORTORDER param is DESC": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "sortOrder": {"desc"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationResolverDescSort()),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationResolverPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationResolverPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationResolverPagination("page")),
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
			},
			"POST": {
				"BAD REQUEST when MISSING IPADDRESS and TYPE FIELDS": {
					ClientSession: TOSession,
					RequestBody:   tc.FederationResolver{},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID IP ADDRESS": {
					ClientSession: TOSession,
					RequestBody: tc.FederationResolver{
						IPAddress: util.Ptr("not a valid IP address"),
						TypeID:    util.Ptr((uint)(totest.GetTypeId(t, TOSession, "RESOLVE4"))),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"DELETE": {
				"NOT FOUND when INVALID ID": {
					EndpointID:    func() int { return 0 },
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetFederationResolvers(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateFederationResolver(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.DeleteFederationResolver(uint(testCase.EndpointID()), testCase.RequestOpts)
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

func validateFederationResolverFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Federation Resolver response to not be nil.")
		frResp := resp.([]tc.FederationResolver)
		for field, expected := range expectedResp {
			for _, fr := range frResp {
				switch field {
				case "ID":
					assert.RequireNotNil(t, fr.ID, "Expected ID to not be nil")
					assert.Equal(t, expected, *fr.ID, "Expected ID to be %v, but got %d", expected, *fr.ID)
				case "IPAddress":
					assert.RequireNotNil(t, fr.IPAddress, "Expected IPAddress to not be nil")
					assert.Equal(t, expected, *fr.IPAddress, "Expected IPAddress to be %v, but got %s", expected, *fr.IPAddress)
				case "Type":
					assert.RequireNotNil(t, fr.Type, "Expected Type to not be nil")
					assert.Equal(t, expected, *fr.Type, "Expected Type to be %v, but got %s", expected, *fr.Type)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateFederationResolverPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.FederationResolver)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetFederationResolvers(opts)
		assert.RequireNoError(t, err, "Cannot get Federation Resolvers: %v - alerts: %+v", err, respBase.Alerts)

		fr := respBase.Response
		assert.RequireGreaterOrEqual(t, len(fr), 3, "Need at least 3 Federation Resolvers in Traffic Ops to test pagination support, found: %d", len(fr))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, fr[:1], paginationResp, "expected GET Federation Resolvers with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, fr[1:2], paginationResp, "expected GET Federation Resolvers with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, fr[1:2], paginationResp, "expected GET Federation Resolvers with limit = 1, page = 2 to return second result")
		}
	}
}

func validateFederationResolverIDSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Federation Resolver response to not be nil.")
		var frIDs []int
		frResp := resp.([]tc.FederationResolver)
		for _, fr := range frResp {
			frIDs = append(frIDs, int(*fr.ID))
		}
		assert.Equal(t, true, sort.IntsAreSorted(frIDs), "List is not sorted by their ids: %v", frIDs)
	}
}

func validateFederationResolverDescSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected FederationResolver response to not be nil.")
		frDescResp := resp.([]tc.FederationResolver)
		var descSortedList []uint
		var ascSortedList []uint
		assert.RequireGreaterOrEqual(t, len(frDescResp), 2, "Need at least 2 Federation Resolvers in Traffic Ops to test desc sort, found: %d", len(frDescResp))
		// Get Federation Resolvers in the default ascending order for comparison.
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		frAscResp, _, err := TOSession.GetFederationResolvers(opts)
		assert.RequireNoError(t, err, "Unexpected error getting Federation Resolvers with default sort order: %v - alerts: %+v", err, frAscResp.Alerts)
		// Verify the response match in length, i.e. equal amount of Federation Resolvers.
		assert.RequireEqual(t, len(frAscResp.Response), len(frDescResp), "Expected descending order response length: %v, to match ascending order response length %v", len(frAscResp.Response), len(frDescResp))
		// Insert Federation Resolvers names to the front of a new list, so they are now reversed to be in ascending order.
		for _, fr := range frDescResp {
			descSortedList = append([]uint{*fr.ID}, descSortedList...)
		}
		// Insert Federation Resolvers IDs by appending to a new list, so they stay in ascending order.
		for _, fr := range frAscResp.Response {
			ascSortedList = append(ascSortedList, *fr.ID)
		}
		assert.Exactly(t, ascSortedList, descSortedList, "Federation Resolver responses are not equal after reversal: %v - %v", ascSortedList, descSortedList)
	}
}

func GetFederationResolverID(t *testing.T, ipAddress string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("ipAddress", ipAddress)
		federationResolvers, _, err := TOSession.GetFederationResolvers(opts)
		assert.RequireNoError(t, err, "Get FederationResolvers Request failed with error:", err)
		assert.RequireEqual(t, 1, len(federationResolvers.Response), "Expected response object length 1, but got %d", len(federationResolvers.Response))
		assert.RequireNotNil(t, federationResolvers.Response[0].ID, "Expected Federation Resolver ID to not be nil")
		return int(*federationResolvers.Response[0].ID)
	}
}
