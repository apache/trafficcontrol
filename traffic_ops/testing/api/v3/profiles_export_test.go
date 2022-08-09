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
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

func TestProfilesExport(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, ProfileParameters}, func() {

		methodTests := utils.V3TestCase{
			"GET": {
				"OK when VALID request": {
					EndpointId:    GetProfileID(t, "EDGE1"),
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateProfilesExportFields(map[string]interface{}{"ProfileCDNName": "cdn1", "ProfileName": "EDGE1",
							"ProfileDescription": "edge1 description", "ProfileType": "ATS_PROFILE"})),
				},
				"NOT FOUND when PROFILE DOESNT EXIST": {
					EndpointId:    func() int { return 1111111111 },
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
							resp, reqInf, err := testCase.ClientSession.ExportProfile(testCase.EndpointId())
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, resp.Alerts, err)
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
			switch field {
			case "ProfileCDNName":
				assert.RequireNotNil(t, profileExport.Profile.CDNName, "Expected Profile CDNName to not be nil.")
				assert.Equal(t, expected, *profileExport.Profile.CDNName, "Expected Profile.CDNName to be %v, but got %d", expected, *profileExport.Profile.CDNName)
			case "ProfileDescription":
				assert.RequireNotNil(t, profileExport.Profile.Description, "Expected Profile Description to not be nil.")
				assert.Equal(t, expected, *profileExport.Profile.Description, "Expected Profile.Description to be %v, but got %d", expected, *profileExport.Profile.Description)
			case "ProfileName":
				assert.RequireNotNil(t, profileExport.Profile.Name, "Expected Profile Name to not be nil.")
				assert.Equal(t, expected, *profileExport.Profile.Name, "Expected Profile.Name to be %v, but got %d", expected, *profileExport.Profile.Name)
			case "ProfileType":
				assert.RequireNotNil(t, profileExport.Profile.Type, "Expected Profile Type to not be nil.")
				assert.Equal(t, expected, *profileExport.Profile.Type, "Expected Profile.Type to be %v, but got %d", expected, *profileExport.Profile.Type)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}
