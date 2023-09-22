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
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestSnapshot(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments}, func() {

		methodTests := utils.V3TestCase{
			"PUT": {
				"VERIFY ANYMAP DELIVERY SERVICE is NOT in CRCONFIG": {
					ClientSession: TOSession,
					RequestParams: url.Values{"cdn": {"cdn1"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateDeliveryServiceNotInCRConfig("cdn1", "anymap-ds"),
						validateCRConfigFields("cdn1", map[string]interface{}{"TMPath": "", "TMHost": "crconfig.tm.url.test.invalid"})),
				},
				"OK when VALID CDN parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"cdn": {"cdn1"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID CDNID parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"cdnID": {strconv.Itoa(GetCDNID(t, "cdn1")())}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"NOT FOUND when NON-EXISTENT CDN": {
					ClientSession: TOSession,
					RequestParams: url.Values{"cdn": {"cdn-invalid"}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"NOT FOUND when NON-EXISTENT CDNID": {
					ClientSession: TOSession,
					RequestParams: url.Values{"cdnID": {"999999"}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "PUT":
						t.Run(name, func(t *testing.T) {
							if name == "OK when VALID CDNID parameter" {
								var cdnID int
								var err error
								if val, ok := testCase.RequestParams["cdnID"]; ok {
									cdnID, err = strconv.Atoi(val[0])
									assert.NoError(t, err, "Error converting string to integer: %v", err)
								}
								resp, reqInf, err := testCase.ClientSession.SnapshotCRConfigByID(cdnID)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else {
								var cdn string
								if val, ok := testCase.RequestParams["cdn"]; ok {
									cdn = val[0]
								}
								reqInf, err := testCase.ClientSession.SnapshotCRConfigWithHdr(cdn, testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, nil, tc.Alerts{}, err)
								}
							}
						})
					}
				}
			})
		}
	})
}

func validateCRConfigFields(cdn string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, _ tc.Alerts, _ error) {
		var crconfig tc.CRConfig
		snapshotResp, _, err := TOSession.GetCRConfig(cdn)
		assert.RequireNoError(t, err, "Unexpected error retrieving Snapshot of CDN '%s': %v", cdn, err)
		err = json.Unmarshal(snapshotResp, &crconfig)
		assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)

		for field, expected := range expectedResp {
			switch field {
			case "TMPath":
				assert.RequireNotNil(t, crconfig.Stats.TMPath, "Expected Stats TM Path to not be nil.")
			case "TMHost":
				assert.RequireNotNil(t, crconfig.Stats.TMHost, "Expected Stats TM Host to not be nil.")
				assert.Equal(t, expected, *crconfig.Stats.TMHost, "Expected Stats TM Host to be %v, but got %s", expected, *crconfig.Stats.TMHost)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}

func validateDeliveryServiceNotInCRConfig(cdn string, deliveryService string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, _ tc.Alerts, _ error) {
		var crconfig tc.CRConfig
		snapshotResp, _, err := TOSession.GetCRConfig(cdn)
		assert.RequireNoError(t, err, "Unexpected error retrieving Snapshot of CDN '%s': %v", cdn, err)
		err = json.Unmarshal(snapshotResp, &crconfig)
		assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)

		for ds := range crconfig.DeliveryServices {
			assert.NotEqual(t, ds, deliveryService, "Found unexpected delivery service: %s in CRConfig Delivery Services.", deliveryService)
		}

		for _, server := range crconfig.ContentServers {
			for ds := range server.DeliveryServices {
				assert.NotEqual(t, ds, deliveryService, "Found unexpected delivery service: %s in CRConfig Content Servers Delivery Services.", deliveryService)
			}
		}
	}
}
