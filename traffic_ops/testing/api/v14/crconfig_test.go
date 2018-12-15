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
	"encoding/json"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestCRConfig(t *testing.T) {
	CreateTestCDNs(t)
	CreateTestTypes(t)
	CreateTestTenants(t)
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
	DeleteTestTenants(t)
	DeleteTestTypes(t)
	DeleteTestCDNs(t)
}

func UpdateTestCRConfigSnapshot(t *testing.T) {
	log.Debugln("UpdateTestCRConfigSnapshot")

	if len(testData.CDNs) < 1 {
		t.Errorf("no cdn test data")
	}
	cdn := testData.CDNs[0].Name

	tmURLParamName := "tm.url"
	tmURLExpected := "crconfig.tm.url.test.invalid"
	_, _, err := TOSession.CreateParameter(tc.Parameter{
		ConfigFile: "global",
		Name:       tmURLParamName,
		Value:      "https://crconfig.tm.url.test.invalid",
	})
	if err != nil {
		t.Fatalf("GetCRConfig CreateParameter error expected: nil, actual: " + err.Error())
	}
	_, err = TOSession.SnapshotCRConfig(cdn)
	if err != nil {
		t.Errorf("SnapshotCRConfig err expected nil, actual %+v", err)
	}
	crcBts, _, err := TOSession.GetCRConfig(cdn)
	if err != nil {
		t.Errorf("GetCRConfig err expected nil, actual %+v", err)
	}
	crc := tc.CRConfig{}
	if err := json.Unmarshal(crcBts, &crc); err != nil {
		t.Errorf("GetCRConfig bytes expected: valid tc.CRConfig, actual JSON unmarshal err: %+v", err)
	}

	if len(crc.DeliveryServices) == 0 {
		t.Errorf("GetCRConfig len(crc.DeliveryServices) expected: >0, actual: 0")
	}

	if crc.Stats.TMPath == nil {
		t.Errorf("GetCRConfig crc.Stats.Path expected: 'snapshot/"+cdn+"', actual: %+v", crc.Stats.TMPath)
	} else if !strings.HasSuffix(*crc.Stats.TMPath, "snapshot/"+cdn) {
		t.Errorf("GetCRConfig crc.Stats.Path expected: '/snapshot"+cdn+"', actual: %+v", *crc.Stats.TMPath)
	}

	if crc.Stats.TMHost == nil {
		t.Errorf("GetCRConfig crc.Stats.Path expected: '"+tmURLExpected+"', actual: %+v", crc.Stats.TMHost)
	} else if *crc.Stats.TMHost != tmURLExpected {
		t.Errorf("GetCRConfig crc.Stats.Path expected: '"+tmURLExpected+"', actual: %+v", *crc.Stats.TMHost)
	}

	paramResp, _, err := TOSession.GetParameterByName(tmURLParamName)
	if err != nil {
		t.Fatalf("cannot GET Parameter by name: %v - %v\n", tmURLParamName, err)
	}
	if len(paramResp) == 0 {
		t.Fatalf("CRConfig create tm.url parameter was successful, but GET returned no parameters")
	}
	tmURLParam := paramResp[0]

	delResp, _, err := TOSession.DeleteParameterByID(tmURLParam.ID)
	if err != nil {
		t.Fatalf("cannot DELETE Parameter by name: %v - %v\n", err, delResp)
	}

	log.Debugln("UpdateTestCRConfigSnapshot() PASSED: ")
}
