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
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v6/lib/go-rfc"
	"github.com/apache/trafficcontrol/v6/lib/go-tc"
	"github.com/apache/trafficcontrol/v6/lib/go-util"
	toclient "github.com/apache/trafficcontrol/v6/traffic_ops/v3-client"
)

func TestDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServerCapabilities, DeliveryServices}, func() {
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		ti := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, ti)
		header.Set(rfc.IfUnmodifiedSince, ti)
		if includeSystemTests {
			SSLDeliveryServiceCDNUpdateTest(t)
			GetTestDeliveryServicesURLSigKeys(t)
		}
		GetTestDeliveryServicesIMS(t)
		GetAccessibleToTest(t)
		UpdateTestDeliveryServices(t)
		UpdateValidateORGServerCacheGroup(t)
		UpdateTestDeliveryServicesWithHeaders(t, header)
		UpdateNullableTestDeliveryServices(t)
		UpdateDeliveryServiceWithInvalidRemapText(t)
		UpdateDeliveryServiceWithInvalidSliceRangeRequest(t)
		UpdateDeliveryServiceWithInvalidTopology(t)
		GetTestDeliveryServicesIMSAfterChange(t, header)
		UpdateDeliveryServiceTopologyHeaderRewriteFields(t)
		GetTestDeliveryServices(t)
		GetInactiveTestDeliveryServices(t)
		GetTestDeliveryServicesCapacity(t)
		DeliveryServiceMinorVersionsTest(t)
		DeliveryServiceTenancyTest(t)
		PostDeliveryServiceTest(t)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestDeliveryServicesWithHeaders(t, header)
	})
}

func UpdateTestDeliveryServicesWithHeaders(t *testing.T, header http.Header) {
	if len(testData.DeliveryServices) > 0 {
		firstDS := testData.DeliveryServices[0]
		if firstDS.XMLID == nil {
			t.Fatalf("couldn't get the xml ID of test DS")
		}
		dses, _, err := TOSession.GetDeliveryServicesNullableWithHdr(header)
		if err != nil {
			t.Errorf("cannot GET Delivery Services: %v", err)
		}

		remoteDS := tc.DeliveryServiceNullable{}
		found := false
		for _, ds := range dses {
			if ds.XMLID != nil && *ds.XMLID == *firstDS.XMLID {
				found = true
				remoteDS = ds
				break
			}
		}
		if !found {
			t.Fatalf("GET Delivery Services missing: %v", *firstDS.XMLID)
		}

		updatedLongDesc := "something different"
		updatedMaxDNSAnswers := 164598
		updatedMaxOriginConnections := 100
		remoteDS.LongDesc = &updatedLongDesc
		remoteDS.MaxDNSAnswers = &updatedMaxDNSAnswers
		remoteDS.MaxOriginConnections = &updatedMaxOriginConnections
		remoteDS.MatchList = nil // verify that this field is optional in a PUT request, doesn't cause nil dereference panic

		_, err = TOSession.UpdateDeliveryServiceNullableWithHdr(strconv.Itoa(*remoteDS.ID), &remoteDS, header)
		if err == nil {
			t.Errorf("expected precondition failed error, got none")
		}
		if !strings.Contains(err.Error(), "412 Precondition Failed") {
			t.Errorf("expected error to be related to 'precondition failed', but instead is realted to %v", err.Error())
		}
	}
}

func createBlankCDN(cdnName string, t *testing.T) tc.CDN {
	_, _, err := TOSession.CreateCDN(tc.CDN{
		DNSSECEnabled: false,
		DomainName:    cdnName + ".ai",
		Name:          cdnName,
	})
	if err != nil {
		t.Fatal("unable to create cdn: " + err.Error())
	}

	originalKeys, _, err := TOSession.GetCDNSSLKeysWithHdr(cdnName, nil)
	if err != nil {
		t.Fatalf("unable to get keys on cdn %v: %v", cdnName, err)
	}

	cdns, _, err := TOSession.GetCDNByNameWithHdr(cdnName, nil)
	if err != nil {
		t.Fatalf("unable to get cdn %v: %v", cdnName, err)
	}
	if len(cdns) < 1 {
		t.Fatal("expected more than 0 cdns")
	}
	keys, _, err := TOSession.GetCDNSSLKeysWithHdr(cdnName, nil)
	if err != nil {
		t.Fatalf("unable to get keys on cdn %v: %v", cdnName, err)
	}
	if len(keys) != len(originalKeys) {
		t.Fatalf("expected %v ssl keys on cdn %v, got %v", len(originalKeys), cdnName, len(keys))
	}
	return cdns[0]
}

func cleanUp(t *testing.T, ds tc.DeliveryServiceNullableV30, oldCDNID int, newCDNID int, sslKeyVersions []string) {
	_, _, err := TOSession.DeleteDeliveryServiceSSLKeysByID(*ds.XMLID)
	if err != nil {
		t.Error(err)
	}
	params := url.Values{}
	for _, version := range sslKeyVersions {
		params.Set("version", version)
		_, _, err := TOSession.DeleteDeliveryServiceSSLKeysByVersion(*ds.XMLID, params)
		if err != nil {
			t.Error(err)
		}
	}
	_, err = TOSession.DeleteDeliveryService(strconv.Itoa(*ds.ID))
	if err != nil {
		t.Error(err)
	}
	_, _, err = TOSession.DeleteCDNByID(oldCDNID)
	if err != nil {
		t.Error(err)
	}
	_, _, err = TOSession.DeleteCDNByID(newCDNID)
	if err != nil {
		t.Error(err)
	}
}

