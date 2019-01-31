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
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
)

func TestUsers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, DeliveryServices, Users}, func() {
		UpdateTestUsers(t)
		RolenameCapitalizationTest(t)
		OpsUpdateAdminTest(t)
		UserSelfUpdateTest(t)
		UserUpdateOwnRoleTest(t)
		GetTestUsers(t)
		GetTestUserCurrent(t)
	})
}

const SessionUserName = "admin" // TODO make dynamic?

func CreateTestUsers(t *testing.T) {
	for _, user := range testData.Users {

		resp, _, err := TOSession.CreateUser(&user)
		if err != nil {
			t.Errorf("could not CREATE user: %v\n", err)
		}
		log.Debugln("Response: ", resp.Alerts)
	}
}

func RolenameCapitalizationTest(t *testing.T) {

	roles, _, _, err := TOSession.GetRoles()
	if err != nil {
		t.Errorf("could not get roles: %v", err)
	}
	if len(roles) == 0 {
		t.Fatalf("there should be at least one role to test the user")
	}

	tenants, _, err := TOSession.Tenants()
	if err != nil {
		t.Errorf("could not get tenants: %v", err)
	}
	if len(tenants) == 0 {
		t.Fatalf("there should be at least one tenant to test the user")
	}

	// this user never does anything, so the role and tenant aren't important
	blob := fmt.Sprintf(`
	{
		"username": "test_user",
		"email": "cooldude6@example.com",
		"fullName": "full name is required",
		"localPasswd": "better_twelve",
		"confirmLocalPasswd": "better_twelve",
		"role": %d,
		"tenantId": %d 
	}`, *roles[0].ID, tenants[0].ID)

	reader := strings.NewReader(blob)
	request, err := http.NewRequest("POST", fmt.Sprintf("%v/api/1.4/users", TOSession.URL), reader)
	if err != nil {
		t.Errorf("could not make new request: %v\n", err)
	}
	resp, err := TOSession.Client.Do(request)
	if err != nil {
		t.Errorf("could not do request: %v\n", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	strResp := buf.String()
	if !strings.Contains(strResp, "roleName") {
		t.Errorf("incorrect json was returned for POST")
	}

	request, err = http.NewRequest("GET", fmt.Sprintf("%v/api/1.4/users?username=test_user", TOSession.URL), nil)
	resp, err = TOSession.Client.Do(request)

	buf = new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	strResp = buf.String()
	if !strings.Contains(strResp, "rolename") {
		t.Errorf("incorrect json was returned for GET")
	}

}

func OpsUpdateAdminTest(t *testing.T) {
	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	opsTOClient, _, err := toclient.LoginWithAgent(TOSession.URL, "opsuser", "pa$$word", true, "to-api-v14-client-tests/opsuser", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to get log in with opsuser: %v", err.Error())
	}

	resp, _, err := TOSession.GetUserByUsername("admin")
	if err != nil {
		t.Errorf("cannot GET user by name: 'admin', %v\n", err)
	}
	user := resp[0]

	fullName := "oops"
	email := "oops@ops.net"
	user.FullName = &fullName
	user.Email = &email

	_, _, err = opsTOClient.UpdateUserByID(*user.ID, &user)
	if err == nil {
		t.Errorf("ops user incorrectly updated an admin\n")
	}
}

func UserSelfUpdateTest(t *testing.T) {
	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	opsTOClient, _, err := toclient.LoginWithAgent(TOSession.URL, "opsuser", "pa$$word", true, "to-api-v14-client-tests/opsuser", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to get log in with opsuser: %v", err.Error())
	}

	resp, _, err := TOSession.GetUserByUsername("opsuser")
	if err != nil {
		t.Errorf("cannot GET user by name: 'opsuser', %v\n", err)
	}
	user := resp[0]

	fullName := "Oops-man"
	email := "operator@example.com"
	user.FullName = &fullName
	user.Email = &email

	var updateResp *tc.UpdateUserResponse
	updateResp, _, err = opsTOClient.UpdateUserByID(*user.ID, &user)
	if err != nil {
		t.Errorf("cannot UPDATE user by id: %v - %v\n", err, updateResp)
	}

	// Make sure it got updated
	resp2, _, err := TOSession.GetUserByID(*user.ID)
	if err != nil {
		t.Errorf("cannot GET user by id: '%d', %v\n", *user.ID, err)
	}
	updatedUser := resp2[0]

	if updatedUser.FullName == nil || updatedUser.Email == nil {
		t.Errorf("user was not correctly updated, field is null")
	}
	if *updatedUser.FullName != fullName {
		t.Errorf("results do not match actual: %s, expected: %s\n", *updatedUser.FullName, fullName)
	}
	if *updatedUser.Email != email {
		t.Errorf("results do not match acutal: %s, expected: %s\n", *updatedUser.Email, email)
	}
}

func UserUpdateOwnRoleTest(t *testing.T) {
	resp, _, err := TOSession.GetUserByUsername(SessionUserName)
	if err != nil {
		t.Errorf("cannot GET user by name: '%s', %v\n", SessionUserName, err)
	}
	user := resp[0]

	*user.Role = *user.Role + 1
	_, _, err = TOSession.UpdateUserByID(*user.ID, &user)
	if err == nil {
		t.Errorf("user incorrectly updated their role\n")
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
	if err != nil {
		t.Errorf("cannot GET user by id: '%d', %v\n", *user.ID, err)
	}
	updatedUser := resp2[0]
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

	// there is a constraint that prevents users from being deleted when they have a log
	q := `DELETE FROM log WHERE NOT tm_user = (SELECT id FROM tm_user WHERE username = 'admin')`
	err = execSQL(db, q, "log")
	if err != nil {
		t.Errorf("cannot execute SQL: %s; SQL is %s", err.Error(), q)
	}

	q = `DELETE FROM tm_user WHERE username IN (` + strings.Join(usernames, ",") + `)`
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
