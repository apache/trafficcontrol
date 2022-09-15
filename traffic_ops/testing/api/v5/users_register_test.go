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
	"encoding/json"
	"net/http"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

func TestUsersRegister(t *testing.T) {
	if includeSystemTests {
		WithObjs(t, []TCObj{Tenants, Parameters}, func() {

			methodTests := utils.V5TestCase{
				"POST": {
					"OK when VALID request": {
						ClientSession: TOSession,
						RequestBody: map[string]interface{}{
							"addressLine1":       "address of ops",
							"addressLine2":       "place",
							"city":               "somewhere",
							"company":            "else",
							"country":            "UK",
							"email":              "opsupdated@example.com",
							"fullName":           "Operations User Updated",
							"localPasswd":        "pa$$word",
							"confirmLocalPasswd": "pa$$word",
							"role":               "operations",
							"tenant":             "root",
							"tenantId":           GetTenantID(t, "root")(),
							"username":           "opsuser",
						},
						Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDeletion("opsupdated@example.com")),
					},
				},
			}

			for method, testCases := range methodTests {
				t.Run(method, func(t *testing.T) {
					for name, testCase := range testCases {
						userRegistration := tc.UserRegistrationRequestV4{}

						if testCase.RequestBody != nil {
							dat, err := json.Marshal(testCase.RequestBody)
							assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
							err = json.Unmarshal(dat, &userRegistration)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						}

						switch method {
						case "POST":
							t.Run(name, func(t *testing.T) {
								alerts, reqInf, err := testCase.ClientSession.RegisterNewUser(userRegistration.TenantID, userRegistration.Role, userRegistration.Email, testCase.RequestOpts)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, alerts, err)
								}
							})
						}
					}
				})
			}
		})
	}
}

func validateDeletion(email string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		db, err := OpenConnection()
		assert.RequireNoError(t, err, "Cannot open db")
		defer db.Close()
		q := `DELETE FROM tm_user WHERE email = '` + email + `'`
		err = execSQL(db, q)
		assert.NoError(t, err, "Cannot execute SQL to delete registered users: %s; SQL is %s", err, q)
	}
}
