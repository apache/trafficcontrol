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
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestServerCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, ServerServerCapabilities}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.ServerCapabilityV5]{
			"GET": {
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerCapabilitiesSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"ram"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"EMPTY RESPONSE when INVALID NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"abcd"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
			},
			"POST": {
				"BAD REQUEST when ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody: tc.ServerCapabilityV5{
						Name:        "foo",
						Description: "foo servers",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID NAME": {
					ClientSession: TOSession,
					RequestBody: tc.ServerCapabilityV5{
						Name:        "b@dname",
						Description: "Server Capability",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"blah"}}},
					RequestBody: tc.ServerCapabilityV5{
						Name:        "newname",
						Description: "Server Capability for new name",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServerCapabilitiesUpdateFields(map[string]interface{}{"Name": "newname"}),
						validateSSCFieldsOnServerCapabilityUpdate("newname", map[string]interface{}{"ServerCapability": "newname"})),
				},
				"BAD REQUEST when NAME DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"invalid"}}},
					RequestBody: tc.ServerCapabilityV5{
						Name:        "newname",
						Description: "Server Capability for new name",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					ClientSession: TOSession,
					RequestOpts: client.RequestOptions{
						QueryParameters: url.Values{"name": {"disk"}},
						Header:          http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					},
					RequestBody: tc.ServerCapabilityV5{
						Name:        "newname",
						Description: "Server Capability for new name",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					ClientSession: TOSession,
					RequestOpts: client.RequestOptions{
						QueryParameters: url.Values{"name": {"disk"}},
						Header:          http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					},
					RequestBody: tc.ServerCapabilityV5{
						Name:        "newname",
						Description: "Server Capability for new name",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"NOT FOUND when NAME DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"invalid"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when EMPTY NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {""}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetServerCapabilities(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateServerCapability(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateServerCapability(testCase.RequestOpts.QueryParameters["name"][0], testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteServerCapability(testCase.RequestOpts.QueryParameters["name"][0], testCase.RequestOpts)
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

func validateServerCapabilitiesUpdateFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Server Capabilities response to not be nil.")
		serverCapabilitiesResp := resp.(tc.ServerCapabilityV5)
		for field, expected := range expectedResp {
			switch field {
			case "Name":
				assert.Equal(t, expected, serverCapabilitiesResp.Name, "Expected Name to be %v, but got %s", expected, serverCapabilitiesResp.Name)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}

func validateServerCapabilitiesSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Server Capabilities response to not be nil.")
		var serverCapabilityNames []string
		serverCapabilitiesResp := resp.([]tc.ServerCapabilityV5)
		for _, serverCapability := range serverCapabilitiesResp {
			serverCapabilityNames = append(serverCapabilityNames, serverCapability.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(serverCapabilityNames), "List is not sorted by their names: %v", serverCapabilityNames)
	}
}

func CreateTestServerCapabilities(t *testing.T) {
	for _, sc := range testData.ServerCapabilities {
		resp, _, err := TOSession.CreateServerCapability(sc, client.RequestOptions{})
		assert.RequireNoError(t, err, "Unexpected error creating Server Capability '%s': %v - alerts: %+v", sc.Name, err, resp.Alerts)
	}
}

func DeleteTestServerCapabilities(t *testing.T) {
	serverCapabilities, _, err := TOSession.GetServerCapabilities(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get Server Capabilities: %v - alerts: %+v", err, serverCapabilities.Alerts)

	for _, serverCapability := range serverCapabilities.Response {
		alerts, _, err := TOSession.DeleteServerCapability(serverCapability.Name, client.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Server Capability '%s': %v - alerts: %+v", serverCapability.Name, err, alerts.Alerts)
		// Retrieve the Server Capability to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", serverCapability.Name)
		getServerCapability, _, err := TOSession.GetServerCapabilities(opts)
		assert.NoError(t, err, "Error getting Server Capability '%s' after deletion: %v - alerts: %+v", serverCapability.Name, err, getServerCapability.Alerts)
		assert.Equal(t, 0, len(getServerCapability.Response), "Expected Server Capability '%s' to be deleted, but it was found in Traffic Ops", serverCapability.Name)
	}
}
