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

func TestUserDeliveryServices(t *testing.T) {
	if Config.NoPerl {
		t.Skip("No Perl instance for proxying")
	}
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, DeliveryServices, UsersDeliveryServices}, func() {
		GetTestUsersDeliveryServices(t)
	})
}

const TestUsersDeliveryServicesUser = "admin" // TODO make dynamic

func CreateTestUsersDeliveryServices(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v - %v", err, dses)
	}
	if len(dses) == 0 {
		t.Error("no delivery services, must have at least 1 ds to test users_deliveryservices")
	}
	users, _, err := TOSession.GetUsers()
	if err != nil {
		t.Errorf("cannot GET users: %v", err)
	}
	if len(users) == 0 {
		t.Error("no users, must have at least 1 user to test users_deliveryservices")
	}

	dsIDs := []int{}
	for _, ds := range dses {
		dsIDs = append(dsIDs, ds.ID)
	}

	userID := 0
	foundUser := false
	for _, user := range users {
		if *user.Username == TestUsersDeliveryServicesUser {
			userID = *user.ID
			foundUser = true
			break
		}
	}
	if !foundUser {
		t.Errorf("get users expected: %v actual: missing", TestUsersDeliveryServicesUser)
	}

	_, err = TOSession.SetDeliveryServiceUser(userID, dsIDs, true)
	if err != nil {
		t.Errorf("failed to set delivery service users: " + err.Error())
	}

	userDSes, _, err := TOSession.GetUserDeliveryServices(userID)
	if err != nil {
		t.Errorf("get user delivery services returned error: " + err.Error())
	}

	if len(userDSes.Response) != len(dsIDs) {
		t.Errorf("get user delivery services expected %v actual %v", len(dsIDs), len(userDSes.Response))
	}

	actualDSIDMap := map[int]struct{}{}
	for _, userDS := range userDSes.Response {
		if userDS.ID == nil {
			t.Error("get user delivery services returned a DS with a nil ID")
		}
		actualDSIDMap[*userDS.ID] = struct{}{}
	}
	for _, dsID := range dsIDs {
		if _, ok := actualDSIDMap[dsID]; !ok {
			t.Errorf("get user delivery services expected %v actual %v", dsID, "missing")
		}
	}
}

func GetTestUsersDeliveryServices(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v - %v", err, dses)
	}
	if len(dses) == 0 {
		t.Error("no delivery services, must have at least 1 ds to test users_deliveryservices")
	}
	users, _, err := TOSession.GetUsers()
	if err != nil {
		t.Errorf("cannot GET users: %v", err)
	}
	if len(users) == 0 {
		t.Error("no users, must have at least 1 user to test users_deliveryservices")
	}

	dsIDs := []int64{}
	for _, ds := range dses {
		dsIDs = append(dsIDs, int64(ds.ID))
	}

	userID := 0
	foundUser := false
	for _, user := range users {
		if *user.Username == TestUsersDeliveryServicesUser {
			userID = *user.ID
			foundUser = true
			break
		}
	}
	if !foundUser {
		t.Errorf("get users expected: %v actual: missing", TestUsersDeliveryServicesUser)
	}

	userDSes, _, err := TOSession.GetUserDeliveryServices(userID)
	if err != nil {
		t.Errorf("get user delivery services returned error: " + err.Error() + "\n")
	}

	if len(userDSes.Response) != len(dsIDs) {
		t.Errorf("get user delivery services expected %v actual %v", len(dsIDs), len(userDSes.Response))
	}

	actualDSIDMap := map[int]struct{}{}
	for _, userDS := range userDSes.Response {
		if userDS.ID == nil {
			t.Error("get user delivery services returned a DS with a nil ID")
		}
		actualDSIDMap[*userDS.ID] = struct{}{}
	}
	for _, dsID := range dsIDs {
		if _, ok := actualDSIDMap[int(dsID)]; !ok {
			t.Errorf("get user delivery services expected %v actual %v", dsID, "missing")
		}
	}
}

func DeleteTestUsersDeliveryServices(t *testing.T) {
	users, _, err := TOSession.GetUsers()
	if err != nil {
		t.Errorf("cannot GET users: %v", err)
	}
	if len(users) == 0 {
		t.Error("no users, must have at least 1 user to test users_deliveryservices")
	}
	userID := 0
	foundUser := false
	for _, user := range users {
		if *user.Username == TestUsersDeliveryServicesUser {
			userID = *user.ID
			foundUser = true
			break
		}
	}
	if !foundUser {
		t.Errorf("get users expected: %v actual: missing", TestUsersDeliveryServicesUser)
	}

	dses, _, err := TOSession.GetUserDeliveryServices(userID)
	if err != nil {
		t.Errorf("get user delivery services returned error: " + err.Error())
	}
	if len(dses.Response) == 0 {
		t.Errorf("get user delivery services expected %v actual %v", ">0", "0")
	}

	for _, ds := range dses.Response {
		if ds.ID == nil {
			t.Error("get user delivery services returned ds with nil ID")
		}
		_, err := TOSession.DeleteDeliveryServiceUser(userID, *ds.ID)
		if err != nil {
			t.Errorf("delete user delivery service returned error: " + err.Error())
		}
	}

	dses, _, err = TOSession.GetUserDeliveryServices(userID)
	if err != nil {
		t.Errorf("get user delivery services returned error: " + err.Error())
	}
	if len(dses.Response) != 0 {
		t.Errorf("get user delivery services after deleting expected %v actual %v", "0", len(dses.Response))
	}
}