func SSLDeliveryServiceCDNUpdateTest(t *testing.T) {
	cdnNameOld := "sslkeytransfer"
	oldCdn := createBlankCDN(cdnNameOld, t)
	cdnNameNew := "sslkeytransfer1"
	newCdn := createBlankCDN(cdnNameNew, t)

	types, _, err := TOSession.GetTypeByNameWithHdr("HTTP", nil)
	if err != nil {
		t.Fatal("unable to get types: " + err.Error())
	}
	if len(types) < 1 {
		t.Fatal("expected at least one type")
	}

	customDS := tc.DeliveryServiceNullableV30{}
	customDS.Active = util.BoolPtr(true)
	customDS.CDNID = util.IntPtr(oldCdn.ID)
	customDS.DSCP = util.IntPtr(0)
	customDS.DisplayName = util.StrPtr("displayName")
	customDS.RoutingName = util.StrPtr("routingName")
	customDS.GeoLimit = util.IntPtr(0)
	customDS.GeoProvider = util.IntPtr(0)
	customDS.IPV6RoutingEnabled = util.BoolPtr(false)
	customDS.InitialDispersion = util.IntPtr(1)
	customDS.LogsEnabled = util.BoolPtr(true)
	customDS.MissLat = util.FloatPtr(0)
	customDS.MissLong = util.FloatPtr(0)
	customDS.MultiSiteOrigin = util.BoolPtr(false)
	customDS.OrgServerFQDN = util.StrPtr("https://test.com")
	customDS.Protocol = util.IntPtr(2)
	customDS.QStringIgnore = util.IntPtr(0)
	customDS.RangeRequestHandling = util.IntPtr(0)
	customDS.RegionalGeoBlocking = util.BoolPtr(false)
	customDS.TenantID = util.IntPtr(1)
	customDS.TypeID = util.IntPtr(types[0].ID)
	customDS.XMLID = util.StrPtr("dsID")
	customDS.MaxRequestHeaderBytes = nil

	ds, _, err := TOSession.CreateDeliveryServiceV30(customDS)
	if err != nil {
		t.Fatal(err)
	}
	ds.CDNName = &oldCdn.Name

	defer cleanUp(t, ds, oldCdn.ID, newCdn.ID, []string{"1"})

	_, _, err = TOSession.GenerateSSLKeysForDS(*ds.XMLID, *ds.CDNName, tc.SSLKeyRequestFields{
		BusinessUnit: util.StrPtr("BU"),
		City:         util.StrPtr("CI"),
		Organization: util.StrPtr("OR"),
		HostName:     util.StrPtr("*.test.com"),
		Country:      util.StrPtr("CO"),
		State:        util.StrPtr("ST"),
	})
	if err != nil {
		t.Fatalf("unable to generate sslkeys for cdn %v: %v", oldCdn.Name, err)
	}

	tries := 0
	var oldCDNKeys []tc.CDNSSLKeys
	for tries < 5 {
		time.Sleep(time.Second)
		oldCDNKeys, _, err = TOSession.GetCDNSSLKeysWithHdr(oldCdn.Name, nil)
		if err == nil && len(oldCDNKeys) > 0 {
			break
		}
		tries++
	}
	if err != nil {
		t.Fatalf("unable to get cdn %v keys: %v", oldCdn.Name, err)
	}
	if len(oldCDNKeys) < 1 {
		t.Fatal("expected at least 1 key")
	}

	newCDNKeys, _, err := TOSession.GetCDNSSLKeysWithHdr(newCdn.Name, nil)
	if err != nil {
		t.Fatalf("unable to get cdn %v keys: %v", newCdn.Name, err)
	}

	ds.RoutingName = util.StrPtr("anothername")
	_, _, err = TOSession.UpdateDeliveryServiceV30WithHdr(*ds.ID, ds, nil)
	if err == nil {
		t.Fatal("should not be able to update delivery service (routing name) as it has ssl keys")
	}
	ds.RoutingName = util.StrPtr("routingName")

	ds.CDNID = &newCdn.ID
	ds.CDNName = &newCdn.Name
	_, _, err = TOSession.UpdateDeliveryServiceV30WithHdr(*ds.ID, ds, nil)
	if err == nil {
		t.Fatal("should not be able to update delivery service (cdn) as it has ssl keys")
	}

	// Check new CDN still has an ssl key
	keys, _, err := TOSession.GetCDNSSLKeysWithHdr(newCdn.Name, nil)
	if err != nil {
		t.Fatalf("unable to get cdn %v keys: %v", newCdn.Name, err)
	}
	if len(keys) != len(newCDNKeys) {
		t.Fatalf("expected %v keys, got %v", len(newCDNKeys), len(keys))
	}

	// Check old CDN does not have ssl key
	keys, _, err = TOSession.GetCDNSSLKeysWithHdr(oldCdn.Name, nil)
	if err != nil {
		t.Fatalf("unable to get cdn %v keys: %v", oldCdn.Name, err)
	}
	if len(keys) != len(oldCDNKeys) {
		t.Fatalf("expected %v key, got %v", len(oldCDNKeys), len(keys))
	}
}

