package v3

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

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestCRConfig(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices}, func() {
		UpdateTestCRConfigSnapshot(t)
		SnapshotTestCDNbyName(t)
		SnapshotTestCDNbyInvalidName(t)
		SnapshotTestCDNbyID(t)
		SnapshotTestCDNbyInvalidID(t)
	})
}

func UpdateTestCRConfigSnapshot(t *testing.T) {
	if len(testData.CDNs) < 1 {
		t.Error("no cdn test data")
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

	// create an ANY_MAP DS assignment to verify that it doesn't show up in the CRConfig
	resp, _, err := TOSession.GetServers(nil, nil)
	if err != nil {
		t.Fatalf("GetServers err expected nil, actual %+v", err)
	}
	servers := resp.Response
	serverID := 0
	for _, server := range servers {
		if server.Type == "EDGE" && server.CDNName != nil && *server.CDNName == "cdn1" && server.ID != nil {
			serverID = *server.ID
			break
		}
	}
	if serverID == 0 {
		t.Errorf("GetServers expected EDGE server in cdn1, actual: %+v", servers)
	}
	res, _, err := TOSession.GetDeliveryServiceByXMLIDNullable("anymap-ds", nil)
	if err != nil {
		t.Errorf("GetDeliveryServiceByXMLIDNullable err expected nil, actual %+v", err)
	}
	if len(res) != 1 {
		t.Error("GetDeliveryServiceByXMLIDNullable expected 1 DS, actual 0")
	}
	if res[0].ID == nil {
		t.Error("GetDeliveryServiceByXMLIDNullable got unknown delivery service id")
	}
	anymapDSID := *res[0].ID
	_, _, err = TOSession.CreateDeliveryServiceServers(anymapDSID, []int{serverID}, true)
	if err != nil {
		t.Errorf("POST delivery service servers: %v", err)
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
		t.Error("GetCRConfig len(crc.DeliveryServices) expected: >0, actual: 0")
	}

	// verify no ANY_MAP delivery services are in the CRConfig
	for ds := range crc.DeliveryServices {
		if ds == "anymap-ds" {
			t.Error("found ANY_MAP delivery service in CRConfig deliveryServices")
		}
	}
	for server := range crc.ContentServers {
		for ds := range crc.ContentServers[server].DeliveryServices {
			if ds == "anymap-ds" {
				t.Error("found ANY_MAP delivery service in contentServers deliveryServices mapping")
			}
		}
	}

	if crc.Stats.TMPath == nil {
		t.Errorf("GetCRConfig crc.Stats.Path expected: '/snapshot', actual: %+v", crc.Stats.TMPath)
	} else if !strings.HasSuffix(*crc.Stats.TMPath, "snapshot") {
		t.Errorf("GetCRConfig crc.Stats.Path expected: '/snapshot', actual: %+v", *crc.Stats.TMPath)
	}

	if crc.Stats.TMHost == nil {
		t.Errorf("GetCRConfig crc.Stats.Path expected: '"+tmURLExpected+"', actual: %+v", crc.Stats.TMHost)
	} else if *crc.Stats.TMHost != tmURLExpected {
		t.Errorf("GetCRConfig crc.Stats.Path expected: '"+tmURLExpected+"', actual: %+v", *crc.Stats.TMHost)
	}

	paramResp, _, err := TOSession.GetParameterByName(tmURLParamName, nil)
	if err != nil {
		t.Fatalf("cannot GET Parameter by name: %v - %v", tmURLParamName, err)
	}
	if len(paramResp) == 0 {
		t.Fatal("CRConfig create tm.url parameter was successful, but GET returned no parameters")
	}
	tmURLParam := paramResp[0]

	delResp, _, err := TOSession.DeleteParameterByID(tmURLParam.ID)
	if err != nil {
		t.Fatalf("cannot DELETE Parameter by name: %v - %v", err, delResp)
	}
}

func SnapshotTestCDNbyName(t *testing.T) {

	firstCDN := testData.CDNs[0]
	_, err := TOSession.SnapshotCRConfig(firstCDN.Name)
	if err != nil {
		t.Errorf("failed to snapshot CDN by name: %v", err)
	}
}

func SnapshotTestCDNbyInvalidName(t *testing.T) {

	invalidCDNName := "cdn-invalid"
	_, err := TOSession.SnapshotCRConfig(invalidCDNName)
	if err == nil {
		t.Errorf("snapshot occurred on invalid cdn name: %v - %v", invalidCDNName, err)
	}
}

func SnapshotTestCDNbyID(t *testing.T) {

	firstCDN := testData.CDNs[0]
	// Retrieve the CDN by name so we can get the id for the snapshot
	resp, _, err := TOSession.GetCDNByName(firstCDN.Name, nil)
	if err != nil {
		t.Errorf("cannot GET CDN by name: '%s', %v", firstCDN.Name, err)
	}
	remoteCDN := resp[0]
	alert, _, err := TOSession.SnapshotCRConfigByID(remoteCDN.ID)
	if err != nil {
		t.Errorf("failed to snapshot CDN by id: %v - %v", err, alert)
	}
}

func SnapshotTestCDNbyInvalidID(t *testing.T) {

	invalidCDNID := 999999
	alert, _, err := TOSession.SnapshotCRConfigByID(invalidCDNID)
	if err == nil {
		t.Errorf("snapshot occurred on invalid cdn id: %v - %v - %v", invalidCDNID, err, alert)
	}
}
