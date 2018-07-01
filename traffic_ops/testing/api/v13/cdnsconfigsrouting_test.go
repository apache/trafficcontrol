package v13

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

func TestCDNsConfigsRouting(t *testing.T) {
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

	GetTestCDNsConfigsRouting(t)

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

func GetTestCDNsConfigsRouting(t *testing.T) {
	log.Debugln("GetTestCDNsConfigsRouting")
	if len(testData.CDNs) < 1 {
		t.Fatalf("cannot test cdns configs routing: not test data CDNs\n")
	}
	cdn := testData.CDNs[0].Name
	cfg, _, err := TOSession.GetRouting(cdn)
	if err != nil {
		t.Fatalf("GET cdns configs routing: %v\n", err)
	}
	if len(cfg.TrafficServers) == 0 {
		t.Fatalf("GET cdns configs routing succeeded, but returned no servers\n")
	}
	if len(cfg.TrafficMonitors) == 0 {
		t.Fatalf("GET cdns configs routing succeeded, but returned no monitors\n")
	}
	if len(cfg.TrafficMonitors) == 0 {
		t.Fatalf("GET cdns configs routing succeeded, but returned no routers\n")
	}
	if len(cfg.TrafficMonitors) == 0 {
		t.Fatalf("GET cdns configs routing succeeded, but returned no delivery services\n")
	}
}
