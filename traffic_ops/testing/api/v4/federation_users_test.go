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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestFederationUsers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices, CDNFederations, FederationUsers}, func() {
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		CreateTestInvalidFederationUsers(t)
		GetTestInvalidFederationIDUsers(t)
		CreateTestValidFederationUsers(t)
		GetTestValidFederationIDUsersIMSAfterChange(t, header)
	})
}

func GetTestValidFederationIDUsersIMSAfterChange(t *testing.T, header http.Header) {
	if len(fedIDs) < 0 {
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
