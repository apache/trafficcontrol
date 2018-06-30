package v13

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

	"github.com/apache/trafficcontrol/lib/go-log"
)

func TestUserDeliveryServicesAvailable(t *testing.T) {
	CreateTestCDNs(t)
	CreateTestTypes(t)
	CreateTestProfiles(t)
	CreateTestStatuses(t)
	CreateTestDivisions(t)
	CreateTestRegions(t)
	CreateTestPhysLocations(t)
	CreateTestCacheGroups(t)
	CreateTestServers(t)
	CreateTestDeliveryServices(t)

	GetTestUserDeliveryServicesAvailable(t)

	DeleteTestDeliveryServices(t)
	DeleteTestServers(t)
	DeleteTestCacheGroups(t)
	DeleteTestPhysLocations(t)
	DeleteTestRegions(t)
	DeleteTestDivisions(t)
	DeleteTestStatuses(t)
	DeleteTestProfiles(t)
	DeleteTestTypes(t)
	DeleteTestCDNs(t)
}

const TestUsersDeliveryServicesUser = "admin" // TODO make dynamic

func GetTestUserDeliveryServicesAvailable(t *testing.T) {
	log.Debugln("GetTestUserDeliveryServicesAvailable")

	users, _, err := TOSession.GetUsers()
	if err != nil {
		t.Fatalf("cannot GET users: %v\n", err)
	}
	if len(users) == 0 {
		t.Fatalf("no users, must have at least 1 user to test users_deliveryservices\n")
	}

	userID := 0
	foundUser := false
	for _, user := range users {
		log.Errorln("DEBUGQ user: " + user.Username)
		if user.Username == TestUsersDeliveryServicesUser {
			userID = user.ID
			foundUser = true
			break
		}
	}
	if !foundUser {
		t.Fatalf("get users expected: %v actual: missing\n", TestUsersDeliveryServicesUser)
	}

	dses, _, err := TOSession.GetUserDeliveryServicesAvailable(userID) // ([]tc.DeliveryServiceAvailableInfo, ReqInf, error) {
	if err != nil {
		t.Fatalf("cannot GET user DeliveryServices available: %v\n", err)
	}
	if len(dses) == 0 {
		t.Fatalf("GET users DeliveryServices available expected len(dses): >0 actual: %v\n", len(dses))
	}
}
