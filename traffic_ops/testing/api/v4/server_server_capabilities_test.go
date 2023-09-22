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
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	totest "github.com/apache/trafficcontrol/v8/lib/go-tc/totestv4"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestServerServerCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments, ServerCapabilities, ServerServerCapabilities, DeliveryServicesRequiredCapabilities}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.ServerServerCapability]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerServerCapabilitiesSort()),
				},
				"OK when VALID SERVERID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"serverId": {strconv.Itoa(totest.GetServerID(t, TOSession, "dtrc-edge-01")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerServerCapabilitiesFields(map[string]interface{}{"ServerID": totest.GetServerID(t, TOSession, "dtrc-edge-01")()})),
				},
				"OK when VALID SERVERHOSTNAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"serverHostName": {"atlanta-edge-16"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerServerCapabilitiesFields(map[string]interface{}{"Server": "atlanta-edge-16"})),
				},
				"OK when VALID SERVERCAPABILITY parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"serverCapability": {"asdf"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServerServerCapabilitiesFields(map[string]interface{}{"ServerCapability": "asdf"})),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"serverId"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateServerServerCapabilitiesPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"serverId"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateServerServerCapabilitiesPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"serverId"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateServerServerCapabilitiesPagination("page")),
				},
				"BAD REQUEST when INVALID LIMIT parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"-2"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID OFFSET parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "offset": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PAGE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "page": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"POST": {
				"BAD REQUEST when ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody: tc.ServerServerCapability{
						ServerID:         util.Ptr(totest.GetServerID(t, TOSession, "dtrc-mid-01")()),
						ServerCapability: util.Ptr("disk"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING SERVER ID": {
					ClientSession: TOSession,
					RequestBody: tc.ServerServerCapability{
						Server: util.Ptr("disk"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING SERVER CAPABILITY": {
					ClientSession: TOSession,
					RequestBody: tc.ServerServerCapability{
						ServerID: util.Ptr(totest.GetServerID(t, TOSession, "dtrc-mid-01")()),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"NOT FOUND when SERVER CAPABILITY DOESNT EXIST": {
					ClientSession: TOSession,
					RequestBody: tc.ServerServerCapability{
						ServerID:         util.Ptr(totest.GetServerID(t, TOSession, "dtrc-mid-01")()),
						ServerCapability: util.Ptr("bogus"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"NOT FOUND when SERVER DOESNT EXIST": {
					ClientSession: TOSession,
					RequestBody: tc.ServerServerCapability{
						ServerID:         util.Ptr(99999999),
						ServerCapability: util.Ptr("bogus"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when SERVER TYPE NOT EDGE or MID": {
					ClientSession: TOSession,
					RequestBody: tc.ServerServerCapability{
						ServerID:         util.Ptr(totest.GetServerID(t, TOSession, "trafficvault")()),
						ServerCapability: util.Ptr("bogus"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"DELETE": {
				"OK when NOT the LAST SERVER of CACHE GROUP of TOPOLOGY DS which has REQUIRED CAPABILITIES": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"serverId": {strconv.Itoa(totest.GetServerID(t, TOSession, "dtrc-edge-01")())}, "serverCapability": {"ram"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when LAST SERVER of CACHE GROUP of TOPOLOGY DS which has REQUIRED CAPABILITIES": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"serverId": {strconv.Itoa(totest.GetServerID(t, TOSession, "edge-in-cdn1-only")())}, "serverCapability": {"ram"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SERVER ASSIGNED TO DS with REQUIRED CAPABILITIES": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"serverId": {strconv.Itoa(totest.GetServerID(t, TOSession, "atlanta-org-2")())}, "serverCapability": {"bar"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"NOT FOUND when SERVER SERVER CAPABILITY DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"serverId": {strconv.Itoa(totest.GetServerID(t, TOSession, "atlanta-org-1")())}, "serverCapability": {"doesntexist"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when MISSING SERVER CAPABILITY": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"serverId": {strconv.Itoa(totest.GetServerID(t, TOSession, "atlanta-org-1")())}, "serverCapability": {""}}},
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
							resp, reqInf, err := testCase.ClientSession.GetServerServerCapabilities(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateServerServerCapability(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							var serverId int
							var serverCapability string
							var err error
							if val, ok := testCase.RequestOpts.QueryParameters["serverId"]; ok {
								serverId, err = strconv.Atoi(val[0])
								assert.RequireNoError(t, err, "Expected no error when converting string to int: %v", err)
							} else {
								t.Fatalf("Query Parameter: \"serverId\" is required for DELETE method tests.")
							}
							if val, ok := testCase.RequestOpts.QueryParameters["serverCapability"]; ok {
								serverCapability = val[0]
							} else {
								t.Fatalf("Query Parameter: \"serverCapability\" is required for DELETE method tests.")
							}
							alerts, reqInf, err := testCase.ClientSession.DeleteServerServerCapability(serverId, serverCapability, testCase.RequestOpts)
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

func validateServerServerCapabilitiesFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Server Server Capabilities response to not be nil.")
		serverServerCapabilityResponse := resp.([]tc.ServerServerCapability)
		for field, expected := range expectedResp {
			for _, serverServerCapability := range serverServerCapabilityResponse {
				switch field {
				case "Server":
					assert.RequireNotNil(t, serverServerCapability.Server, "Expected Server to not be nil.")
					assert.Equal(t, expected, *serverServerCapability.Server, "Expected Server to be %v, but got %s", expected, *serverServerCapability.Server)
				case "ServerCapability":
					assert.RequireNotNil(t, serverServerCapability.ServerCapability, "Expected Server Capability to not be nil.")
					assert.Equal(t, expected, *serverServerCapability.ServerCapability, "Expected ServerCapability to be %v, but got %s", expected, *serverServerCapability.ServerCapability)
				case "ServerID":
					assert.RequireNotNil(t, serverServerCapability.ServerID, "Expected Server ID to not be nil.")
					assert.Equal(t, expected, *serverServerCapability.ServerID, "Expected ServerID to be %v, but got %d", expected, *serverServerCapability.ServerID)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateSSCFieldsOnServerCapabilityUpdate(name string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("serverCapability", name)
		ssc, _, err := TOSession.GetServerServerCapabilities(opts)
		assert.RequireNoError(t, err, "Error getting Server Server Capabilities: %v - alerts: %+v", err, ssc.Alerts)
		assert.RequireEqual(t, 1, len(ssc.Response), "Expected one Server Server Capability returned Got: %d", len(ssc.Response))
		validateServerServerCapabilitiesFields(expectedResp)(t, toclientlib.ReqInf{}, ssc.Response, tc.Alerts{}, nil)
	}
}

func validateServerServerCapabilitiesSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Server Server Capabilities response to not be nil.")
		var serverNames []string
		serverServerCapabilityResponse := resp.([]tc.ServerServerCapability)
		for _, serverServerCapability := range serverServerCapabilityResponse {
			assert.RequireNotNil(t, serverServerCapability.Server, "Expected Server to not be nil.")
			serverNames = append(serverNames, *serverServerCapability.Server)
		}
		assert.Equal(t, true, sort.StringsAreSorted(serverNames), "List is not sorted by server names: %v", serverNames)
	}
}

func validateServerServerCapabilitiesPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.ServerServerCapability)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "serverId")
		respBase, _, err := TOSession.GetServerServerCapabilities(opts)
		assert.RequireNoError(t, err, "Cannot get Server Server Capabilities: %v - alerts: %+v", err, respBase.Alerts)

		ssc := respBase.Response
		assert.RequireGreaterOrEqual(t, len(ssc), 3, "Need at least 3 Server Server Capabilities in Traffic Ops to test pagination support, found: %d", len(ssc))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, ssc[:1], paginationResp, "expected GET Server Server Capabilities with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, ssc[1:2], paginationResp, "expected GET Server Server Capabilities with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, ssc[1:2], paginationResp, "expected GET Server Server Capabilities with limit = 1, page = 2 to return second result")
		}
	}
}
