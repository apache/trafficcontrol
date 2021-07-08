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
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestCRConfig(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices}, func() {
		UpdateTestCRConfigSnapshot(t)
		MonitoringConfig(t)
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

	tenantOpts := client.NewRequestOptions()
	tenantOpts.QueryParameters.Set("name", "root")
	resp, _, err := TOSession.GetTenants(tenantOpts)
	if err != nil {
		t.Fatalf("couldn't get the root tenant ID: %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Tenant to have the name 'root', found: %d", len(resp.Response))
	}

	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	user := tc.UserV40{
		User: tc.User{
			Username:             util.StrPtr("test_user"),
			RegistrationSent:     tc.TimeNoModFromTime(time.Now()),
			LocalPassword:        util.StrPtr("test_pa$$word"),
			ConfirmLocalPassword: util.StrPtr("test_pa$$word"),
			RoleName:             util.StrPtr("read-only user"),
		},
	}
	user.Email = util.StrPtr("email@domain.com")
	user.TenantID = util.IntPtr(resp.Response[0].ID)
	user.FullName = util.StrPtr("firstName LastName")

	u, _, err := TOSession.CreateUser(user, client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not create read-only user: %v - alerts: %+v", err, u.Alerts)
	}
	client, _, err := toclient.LoginWithAgent(TOSession.URL, "test_user", "test_pa$$word", true, "to-api-v4-client-tests/tenant4user", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to log in with test_user: %v", err.Error())
	}
	opts := toclient.NewRequestOptions()
	opts.QueryParameters.Set("cdn", testData.CDNs[0].Name)
	_, reqInf, err := client.SnapshotCRConfig(opts)
	if err == nil {
		t.Errorf("expected to get an error about a read-only client trying to snap a CDN, but got none")
	}
	if reqInf.StatusCode != http.StatusForbidden {
		t.Errorf("expected a 403 forbidden status code, but got %d", reqInf.StatusCode)
	}
	if u.Response.Username != nil {
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
	paramAlerts, _, err := TOSession.CreateParameter(tc.Parameter{
		ConfigFile: "global",
		Name:       tmURLParamName,
		Value:      "https://crconfig.tm.url.test.invalid",
	}, client.RequestOptions{})
	if err != nil {
		t.Fatalf("GetCRConfig CreateParameter error expected: nil, actual: %v - alerts: %+v", err, paramAlerts.Alerts)
	}

	// create an ANY_MAP DS assignment to verify that it doesn't show up in the CRConfig
	resp, _, err := TOSession.GetServers(client.RequestOptions{})
	if err != nil {
		t.Fatalf("GetServers err expected nil, actual: %v - alerts: %+v", err, resp.Alerts)
	}
	servers := resp.Response
	serverID := 0
	for _, server := range servers {
		if server.CDNName == nil || server.ID == nil {
			t.Error("Traffic Ops returned a representation for a servver with null or undefined ID and/or CDN name")
			continue
		}
		if server.Type == "EDGE" && *server.CDNName == "cdn1" {
			serverID = *server.ID
			break
		}
	}
	if serverID == 0 {
		t.Errorf("GetServers expected EDGE server in cdn1, actual: %+v", servers)
	}
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("xmlId", "anymap-ds")
	res, _, err := TOSession.GetDeliveryServices(opts)
	if err != nil {
		t.Errorf("Unexpected error getting Delivery Services filtered by XMLID 'anymap-ds': %v - alerts: %+v", err, res.Alerts)
	}
	if len(res.Response) != 1 {
		t.Fatalf("Expected exactly 1 Delivery Service to exist with XMLID 'anymap-ds', actual %d", len(res.Response))
	}
	if res.Response[0].ID == nil {
		t.Fatal("Traffic Ops returned a representation of Delivery Service 'anymap-ds' that had a null or undefined ID")
	}
	anymapDSID := *res.Response[0].ID
	alerts, _, err := TOSession.CreateDeliveryServiceServers(anymapDSID, []int{serverID}, true, client.RequestOptions{})
	if err != nil {
		t.Fatalf("Unexpected error assigning server #%d to Delivery Service #%d: %v - alerts: %+v", serverID, anymapDSID, err, alerts.Alerts)
	}

	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("cdn", cdn)
	snapshotResp, _, err := TOSession.SnapshotCRConfig(opts)
	if err != nil {
		t.Errorf("Unexpected error taking Snapshot of CDN '%s': %v - alerts: %+v", cdn, err, snapshotResp.Alerts)
	}
	crcResp, _, err := TOSession.GetCRConfig(cdn, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error retrieving Snapshot of CDN '%s': %v - alerts: %+v", cdn, err, crcResp.Alerts)
	}
	crc := crcResp.Response

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

	opts.QueryParameters.Del("cdn")
	opts.QueryParameters.Set("name", tmURLParamName)
	paramResp, _, err := TOSession.GetParameters(opts)
	if err != nil {
		t.Fatalf("cannot get Parameter by name '%s': %v - alerts: %+v", tmURLParamName, err, paramResp.Alerts)
	}
	if len(paramResp.Response) == 0 {
		t.Fatal("CRConfig create tm.url parameter was successful, but GET returned no parameters")
	}
	tmURLParam := paramResp.Response[0]

	delResp, _, err := TOSession.DeleteParameter(tmURLParam.ID, client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot DELETE Parameter by name: %v - %v", err, delResp)
	}

	crcResp, _, err = TOSession.GetCRConfigNew(cdn, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error getting new Snapshot for CDN '%s': %v - alerts: %+v", cdn, err, crcResp.Alerts)
	}
	crcNew := crcResp.Response

	if len(crcNew.DeliveryServices) != len(crc.DeliveryServices) {
		t.Errorf("/new endpoint returned a different snapshot. DeliveryServices length expected %v, was %v", len(crc.DeliveryServices), len(crcNew.DeliveryServices))
	}

	if *crcNew.Stats.TMHost != "" {
		t.Errorf("update to snapshot not captured in /new endpoint")
	}
}

func MonitoringConfig(t *testing.T) {
	if len(testData.CDNs) < 1 {
		t.Fatalf("no cdn test data")
	}
	const cdnName = "cdn1"
	const profileName = "EDGE1"
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", cdnName)
	cdns, _, err := TOSession.GetCDNs(opts)
	if err != nil {
		t.Fatalf("getting CDNs with name '%s': %v - alerts: %+v", cdnName, err, cdns.Alerts)
	}
	if len(cdns.Response) != 1 {
		t.Fatalf("expected exactly 1 CDN named '%s' but found %d CDNs", cdnName, len(cdns.Response))
	}
	opts.QueryParameters.Set("name", profileName)
	profiles, _, err := TOSession.GetProfiles(opts)
	if err != nil {
		t.Fatalf("getting Profiles with name '%s': %v - alerts: %+v", profileName, err, profiles.Alerts)
	}
	if len(profiles.Response) != 1 {
		t.Fatalf("expected exactly 1 Profiles named %s but found %d Profiles", profileName, len(profiles.Response))
	}
	parameters, _, err := TOSession.GetParametersByProfileName(profileName, client.RequestOptions{})
	if err != nil {
		t.Fatalf("getting Parameters by Profile name '%s': %v - alerts: %+v", profileName, err, parameters.Alerts)
	}
	parameterMap := map[string]tc.HealthThreshold{}
	parameterFound := map[string]bool{}
	const thresholdPrefixLength = len(tc.ThresholdPrefix)
	for _, parameter := range parameters.Response {
		if !strings.HasPrefix(parameter.Name, tc.ThresholdPrefix) {
			continue
		}
		parameterName := parameter.Name[thresholdPrefixLength:]
		parameterMap[parameterName], err = tc.StrToThreshold(parameter.Value)
		if err != nil {
			t.Fatalf("converting string '%s' to HealthThreshold: %s", parameter.Value, err.Error())
		}
		parameterFound[parameterName] = false
	}
	const expectedThresholdParameters = 3
	if len(parameterMap) != expectedThresholdParameters {
		t.Fatalf("expected Profile '%s' to contain %d Parameters with names starting with '%s' but %d such Parameters were found", profileName, expectedThresholdParameters, tc.ThresholdPrefix, len(parameterMap))
	}
	tmConfig, _, err := TOSession.GetTrafficMonitorConfig(cdnName, client.RequestOptions{})
	if err != nil {
		t.Fatalf("getting Traffic Monitor Config: %v - alerts: %+v", err, tmConfig.Alerts)
	}
	profileFound := false
	var profile tc.TMProfile
	for _, profile = range tmConfig.Response.Profiles {
		if profile.Name == profileName {
			profileFound = true
			break
		}
	}
	if !profileFound {
		t.Fatalf("Traffic Monitor Config contained no Profile named '%s", profileName)
	}
	for parameterName, value := range profile.Parameters.Thresholds {
		if _, ok := parameterFound[parameterName]; !ok {
			t.Fatalf("unexpected Threshold Parameter name '%s' found in Profile '%s' in Traffic Monitor Config", parameterName, profileName)
		}
		parameterFound[parameterName] = true
		if parameterMap[parameterName].String() != value.String() {
			t.Fatalf("expected '%s' but received '%s' for Threshold Parameter '%s' in Profile '%s' in Traffic Monitor Config", parameterMap[parameterName].String(), value.String(), parameterName, profileName)
		}
	}
	missingParameters := []string{}
	for parameterName, found := range parameterFound {
		if !found {
			missingParameters = append(missingParameters, parameterName)
		}
	}
	if len(missingParameters) != 0 {
		t.Fatalf("Threshold parameters defined for Profile '%s' but missing for Profile '%s' in Traffic Monitor Config: %s", profileName, profileName, strings.Join(missingParameters, ", "))
	}
}

func SnapshotTestCDNbyName(t *testing.T) {
	if len(testData.CDNs) < 1 {
		t.Fatal("Need at least one CDN to test taking CDN Snapshot using CDN name")
	}
	firstCDN := testData.CDNs[0].Name
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("cdn", firstCDN)
	resp, _, err := TOSession.SnapshotCRConfig(opts)
	if err != nil {
		t.Errorf("failed to snapshot CDN '%s' by name: %v - alerts: %+v", firstCDN, err, resp.Alerts)
	}
}

// Note that this test will break if anyone adds a CDN to the fixture data with
// the name "cdn-invalid".
func SnapshotTestCDNbyInvalidName(t *testing.T) {
	invalidCDNName := "cdn-invalid"
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("cdn", invalidCDNName)
	_, _, err := TOSession.SnapshotCRConfig(opts)
	if err == nil {
		t.Errorf("snapshot occurred without error on (presumed) invalid CDN '%s'", invalidCDNName)
	}
}

func SnapshotTestCDNbyID(t *testing.T) {
	if len(testData.CDNs) < 1 {
		t.Fatal("Need at least one CDN to test Snapshotting CDNs")
	}
	firstCDNName := testData.CDNs[0].Name
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", firstCDNName)
	// Retrieve the CDN by name so we can get the id for the snapshot
	resp, _, err := TOSession.GetCDNs(opts)
	if err != nil {
		t.Errorf("cannot get CDN '%s': %v - alerts: %+v", firstCDNName, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one CDN to exist with name '%s', found: %d", firstCDNName, len(resp.Response))
	}
	remoteCDNID := resp.Response[0].ID
	opts.QueryParameters.Del("name")
	opts.QueryParameters.Set("cdnID", strconv.Itoa(remoteCDNID))
	alert, _, err := TOSession.SnapshotCRConfig(opts)
	if err != nil {
		t.Errorf("failed to snapshot CDN '%s' (#%d) by id: %v - alerts: %+v", firstCDNName, remoteCDNID, err, alert.Alerts)
	}
}

// Note that this test will break in the event that 1,000,000 CDNs are created
// in the TO instance at any time (they don't need to exist concurrently, just
// that many successful CDN creations have to happen, even if they are
// all immediately deleted except the 999999th one).
func SnapshotTestCDNbyInvalidID(t *testing.T) {
	opts := client.NewRequestOptions()
	invalidCDNID := 999999
	opts.QueryParameters.Set("cdnID", strconv.Itoa(invalidCDNID))
	alert, _, err := TOSession.SnapshotCRConfig(opts)
	if err == nil {
		t.Errorf("snapshot occurred on (presumed) invalid CDN #%d: %v - alerts: %+v", invalidCDNID, err, alert.Alerts)
	}
}
