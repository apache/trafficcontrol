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
	"net/url"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	totest "github.com/apache/trafficcontrol/v8/lib/go-tc/totestv4"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestSnapshot(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments}, func() {

		readOnlyUserSession := utils.CreateV4Session(t, Config.TrafficOps.URL, "readonlyuser", "pa$$word", Config.Default.Session.TimeoutInSecs)

		methodTests := utils.V4TestCase{
			"PUT": {
				"VERIFY ANYMAP DELIVERY SERVICE is NOT in CRCONFIG": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateDeliveryServiceNotInCRConfig("cdn1", "anymap-ds"),
						validateCRConfigFields("cdn1", map[string]interface{}{"TMPath": (*string)(nil), "TMHost": "crconfig.tm.url.test.invalid"})),
				},
				"OK when VALID CDN parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID CDNID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdnID": {strconv.Itoa(totest.GetCDNID(t, TOSession, "cdn1")())}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"NOT FOUND when NON-EXISTENT CDN": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn-invalid"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"NOT FOUND when NON-EXISTENT CDNID": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdnID": {"999999"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"FORBIDDEN when READ-ONLY user": {
					ClientSession: readOnlyUserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn1"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.SnapshotCRConfig(testCase.RequestOpts)
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

func validateCRConfigFields(cdn string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, _ tc.Alerts, _ error) {
		snapshotResp, _, err := TOSession.GetCRConfig(cdn, client.RequestOptions{})
		assert.RequireNoError(t, err, "Unexpected error retrieving Snapshot of CDN '%s': %v - alerts: %+v", cdn, err, snapshotResp.Alerts)
		crconfig := snapshotResp.Response

		for field, expected := range expectedResp {
			switch field {
			case "TMPath":
				assert.Equal(t, expected, crconfig.Stats.TMPath, "Expected no TMPath in APIv4, but it was: %v", crconfig.Stats.TMPath)
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
		snapshotResp, _, err := TOSession.GetCRConfig(cdn, client.RequestOptions{})
		assert.RequireNoError(t, err, "Unexpected error retrieving Snapshot of CDN '%s': %v - alerts: %+v", cdn, err, snapshotResp.Alerts)

		for ds := range snapshotResp.Response.DeliveryServices {
			assert.NotEqual(t, ds, deliveryService, "Found unexpected delivery service: %s in CRConfig Delivery Services.", deliveryService)
		}

		for _, server := range snapshotResp.Response.ContentServers {
			for ds := range server.DeliveryServices {
				assert.NotEqual(t, ds, deliveryService, "Found unexpected delivery service: %s in CRConfig Content Servers Delivery Services.", deliveryService)
			}
		}
	}
}
