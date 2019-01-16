package v14

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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestUsers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, DeliveryServices, Users}, func() {
		UpdateTestUsers(t)
		GetTestUsers(t)
		GetTestUserCurrent(t)
	})
}

const SessionUserName = "admin" // TODO make dynamic?

func CreateTestUsers(t *testing.T) {
	for _, user := range testData.Users {

		resp, _, err := TOSession.CreateUser(&user)
		if err != nil {
			t.Fatalf("could not CREATE user: %v\n", err)
		}
		log.Debugln("Response: ", resp.Alerts)
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
		t.Errorf("cannot GET users: %v\n", err)
	}
}

func GetTestUserCurrent(t *testing.T) {
	user, _, err := TOSession.GetUserCurrent()
	if err != nil {
		t.Errorf("cannot GET current user: %v\n", err)
	}
	if user.UserName == nil {
		t.Errorf("current user expected: %v actual: %v\n", SessionUserName, nil)
	}
	if *user.UserName != SessionUserName {
		t.Errorf("current user expected: %v actual: %v\n", SessionUserName, *user.UserName)
	}
}

// ForceDeleteTestUsers forcibly deletes the users from the db.
func ForceDeleteTestUsers(t *testing.T) {

	// NOTE: Special circumstances!  This should *NOT* be done without a really good reason!
	//  Connects directly to the DB to remove users rather than going thru the client.
	//  This is required here because the DeleteUser action does not really delete users,  but disables them.
	db, err := OpenConnection()
	if err != nil {
		t.Errorf("cannot open db")
	}
	defer db.Close()

	var usernames []string
	for _, user := range testData.Users {
		usernames = append(usernames, `'`+*user.Username+`'`)
	}

	q := `DELETE FROM tm_user WHERE username IN (` + strings.Join(usernames, ",") + `)`

	err = execSQL(db, q, "tm_user")

	if err != nil {
		t.Errorf("cannot execute SQL: %s; SQL is %s", err.Error(), q)
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
