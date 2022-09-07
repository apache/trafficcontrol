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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestServerUpdateApplyTimeLocked(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServiceCategories, Topologies, DeliveryServices}, func() {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("hostName", "atlanta-edge-01")
		resp, _, err := TOSession.GetServers(opts)
		if err != nil {
			t.Fatalf("cannot get server by hostname: %v", err)
		}
		if len(resp.Response) != 1 {
			t.Fatalf("Expected a server named 'atlanta-edge-01' to exist")
		}
		testServer := resp.Response[0]

		testVals := func(configApply, revalApply *time.Time) {
			resp, _, err := TOSession.GetServers(opts)
			if err != nil {
				t.Errorf("cannot get Server by name '%s': %v - alerts: %+v", *testServer.HostName, err, resp.Alerts)
			} else if len(resp.Response) != 1 {
				t.Fatalf("GET Server expected 1, actual %v", len(resp.Response))
			}

			beforeServer := resp.Response[0]

			// Ensure baseline
			if beforeServer.UpdPending == nil {
				t.Fatalf("Server '%s' had nil UpdPending before update status change", *testServer.HostName)
			}
			if beforeServer.RevalPending == nil {
				t.Fatalf("Server '%s' had nil RevalPending before update status change", *testServer.HostName)
			}
			if beforeServer.ConfigUpdateTime == nil {
				t.Fatalf("Server '%s' had nil ConfigUpdateTime before update status change", *testServer.HostName)
			}
			if beforeServer.ConfigApplyTime == nil {
				t.Fatalf("Server '%s' had nil ConfigApplyTime before update status change", *testServer.HostName)
			}
			if beforeServer.RevalUpdateTime == nil {
				t.Fatalf("Server '%s' had nil RevalUpdateTime before update status change", *testServer.HostName)
			}
			if beforeServer.RevalApplyTime == nil {
				t.Fatalf("Server '%s' had nil RevalApplyTime before update status change", *testServer.HostName)
			}

			_, _, err = TOSession.CreateCDNLock(tc.CDNLock{CDN: *testServer.CDNName, Soft: util.BoolPtr(false)}, client.NewRequestOptions())
			if err != nil {
				t.Errorf("cannont acquire lock on the Server '%s' CDN '%s: %v", *testServer.HostName, *testServer.CDNName, err)
			}

			opsSession, _, err := client.LoginWithAgent(Config.TrafficOps.URL, Config.TrafficOps.Users.Operations,
				Config.TrafficOps.UserPassword, true, "to-api-v4-locks-ops",
				true, time.Second*time.Duration(Config.Default.Session.TimeoutInSecs))
			if err != nil {
				t.Errorf("cannot login as '%s' user: %v", Config.TrafficOps.Users.Operations, err)
			}

			badQueryOpts := client.NewRequestOptions()
			badQueryOpts.QueryParameters.Set("updated", "true")
			if _, _, err := opsSession.SetUpdateServerStatusTimes(*testServer.HostName, configApply, revalApply, badQueryOpts); err == nil {
				t.Fatalf("expected SetUpdateServerStatusTimes to not allow updates due to `updated` field")
			}

			badQueryOpts = client.NewRequestOptions()
			badQueryOpts.QueryParameters.Set("reval_updated", "true")
			if _, _, err := opsSession.SetUpdateServerStatusTimes(*testServer.HostName, configApply, revalApply, badQueryOpts); err == nil {
				t.Fatalf("expected SetUpdateServerStatusTimes to not allow updates due to `reval_updated` field")
			}

			if alerts, _, err := opsSession.SetUpdateServerStatusTimes(*testServer.HostName, configApply, revalApply, client.RequestOptions{}); err != nil {
				t.Fatalf("SetUpdateServerStatusTimes error. expected: nil, actual: %v - alerts: %+v", err, alerts.Alerts)
			}

			resp, _, err = opsSession.GetServers(opts)
			if err != nil {
				t.Errorf("cannot GET Server by name '%s': %v - alerts: %+v", *testServer.HostName, err, resp.Alerts)
			} else if len(resp.Response) != 1 {
				t.Fatalf("GET Server expected 1, actual %v", len(resp.Response))
			}
			afterServer := resp.Response[0]

			opts := client.NewRequestOptions()
			opts.QueryParameters.Set("cdn", *testServer.CDNName)
			_, _, err = TOSession.DeleteCDNLocks(opts)
			if err != nil {
				t.Errorf("cannont delete acquired lock on the Server '%s' CDN '%s: %v", *testServer.HostName, *testServer.CDNName, err)
			}

			if afterServer.UpdPending == nil {
				t.Fatalf("Server '%s' had nil UpdPending after update status change", *testServer.HostName)
			}
			if afterServer.RevalPending == nil {
				t.Fatalf("Server '%s' had nil RevalPending after update status change", *testServer.HostName)
			}

			// Ensure values were actually set
			if configApply != nil {
				if afterServer.ConfigApplyTime == nil || !afterServer.ConfigApplyTime.Equal(*configApply) {
					t.Errorf("Failed to set server's ConfigApplyTime. expected: %v actual: %v", *configApply, afterServer.ConfigApplyTime)
				}
			}
			if revalApply != nil {
				if afterServer.RevalApplyTime == nil || !afterServer.RevalApplyTime.Equal(*revalApply) {
					t.Errorf("Failed to set server's RevalApplyTime. expected: %v actual: %v", *revalApply, afterServer.RevalApplyTime)
				}
			}

		}

		now := time.Now().Round(time.Microsecond)
		testVals(util.TimePtr(now), nil)
		testVals(nil, util.TimePtr(now))
	})
}

