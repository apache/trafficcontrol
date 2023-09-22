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

package v3

import (
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func TestMonitoring(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, ProfileParameters, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices}, func() {
		GetTestMonitoringConfigNoSnapshotOnTheFly(t) // MUST run first
		AllCDNsCanSnapshot(t)
	})
}

// GetTestMonitoringConfigNoSnapshotOnTheFly verifies that Traffic Ops generates a monitoring.json on-the-fly rather than returning "" or "{}" if no snapshot exists.
// This MUST NOT be run after a different function in the same Test creates a Snapshot, or the test will be invalid.
// This prevents a critical bug of upgrading to 4.x bringing a CDN down until a Snapshot is performed.
func GetTestMonitoringConfigNoSnapshotOnTheFly(t *testing.T) {
	server := tc.ServerV30{}
	for _, sv := range testData.Servers {
		if sv.Type != "EDGE" {
			continue
		}
		server = sv
		break
	}
	if server.CDNName == nil || *server.CDNName == "" {
		t.Fatal("No edge server found in test data, cannot test")
	}

	tmConfig, _, err := TOSession.GetTrafficMonitorConfigMap(*server.CDNName)
	if err != nil {
		t.Error("getting monitoring: " + err.Error())
	} else if len(tmConfig.TrafficServer) == 0 {
		t.Error("Expected Monitoring without a snapshot to generate on-the-fly, actual: empty monitoring object for cdn '" + *server.CDNName + "'")
	}
}

func AllCDNsCanSnapshot(t *testing.T) {

	serversByHost := make(map[string]tc.ServerV30)

	for _, server := range testData.Servers {
		serversByHost[*server.HostName] = server
	}

	for _, cdn := range testData.CDNs {
		_, err := TOSession.SnapshotCRConfig(cdn.Name)
		if err != nil {
			t.Error(err)
			continue
		}

		tmConfig, _, err := TOSession.GetTrafficMonitorConfigMap(cdn.Name)
		if err != nil && tmConfig == nil {
			t.Error(err)
			continue
		}

		for hostName, server := range tmConfig.TrafficServer {
			if _, ok := serversByHost[hostName]; !ok {
				t.Errorf("Server %v not found in test data", hostName)
				continue
			}
			if len(server.Interfaces) < 1 {
				t.Errorf("Server %v expected to get more than 1 interface(s), got %v", hostName, len(server.Interfaces))
			}
		}
	}
}
