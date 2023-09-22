package v4

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
	"net/mail"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestUsersRegister(t *testing.T) {
	if includeSystemTests {
		WithObjs(t, []TCObj{Tenants, Parameters}, func() {

			methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.UserRegistrationRequestV4]{
				"POST": {
					"OK when VALID request": {
						ClientSession: TOSession,
						RequestBody: tc.UserRegistrationRequestV4{
							Email:    rfc.EmailAddress{Address: mail.Address{Address: "opsupdated@example.com"}},
							Role:     "operations",
							TenantID: uint(GetTenantID(t, "root")()),
						},
						Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDeletion("opsupdated@example.com")),
					},
				},
			}

			for method, testCases := range methodTests {
				t.Run(method, func(t *testing.T) {
					for name, testCase := range testCases {
						switch method {
						case "POST":
							t.Run(name, func(t *testing.T) {
								userRegistration := testCase.RequestBody
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
