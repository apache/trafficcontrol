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
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestUserCurrent(t *testing.T) {
	WithObjs(t, []TCObj{Tenants, Parameters, Users}, func() {

		opsUserSession := utils.CreateV4Session(t, Config.TrafficOps.URL, "opsuser", "pa$$word", Config.Default.Session.TimeoutInSecs)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.UserV4]{
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
				"BAD REQUEST when EMPTY EMAIL field": {
					ClientSession: opsUserSession,
					RequestBody: tc.UserV4{
						AddressLine1:  util.Ptr("address of ops"),
						AddressLine2:  util.Ptr("place"),
						City:          util.Ptr("somewhere"),
						Company:       util.Ptr("else"),
						Country:       util.Ptr("UK"),
						Email:         util.Ptr(""),
						FullName:      util.Ptr("Operations User Updated"),
						LocalPassword: util.Ptr("pa$$word"),
						Role:          "operations",
						Tenant:        util.Ptr("root"),
						TenantID:      GetTenantID(t, "root")(),
						Username:      "opsuser",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetUserCurrent(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateCurrentUser(testCase.RequestBody, testCase.RequestOpts)
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
