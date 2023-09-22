package deliveryservice

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
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

func TestAddHealth(t *testing.T) {
	crStates := tc.CRStates{
		Caches: map[tc.CacheName]tc.IsAvailable{
			"cache1": {
				IsAvailable:   true,
				Ipv4Available: true,
				Ipv6Available: true,
				Status:        "REPORTED - available",
				LastPoll:      time.Now(),
			},
			"cache2": {
				IsAvailable:   true,
				Ipv4Available: true,
				Ipv6Available: true,
				Status:        "REPORTED - available",
				LastPoll:      time.Now(),
			},
			"cache3": {
				IsAvailable:   true,
				Ipv4Available: true,
				Ipv6Available: true,
				Status:        "REPORTED - available",
				LastPoll:      time.Now(),
			},
		},
		DeliveryService: map[tc.DeliveryServiceName]tc.CRStatesDeliveryService{
			"ds1": {
				DisabledLocations: []tc.CacheGroupName{},
				IsAvailable:       true,
			},
			"ds2-topology": {
				DisabledLocations: []tc.CacheGroupName{},
				IsAvailable:       true,
			},
		},
	}

	status := tc.CRConfigServerStatus("REPORTED")
	crConfig := tc.CRConfig{
		Config: nil,
		ContentServers: map[string]tc.CRConfigTrafficOpsServer{
			"cache1": {
				CacheGroup:   util.StrPtr("cg1"),
				Capabilities: []string{"cap1", "cap2"},
				ServerStatus: &status,
				ServerType:   util.StrPtr("EDGE"),
				DeliveryServices: map[string][]string{
					"ds1": {"edge.ds1.test.net"},
				},
			},
			"cache2": {
				CacheGroup:   util.StrPtr("cg2"),
				Capabilities: []string{"cap2", "cap3"},
				ServerStatus: &status,
				ServerType:   util.StrPtr("EDGE"),
			},
			"cache3": {
				CacheGroup:   util.StrPtr("cg2"),
				Capabilities: []string{"cap3", "cap4"},
				ServerStatus: &status,
				ServerType:   util.StrPtr("EDGE"),
			},
		},
		ContentRouters: nil,
		DeliveryServices: map[string]tc.CRConfigDeliveryService{
			"ds1": {},
			"ds2-topology": {
				Topology:             util.StrPtr("test_topology"),
				RequiredCapabilities: []string{"cap2"},
			},
		},
		EdgeLocations:   nil,
		RouterLocations: nil,
		Monitors:        nil,
		Stats:           tc.CRConfigStats{},
		Topologies: map[string]tc.CRConfigTopology{
			"test_topology": {Nodes: []string{"cg2"}},
		},
	}
	data := make(map[tc.CacheGroupName]tc.HealthDataCacheGroup)
	data[tc.CacheGroupName("cache1")] = tc.HealthDataCacheGroup{
		Offline: 0,
		Online:  0,
		Name:    "cg1",
	}
	data[tc.CacheGroupName("cache2")] = tc.HealthDataCacheGroup{
		Offline: 0,
		Online:  0,
		Name:    "cg2",
	}
	data[tc.CacheGroupName("cache3")] = tc.HealthDataCacheGroup{
		Offline: 0,
		Online:  0,
		Name:    "cg2",
	}
	_, available, unAvailable, err := addHealth("ds1", data, 0, 0, crStates, crConfig)
	if err != nil {
		t.Fatalf("expected no error while adding health of ds1, but got %v", err)
	}
	if available != 1 || unAvailable != 0 {
		t.Errorf("expected ds1 to have 1 online and 0 offline caches, but got %d online and %d offline instead", available, unAvailable)
	}
	// Even though there are 2 REPORTED EDGE caches in cg2, the result should just include 1, because one of them should get filtered out because it's missing a required capability (cap2)
	_, available, unAvailable, err = addHealth("ds2-topology", data, 0, 0, crStates, crConfig)
	if err != nil {
		t.Fatalf("expected no error while adding health of ds2, but got %v", err)
	}
	if available != 1 || unAvailable != 0 {
		t.Errorf("expected ds2-topology to have 1 online and 0 offline caches, but got %d online and %d offline instead", available, unAvailable)
	}
}
