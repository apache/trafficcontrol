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

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestFederationUsers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations, FederationUsers}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V4TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"BAD REQUEST when INVALID FEDERATION ID": {
					EndpointId:    func() int { return -1 },
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"SORTED by ID when ORDERBY=USERID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"userID"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationUserIDSort()),
				},
				"VALID when SORTORDER param is DESC": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"userID"}, "sortOrder": {"desc"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationUserIDSort()),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationUsersPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationUsersPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateFederationUsersPagination("page")),
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
				"OK when VALID request": {
					EndpointId:    GetFederationID(),
					ClientSession: TOSession,
					RequestBody:   map[string]interface{}{"userIds": []int{}, "replace": false},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when REPLACING USERS": {
					EndpointId:    GetFederationID(),
					ClientSession: TOSession,
					RequestBody:   map[string]interface{}{"userIds": []int{}, "replace": false},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when ADDING USER": {
					EndpointId:    GetFederationID(),
					ClientSession: TOSession,
					RequestBody:   map[string]interface{}{"userIds": []int{}, "replace": false},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when INVALID FEDERATION ID": {
					EndpointId:    func() int { return -1 },
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when INVALID USER ID": {
					EndpointId:    GetFederationID(),
					ClientSession: TOSession,
					RequestBody:   map[string]interface{}{"userIds": []int{-1}, "replace": false},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}
	})
}

func validateFederationUsersPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.FederationUser)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetFederationUsers(opts)
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

func validateFederationUserIDSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Federation User response to not be nil.")
		var federationUserIDs []int
		federationUserResp := resp.([]tc.FederationUser)
		for _, fedUser := range federationUserResp {
			federationUserIDs = append(federationUserIDs, *fedUser.ID)
		}
		assert.Equal(t, true, sort.IntsAreSorted(federationUserIDs), "List is not sorted by their ids: %v", federationUserIDs)
	}
}

func CreateTestFederationUsers(t *testing.T) {
	if len(testData.Federations) == 0 {
		t.Error("no federations test data")
	}

	if len(fedIDs) < 1 {
		t.Fatal("need at least one stored Federation ID to test Federations")
	}
	fedID := fedIDs[0]

	// Get Users
	users, _, err := TOSession.GetUsers(client.RequestOptions{})
	if err != nil {
		t.Fatalf("getting users: %v - alerts: %+v", err, users.Alerts)
	}
	if len(users.Response) < 3 {
		t.Fatal("need > 3 users to create federation users")
	}

	u1 := users.Response[0].ID
	u2 := users.Response[1].ID
	u3 := users.Response[2].ID
	if u1 == nil || u2 == nil || u3 == nil {
		t.Fatal("Traffic Ops returned at least one representation of a relationship between a user and a Federation that had a null or undefined ID")
	}

	// Associate one user to federation
	resp, _, err := TOSession.CreateFederationUsers(fedID, []int{*u1}, false, client.RequestOptions{})
	if err != nil {
		t.Fatalf("assigning users %v to federation %d: %v - alerts: %+v", []int{*u1}, fedID, err, resp.Alerts)
	}
}

func DeleteTestFederationUsers(t *testing.T) {
	if len(testData.Federations) == 0 {
		t.Error("no federations test data")
	}

	if len(fedIDs) < 1 {
		t.Fatal("need at least one stored Federation ID to test Federations")
	}
	fedID := fedIDs[0]

	fedUsers, _, err := TOSession.GetFederationUsers(fedID, client.RequestOptions{})
	if err != nil {
		t.Fatalf("gettings users for federation %d: %v - alerts: %+v", fedID, err, fedUsers.Alerts)
	}
	if len(fedUsers.Response) != 3 {
		t.Errorf("federation users expected 3, actual: %d", len(fedUsers.Response))
	}

	for _, fedUser := range fedUsers.Response {
		if fedUser.ID == nil {
			t.Error("Traffic Ops returned a representation of a relationship between a user and a Federation that had null or undefined ID")
			continue
		}
		alerts, _, err := TOSession.DeleteFederationUser(fedID, *fedUser.ID, client.RequestOptions{})
		if err != nil {
			t.Fatalf("deleting user #%d from federation #%d: %v - alerts: %+v", *fedUser.ID, fedID, err, alerts.Alerts)
		}
	}

	fedUsers, _, err = TOSession.GetFederationUsers(fedID, client.RequestOptions{})
	if err != nil {
		t.Fatalf("gettings users for federation %d: %v - alerts: %+v", fedID, err, fedUsers.Alerts)
	}
	if len(fedUsers.Response) != 0 {
		t.Errorf("federation users expected 0, actual: %+v", len(fedUsers.Response))
	}
}
