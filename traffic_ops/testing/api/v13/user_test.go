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
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"testing"
)

func TestUsers(t *testing.T) {
	CreateTestCDNs(t)
	CreateTestTypes(t)
	CreateTestProfiles(t)
	CreateTestStatuses(t)
	CreateTestDivisions(t)
	CreateTestRegions(t)
	CreateTestPhysLocations(t)
	CreateTestCacheGroups(t)
	CreateTestDeliveryServices(t)

	CreateTestUsers(t)
	UpdateTestUsers(t)
	GetTestUsers(t)
	GetTestUserCurrent(t)
	//Delete will be new functionality to 1.4, ignore for now
	//DeleteTestUsers(t)

	DeleteTestDeliveryServices(t)
	DeleteTestCacheGroups(t)
	DeleteTestPhysLocations(t)
	DeleteTestRegions(t)
	DeleteTestDivisions(t)
	DeleteTestStatuses(t)
	DeleteTestProfiles(t)
	DeleteTestTypes(t)
	DeleteTestCDNs(t)

}

const SessionUserName = "admin" // TODO make dynamic?

func CreateTestUsers(t *testing.T) {
	for _, user := range testData.Users {

		resp, _, err := TOSession.CreateUser(&user)
		log.Debugln("Response: ", resp.Alerts)

		if err != nil {
			t.Errorf("could not CREATE user: %v\n", err)
		}
	}
}

func UpdateTestUsers(t *testing.T) {
	firstUsername := *testData.Users[0].Username
	resp, _, err := TOSession.GetUserByUsername(firstUsername)
	if err != nil {
		t.Errorf("cannot GET user by name: '%s', %v\n", firstUsername, err)
	}
	user := resp[0]
	newCity := "kidz kable kown"
	*user.City = newCity

	var updateResp *tc.UpdateUserResponse
	updateResp, _, err = TOSession.UpdateUserByID(*user.ID, &user)
	if err != nil {
		t.Errorf("cannot UPDATE user by id: %v - %v\n", err, updateResp.Alerts)
	}

	// Make sure it got updated
	resp2, _, err := TOSession.GetUserByID(*user.ID)
	updatedUser := resp2[0]

	if err != nil {
		t.Errorf("cannot GET user by id: '%d', %v\n", *user.ID, err)
	}
	if *updatedUser.City != newCity {
		t.Errorf("results do not match actual: %s, expected: %s\n", *updatedUser.City, newCity)
	}
}

func GetTestUsers(t *testing.T) {
	_, _, err := TOSession.GetUsers()
	if err != nil {
		t.Fatalf("cannot GET users: %v\n", err)
	}
}

func GetTestUserCurrent(t *testing.T) {
	user, _, err := TOSession.GetUserCurrent()
	if err != nil {
		t.Fatalf("cannot GET current user: %v\n", err)
	}
	if user.UserName == nil {
		t.Fatalf("current user expected: %v actual: %v\n", SessionUserName, nil)
	}
	if *user.UserName != SessionUserName {
		t.Fatalf("current user expected: %v actual: %v\n", SessionUserName, *user.UserName)
	}
}

func DeleteTestUsers(t *testing.T) {
	for _, user := range testData.Users {

		resp, _, err := TOSession.GetUserByUsername(*user.Username)
		if err != nil {
			t.Errorf("cannot GET user by name: %v - %v\n", *user.Username, err)
		}

		if resp != nil {
			respUser := resp[0]

			_, _, err := TOSession.DeleteUserByID(*respUser.ID)
			if err != nil {
				t.Errorf("cannot DELETE user by name: '%s' %v\n", *respUser.Username, err)
			}

			// Make sure it got deleted
			resp, _, err := TOSession.GetUserByUsername(*user.Username)
			if err != nil {
				t.Errorf("error deleting user by name: %s\n", err.Error())
			}
			if len(resp) > 0 {
				t.Errorf("expected user: %s to be deleted\n", *user.Username)
			}
		}
	}
}
