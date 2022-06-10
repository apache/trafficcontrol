package v2

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
	"net/mail"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v2-client"
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
		UserTenancyTest(t)
		if includeSystemTests {
			// UserRegistrationTest deletes test users before registering new users, so it must come after the other user tests.
			UserRegistrationTest(t)
		}
	})
}

const SessionUserName = "admin" // TODO make dynamic?

func CreateTestUsers(t *testing.T) {
	for _, user := range testData.Users {

		resp, _, err := TOSession.CreateUser(&user)
		if err != nil {
			t.Errorf("could not CREATE user: %v", err)
		}
		t.Log("Response: ", resp.Alerts)
	}
}

func RolenameCapitalizationTest(t *testing.T) {

	roles, _, _, err := TOSession.GetRoles()
	if err != nil {
		t.Errorf("could not get roles: %v", err)
	}
	if len(roles) == 0 {
		t.Fatal("there should be at least one role to test the user")
	}

	tenants, _, err := TOSession.Tenants()
	if err != nil {
		t.Errorf("could not get tenants: %v", err)
	}
	if len(tenants) == 0 {
		t.Fatal("there should be at least one tenant to test the user")
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
	request, err := http.NewRequest("POST", fmt.Sprintf("%v/api/2.0/users", TOSession.URL), reader)
	if err != nil {
		t.Errorf("could not make new request: %v", err)
	}
	resp, err := TOSession.Client.Do(request)
	if err != nil {
		t.Errorf("could not do request: %v", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	strResp := buf.String()
	if !strings.Contains(strResp, "roleName") {
		t.Error("incorrect json was returned for POST")
	}

	request, err = http.NewRequest("GET", fmt.Sprintf("%v/api/2.0/users?username=test_user", TOSession.URL), nil)
	resp, err = TOSession.Client.Do(request)

	buf = new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	strResp = buf.String()
	if !strings.Contains(strResp, "rolename") {
		t.Error("incorrect json was returned for GET")
	}

}

func OpsUpdateAdminTest(t *testing.T) {
	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	opsTOClient, _, err := toclient.LoginWithAgent(TOSession.URL, "opsuser", "pa$$word", true, "to-api-v2-client-tests/opsuser", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to get log in with opsuser: %v", err.Error())
	}

	resp, _, err := TOSession.GetUserByUsername("admin")
	if err != nil {
		t.Errorf("cannot GET user by name: 'admin', %v", err)
	}
	user := resp[0]

	fullName := "oops"
	email := "oops@ops.net"
	user.FullName = &fullName
	user.Email = &email

	_, _, err = opsTOClient.UpdateUserByID(*user.ID, &user)
	if err == nil {
		t.Error("ops user incorrectly updated an admin")
	}
}

func UserRegistrationTest(t *testing.T) {
	ForceDeleteTestUsers(t)
	var emails []string
	for _, user := range testData.Users {
		tenant, _, err := TOSession.TenantByName(*user.Tenant)
		if err != nil {
			t.Fatalf("could not get tenant %v: %v", *user.Tenant, err)
		}
		resp, _, err := TOSession.RegisterNewUser(uint(tenant.ID), uint(*user.Role), rfc.EmailAddress{Address: mail.Address{Address: *user.Email}})
		if err != nil {
			t.Fatalf("could not register user: %v", err)
		}
		t.Log("Response: ", resp.Alerts)
		emails = append(emails, fmt.Sprintf(`'%v'`, *user.Email))
	}

	db, err := OpenConnection()
	if err != nil {
		t.Error("cannot open db")
	}
	defer db.Close()
	q := `DELETE FROM tm_user WHERE email IN (` + strings.Join(emails, ",") + `)`
	if err := execSQL(db, q); err != nil {
		t.Errorf("cannot execute SQL to delete registered users: %s; SQL is %s", err.Error(), q)
	}
}

func UserSelfUpdateTest(t *testing.T) {
	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	opsTOClient, _, err := toclient.LoginWithAgent(TOSession.URL, "opsuser", "pa$$word", true, "to-api-v2-client-tests/opsuser", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to get log in with opsuser: %v", err.Error())
	}

	resp, _, err := TOSession.GetUserByUsername("opsuser")
	if err != nil {
		t.Fatalf("cannot GET user by name: 'opsuser', %v\n", err)
	}
	if len(resp) < 1 {
		t.Fatalf("no users returned when requesting user 'opsuser'")
	}
	user := resp[0]

	if user.ID == nil {
		t.Fatalf("user 'opsuser' has a null or missing ID - cannot proceed")
	}

	user.FullName = util.StrPtr("Oops-man")
	user.Email = util.StrPtr("operator@not.example.com")

	var updateResp *tc.UpdateUserResponse
	updateResp, _, err = opsTOClient.UpdateUserByID(*user.ID, &user)
	if err != nil {
		t.Fatalf("cannot UPDATE user by id: %v - %v\n", err, updateResp)
	}

	// Make sure it got updated
	resp2, _, err := TOSession.GetUserByID(*user.ID)
	if err != nil {
		t.Fatalf("cannot GET user by id: '%d', %v\n", *user.ID, err)
	}
	if len(resp2) < 1 {
		t.Fatalf("no results returned when requesting user #%d", *user.ID)
	}
	updatedUser := resp2[0]

	if updatedUser.FullName == nil {
		t.Errorf("user was not correctly updated, FullName is null or missing")
	} else if *updatedUser.FullName != "Oops-man" {
		t.Errorf("results do not match actual: '%s', expected: 'Oops-man'\n", *updatedUser.FullName)
	}

	if updatedUser.Email == nil {
		t.Errorf("user was not correctly updated, Email is null or missing")
	} else if *updatedUser.Email != "operator@not.example.com" {
		t.Errorf("results do not match actual: '%s', expected: 'operator@not.example.com'\n", *updatedUser.Email)
	}

	// Same thing using /user/current
	user.FullName = util.StrPtr("ops-man")
	user.Email = util.StrPtr("operator@example.com")
	updateResp, _, err = opsTOClient.UpdateCurrentUser(user)
	if err != nil {
		t.Fatalf("error updating current user: %v - %v", err, updateResp)
	}

	// Make sure it got updated
	resp2, _, err = TOSession.GetUserByID(*user.ID)
	if err != nil {
		t.Fatalf("error getting user #%d: %v", *user.ID, err)
	}

	if len(resp2) < 1 {
		t.Fatalf("no user returned when requesting user #%d", *user.ID)
	}

	if resp2[0].FullName == nil {
		t.Errorf("FullName missing or null after update")
	} else if *resp2[0].FullName != "ops-man" {
		t.Errorf("Expected FullName to be 'ops-man', but it was '%s'", *resp2[0].FullName)
	}

	if resp2[0].Email == nil {
		t.Errorf("Email missing or null after update")
	} else if *resp2[0].Email != "operator@example.com" {
		t.Errorf("Expected Email to be restored to 'operator@example.com', but it was '%s'", *resp2[0].Email)
	}

	// now test using an invalid email address
	currentEmail := *user.Email
	user.Email = new(string)
	updateResp, _, err = TOSession.UpdateCurrentUser(user)
	if err == nil {
		t.Fatal("error was expected updating user with email: '' - got none")
	}

	// Ensure it wasn't actually updated
	resp2, _, err = TOSession.GetUserByID(*user.ID)
	if err != nil {
		t.Fatalf("error getting user #%d: %v", *user.ID, err)
	}

	if len(resp2) < 1 {
		t.Fatalf("no user returned when requesting user #%d", *user.ID)
	}

	if resp2[0].Email == nil {
		t.Errorf("Email missing or null after update")
	} else if *resp2[0].Email != currentEmail {
		t.Errorf("Expected Email to still be '%s', but it was '%s'", currentEmail, *resp2[0].Email)
	}

	// Now test using an invalid username
	currentUsername := *user.Username
	user.Username = new("ops man")
	updateResp, _, err = TOSession.UpdateCurrentUser(user)
	if err == nil {
		t.Fatal("error was expected updating user with username: 'ops man' - got none")
	}

	// Ensure it wasn't actually updated
	resp2, _, err = TOSession.GetUserByID(*user.ID)
	if err != nil {
		t.Fatalf("error getting user #%d: %v", *user.ID, err)
	}

	if len(resp2) < 1 {
		t.Fatalf("no user returned when requesting user #%d", *user.ID)
	}

	if resp2[0].Username == nil {
		t.Errorf("Username missing or null after update")
	} else if *resp2[0].Username != currentUsername {
		t.Errorf("Expected Username to still be '%s', but it was '%s'", currentUsername, *resp2[0].Username)
	}
}

func UserUpdateOwnRoleTest(t *testing.T) {
	resp, _, err := TOSession.GetUserByUsername(SessionUserName)
	if err != nil {
		t.Errorf("cannot GET user by name: '%s', %v", SessionUserName, err)
	}
	user := resp[0]

	*user.Role = *user.Role + 1
	_, _, err = TOSession.UpdateUserByID(*user.ID, &user)
	if err == nil {
		t.Error("user incorrectly updated their role")
	}
}

func UpdateTestUsers(t *testing.T) {
	firstUsername := *testData.Users[0].Username
	resp, _, err := TOSession.GetUserByUsername(firstUsername)
	if err != nil {
		t.Errorf("cannot GET user by name: '%s', %v", firstUsername, err)
	}
	user := resp[0]
	newCity := "kidz kable kown"
	*user.City = newCity

	var updateResp *tc.UpdateUserResponse
	updateResp, _, err = TOSession.UpdateUserByID(*user.ID, &user)
	if err != nil {
		t.Errorf("cannot UPDATE user by id: %v - %v", err, updateResp.Alerts)
	}

	// Make sure it got updated
	resp2, _, err := TOSession.GetUserByID(*user.ID)
	if err != nil {
		t.Errorf("cannot GET user by id: '%d', %v", *user.ID, err)
	}
	updatedUser := resp2[0]
	if *updatedUser.City != newCity {
		t.Errorf("results do not match actual: %s, expected: %s", *updatedUser.City, newCity)
	}
}

func GetTestUsers(t *testing.T) {
	_, _, err := TOSession.GetUsers()
	if err != nil {
		t.Errorf("cannot GET users: %v", err)
	}
}

func GetTestUserCurrent(t *testing.T) {
	user, _, err := TOSession.GetUserCurrent()
	if err != nil {
		t.Errorf("cannot GET current user: %v", err)
	}
	if user.UserName == nil {
		t.Errorf("current user expected: %v actual: %v", SessionUserName, nil)
	}
	if *user.UserName != SessionUserName {
		t.Errorf("current user expected: %v actual: %v", SessionUserName, *user.UserName)
	}
}

func UserTenancyTest(t *testing.T) {
	users, _, err := TOSession.GetUsers()
	if err != nil {
		t.Errorf("cannot GET users: %v", err)
	}
	tenant3Found := false
	tenant4Found := false
	tenant3Username := "tenant3user"
	tenant4Username := "tenant4user"
	tenant3User := tc.User{}

	// assert admin user can view tenant3user and tenant4user
	for _, user := range users {
		if *user.Username == tenant3Username {
			tenant3Found = true
			tenant3User = user
		} else if *user.Username == tenant4Username {
			tenant4Found = true
		}
		if tenant3Found && tenant4Found {
			break
		}
	}
	if !tenant3Found || !tenant4Found {
		t.Error("expected admin to be able to view tenants: tenant3 and tenant4")
	}

	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	tenant4TOClient, _, err := toclient.LoginWithAgent(TOSession.URL, "tenant4user", "pa$$word", true, "to-api-v2-client-tests/tenant4user", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to log in with tenant4user: %v", err.Error())
	}

	usersReadableByTenant4, _, err := tenant4TOClient.GetUsers()
	if err != nil {
		t.Error("tenant4user cannot GET users")
	}

	tenant4canReadItself := false
	for _, user := range usersReadableByTenant4 {
		// assert that tenant4user cannot read tenant3user
		if *user.Username == tenant3Username {
			t.Error("expected tenant4user to be unable to read tenant3user")
		}
		// assert that tenant4user can read itself
		if *user.Username == tenant4Username {
			tenant4canReadItself = true
		}
	}
	if !tenant4canReadItself {
		t.Error("expected tenant4user to be able to read itself")
	}

	// assert that tenant4user cannot update tenant3user
	if _, _, err = tenant4TOClient.UpdateUserByID(*tenant3User.ID, &tenant3User); err == nil {
		t.Error("expected tenant4user to be unable to update tenant4user")
	}

	// assert that tenant4user cannot create a user outside of its tenant
	rootTenant, _, err := TOSession.TenantByName("root")
	if err != nil {
		t.Error("expected to be able to GET the root tenant")
	}
	newUser := testData.Users[0]
	newUser.Email = util.StrPtr("testusertenancy@example.com")
	newUser.Username = util.StrPtr("testusertenancy")
	newUser.TenantID = &rootTenant.ID
	if _, _, err = tenant4TOClient.CreateUser(&newUser); err == nil {
		t.Error("expected tenant4user to be unable to create a new user in the root tenant")
	}
}

// ForceDeleteTestUsers forcibly deletes the users from the db.
func ForceDeleteTestUsers(t *testing.T) {

	// NOTE: Special circumstances!  This should *NOT* be done without a really good reason!
	//  Connects directly to the DB to remove users rather than going thru the client.
	//  This is required here because the DeleteUser action does not really delete users,  but disables them.
	db, err := OpenConnection()
	if err != nil {
		t.Error("cannot open db")
	}
	defer db.Close()

	var usernames []string
	for _, user := range testData.Users {
		usernames = append(usernames, `'`+*user.Username+`'`)
	}

	// there is a constraint that prevents users from being deleted when they have a log
	q := `DELETE FROM log WHERE NOT tm_user = (SELECT id FROM tm_user WHERE username = 'admin')`
	err = execSQL(db, q)
	if err != nil {
		t.Errorf("cannot execute SQL: %s; SQL is %s", err.Error(), q)
	}

	q = `DELETE FROM tm_user WHERE username IN (` + strings.Join(usernames, ",") + `)`
	err = execSQL(db, q)
	if err != nil {
		t.Errorf("cannot execute SQL: %s; SQL is %s", err.Error(), q)
	}
}

func DeleteTestUsers(t *testing.T) {
	for _, user := range testData.Users {

		resp, _, err := TOSession.GetUserByUsername(*user.Username)
		if err != nil {
			t.Errorf("cannot GET user by name: %v - %v", *user.Username, err)
		}

		if resp != nil {
			respUser := resp[0]

			_, _, err := TOSession.DeleteUserByID(*respUser.ID)
			if err != nil {
				t.Errorf("cannot DELETE user by name: '%s' %v", *respUser.Username, err)
			}

			// Make sure it got deleted
			resp, _, err := TOSession.GetUserByUsername(*user.Username)
			if err != nil {
				t.Errorf("error deleting user by name: %s", err.Error())
			}
			if len(resp) > 0 {
				t.Errorf("expected user: %s to be deleted", *user.Username)
			}
		}
	}
}
