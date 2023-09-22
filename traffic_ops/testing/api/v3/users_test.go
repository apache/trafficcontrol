package v3

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
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestUsers(t *testing.T) {
	WithObjs(t, []TCObj{Tenants, Parameters, Users}, func() {

		opsUserSession := utils.CreateV3Session(t, Config.TrafficOps.URL, "opsuser", "pa$$word", Config.Default.Session.TimeoutInSecs)
		tenant4UserSession := utils.CreateV3Session(t, Config.TrafficOps.URL, "tenant4user", "pa$$word", Config.Default.Session.TimeoutInSecs)

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateUsersSort()),
				},
				"ADMIN can view CHILD TENANTS": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateTenants(map[string]bool{"tenant3": true, "tenant4": true})),
				},
				"CHILD TENANT should NOT read PARENT TENANT": {
					ClientSession: tenant4UserSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateTenants(map[string]bool{"tenant3": false, "tenant4": true})),
				},
			},
			"POST": {
				"FORBIDDEN when CHILD TENANT creates USER with PARENT TENANCY": {
					ClientSession: tenant4UserSession,
					RequestBody: map[string]interface{}{
						"email":              "outsidetenancy@example.com",
						"fullName":           "Outside Tenancy",
						"localPasswd":        "pa$$word",
						"confirmLocalPasswd": "pa$$word",
						"role":               3,
						"tenantId":           GetTenantID(t, "tenant3")(),
						"username":           "outsideTenantUser",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetUserID(t, "steering"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"addressLine1":       "updated line 1",
						"addressLine2":       "updated line 2",
						"city":               "updated city name",
						"company":            "new company",
						"country":            "US",
						"email":              "steeringupdated@example.com",
						"fullName":           "Steering User Updated",
						"localPasswd":        "pa$$word",
						"confirmLocalPasswd": "pa$$word",
						"newUser":            false,
						"role":               6,
						"tenant":             "root",
						"tenantId":           GetTenantID(t, "root")(),
						"username":           "steering",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateUsersUpdateCreateFields(map[string]interface{}{"AddressLine1": "updated line 1",
							"AddressLine2": "updated line 2", "City": "updated city name", "Company": "new company",
							"Country": "US", "Email": "steeringupdated@example.com", "FullName": "Steering User Updated"})),
				},
				"OK when UPDATING SELF": {
					EndpointID:    GetUserID(t, "opsuser"),
					ClientSession: opsUserSession,
					RequestBody: map[string]interface{}{
						"addressLine1":       "address of ops",
						"addressLine2":       "place",
						"city":               "somewhere",
						"company":            "else",
						"country":            "UK",
						"email":              "ops-updated@example.com",
						"fullName":           "Operations User Updated",
						"localPasswd":        "pa$$word",
						"confirmLocalPasswd": "pa$$word",
						"role":               3,
						"tenant":             "root",
						"tenantId":           GetTenantID(t, "root")(),
						"username":           "opsuser",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateUsersUpdateCreateFields(map[string]interface{}{"Email": "ops-updated@example.com", "FullName": "Operations User Updated"})),
				},
				"BAD REQUEST when updating OWN ROLE": {
					EndpointID:    GetUserID(t, "opsuser"),
					ClientSession: opsUserSession,
					RequestBody: map[string]interface{}{
						"addressLine1":       "address of ops",
						"addressLine2":       "place",
						"city":               "somewhere",
						"company":            "else",
						"country":            "UK",
						"email":              "ops-updated@example.com",
						"fullName":           "Operations User Updated",
						"localPasswd":        "pa$$word",
						"confirmLocalPasswd": "pa$$word",
						"role":               9999,
						"tenant":             "root",
						"tenantId":           GetTenantID(t, "root")(),
						"username":           "opsuser",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"FORBIDDEN when OPERATIONS USER updates ADMIN USER": {
					EndpointID:    GetUserID(t, "admin"),
					ClientSession: opsUserSession,
					RequestBody: map[string]interface{}{
						"email":              "oops@ops.net",
						"fullName":           "oops",
						"localPasswd":        "pa$$word",
						"confirmLocalPasswd": "pa$$word",
						"role":               4,
						"tenant":             "root",
						"tenantId":           GetTenantID(t, "root")(),
						"username":           "admin",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"FORBIDDEN when CHILD TENANT USER updates PARENT TENANT USER": {
					EndpointID:    GetUserID(t, "tenant3user"),
					ClientSession: tenant4UserSession,
					RequestBody: map[string]interface{}{
						"email":              "tenant3user@example.com",
						"fullName":           "Parent tenant test",
						"localPasswd":        "pa$$word",
						"confirmLocalPasswd": "pa$$word",
						"role":               4,
						"tenant":             "tenant2",
						"tenantId":           GetTenantID(t, "tenant2")(),
						"username":           "tenant3user",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					user := tc.User{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &user)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetUsersWithHdr(testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateUser(&user)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateUserByID(testCase.EndpointID(), &user)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func validateUsersFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Users response to not be nil.")
		userResp := resp.([]tc.User)
		for field, expected := range expectedResp {
			for _, user := range userResp {
				switch field {
				case "AddressLine1":
					assert.RequireNotNil(t, user.AddressLine1, "Expected AddressLine1 to not be nil.")
					assert.Equal(t, expected, *user.AddressLine1, "Expected AddressLine1 to be %v, but got %s", expected, *user.AddressLine1)
				case "AddressLine2":
					assert.RequireNotNil(t, user.AddressLine2, "Expected AddressLine2 to not be nil.")
					assert.Equal(t, expected, *user.AddressLine2, "Expected AddressLine2 to be %v, but got %s", expected, *user.AddressLine2)
				case "City":
					assert.RequireNotNil(t, user.City, "Expected City to not be nil.")
					assert.Equal(t, expected, *user.City, "Expected City to be %v, but got %s", expected, *user.City)
				case "Company":
					assert.RequireNotNil(t, user.Company, "Expected Company to not be nil.")
					assert.Equal(t, expected, *user.Company, "Expected Company to be %v, but got %s", expected, *user.Company)
				case "Country":
					assert.RequireNotNil(t, user.Country, "Expected Country to not be nil.")
					assert.Equal(t, expected, *user.Country, "Expected Country to be %v, but got %s", expected, *user.Country)
				case "Email":
					assert.RequireNotNil(t, user.Email, "Expected Email to not be nil.")
					assert.Equal(t, expected, *user.Email, "Expected Email to be %v, but got %s", expected, *user.Email)
				case "FullName":
					assert.RequireNotNil(t, user.FullName, "Expected FullName to not be nil.")
					assert.Equal(t, expected, *user.FullName, "Expected FullName to be %v, but got %s", expected, *user.FullName)
				case "ID":
					assert.RequireNotNil(t, user.ID, "Expected ID to not be nil.")
					assert.Equal(t, expected, *user.ID, "Expected ID to be %v, but got %d", expected, user.ID)
				case "PhoneNumber":
					assert.RequireNotNil(t, user.PhoneNumber, "Expected PhoneNumber to not be nil.")
					assert.Equal(t, expected, *user.PhoneNumber, "Expected PhoneNumber to be %v, but got %s", expected, *user.PhoneNumber)
				case "PostalCode":
					assert.RequireNotNil(t, user.PostalCode, "Expected PostalCode to not be nil.")
					assert.Equal(t, expected, *user.PostalCode, "Expected PostalCode to be %v, but got %s", expected, *user.PostalCode)
				case "RegistrationSent":
					assert.RequireNotNil(t, user.RegistrationSent, "Expected RegistrationSent to not be nil.")
					assert.Equal(t, expected, *user.RegistrationSent, "Expected RegistrationSent to be %v, but got %v", expected, *user.RegistrationSent)
				case "Role":
					assert.RequireNotNil(t, user.Role, "Expected Role to not be nil.")
					assert.Equal(t, expected, *user.Role, "Expected Role to be %v, but got %s", expected, *user.Role)
				case "StateOrProvince":
					assert.RequireNotNil(t, user.StateOrProvince, "Expected StateOrProvince to not be nil.")
					assert.Equal(t, expected, *user.StateOrProvince, "Expected StateOrProvince to be %v, but got %s", expected, *user.StateOrProvince)
				case "Tenant":
					assert.RequireNotNil(t, user.Tenant, "Expected Tenant to not be nil.")
					assert.Equal(t, expected, *user.Tenant, "Expected Tenant to be %v, but got %s", expected, *user.Tenant)
				case "TenantID":
					assert.RequireNotNil(t, user.TenantID, "Expected Tenant to not be nil.")
					assert.Equal(t, expected, *user.TenantID, "Expected TenantID to be %v, but got %d", expected, *user.TenantID)
				case "Username":
					assert.RequireNotNil(t, user.Username, "Expected Username to not be nil.")
					assert.Equal(t, expected, *user.Username, "Expected Username to be %v, but got %s", expected, *user.Username)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateTenants(expectedTenants map[string]bool) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Users response to not be nil.")
		userResp := resp.([]tc.User)

		for _, user := range userResp {
			for tenant, expected := range expectedTenants {
				assert.RequireNotNil(t, user.Tenant, "Expected Users response to not be nil.")
				if *user.Tenant == tenant && !expected {
					t.Errorf("Tenant: %s was not expected", *user.Tenant)
				}
			}
		}
	}
}

func validateUsersUpdateCreateFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Users response to not be nil.")
		assert.RequireNotEqual(t, resp.(tc.User), tc.User{}, "Expected a non empty response.")
		userResp := resp.(tc.User)
		users := []tc.User{userResp}
		validateUsersFields(expectedResp)(t, toclientlib.ReqInf{}, users, tc.Alerts{}, nil)
	}
}

func validateUsersSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Users response to not be nil.")
		var usernames []string
		usersResp := resp.([]tc.User)
		for _, user := range usersResp {
			assert.RequireNotNil(t, user.Username, "Expected Username to not be nil.")
			usernames = append(usernames, *user.Username)
		}
		assert.Equal(t, true, sort.StringsAreSorted(usernames), "List is not sorted by their usernames: %v", usernames)
	}
}

func GetUserID(t *testing.T, username string) func() int {
	return func() int {
		users, _, err := TOSession.GetUserByUsernameWithHdr(username, nil)
		assert.RequireNoError(t, err, "Get Users Request failed with error:", err)
		assert.RequireEqual(t, 1, len(users), "Expected response object length 1, but got %d", len(users))
		assert.RequireNotNil(t, users[0].ID, "Expected ID to not be nil.")
		return *users[0].ID
	}
}

func CreateTestUsers(t *testing.T) {
	for _, user := range testData.Users {
		resp, _, err := TOSession.CreateUser(&user)
		assert.RequireNoError(t, err, "Could not create user: %v - alerts: %+v", err, resp.Alerts)
	}
}

// ForceDeleteTestUsers forcibly deletes the users from the db.
// NOTE: Special circumstances!  This should *NOT* be done without a really good reason!
// Connects directly to the DB to remove users rather than going through the client.
// This is required here because the DeleteUser action does not really delete users,  but disables them.
func ForceDeleteTestUsers(t *testing.T) {

	db, err := OpenConnection()
	assert.RequireNoError(t, err, "Cannot open db")
	defer db.Close()

	var usernames []string
	for _, user := range testData.Users {
		usernames = append(usernames, `'`+*user.Username+`'`)
	}

	// there is a constraint that prevents users from being deleted when they have a log
	q := `DELETE FROM log WHERE NOT tm_user = (SELECT id FROM tm_user WHERE username = 'admin')`
	err = execSQL(db, q)
	assert.RequireNoError(t, err, "Cannot execute SQL: %v; SQL is %s", err, q)

	q = `DELETE FROM tm_user WHERE username IN (` + strings.Join(usernames, ",") + `)`
	err = execSQL(db, q)
	assert.NoError(t, err, "Cannot execute SQL: %v; SQL is %s", err, q)
}
