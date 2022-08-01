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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
)

func TestUserCurrent(t *testing.T) {
	WithObjs(t, []TCObj{Tenants, Parameters, Users}, func() {

		opsUserSession := utils.CreateV3Session(t, Config.TrafficOps.URL, "opsuser", "pa$$word", Config.Default.Session.TimeoutInSecs)

		methodTests := utils.V3TestCase{
			"GET": {
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateUsersUpdateCreateFields(map[string]interface{}{"Username": "admin"})),
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
						"localPasswd":        "pa$$word",
						"confirmLocalPasswd": "pa$$word",
						"role":               "operations",
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
						"localPasswd":        "pa$$word",
						"confirmLocalPasswd": "pa$$word",
						"role":               "operations",
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