func GetTestDeliveryServicesIMSAfterChange(t *testing.T, header http.Header) {
	_, reqInf, err := TOSession.GetDeliveryServicesV30WithHdr(header, nil)
	if err != nil {
		t.Fatalf("could not GET Delivery Services: %v", err)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	_, reqInf, err = TOSession.GetDeliveryServicesV30WithHdr(header, nil)
	if err != nil {
		t.Fatalf("could not GET Delivery Services: %v", err)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func PostDeliveryServiceTest(t *testing.T) {
	if len(testData.DeliveryServices) < 1 {
		t.Fatal("Need at least one testing Delivery Service to test creating Delivery Services")
	}
	ds := testData.DeliveryServices[0]
	if ds.XMLID == nil {
		t.Fatal("Testing Delivery Service had no XMLID")
	}
	xmlid := *ds.XMLID + "-topology-test"

	ds.XMLID = util.StrPtr("")
	_, _, err := TOSession.CreateDeliveryServiceV30(ds)
	if err == nil {
		t.Error("Expected error with empty xmlid")
	}
	ds.XMLID = nil
	_, _, err = TOSession.CreateDeliveryServiceV30(ds)
	if err == nil {
		t.Error("Expected error with nil xmlid")
	}

	ds.Topology = new(string)
	ds.XMLID = &xmlid

	_, reqInf, err := TOSession.CreateDeliveryServiceV30(ds)
	if err == nil {
		t.Error("Expected error with non-existent Topology")
	}
	if reqInf.StatusCode < 400 || reqInf.StatusCode >= 500 {
		t.Errorf("Expected client-level error creating DS with non-existent Topology, got: %d", reqInf.StatusCode)
	}
}

func CreateTestDeliveryServices(t *testing.T) {
	pl := tc.Parameter{
		ConfigFile: "remap.config",
		Name:       "location",
		Value:      "/remap/config/location/parameter/",
	}
	_, _, err := TOSession.CreateParameter(pl)
	if err != nil {
		t.Errorf("cannot create parameter: %v", err)
	}
	for _, ds := range testData.DeliveryServices {
		_, _, err = TOSession.CreateDeliveryServiceV30(ds)
		if err != nil {
			t.Errorf("could not CREATE delivery service '%s': %v", *ds.XMLID, err)
		}
	}
}

func GetTestDeliveryServicesIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	_, reqInf, err := TOSession.GetDeliveryServicesV30WithHdr(header, nil)
	if err != nil {
		t.Fatalf("could not GET Delivery Services: %v", err)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestDeliveryServices(t *testing.T) {
	actualDSes, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v - %v", err, actualDSes)
	}
	actualDSMap := make(map[string]tc.DeliveryServiceNullableV30, len(actualDSes))
	for _, ds := range actualDSes {
		actualDSMap[*ds.XMLID] = ds
	}
	cnt := 0
	for _, ds := range testData.DeliveryServices {
		if _, ok := actualDSMap[*ds.XMLID]; !ok {
			t.Errorf("GET DeliveryService missing: %v", ds.XMLID)
		}
		// exactly one ds should have exactly 3 query params. the rest should have none
		if c := len(ds.ConsistentHashQueryParams); c > 0 {
			if c != 3 {
				t.Errorf("deliveryservice %s has %d query params; expected %d or %d", *ds.XMLID, c, 3, 0)
			}
			cnt++
		}
	}
	if cnt > 2 {
		t.Errorf("exactly 2 deliveryservices should have more than one query param; found %d", cnt)
	}
}

func GetInactiveTestDeliveryServices(t *testing.T) {
	params := url.Values{}
	params.Set("active", strconv.FormatBool(false))
	inactiveDSes, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, params)
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v - %v", err, inactiveDSes)
	}
	for _, ds := range inactiveDSes {
		if *ds.Active != false {
			t.Errorf("expected all delivery services to be inactive, but got atleast one active DS")
		}
	}
	params.Set("active", strconv.FormatBool(true))
	activeDSes, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, params)
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v - %v", err, activeDSes)
	}
	for _, ds := range activeDSes {
		if *ds.Active != true {
			t.Errorf("expected all delivery services to be active, but got atleast one inactive DS")
		}
	}
}

func GetTestDeliveryServicesCapacity(t *testing.T) {
	actualDSes, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v - %v", err, actualDSes)
	}
	actualDSMap := map[string]tc.DeliveryServiceNullableV30{}
	for _, ds := range actualDSes {
		actualDSMap[*ds.XMLID] = ds
		capDS, _, err := TOSession.GetDeliveryServiceCapacityWithHdr(strconv.Itoa(*ds.ID), nil)
		if err != nil {
			t.Errorf("cannot GET DeliveryServices: %v's Capacity: %v - %v", ds, err, capDS)
		}
	}

}

func UpdateTestDeliveryServices(t *testing.T) {
	firstDS := testData.DeliveryServices[0]

	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		t.Errorf("cannot GET Delivery Services: %v", err)
	}

	remoteDS := tc.DeliveryServiceNullableV30{}
	found := false
	for _, ds := range dses {
		if *ds.XMLID == *firstDS.XMLID {
			found = true
			remoteDS = ds
			break
		}
	}
	if !found {
		t.Errorf("GET Delivery Services missing: %v", firstDS.XMLID)
	}

	updatedMaxRequestHeaderSize := 131080
	updatedLongDesc := "something different"
	updatedMaxDNSAnswers := 164598
	updatedMaxOriginConnections := 100
	remoteDS.LongDesc = &updatedLongDesc
	remoteDS.MaxDNSAnswers = &updatedMaxDNSAnswers
	remoteDS.MaxOriginConnections = &updatedMaxOriginConnections
	remoteDS.MatchList = nil // verify that this field is optional in a PUT request, doesn't cause nil dereference panic
	remoteDS.MaxRequestHeaderBytes = &updatedMaxRequestHeaderSize

	if updateResp, _, err := TOSession.UpdateDeliveryServiceV30WithHdr(*remoteDS.ID, remoteDS, nil); err != nil {
		t.Errorf("cannot UPDATE DeliveryService by ID: %v - %v", err, updateResp)
	}

	// Retrieve the server to check rack and interfaceName values were updated
	params := url.Values{}
	params.Set("id", strconv.Itoa(*remoteDS.ID))
	apiResp, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, params)
	if err != nil {
		t.Fatalf("cannot GET Delivery Service by ID: %v - %v", remoteDS.XMLID, err)
	}
	if len(apiResp) < 1 {
		t.Fatalf("cannot GET Delivery Service by ID: %v - nil", remoteDS.XMLID)
	}
	resp := apiResp[0]

	if *resp.LongDesc != updatedLongDesc || *resp.MaxDNSAnswers != updatedMaxDNSAnswers || *resp.MaxOriginConnections != updatedMaxOriginConnections || *resp.MaxRequestHeaderBytes != updatedMaxRequestHeaderSize {
		t.Errorf("long description do not match actual: %s, expected: %s", *resp.LongDesc, updatedLongDesc)
		t.Errorf("max DNS answers do not match actual: %v, expected: %v", resp.MaxDNSAnswers, updatedMaxDNSAnswers)
		t.Errorf("max origin connections do not match actual: %v, expected: %v", resp.MaxOriginConnections, updatedMaxOriginConnections)
		t.Errorf("max request header sizes do not match actual: %v, expected: %v", resp.MaxRequestHeaderBytes, updatedMaxRequestHeaderSize)
	}
}

