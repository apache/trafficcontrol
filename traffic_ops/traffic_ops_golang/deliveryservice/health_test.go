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
	"github.com/apache/trafficcontrol/lib/go-tc"

	"encoding/json"
	"testing"
)

const crConfigJson = `
{
  "config": {},
  "contentServers": {
    "cache1": {
      "cacheGroup": "cg1",
      "capabilities": [
        "cap1",
        "cap2"
      ],
      "status": "REPORTED",
      "type": "EDGE",
      "deliveryServices": {
        "ds1": [
          "edge.ds1.test.net"
        ]
       }
    },
    "cache2": {
      "cacheGroup": "cg2",
      "capabilities": [
        "cap2",
        "cap3"
      ],
      "status": "REPORTED",
      "type": "EDGE"
    },
    "cache3": {
      "cacheGroup": "cg2",
      "capabilities": [
        "cap3",
        "cap4"
      ],
      "status": "REPORTED",
      "type": "EDGE"
    }
  },
  "deliveryServices": {
    "ds1": {},
    "ds2-topology": {
      "topology": "test_topology",
      "requiredCapabilities": ["cap2"]
    }
  },
  "contentRouters": {
    "tr1": {}
  },
  "edgeLocations": {
    "edge1": {}
  },
  "trafficRouterLocations": {
    "tr-loc": {}
  },
  "monitors": {
    "tm-host": {}
  },
  "stats": {},
  "topologies": {
    "test_topology": {
      "nodes": [
        "cg2"
      ]
    }
  }
}
`
const crStatesJson = `
{
  "caches": {
    "cache1": {
      "isAvailable": true,
      "ipv4Available": true,
      "ipv6Available": true,
      "status": "REPORTED - available",
      "lastPoll": "2022-05-11T19:50:55.036253631Z"
    },
    "cache2": {
      "isAvailable": true,
      "ipv4Available": true,
      "ipv6Available": true,
      "status": "REPORTED - available",
      "lastPoll": "2022-05-11T19:51:06.965095596Z"
    },
    "cache3": {
      "isAvailable": true,
      "ipv4Available": true,
      "ipv6Available": true,
      "status": "REPORTED - available",
      "lastPoll": "2022-05-11T19:51:06.965095596Z"
    }
  },
  "deliveryServices": {
    "ds1": {
      "disabledLocations": [],
      "isAvailable": true
    },
    "ds2-topology": {
      "disabledLocations": [],
      "isAvailable": true
    }
  }
}
`

func TestAddHealth(t *testing.T) {
	crStates := tc.CRStates{}
	crConfig := tc.CRConfig{}
	err := json.Unmarshal([]byte(crStatesJson), &crStates)
	if err != nil {
		t.Fatalf("error unmarshalling crStates: %v", err)
	}
	err = json.Unmarshal([]byte(crConfigJson), &crConfig)
	if err != nil {
		t.Fatalf("error unmarshalling crConfig: %v", err)
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
	_, available, unAvailable := addHealth("ds1", data, 0, 0, crStates, crConfig)
	if available != 1 || unAvailable != 0 {
		t.Errorf("expected ds1 to have 1 online and 0 offline caches, but got %d online and %d offline instead", available, unAvailable)
	}
	// Even though there are 2 REPORTED EDGE caches in cg2, the result should just include 1, because one of them should get filtered out because it's missing a required capability (cap2)
	_, available, unAvailable = addHealth("ds2-topology", data, 0, 0, crStates, crConfig)
	if available != 1 || unAvailable != 0 {
		t.Errorf("expected ds2-topology to have 1 online and 0 offline caches, but got %d online and %d offline instead", available, unAvailable)
	}
}
