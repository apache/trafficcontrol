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
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestUserCurrent(t *testing.T) {
	WithObjs(t, []TCObj{Tenants, Parameters, Users}, func() {

		opsUserSession := utils.CreateV3Session(t, Config.TrafficOps.URL, "opsuser", "pa$$word", Config.Default.Session.TimeoutInSecs)

		methodTests := utils.V3TestCase{
			"GET": {
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateUserCurrentFields(map[string]interface{}{"Username": "admin"})),
				},
			},
			"PUT": {
				"OK when VALID request": {
					ClientSession: opsUserSession,
					RequestBody: map[string]interface{}{
						"addressLine1":       "address of ops",
						"addressLine2":       "place",
						"city":               "somewhere",
						"company":            "else",
						"country":            "UK",
						"email":              "ops-updated@example.com",
						"fullName":           "Operations User Updated",
						"id":                 GetUserID(t, "opsuser")(),
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
				"BAD REQUEST when EMPTY EMAIL field": {
					ClientSession: opsUserSession,
					RequestBody: map[string]interface{}{
						"addressLine1":       "address of ops",
						"addressLine2":       "place",
						"city":               "somewhere",
						"company":            "else",
						"country":            "UK",
						"email":              "",
						"fullName":           "Operations User Updated",
						"id":                 GetUserID(t, "opsuser")(),
						"localPasswd":        "pa$$word",
						"confirmLocalPasswd": "pa$$word",
						"role":               3,
						"tenant":             "root",
						"tenantId":           GetTenantID(t, "root")(),
						"username":           "opsuser",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
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
							resp, reqInf, err := testCase.ClientSession.GetUserCurrentWithHdr(testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateCurrentUser(user)
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

func validateUserCurrentFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Users response to not be nil.")
		user := resp.(*tc.UserCurrent)
		for field, expected := range expectedResp {
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
				assert.RequireNotNil(t, user.UserName, "Expected Username to not be nil.")
				assert.Equal(t, expected, *user.UserName, "Expected Username to be %v, but got %s", expected, *user.UserName)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}
