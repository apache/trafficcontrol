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

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestMonitoring(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices}, func() {
		AllCDNsCanSnapshot(t)
	})
}

func AllCDNsCanSnapshot(t *testing.T) {

	serversByHost := make(map[string]tc.ServerNullable)

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
		if err != nil {
			t.Error(err)
			continue
		}

		for hostName, server := range tmConfig.TrafficServer {
			if _, ok := serversByHost[hostName]; !ok {
				t.Errorf("Server %v not found in test data", hostName)
				continue
			}
			if len(server.Interfaces) != len(serversByHost[hostName].Interfaces) {
				t.Errorf("Server %v expected to get %v interfaces, got %v", hostName, len(server.Interfaces), len(serversByHost[hostName].Interfaces))
			}
		}
	}
}
