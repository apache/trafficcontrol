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

func TestFederationUsers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations, FederationUsers}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.FederationUserPost]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"BAD REQUEST when INVALID FEDERATION ID": {
					EndpointID:    func() int { return -1 },
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"SORTED by ID when ORDERBY=USERID parameter": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"userID"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationUserIDSort(false)),
				},
				"VALID when SORTORDER param is DESC": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"userID"}, "sortOrder": {"desc"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationUserIDSort(true)),
				},
				"FIRST RESULT when LIMIT=1": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"userID"}, "limit": {"1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateFederationUsersPagination(totest.GetFederationID(t, "the.cname.com.")(), "limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"userID"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateFederationUsersPagination(totest.GetFederationID(t, "the.cname.com.")(), "offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"userID"}, "limit": {"1"}, "page": {"2"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateFederationUsersPagination(totest.GetFederationID(t, "the.cname.com.")(), "page")),
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
				"OK when CHANGES made": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"OK when VALID request": {
					EndpointID:    totest.GetFederationID(t, "google.com."),
					ClientSession: TOSession,
					RequestBody: tc.FederationUserPost{
						IDs: []int{
							totest.GetUserID(t, TOSession, "readonlyuser")(),
							totest.GetUserID(t, TOSession, "disalloweduser")(),
						},
						Replace: util.Ptr(false),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when REPLACING USERS": {
					EndpointID:    totest.GetFederationID(t, "booya.com."),
					ClientSession: TOSession,
					RequestBody: tc.FederationUserPost{
						IDs:     []int{totest.GetUserID(t, TOSession, "readonlyuser")()},
						Replace: util.Ptr(true),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when ADDING USER": {
					EndpointID:    totest.GetFederationID(t, "booya.com."),
					ClientSession: TOSession,
					RequestBody: tc.FederationUserPost{
						IDs:     []int{totest.GetUserID(t, TOSession, "disalloweduser")()},
						Replace: util.Ptr(false),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when INVALID FEDERATION ID": {
					EndpointID:    func() int { return -1 },
					ClientSession: TOSession,
					RequestBody: tc.FederationUserPost{
						IDs:     []int{},
						Replace: util.Ptr(false),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when INVALID USER ID": {
					EndpointID:    totest.GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestBody: tc.FederationUserPost{
						IDs:     []int{-1},
						Replace: util.Ptr(false),
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
							resp, reqInf, err := testCase.ClientSession.GetFederationUsers(testCase.EndpointID(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							federationUser := testCase.RequestBody
							alerts, reqInf, err := testCase.ClientSession.CreateFederationUsers(testCase.EndpointID(), federationUser.IDs, *federationUser.Replace, testCase.RequestOpts)
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

func validateFederationUsersPagination(federationID int, paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.FederationUser)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "userID")
		respBase, _, err := TOSession.GetFederationUsers(federationID, opts)
		assert.RequireNoError(t, err, "Cannot get Federation Users: %v - alerts: %+v", err, respBase.Alerts)

		federationUsers := respBase.Response
		assert.RequireGreaterOrEqual(t, len(federationUsers), 3, "Need at least 3 Federation Users in Traffic Ops to test pagination support, found: %d", len(federationUsers))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, federationUsers[:1], paginationResp, "expected GET Federation Users with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, federationUsers[1:2], paginationResp, "expected GET Federation Users with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, federationUsers[1:2], paginationResp, "expected GET Federation Users with limit = 1, page = 2 to return second result")
		}
	}
}

func validateFederationUserIDSort(desc bool) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Federation User response to not be nil.")
		var federationUserIDs []int
		federationUserResp := resp.([]tc.FederationUser)
		for _, fedUser := range federationUserResp {
			if desc {
				federationUserIDs = append([]int{*fedUser.ID}, federationUserIDs...)
			} else {
				federationUserIDs = append(federationUserIDs, *fedUser.ID)
			}
		}
		assert.Equal(t, true, sort.IntsAreSorted(federationUserIDs), "List is not sorted by their ids: %v", federationUserIDs)
	}
}
