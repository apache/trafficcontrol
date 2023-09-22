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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestServersHostnameUpdate(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers}, func() {

		// Postgres stores microsecond precision. There is also some discussion around MacOS losing
		// precision as well. The nanosecond precision is accurate within go one linux however,
		// but round trips to and from the database may result in an inaccurate Equals comparison
		// with the loss of precision. Also, it appears to Round and not Truncate.
		now := time.Now().Round(time.Microsecond)

		methodTests := utils.V4TestCase{
			"POST": {
				"OK when VALID CONFIG_APPLY_TIME PARAMETER": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"hostName": {"atlanta-edge-01"}}},
					RequestBody: map[string]interface{}{
						"config_apply_time": util.TimePtr(now),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServerApplyTimes("atlanta-edge-01", map[string]interface{}{"ConfigApplyTime": now})),
				},
				"OK when VALID REVALIDATE_APPLY_TIME PARAMETER": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"hostName": {"cdn2-test-edge"}}},
					RequestBody: map[string]interface{}{
						"revalidate_apply_time": util.TimePtr(now),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServerApplyTimes("cdn2-test-edge", map[string]interface{}{"RevalApplyTime": now})),
				},
				"BAD REQUEST when UPDATED AND CONFIG_APPLY_TIME PARAMETER": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"hostName": {"atlanta-edge-01"}, "updated": {"true"}}},
					RequestBody: map[string]interface{}{
						"config_apply_time": util.TimePtr(now),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when REVAL_UPDATED AND REVALIDATE_APPLY_TIME PARAMETER": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"hostName": {"atlanta-edge-01"}, "reval_updated": {"true"}}},
					RequestBody: map[string]interface{}{
						"revalidate_apply_time": util.TimePtr(now),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var hostName string
					var configApplyTime *time.Time
					var revalApplyTime *time.Time

					if hostNameParam, ok := testCase.RequestOpts.QueryParameters["hostName"]; ok {
						hostName = hostNameParam[0]
					}

					if configApplyTimeVal, ok := testCase.RequestBody["config_apply_time"]; ok {
						configApplyTime = configApplyTimeVal.(*time.Time)
					}

					if revalApplyTimeVal, ok := testCase.RequestBody["revalidate_apply_time"]; ok {
						revalApplyTime = revalApplyTimeVal.(*time.Time)
					}

					switch method {
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.SetUpdateServerStatusTimes(hostName, configApplyTime, revalApplyTime, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func validateServerApplyTimes(hostName string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, _ tc.Alerts, _ error) {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Add("hostName", hostName)
		resp, _, err := TOSession.GetServers(opts)
		assert.RequireNoError(t, err, "Cannot GET Server by name '%s': %v - alerts: %+v", hostName, err, resp.Alerts)
		assert.RequireEqual(t, 1, len(resp.Response), "GET Server expected 1, actual %v", len(resp.Response))
		assert.RequireNotNil(t, resp.Response[0].UpdPending, "Server '%s' had nil UpdPending after update status change", hostName)
		assert.RequireNotNil(t, resp.Response[0].RevalPending, "Server '%s' had nil RevalPending after update status change", hostName)

		for field, expected := range expectedResp {
			for _, server := range resp.Response {
				switch field {
				case "ConfigApplyTime":
					assert.RequireNotNil(t, resp.Response[0].ConfigApplyTime, "Expected ConfigApplyTime to not be nil.")
					assert.Equal(t, true, server.ConfigApplyTime.Equal(expected.(time.Time)), "Expected ConfigApplyTime to be %v, but got %v", expected, server.ConfigApplyTime)
				case "RevalApplyTime":
					assert.RequireNotNil(t, resp.Response[0].RevalApplyTime, "Expected RevalApplyTime to not be nil.")
					assert.Equal(t, true, server.RevalApplyTime.Equal(expected.(time.Time)), "Expected RevalApplyTime to be %v, but got %v", expected, server.RevalApplyTime)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}
