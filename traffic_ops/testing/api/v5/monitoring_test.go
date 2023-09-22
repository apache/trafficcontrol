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
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestMonitoring(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, ProfileParameters, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, ServerCapabilities, ServerServerCapabilities, DeliveryServices}, func() {
		GetTestMonitoringConfigNoSnapshotOnTheFly(t) // MUST run first
		AllCDNsCanSnapshot(t)
	})
}

// GetTestMonitoringConfigNoSnapshotOnTheFly verifies that Traffic Ops generates a monitoring.json on-the-fly rather than returning "" or "{}" if no snapshot exists.
// This MUST NOT be run after a different function in the same Test creates a Snapshot, or the test will be invalid.
// This prevents a critical bug of upgrading to 4.x bringing a CDN down until a Snapshot is performed.
func GetTestMonitoringConfigNoSnapshotOnTheFly(t *testing.T) {
	var server tc.ServerV5
	for _, sv := range testData.Servers {
		if sv.Type != "EDGE" {
			continue
		}
		server = sv
		break
	}
	if server.CDN == "" {
		t.Fatal("No edge server found in test data, cannot test")
	}

	resp, _, err := TOSession.GetTrafficMonitorConfig(server.CDN, client.RequestOptions{})
	if err != nil {
		t.Errorf("getting monitoring: %v - alerts: %+v", err, resp.Alerts)
	} else if len(resp.Response.TrafficServers) == 0 {
		t.Errorf("Expected Monitoring without a snapshot to generate on-the-fly, actual: empty monitoring object for CDN '%s'", server.CDN)
	}
}

func AllCDNsCanSnapshot(t *testing.T) {

	serversByHost := make(map[string]tc.ServerV5, len(testData.Servers))

	for _, server := range testData.Servers {
		serversByHost[server.HostName] = server
	}

	opts := client.NewRequestOptions()
	for _, cdn := range testData.CDNs {
		opts.QueryParameters.Set("cdn", cdn.Name)
		resp, _, err := TOSession.SnapshotCRConfig(opts)
		if err != nil {
			t.Errorf("Unexpected error making Snapshot for CDN '%s': %v - alerts: %+v", cdn.Name, err, resp.Alerts)
			continue
		}

		tmConfig, _, err := TOSession.GetTrafficMonitorConfig(cdn.Name, client.RequestOptions{})
		if err != nil {
			t.Errorf("Unexpected error fetching Traffic Monitor Config for CDN '%s': %v - alerts: %+v", cdn.Name, err, tmConfig.Alerts)
			continue
		}

		for _, server := range tmConfig.Response.TrafficServers {
			if _, ok := serversByHost[server.HostName]; !ok {
				t.Errorf("Server '%s' not found in test data", server.HostName)
				continue
			}
			if len(server.Interfaces) < 1 {
				t.Errorf("Server '%s' expected to get more than 1 interface(s), got %d", server.HostName, len(server.Interfaces))
			}
		}
	}
}
