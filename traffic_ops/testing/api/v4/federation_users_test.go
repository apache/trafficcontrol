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
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestFederationUsers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations, FederationUsers}, func() {
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		CreateTestInvalidFederationUsers(t)
		GetTestInvalidFederationIDUsers(t)
		CreateTestValidFederationUsers(t)
		GetTestValidFederationIDUsersIMSAfterChange(t, header)
		GetTestPaginationSupportFedUsers(t)
		SortTestFederationUsers(t)
		SortTestFederationsUsersDesc(t)
	})
}

func GetTestValidFederationIDUsersIMSAfterChange(t *testing.T, header http.Header) {
	if len(fedIDs) < 1 {
		t.Fatal("Need at least one Federation ID to test Federation ID Users change")
	}
	opts := client.RequestOptions{Header: header}
	resp, reqInf, err := TOSession.GetFederationUsers(fedIDs[0], opts)
	if err != nil {
		t.Fatalf("No error expected, but got: %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	opts.Header.Set(rfc.IfModifiedSince, timeStr)
	resp, reqInf, err = TOSession.GetFederationUsers(fedIDs[0], opts)
	if err != nil {
		t.Fatalf("No error expected, but got: %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
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

	fedUsers, _, err := TOSession.GetFederationUsers(fedID, client.RequestOptions{})
	if err != nil {
		t.Fatalf("gettings users for federation %v: %v - alerts: %+v", fedID, err, fedUsers.Alerts)
	}
	if len(fedUsers.Response) != 1 {
		t.Errorf("federation users expected 1, actual: %d", len(fedUsers.Response))
	}

	// Associate two users to federation and replace first one
	resp, _, err = TOSession.CreateFederationUsers(fedID, []int{*u2, *u3}, true, client.RequestOptions{})
	if err != nil {
		t.Fatalf("assigning users %v to federation %d: %v - alerts: %+v", []int{*u2, *u3}, fedID, err, resp.Alerts)
	}

	fedUsers, _, err = TOSession.GetFederationUsers(fedID, client.RequestOptions{})
	if err != nil {
		t.Fatalf("gettings users for federation %d: %v - alerts: %+v", fedID, err, fedUsers.Alerts)
	}
	if len(fedUsers.Response) != 2 {
		t.Errorf("federation users expected 2, actual: %d", len(fedUsers.Response))
	}

	// Associate one more user to federation
	resp, _, err = TOSession.CreateFederationUsers(fedID, []int{*u1}, false, client.RequestOptions{})
	if err != nil {
		t.Fatalf("assigning users %v to federation %d: %v - alerts: %+v", []int{*u1}, fedID, err, resp.Alerts)
	}

	fedUsers, _, err = TOSession.GetFederationUsers(fedID, client.RequestOptions{})
	if err != nil {
		t.Fatalf("gettings users for federation %d: %v - alerts: %+v", fedID, err, fedUsers.Alerts)
	}
	if len(fedUsers.Response) != 3 {
		t.Errorf("federation users expected 2, actual: %d", len(fedUsers.Response))
	}
}

func GetTestInvalidFederationIDUsers(t *testing.T) {
	_, _, err := TOSession.GetFederationUsers(-1, client.RequestOptions{})
	if err == nil {
		t.Fatalf("expected to get an error when requesting federation users for a non-existent federation")
	}
}

func CreateTestValidFederationUsers(t *testing.T) {
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
	if len(users.Response) == 0 {
		t.Fatal("need at least 1 user to test invalid federation user create")
	}

	// Associate with invalid federdation id
	_, _, err = TOSession.CreateFederationUsers(fedID, []int{*users.Response[0].ID}, false, client.RequestOptions{})
	if err == nil {
		t.Error("expected to get error back from associating non existent federation id")
	}
}
func CreateTestInvalidFederationUsers(t *testing.T) {
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
	if len(users.Response) == 0 {
		t.Fatal("need at least 1 user to test invalid federation user create")
	}
	if users.Response[0].ID == nil {
		t.Fatal("Traffic Ops returned a representation of a user with null or undefined ID")
	}

	// Associate with invalid federdation id
	_, _, err = TOSession.CreateFederationUsers(-1, []int{*users.Response[0].ID}, false, client.RequestOptions{})
	if err == nil {
		t.Error("expected to get error back from associating non existent federation id")
	}

	// Associate with invalid user id
	_, _, err = TOSession.CreateFederationUsers(fedID, []int{-1}, false, client.RequestOptions{})
	if err == nil {
		t.Error("expected to get error back from associating non existent user id")
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

func GetTestPaginationSupportFedUsers(t *testing.T) {

	if len(fedIDs) < 1 {
		t.Fatal("need at least one stored Federation ID to test Federations")
	}
	fedID := fedIDs[0]

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "userID")
	resp, _, err := TOSession.GetFederationUsers(fedID, opts)
	if err != nil {
		t.Fatalf("cannot get Federation #%d Users: %v - alerts: %+v", fedID, err, resp.Alerts)
	}
	fedUser := resp.Response
	if len(fedUser) < 3 {
		t.Fatalf("Need at least 3 Federation users in Traffic Ops to test pagination support, found: %d", len(fedUser))
	}

	opts.QueryParameters.Set("limit", "1")
	fedUserWithLimit, _, err := TOSession.GetFederationUsers(fedID, opts)
	if err != nil {
		t.Fatalf("cannot Get Federation user with Limit: %v - alerts: %+v", err, fedUserWithLimit.Alerts)
	}
	if !reflect.DeepEqual(fedUser[:1], fedUserWithLimit.Response) {
		t.Error("expected GET Federation user with limit = 1 to return first result")
	}

	opts.QueryParameters.Set("offset", "1")
	fedUserWithOffset, _, err := TOSession.GetFederationUsers(fedID, opts)
	if err != nil {
		t.Fatalf("cannot Get Federation user with Limit and Offset: %v - alerts: %+v", err, fedUserWithOffset.Alerts)
	}
	if !reflect.DeepEqual(fedUser[1:2], fedUserWithOffset.Response) {
		t.Error("expected GET Federation user with limit = 1, offset = 1 to return second result")
	}

	opts.QueryParameters.Del("offset")
	opts.QueryParameters.Set("page", "2")
	fedUserWithPage, _, err := TOSession.GetFederationUsers(fedID, opts)
	if err != nil {
		t.Fatalf("cannot Get Federation user with Limit and Page: %v - alerts: %+v", err, fedUserWithPage.Alerts)
	}
	if !reflect.DeepEqual(fedUser[1:2], fedUserWithPage.Response) {
		t.Error("expected GET Federation user with limit = 1, page = 2 to return second result")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "-2")
	resp, reqInf, err := TOSession.GetFederationUsers(fedID, opts)
	if err == nil {
		t.Error("expected GET Federation user to return an error when limit is not bigger than -1")
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
	}

	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "0")
	resp, reqInf, err = TOSession.GetFederationUsers(fedID, opts)
	if err == nil {
		t.Error("expected GET Federation user to return an error when offset is not a positive integer")
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "0")
	resp, reqInf, err = TOSession.GetFederationUsers(fedID, opts)
	if err == nil {
		t.Error("expected GET Federation user to return an error when page is not a positive integer")
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
	}
}

func SortTestFederationUsers(t *testing.T) {
	var sortedList []int
	if len(fedIDs) == 0 {
		t.Fatalf("no federations, must have at least 1 federation to test federations users")
	}
	fedID := fedIDs[0]
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "userID")
	resp, _, err := TOSession.GetFederationUsers(fedID, opts)
	if err != nil {
		t.Fatalf("Expected no error, but got %v - alerts: %+v", err, resp.Alerts)
	}
	for _, fedRes := range resp.Response {
		sortedList = append(sortedList, *fedRes.ID)
	}
	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
		t.Errorf("list is not sorted by their ID: %v", sortedList)
	}
}

func SortTestFederationsUsersDesc(t *testing.T) {

	if len(fedIDs) == 0 {
		t.Fatalf("no federations, must have at least 1 federation to test federations users")
	}
	fedID := fedIDs[0]

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "userID")
	resp, _, err := TOSession.GetFederationUsers(fedID, opts)
	if err != nil {
		t.Fatalf("Expected no error, but got error in Federation users default ordering %v - alerts: %+v", err, resp.Alerts)
	}
	respAsc := resp.Response
	if len(respAsc) < 1 {
		t.Fatal("Need at least one Federation users in Traffic Ops to test Federation users sort ordering")
	}
	opts.QueryParameters.Set("sortOrder", "desc")
	resp, _, err = TOSession.GetFederationUsers(fedID, opts)
	if err != nil {
		t.Errorf("Expected no error, but got error in Federation users with Descending ordering: %v - alerts: %+v", err, resp.Alerts)
	}
	respDesc := resp.Response
	if len(respDesc) < 1 {
		t.Fatal("Need at least one Federation users in Traffic Ops to test Federation users sort ordering")
	}
	if len(respAsc) != len(respDesc) {
		t.Fatalf("Traffic Ops returned %d Federation users using default sort order, but %d Federation users when sort order was explicitly set to descending", len(respAsc), len(respDesc))
	}
	for start, end := 0, len(respDesc)-1; start < end; start, end = start+1, end-1 {
		respDesc[start], respDesc[end] = respDesc[end], respDesc[start]
	}
	if *respDesc[0].ID != *respAsc[0].ID {
		t.Errorf("Federation users responses are not equal after reversal: Asc: %d - Desc: %d", *respDesc[0].ID, *respAsc[0].ID)
	}
}
