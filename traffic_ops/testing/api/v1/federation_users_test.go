package v1

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
	"testing"
)

func TestFederationUsers(t *testing.T) {
	if Config.NoPerl {
		t.Skip("No Perl instance for proxying")
	}
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, DeliveryServices, UsersDeliveryServices, CDNFederations, FederationUsers}, func() {
		CreateTestInvalidFederationUsers(t)
		GetTestInvalidFederationIDUsers(t)
	})
}

func CreateTestFederationUsers(t *testing.T) {
	if len(testData.Federations) == 0 {
		t.Error("no federations test data")
	}

	fedID := fedIDs[0]

	// Get Users
	users, _, err := TOSession.GetUsers()
	if err != nil {
		t.Fatalf("getting users: " + err.Error())
	}
	if len(users) < 3 {
		t.Fatal("need > 3 users to create federation users")
	}

	u1 := users[0].ID
	u2 := users[1].ID
	u3 := users[2].ID

	// Associate one user to federation
	_, _, err = TOSession.CreateFederationUsers(fedID, []int{*u1}, false)
	if err != nil {
		t.Fatalf("assigning users %v to federation %v: %v", []int{*u1}, fedID, err.Error())
	}

	fedUsers, _, err := TOSession.GetFederationUsers(fedID)
	if err != nil {
		t.Fatalf("gettings users for federation %v: %v", fedID, err.Error())
	}
	if len(fedUsers) != 1 {
		t.Errorf("federation users expected 1, actual: %+v", len(fedUsers))
	}

	// Associate two users to federation and replace first one
	_, _, err = TOSession.CreateFederationUsers(fedID, []int{*u2, *u3}, true)
	if err != nil {
		t.Fatalf("assigning users %v to federation %v: %v", []int{*u2, *u3}, fedID, err.Error())
	}

	fedUsers, _, err = TOSession.GetFederationUsers(fedID)
	if err != nil {
		t.Fatalf("gettings users for federation %v: %v", fedID, err.Error())
	}
	if len(fedUsers) != 2 {
		t.Errorf("federation users expected 2, actual: %+v", len(fedUsers))
	}

	// Associate one more user to federation
	_, _, err = TOSession.CreateFederationUsers(fedID, []int{*u1}, false)
	if err != nil {
		t.Fatalf("assigning users %v to federation %v: %v", []int{*u1}, fedID, err.Error())
	}

	fedUsers, _, err = TOSession.GetFederationUsers(fedID)
	if err != nil {
		t.Fatalf("gettings users for federation %v: %v", fedID, err.Error())
	}
	if len(fedUsers) != 3 {
		t.Errorf("federation users expected 2, actual: %+v", len(fedUsers))
	}
}

func GetTestInvalidFederationIDUsers(t *testing.T) {
	_, _, err := TOSession.GetFederationUsers(-1)
	if err == nil {
		t.Fatalf("expected to get an error when requesting federation users for a non-existent federation")
	}
}

func CreateTestInvalidFederationUsers(t *testing.T) {
	if len(testData.Federations) == 0 {
		t.Error("no federations test data")
	}

	fedID := fedIDs[0]

	// Get Users
	users, _, err := TOSession.GetUsers()
	if err != nil {
		t.Fatalf("getting users: " + err.Error())
	}
	if len(users) == 0 {
		t.Fatal("need at least 1 user to test invalid federation user create")
	}

	// Associate with invalid federdation id
	_, _, err = TOSession.CreateFederationUsers(-1, []int{*users[0].ID}, false)
	if err == nil {
		t.Error("expected to get error back from associating non existent federation id")
	}

	// Associate with invalid user id
	_, _, err = TOSession.CreateFederationUsers(fedID, []int{-1}, false)
	if err == nil {
		t.Error("expected to get error back from associating non existent user id")
	}
}

func DeleteTestFederationUsers(t *testing.T) {
	if len(testData.Federations) == 0 {
		t.Error("no federations test data")
	}

	fedID := fedIDs[0]

	fedUsers, _, err := TOSession.GetFederationUsers(fedID)
	if err != nil {
		t.Fatalf("gettings users for federation %v: %v", fedID, err.Error())
	}
	if len(fedUsers) != 3 {
		t.Errorf("federation users expected 3, actual: %+v", len(fedUsers))
	}

	for _, fedUser := range fedUsers {
		_, _, err = TOSession.DeleteFederationUser(fedID, *fedUser.ID)
		if err != nil {
			t.Fatalf("deleting user %v from federation %v: %v", *fedUser.ID, fedID, err.Error())
		}
	}

	fedUsers, _, err = TOSession.GetFederationUsers(fedID)
	if err != nil {
		t.Fatalf("gettings users for federation %v: %v", fedID, err.Error())
	}
	if len(fedUsers) != 0 {
		t.Errorf("federation users expected 0, actual: %+v", len(fedUsers))
	}
}
