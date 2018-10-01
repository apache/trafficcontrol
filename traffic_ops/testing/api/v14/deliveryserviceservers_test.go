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
	"testing"

	"github.com/apache/trafficcontrol/lib/go-log"
)

func TestDeliveryServiceServers(t *testing.T) {
	CreateTestCDNs(t)
	CreateTestTypes(t)
	CreateTestProfiles(t)
	CreateTestStatuses(t)
	CreateTestDivisions(t)
	CreateTestRegions(t)
	CreateTestPhysLocations(t)
	CreateTestCacheGroups(t)
	CreateTestServers(t)
	CreateTestDeliveryServices(t)

	DeleteTestDeliveryServiceServers(t)

	DeleteTestDeliveryServices(t)
	DeleteTestServers(t)
	DeleteTestCacheGroups(t)
	DeleteTestPhysLocations(t)
	DeleteTestRegions(t)
	DeleteTestDivisions(t)
	DeleteTestStatuses(t)
	DeleteTestProfiles(t)
	DeleteTestTypes(t)
	DeleteTestCDNs(t)
}

func DeleteTestDeliveryServiceServers(t *testing.T) {
	log.Debugln("DeleteTestDeliveryServiceServers")

	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Fatalf("cannot GET DeliveryServices: %v\n", err)
	}
	if len(dses) < 1 {
		t.Fatalf("GET DeliveryServices returned no dses, must have at least 1 to test ds-servers")
	}
	ds := dses[0]

	servers, _, err := TOSession.GetServers()
	if err != nil {
		t.Fatalf("cannot GET Servers: %v\n", err)
	}
	if len(servers) < 1 {
		t.Fatalf("GET Servers returned no dses, must have at least 1 to test ds-servers")
	}
	server := servers[0]

	_, err = TOSession.CreateDeliveryServiceServers(ds.ID, []int{server.ID}, true)
	if err != nil {
		t.Fatalf("POST delivery service servers: %v\n", err)
	}

	dsServers, _, err := TOSession.GetDeliveryServiceServers()
	if err != nil {
		t.Fatalf("GET delivery service servers: %v\n", err)
	}

	found := false
	for _, dss := range dsServers.Response {
		if *dss.DeliveryService == ds.ID && *dss.Server == server.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("POST delivery service servers returned success, but ds-server not in GET")
	}

	if _, _, err := TOSession.DeleteDeliveryServiceServer(ds.ID, server.ID); err != nil {
		t.Fatalf("DELETE delivery service server: %v\n", err)
	}

	dsServers, _, err = TOSession.GetDeliveryServiceServers()
	if err != nil {
		t.Fatalf("GET delivery service servers: %v\n", err)
	}

	found = false
	for _, dss := range dsServers.Response {
		if *dss.DeliveryService == ds.ID && *dss.Server == server.ID {
			found = true
			break
		}
	}
	if found {
		t.Fatalf("DELETE delivery service servers returned success, but still in GET")
	}
}
