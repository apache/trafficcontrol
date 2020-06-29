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
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
)

func TestDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices}, func() {
		GetTestDeliveryServicesIMS(t)
		GetAccessibleToTest(t)
		currentTime := time.Now().Add(-1 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		UpdateTestDeliveryServices(t)
		UpdateNullableTestDeliveryServices(t)
		UpdateDeliveryServiceWithInvalidRemapText(t)
		UpdateDeliveryServiceWithInvalidSliceRangeRequest(t)
		UpdateDeliveryServiceWithInvalidTopology(t)
		GetTestDeliveryServicesIMSAfterChange(t, header)
		UpdateDeliveryServiceTopologyHeaderRewriteFields(t)
		GetTestDeliveryServices(t)
		DeliveryServiceMinorVersionsTest(t)
		DeliveryServiceTenancyTest(t)
		PostDeliveryServiceTest(t)
	})
}

func GetTestDeliveryServicesIMSAfterChange(t *testing.T, header http.Header) {
	_, reqInf, err := TOSession.GetDeliveryServicesNullable(header)
	if err != nil {
		t.Fatalf("could not GET Delivery Services: %v", err)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
	}
	currentTime := time.Now()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	_, reqInf, err = TOSession.GetDeliveryServicesNullable(header)
	if err != nil {
		t.Fatalf("could not GET Delivery Services: %v", err)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func PostDeliveryServiceTest(t *testing.T) {
	ds := testData.DeliveryServices[0]
	ds.XMLID = util.StrPtr("")
	_, err := TOSession.CreateDeliveryServiceNullable(&ds)
	if err == nil {
		t.Fatal("Expected error with empty xmlid")
	}
	ds.XMLID = nil
	_, err = TOSession.CreateDeliveryServiceNullable(&ds)
	if err == nil {
		t.Fatal("Expected error with nil xmlid")
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
		_, err = TOSession.CreateDeliveryServiceNullable(&ds)
		if err != nil {
			t.Errorf("could not CREATE delivery service '%s': %v", *ds.XMLID, err)
		}
	}
}

func GetTestDeliveryServicesIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0,0,1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	_, reqInf, err := TOSession.GetDeliveryServicesNullable(header)
	if err != nil {
		t.Fatalf("could not GET Delivery Services: %v", err)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestDeliveryServices(t *testing.T) {
	actualDSes, _, err := TOSession.GetDeliveryServicesNullable(nil)
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v - %v", err, actualDSes)
	}
	actualDSMap := map[string]tc.DeliveryServiceNullable{}
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

func UpdateTestDeliveryServices(t *testing.T) {
	firstDS := testData.DeliveryServices[0]

	dses, _, err := TOSession.GetDeliveryServicesNullable(nil)
	if err != nil {
		t.Errorf("cannot GET Delivery Services: %v", err)
	}

	remoteDS := tc.DeliveryServiceNullable{}
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

	updatedLongDesc := "something different"
	updatedMaxDNSAnswers := 164598
	updatedMaxOriginConnections := 100
	remoteDS.LongDesc = &updatedLongDesc
	remoteDS.MaxDNSAnswers = &updatedMaxDNSAnswers
	remoteDS.MaxOriginConnections = &updatedMaxOriginConnections
	remoteDS.MatchList = nil // verify that this field is optional in a PUT request, doesn't cause nil dereference panic

	if updateResp, err := TOSession.UpdateDeliveryServiceNullable(strconv.Itoa(*remoteDS.ID), &remoteDS); err != nil {
		t.Errorf("cannot UPDATE DeliveryService by ID: %v - %v", err, updateResp)
	}

	// Retrieve the server to check rack and interfaceName values were updated
	resp, _, err := TOSession.GetDeliveryServiceNullable(strconv.Itoa(*remoteDS.ID), nil)
	if err != nil {
		t.Errorf("cannot GET Delivery Service by ID: %v - %v", remoteDS.XMLID, err)
	}
	if resp == nil {
		t.Errorf("cannot GET Delivery Service by ID: %v - nil", remoteDS.XMLID)
	}

	if *resp.LongDesc != updatedLongDesc || *resp.MaxDNSAnswers != updatedMaxDNSAnswers || *resp.MaxOriginConnections != updatedMaxOriginConnections {
		t.Errorf("results do not match actual: %s, expected: %s", *resp.LongDesc, updatedLongDesc)
		t.Errorf("results do not match actual: %v, expected: %v", resp.MaxDNSAnswers, updatedMaxDNSAnswers)
		t.Errorf("results do not match actual: %v, expected: %v", resp.MaxOriginConnections, updatedMaxOriginConnections)
	}
}

func UpdateNullableTestDeliveryServices(t *testing.T) {
	firstDS := testData.DeliveryServices[0]

	dses, _, err := TOSession.GetDeliveryServicesNullable(nil)
	if err != nil {
		t.Fatalf("cannot GET Delivery Services: %v", err)
	}

	remoteDS := tc.DeliveryServiceNullable{}
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

	if updateResp, err := TOSession.UpdateDeliveryServiceNullable(strconv.Itoa(*remoteDS.ID), &remoteDS); err != nil {
		t.Fatalf("cannot UPDATE DeliveryService by ID: %v - %v", err, updateResp)
	}

	// Retrieve the server to check rack and interfaceName values were updated
	resp, _, err := TOSession.GetDeliveryServiceNullable(strconv.Itoa(*remoteDS.ID), nil)
	if err != nil {
		t.Fatalf("cannot GET Delivery Service by ID: %v - %v", remoteDS.XMLID, err)
	}
	if resp == nil {
		t.Fatalf("cannot GET Delivery Service by ID: %v - nil", remoteDS.XMLID)
	}

	if resp.LongDesc == nil || resp.MaxDNSAnswers == nil {
		t.Errorf("results do not match actual: %v, expected: %s", resp.LongDesc, updatedLongDesc)
		t.Fatalf("results do not match actual: %v, expected: %d", resp.MaxDNSAnswers, updatedMaxDNSAnswers)
	}

	if *resp.LongDesc != updatedLongDesc || *resp.MaxDNSAnswers != updatedMaxDNSAnswers {
		t.Errorf("results do not match actual: %s, expected: %s", *resp.LongDesc, updatedLongDesc)
		t.Fatalf("results do not match actual: %d, expected: %d", *resp.MaxDNSAnswers, updatedMaxDNSAnswers)
	}
}

// UpdateDeliveryServiceWithInvalidTopology ensures that a topology cannot be assigned to (CLIENT_)STEERING delivery services.
func UpdateDeliveryServiceWithInvalidTopology(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServicesNullable(nil)
	if err != nil {
		t.Fatalf("cannot GET Delivery Services: %v", err)
	}

	found := false
	for _, ds := range dses {
		if *ds.Type == tc.DSTypeClientSteering {
			found = true
			ds.Topology = util.StrPtr("my-topology")
			if _, err := TOSession.UpdateDeliveryServiceNullable(strconv.Itoa(*ds.ID), &ds); err == nil {
				t.Errorf("assigning topology to CLIENT_STEERING delivery service - expected: error, actual: no error")
			}
		}
	}
	if !found {
		t.Fatalf("expected at least one CLIENT_STEERING delivery service")
	}
}

// UpdateDeliveryServiceTopologyHeaderRewriteFields ensures that a delivery service can only use firstHeaderRewrite,
// innerHeaderRewrite, or lastHeadeRewrite if a topology is assigned.
func UpdateDeliveryServiceTopologyHeaderRewriteFields(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServicesNullable(nil)
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
		_, err := TOSession.UpdateDeliveryServiceNullable(strconv.Itoa(*ds.ID), &ds)
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
		_, err = TOSession.UpdateDeliveryServiceNullable(strconv.Itoa(*ds.ID), &ds)
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

	dses, _, err := TOSession.GetDeliveryServicesNullable(nil)
	if err != nil {
		t.Fatalf("cannot GET Delivery Services: %v", err)
	}

	remoteDS := tc.DeliveryServiceNullable{}
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

	if _, err := TOSession.UpdateDeliveryServiceNullable(strconv.Itoa(*remoteDS.ID), &remoteDS); err == nil {
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

	dses, _, err := TOSession.GetDeliveryServicesNullable(nil)
	if err != nil {
		t.Fatalf("cannot GET Delivery Services: %v", err)
	}

	remoteDS := tc.DeliveryServiceNullable{}
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
			if _, err := TOSession.UpdateDeliveryServiceNullable(strconv.Itoa(*remoteDS.ID), &remoteDS); err == nil {
				t.Error("Delivery service updated with invalid slice plugin configuration")
			}
		})
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
	childTenant, _, err := TOSession.TenantByName("tenant1", nil)
	if err != nil {
		t.Fatal("unable to get tenant " + err.Error())
	}
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
	deliveryServices, _, err := TOSession.GetAccessibleDeliveryServicesByTenant(tenantID)
	if err != nil {
		return err
	}
	if len(deliveryServices) != expectedCount {
		return errors.New(fmt.Sprintf("expected %v delivery service, got %v", expectedCount, len(deliveryServices)))
	}
	return nil
}

func DeleteTestDeliveryServices(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServicesNullable(nil)
	if err != nil {
		t.Errorf("cannot GET deliveryservices: %v", err)
	}
	for _, testDS := range testData.DeliveryServices {
		ds := tc.DeliveryServiceNullable{}
		found := false
		for _, realDS := range dses {
			if *realDS.XMLID == *testDS.XMLID {
				ds = realDS
				found = true
				break
			}
		}
		if !found {
			t.Errorf("DeliveryService not found in Traffic Ops: %v", ds.XMLID)
		}

		delResp, err := TOSession.DeleteDeliveryService(strconv.Itoa(*ds.ID))
		if err != nil {
			t.Errorf("cannot DELETE DeliveryService by ID: %v - %v", err, delResp)
		}

		// Retrieve the Server to see if it got deleted
		foundDS, _, err := TOSession.GetDeliveryServiceNullable(strconv.Itoa(*ds.ID), nil)
		if err == nil && foundDS != nil {
			t.Errorf("expected Delivery Service: %s to be deleted", *ds.XMLID)
		}
	}

	// clean up parameter created in CreateTestDeliveryServices()
	params, _, err := TOSession.GetParameterByNameAndConfigFile("location", "remap.config", nil)
	for _, param := range params {
		deleted, _, err := TOSession.DeleteParameterByID(param.ID)
		if err != nil {
			t.Errorf("cannot DELETE parameter by ID (%d): %v - %v", param.ID, err, deleted)
		}
	}
}

func DeliveryServiceMinorVersionsTest(t *testing.T) {
	testDS := testData.DeliveryServices[4]
	if *testDS.XMLID != "ds-test-minor-versions" {
		t.Errorf("expected XMLID: ds-test-minor-versions, actual: %s", *testDS.XMLID)
	}

	dses, _, err := TOSession.GetDeliveryServicesNullable(nil)
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v - %v", err, dses)
	}
	ds := tc.DeliveryServiceNullable{}
	for _, d := range dses {
		if *d.XMLID == *testDS.XMLID {
			ds = d
			break
		}
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
	dses, _, err := TOSession.GetDeliveryServicesNullable(nil)
	if err != nil {
		t.Errorf("cannot GET deliveryservices: %v", err)
	}
	tenant3DS := tc.DeliveryServiceNullable{}
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

	dsesReadableByTenant4, _, err := tenant4TOClient.GetDeliveryServicesNullable(nil)
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
	if _, err = tenant4TOClient.UpdateDeliveryServiceNullable(string(*tenant3DS.ID), &tenant3DS); err == nil {
		t.Errorf("expected tenant4user to be unable to update tenant3's deliveryservice (%s)", *tenant3DS.XMLID)
	}

	// assert that tenant4user cannot delete tenant3user's deliveryservice
	if _, err = tenant4TOClient.DeleteDeliveryService(string(*tenant3DS.ID)); err == nil {
		t.Errorf("expected tenant4user to be unable to delete tenant3's deliveryservice (%s)", *tenant3DS.XMLID)
	}

	// assert that tenant4user cannot create a deliveryservice outside of its tenant
	tenant3DS.XMLID = util.StrPtr("deliveryservicetenancytest")
	tenant3DS.DisplayName = util.StrPtr("deliveryservicetenancytest")
	if _, err = tenant4TOClient.CreateDeliveryServiceNullable(&tenant3DS); err == nil {
		t.Error("expected tenant4user to be unable to create a deliveryservice outside of its tenant")
	}

}
