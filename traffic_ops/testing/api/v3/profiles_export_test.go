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

package v3

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestProfilesExport(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, ProfileParameters}, func() {

		methodTests := utils.V3TestCaseT[struct{}]{
			"GET": {
				"OK when VALID request": {
					EndpointID:    GetProfileID(t, "EDGE1"),
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateProfilesExportFields(map[string]interface{}{"CDNName": "cdn1", "Name": "EDGE1",
							"Description": "edge1 description", "Type": "ATS_PROFILE"})),
				},
				"NOT FOUND when PROFILE DOESNT EXIST": {
					EndpointID:    func() int { return 1111111111 },
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.ExportProfile(testCase.EndpointID())
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					}
				}
			})
		}

	})
}

func validateProfilesExportFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Profiles Export response to not be nil.")
		profileExport := resp.(*tc.ProfileExportResponse)
		for field, expected := range expectedResp {
			fieldValue := reflect.Indirect(reflect.ValueOf(profileExport.Profile).FieldByName(field)).String()
			assert.RequireNotNil(t, fieldValue, "Expected %s to not be nil.", field)
			assert.Equal(t, expected, fieldValue, "Expected %s to be %v, but got %s", field, expected, fieldValue)
		}
	}
}