func UpdateNullableTestDeliveryServices(t *testing.T) {
	firstDS := testData.DeliveryServices[0]

	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		t.Fatalf("cannot GET Delivery Services: %v", err)
	}

	var remoteDS tc.DeliveryServiceNullableV30
	found := false
	for _, ds := range dses {
		if ds.XMLID == nil || ds.ID == nil {
			continue
		}
		if *ds.XMLID == *firstDS.XMLID {
			found = true
			remoteDS = ds
			break
		}
	}
	if !found {
		t.Fatalf("GET Delivery Services missing: %v", firstDS.XMLID)
	}

	updatedLongDesc := "something else different"
	updatedMaxDNSAnswers := 164599
	remoteDS.LongDesc = &updatedLongDesc
	remoteDS.MaxDNSAnswers = &updatedMaxDNSAnswers

	if updateResp, _, err := TOSession.UpdateDeliveryServiceV30WithHdr(*remoteDS.ID, remoteDS, nil); err != nil {
		t.Fatalf("cannot UPDATE DeliveryService by ID: %v - %v", err, updateResp)
	}

	params := url.Values{}
	params.Set("id", strconv.Itoa(*remoteDS.ID))
	apiResp, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, params)
	if err != nil {
		t.Fatalf("cannot GET Delivery Service by ID: %v - %v", remoteDS.XMLID, err)
	}
	if apiResp == nil {
		t.Fatalf("cannot GET Delivery Service by ID: %v - nil", remoteDS.XMLID)
	}
	resp := apiResp[0]

	if resp.LongDesc == nil || resp.MaxDNSAnswers == nil {
		t.Errorf("results do not match actual: %v, expected: %s", resp.LongDesc, updatedLongDesc)
		t.Fatalf("results do not match actual: %v, expected: %d", resp.MaxDNSAnswers, updatedMaxDNSAnswers)
	}

	if *resp.LongDesc != updatedLongDesc || *resp.MaxDNSAnswers != updatedMaxDNSAnswers {
		t.Errorf("results do not match actual: %s, expected: %s", *resp.LongDesc, updatedLongDesc)
		t.Fatalf("results do not match actual: %d, expected: %d", *resp.MaxDNSAnswers, updatedMaxDNSAnswers)
	}
}

