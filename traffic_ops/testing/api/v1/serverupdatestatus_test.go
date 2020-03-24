package v1

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
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestServerUpdateStatus(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers}, func() {
		edge1cdn1 := tc.ServerV1{}
		edge2cdn1 := tc.ServerV1{}
		mid1cdn1 := tc.ServerV1{}
		edge1cdn2 := tc.ServerV1{}

		getServers := func() {
			for _, s := range []struct {
				name   string
				server *tc.ServerV1
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
				resp, _, err := TOSession.GetServerByHostName(s.name)
				if err != nil {
					t.Errorf("cannot GET Server by hostname: %v - %v", s.name, err)
				}
				*s.server = resp[0]
			}
		}
		getServers()

		// assert that servers don't have updates pending
		for _, s := range []tc.ServerV1{
			edge1cdn1,
			edge2cdn1,
			mid1cdn1,
			edge1cdn2,
		} {
			if s.UpdPending {
				t.Error("expected UpdPending: false, actual: true")
			}
		}

		// update status of MID server to OFFLINE
		_, _, err := TOSession.UpdateServerStatus(mid1cdn1.ID, tc.ServerPutStatus{
			Status:        util.JSONNameOrIDStr{Name: util.StrPtr("OFFLINE")},
			OfflineReason: util.StrPtr("testing")})
		if err != nil {
			t.Errorf("cannot update server status: %v", err)
		}

		// assert that updates were queued for the proper EDGE servers
		getServers()
		if !edge1cdn1.UpdPending {
			t.Errorf("expected: child %s to have updates pending, actual: no updates pending", edge1cdn1.HostName)
		}
		if !edge2cdn1.UpdPending {
			t.Errorf("expected: child %s to have updates pending, actual: no updates pending", edge2cdn1.HostName)
		}
		if mid1cdn1.UpdPending {
			t.Errorf("expected: server %s with updated status to have no updates pending, actual: updates pending", mid1cdn1.HostName)
		}
		if edge1cdn2.UpdPending {
			t.Errorf("expected: server %s in different CDN than server with updated status to have no updates pending, actual: updates pending", edge2cdn1.HostName)
		}

		// update status of MID server to OFFLINE via status ID
		status, _, err := TOSession.GetStatusByName("OFFLINE")
		if err != nil {
			t.Fatalf("cannot GET status by name: %v", err)
		}
		_, _, err = TOSession.UpdateServerStatus(
			mid1cdn1.ID,
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
			mid1cdn1.ID,
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
			mid1cdn1.ID,
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
			mid1cdn1.ID,
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
	WithObjs(t, []TCObj{Divisions, Regions, PhysLocations, Statuses, Types, CacheGroups, CDNs, Profiles, Servers}, func() {
		const serverName = "atlanta-edge-01"

		queueUpdateActions := map[bool]string{
			false: "dequeue",
			true:  "queue",
		}

		var s tc.ServerV1
		resp, _, err := TOSession.GetServerByHostName(serverName)
		if err != nil {
			t.Fatalf("failed to GET Server by hostname: %v - %v", serverName, err)
		}
		s = resp[0]

		// assert that servers don't have updates pending
		if got, want := s.UpdPending, false; got != want {
			t.Fatalf("unexpected UpdPending, got: %v, want: %v", got, want)
		}

		for _, setVal := range [...]bool{true, false} {
			t.Run(fmt.Sprint(setVal), func(t *testing.T) {
				// queue update and check response
				quResp, _, err := TOSession.SetServerQueueUpdate(s.ID, setVal)
				if err != nil {
					t.Fatalf("failed to set queue update for server with ID %v to %v: %v", s.ID, setVal, err)
				}
				if got, want := int(quResp.Response.ServerID), s.ID; got != want {
					t.Errorf("wrong serverId in response, got: %v, want: %v", got, want)
				}
				if got, want := quResp.Response.Action, queueUpdateActions[setVal]; got != want {
					t.Errorf("wrong action in response, got: %v, want: %v", got, want)
				}

				// assert that the server has updates queued
				resp, _, err = TOSession.GetServerByID(s.ID)
				if err != nil {
					t.Fatalf("failed to GET Server by ID: %v - %v", s.ID, err)
				}
				s = resp[0]
				if got, want := s.UpdPending, setVal; got != want {
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
			path := fmt.Sprintf("/api/1.4/servers/%d/queue_update", s.ID)
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
