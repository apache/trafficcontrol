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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestServerUpdateStatus(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers}, func() {
		//TODO: DON'T hard-code server hostnames!
		var edge1cdn1 tc.ServerNullable
		var edge2cdn1 tc.ServerNullable
		var mid1cdn1 tc.ServerNullable
		var edge1cdn2 tc.ServerNullable

		params := url.Values{}

		getServers := func() {
			for _, s := range []struct {
				name   string
				server *tc.ServerNullable
			}{
				{
					"atlanta-edge-01",
					&edge1cdn1,
				},
				{
					"atlanta-edge-03",
					&edge2cdn1,
				},
				{
					"atlanta-mid-16",
					&mid1cdn1,
				},
				{
					"edge1-cdn2",
					&edge1cdn2,
				},
			} {
				params.Set("hostName", s.name)
				resp, _, err := TOSession.GetServers(&params)
				if err != nil {
					t.Errorf("cannot GET Server by hostname '%s': %v - %v", s.name, err, resp.Alerts)
				}
				if len(resp.Response) < 1 {
					t.Fatalf("Expected a server named '%s' to exist", s.name)
				}
				if len(resp.Response) > 1 {
					t.Errorf("Expected exactly one server named '%s' to exist - actual: %d", s.name, len(resp.Response))
					t.Logf("Testing will proceed with server: %+v", resp.Response[0])
				}
				*s.server = resp.Response[0]
				if s.server.ID == nil {
					t.Fatalf("server '%s' was returned with nil ID", s.name)
				}
				if s.server.HostName == nil {
					t.Fatalf("server '%s' was returned with nil HostName", s.name)
				}
			}
		}
		getServers()

		// assert that servers don't have updates pending
		for _, s := range []tc.ServerNullable{
			edge1cdn1,
			edge2cdn1,
			mid1cdn1,
			edge1cdn2,
		} {
			if s.UpdPending == nil {
				t.Error("expected UpdPending: false, actual: null")
			} else if *s.UpdPending {
				t.Error("expected UpdPending: false, actual: true")
			}
		}

		// update status of MID server to OFFLINE
		_, _, err := TOSession.UpdateServerStatus(*mid1cdn1.ID, tc.ServerPutStatus{
			Status:        util.JSONNameOrIDStr{Name: util.StrPtr("OFFLINE")},
			OfflineReason: util.StrPtr("testing")})
		if err != nil {
			t.Errorf("cannot update server status: %v", err)
		}

		// assert that updates were queued for the proper EDGE servers
		getServers()
		if edge1cdn1.UpdPending == nil {
			t.Errorf("expected: child %s (%d) to have updates pending, actual: property was null (or missing)", *edge1cdn1.HostName, *edge1cdn1.ID)
		} else if !*edge1cdn1.UpdPending {
			t.Errorf("expected: child %s (%d) to have updates pending, actual: no updates pending", *edge1cdn1.HostName, *edge1cdn1.ID)
		}

		if edge2cdn1.UpdPending == nil {
			t.Errorf("expected: child %s (%d) to have updates pending, actual: property was null (or missing)", *edge2cdn1.HostName, *edge2cdn1.ID)
		} else if !*edge2cdn1.UpdPending {
			t.Errorf("expected: child %s (%d) to have updates pending, actual: no updates pending", *edge2cdn1.HostName, *edge2cdn1.ID)
		}
		if mid1cdn1.UpdPending == nil {
			t.Errorf("expected: server %s (%d) with updated status to have no updates pending, actual: property was null (or missing)", *mid1cdn1.HostName, *mid1cdn1.ID)
		} else if *mid1cdn1.UpdPending {
			t.Errorf("expected: server %s (%d) with updated status to have no updates pending, actual: updates pending", *mid1cdn1.HostName, *mid1cdn1.ID)
		}

		if edge1cdn2.UpdPending == nil {
			t.Errorf("expected: server %s (%d) in different CDN than server with updated status to have no updates pending, actual: updates pending", *edge1cdn2.HostName, *edge1cdn2.ID)
		} else if *edge1cdn2.UpdPending {
			t.Errorf("expected: server %s (%d) in different CDN than server with updated status to have no updates pending, actual: updates pending", *edge1cdn2.HostName, *edge1cdn2.ID)
		}

		// update status of MID server to OFFLINE via status ID
		status, _, err := TOSession.GetStatusByName("OFFLINE")
		if err != nil {
			t.Fatalf("cannot GET status by name: %v", err)
		}
		_, _, err = TOSession.UpdateServerStatus(
			*mid1cdn1.ID,
			tc.ServerPutStatus{
				Status:        util.JSONNameOrIDStr{ID: util.IntPtr(status[0].ID)},
				OfflineReason: util.StrPtr("testing"),
			},
		)
		if err != nil {
			t.Errorf("cannot update server status: %v", err)
		}

		// negative cases:
		// server doesn't exist
		_, _, err = TOSession.UpdateServerStatus(
			-1,
			tc.ServerPutStatus{
				Status:        util.JSONNameOrIDStr{Name: util.StrPtr("OFFLINE")},
				OfflineReason: util.StrPtr("testing"),
			},
		)
		if err == nil {
			t.Error("update server status exected: err, actual: nil")
		}

		// status does not exist
		_, _, err = TOSession.UpdateServerStatus(
			*mid1cdn1.ID,
			tc.ServerPutStatus{
				Status:        util.JSONNameOrIDStr{Name: util.StrPtr("NOT_A_REAL_STATUS")},
				OfflineReason: util.StrPtr("testing"),
			},
		)
		if err == nil {
			t.Error("update server status exected: err, actual: nil")
		}

		// offlineReason required for OFFLINE status
		_, _, err = TOSession.UpdateServerStatus(
			*mid1cdn1.ID,
			tc.ServerPutStatus{
				Status:        util.JSONNameOrIDStr{Name: util.StrPtr("OFFLINE")},
				OfflineReason: nil,
			},
		)
		if err == nil {
			t.Error("update server status exected: err, actual: nil")
		}

		// offlineReason required for ADMIN_DOWN status
		_, _, err = TOSession.UpdateServerStatus(
			*mid1cdn1.ID,
			tc.ServerPutStatus{
				Status:        util.JSONNameOrIDStr{Name: util.StrPtr("ADMIN_DOWN")},
				OfflineReason: nil,
			},
		)
		if err == nil {
			t.Error("update server status exected: err, actual: nil")
		}
	})
}

