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
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestServersIDStatus(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServiceCategories, Topologies, ServerCapabilities, ServerServerCapabilities, DeliveryServices, DeliveryServiceServerAssignments}, func() {

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.ServerPutStatus]{
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
				"CONFLICT when SERVER STATUS OFFLINE when ONLY ORIGIN SERVER ASSIGNED": {
					EndpointID:    GetServerID(t, "test-mso-org-01"),
					ClientSession: TOSession,
					RequestBody: tc.ServerPutStatus{
						Status:        util.JSONNameOrIDStr{Name: util.Ptr("OFFLINE")},
						OfflineReason: util.Ptr("test last origin"),
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
							alerts, reqInf, err := testCase.ClientSession.UpdateServerStatus(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
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

func validateUpdPending(hostName string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		descendants := make(map[string]struct{})

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("hostName", hostName)
		servers, _, err := TOSession.GetServers(opts)
		assert.RequireNoError(t, err, "Expected no error when getting servers: %v", err)
		assert.RequireEqual(t, 1, len(servers.Response), "Expected exactly one server returned from response, Got: %d", len(servers.Response))

		updatedServer := servers.Response[0]

		opts.QueryParameters.Del("hostName")
		cacheGroups, _, err := TOSession.GetCacheGroups(opts)
		assert.RequireNoError(t, err, "Expected no error when getting cache groups: %v", err)
		for _, cacheGroup := range cacheGroups.Response {
			if cacheGroup.ParentCachegroupID != nil {
				if *cacheGroup.ParentCachegroupID == servers.Response[0].CacheGroupID {
					assert.RequireNotNil(t, cacheGroup.Name, "Expected Cachegroup's Name to NOT be nil.")
					descendants[*cacheGroup.Name] = struct{}{}
					if cacheGroup.SecondaryParentCachegroupID != nil {
						assert.RequireNotNil(t, cacheGroup.SecondaryParentName, "Expected Cachegroup's Secondary Parent's Name to NOT be nil.")
						descendants[*cacheGroup.SecondaryParentName] = struct{}{}
					}
				}
			}
			if cacheGroup.SecondaryParentCachegroupID != nil {
				if *cacheGroup.SecondaryParentCachegroupID == servers.Response[0].CacheGroupID {
					assert.RequireNotNil(t, cacheGroup.Name, "Expected Cachegroup's Name to NOT be nil.")
					descendants[*cacheGroup.Name] = struct{}{}
				}
			}
		}

		allServers, _, err := TOSession.GetServers(opts)
		assert.RequireNoError(t, err, "Expected no error when getting servers: %v", err)
		for _, server := range allServers.Response {
			_, ok := descendants[server.CacheGroup]
			if ok && server.CDN == updatedServer.CDN {
				assert.Equal(t, true, server.UpdatePending(), "Expected server %s with cachegroup %s to have updates pending.", server.HostName, server.CacheGroup)
			} else {
				assert.Equal(t, false, server.UpdatePending(), "Expected server %s with cachegroup %s to NOT have updates pending.", server.HostName, server.CacheGroup)
			}
		}
	}
}

func clearUpdates(t *testing.T) {
	cdns, _, err := TOSession.GetCDNs(client.RequestOptions{})
	assert.RequireNoError(t, err, "Error getting CDNs: %v", err)
	for _, cdn := range cdns.Response {
		_, _, err := TOSession.QueueUpdatesForCDN(cdn.ID, false, client.RequestOptions{})
		assert.RequireNoError(t, err, "Error Dequeing Updates for CDN %s: %v", cdn.Name, err)
	}
}

func validateUpdPendingSpecificServers(expected map[string]bool) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, _ tc.Alerts, _ error) {
		for expectedServer, expectedUpdPending := range expected {
			resp, _, err := TOSession.GetServerUpdateStatus(expectedServer, client.RequestOptions{})
			assert.RequireNoError(t, err, "Expected no error when getting server's update status: %v", err)
			assert.RequireEqual(t, 1, len(resp.Response), "Expected exactly one server's update status returned from response, Got: %d", len(resp.Response))
			actualUpdPending := resp.Response[0].UpdatePending
			assert.Equal(t, expectedUpdPending, actualUpdPending, "Expected Update Pending for server: %s to be %t, Got: %t", expectedServer, expectedUpdPending, actualUpdPending)
		}
	}
}

func validateParentPendingSpecificServers(expected map[string]bool) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, _ tc.Alerts, _ error) {
		for expectedServer, expectedParentPending := range expected {
			resp, _, err := TOSession.GetServerUpdateStatus(expectedServer, client.RequestOptions{})
			assert.RequireNoError(t, err, "Expected no error when getting server's update status: %v", err)
			assert.RequireEqual(t, 1, len(resp.Response), "Expected exactly one server's update status returned from response, Got: %d", len(resp.Response))
			actualParentPending := resp.Response[0].ParentPending
			assert.Equal(t, expectedParentPending, actualParentPending, "Expected Parent Pending for server: %s to be %t, Got: %t", expectedServer, expectedParentPending, actualParentPending)
		}
	}
}
