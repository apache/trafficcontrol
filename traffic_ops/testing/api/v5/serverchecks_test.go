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
	"net/url"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestServerChecks(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCheckExtensions, ServerChecks}, func() {

		extensionSession := utils.CreateV5Session(t, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.Default.Session.TimeoutInSecs)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.ServercheckRequestNullable]{
			"GET": {
				"OK when VALID request": {
					ClientSession: extensionSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID HOSTNAME parameter": {
					ClientSession: extensionSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"hostName": {"atlanta-edge-01"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateServerCheckFields(map[string]interface{}{"HostName": "atlanta-edge-01", "Checks": map[string]int{"ORT": 13}})),
				},
				"OK when VALID ID parameter": {
					ClientSession: extensionSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"id": {strconv.Itoa(GetServerID(t, "atlanta-edge-01")())}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1)),
				},
			},
			"POST": {
				"OK when UPDATING EXISTING SERVER CHECK": {
					ClientSession: extensionSession,
					RequestBody: tc.ServercheckRequestNullable{
						Name:     util.Ptr("ILO"),
						HostName: util.Ptr("atlanta-edge-01"),
						Value:    util.Ptr(0),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServerCheckCreateFields("atlanta-edge-01", map[string]interface{}{"HostName": "atlanta-edge-01", "Checks": map[string]int{"ORT": 13, "ILO": 0}})),
				},
				"BAD REQUEST when NO SERVER ID": {
					ClientSession: extensionSession,
					RequestBody: tc.ServercheckRequestNullable{
						ID:    nil,
						Name:  util.Ptr("ILO"),
						Value: util.Ptr(1),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID SERVER ID": {
					ClientSession: extensionSession,
					RequestBody: tc.ServercheckRequestNullable{
						HostName: util.Ptr("atlanta-edge-01"),
						ID:       util.Ptr(-1),
						Name:     util.Ptr("ILO"),
						Value:    util.Ptr(1),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID SERVERCHECK SHORT NAME": {
					ClientSession: extensionSession,
					RequestBody: tc.ServercheckRequestNullable{
						HostName: util.Ptr("atlanta-edge-01"),
						Name:     util.Ptr("BOGUS"),
						Value:    util.Ptr(1),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"FORBIDDEN when NON EXTENSION USER": {
					ClientSession: TOSession,
					RequestBody: tc.ServercheckRequestNullable{
						HostName: util.Ptr("atlanta-edge-01"),
						ID:       util.Ptr(GetServerID(t, "atlanta-edge-01")()),
						Name:     util.Ptr("TEST"),
						Value:    util.Ptr(1),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetServersChecks(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.InsertServerCheckStatus(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, tc.Alerts{Alerts: resp.Alerts}, err)
							}
						})
					}
				}
			})
		}
	})
}

func validateServerCheckFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Server Check response to not be nil.")
		serverCheckResp := resp.([]tc.GenericServerCheck)
		for field, expected := range expectedResp {
			for _, serverCheck := range serverCheckResp {
				switch field {
				case "HostName":
					assert.Equal(t, expected, serverCheck.HostName, "Expected HostName to be %v, but got %s", expected, serverCheck.HostName)
				case "Checks":
					for name, value := range expected.(map[string]int) {
						assert.RequireNotNil(t, serverCheck.Checks[name], "Expected Checks[%s] value to not be nil.", name)
						assert.Equal(t, value, *serverCheck.Checks[name], "Expected Checks ILO Value to be %d, but got %s", value, *serverCheck.Checks[name])
					}
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateServerCheckCreateFields(hostName string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("hostName", hostName)
		serverChecks, _, err := TOSession.GetServersChecks(opts)
		assert.RequireNoError(t, err, "Error getting Server Checks: %v - alerts: %+v", err, serverChecks.Alerts)
		assert.RequireEqual(t, 1, len(serverChecks.Response), "Expected one Server Check returned Got: %d", len(serverChecks.Response))
		validateServerCheckFields(expectedResp)(t, toclientlib.ReqInf{}, serverChecks.Response, tc.Alerts{}, nil)
	}
}

func CreateTestServerChecks(t *testing.T) {
	extensionSession := utils.CreateV5Session(t, Config.TrafficOps.URL, Config.TrafficOps.Users.Extension, Config.TrafficOps.UserPassword, Config.Default.Session.TimeoutInSecs)

	for _, servercheck := range testData.Serverchecks {
		resp, _, err := extensionSession.InsertServerCheckStatus(servercheck, client.RequestOptions{})
		assert.RequireNoError(t, err, "Could not insert Servercheck: %v - alerts: %+v", err, resp.Alerts)
	}
}

// Need to define no-op function as TCObj interface expects a delete function
// There is no delete path for serverchecks
func DeleteTestServerChecks(*testing.T) {
	return
}