// UpdateDeliveryServiceWithInvalidTopology ensures that a topology cannot be:
// - assigned to (CLIENT_)STEERING delivery services
// - assigned to any delivery services which have required capabilities that the topology can't satisfy
func UpdateDeliveryServiceWithInvalidTopology(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		t.Fatalf("cannot GET Delivery Services: %v", err)
	}

	found := false
	var nonCSDS *tc.DeliveryServiceNullableV30
	for _, ds := range dses {
		if ds.Type == nil || ds.ID == nil {
			continue
		}
		if *ds.Type == tc.DSTypeClientSteering {
			found = true
			ds.Topology = util.StrPtr("my-topology")
			if _, _, err := TOSession.UpdateDeliveryServiceV30WithHdr(*ds.ID, ds, nil); err == nil {
				t.Errorf("assigning topology to CLIENT_STEERING delivery service - expected: error, actual: no error")
			}
		} else if nonCSDS == nil {
			nonCSDS = new(tc.DeliveryServiceNullableV30)
			*nonCSDS = ds
		}
	}
	if !found {
		t.Error("expected at least one CLIENT_STEERING delivery service")
	}
	if nonCSDS == nil {
		t.Fatal("Expected at least on non-CLIENT_STEERING DS to exist")
	}

	nonCSDS.Topology = new(string)
	_, inf, err := TOSession.UpdateDeliveryServiceV30WithHdr(*nonCSDS.ID, *nonCSDS, nil)
	if err == nil {
		t.Error("Expected an error assigning a non-existent topology")
	}
	if inf.StatusCode < 400 || inf.StatusCode >= 500 {
		t.Errorf("Expected client-level error assigning a non-existent topology, got: %d", inf.StatusCode)
	}

	params := url.Values{}
	params.Add("xmlId", "ds-top-req-cap")
	dses, _, err = TOSession.GetDeliveryServicesV30WithHdr(nil, params)
	if err != nil {
		t.Fatalf("cannot GET delivery service: %v", err)
	}
	if len(dses) != 1 {
		t.Fatalf("expected: 1 DS, actual: %d", len(dses))
	}
	ds := dses[0]
	// unassign its topology, add a required capability that its topology
	// can't satisfy, then attempt to reassign its topology
	top := *ds.Topology
	ds.Topology = nil
	_, _, err = TOSession.UpdateDeliveryServiceV30WithHdr(*ds.ID, ds, nil)
	if err != nil {
		t.Fatalf("updating DS to remove topology, expected: no error, actual: %v", err)
	}
	reqCap := tc.DeliveryServicesRequiredCapability{
		DeliveryServiceID:  ds.ID,
		RequiredCapability: util.StrPtr("asdf"),
	}
	_, _, err = TOSession.CreateDeliveryServicesRequiredCapability(reqCap)
	if err != nil {
		t.Fatalf("adding 'asdf' required capability to '%s', expected: no error, actual: %v", *ds.XMLID, err)
	}
	ds.Topology = &top
	_, reqInf, err := TOSession.UpdateDeliveryServiceV30WithHdr(*ds.ID, ds, nil)
	if err == nil {
		t.Errorf("updating DS topology which doesn't meet the DS required capabilities - expected: error, actual: nil")
	}
	if reqInf.StatusCode < http.StatusBadRequest || reqInf.StatusCode >= http.StatusInternalServerError {
		t.Errorf("updating DS topology which doesn't meet the DS required capabilities - expected: 400-level status code, actual: %d", reqInf.StatusCode)
	}
	_, _, err = TOSession.DeleteDeliveryServicesRequiredCapability(*ds.ID, "asdf")
	if err != nil {
		t.Fatalf("removing 'asdf' required capability from '%s', expected: no error, actual: %v", *ds.XMLID, err)
	}
	_, _, err = TOSession.UpdateDeliveryServiceV30WithHdr(*ds.ID, ds, nil)
	if err != nil {
		t.Errorf("updating DS topology - expected: no error, actual: %v", err)
	}

	const xmlId = "top-ds-in-cdn2"
	dses, _, err = TOSession.GetDeliveryServicesV30WithHdr(nil, url.Values{"xmlId": {xmlId}})
	if err != nil {
		t.Fatalf("getting Delivery Service %s: %s", xmlId, err.Error())
	}
	const expectedSize = 1
	if len(dses) != expectedSize {
		t.Fatalf("expected %d Delivery Service with xmlId %s but instead received %d Delivery Services", expectedSize, xmlId, len(dses))
	}
	ds = dses[0]
	dsTopology := ds.Topology
	ds.Topology = nil
	ds, _, err = TOSession.UpdateDeliveryServiceV30WithHdr(*ds.ID, ds, nil)
	if err != nil {
		t.Fatalf("updating Delivery Service %s: %s", xmlId, err.Error())
	}
	const cdn1Name = "cdn1"
	cdns, _, err := TOSession.GetCDNByNameWithHdr(cdn1Name, nil)
	if err != nil {
		t.Fatalf("getting CDN %s: %s", cdn1Name, err.Error())
	}
	if len(cdns) != expectedSize {
		t.Fatalf("expected %d CDN with name %s but instead received %d CDNs", expectedSize, cdn1Name, len(cdns))
	}
	cdn1 := cdns[0]
	const cacheGroupName = "dtrc1"
	cachegroups, _, err := TOSession.GetCacheGroupsByQueryParamsWithHdr(url.Values{"name": {cacheGroupName}}, nil)
	if err != nil {
		t.Fatalf("getting Cache Group %s: %s", cacheGroupName, err.Error())
	}
	if len(cachegroups) != expectedSize {
		t.Fatalf("expected %d Cache Group with name %s but instead received %d Cache Groups", expectedSize, cacheGroupName, len(cachegroups))
	}
	cachegroup := cachegroups[0]
	params = url.Values{"cdn": {strconv.Itoa(*ds.CDNID)}, "cachegroup": {strconv.Itoa(*cachegroup.ID)}}
	servers, _, err := TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("getting Server with params %v: %s", params, err.Error())
	}
	if len(servers.Response) != expectedSize {
		t.Fatalf("expected %d Server returned for query params %v but instead received %d Servers", expectedSize, params, len(servers.Response))
	}
	server := servers.Response[0]
	*server.CDNID = cdn1.ID

	// A profile specific to CDN 1 is required
	profileCopy := tc.ProfileCopy{
		Name:         *server.Profile + "_BUT_IN_CDN1",
		ExistingID:   *server.ProfileID,
		ExistingName: *server.Profile,
		Description:  *server.ProfileDesc,
	}
	_, _, err = TOSession.CopyProfile(profileCopy)
	if err != nil {
		t.Fatalf("copying Profile %s: %s", *server.Profile, err.Error())
	}

	profiles, _, err := TOSession.GetProfileByNameWithHdr(profileCopy.Name, nil)
	if err != nil {
		t.Fatalf("getting Profile %s: %s", profileCopy.Name, err.Error())
	}
	if len(profiles) != expectedSize {
		t.Fatalf("expected %d Profile with name %s but instead received %d Profiles", expectedSize, profileCopy.Name, len(profiles))
	}
	profile := profiles[0]
	profile.CDNID = cdn1.ID
	_, _, err = TOSession.UpdateProfileByIDWithHdr(profile.ID, profile, nil)
	if err != nil {
		t.Fatalf("updating Profile %s: %s", profile.Name, err.Error())
	}
	*server.ProfileID = profile.ID

	// Empty Cache Group dtrc1 with respect to CDN 2
	_, _, err = TOSession.UpdateServerByIDWithHdr(*server.ID, server, nil)
	if err != nil {
		t.Fatalf("updating Server %s: %s", *server.HostName, err.Error())
	}
	ds.Topology = dsTopology
	_, reqInf, err = TOSession.UpdateDeliveryServiceV30WithHdr(*ds.ID, ds, nil)
	if err == nil {
		t.Fatalf("expected 400-level error assigning Topology %s to Delivery Service %s because Cache Group %s has no Servers in it in CDN %d, no error received", *dsTopology, xmlId, cacheGroupName, *ds.CDNID)
	}
	if reqInf.StatusCode < http.StatusBadRequest || reqInf.StatusCode >= http.StatusInternalServerError {
		t.Fatalf("expected %d-level status code but received status code %d", http.StatusBadRequest, reqInf.StatusCode)
	}
	*server.CDNID = *ds.CDNID
	*server.ProfileID = profileCopy.ExistingID

	// Put things back the way they were
	_, _, err = TOSession.UpdateServerByIDWithHdr(*server.ID, server, nil)
	if err != nil {
		t.Fatalf("updating Server %s: %s", *server.HostName, err.Error())
	}

	_, _, err = TOSession.DeleteProfileByID(profile.ID)
	if err != nil {
		t.Fatalf("deleting Profile %s: %s", profile.Name, err.Error())
	}

	ds, reqInf, err = TOSession.UpdateDeliveryServiceV30WithHdr(*ds.ID, ds, nil)
	if err != nil {
		t.Fatalf("updating Delivery Service %s: %s", xmlId, err.Error())
	}
}

