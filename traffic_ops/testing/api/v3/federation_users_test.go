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
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
)

func TestFederationUsers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations, FederationUsers}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					EndpointId:     GetFederationID(t, "the.cname.com."),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"BAD REQUEST when INVALID FEDERATION ID": {
					EndpointId:    func() int { return -1 },
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
			"POST": {
				"OK when VALID request": {
					EndpointId:    GetFederationID(t, "google.com."),
					ClientSession: TOSession,
					RequestBody:   map[string]interface{}{"userIds": []int{GetUserID(t, "readonlyuser")(), GetUserID(t, "disalloweduser")()}, "replace": false},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when REPLACING USERS": {
					EndpointId:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestBody:   map[string]interface{}{"userIds": []int{GetUserID(t, "readonlyuser")()}, "replace": true},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when ADDING USER": {
					EndpointId:    GetFederationID(t, "booya.com."),
					ClientSession: TOSession,
					RequestBody:   map[string]interface{}{"userIds": []int{GetUserID(t, "disalloweduser")()}, "replace": false},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when INVALID FEDERATION ID": {
					EndpointId:    func() int { return -1 },
					ClientSession: TOSession,
					RequestBody:   map[string]interface{}{"userIds": []int{}, "replace": false},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when INVALID USER ID": {
					EndpointId:    GetFederationID(t, "the.cname.com."),
					ClientSession: TOSession,
					RequestBody:   map[string]interface{}{"userIds": []int{-1}, "replace": false},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					EndpointId:     GetFederationID(t, "the.cname.com."),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					federationUser := tc.FederationUserPost{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &federationUser)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetFederationUsersWithHdr(testCase.EndpointId(), testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateFederationUsers(testCase.EndpointId(), federationUser.IDs, *federationUser.Replace)
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