func TestServerQueueUpdate(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers}, func() {
		// TODO: DON'T hard-code server hostnames!
		const serverName = "atlanta-edge-01"

		queueUpdateActions := map[bool]string{
			false: "dequeue",
			true:  "queue",
		}

		var s tc.ServerNullable
		params := url.Values{}
		params.Add("hostName", serverName)
		resp, _, err := TOSession.GetServers(&params)
		if err != nil {
			t.Fatalf("failed to GET Server by hostname '%s': %v - %v", serverName, err, resp.Alerts)
		}
		if len(resp.Response) < 1 {
			t.Fatalf("Expected a server named '%s' to exist", serverName)
		}
		if len(resp.Response) > 1 {
			t.Errorf("Expected exactly one server named '%s' to exist", serverName)
			t.Logf("Testing will proceed with server: %+v", resp.Response[0])
		}
		s = resp.Response[0]

		// assert that servers don't have updates pending
		if s.UpdPending == nil {
			t.Fatalf("Server '%s' had null (or missing) updPending property", serverName)
		}
		if got, want := *s.UpdPending, false; got != want {
			t.Fatalf("unexpected UpdPending, got: %v, want: %v", got, want)
		}

		if s.ID == nil {
			t.Fatalf("Server '%s' had nil ID", serverName)
		}

		for _, setVal := range [...]bool{true, false} {
			t.Run(fmt.Sprint(setVal), func(t *testing.T) {
				// queue update and check response
				quResp, _, err := TOSession.SetServerQueueUpdate(*s.ID, setVal)
				if err != nil {
					t.Fatalf("failed to set queue update for server with ID %v to %v: %v", s.ID, setVal, err)
				}
				if got, want := int(quResp.Response.ServerID), *s.ID; got != want {
					t.Errorf("wrong serverId in response, got: %v, want: %v", got, want)
				}
				if got, want := quResp.Response.Action, queueUpdateActions[setVal]; got != want {
					t.Errorf("wrong action in response, got: %v, want: %v", got, want)
				}

				// assert that the server has updates queued
				resp, _, err = TOSession.GetServers(&params)
				if err != nil {
					t.Fatalf("failed to GET Server by hostname '%s': %v - %v", serverName, err, resp.Alerts)
				}
				if len(resp.Response) < 1 {
					t.Fatalf("Expected a server named '%s' to exist", serverName)
				}
				if len(resp.Response) > 1 {
					t.Errorf("Expected exactly one server named '%s' to exist", serverName)
					t.Logf("Testing will proceed with server: %+v", resp.Response[0])
				}
				s = resp.Response[0]
				if s.UpdPending == nil {
					t.Fatalf("Server '%s' had null (or missing) updPending property", serverName)
				}
				if got, want := *s.UpdPending, setVal; got != want {
					t.Errorf("unexpected UpdPending, got: %v, want: %v", got, want)
				}
			})
		}

		t.Run("validations", func(t *testing.T) {
			// server doesn't exist
			_, _, err = TOSession.SetServerQueueUpdate(-1, true)
			if err == nil {
				t.Error("update server status expected: error, actual: nil")
			}

			// invalid action
			req, err := json.Marshal(tc.ServerQueueUpdateRequest{Action: "foobar"})
			if err != nil {
				t.Fatalf("failed to encode request body: %v", err)
			}

			// TODO: don't construct URLs like this, nor use "RawRequest"
			path := fmt.Sprintf(TestAPIBase+"/servers/%d/queue_update", *s.ID)
			httpResp, _, err := TOSession.RawRequest(http.MethodPost, path, req)
			if err != nil {
				t.Fatalf("POST request failed: %v", err)
			}
			if httpResp.StatusCode >= 200 && httpResp.StatusCode <= 299 {
				t.Errorf("unexpected status code: got %v, want something outside the range [200, 299]", httpResp.StatusCode)
			}
		})
	})
}

