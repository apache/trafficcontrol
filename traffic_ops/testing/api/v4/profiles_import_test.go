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

package v4

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestProfilesImport(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, ProfileParameters}, func() {

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.ProfileImportRequest]{
			"POST": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestBody: tc.ProfileImportRequest{
						Profile: tc.ProfileExportImportNullable{
							Name:        util.Ptr("GLOBAL"),
							Description: util.Ptr("Global Traffic Ops profile"),
							CDNName:     util.Ptr("cdn1"),
							Type:        util.Ptr("UNK_PROFILE"),
						},
						Parameters: []tc.ProfileExportImportParameterNullable{
							{
								ConfigFile: util.Ptr("global"),
								Name:       util.Ptr("tm.instance_name"),
								Value:      util.Ptr("Traffic Ops CDN"),
							},
							{
								ConfigFile: util.Ptr("global"),
								Name:       util.Ptr("tm.toolname"),
								Value:      util.Ptr("Traffic Ops"),
							},
						},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateProfilesImport(map[string]interface{}{"Name": "GLOBAL", "CDNName": "cdn1",
							"Description": "Global Traffic Ops profile", "Type": "UNK_PROFILE"})),
				},
				"BAD REQUEST when SPACE in PROFILE NAME": {
					ClientSession: TOSession,
					RequestBody: tc.ProfileImportRequest{
						Profile: tc.ProfileExportImportNullable{
							Name:        util.Ptr("GLOBAL SPACES"),
							Description: util.Ptr("Global Traffic Ops profile"),
							CDNName:     util.Ptr("cdn1"),
							Type:        util.Ptr("UNK_PROFILE"),
						},
						Parameters: []tc.ProfileExportImportParameterNullable{
							{
								ConfigFile: util.Ptr("global"),
								Name:       util.Ptr("tm.instance_name"),
								Value:      util.Ptr("Traffic Ops CDN"),
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
					switch method {
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.ImportProfile(testCase.RequestBody, testCase.RequestOpts)
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
			fieldValue := reflect.Indirect(reflect.ValueOf(profileImport).FieldByName(field)).String()
			assert.RequireNotNil(t, fieldValue, "Expected %s to not be nil.", field)
			assert.Equal(t, expected, fieldValue, "Expected %s to be %v, but got %s", field, expected, fieldValue)
		}
	}
}
