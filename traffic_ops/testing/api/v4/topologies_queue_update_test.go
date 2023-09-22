package v4

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
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

func TestTopologiesQueueUpdate(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices}, func() {

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.TopologiesQueueUpdateRequest]{
			"POST": {
				"OK when VALID REQUEST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"mso-topology"}}},
					RequestBody: tc.TopologiesQueueUpdateRequest{
						Action: "queue",
						CDNID:  int64(totest.GetCDNID(t, TOSession, "cdn1")()),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateTopologiesQueueUpdateFields(map[string]interface{}{"Action": "queue", "CDNID": int64(totest.GetCDNID(t, TOSession, "cdn1")()), "Topology": tc.TopologyName("mso-topology")}),
						validateServerUpdatesAreQueued("ds-top")),
				},
				"BAD REQUEST when INVALID CDNID": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"mso-topology"}}},
					RequestBody: tc.TopologiesQueueUpdateRequest{
						Action: "queue",
						CDNID:  -1,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID ACTION": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"mso-topology"}}},
					RequestBody: tc.TopologiesQueueUpdateRequest{
						Action: "requeue",
						CDNID:  int64(totest.GetCDNID(t, TOSession, "cdn1")()),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when TOPOLOGY DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"nonexistent"}}},
					RequestBody: tc.TopologiesQueueUpdateRequest{
						Action: "queue",
						CDNID:  int64(totest.GetCDNID(t, TOSession, "cdn1")()),
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
							topologyQueueUpdate := testCase.RequestBody
							if _, ok := testCase.RequestOpts.QueryParameters["name"]; !ok {
								t.Fatalf("Query Parameter: \"name\" is required for POST method tests.")
							}
							resp, reqInf, err := testCase.ClientSession.TopologiesQueueUpdate(testCase.RequestOpts.QueryParameters["name"][0], topologyQueueUpdate, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.TopologiesQueueUpdate, resp.Alerts, err)
							}
						})
					}
				}
			})
		}

	})
}

func validateTopologiesQueueUpdateFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Topologies Queue Update response to not be nil.")
		topQueueUpdateResp := resp.(tc.TopologiesQueueUpdate)
		for field, expected := range expectedResp {
			switch field {
			case "Action":
				assert.Equal(t, expected, topQueueUpdateResp.Action, "Expected Action to be %v, but got %s", expected, topQueueUpdateResp.Action)
			case "CDNID":
				assert.Equal(t, expected, topQueueUpdateResp.CDNID, "Expected CDNID to be %v, but got %d", expected, topQueueUpdateResp.CDNID)
			case "Topology":
				assert.Equal(t, expected, topQueueUpdateResp.Topology, "Expected Topology to be %v, but got %s", expected, topQueueUpdateResp.Topology)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}

func validateServerUpdatesAreQueued(topologyDS string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Topologies Queue Update response to not be nil.")
		topQueueUpdateResp := resp.(tc.TopologiesQueueUpdate)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("dsId", strconv.Itoa(totest.GetDeliveryServiceId(t, TOSession, topologyDS)()))
		serversResponse, _, err := TOSession.GetServers(opts)
		assert.RequireNoError(t, err, "Expected no error when getting servers: %v", err)

		for _, server := range serversResponse.Response {
			assert.RequireNotNil(t, server.CDNID, "Expected Server CDNID to not be nil.")
			assert.RequireNotNil(t, server.HostName, "Expected Server HostName to not be nil.")
			assert.RequireNotNil(t, server.UpdPending, "Expected Server UpdPending to not be nil.")
			if *server.CDNID != int(topQueueUpdateResp.CDNID) {
				continue
			}
			assert.Equal(t, true, *server.UpdPending, "Expected Server %s Update Pending flag to be set to true.", *server.HostName)
		}
	}
}
