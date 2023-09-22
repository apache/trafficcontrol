package v3

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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestServersIDStatus(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServiceCategories, Topologies, DeliveryServices, DeliveryServiceServerAssignments}, func() {

		methodTests := utils.V3TestCaseT[tc.ServerPutStatus]{
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetServerID(t, "atlanta-mid-01"),
					ClientSession: TOSession,
					RequestBody: tc.ServerPutStatus{
						Status:        util.JSONNameOrIDStr{ID: util.Ptr(GetStatusID(t, "OFFLINE")())},
						OfflineReason: util.Ptr("test mid"),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateUpdPending("atlanta-mid-01")),
				},
				"OK when using STATUS ID FIELD": {
					EndpointID:    GetServerID(t, "atlanta-mid-16"),
					ClientSession: TOSession,
					RequestBody: tc.ServerPutStatus{
						Status:        util.JSONNameOrIDStr{ID: util.Ptr(GetStatusID(t, "OFFLINE")())},
						OfflineReason: util.Ptr("test mid"),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateUpdPending("atlanta-mid-16")),
				},
				"VALIDATE TOPOLOGY DESCENDANTS receive STATUS UPDATES": {
					EndpointID:    GetServerID(t, "topology-mid-04"),
					ClientSession: TOSession,
					RequestBody: tc.ServerPutStatus{
						Status:        util.JSONNameOrIDStr{ID: util.Ptr(GetStatusID(t, "OFFLINE")())},
						OfflineReason: util.Ptr("test topology mid cachegroup"),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateUpdPendingSpecificServers(map[string]bool{"topology-mid-04": false, "midInTopologyMidCg04": false,
							"topology-edge-01": true, "edgeInTopologyEdgeCg01": false, "topology-edge-02": false, "edgeInTopologyEdgeCg02": false}),
						validateParentPendingSpecificServers(map[string]bool{"topology-edge-01": true, "edgeInTopologyEdgeCg02": false})),
				},
				"NOT FOUND when SERVER DOESNT EXIST": {
					EndpointID:    func() int { return 11111111 },
					ClientSession: TOSession,
					RequestBody: tc.ServerPutStatus{
						Status:        util.JSONNameOrIDStr{Name: util.Ptr("OFFLINE")},
						OfflineReason: util.Ptr("test last edge"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when STATUS DOESNT EXIST": {
					EndpointID:    GetServerID(t, "atlanta-mid-16"),
					ClientSession: TOSession,
					RequestBody: tc.ServerPutStatus{
						Status:        util.JSONNameOrIDStr{Name: util.Ptr("NOT_A_REAL_STATUS")},
						OfflineReason: util.Ptr("test last edge"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING OFFLINE REASON when OFFLINE STATUS": {
					EndpointID:    GetServerID(t, "atlanta-mid-16"),
					ClientSession: TOSession,
					RequestBody: tc.ServerPutStatus{
						Status: util.JSONNameOrIDStr{Name: util.Ptr("OFFLINE")},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING OFFLINE REASON when ADMIN_DOWN STATUS": {
					EndpointID:    GetServerID(t, "atlanta-mid-16"),
					ClientSession: TOSession,
					RequestBody: tc.ServerPutStatus{
						Status: util.JSONNameOrIDStr{Name: util.Ptr("ADMIN_DOWN")},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"CONFLICT when SERVER STATUS OFFLINE when ONLY EDGE SERVER ASSIGNED": {
					EndpointID:    GetServerID(t, "test-ds-server-assignments"),
					ClientSession: TOSession,
					RequestBody: tc.ServerPutStatus{
						Status:        util.JSONNameOrIDStr{Name: util.Ptr("OFFLINE")},
						OfflineReason: util.Ptr("test last edge"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "PUT":
						t.Run(name, func(t *testing.T) {
							clearUpdates(t)
							alerts, reqInf, err := testCase.ClientSession.UpdateServerStatus(testCase.EndpointID(), testCase.RequestBody)
							for _, check := range testCase.Expectations {
								if alerts != nil {
									check(t, reqInf, nil, *alerts, err)
								}
							}
						})
					}
				}
			})
		}
	})
}

func validateUpdPending(hostName string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		descendants := make(map[string]struct{})

		params := url.Values{}
		params.Set("hostName", hostName)
		servers, _, err := TOSession.GetServersWithHdr(&params, nil)
		assert.RequireNoError(t, err, "Expected no error when getting servers: %v", err)
		assert.RequireEqual(t, 1, len(servers.Response), "Expected exactly one server returned from response, Got: %d", len(servers.Response))

		updatedServer := servers.Response[0]
		assert.RequireNotNil(t, updatedServer.CachegroupID, "Expected Server's CachegroupID to NOT be nil.")
		assert.RequireNotNil(t, updatedServer.Cachegroup, "Expected Server's Cachegroup to NOT be nil.")

		cacheGroups, _, err := TOSession.GetCacheGroupsNullableWithHdr(nil)
		assert.RequireNoError(t, err, "Expected no error when getting cache groups: %v", err)
		for _, cacheGroup := range cacheGroups {
			if cacheGroup.ParentCachegroupID != nil {
				if *cacheGroup.ParentCachegroupID == *servers.Response[0].CachegroupID {
					assert.RequireNotNil(t, cacheGroup.Name, "Expected Cachegroup's Name to NOT be nil.")
					descendants[*cacheGroup.Name] = struct{}{}
					if cacheGroup.SecondaryParentCachegroupID != nil {
						assert.RequireNotNil(t, cacheGroup.SecondaryParentName, "Expected Cachegroup's Secondary Parent's Name to NOT be nil.")
						descendants[*cacheGroup.SecondaryParentName] = struct{}{}
					}
				}
			}
			if cacheGroup.SecondaryParentCachegroupID != nil {
				if *cacheGroup.SecondaryParentCachegroupID == *servers.Response[0].CachegroupID {
					assert.RequireNotNil(t, cacheGroup.Name, "Expected Cachegroup's Name to NOT be nil.")
					descendants[*cacheGroup.Name] = struct{}{}
				}
			}
		}

		params.Del(hostName)
		allServers, _, err := TOSession.GetServersWithHdr(&params, nil)
		assert.RequireNoError(t, err, "Expected no error when getting servers: %v", err)
		for _, server := range allServers.Response {
			assert.RequireNotNil(t, server.HostName, "Expected Hostname to NOT be nil.")
			assert.RequireNotNil(t, server.Cachegroup, "Expected Cachegroup to NOT be nil.")
			assert.RequireNotNil(t, server.UpdPending, "Expected UpdPending to NOT be nil.")
			_, ok := descendants[*server.Cachegroup]
			if ok && *server.CDNName == *updatedServer.CDNName {
				assert.Equal(t, true, *server.UpdPending, "Expected server %s with cachegroup %s to have updates pending.", *server.HostName, *server.Cachegroup)
			} else {
				assert.Equal(t, false, *server.UpdPending, "Expected server %s with cachegroup %s to NOT have updates pending.", *server.HostName, *server.Cachegroup)
			}
		}
	}
}

func clearUpdates(t *testing.T) {
	cdns, _, err := TOSession.GetCDNsWithHdr(nil)
	assert.RequireNoError(t, err, "Error getting CDNs: %v", err)
	cdnQueueUpdate := tc.CDNQueueUpdateRequest{
		Action: "dequeue",
	}
	for _, cdn := range cdns {
		_, _, err := TOSession.QueueUpdatesForCDN(cdn.ID, cdnQueueUpdate)
		assert.RequireNoError(t, err, "Error Dequeing Updates for CDN %s: %v", cdn.Name, err)
	}
}

func validateUpdPendingSpecificServers(expected map[string]bool) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, _ tc.Alerts, _ error) {
		for expectedServer, expectedUpdPending := range expected {
			resp, _, err := TOSession.GetServerUpdateStatusWithHdr(expectedServer, nil)
			assert.RequireNoError(t, err, "Expected no error when getting server's update status: %v", err)
			actualUpdPending := resp.UpdatePending
			assert.Equal(t, expectedUpdPending, actualUpdPending, "Expected Update Pending for server: %s to be %t, Got: %t", expectedServer, expectedUpdPending, actualUpdPending)
		}
	}
}

func validateParentPendingSpecificServers(expected map[string]bool) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, _ tc.Alerts, _ error) {
		for expectedServer, expectedParentPending := range expected {
			resp, _, err := TOSession.GetServerUpdateStatusWithHdr(expectedServer, nil)
			assert.RequireNoError(t, err, "Expected no error when getting server's update status: %v", err)
			actualParentPending := resp.ParentPending
			assert.Equal(t, expectedParentPending, actualParentPending, "Expected Parent Pending for server: %s to be %t, Got: %t", expectedServer, expectedParentPending, actualParentPending)
		}
	}
}
