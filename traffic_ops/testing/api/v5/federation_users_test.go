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

func TestFederationUsers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, ServerCapabilities, ServerServerCapabilities, DeliveryServices, CDNFederations, FederationUsers}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.FederationUserPost]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					EndpointID:    GetFederationID(t, "the.cname.com."),
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
					EndpointID:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"userID"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationUserIDSort(false)),
				},
				"VALID when SORTORDER param is DESC": {
					EndpointID:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"userID"}, "sortOrder": {"desc"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationUserIDSort(true)),
				},
				"FIRST RESULT when LIMIT=1": {
					EndpointID:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"userID"}, "limit": {"1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateFederationUsersPagination(GetFederationID(t, "the.cname.com.")(), "limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					EndpointID:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"userID"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateFederationUsersPagination(GetFederationID(t, "the.cname.com.")(), "offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					EndpointID:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"userID"}, "limit": {"1"}, "page": {"2"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateFederationUsersPagination(GetFederationID(t, "the.cname.com.")(), "page")),
				},
				"BAD REQUEST when INVALID LIMIT parameter": {
					EndpointID:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"-2"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID OFFSET parameter": {
					EndpointID:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "offset": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PAGE parameter": {
					EndpointID:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "page": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when CHANGES made": {
					EndpointID:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"OK when VALID request": {
					EndpointID:    GetFederationID(t, "google.com."),
					ClientSession: TOSession,
					RequestBody: tc.FederationUserPost{
						IDs: []int{
							GetUserID(t, "readonlyuser")(),
							GetUserID(t, "disalloweduser")(),
						},
						Replace: util.Ptr(false),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when REPLACING USERS": {
					EndpointID:    GetFederationID(t, "booya.com."),
					ClientSession: TOSession,
					RequestBody: tc.FederationUserPost{
						IDs:     []int{GetUserID(t, "readonlyuser")()},
						Replace: util.Ptr(true),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when ADDING USER": {
					EndpointID:    GetFederationID(t, "booya.com."),
					ClientSession: TOSession,
					RequestBody: tc.FederationUserPost{
						IDs:     []int{GetUserID(t, "disalloweduser")()},
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
					EndpointID:    GetFederationID(t, "the.cname.com."),
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

func CreateTestFederationUsers(t *testing.T) {
	// Prerequisite Federation Users
	federationUsers := map[string]tc.FederationUserPost{
		"the.cname.com.": {
			IDs:     []int{GetUserID(t, "admin")(), GetUserID(t, "adminuser")(), GetUserID(t, "disalloweduser")(), GetUserID(t, "readonlyuser")()},
			Replace: util.BoolPtr(false),
		},
		"booya.com.": {
			IDs:     []int{GetUserID(t, "adminuser")()},
			Replace: util.BoolPtr(false),
		},
	}

	for cname, federationUser := range federationUsers {
		fedID := GetFederationID(t, cname)()
		resp, _, err := TOSession.CreateFederationUsers(fedID, federationUser.IDs, *federationUser.Replace, client.RequestOptions{})
		assert.RequireNoError(t, err, "Assigning users %v to federation %d: %v - alerts: %+v", federationUser.IDs, fedID, err, resp.Alerts)
	}
}

func DeleteTestFederationUsers(t *testing.T) {
	for _, fedID := range fedIDs {
		fedUsers, _, err := TOSession.GetFederationUsers(fedID, client.RequestOptions{})
		assert.RequireNoError(t, err, "Error getting users for federation %d: %v - alerts: %+v", fedID, err, fedUsers.Alerts)
		for _, fedUser := range fedUsers.Response {
			if fedUser.ID == nil {
				t.Error("Traffic Ops returned a representation of a relationship between a user and a Federation that had null or undefined ID")
				continue
			}
			alerts, _, err := TOSession.DeleteFederationUser(fedID, *fedUser.ID, client.RequestOptions{})
			assert.NoError(t, err, "Error deleting user #%d from federation #%d: %v - alerts: %+v", *fedUser.ID, fedID, err, alerts.Alerts)
		}
		fedUsers, _, err = TOSession.GetFederationUsers(fedID, client.RequestOptions{})
		assert.NoError(t, err, "Error getting users for federation %d: %v - alerts: %+v", fedID, err, fedUsers.Alerts)
		assert.Equal(t, 0, len(fedUsers.Response), "Federation users expected 0, actual: %+v", len(fedUsers.Response))
	}
}
