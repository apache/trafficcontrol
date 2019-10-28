package v14

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
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestServerStatusCounts(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers}, func() {
		GetTestServerStatusCounts(t)
	})
}

func GetTestServerStatusCounts(t *testing.T) {
	servers, _, err := TOSession.GetServers()
	if err != nil {
		t.Errorf("cannot GET Servers: %v", err)
	}
	manualCounts := make(map[string]int)
	manualEdgeCounts := make(map[string]int)
	for _, server := range servers {
		manualCounts[server.Status]++
		if server.Type == "EDGE" {
			manualEdgeCounts[server.Status]++
		}
	}
	apiCounts, _, err := TOSession.GetServerStatusCounts(nil)
	if !reflect.DeepEqual(manualCounts, apiCounts) {
		t.Errorf("expected server status counts: %v, actual: %v", manualCounts, apiCounts)
	}
	apiEdgeCounts, _, err := TOSession.GetServerStatusCounts(util.StrPtr("EDGE"))
	if !reflect.DeepEqual(manualEdgeCounts, apiEdgeCounts) {
		t.Errorf("expected EDGE server status counts: %v, actual: %v", manualEdgeCounts, apiEdgeCounts)
	}
}
