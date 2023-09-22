package v5

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
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestUsers(t *testing.T) {
	WithObjs(t, []TCObj{Tenants, Parameters, Users}, func() {

		opsUserSession := utils.CreateV5Session(t, Config.TrafficOps.URL, "opsuser", "pa$$word", Config.Default.Session.TimeoutInSecs)
		tenant4UserSession := utils.CreateV5Session(t, Config.TrafficOps.URL, "tenant4user", "pa$$word", Config.Default.Session.TimeoutInSecs)

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.UserV4]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateUsersSort()),
				},
				"ADMIN can view CHILD TENANT": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"tenant": {"tenant4"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateUsersFields(map[string]interface{}{"Tenant": "tenant4"})),
				},
				"EMPTY RESPONSE when CHILD TENANT reads PARENT TENANT": {
					ClientSession: tenant4UserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"tenant": {"tenant3"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
			},
			"POST": {
				"FORBIDDEN when CHILD TENANT creates USER with PARENT TENANCY": {
					ClientSession: tenant4UserSession,
					RequestBody: tc.UserV4{
						Email:         util.Ptr("outsidetenancy@example.com"),
						FullName:      util.Ptr("Outside Tenancy"),
						LocalPassword: util.Ptr("pa$$word"),
						Role:          "operations",
						TenantID:      GetTenantID(t, "tenant3")(),
						Username:      "outsideTenantUser",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetUserID(t, "steering"),
					ClientSession: TOSession,
					RequestBody: tc.UserV4{
						AddressLine1:  util.Ptr("updated line 1"),
						AddressLine2:  util.Ptr("updated line 2"),
						City:          util.Ptr("updated city name"),
						Company:       util.Ptr("new company"),
						Country:       util.Ptr("US"),
						Email:         util.Ptr("steeringupdated@example.com"),
						FullName:      util.Ptr("Steering User Updated"),
						LocalPassword: util.Ptr("pa$$word"),
						NewUser:       false,
						Role:          "steering",
						Tenant:        util.Ptr("root"),
						TenantID:      GetTenantID(t, "root")(),
						Username:      "steering",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateUsersUpdateCreateFields(map[string]interface{}{"AddressLine1": "updated line 1",
							"AddressLine2": "updated line 2", "City": "updated city name", "Company": "new company",
							"Country": "US", "Email": "steeringupdated@example.com", "FullName": "Steering User Updated"})),
				},
				"OK when UPDATING SELF": {
					EndpointID:    GetUserID(t, "opsuser"),
					ClientSession: opsUserSession,
					RequestBody: tc.UserV4{
						AddressLine1:  util.Ptr("address of ops"),
						AddressLine2:  util.Ptr("place"),
						City:          util.Ptr("somewhere"),
						Company:       util.Ptr("else"),
						Country:       util.Ptr("UK"),
						Email:         util.Ptr("ops-updated@example.com"),
						FullName:      util.Ptr("Operations User Updated"),
						LocalPassword: util.Ptr("pa$$word"),
						Role:          "operations",
						Tenant:        util.Ptr("root"),
						TenantID:      GetTenantID(t, "root")(),
						Username:      "opsuser",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateUsersUpdateCreateFields(map[string]interface{}{"Email": "ops-updated@example.com", "FullName": "Operations User Updated"})),
				},
				"NOT FOUND when UPDATING SELF with ROLE that DOESNT EXIST": {
					EndpointID:    GetUserID(t, "opsuser"),
					ClientSession: opsUserSession,
					RequestBody: tc.UserV4{
						AddressLine1:  util.Ptr("address of ops"),
						AddressLine2:  util.Ptr("place"),
						City:          util.Ptr("somewhere"),
						Company:       util.Ptr("else"),
						Country:       util.Ptr("UK"),
						Email:         util.Ptr("ops-updated@example.com"),
						FullName:      util.Ptr("Operations User Updated"),
						LocalPassword: util.Ptr("pa$$word"),
						Role:          "operations_updated",
						Tenant:        util.Ptr("root"),
						TenantID:      GetTenantID(t, "root")(),
						Username:      "opsuser",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"FORBIDDEN when OPERATIONS USER updates ADMIN USER": {
					EndpointID:    GetUserID(t, "admin"),
					ClientSession: opsUserSession,
					RequestBody: tc.UserV4{
						Email:         util.Ptr("oops@ops.net"),
						FullName:      util.Ptr("oops"),
						LocalPassword: util.Ptr("pa$$word"),
						Role:          "admin",
						Tenant:        util.Ptr("root"),
						TenantID:      GetTenantID(t, "root")(),
						Username:      "admin",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"FORBIDDEN when CHILD TENANT USER updates PARENT TENANT USER": {
					EndpointID:    GetUserID(t, "tenant3user"),
					ClientSession: tenant4UserSession,
					RequestBody: tc.UserV4{
						Email:         util.Ptr("tenant3user@example.com"),
						FullName:      util.Ptr("Parent tenant test"),
						LocalPassword: util.Ptr("pa$$word"),
						Role:          "admin",
						Tenant:        util.Ptr("tenant2"),
						TenantID:      GetTenantID(t, "tenant2")(),
						Username:      "tenant3user",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetUsers(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateUser(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateUser(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
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
		userResp := resp.([]tc.UserV4)
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
					assert.Equal(t, expected, user.Role, "Expected Role to be %v, but got %s", expected, user.Role)
				case "StateOrProvince":
					assert.RequireNotNil(t, user.StateOrProvince, "Expected StateOrProvince to not be nil.")
					assert.Equal(t, expected, *user.StateOrProvince, "Expected StateOrProvince to be %v, but got %s", expected, *user.StateOrProvince)
				case "Tenant":
					assert.RequireNotNil(t, user.Tenant, "Expected Tenant to not be nil.")
					assert.Equal(t, expected, *user.Tenant, "Expected Tenant to be %v, but got %s", expected, *user.Tenant)
				case "TenantID":
					assert.Equal(t, expected, user.TenantID, "Expected TenantID to be %v, but got %d", expected, user.TenantID)
				case "Username":
					assert.Equal(t, expected, user.Username, "Expected Username to be %v, but got %s", expected, user.Username)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateUsersUpdateCreateFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Users response to not be nil.")
		assert.RequireNotEqual(t, resp.(tc.UserV4), tc.UserV4{}, "Expected a non empty response.")
		userResp := resp.(tc.UserV4)
		users := []tc.UserV4{userResp}
		validateUsersFields(expectedResp)(t, toclientlib.ReqInf{}, users, tc.Alerts{}, nil)
	}
}

func validateUsersSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Users response to not be nil.")
		var usernames []string
		usersResp := resp.([]tc.UserV4)
		for _, user := range usersResp {
			usernames = append(usernames, user.Username)
		}
		assert.Equal(t, true, sort.StringsAreSorted(usernames), "List is not sorted by their usernames: %v", usernames)
	}
}

func GetUserID(t *testing.T, username string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("username", username)
		users, _, err := TOSession.GetUsers(opts)
		assert.RequireNoError(t, err, "Get Users Request failed with error:", err)
		assert.RequireEqual(t, 1, len(users.Response), "Expected response object length 1, but got %d", len(users.Response))
		assert.RequireNotNil(t, users.Response[0].ID, "Expected ID to not be nil.")
		return *users.Response[0].ID
	}
}

func CreateTestUsers(t *testing.T) {
	for _, user := range testData.Users {
		resp, _, err := TOSession.CreateUser(user, client.RequestOptions{})
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
		usernames = append(usernames, `'`+user.Username+`'`)
	}

	// there is a constraint that prevents users from being deleted when they have a log
	q := `DELETE FROM log WHERE NOT tm_user = (SELECT id FROM tm_user WHERE username = 'admin')`
	err = execSQL(db, q)
	assert.RequireNoError(t, err, "Cannot execute SQL: %v; SQL is %s", err, q)

	q = `DELETE FROM tm_user WHERE username IN (` + strings.Join(usernames, ",") + `)`
	err = execSQL(db, q)
	assert.NoError(t, err, "Cannot execute SQL: %v; SQL is %s", err, q)
}

// ForceDeleteTestUsersByUsernames forcibly deletes the users passed in from a slice of usernames from the db.
// NOTE: Special circumstances!  This should *NOT* be done without a really good reason!
// Connects directly to the DB to remove users rather than going through the client.
// This is required here because the DeleteUser action does not really delete users, but disables them.
func ForceDeleteTestUsersByUsernames(t *testing.T, usernames []string) {

	db, err := OpenConnection()
	assert.RequireNoError(t, err, "Cannot open db")
	defer db.Close()

	for i, u := range usernames {
		usernames[i] = `'` + u + `'`
	}
	// there is a constraint that prevents users from being deleted when they have a log
	q := `DELETE FROM log WHERE NOT tm_user = (SELECT id FROM tm_user WHERE username = 'admin')`
	err = execSQL(db, q)
	assert.RequireNoError(t, err, "Cannot execute SQL: %s; SQL is %s", err, q)

	q = `DELETE FROM tm_user WHERE username IN (` + strings.Join(usernames, ",") + `)`
	err = execSQL(db, q)
	assert.NoError(t, err, "Cannot execute SQL: %s; SQL is %s", err, q)
}