// UpdateDeliveryServiceTopologyHeaderRewriteFields ensures that a delivery service can only use firstHeaderRewrite,
// innerHeaderRewrite, or lastHeadeRewrite if a topology is assigned.
func UpdateDeliveryServiceTopologyHeaderRewriteFields(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		t.Fatalf("cannot GET Delivery Services: %v", err)
	}
	foundTopology := false
	for _, ds := range dses {
		if ds.Topology != nil {
			foundTopology = true
		}
		ds.FirstHeaderRewrite = util.StrPtr("foo")
		ds.InnerHeaderRewrite = util.StrPtr("bar")
		ds.LastHeaderRewrite = util.StrPtr("baz")
		_, _, err := TOSession.UpdateDeliveryServiceV30WithHdr(*ds.ID, ds, nil)
		if ds.Topology != nil && err != nil {
			t.Errorf("expected: no error updating topology-based header rewrite fields for topology-based DS, actual: %v", err)
		}
		if ds.Topology == nil && err == nil {
			t.Errorf("expected: error updating topology-based header rewrite fields for non-topology-based DS, actual: nil")
		}
		ds.FirstHeaderRewrite = nil
		ds.InnerHeaderRewrite = nil
		ds.LastHeaderRewrite = nil
		ds.EdgeHeaderRewrite = util.StrPtr("foo")
		ds.MidHeaderRewrite = util.StrPtr("bar")
		_, _, err = TOSession.UpdateDeliveryServiceV30WithHdr(*ds.ID, ds, nil)
		if ds.Topology != nil && err == nil {
			t.Errorf("expected: error updating legacy header rewrite fields for topology-based DS, actual: nil")
		}
		if ds.Topology == nil && err != nil {
			t.Errorf("expected: no error updating legacy header rewrite fields for non-topology-based DS, actual: %v", err)
		}
	}
	if !foundTopology {
		t.Errorf("expected: at least one topology-based delivery service, actual: none found")
	}
}

// UpdateDeliveryServiceWithInvalidRemapText ensures that a delivery service can't be updated with a remap text value with a line break in it.
func UpdateDeliveryServiceWithInvalidRemapText(t *testing.T) {
	firstDS := testData.DeliveryServices[0]

	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		t.Fatalf("cannot GET Delivery Services: %v", err)
	}

	var remoteDS tc.DeliveryServiceNullableV30
	found := false
	for _, ds := range dses {
		if ds.XMLID == nil || ds.ID == nil {
			continue
		}
		if *ds.XMLID == *firstDS.XMLID {
			found = true
			remoteDS = ds
			break
		}
	}
	if !found {
		t.Fatalf("GET Delivery Services missing: %v", firstDS.XMLID)
	}

	updatedRemapText := "@plugin=tslua.so @pparam=/opt/trafficserver/etc/trafficserver/remapPlugin1.lua\nline2"
	remoteDS.RemapText = &updatedRemapText

	if _, _, err := TOSession.UpdateDeliveryServiceV30WithHdr(*remoteDS.ID, remoteDS, nil); err == nil {
		t.Errorf("Delivery service updated with invalid remap text: %v", updatedRemapText)
	}
}

// UpdateDeliveryServiceWithInvalidSliceRangeRequest ensures that a delivery service can't be updated with a invalid slice range request handler setting.
func UpdateDeliveryServiceWithInvalidSliceRangeRequest(t *testing.T) {
	// GET a HTTP / DNS type DS
	var dsXML *string
	for _, ds := range testData.DeliveryServices {
		if ds.Type.IsDNS() || ds.Type.IsHTTP() {
			dsXML = ds.XMLID
			break
		}
	}
	if dsXML == nil {
		t.Fatal("no HTTP or DNS Delivery Services to test with")
	}

	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		t.Fatalf("cannot GET Delivery Services: %v", err)
	}

	var remoteDS tc.DeliveryServiceNullableV30
	found := false
	for _, ds := range dses {
		if ds.XMLID == nil || ds.ID == nil {
			continue
		}
		if *ds.XMLID == *dsXML {
			found = true
			remoteDS = ds
			break
		}
	}
	if !found {
		t.Fatalf("GET Delivery Services missing: %v", *dsXML)
	}
	testCases := []struct {
		description         string
		rangeRequestSetting *int
		slicePluginSize     *int
	}{
		{
			description:         "Missing slice plugin size",
			rangeRequestSetting: util.IntPtr(tc.RangeRequestHandlingSlice),
			slicePluginSize:     nil,
		},
		{
			description:         "Slice plugin size set with incorrect range request setting",
			rangeRequestSetting: util.IntPtr(tc.RangeRequestHandlingBackgroundFetch),
			slicePluginSize:     util.IntPtr(262144),
		},
		{
			description:         "Slice plugin size set to small",
			rangeRequestSetting: util.IntPtr(tc.RangeRequestHandlingSlice),
			slicePluginSize:     util.IntPtr(0),
		},
		{
			description:         "Slice plugin size set to large",
			rangeRequestSetting: util.IntPtr(tc.RangeRequestHandlingSlice),
			slicePluginSize:     util.IntPtr(40000000),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			remoteDS.RangeSliceBlockSize = tc.slicePluginSize
			remoteDS.RangeRequestHandling = tc.rangeRequestSetting
			if _, _, err := TOSession.UpdateDeliveryServiceV30WithHdr(*remoteDS.ID, remoteDS, nil); err == nil {
				t.Error("Delivery service updated with invalid slice plugin configuration")
			}
		})
	}

}

// UpdateValidateORGServerCacheGroup validates ORG server's cachegroup are part of topology's cachegroup
func UpdateValidateORGServerCacheGroup(t *testing.T) {
	params := url.Values{}
	params.Set("xmlId", "ds-top")

	//Get the correct DS
	remoteDS, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, params)
	if err != nil {
		t.Errorf("cannot GET Delivery Services: %v", err)
	}

	//Assign ORG server to DS
	assignServer := []string{"denver-mso-org-01"}
	_, _, err = TOSession.AssignServersToDeliveryService(assignServer, *remoteDS[0].XMLID)
	if err != nil {
		t.Errorf("cannot assign server to Delivery Services: %v", err)
	}

	//Update DS's Topology to a non-ORG server's cachegroup
	origTopo := *remoteDS[0].Topology
	remoteDS[0].Topology = util.StrPtr("another-topology")
	ds, reqInf, err := TOSession.UpdateDeliveryServiceV30WithHdr(*remoteDS[0].ID, remoteDS[0], nil)
	if err == nil {
		t.Errorf("shouldnot UPDATE DeliveryService by ID: %v, but update was successful", ds.XMLID)
	} else if !strings.Contains(err.Error(), "the following ORG server cachegroups are not in the delivery service's topology") {
		t.Errorf("expected: error containing \"the following ORG server cachegroups are not in the delivery service's topology\", actual: %s", err.Error())
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected to fail since ORG server's topology not part of DS. Expected:%v, Got: :%v", http.StatusBadRequest, reqInf.StatusCode)
	}

	// Retrieve the DS to check if topology was updated with missing ORG server
	params.Set("id", strconv.Itoa(*remoteDS[0].ID))
	apiResp, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, params)
	if err != nil {
		t.Fatalf("cannot GET Delivery Service by ID: %v - %v", *remoteDS[0].XMLID, err)
	}
	if len(apiResp) < 1 {
		t.Fatalf("cannot GET Delivery Service by ID: %v - nil", *remoteDS[0].XMLID)
	}

	//Set topology back to as it was for further testing
	remoteDS[0].Topology = &origTopo
	_, _, err = TOSession.UpdateDeliveryServiceV30WithHdr(*remoteDS[0].ID, remoteDS[0], nil)
	if err != nil {
		t.Fatalf("couldn't update topology:%v, %v", *remoteDS[0].Topology, err)
	}
}

