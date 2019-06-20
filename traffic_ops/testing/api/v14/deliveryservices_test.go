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
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
)

func TestDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		UpdateTestDeliveryServices(t)
		UpdateNullableTestDeliveryServices(t)
		GetTestDeliveryServices(t)
		DeliveryServiceTenancyTest(t)
	})
}

func CreateTestDeliveryServices(t *testing.T) {
	pl := tc.Parameter{
		ConfigFile: "remap.config",
		Name:       "location",
		Value:      "/remap/config/location/parameter/",
	}
	_, _, err := TOSession.CreateParameter(pl)
	if err != nil {
		t.Errorf("cannot create parameter: %v\n", err)
	}
	for _, ds := range testData.DeliveryServices {
		_, err = TOSession.CreateDeliveryService(&ds)
		if err != nil {
			t.Errorf("could not CREATE delivery service '%s': %v\n", ds.XMLID, err)
		}
	}
}

func GetTestDeliveryServices(t *testing.T) {
	actualDSes, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v - %v\n", err, actualDSes)
	}
	actualDSMap := map[string]tc.DeliveryService{}
	for _, ds := range actualDSes {
		actualDSMap[ds.XMLID] = ds
	}
	cnt := 0
	for _, ds := range testData.DeliveryServices {
		if _, ok := actualDSMap[ds.XMLID]; !ok {
			t.Errorf("GET DeliveryService missing: %v\n", ds.XMLID)
		}
		// exactly one ds should have exactly 3 query params. the rest should have none
		if c := len(ds.ConsistentHashQueryParams); c > 0 {
			if c != 3 {
				t.Errorf("deliveryservice %s has %d query params; expected %d or %d", ds.XMLID, c, 3, 0)
			}
			cnt++
		}
	}
	if cnt > 1 {
		t.Errorf("exactly 1 deliveryservice should have more than one query param; found %d", cnt)
	}
}

func UpdateTestDeliveryServices(t *testing.T) {
	firstDS := testData.DeliveryServices[0]

	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Errorf("cannot GET Delivery Services: %v\n", err)
	}

	remoteDS := tc.DeliveryService{}
	found := false
	for _, ds := range dses {
		if ds.XMLID == firstDS.XMLID {
			found = true
			remoteDS = ds
			break
		}
	}
	if !found {
		t.Errorf("GET Delivery Services missing: %v\n", firstDS.XMLID)
	}

	updatedLongDesc := "something different"
	updatedMaxDNSAnswers := 164598
	updatedMaxOriginConnections := 100
	remoteDS.LongDesc = updatedLongDesc
	remoteDS.MaxDNSAnswers = updatedMaxDNSAnswers
	remoteDS.MaxOriginConnections = updatedMaxOriginConnections
	remoteDS.MatchList = nil // verify that this field is optional in a PUT request, doesn't cause nil dereference panic

	if updateResp, err := TOSession.UpdateDeliveryService(strconv.Itoa(remoteDS.ID), &remoteDS); err != nil {
		t.Errorf("cannot UPDATE DeliveryService by ID: %v - %v\n", err, updateResp)
	}

	// Retrieve the server to check rack and interfaceName values were updated
	resp, _, err := TOSession.GetDeliveryService(strconv.Itoa(remoteDS.ID))
	if err != nil {
		t.Errorf("cannot GET Delivery Service by ID: %v - %v\n", remoteDS.XMLID, err)
	}
	if resp == nil {
		t.Errorf("cannot GET Delivery Service by ID: %v - nil\n", remoteDS.XMLID)
	}

	if resp.LongDesc != updatedLongDesc || resp.MaxDNSAnswers != updatedMaxDNSAnswers || resp.MaxOriginConnections != updatedMaxOriginConnections {
		t.Errorf("results do not match actual: %s, expected: %s\n", resp.LongDesc, updatedLongDesc)
		t.Errorf("results do not match actual: %v, expected: %v\n", resp.MaxDNSAnswers, updatedMaxDNSAnswers)
		t.Errorf("results do not match actual: %v, expected: %v\n", resp.MaxOriginConnections, updatedMaxOriginConnections)
	}
}

