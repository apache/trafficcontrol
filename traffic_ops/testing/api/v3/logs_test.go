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
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
)

func TestLogs(t *testing.T) {
	WithObjs(t, []TCObj{Roles, Tenants, Users}, func() { // Objs added to create logs when this test is run alone

		methodTests := utils.V3TestCaseT[struct{}]{
			"GET": {
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"RESPONSE LENGTH matches LIMIT parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"limit": {"10"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(10)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							if name == "RESPONSE LENGTH matches LIMIT parameter" {
								limit, _ := strconv.Atoi(testCase.RequestParams["limit"][0])
								resp, reqInf, err := testCase.ClientSession.GetLogsByLimit(limit)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else {
								resp, reqInf, err := testCase.ClientSession.GetLogs()
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							}
						})
					}
				}
			})
		}
	})
}