func TestServerUpdateStatusLastAssigned(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServiceCategories, Topologies, DeliveryServices}, func() {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("hostName", "atlanta-edge-01")
		resp, _, err := TOSession.GetServers(opts)
		if err != nil {
			t.Fatalf("cannot get server by hostname: %v", err)
		}
		if len(resp.Response) != 1 {
			t.Fatalf("Expected a server named 'atlanta-edge-01' to exist")
		}
		edge := resp.Response[0]
		opts = client.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", "ds-top")
		dsResp, _, err := TOSession.GetDeliveryServices(opts)
		if err != nil {
			t.Fatalf("cannot get delivery service by xmlId: %v", err)
		}
		if len(resp.Response) != 1 {
			t.Fatalf("Expected one delivery service with xmlId 'ds-top' to exist")
		}
		// temporarily unassign the topology in order to assign an EDGE
		ds := dsResp.Response[0]
		tmpTop := *ds.Topology
		ds.Topology = nil
		ds.FirstHeaderRewrite = nil
		ds.LastHeaderRewrite = nil
		ds.InnerHeaderRewrite = nil
		_, _, err = TOSession.UpdateDeliveryService(*ds.ID, ds, client.RequestOptions{})
		if err != nil {
			t.Fatalf("cannot update delivery service 'ds-top': %v", err)
		}
		_, _, err = TOSession.CreateDeliveryServiceServers(*ds.ID, []int{*edge.ID}, true, client.RequestOptions{})
		if err != nil {
			t.Fatalf("cannot create delivery service server: %v", err)
		}
		// reassign the topology
		ds.Topology = &tmpTop
		_, _, err = TOSession.UpdateDeliveryService(*ds.ID, ds, client.RequestOptions{})
		if err != nil {
			t.Fatalf("cannot update delivery service 'ds-top': %v", err)
		}
		// attempt to set the edge to OFFLINE
		_, _, err = TOSession.UpdateServerStatus(*edge.ID, tc.ServerPutStatus{
			Status:        util.JSONNameOrIDStr{Name: util.StrPtr("OFFLINE")},
			OfflineReason: util.StrPtr("testing")}, client.RequestOptions{})
		if err != nil {
			t.Errorf("setting edge to OFFLINE when it's the only edge assigned to a topology-based delivery service - expected: no error, actual: %v", err)
		}
		// remove EDGE assignment
		_, _, err = TOSession.CreateDeliveryServiceServers(*ds.ID, []int{}, true, client.RequestOptions{})
		if err != nil {
			t.Errorf("removing delivery service servers from topology-based delivery service - expected: no error, actual: %v", err)
		}
	})
}