func UpdateNullableTestDeliveryServices(t *testing.T) {
	firstDS := testData.DeliveryServices[0]

	dses, _, err := TOSession.GetDeliveryServicesNullable()
	if err != nil {
		t.Fatalf("cannot GET Delivery Services: %v\n", err)
	}

	remoteDS := tc.DeliveryServiceNullable{}
	found := false
	for _, ds := range dses {
		if ds.XMLID == nil || ds.ID == nil {
			continue
		}
		if *ds.XMLID == firstDS.XMLID {
			found = true
			remoteDS = ds
			break
		}
	}
	if !found {
		t.Fatalf("GET Delivery Services missing: %v\n", firstDS.XMLID)
	}

	updatedLongDesc := "something else different"
	updatedMaxDNSAnswers := 164599
	remoteDS.LongDesc = &updatedLongDesc
	remoteDS.MaxDNSAnswers = &updatedMaxDNSAnswers

	if updateResp, err := TOSession.UpdateDeliveryServiceNullable(strconv.Itoa(*remoteDS.ID), &remoteDS); err != nil {
		t.Fatalf("cannot UPDATE DeliveryService by ID: %v - %v\n", err, updateResp)
	}

	// Retrieve the server to check rack and interfaceName values were updated
	resp, _, err := TOSession.GetDeliveryServiceNullable(strconv.Itoa(*remoteDS.ID))
	if err != nil {
		t.Fatalf("cannot GET Delivery Service by ID: %v - %v\n", remoteDS.XMLID, err)
	}
	if resp == nil {
		t.Fatalf("cannot GET Delivery Service by ID: %v - nil\n", remoteDS.XMLID)
	}

	if resp.LongDesc == nil || resp.MaxDNSAnswers == nil {
		t.Errorf("results do not match actual: %v, expected: %s\n", resp.LongDesc, updatedLongDesc)
		t.Fatalf("results do not match actual: %v, expected: %d\n", resp.MaxDNSAnswers, updatedMaxDNSAnswers)
	}

	if *resp.LongDesc != updatedLongDesc || *resp.MaxDNSAnswers != updatedMaxDNSAnswers {
		t.Errorf("results do not match actual: %s, expected: %s\n", *resp.LongDesc, updatedLongDesc)
		t.Fatalf("results do not match actual: %d, expected: %d\n", *resp.MaxDNSAnswers, updatedMaxDNSAnswers)
	}
}

func DeleteTestDeliveryServices(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Errorf("cannot GET deliveryservices: %v\n", err)
	}
	for _, testDS := range testData.DeliveryServices {
		ds := tc.DeliveryService{}
		found := false
		for _, realDS := range dses {
			if realDS.XMLID == testDS.XMLID {
				ds = realDS
				found = true
				break
			}
		}
		if !found {
			t.Errorf("DeliveryService not found in Traffic Ops: %v\n", ds.XMLID)
		}

		delResp, err := TOSession.DeleteDeliveryService(strconv.Itoa(ds.ID))
		if err != nil {
			t.Errorf("cannot DELETE DeliveryService by ID: %v - %v\n", err, delResp)
		}

		// Retrieve the Server to see if it got deleted
		foundDS, err := TOSession.DeliveryService(strconv.Itoa(ds.ID))
		if err == nil && foundDS != nil {
			t.Errorf("expected Delivery Service: %s to be deleted\n", ds.XMLID)
		}
	}

	// clean up parameter created in CreateTestDeliveryServices()
	params, _, err := TOSession.GetParameterByNameAndConfigFile("location", "remap.config")
	for _, param := range params {
		deleted, _, err := TOSession.DeleteParameterByID(param.ID)
		if err != nil {
			t.Errorf("cannot DELETE parameter by ID (%d): %v - %v\n", param.ID, err, deleted)
		}
	}
}

func DeliveryServiceTenancyTest(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServicesNullable()
	if err != nil {
		t.Errorf("cannot GET deliveryservices: %v\n", err)
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
	tenant4TOClient, _, err := toclient.LoginWithAgent(TOSession.URL, "tenant4user", "pa$$word", true, "to-api-v14-client-tests/tenant4user", true, toReqTimeout)
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
		t.Errorf("expected tenant4user to be unable to create a deliveryservice outside of its tenant")
	}

}
