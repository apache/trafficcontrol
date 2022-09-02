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
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

func TestTopologies(t *testing.T) {
	WithObjs(t, []TCObj{Tenants, Users, Types, CacheGroups, CDNs, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, Servers, ServerCapabilities, ServerServerCapabilities, Topologies, ServiceCategories, DeliveryServices, DeliveryServicesRequiredCapabilities, DeliveryServiceServerAssignments}, func() {

		methodTests := utils.V3TestCase{
			"GET": {
				"OK when VALID REQUEST": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
			},
			"POST": {
				"OK when MISSING DESCRIPTION": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":  "topology-missing-description",
						"nodes": []map[string]interface{}{{"cachegroup": "cachegroup1", "parents": []int{}}},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when EMPTY NODES": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":  "topology-no-nodes",
						"nodes": []map[string]interface{}{},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when NODE PARENT of ITSELF": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":  "self-parent",
						"nodes": []map[string]interface{}{{"cachegroup": "cachegroup1", "parents": []int{0}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when TOO MANY PARENTS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name": "too-many-parents",
						"nodes": []map[string]interface{}{
							{
								"cachegroup": "parentCachegroup",
								"parents":    []int{},
							},
							{
								"cachegroup": "secondaryCachegroup",
								"parents":    []int{},
							},
							{
								"cachegroup": "parentCachegroup2",
								"parents":    []int{},
							},
							{
								"cachegroup": "cachegroup1",
								"parents":    []int{0, 1, 2},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when EDGE_LOC PARENTS MID_LOC": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name": "edge-parents-mid",
						"nodes": []map[string]interface{}{
							{
								"cachegroup": "parentCachegroup",
								"parents":    []int{1},
							},
							{
								"cachegroup": "cachegroup2",
								"parents":    []int{},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CYCLICAL NODES": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name": "cyclical-nodes",
						"nodes": []map[string]interface{}{
							{
								"cachegroup": "cachegroup1",
								"parents":    []int{1, 2},
							},
							{
								"cachegroup": "parentCachegroup",
								"parents":    []int{2},
							},
							{
								"cachegroup": "secondaryCachegroup",
								"parents":    []int{1},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CYCLES ACROSS TOPOLOGIES": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name": "cyclical-nodes-tiered",
						"nodes": []map[string]interface{}{
							{
								"cachegroup": "parentCachegroup",
								"parents":    []int{1},
							},
							{
								"cachegroup": "parentCachegroup2",
								"parents":    []int{},
							},
							{
								"cachegroup": "cachegroup1",
								"parents":    []int{0},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CYCLICAL NODES BUT EMPTY CACHE GROUPS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name": "cyclical-nodes-nontopology",
						"nodes": []map[string]interface{}{
							{
								"cachegroup": "edge-parent1",
								"parents":    []int{1},
							},
							{
								"cachegroup": "has-edge-parent1",
								"parents":    []int{},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when OUT-OF-BOUNDS PARENT INDEX": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":  "outofbounds",
						"nodes": []map[string]interface{}{{"cachegroup": "cachegroup1", "parents": []int{7}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CACHEGROUP DOESNT EXIST": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":  "topology-nonexistent-cachegroup",
						"nodes": []map[string]interface{}{{"cachegroup": "doesntexist", "parents": []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING NAME": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"description": "missing name",
						"nodes":       []map[string]interface{}{{"cachegroup": "cachegroup1", "parents": []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":  "mso-topology",
						"nodes": []map[string]interface{}{{"cachegroup": "cachegroup1", "parents": []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CACHEGROUP has NO SERVERS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":  "topology-empty-cg",
						"nodes": []map[string]interface{}{{"cachegroup": "noServers", "parents": []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when DUPLICATE PARENTS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name": "topology-duplicate-parents",
						"nodes": []map[string]interface{}{
							{
								"cachegroup": "parentCachegroup",
								"parents":    []int{},
							},
							{
								"cachegroup": "cachegroup1",
								"parents":    []int{0, 0},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ORG_LOC is CHILD NODE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name": "topology-orgloc-child",
						"nodes": []map[string]interface{}{
							{
								"cachegroup": "cachegroup1",
								"parents":    []int{},
							},
							{
								"cachegroup": "multiOriginCachegroup",
								"parents":    []int{0},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when LEAF NODE is a MID_LOC": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":  "topology-midloc-leaf",
						"nodes": []map[string]interface{}{{"cachegroup": "parentCachegroup", "parents": []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"WARNING LEVEL ALERT when MID PARENTING EDGE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"name":        "topology-mid-parent",
						"description": "mid parent to edge",
						"nodes": []map[string]interface{}{
							{
								"cachegroup": "cachegroup1",
								"parents":    []int{1},
							},
							{
								"cachegroup": "cachegroup2",
								"parents":    []int{},
							},
						},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.HasAlertLevel(tc.WarnLevel.String())),
				},
			},
			"PUT": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"top-with-no-mids"}},
					RequestBody: map[string]interface{}{
						"name":        "top-with-no-mids-updated",
						"description": "Updating fields",
						"nodes":       []map[string]interface{}{{"cachegroup": "cachegroup2", "parents": []int{}}},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateTopologiesUpdateCreateFields(map[string]interface{}{"Name": "top-with-no-mids-updated", "Description": "Updating fields"})),
				},
				"BAD REQUEST when OUT-OF-BOUNDS PARENT INDEX": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"another-topology"}},
					RequestBody: map[string]interface{}{
						"name":  "topology-invalid-parent",
						"nodes": []map[string]interface{}{{"cachegroup": "cachegroup1", "parents": []int{100}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CACHEGROUP has NO SERVERS": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"another-topology"}},
					RequestBody: map[string]interface{}{
						"name":  "topology-empty-cg",
						"nodes": []map[string]interface{}{{"cachegroup": "noServers", "parents": []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CACHEGROUP SERVERS DO NOT HAVE REQUIRED CAPABILITIES": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"top-for-ds-req"}},
					RequestBody: map[string]interface{}{
						"name":  "top-for-ds-req",
						"nodes": []map[string]interface{}{{"cachegroup": "cachegroup1", "parents": []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when NODE PARENT of ITSELF": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"another-topology"}},
					RequestBody: map[string]interface{}{
						"name":  "another-topology",
						"nodes": []map[string]interface{}{{"cachegroup": "cachegroup1", "parents": []int{0}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CACHEGROUP HAS NO SERVERS IN TOPOLOGY CDN": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"top-used-by-cdn1-and-cdn2"}},
					RequestBody: map[string]interface{}{
						"name":  "top-used-by-cdn1-and-cdn2",
						"nodes": []map[string]interface{}{{"cachegroup": "cdn1-only", "parents": []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ORG_LOC is CHILD NODE": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"another-topology"}},
					RequestBody: map[string]interface{}{
						"name": "topology-orgloc-child",
						"nodes": []map[string]interface{}{
							{
								"cachegroup": "cachegroup1",
								"parents":    []int{},
							},
							{
								"cachegroup": "multiOriginCachegroup",
								"parents":    []int{0},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when LEAF NODE is a MID_LOC": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"another-topology"}},
					RequestBody: map[string]interface{}{
						"name":        "topology-child-midloc",
						"description": "child mid_loc",
						"nodes":       []map[string]interface{}{{"cachegroup": "parentCachegroup", "parents": []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when DUPLICATE PARENTS": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"another-topology"}},
					RequestBody: map[string]interface{}{
						"name": "topology-same-parents",
						"nodes": []map[string]interface{}{
							{
								"cachegroup": "cachegroup1",
								"parents":    []int{},
							},
							{
								"cachegroup": "cachegroup2",
								"parents":    []int{0, 0},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when REMOVING ORG ASSIGNED DS": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"mso-topology"}},
					RequestBody: map[string]interface{}{
						"name":  "mso-topology",
						"nodes": []map[string]interface{}{{"cachegroup": "topology-edge-cg-01", "parents": []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					topology := tc.Topology{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &topology)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetTopologiesWithHdr(testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateTopology(topology)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							if _, ok := testCase.RequestParams["name"]; !ok {
								t.Fatalf("Query Parameter: \"name\" is required for PUT method tests.")
							}
							resp, reqInf, err := testCase.ClientSession.UpdateTopology(testCase.RequestParams["name"][0], topology)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							if _, ok := testCase.RequestParams["name"]; !ok {
								t.Fatalf("Query Parameter: \"name\" is required for DELETE method tests.")
							}
							alerts, reqInf, err := testCase.ClientSession.DeleteTopology(testCase.RequestParams["name"][0])
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

func validateTopologiesFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Topology response to not be nil.")
		topologyResp := resp.([]tc.Topology)
		for field, expected := range expectedResp {
			for _, topology := range topologyResp {
				switch field {
				case "Description":
					assert.Equal(t, expected, topology.Description, "Expected Description to be %v, but got %s", expected, topology.Description)
				case "Name":
					assert.Equal(t, expected, topology.Name, "Expected Name to be %v, but got %s", expected, topology.Name)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateTopologiesUpdateCreateFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Topology response to not be nil.")
		topology := resp.(tc.Topology)
		topologyResp := []tc.Topology{topology}
		validateTopologiesFields(expectedResp)(t, toclientlib.ReqInf{}, topologyResp, tc.Alerts{}, nil)
	}
}

func CreateTestTopologies(t *testing.T) {
	for _, topology := range testData.Topologies {
		resp, _, err := TOSession.CreateTopology(topology)
		assert.RequireNoError(t, err, "Could not create Topology: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestTopologies(t *testing.T) {
	topologies, _, err := TOSession.GetTopologiesWithHdr(nil)
	assert.NoError(t, err, "Cannot get Topologies: %v", err)

	for _, topology := range topologies {
		alerts, _, err := TOSession.DeleteTopology(topology.Name)
		assert.NoError(t, err, "Cannot delete Topology: %v - alerts: %+v", err, alerts.Alerts)
	}
	// Retrieve the Topologies to see if they were deleted
	resp, _, err := TOSession.GetTopologiesWithHdr(nil)
	assert.NoError(t, err, "Unexpected error trying to fetch Topologies after deletion: %v", err)
	assert.Equal(t, 0, len(resp), "Expected Topologies to be deleted, found: %d", len(resp))
}
