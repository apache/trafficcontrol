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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestCDNNameConfigsMonitoring(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, ProfileParameters, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices}, func() {

		methodTests := utils.V4TestCase{
			"GET": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateHealthThresholdParameters("EDGE1", map[string]string{"loadavg": "25.0", "availableBandwidthInKbps": ">1750000", "queryTime": "1000"})),
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
							if val, ok := testCase.RequestOpts.QueryParameters["cdn"]; ok {
								cdn = val[0]
							}
							resp, reqInf, err := testCase.ClientSession.GetTrafficMonitorConfig(cdn, testCase.RequestOpts)
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

func validateHealthThresholdParameters(profileName string, healthThresholdParams map[string]string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Traffic Monitor Config response to not be nil.")
		tmConfig := resp.(tc.TrafficMonitorConfig)
		parameterMap := map[string]tc.HealthThreshold{}
		parameterFound := map[string]bool{}
		for parameter, parameterValue := range healthThresholdParams {
			threshold, err := tc.StrToThreshold(parameterValue)
			parameterMap[parameter] = threshold
			assert.RequireNoError(t, err, "Error: converting string '%s' to HealthThreshold: %v", parameterValue, err)
			parameterFound[parameter] = false
		}

		profileFound := false
		var profile tc.TMProfile
		for _, profile = range tmConfig.Profiles {
			if profile.Name == profileName {
				profileFound = true
				break
			}
		}
		assert.RequireEqual(t, true, profileFound, "Traffic Monitor Config contained no Profile named '%s", profileName)

		for parameterName, value := range profile.Parameters.Thresholds {
			_, ok := parameterFound[parameterName]
			assert.Equal(t, true, ok, "Unexpected Threshold Parameter name '%s' found in Profile '%s' in Traffic Monitor Config", parameterName, profileName)
			parameterFound[parameterName] = true
			assert.Equal(t, parameterMap[parameterName].String(), value.String(), "Expected '%s' but received '%s' for Threshold Parameter '%s' in Profile '%s' in Traffic Monitor Config", parameterMap[parameterName].String(), value.String(), parameterName, profileName)
		}
		missingParameters := []string{}
		for parameterName, found := range parameterFound {
			if !found {
				missingParameters = append(missingParameters, parameterName)
			}
		}
		assert.Equal(t, 0, len(missingParameters), "Threshold parameters defined for Profile '%s' but missing for Profile '%s' in Traffic Monitor Config: %s", profileName, profileName, strings.Join(missingParameters, ", "))
	}
}
