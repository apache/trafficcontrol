package v4

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
	"net/http"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestCRConfig(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices}, func() {
		UpdateTestCRConfigSnapshot(t)
		SnapshotTestCDNbyName(t)
		SnapshotTestCDNbyInvalidName(t)
		SnapshotTestCDNbyID(t)
		SnapshotTestCDNbyInvalidID(t)
		SnapshotWithReadOnlyUser(t)
	})
}

func SnapshotWithReadOnlyUser(t *testing.T) {
	if len(testData.CDNs) == 0 {
		t.Fatalf("expected one or more valid CDNs, but got none")
	}
	resp, _, err := TOSession.TenantByNameWithHdr("root", nil)
	if err != nil {
		t.Fatalf("couldn't get the root tenant ID: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected a valid tenant response, but got nothing")
	}

	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	user := tc.User{
		Username:             util.StrPtr("test_user"),
		RegistrationSent:     tc.TimeNoModFromTime(time.Now()),
		LocalPassword:        util.StrPtr("test_pa$$word"),
		ConfirmLocalPassword: util.StrPtr("test_pa$$word"),
		RoleName:             util.StrPtr("read-only user"),
	}
	user.Email = util.StrPtr("email@domain.com")
	user.TenantID = util.IntPtr(resp.ID)
	user.FullName = util.StrPtr("firstName LastName")

	u, _, err := TOSession.CreateUser(&user)
	if err != nil {
		t.Fatalf("could not create read-only user: %v", err)
	}
	client, _, err := toclient.LoginWithAgent(TOSession.URL, "test_user", "test_pa$$word", true, "to-api-v4-client-tests/tenant4user", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to log in with test_user: %v", err.Error())
	}
	reqInf, err := client.SnapshotCRConfigWithHdr(testData.CDNs[0].Name, nil)
	if err == nil {
		t.Errorf("expected to get an error about a read-only client trying to snap a CDN, but got none")
	}
	if reqInf.StatusCode != http.StatusForbidden {
		t.Errorf("expected a 403 forbidden status code, but got %d", reqInf.StatusCode)
	}
	if u != nil && u.Response.Username != nil {
		ForceDeleteTestUsersByUsernames(t, []string{"test_user"})
	}
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
	servers := resp
	serverID := 0
	for _, server := range servers {
		if server.Type == "EDGE" && *server.CDNName == "cdn1" {
			serverID = *server.ID
			break
		}
	}
	if serverID == 0 {
		t.Errorf("GetServers expected EDGE server in cdn1, actual: %+v", servers)
	}
	res, _, err := TOSession.GetDeliveryServiceByXMLIDNullable("anymap-ds")
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

	if crc.Stats.TMPath != nil {
		t.Errorf("Expected no TMPath in APIv4, but it was: %v", *crc.Stats.TMPath)
	}

	if crc.Stats.TMHost == nil {
		t.Errorf("GetCRConfig crc.Stats.Path expected: '"+tmURLExpected+"', actual: %+v", crc.Stats.TMHost)
	} else if *crc.Stats.TMHost != tmURLExpected {
		t.Errorf("GetCRConfig crc.Stats.Path expected: '"+tmURLExpected+"', actual: %+v", *crc.Stats.TMHost)
	}

	paramResp, _, err := TOSession.GetParameterByName(tmURLParamName)
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

	crcBtsNew, _, err := TOSession.GetCRConfigNew(cdn)
	if err != nil {
		t.Errorf("GetCRConfig err expected nil, actual %+v", err)
	}
	crcNew := tc.CRConfig{}
	if err := json.Unmarshal(crcBtsNew, &crcNew); err != nil {
		t.Errorf("GetCRConfig bytes expected: valid tc.CRConfig, actual JSON unmarshal err: %+v", err)
	}

	if len(crcNew.DeliveryServices) != len(crc.DeliveryServices) {
		t.Errorf("/new endpoint returned a different snapshot. DeliveryServices length expected %v, was %v", len(crc.DeliveryServices), len(crcNew.DeliveryServices))
	}

	if *crcNew.Stats.TMHost != "" {
		t.Errorf("update to snapshot not captured in /new endpoint")
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
	resp, _, err := TOSession.GetCDNByName(firstCDN.Name)
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
