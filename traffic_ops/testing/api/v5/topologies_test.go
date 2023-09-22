package v5

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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestTopologies(t *testing.T) {
	WithObjs(t, []TCObj{Tenants, Users, Types, CacheGroups, CDNs, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, Servers, ServerCapabilities, ServerServerCapabilities, Topologies, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments}, func() {

		readOnlyUserSession := utils.CreateV5Session(t, Config.TrafficOps.URL, "readonlyuser", "pa$$word", Config.Default.Session.TimeoutInSecs)

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.TopologyV5]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when READ ONLY USER": {
					ClientSession: readOnlyUserSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"EMPTY RESPONSE when TOPOLOGY DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"non-existent-topology"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
			},
			"POST": {
				"OK when MISSING DESCRIPTION": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name:  "topology-missing-description",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "cachegroup1", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when EMPTY NODES": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name:  "topology-no-nodes",
						Nodes: []tc.TopologyNodeV5{},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when NODE PARENT of ITSELF": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name:  "self-parent",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "cachegroup1", Parents: []int{0}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when TOO MANY PARENTS": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name: "too-many-parents",
						Nodes: []tc.TopologyNodeV5{
							{
								Cachegroup: "parentCachegroup",
								Parents:    []int{},
							},
							{
								Cachegroup: "secondaryCachegroup",
								Parents:    []int{},
							},
							{
								Cachegroup: "parentCachegroup2",
								Parents:    []int{},
							},
							{
								Cachegroup: "cachegroup1",
								Parents:    []int{0, 1, 2},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when EDGE_LOC PARENTS MID_LOC": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name: "edge-parents-mid",
						Nodes: []tc.TopologyNodeV5{
							{
								Cachegroup: "parentCachegroup",
								Parents:    []int{1},
							},
							{
								Cachegroup: "cachegroup2",
								Parents:    []int{},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CYCLICAL NODES": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name: "cyclical-nodes",
						Nodes: []tc.TopologyNodeV5{
							{
								Cachegroup: "cachegroup1",
								Parents:    []int{1, 2},
							},
							{
								Cachegroup: "parentCachegroup",
								Parents:    []int{2},
							},
							{
								Cachegroup: "secondaryCachegroup",
								Parents:    []int{1},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CYCLES ACROSS TOPOLOGIES": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name: "cyclical-nodes-tiered",
						Nodes: []tc.TopologyNodeV5{
							{
								Cachegroup: "parentCachegroup",
								Parents:    []int{1},
							},
							{
								Cachegroup: "parentCachegroup2",
								Parents:    []int{},
							},
							{
								Cachegroup: "cachegroup1",
								Parents:    []int{0},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CYCLICAL NODES BUT EMPTY CACHE GROUPS": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name: "cyclical-nodes-nontopology",
						Nodes: []tc.TopologyNodeV5{
							{
								Cachegroup: "edge-parent1",
								Parents:    []int{1},
							},
							{
								Cachegroup: "has-edge-parent1",
								Parents:    []int{},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when OUT-OF-BOUNDS PARENT INDEX": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name:  "outofbounds",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "cachegroup1", Parents: []int{7}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CACHEGROUP DOESNT EXIST": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name:  "topology-nonexistent-cachegroup",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "doesntexist", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING NAME": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Description: "missing name",
						Nodes:       []tc.TopologyNodeV5{{Cachegroup: "cachegroup1", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name:  "mso-topology",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "cachegroup1", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CACHEGROUP has NO SERVERS": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name:  "topology-empty-cg",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "noServers", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when DUPLICATE PARENTS": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name: "topology-duplicate-parents",
						Nodes: []tc.TopologyNodeV5{
							{
								Cachegroup: "parentCachegroup",
								Parents:    []int{},
							},
							{
								Cachegroup: "cachegroup1",
								Parents:    []int{0, 0},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ORG_LOC is CHILD NODE": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name: "topology-orgloc-child",
						Nodes: []tc.TopologyNodeV5{
							{
								Cachegroup: "cachegroup1",
								Parents:    []int{},
							},
							{
								Cachegroup: "multiOriginCachegroup",
								Parents:    []int{0},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when LEAF NODE is a MID_LOC": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name:  "topology-midloc-leaf",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "parentCachegroup", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"WARNING LEVEL ALERT when MID PARENTING EDGE": {
					ClientSession: TOSession,
					RequestBody: tc.TopologyV5{
						Name:        "topology-mid-parent",
						Description: "mid parent to edge",
						Nodes: []tc.TopologyNodeV5{
							{
								Cachegroup: "cachegroup1",
								Parents:    []int{1},
							},
							{
								Cachegroup: "cachegroup2",
								Parents:    []int{},
							},
						},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.HasAlertLevel(tc.WarnLevel.String())),
				},
				"FORBIDDEN when READ-ONLY USER": {
					ClientSession: readOnlyUserSession,
					RequestBody: tc.TopologyV5{
						Name:  "topology-ro",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "cachegroup1", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"top-with-no-mids"}}},
					RequestBody: tc.TopologyV5{
						Name:        "top-with-no-mids-updated",
						Description: "Updating fields",
						Nodes:       []tc.TopologyNodeV5{{Cachegroup: "cachegroup2", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateTopologiesUpdateCreateFields(map[string]interface{}{"Name": "top-with-no-mids-updated", "Description": "Updating fields"})),
				},
				"BAD REQUEST when OUT-OF-BOUNDS PARENT INDEX": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"another-topology"}}},
					RequestBody: tc.TopologyV5{
						Name:  "topology-invalid-parent",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "cachegroup1", Parents: []int{100}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CACHEGROUP has NO SERVERS": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"another-topology"}}},
					RequestBody: tc.TopologyV5{
						Name:  "topology-empty-cg",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "noServers", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CACHEGROUP SERVERS DO NOT HAVE REQUIRED CAPABILITIES": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"top-for-ds-req"}}},
					RequestBody: tc.TopologyV5{
						Name:  "top-for-ds-req",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "cachegroup1", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when NODE PARENT of ITSELF": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"another-topology"}}},
					RequestBody: tc.TopologyV5{
						Name:  "another-topology",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "cachegroup1", Parents: []int{0}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CACHEGROUP HAS NO SERVERS IN TOPOLOGY CDN": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"top-used-by-cdn1-and-cdn2"}}},
					RequestBody: tc.TopologyV5{
						Name:  "top-used-by-cdn1-and-cdn2",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "cdn1-only", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ORG_LOC is CHILD NODE": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"another-topology"}}},
					RequestBody: tc.TopologyV5{
						Name: "topology-orgloc-child",
						Nodes: []tc.TopologyNodeV5{
							{
								Cachegroup: "cachegroup1",
								Parents:    []int{},
							},
							{
								Cachegroup: "multiOriginCachegroup",
								Parents:    []int{0},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when LEAF NODE is a MID_LOC": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"another-topology"}}},
					RequestBody: tc.TopologyV5{
						Name:        "topology-child-midloc",
						Description: "child mid_loc",
						Nodes:       []tc.TopologyNodeV5{{Cachegroup: "parentCachegroup", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when DUPLICATE PARENTS": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"another-topology"}}},
					RequestBody: tc.TopologyV5{
						Name: "topology-same-parents",
						Nodes: []tc.TopologyNodeV5{
							{
								Cachegroup: "cachegroup1",
								Parents:    []int{},
							},
							{
								Cachegroup: "cachegroup2",
								Parents:    []int{0, 0},
							},
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when REMOVING ORG ASSIGNED DS": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"mso-topology"}}},
					RequestBody: tc.TopologyV5{
						Name:  "mso-topology",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "topology-edge-cg-01", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"FORBIDDEN when READ-ONLY USER": {
					ClientSession: readOnlyUserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"another-topology"}}},
					RequestBody: tc.TopologyV5{
						Name:  "topology-ro",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "cachegroup1", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					ClientSession: TOSession,
					RequestOpts: client.RequestOptions{
						QueryParameters: url.Values{"name": {"another-topology"}},
						Header:          http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					},
					RequestBody: tc.TopologyV5{
						Name:  "another-topology",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "cachegroup1", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					ClientSession: TOSession,
					RequestOpts: client.RequestOptions{
						QueryParameters: url.Values{"name": {"another-topology"}},
						Header:          http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					},
					RequestBody: tc.TopologyV5{
						Name:  "another-topology",
						Nodes: []tc.TopologyNodeV5{{Cachegroup: "cachegroup1", Parents: []int{}}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"BAD REQUEST when TOPOLOGY DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"non-existent-topology"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when TOPOLOGY IN USE by DELIVERY SERVICE": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"mso-topology"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"FORBIDDEN when READ-ONLY USER": {
					ClientSession: readOnlyUserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"mso-topology"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					topology := testCase.RequestBody

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetTopologies(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateTopology(topology, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							if _, ok := testCase.RequestOpts.QueryParameters["name"]; !ok {
								t.Fatalf("Query Parameter: \"name\" is required for PUT method tests.")
							}
							resp, reqInf, err := testCase.ClientSession.UpdateTopology(testCase.RequestOpts.QueryParameters["name"][0], topology, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							if _, ok := testCase.RequestOpts.QueryParameters["name"]; !ok {
								t.Fatalf("Query Parameter: \"name\" is required for DELETE method tests.")
							}
							alerts, reqInf, err := testCase.ClientSession.DeleteTopology(testCase.RequestOpts.QueryParameters["name"][0], testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					default:
						t.Errorf("Method: %s, is not a valid test method.", method)
					}
				}
			})
		}
	})
}

func validateTopologiesFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Topology response to not be nil.")
		topologyResp := resp.([]tc.TopologyV5)
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
		topology := resp.(tc.TopologyV5)
		topologyResp := []tc.TopologyV5{topology}
		validateTopologiesFields(expectedResp)(t, toclientlib.ReqInf{}, topologyResp, tc.Alerts{}, nil)
	}
}

func CreateTestTopologies(t *testing.T) {
	for _, topology := range testData.Topologies {
		resp, _, err := TOSession.CreateTopology(topology, client.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Topology: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestTopologies(t *testing.T) {
	topologies, _, err := TOSession.GetTopologies(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get Topologies: %v - alerts: %+v", err, topologies.Alerts)

	for _, topology := range topologies.Response {
		alerts, _, err := TOSession.DeleteTopology(topology.Name, client.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Topology: %v - alerts: %+v", err, alerts.Alerts)
		// Retrieve the Topology to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", topology.Name)
		resp, _, err := TOSession.GetTopologies(opts)
		assert.NoError(t, err, "Unexpected error trying to fetch Topologies after deletion: %v - alerts: %+v", err, resp.Alerts)
		assert.Equal(t, 0, len(resp.Response), "Expected Topology '%s' to be deleted, but it was found in Traffic Ops", topology.Name)
	}
}
