package v5

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
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestSteering(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, ServerCapabilities, ServerServerCapabilities, DeliveryServices, Users, SteeringTargets}, func() {

		methodTests := utils.TestCase[client.Session, client.RequestOptions, struct{}]{
			"GET": {
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(2),
						validateSteeringFields(map[string]interface{}{"TargetsLength": 1, "TargetsOrder": int32(0),
							"TargetsGeoOrderPtr": (*int)(nil), "TargetsLongitudePtr": (*float64)(nil), "TargetsLatitudePtr": (*float64)(nil), "TargetsWeight": int32(42)})),
				},
			},
		}
		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.Steering(testCase.RequestOpts)
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

func validateSteeringFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Steering response to not be nil.")
		steeringResp := resp.([]tc.Steering)
		for field, expected := range expectedResp {
			for _, steering := range steeringResp {
				switch field {
				case "TargetsLength":
					assert.Equal(t, expected, len(steering.Targets), "Expected Targets Length to be %v, but got %d", expected, len(steering.Targets))
				case "TargetsOrder":
					assert.RequireEqual(t, 1, len(steering.Targets), "Expected Targets Length to be %d, but got %d", 1, len(steering.Targets))
					assert.Equal(t, expected, steering.Targets[0].Order, "Expected Targets Order to be %v, but got %d", expected, steering.Targets[0].Order)
				case "TargetsGeoOrderPtr":
					assert.RequireEqual(t, 1, len(steering.Targets), "Expected Targets Length to be %d, but got %d", 1, len(steering.Targets))
					assert.Equal(t, expected, steering.Targets[0].GeoOrder, "Expected Targets GeoOrder to be %v, but got %v", nil, steering.Targets[0].GeoOrder)
				case "TargetsLongitudePtr":
					assert.RequireEqual(t, 1, len(steering.Targets), "Expected Targets Length to be %d, but got %d", 1, len(steering.Targets))
					assert.Equal(t, expected, steering.Targets[0].Longitude, "Expected Targets Longitude to be %v, but got %v", nil, steering.Targets[0].Longitude)
				case "TargetsLatitudePtr":
					assert.RequireEqual(t, 1, len(steering.Targets), "Expected Targets Length to be %d, but got %d", 1, len(steering.Targets))
					assert.Equal(t, expected, steering.Targets[0].Latitude, "Expected Targets Latitude to be %v, but got %v", nil, steering.Targets[0].Latitude)
				case "TargetsWeight":
					assert.RequireEqual(t, 1, len(steering.Targets), "Expected Targets Length to be %d, but got %d", 1, len(steering.Targets))
					assert.Equal(t, expected, steering.Targets[0].Weight, "Expected Targets Weight to be %v, but got %v", expected, steering.Targets[0].Weight)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}
