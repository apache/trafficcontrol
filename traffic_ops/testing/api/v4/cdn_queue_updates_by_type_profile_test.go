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
	"strconv"
	"strings"
	"testing"

	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestCDNQueueUpdateByProfileAndType(t *testing.T) {
	WithObjs(t, []TCObj{Types, CDNs, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers}, func() {
		QueueUpdatesByType(t)
		QueueUpdatesByProfile(t)
	})
}

func QueueUpdatesByType(t *testing.T) {
	allServersResp, _, err := TOSession.GetServers(client.NewRequestOptions())
	if err != nil {
		t.Fatalf("couldn't get all servers: %v", err)
	}

	// Clear updates on all servers to begin with
	for _, s := range allServersResp.Response {
		if s.ID != nil {
			_, _, err = TOSession.SetServerQueueUpdate(*s.ID, false, client.NewRequestOptions())
			if err != nil {
				t.Errorf("couldn't clear updates on server with ID: %d, err: %v", *s.ID, err.Error())
			}
		}
	}

	allServersResp, _, err = TOSession.GetServers(client.NewRequestOptions())
	if err != nil {
		t.Fatalf("couldn't get all servers: %v", err)
	}
	queryOpts := client.NewRequestOptions()
	if len(testData.Servers) < 1 {
		t.Fatalf("no servers to run the tests on...quitting.")
	}
	server := testData.Servers[0]
	opts := client.NewRequestOptions()
	if server.CDNName == nil {
		t.Fatalf("server doesn't have a CDN name...quitting")
	}
	opts.QueryParameters.Set("name", *server.CDNName)

	// Get the first server's CDN ID
	cdns, _, err := TOSession.GetCDNs(opts)
	if err != nil {
		t.Fatalf("error while getting CDNs: %v", err)
	}
	if len(cdns.Response) < 1 {
		t.Fatalf("expected 1 CDN in response, got %d", len(cdns.Response))
	}
	opts.QueryParameters.Del("name")
	queryOpts.QueryParameters.Set("type", server.Type)
	// Queue updates by type (and CDN)
	_, _, err = TOSession.QueueUpdatesForCDN(cdns.Response[0].ID, true, queryOpts)
	if err != nil {
		t.Errorf("couldn't queue updates by type (and CDN): %v", err)
	}

	// Get all the servers for the same CDN and type as that of the first server
	opts.QueryParameters.Set("cdn", strconv.Itoa(cdns.Response[0].ID))
	opts.QueryParameters.Set("type", server.Type)
	serverIDMap := make(map[int]bool, 0)
	resp, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("couldn't get servers by cdn and type: %v", err)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("expected atleast one server in response, got %d", len(resp.Response))
	}
	for _, s := range resp.Response {
		if s.UpdPending == nil || !*s.UpdPending {
			t.Errorf("expected updates to be queued on all the servers filtered by type and CDN, but %s didn't queue updates", *s.HostName)
		}
		if s.ID != nil {
			serverIDMap[*s.ID] = true
		}
	}

	// Make sure that the servers that are not filtered by the above criteria do not have updates queued
	allServersResp, _, err = TOSession.GetServers(client.NewRequestOptions())
	if err != nil {
		t.Fatalf("couldn't get all servers: %v", err)
	}
	for _, s := range allServersResp.Response {
		if s.ID != nil {
			if _, ok := serverIDMap[*s.ID]; !ok {
				if s.UpdPending != nil && *s.UpdPending {
					t.Errorf("did not expect server with ID: %d to have queued updates", *s.ID)
				}
			}

		}
	}
	_, _, err = TOSession.QueueUpdatesForCDN(cdns.Response[0].ID, false, queryOpts)
	if err != nil {
		t.Errorf("couldn't queue updates by type (and CDN): %v", err)
	}
}

func QueueUpdatesByProfile(t *testing.T) {
	allServersResp, _, err := TOSession.GetServers(client.NewRequestOptions())
	if err != nil {
		t.Fatalf("couldn't get all servers: %v", err)
	}

	// Clear updates on all servers to begin with
	for _, s := range allServersResp.Response {
		if s.ID != nil {
			_, _, err = TOSession.SetServerQueueUpdate(*s.ID, false, client.NewRequestOptions())
			if err != nil {
				t.Errorf("couldn't clear updates on server with ID: %d, err: %v", *s.ID, err.Error())
			}
		}
	}

	queryOpts := client.NewRequestOptions()
	if len(testData.Servers) < 1 {
		t.Fatalf("no servers to run the tests on...quitting.")
	}
	server := testData.Servers[0]
	opts := client.NewRequestOptions()
	if server.CDNName == nil || len(server.ProfileNames) == 0 {
		t.Fatalf("server doesn't have a CDN name or a profile name...quitting")
	}

	//Get the first server's CDN ID
	opts.QueryParameters.Set("name", strings.TrimSpace(*server.CDNName))

	cdns, _, err := TOSession.GetCDNs(opts)
	if err != nil {
		t.Fatalf("error while getting CDNs: %v", err)
	}
	if len(cdns.Response) < 1 {
		t.Fatalf("expected 1 CDN in response, got %d", len(cdns.Response))
	}
	opts.QueryParameters.Del("name")

	// Get the first server's Profile ID
	opts.QueryParameters.Set("name", server.ProfileNames[0])
	profiles, _, err := TOSession.GetProfiles(opts)
	if err != nil {
		t.Fatalf("error while getting profiles: %v", err)
	}
	if len(profiles.Response) < 1 {
		t.Fatalf("expected 1 profile in response, got %d", len(profiles.Response))
	}
	opts.QueryParameters.Del("name")
	queryOpts.QueryParameters.Set("profile", profiles.Response[0].Name)
	// Queue updates by profile (and CDN)
	_, _, err = TOSession.QueueUpdatesForCDN(cdns.Response[0].ID, true, queryOpts)
	if err != nil {
		t.Errorf("couldn't queue updates by profile (and CDN): %v", err)
	}

	// Get all the servers for the same CDN and profile as that of the first server
	opts.QueryParameters.Set("cdn", strconv.Itoa(cdns.Response[0].ID))
	opts.QueryParameters.Set("profileId", strconv.Itoa(profiles.Response[0].ID))
	serverIDMap := make(map[int]bool, 0)
	resp, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("couldn't get servers by cdn and profile: %v", err)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("expected atleast one server in response, got %d", len(resp.Response))
	}
	for _, s := range resp.Response {
		if s.UpdPending == nil || !*s.UpdPending {
			t.Errorf("expected updates to be queued on all the servers filtered by profile and CDN, but %s didn't queue updates", *s.HostName)
		}
		if s.ID != nil {
			serverIDMap[*s.ID] = true
		}
	}

	// Make sure that the servers that are not filtered by the above criteria do not have updates queued
	allServersResp, _, err = TOSession.GetServers(client.NewRequestOptions())
	if err != nil {
		t.Fatalf("couldn't get all servers: %v", err)
	}
	for _, s := range allServersResp.Response {
		if s.ID != nil {
			if _, ok := serverIDMap[*s.ID]; !ok {
				if s.UpdPending != nil && *s.UpdPending {
					t.Errorf("did not expect server with ID: %d to have queued updates", *s.ID)
				}
			}

		}
	}
	_, _, err = TOSession.QueueUpdatesForCDN(cdns.Response[0].ID, false, queryOpts)
	if err != nil {
		t.Errorf("couldn't queue updates by type (and CDN): %v", err)
	}
}
