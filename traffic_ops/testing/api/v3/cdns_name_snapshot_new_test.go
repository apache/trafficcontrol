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
	"net/url"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

var baselineCRConfig tc.CRConfig

func TestCDNNameSnapshotNew(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices}, func() {

		methodTests := utils.V3TestCase{
			"GET": {
				"VERIFY SNAPSHOT UPDATE CAPTURED CORRECTLY": {
					ClientSession: TOSession,
					RequestParams: url.Values{"cdn": {"cdn1"}},
					PreReqFuncs:   []func(){getBaselineCRConfig(t, "cdn1"), deleteParameter(t, "tm.url")},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateCRConfigNewFields("cdn1", map[string]interface{}{"TMHost": ""}), validateDeliveryServicesUnchanged()),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							var cdn string
							if val, ok := testCase.RequestParams["cdn"]; ok {
								cdn = val[0]
							}
							for _, prerequisite := range testCase.PreReqFuncs {
								prerequisite()
							}
							resp, reqInf, err := testCase.ClientSession.GetCRConfigNew(cdn)
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

func validateCRConfigNewFields(cdn string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected CRConfigNew response to not be nil.")
		var newCRConfig tc.CRConfig
		err := json.Unmarshal(resp.([]byte), &newCRConfig)
		assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)

		for field, expected := range expectedResp {
			switch field {
			case "TMPath":
				assert.RequireNotNil(t, newCRConfig.Stats.TMPath, "Expected Stats TM Path to not be nil.")
			case "TMHost":
				assert.RequireNotNil(t, newCRConfig.Stats.TMHost, "Expected Stats TM Host to not be nil.")
				assert.Equal(t, expected, *newCRConfig.Stats.TMHost, "Expected Stats TM Host to be %v, but got %s", expected, *newCRConfig.Stats.TMHost)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}

func validateDeliveryServicesUnchanged() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected new snapshot response to not be nil.")
		var newCRConfig tc.CRConfig
		err := json.Unmarshal(resp.([]byte), &newCRConfig)
		assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
		assert.Exactly(t, newCRConfig.DeliveryServices, baselineCRConfig.DeliveryServices, "Expected Delivery Services to be unchanged.")
	}
}

func getBaselineCRConfig(t *testing.T, cdn string) func() {
	return func() {
		_, err := TOSession.SnapshotCRConfigWithHdr(cdn, nil)
		assert.RequireNoError(t, err, "Unexpected error taking Snapshot of CDN '%s': %v", cdn, err)
		getCRConfig, _, err := TOSession.GetCRConfig(cdn)
		assert.RequireNoError(t, err, "Unexpected error retrieving Snapshot of CDN '%s': %v", cdn, err)
		err = json.Unmarshal(getCRConfig, &baselineCRConfig)
		assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
	}
}

func deleteParameter(t *testing.T, paramName string) func() {
	return func() {
		paramResp, _, err := TOSession.GetParameterByNameWithHdr(paramName, nil)
		assert.RequireNoError(t, err, "Cannot get Parameter by name '%s': %v", paramName, err)
		assert.RequireGreaterOrEqual(t, len(paramResp), 1, "Expected at least one parameter to be returned.")
		delResp, _, err := TOSession.DeleteParameterByID(paramResp[0].ID)
		assert.RequireNoError(t, err, "Cannot DELETE Parameter by name: %v - %v", err, delResp)
	}
}