func TestSetServerUpdateStatuses(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers}, func() {
		if len(testData.Servers) < 1 {
			t.Fatal("cannot GET Server: no test data")
		}
		testServer := testData.Servers[0]
		if testServer.HostName == nil {
			t.Fatalf("First test server had nil hostname: %+v", testServer)
		}

		params := url.Values{}
		params.Add("hostName", *testServer.HostName)
		testVals := func(queue *bool, reval *bool) {
			resp, _, err := TOSession.GetServers(&params)
			if err != nil {
				t.Errorf("cannot GET Server by name '%s': %v - %v", *testServer.HostName, err, resp.Alerts)
			} else if len(resp.Response) != 1 {
				t.Fatalf("GET Server expected 1, actual %v", len(resp.Response))
			}

			existingServer := resp.Response

			if existingServer[0].UpdPending == nil {
				t.Fatalf("Server '%s' had nil UpdPending before update status change", *testServer.HostName)
			}
			if existingServer[0].RevalPending == nil {
				t.Fatalf("Server '%s' had nil RevalPending before update status change", *testServer.HostName)
			}

			if _, err := TOSession.SetUpdateServerStatuses(*testServer.HostName, queue, reval); err != nil {
				t.Fatalf("UpdateServerStatuses error expected: nil, actual: %v", err)
			}

			resp, _, err = TOSession.GetServers(&params)
			if err != nil {
				t.Errorf("cannot GET Server by name '%s': %v - %v", *testServer.HostName, err, resp.Alerts)
			} else if len(resp.Response) != 1 {
				t.Fatalf("GET Server expected 1, actual %v", len(resp.Response))
			}

			newServer := resp.Response

			if newServer[0].UpdPending == nil {
				t.Fatalf("Server '%s' had nil UpdPending after update status change", *testServer.HostName)
			}
			if newServer[0].RevalPending == nil {
				t.Fatalf("Server '%s' had nil RevalPending after update status change", *testServer.HostName)
			}

			if queue != nil {
				if *newServer[0].UpdPending != *queue {
					t.Errorf("set queue update pending to %v, but then got server %v", *queue, *newServer[0].UpdPending)
				}
			} else {
				if *newServer[0].UpdPending != *existingServer[0].UpdPending {
					t.Errorf("set queue update pending with nil (don't update), but then got server %v which didn't match pre-update value %v", *newServer[0].UpdPending, *existingServer[0].UpdPending)
				}
			}
			if reval != nil {
				if *newServer[0].RevalPending != *reval {
					t.Errorf("set reval update pending to %v, but then got server %v", *reval, *newServer[0].RevalPending)
				}
			} else {
				if *newServer[0].RevalPending != *existingServer[0].RevalPending {
					t.Errorf("set reval update pending with nil (don't update), but then got server %v which didn't match pre-update value %v", *newServer[0].RevalPending, *existingServer[0].RevalPending)
				}
			}
		}

		testVals(util.BoolPtr(true), util.BoolPtr(true))
		testVals(util.BoolPtr(true), util.BoolPtr(false))
		testVals(util.BoolPtr(false), util.BoolPtr(false))
		testVals(nil, util.BoolPtr(true))
		testVals(nil, util.BoolPtr(false))
		testVals(util.BoolPtr(true), nil)
		testVals(util.BoolPtr(false), nil)

		if _, err := TOSession.SetUpdateServerStatuses(*testServer.HostName, nil, nil); err == nil {
			t.Errorf("UpdateServerStatuses with (nil,nil) expected error, actual nil")
		}
	})
}

