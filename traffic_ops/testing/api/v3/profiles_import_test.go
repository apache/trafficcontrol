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
	"encoding/json"
	"net/http"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

func TestProfilesImport(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, ProfileParameters}, func() {

		methodTests := utils.V3TestCase{
			"POST": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"importProfile": map[string]interface{}{
							"profile": map[string]interface{}{
								"name":        "GLOBAL",
								"description": "Global Traffic Ops profile",
								"cdn":         "cdn1",
								"type":        "UNK_PROFILE",
							},
							"parameters": []map[string]interface{}{
								{
									"config_file": "global",
									"name":        "tm.instance_name",
									"value":       "Traffic Ops CDN",
								},
								{
									"config_file": "global",
									"name":        "tm.toolname",
									"value":       "Traffic Ops",
								},
							},
						},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateProfilesImport(map[string]interface{}{"ProfileName": "GLOBAL", "ProfileCDNName": "cdn1",
							"ProfileDescription": "Global Traffic Ops profile", "ProfileType": "UNK_PROFILE"})),
				},
				"BAD REQUEST when SPACE in PROFILE NAME": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"importProfile": map[string]interface{}{
							"profile": map[string]interface{}{
								"name":        "GLOBAL SPACES",
								"description": "Global Traffic Ops profile",
								"cdn":         "cdn1",
								"type":        "UNK_PROFILE",
							},
							"parameters": []map[string]interface{}{
								{
									"config_file": "global",
									"name":        "tm.instance_name",
									"value":       "Traffic Ops CDN",
								},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					importProfile := tc.ProfileImportRequest{}

					if testCase.RequestBody != nil {
						if importProfileBody, ok := testCase.RequestBody["importProfile"]; ok {
							dat, err := json.Marshal(importProfileBody)
							assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
							err = json.Unmarshal(dat, &importProfile)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						}
					}

					switch method {
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.ImportProfile(&importProfile)
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

func validateProfilesImport(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Profiles Export response to not be nil.")
		profileImportResp := resp.(tc.ProfileImportResponseObj)
		profileImport := profileImportResp.ProfileExportImportNullable
		for field, expected := range expectedResp {
			switch field {
			case "ProfileCDNName":
				assert.RequireNotNil(t, profileImport.CDNName, "Expected Profile CDNName to not be nil.")
				assert.Equal(t, expected, *profileImport.CDNName, "Expected Profile CDNName to be %v, but got %d", expected, *profileImport.CDNName)
			case "ProfileDescription":
				assert.RequireNotNil(t, profileImport.Description, "Expected Profile Description to not be nil.")
				assert.Equal(t, expected, *profileImport.Description, "Expected Profile Description to be %v, but got %d", expected, *profileImport.Description)
			case "ProfileName":
				assert.RequireNotNil(t, profileImport.Name, "Expected Profile Name to not be nil.")
				assert.Equal(t, expected, *profileImport.Name, "Expected Profile Name to be %v, but got %d", expected, *profileImport.Name)
			case "ProfileType":
				assert.RequireNotNil(t, profileImport.Type, "Expected Profile Type to not be nil.")
				assert.Equal(t, expected, *profileImport.Type, "Expected Profile Type to be %v, but got %d", expected, *profileImport.Type)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}