func GetAccessibleToTest(t *testing.T) {
	//Every delivery service is associated with the root tenant
	err := getByTenants(1, len(testData.DeliveryServices))
	if err != nil {
		t.Fatal(err.Error())
	}

	tenant := &tc.Tenant{
		Active:     true,
		Name:       "the strongest",
		ParentID:   1,
		ParentName: "root",
	}

	resp, err := TOSession.CreateTenant(tenant)
	if err != nil {
		t.Fatal(err.Error())
	}
	if resp == nil {
		t.Fatal("unexpected null response when creating tenant")
	}
	tenant = &resp.Response

	//No delivery services are associated with this new tenant
	err = getByTenants(tenant.ID, 0)
	if err != nil {
		t.Fatal(err.Error())
	}

	//First and only child tenant, no access to root
	childTenant, _, err := TOSession.TenantByName("tenant1")
	if err != nil {
		t.Fatal("unable to get tenant " + err.Error())
	}
	// TODO: document that all DSes added to the fixture data need to have the
	// Tenant 'tenant1' unless you change this code
	err = getByTenants(childTenant.ID, len(testData.DeliveryServices)-1)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = TOSession.DeleteTenant(strconv.Itoa(tenant.ID))
	if err != nil {
		t.Fatalf("unable to clean up tenant %v", err.Error())
	}
}

func getByTenants(tenantID int, expectedCount int) error {
	params := url.Values{}
	params.Set("accessibleTo", strconv.Itoa(tenantID))
	deliveryServices, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, params)
	if err != nil {
		return err
	}
	if len(deliveryServices) != expectedCount {
		return fmt.Errorf("expected %v delivery service, got %v", expectedCount, len(deliveryServices))
	}
	return nil
}

func DeleteTestDeliveryServices(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		t.Errorf("cannot GET deliveryservices: %v", err)
	}
	for _, testDS := range testData.DeliveryServices {
		if testDS.XMLID == nil {
			t.Errorf("testing Delivery Service has no XMLID")
			continue
		}
		var ds tc.DeliveryServiceNullableV30
		found := false
		for _, realDS := range dses {
			if realDS.XMLID != nil && *realDS.XMLID == *testDS.XMLID {
				ds = realDS
				found = true
				break
			}
		}
		if !found {
			t.Errorf("DeliveryService not found in Traffic Ops: %v", *testDS.XMLID)
			continue
		}

		delResp, err := TOSession.DeleteDeliveryService(strconv.Itoa(*ds.ID))
		if err != nil {
			t.Errorf("cannot DELETE DeliveryService by ID: %v - %v", err, delResp)
			continue
		}

		// Retrieve the Server to see if it got deleted
		params := url.Values{}
		params.Set("id", strconv.Itoa(*ds.ID))
		foundDS, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, params)
		if err != nil {
			t.Errorf("Unexpected error deleting Delivery Service '%s': %v", *ds.XMLID, err)
		}
		if len(foundDS) > 0 {
			t.Errorf("expected Delivery Service: %s to be deleted, but %d exist with same ID (#%d)", *ds.XMLID, len(foundDS), *ds.ID)
		}
	}

	// clean up parameter created in CreateTestDeliveryServices()
	params, _, err := TOSession.GetParameterByNameAndConfigFile("location", "remap.config")
	for _, param := range params {
		deleted, _, err := TOSession.DeleteParameterByID(param.ID)
		if err != nil {
			t.Errorf("cannot DELETE parameter by ID (%d): %v - %v", param.ID, err, deleted)
		}
	}
}