func TestSetTopologiesServerUpdateStatuses(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies}, func() {
		const (
			topologyName        = "forked-topology"
			edgeCacheGroup      = "topology-edge-cg-01"
			otherEdgeCacheGroup = "topology-edge-cg-02"
			midCacheGroup       = "topology-mid-cg-04"
		)
		cacheGroupNames := []string{edgeCacheGroup, otherEdgeCacheGroup, midCacheGroup}
		cachesByCacheGroup := map[string]tc.ServerNullable{}
		updateStatusByCacheGroup := map[string]tc.ServerUpdateStatus{}

		forkedTopology, _, err := TOSession.GetTopology(topologyName)
		if err != nil {
			t.Fatalf("topology %s was not found", topologyName)
		}
		for _, cacheGroupName := range cacheGroupNames {
			foundNode := false
			for _, node := range forkedTopology.Nodes {
				if node.Cachegroup == cacheGroupName {
					foundNode = true
					break
				}
			}
			if !foundNode {
				t.Fatalf("unable to find topology node with cachegroup %s", cacheGroupName)
			}

			cacheGroups, _, err := TOSession.GetCacheGroupNullableByName(cacheGroupName)
			if err != nil {
				t.Fatalf("unable to get cachegroup %s: %s", cacheGroupName, err.Error())
			}
			if len(cacheGroups) != 1 {
				t.Fatalf("incorrect number of cachegroups. expected: %d actual: %d", 1, len(cacheGroups))
			}
			cacheGroup := cacheGroups[0]

			params := url.Values{"cachegroup": []string{strconv.Itoa(*cacheGroup.ID)}}
			cachesByCacheGroup[cacheGroupName], _, err = TOSession.GetFirstServer(&params)
			if err != nil {
				t.Fatalf("unable to get a server from cachegroup %s: %s", cacheGroupName, err.Error())
			}
		}

		// update status of MID server to OFFLINE
		_, _, err = TOSession.UpdateServerStatus(*cachesByCacheGroup[midCacheGroup].ID, tc.ServerPutStatus{
			Status:        util.JSONNameOrIDStr{Name: util.StrPtr("OFFLINE")},
			OfflineReason: util.StrPtr("testing")})
		if err != nil {
			t.Fatalf("cannot update server status: %s", err.Error())
		}

		for _, cacheGroupName := range cacheGroupNames {
			params := url.Values{"cachegroup": []string{strconv.Itoa(*cachesByCacheGroup[cacheGroupName].CachegroupID)}}
			cachesByCacheGroup[cacheGroupName], _, err = TOSession.GetFirstServer(&params)
			if err != nil {
				t.Fatalf("unable to get a server from cachegroup %s: %s", cacheGroupName, err.Error())
			}
		}
		for _, cacheGroupName := range cacheGroupNames {
			updateStatusByCacheGroup[cacheGroupName], _, err = TOSession.GetServerUpdateStatus(*cachesByCacheGroup[cacheGroupName].HostName)
			if err != nil {
				t.Fatalf("unable to get a server from cachegroup %s: %s", cacheGroupName, err.Error())
			}
		}
		// updating the server status does not queue updates within the same cachegroup
		if *cachesByCacheGroup[midCacheGroup].UpdPending {
			t.Fatalf("expected UpdPending: %t, actual: %t", false, *cachesByCacheGroup[midCacheGroup].UpdPending)
		}
		// edgeCacheGroup is a descendant of midCacheGroup
		if !*cachesByCacheGroup[edgeCacheGroup].UpdPending {
			t.Fatalf("expected UpdPending: %t, actual: %t", true, *cachesByCacheGroup[edgeCacheGroup].UpdPending)
		}
		if !updateStatusByCacheGroup[edgeCacheGroup].UpdatePending {
			t.Fatalf("expected UpdPending: %t, actual: %t", true, updateStatusByCacheGroup[edgeCacheGroup].UpdatePending)
		}
		// otherEdgeCacheGroup is not a descendant of midCacheGroup but is still in the same topology
		if *cachesByCacheGroup[otherEdgeCacheGroup].UpdPending {
			t.Fatalf("expected UpdPending: %t, actual: %t", false, *cachesByCacheGroup[otherEdgeCacheGroup].UpdPending)
		}
		if updateStatusByCacheGroup[otherEdgeCacheGroup].UpdatePending {
			t.Fatalf("expected UpdPending: %t, actual: %t", false, updateStatusByCacheGroup[otherEdgeCacheGroup].UpdatePending)
		}

		_, _, err = TOSession.SetServerQueueUpdate(*cachesByCacheGroup[midCacheGroup].ID, true)
		if err != nil {
			t.Fatalf("cannot update server status on %s: %s", *cachesByCacheGroup[midCacheGroup].HostName, err.Error())
		}
		for _, cacheGroupName := range cacheGroupNames {
			updateStatusByCacheGroup[cacheGroupName], _, err = TOSession.GetServerUpdateStatus(*cachesByCacheGroup[cacheGroupName].HostName)
			if err != nil {
				t.Fatalf("unable to get a server from cachegroup %s: %s", cacheGroupName, err.Error())
			}
		}

		// edgeCacheGroup is a descendant of midCacheGroup
		if !updateStatusByCacheGroup[edgeCacheGroup].ParentPending {
			t.Fatalf("expected UpdPending: %t, actual: %t", true, updateStatusByCacheGroup[edgeCacheGroup].ParentPending)
		}
		// otherEdgeCacheGroup is not a descendant of midCacheGroup but is still in the same topology
		if updateStatusByCacheGroup[otherEdgeCacheGroup].ParentPending {
			t.Fatalf("expected UpdPending: %t, actual: %t", false, updateStatusByCacheGroup[otherEdgeCacheGroup].ParentPending)
		}

		edgeHostName := *cachesByCacheGroup[edgeCacheGroup].HostName
		*cachesByCacheGroup[edgeCacheGroup].HostName = *cachesByCacheGroup[midCacheGroup].HostName
		_, _, err = TOSession.UpdateServerByID(*cachesByCacheGroup[edgeCacheGroup].ID, cachesByCacheGroup[edgeCacheGroup])
		if err != nil {
			t.Fatalf("unable to update %s's hostname to %s: %s", edgeHostName, *cachesByCacheGroup[midCacheGroup].HostName, err)
		}

		_, _, err = TOSession.GetServerUpdateStatus(*cachesByCacheGroup[midCacheGroup].HostName)
		if err != nil {
			t.Fatalf("expected no error getting server updates for a non-unique hostname %s, got %s", *cachesByCacheGroup[midCacheGroup].HostName, err)
		}

		*cachesByCacheGroup[edgeCacheGroup].HostName = edgeHostName
		_, _, err = TOSession.UpdateServerByID(*cachesByCacheGroup[edgeCacheGroup].ID, cachesByCacheGroup[edgeCacheGroup])
		if err != nil {
			t.Fatalf("unable to revert %s's hostname back to %s: %s", edgeHostName, edgeHostName, err)
		}
	})
}
