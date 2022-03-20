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
		testVals := func(configUpdate, configApply, revalUpdate, revalApply *time.Time, name string) {
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
			if alerts, _, err := TOSession.SetUpdateServerStatusTimes(*testServer.HostName, configUpdate, configApply, revalUpdate, revalApply, client.RequestOptions{}); err != nil {
				t.Errorf("%v, %v, %v, %v, %s", configUpdate, configApply, revalUpdate, revalApply, name)
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
			if configUpdate != nil {
				if afterServer.ConfigUpdateTime == nil || !afterServer.ConfigUpdateTime.Equal(*configUpdate) {
					t.Errorf("Faild to set server's ConfigUpdateTime. expected: %v actual: %v", *configUpdate, afterServer.ConfigUpdateTime)
				}
			}
			if configApply != nil {
				if afterServer.ConfigApplyTime == nil || !afterServer.ConfigApplyTime.Equal(*configApply) {
					t.Errorf("Faild to set server's ConfigApplyTime. expected: %v actual: %v", *configApply, afterServer.ConfigApplyTime)
				}
			}
			if revalUpdate != nil {
				if afterServer.RevalUpdateTime == nil || !afterServer.RevalUpdateTime.Equal(*revalUpdate) {
					t.Errorf("Faild to set server's RevalUpdateTime. expected: %v actual: %v", *revalUpdate, afterServer.RevalUpdateTime)
				}
			}
			if revalApply != nil {
				if afterServer.RevalApplyTime == nil || !afterServer.RevalApplyTime.Equal(*revalApply) {
					t.Errorf("Faild to set server's RevalApplyTime. expected: %v actual: %v", *revalApply, afterServer.RevalApplyTime)
				}
			}

			// Ensure boolean logic continues to work as expected
			if configUpdate != nil && configApply != nil {
				if ((*configUpdate).Before(*configApply) || (*configUpdate).Equal(*configApply)) &&
					*afterServer.UpdPending {
					t.Error("The configUpdateTime <= configApplyTime. UpdPending should be false")
				} else if (*configUpdate).After(*configApply) && !*afterServer.UpdPending {
					t.Error("The configUpdateTime > configApplyTime. UpdPending should be true")
				}
			}
			if revalUpdate != nil && revalApply != nil {
				if ((*revalUpdate).Before(*revalApply) || (*revalUpdate).Equal(*revalApply)) &&
					*afterServer.RevalPending {
					t.Error("The configUpdateTime <= configApplyTime. RevalPending should be false")
				} else if (*revalUpdate).After(*revalApply) && !*afterServer.RevalPending {
					t.Error("The configUpdateTime > configApplyTime. RevalPending should be true")
				}
			}
		}

		// Postgres stores microsecond precision. There is also some discussion around MacOS losing
		// precision as well. The nanosecond precision is accurate within go one linux however,
		// but round trips to and from the database may result in an inaccurate Equals comparison
		// with the loss of precision. Also, it appears to Round and not Truncate.
		now := time.Now().Round(time.Microsecond)
		later := time.Now().Add(time.Hour * 6)

		// Test setting the values works as expected
		testVals(util.TimePtr(now), nil, nil, nil, "configUpdate")
		testVals(nil, util.TimePtr(now), nil, nil, "configApply")
		testVals(nil, nil, util.TimePtr(now), nil, "revalUpdate")
		testVals(nil, nil, nil, util.TimePtr(now), "revalApply")

		// Test the boolean logic works as expected
		testVals(util.TimePtr(now), util.TimePtr(now), nil, nil, "configUpdate = configApply")
		testVals(util.TimePtr(now), util.TimePtr(later), nil, nil, "configUpdate < configApply")
		testVals(nil, nil, util.TimePtr(now), util.TimePtr(now), "revalUpdate = revalApply")
		testVals(nil, nil, util.TimePtr(now), util.TimePtr(later), "revalUpdate < revalApply")

		// Test sending all nils. Should fail
		if _, _, err := TOSession.SetUpdateServerStatusTimes(*testServer.HostName, nil, nil, nil, nil, client.RequestOptions{}); err == nil {
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
		cachesByCacheGroup := map[string]tc.ServerV40{}
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
			cachesByCacheGroup[cacheGroupName] = srvs.Response[0]
		}

		// update status of MID server to OFFLINE
		resp, _, err := TOSession.UpdateServerStatus(*cachesByCacheGroup[midCacheGroup].ID, tc.ServerPutStatus{
			Status:        util.JSONNameOrIDStr{Name: util.StrPtr("OFFLINE")},
			OfflineReason: util.StrPtr("testing")}, client.RequestOptions{})
		if err != nil {
			t.Fatalf("cannot update server status: %v - alerts: %+v", err, resp.Alerts)
		}

		opts = client.NewRequestOptions()
		for _, cacheGroupName := range cacheGroupNames {
			cgID := *cachesByCacheGroup[cacheGroupName].CachegroupID
			opts.QueryParameters.Set("cachegroup", strconv.Itoa(cgID))
			srvs, _, err := TOSession.GetServers(opts)
			if err != nil {
				t.Fatalf("unable to get a server from cachegroup %s: %v - alerts: %+v", cacheGroupName, err, srvs.Alerts)
			}
			if len(srvs.Response) < 1 {
				t.Fatalf("Expected at least one Server in Cache Group #%d, found none", cgID)
			}
			srv := srvs.Response[0]
			if srv.HostName == nil || srv.UpdPending == nil || srv.ID == nil {
				t.Fatal("Traffic Ops returned a representation of a server with null or undefined Host Name and/or ID and/or Update Pending flag")
			}
			cachesByCacheGroup[cacheGroupName] = srvs.Response[0]
		}
		for _, cacheGroupName := range cacheGroupNames {
			updResp, _, err := TOSession.GetServerUpdateStatus(*cachesByCacheGroup[cacheGroupName].HostName, client.RequestOptions{})
			if err != nil {
				t.Fatalf("unable to get update status for a server from Cache Group '%s': %v - alerts: %+v", cacheGroupName, err, updResp.Alerts)
			}
			if len(updResp.Response) < 1 {
				t.Fatalf("Expected at least one server with Host Name '%s' to have an update status", *cachesByCacheGroup[cacheGroupName].HostName)
			}
			updateStatusByCacheGroup[cacheGroupName] = updResp.Response[0]
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

		squResp, _, err := TOSession.SetServerQueueUpdate(*cachesByCacheGroup[midCacheGroup].ID, true, client.RequestOptions{})
		if err != nil {
			t.Fatalf("cannot update server status on %s: %v - alerts: %+v", *cachesByCacheGroup[midCacheGroup].HostName, err, squResp.Alerts)
		}
		for _, cacheGroupName := range cacheGroupNames {
			updResp, _, err := TOSession.GetServerUpdateStatus(*cachesByCacheGroup[cacheGroupName].HostName, client.RequestOptions{})
			if err != nil {
				t.Fatalf("unable to get an update status for a server from Cache Group '%s': %v - alerts: %+v", cacheGroupName, err, updResp.Alerts)
			}
			if len(updResp.Response) < 1 {
				t.Fatalf("Expected at least one server with Host Name '%s' to have an update status", *cachesByCacheGroup[cacheGroupName].HostName)
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

		edgeHostName := *cachesByCacheGroup[edgeCacheGroup].HostName
		*cachesByCacheGroup[edgeCacheGroup].HostName = *cachesByCacheGroup[midCacheGroup].HostName
		_, _, err = TOSession.UpdateServer(*cachesByCacheGroup[edgeCacheGroup].ID, cachesByCacheGroup[edgeCacheGroup], client.RequestOptions{})
		if err != nil {
			t.Fatalf("unable to update %s's hostname to %s: %s", edgeHostName, *cachesByCacheGroup[midCacheGroup].HostName, err)
		}

		updResp, _, err := TOSession.GetServerUpdateStatus(*cachesByCacheGroup[midCacheGroup].HostName, client.RequestOptions{})
		if err != nil {
			t.Fatalf("expected no error getting server updates for a non-unique hostname %s, got: %v - alerts: %+v", *cachesByCacheGroup[midCacheGroup].HostName, err, updResp.Alerts)
		}

		*cachesByCacheGroup[edgeCacheGroup].HostName = edgeHostName
		_, _, err = TOSession.UpdateServer(*cachesByCacheGroup[edgeCacheGroup].ID, cachesByCacheGroup[edgeCacheGroup], client.RequestOptions{})
		if err != nil {
			t.Fatalf("unable to revert %s's hostname back to %s: %s", edgeHostName, edgeHostName, err)
		}
	})
}
