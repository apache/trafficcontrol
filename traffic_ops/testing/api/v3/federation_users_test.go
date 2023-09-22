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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
)

func TestFederationUsers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations, FederationUsers}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.FederationUserPost]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					EndpointID:     GetFederationID(t, "the.cname.com."),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"BAD REQUEST when INVALID FEDERATION ID": {
					EndpointID:    func() int { return -1 },
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
			"POST": {
				"OK when VALID request": {
					EndpointID:    GetFederationID(t, "google.com."),
					ClientSession: TOSession,
					RequestBody: tc.FederationUserPost{
						IDs:     []int{GetUserID(t, "readonlyuser")(), GetUserID(t, "disalloweduser")()},
						Replace: util.Ptr(false),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when REPLACING USERS": {
					EndpointID:    GetFederationID(t, "the.cname.com."),
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
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					EndpointID:     GetFederationID(t, "the.cname.com."),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetFederationUsersWithHdr(testCase.EndpointID(), testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateFederationUsers(testCase.EndpointID(), testCase.RequestBody.IDs, *testCase.RequestBody.Replace)
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
		resp, _, err := TOSession.CreateFederationUsers(fedID, federationUser.IDs, *federationUser.Replace)
		assert.RequireNoError(t, err, "Assigning users %v to federation %d: %v - alerts: %+v", federationUser.IDs, fedID, err, resp.Alerts)
	}
}

func DeleteTestFederationUsers(t *testing.T) {
	for _, fedID := range fedIDs {
		fedUsers, _, err := TOSession.GetFederationUsersWithHdr(fedID, nil)
		assert.RequireNoError(t, err, "Error getting users for federation %d: %v", fedID, err)
		for _, fedUser := range fedUsers {
			if fedUser.ID == nil {
				t.Error("Traffic Ops returned a representation of a relationship between a user and a Federation that had null or undefined ID")
				continue
			}
			alerts, _, err := TOSession.DeleteFederationUser(fedID, *fedUser.ID)
			assert.NoError(t, err, "Error deleting user #%d from federation #%d: %v - alerts: %+v", *fedUser.ID, fedID, err, alerts.Alerts)
		}
	}
	for _, fedID := range fedIDs {
		fedUsers, _, err := TOSession.GetFederationUsersWithHdr(fedID, nil)
		assert.NoError(t, err, "Error getting users for federation %d: %v", fedID, err)
		assert.Equal(t, 0, len(fedUsers), "Federation users expected 0, actual: %+v", len(fedUsers))
	}
}