func DeliveryServiceMinorVersionsTest(t *testing.T) {
	if len(testData.DeliveryServices) < 5 {
		t.Fatalf("Need at least 5 DSes to test minor versions; got: %d", len(testData.DeliveryServices))
	}
	testDS := testData.DeliveryServices[4]
	if testDS.XMLID == nil {
		t.Fatal("expected XMLID: ds-test-minor-versions, actual: <nil>")
	}
	if *testDS.XMLID != "ds-test-minor-versions" {
		t.Errorf("expected XMLID: ds-test-minor-versions, actual: %s", *testDS.XMLID)
	}

	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v - %v", err, dses)
	}

	var ds tc.DeliveryServiceNullableV30
	found := false
	for _, d := range dses {
		if d.XMLID != nil && *d.XMLID == *testDS.XMLID {
			ds = d
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Delivery Service '%s' not found in Traffic Ops", *testDS.XMLID)
	}

	// GET latest, verify expected values for 1.3 and 1.4 fields
	if ds.DeepCachingType == nil {
		t.Errorf("expected DeepCachingType: %s, actual: nil", testDS.DeepCachingType.String())
	} else if *ds.DeepCachingType != *testDS.DeepCachingType {
		t.Errorf("expected DeepCachingType: %s, actual: %s", testDS.DeepCachingType.String(), ds.DeepCachingType.String())
	}
	if ds.FQPacingRate == nil {
		t.Errorf("expected FQPacingRate: %d, actual: nil", testDS.FQPacingRate)
	} else if *ds.FQPacingRate != *testDS.FQPacingRate {
		t.Errorf("expected FQPacingRate: %d, actual: %d", testDS.FQPacingRate, *ds.FQPacingRate)
	}
	if ds.SigningAlgorithm == nil {
		t.Errorf("expected SigningAlgorithm: %s, actual: nil", *testDS.SigningAlgorithm)
	} else if *ds.SigningAlgorithm != *testDS.SigningAlgorithm {
		t.Errorf("expected SigningAlgorithm: %s, actual: %s", *testDS.SigningAlgorithm, *ds.SigningAlgorithm)
	}
	if ds.Tenant == nil {
		t.Errorf("expected Tenant: %s, actual: nil", *testDS.Tenant)
	} else if *ds.Tenant != *testDS.Tenant {
		t.Errorf("expected Tenant: %s, actual: %s", *testDS.Tenant, *ds.Tenant)
	}
	if ds.TRRequestHeaders == nil {
		t.Errorf("expected TRRequestHeaders: %s, actual: nil", *testDS.TRRequestHeaders)
	} else if *ds.TRRequestHeaders != *testDS.TRRequestHeaders {
		t.Errorf("expected TRRequestHeaders: %s, actual: %s", *testDS.TRRequestHeaders, *ds.TRRequestHeaders)
	}
	if ds.TRResponseHeaders == nil {
		t.Errorf("expected TRResponseHeaders: %s, actual: nil", *testDS.TRResponseHeaders)
	} else if *ds.TRResponseHeaders != *testDS.TRResponseHeaders {
		t.Errorf("expected TRResponseHeaders: %s, actual: %s", *testDS.TRResponseHeaders, *ds.TRResponseHeaders)
	}
	if ds.ConsistentHashRegex == nil {
		t.Errorf("expected ConsistentHashRegex: %s, actual: nil", *testDS.ConsistentHashRegex)
	} else if *ds.ConsistentHashRegex != *testDS.ConsistentHashRegex {
		t.Errorf("expected ConsistentHashRegex: %s, actual: %s", *testDS.ConsistentHashRegex, *ds.ConsistentHashRegex)
	}
	if ds.ConsistentHashQueryParams == nil {
		t.Errorf("expected ConsistentHashQueryParams: %v, actual: nil", testDS.ConsistentHashQueryParams)
	} else if !reflect.DeepEqual(ds.ConsistentHashQueryParams, testDS.ConsistentHashQueryParams) {
		t.Errorf("expected ConsistentHashQueryParams: %v, actual: %v", testDS.ConsistentHashQueryParams, ds.ConsistentHashQueryParams)
	}
	if ds.MaxOriginConnections == nil {
		t.Errorf("expected MaxOriginConnections: %d, actual: nil", testDS.MaxOriginConnections)
	} else if *ds.MaxOriginConnections != *testDS.MaxOriginConnections {
		t.Errorf("expected MaxOriginConnections: %d, actual: %d", testDS.MaxOriginConnections, *ds.MaxOriginConnections)
	}

	ds.ID = nil
	_, err = json.Marshal(ds)
	if err != nil {
		t.Errorf("cannot POST deliveryservice, failed to marshal JSON: %s", err.Error())
	}
}

func DeliveryServiceTenancyTest(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		t.Errorf("cannot GET deliveryservices: %v", err)
	}
	var tenant3DS tc.DeliveryServiceNullableV30
	foundTenant3DS := false
	for _, d := range dses {
		if *d.XMLID == "ds3" {
			tenant3DS = d
			foundTenant3DS = true
		}
	}
	if !foundTenant3DS || *tenant3DS.Tenant != "tenant3" {
		t.Error("expected to find deliveryservice 'ds3' with tenant 'tenant3'")
	}

	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	tenant4TOClient, _, err := toclient.LoginWithAgent(TOSession.URL, "tenant4user", "pa$$word", true, "to-api-v3-client-tests/tenant4user", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to log in with tenant4user: %v", err.Error())
	}

	dsesReadableByTenant4, _, err := tenant4TOClient.GetDeliveryServicesNullable()
	if err != nil {
		t.Error("tenant4user cannot GET deliveryservices")
	}

	// assert that tenant4user cannot read deliveryservices outside of its tenant
	for _, ds := range dsesReadableByTenant4 {
		if *ds.XMLID == "ds3" {
			t.Error("expected tenant4 to be unable to read delivery services from tenant 3")
		}
	}

	// assert that tenant4user cannot update tenant3user's deliveryservice
	if _, _, err = tenant4TOClient.UpdateDeliveryServiceV30WithHdr(*tenant3DS.ID, tenant3DS, nil); err == nil {
		t.Errorf("expected tenant4user to be unable to update tenant3's deliveryservice (%s)", *tenant3DS.XMLID)
	}

	// assert that tenant4user cannot delete tenant3user's deliveryservice
	if _, err = tenant4TOClient.DeleteDeliveryService(strconv.Itoa(*tenant3DS.ID)); err == nil {
		t.Errorf("expected tenant4user to be unable to delete tenant3's deliveryservice (%s)", *tenant3DS.XMLID)
	}

	// assert that tenant4user cannot create a deliveryservice outside of its tenant
	tenant3DS.XMLID = util.StrPtr("deliveryservicetenancytest")
	tenant3DS.DisplayName = util.StrPtr("deliveryservicetenancytest")
	if _, _, err = tenant4TOClient.CreateDeliveryServiceV30(tenant3DS); err == nil {
		t.Error("expected tenant4user to be unable to create a deliveryservice outside of its tenant")
	}
}

func GetTestDeliveryServicesURLSigKeys(t *testing.T) {
	if len(testData.DeliveryServices) == 0 {
		t.Fatal("couldn't get the xml ID of test DS")
	}
	firstDS := testData.DeliveryServices[0]
	if firstDS.XMLID == nil {
		t.Fatal("couldn't get the xml ID of test DS")
	}

	_, _, err := TOSession.GetDeliveryServiceURLSigKeysWithHdr(*firstDS.XMLID, nil)
	if err != nil {
		t.Error("failed to get url sig keys: " + err.Error())
	}
}
