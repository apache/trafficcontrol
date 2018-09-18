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
	"encoding/json"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestCRConfig(t *testing.T) {
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

	UpdateTestCRConfigSnapshot(t)

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

func UpdateTestCRConfigSnapshot(t *testing.T) {
	log.Debugln("UpdateTestCRConfigSnapshot")

	if len(testData.CDNs) < 1 {
		t.Fatalf("no cdn test data")
	}
	cdn := testData.CDNs[0].Name
	_, err := TOSession.SnapshotCRConfig(cdn)
	if err != nil {
		t.Fatalf("SnapshotCRConfig err expected nil, actual %+v", err)
	}
	crcBts, _, err := TOSession.GetCRConfig(cdn)
	if err != nil {
		t.Fatalf("GetCRConfig err expected nil, actual %+v", err)
	}
	crc := tc.CRConfig{}
	if err := json.Unmarshal(crcBts, &crc); err != nil {
		t.Fatalf("GetCRConfig bytes expected: valid tc.CRConfig, actual JSON unmarshal err: %+v", err)
	}

	if len(crc.DeliveryServices) == 0 {
		t.Fatalf("GetCRConfig len(crc.DeliveryServices) expected: >0, actual: 0")
	}

	log.Debugln("UpdateTestCRConfigSnapshot() PASSED: ")
}
