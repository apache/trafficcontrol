package v3

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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

func TestTopologiesQueueUpdate(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices}, func() {

		methodTests := utils.V3TestCaseT[tc.TopologiesQueueUpdateRequest]{
			"POST": {
				"OK when VALID REQUEST": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"mso-topology"}},
					RequestBody: tc.TopologiesQueueUpdateRequest{
						Action: "queue",
						CDNID:  int64(GetCDNID(t, "cdn1")()),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateTopologiesQueueUpdateFields(map[string]interface{}{"Action": "queue", "CDNID": int64(GetCDNID(t, "cdn1")()), "Topology": tc.TopologyName("mso-topology")}),
						validateServerUpdatesAreQueued("ds-top")),
				},
				"BAD REQUEST when INVALID CDNID": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"mso-topology"}},
					RequestBody: tc.TopologiesQueueUpdateRequest{
						Action: "queue",
						CDNID:  -1,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID ACTION": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"mso-topology"}},
					RequestBody: tc.TopologiesQueueUpdateRequest{
						Action: "requeue",
						CDNID:  int64(GetCDNID(t, "cdn1")()),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when TOPOLOGY DOESNT EXIST": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"nonexistent"}},
					RequestBody: tc.TopologiesQueueUpdateRequest{
						Action: "queue",
						CDNID:  int64(GetCDNID(t, "cdn1")()),
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
							resp, reqInf, err := testCase.ClientSession.TopologiesQueueUpdate(tc.TopologyName(testCase.RequestParams["name"][0]), testCase.RequestBody)
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

		params := url.Values{}
		params.Set("dsId", strconv.Itoa(GetDeliveryServiceId(t, topologyDS)()))
		serversResponse, _, err := TOSession.GetServersWithHdr(&params, nil)
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