func TestServerUpdateStatus(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers}, func() {
		//TODO: DON'T hard-code server hostnames!
		var edge1cdn1 tc.ServerV4
		var edge2cdn1 tc.ServerV4
		var mid1cdn1 tc.ServerV4
		var edge1cdn2 tc.ServerV4

		opts := client.NewRequestOptions()

		getServers := func() {
			for _, s := range []struct {
				name   string
				server *tc.ServerV40
			}{
				{
					"atlanta-edge-01",
					&edge1cdn1.ServerV40,
				},
				{
					"atlanta-edge-03",
					&edge2cdn1.ServerV40,
				},
				{
					"atlanta-mid-16",
					&mid1cdn1.ServerV40,
				},
				{
					"edge1-cdn2",
					&edge1cdn2.ServerV40,
				},
			} {
				opts.QueryParameters.Set("hostName", s.name)
				resp, _, err := TOSession.GetServers(opts)
				if err != nil {
					t.Errorf("cannot get Server by hostname '%s': %v - alerts: %+v", s.name, err, resp.Alerts)
				}
				if len(resp.Response) < 1 {
					t.Fatalf("Expected a server named '%s' to exist", s.name)
				}
				if len(resp.Response) > 1 {
					t.Errorf("Expected exactly one server named '%s' to exist - actual: %d", s.name, len(resp.Response))
					t.Logf("Testing will proceed with server: %+v", resp.Response[0])
				}
				*s.server = resp.Response[0].ServerV40
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
		for _, s := range []tc.ServerV4{
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
		alerts, _, err := TOSession.UpdateServerStatus(*mid1cdn1.ID, tc.ServerPutStatus{
			Status:        util.JSONNameOrIDStr{Name: util.StrPtr("OFFLINE")},
			OfflineReason: util.StrPtr("testing")}, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot update server status: %v - alerts: %+v", err, alerts)
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
		opts = client.NewRequestOptions()
		opts.QueryParameters.Set("name", "OFFLINE")
		status, _, err := TOSession.GetStatuses(opts)
		if err != nil {
			t.Fatalf("cannot get Status 'OFFLINE': %v - alerts: %+v", err, status.Alerts)
		}
		if len(status.Response) != 1 {
			t.Fatalf("Expected exactly one Status to exist with name 'OFFLINE', found: %d", len(status.Response))
		}
		alerts, _, err = TOSession.UpdateServerStatus(
			*mid1cdn1.ID,
			tc.ServerPutStatus{
				Status:        util.JSONNameOrIDStr{ID: util.IntPtr(status.Response[0].ID)},
				OfflineReason: util.StrPtr("testing"),
			},
			client.RequestOptions{},
		)
		if err != nil {
			t.Errorf("cannot update server status: %v - alerts: %+v", err, alerts.Alerts)
		}

		// negative cases:
		// server doesn't exist
		_, _, err = TOSession.UpdateServerStatus(
			-1,
			tc.ServerPutStatus{
				Status:        util.JSONNameOrIDStr{Name: util.StrPtr("OFFLINE")},
				OfflineReason: util.StrPtr("testing"),
			},
			client.RequestOptions{},
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
			client.RequestOptions{},
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
			client.RequestOptions{},
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
			client.RequestOptions{},
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

		var s tc.ServerV4
		opts := client.NewRequestOptions()
		opts.QueryParameters.Add("hostName", serverName)
		resp, _, err := TOSession.GetServers(opts)
		if err != nil {
			t.Fatalf("failed to get Server by hostname '%s': %v - alerts: %+v", serverName, err, resp.Alerts)
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
				quResp, _, err := TOSession.SetServerQueueUpdate(*s.ID, setVal, client.RequestOptions{})
				if err != nil {
					t.Fatalf("failed to set queue update for server with ID %d to %t: %v - alerts: %+v", s.ID, setVal, err, quResp.Alerts)
				}
				if got, want := int(quResp.Response.ServerID), *s.ID; got != want {
					t.Errorf("wrong serverId in response, got: %v, want: %v", got, want)
				}
				if got, want := quResp.Response.Action, queueUpdateActions[setVal]; got != want {
					t.Errorf("wrong action in response, got: %v, want: %v", got, want)
				}

				// assert that the server has updates queued
				resp, _, err = TOSession.GetServers(opts)
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
			_, _, err = TOSession.SetServerQueueUpdate(-1, true, client.RequestOptions{})
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

		opts := client.NewRequestOptions()
		opts.QueryParameters.Add("hostName", *testServer.HostName)
		testVals := func(configApply, revalApply *time.Time) {
			resp, _, err := TOSession.GetServers(opts)
			if err != nil {
				t.Errorf("cannot get Server by name '%s': %v - alerts: %+v", *testServer.HostName, err, resp.Alerts)
			} else if len(resp.Response) != 1 {
				t.Fatalf("GET Server expected 1, actual %v", len(resp.Response))
			}

			beforeServer := resp.Response[0]

			// Ensure baseline
			if beforeServer.UpdPending == nil {
				t.Fatalf("Server '%s' had nil UpdPending before update status change", *testServer.HostName)
			}
			if beforeServer.RevalPending == nil {
				t.Fatalf("Server '%s' had nil RevalPending before update status change", *testServer.HostName)
			}
			if beforeServer.ConfigUpdateTime == nil {
				t.Fatalf("Server '%s' had nil ConfigUpdateTime before update status change", *testServer.HostName)
			}
			if beforeServer.ConfigApplyTime == nil {
				t.Fatalf("Server '%s' had nil ConfigApplyTime before update status change", *testServer.HostName)
			}
			if beforeServer.RevalUpdateTime == nil {
				t.Fatalf("Server '%s' had nil RevalUpdateTime before update status change", *testServer.HostName)
			}
			if beforeServer.RevalApplyTime == nil {
				t.Fatalf("Server '%s' had nil RevalApplyTime before update status change", *testServer.HostName)
			}

			// Make change
			if alerts, _, err := TOSession.SetUpdateServerStatusTimes(*testServer.HostName, configApply, revalApply, client.RequestOptions{}); err != nil {
				t.Fatalf("SetUpdateServerStatusTimes error. expected: nil, actual: %v - alerts: %+v", err, alerts.Alerts)
			}

			resp, _, err = TOSession.GetServers(opts)
			if err != nil {
				t.Errorf("cannot GET Server by name '%s': %v - alerts: %+v", *testServer.HostName, err, resp.Alerts)
			} else if len(resp.Response) != 1 {
				t.Fatalf("GET Server expected 1, actual %v", len(resp.Response))
			}

			afterServer := resp.Response[0]

			if afterServer.UpdPending == nil {
				t.Fatalf("Server '%s' had nil UpdPending after update status change", *testServer.HostName)
			}
			if afterServer.RevalPending == nil {
				t.Fatalf("Server '%s' had nil RevalPending after update status change", *testServer.HostName)
			}

			// Ensure values were actually set
			if configApply != nil {
				if afterServer.ConfigApplyTime == nil || !afterServer.ConfigApplyTime.Equal(*configApply) {
					t.Errorf("Failed to set server's ConfigApplyTime. expected: %v actual: %v", *configApply, afterServer.ConfigApplyTime)
				}
			}
			if revalApply != nil {
				if afterServer.RevalApplyTime == nil || !afterServer.RevalApplyTime.Equal(*revalApply) {
					t.Errorf("Failed to set server's RevalApplyTime. expected: %v actual: %v", *revalApply, afterServer.RevalApplyTime)
				}
			}

		}

		// Postgres stores microsecond precision. There is also some discussion around MacOS losing
		// precision as well. The nanosecond precision is accurate within go one linux however,
		// but round trips to and from the database may result in an inaccurate Equals comparison
		// with the loss of precision. Also, it appears to Round and not Truncate.
		now := time.Now().Round(time.Microsecond)

		// Test setting the values works as expected
		testVals(util.TimePtr(now), nil) // configApply
		testVals(nil, util.TimePtr(now)) // revalApply

		// Test sending all nils. Should fail
		if _, _, err := TOSession.SetUpdateServerStatusTimes(*testServer.HostName, nil, nil, client.RequestOptions{}); err == nil {
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
		cachesByCDNCacheGroup := make(map[string]map[string][]tc.ServerV4)
		updateStatusByCacheGroup := map[string]tc.ServerUpdateStatusV40{}

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", topologyName)
		forkedTopology, _, err := TOSession.GetTopologies(opts)
		if err != nil {
			t.Fatalf("Topology '%s' was not found: %v - alerts: %+v", topologyName, err, forkedTopology.Alerts)
		}
		if len(forkedTopology.Response) != 1 {
			t.Fatalf("Expected exactly one Topology to exist with name '%s', found: %d", topologyName, len(forkedTopology.Response))
		}
		for _, cacheGroupName := range cacheGroupNames {
			foundNode := false
			for _, node := range forkedTopology.Response[0].Nodes {
				if node.Cachegroup == cacheGroupName {
					foundNode = true
					break
				}
			}
			if !foundNode {
				t.Fatalf("unable to find topology node with cachegroup %s", cacheGroupName)
			}

			opts = client.NewRequestOptions()
			opts.QueryParameters.Set("name", cacheGroupName)
			cacheGroups, _, err := TOSession.GetCacheGroups(opts)
			if err != nil {
				t.Fatalf("unable to get cachegroup %s: %s", cacheGroupName, err.Error())
			}
			if len(cacheGroups.Response) != 1 {
				t.Fatalf("incorrect number of cachegroups. expected: 1 actual: %d", len(cacheGroups.Response))
			}
			cacheGroup := cacheGroups.Response[0]

			opts.QueryParameters = url.Values{"cachegroup": []string{strconv.Itoa(*cacheGroup.ID)}}
			srvs, _, err := TOSession.GetServers(opts)
			if err != nil {
				t.Fatalf("unable to get a server from cachegroup %s: %v - alerts: %+v", cacheGroupName, err, srvs.Alerts)
			}
			if len(srvs.Response) < 1 {
				t.Fatalf("Expected at least one server in Cache Group #%d - found none", *cacheGroup.ID)
			}
			for _, s := range srvs.Response {
				if _, ok := cachesByCDNCacheGroup[*s.CDNName]; !ok {
					cachesByCDNCacheGroup[*s.CDNName] = make(map[string][]tc.ServerV4)
				}
				cachesByCDNCacheGroup[*s.CDNName][cacheGroupName] = append(cachesByCDNCacheGroup[*s.CDNName][cacheGroupName], s)
			}
		}
		cdnNames := make([]string, 0, len(cachesByCDNCacheGroup))
		for cdn := range cachesByCDNCacheGroup {
			cdnNames = append(cdnNames, cdn)
		}
		if len(cdnNames) < 2 {
			t.Fatalf("expected servers in at least two CDNs, actual number of CDNs: %d", len(cdnNames))
		}
		cdn1 := cdnNames[0]
		cdn2 := cdnNames[1]

		// update status of MID server to OFFLINE
		resp, _, err := TOSession.UpdateServerStatus(*cachesByCDNCacheGroup[cdn1][midCacheGroup][0].ID, tc.ServerPutStatus{
			Status:        util.JSONNameOrIDStr{Name: util.StrPtr("OFFLINE")},
			OfflineReason: util.StrPtr("testing")}, client.RequestOptions{})
		if err != nil {
			t.Fatalf("cannot update server status: %v - alerts: %+v", err, resp.Alerts)
		}

		opts = client.NewRequestOptions()
		for _, cacheGroupName := range cacheGroupNames {
			cgID := *cachesByCDNCacheGroup[cdn1][cacheGroupName][0].CachegroupID
			opts.QueryParameters.Set("cachegroup", strconv.Itoa(cgID))
			srvs, _, err := TOSession.GetServers(opts)
			if err != nil {
				t.Fatalf("unable to get a server from cachegroup %s: %v - alerts: %+v", cacheGroupName, err, srvs.Alerts)
			}
			if len(srvs.Response) < 1 {
				t.Fatalf("Expected at least one Server in Cache Group #%d, found none", cgID)
			}
			for _, s := range srvs.Response {
				if s.HostName == nil || s.UpdPending == nil || s.ID == nil {
					t.Fatal("Traffic Ops returned a representation of a server with null or undefined Host Name and/or ID and/or Update Pending flag")
				}
				if len(cachesByCDNCacheGroup[*s.CDNName][cacheGroupName]) > 0 {
					cachesByCDNCacheGroup[*s.CDNName][cacheGroupName] = []tc.ServerV4{}
				}
				cachesByCDNCacheGroup[*s.CDNName][cacheGroupName] = append(cachesByCDNCacheGroup[*s.CDNName][cacheGroupName], s)
			}
		}
		for _, cacheGroupName := range cacheGroupNames {
			updResp, _, err := TOSession.GetServerUpdateStatus(*cachesByCDNCacheGroup[cdn1][cacheGroupName][0].HostName, client.RequestOptions{})
			if err != nil {
				t.Fatalf("unable to get update status for a server from Cache Group '%s': %v - alerts: %+v", cacheGroupName, err, updResp.Alerts)
			}
			if len(updResp.Response) < 1 {
				t.Fatalf("Expected at least one server with Host Name '%s' to have an update status", *cachesByCDNCacheGroup[cdn1][cacheGroupName][0].HostName)
			}
			updateStatusByCacheGroup[cacheGroupName] = updResp.Response[0]
		}
		// updating the server status does not queue updates within the same cachegroup in same CDN
		if *cachesByCDNCacheGroup[cdn1][midCacheGroup][0].UpdPending {
			t.Fatalf("expected UpdPending: %t, actual: %t", false, *cachesByCDNCacheGroup[cdn1][midCacheGroup][0].UpdPending)
		}
		// updating the server status does not queue updates within the same cachegroup in different CDN
		if *cachesByCDNCacheGroup[cdn2][midCacheGroup][0].UpdPending {
			t.Fatalf("expected UpdPending: %t, actual: %t", false, *cachesByCDNCacheGroup[cdn2][midCacheGroup][0].UpdPending)
		}
		// edgeCacheGroup is a descendant of midCacheGroup
		if !*cachesByCDNCacheGroup[cdn1][edgeCacheGroup][0].UpdPending {
			t.Fatalf("expected UpdPending: %t, actual: %t", true, *cachesByCDNCacheGroup[cdn1][edgeCacheGroup][0].UpdPending)
		}
		// descendant of midCacheGroup in different CDN should not be queued
		if *cachesByCDNCacheGroup[cdn2][edgeCacheGroup][0].UpdPending {
			t.Fatalf("expected UpdPending: %t, actual: %t", false, *cachesByCDNCacheGroup[cdn2][edgeCacheGroup][0].UpdPending)
		}
		if !updateStatusByCacheGroup[edgeCacheGroup].UpdatePending {
			t.Fatalf("expected UpdPending: %t, actual: %t", true, updateStatusByCacheGroup[edgeCacheGroup].UpdatePending)
		}
		// otherEdgeCacheGroup is not a descendant of midCacheGroup but is still in the same topology
		if *cachesByCDNCacheGroup[cdn1][otherEdgeCacheGroup][0].UpdPending {
			t.Fatalf("expected UpdPending: %t, actual: %t", false, *cachesByCDNCacheGroup[cdn1][otherEdgeCacheGroup][0].UpdPending)
		}
		if updateStatusByCacheGroup[otherEdgeCacheGroup].UpdatePending {
			t.Fatalf("expected UpdPending: %t, actual: %t", false, updateStatusByCacheGroup[otherEdgeCacheGroup].UpdatePending)
		}

		squResp, _, err := TOSession.SetServerQueueUpdate(*cachesByCDNCacheGroup[cdn1][midCacheGroup][0].ID, true, client.RequestOptions{})
		if err != nil {
			t.Fatalf("cannot update server status on %s: %v - alerts: %+v", *cachesByCDNCacheGroup[cdn1][midCacheGroup][0].HostName, err, squResp.Alerts)
		}
		for _, cacheGroupName := range cacheGroupNames {
			updResp, _, err := TOSession.GetServerUpdateStatus(*cachesByCDNCacheGroup[cdn1][cacheGroupName][0].HostName, client.RequestOptions{})
			if err != nil {
				t.Fatalf("unable to get an update status for a server from Cache Group '%s': %v - alerts: %+v", cacheGroupName, err, updResp.Alerts)
			}
			if len(updResp.Response) < 1 {
				t.Fatalf("Expected at least one server with Host Name '%s' to have an update status", *cachesByCDNCacheGroup[cdn1][cacheGroupName][0].HostName)
			}
			updateStatusByCacheGroup[cacheGroupName] = updResp.Response[0]
		}

		// edgeCacheGroup is a descendant of midCacheGroup
		if !updateStatusByCacheGroup[edgeCacheGroup].ParentPending {
			t.Fatalf("expected UpdPending: %t, actual: %t", true, updateStatusByCacheGroup[edgeCacheGroup].ParentPending)
		}
		// otherEdgeCacheGroup is not a descendant of midCacheGroup but is still in the same topology
		if updateStatusByCacheGroup[otherEdgeCacheGroup].ParentPending {
			t.Fatalf("expected UpdPending: %t, actual: %t", false, updateStatusByCacheGroup[otherEdgeCacheGroup].ParentPending)
		}

		edgeHostName := *cachesByCDNCacheGroup[cdn1][edgeCacheGroup][0].HostName
		*cachesByCDNCacheGroup[cdn1][edgeCacheGroup][0].HostName = *cachesByCDNCacheGroup[cdn1][midCacheGroup][0].HostName
		_, _, err = TOSession.UpdateServer(*cachesByCDNCacheGroup[cdn1][edgeCacheGroup][0].ID, cachesByCDNCacheGroup[cdn1][edgeCacheGroup][0], client.RequestOptions{})
		if err != nil {
			t.Fatalf("unable to update %s's hostname to %s: %s", edgeHostName, *cachesByCDNCacheGroup[cdn1][midCacheGroup][0].HostName, err)
		}

		updResp, _, err := TOSession.GetServerUpdateStatus(*cachesByCDNCacheGroup[cdn1][midCacheGroup][0].HostName, client.RequestOptions{})
		if err != nil {
			t.Fatalf("expected no error getting server updates for a non-unique hostname %s, got: %v - alerts: %+v", *cachesByCDNCacheGroup[cdn1][midCacheGroup][0].HostName, err, updResp.Alerts)
		}

		*cachesByCDNCacheGroup[cdn1][edgeCacheGroup][0].HostName = edgeHostName
		_, _, err = TOSession.UpdateServer(*cachesByCDNCacheGroup[cdn1][edgeCacheGroup][0].ID, cachesByCDNCacheGroup[cdn1][edgeCacheGroup][0], client.RequestOptions{})
		if err != nil {
			t.Fatalf("unable to revert %s's hostname back to %s: %s", edgeHostName, edgeHostName, err)
		}
	})
}
