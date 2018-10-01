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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"testing"
)

func TestDeliveryServices(t *testing.T) {
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
	UpdateTestDeliveryServices(t)
	GetTestDeliveryServices(t)
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

func CreateTestDeliveryServices(t *testing.T) {
	log.Debugln("CreateTestDeliveryServices")

	pl := tc.Parameter{
		ConfigFile: "remap.config",
		Name:       "location",
		Value:      "/remap/config/location/parameter/",
	}
	_, _, err := TOSession.CreateParameter(pl)
	if err != nil {
		t.Fatalf("cannot create parameter: %v\n", err)
	}
	for _, ds := range testData.DeliveryServices {
		respCDNs, _, err := TOSession.GetCDNByName(ds.CDNName)
		if err != nil {
			t.Fatalf("cannot GET CDN - %v\n", err)
		}
		if len(respCDNs) < 1 {
			t.Fatalf("cannot GET CDN - no CDNs\n")
		}
		ds.CDNID = respCDNs[0].ID

		respTypes, _, err := TOSession.GetTypeByName(string(ds.Type))
		if err != nil {
			t.Fatalf("cannot GET Type by name: %v\n", err)
		}
		if len(respTypes) < 1 {
			t.Fatalf("cannot GET Type - no Types\n")
		}
		ds.TypeID = respTypes[0].ID

		if ds.ProfileName != "" {
			respProfiles, _, err := TOSession.GetProfileByName(ds.ProfileName)
			if err != nil {
				t.Fatalf("cannot GET Profile by name: %v\n", err)
			}
			if len(respProfiles) < 1 {
				t.Fatalf("cannot GET Profile - no Profiles\n")
			}
			ds.ProfileID = respProfiles[0].ID
		}

		respTenants, _, err := TOSession.Tenants()
		if err != nil {
			t.Fatalf("cannot GET tenants: %v\n", err)
		}
		if len(respTenants) < 1 {
			t.Fatalf("cannot GET tenants: no tenants returned from Traffic Ops\n")
		}
		ds.TenantID = respTenants[0].ID

		_, err = TOSession.CreateDeliveryService(&ds)
		if err != nil {
			t.Fatalf("could not CREATE delivery service '%s': %v\n", ds.XMLID, err)
		}
	}
}

func GetTestDeliveryServices(t *testing.T) {
	failed := false
	actualDSes, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Fatalf("cannot GET DeliveryServices: %v - %v\n", err, actualDSes)
		failed = true
	}
	actualDSMap := map[string]tc.DeliveryService{}
	for _, ds := range actualDSes {
		actualDSMap[ds.XMLID] = ds
	}
	for _, ds := range testData.DeliveryServices {
		if _, ok := actualDSMap[ds.XMLID]; !ok {
			t.Errorf("GET DeliveryService missing: %v\n", ds.XMLID)
			failed = true
		}
	}
	if !failed {
		log.Debugln("GetTestDeliveryServices() PASSED: ")
	}
}

func UpdateTestDeliveryServices(t *testing.T) {
	failed := false
	firstDS := testData.DeliveryServices[0]

	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		failed = true
		t.Fatalf("cannot GET Delivery Services: %v\n", err)
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
		failed = true
		t.Fatalf("GET Delivery Services missing: %v\n", firstDS.XMLID)
	}

	updatedLongDesc := "something different"
	updatedMaxDNSAnswers := 164598
	remoteDS.LongDesc = updatedLongDesc
	remoteDS.MaxDNSAnswers = updatedMaxDNSAnswers

	if updateResp, err := TOSession.UpdateDeliveryService(strconv.Itoa(remoteDS.ID), &remoteDS); err != nil {
		t.Errorf("cannot UPDATE DeliveryService by ID: %v - %v\n", err, updateResp)
	}

	// Retrieve the server to check rack and interfaceName values were updated
	resp, _, err := TOSession.GetDeliveryService(strconv.Itoa(remoteDS.ID))
	if err != nil {
		failed = true
		t.Fatalf("cannot GET Delivery Service by ID: %v - %v\n", remoteDS.XMLID, err)
	}
	if resp == nil {
		failed = true
		t.Fatalf("cannot GET Delivery Service by ID: %v - nil\n", remoteDS.XMLID)
	}

	if resp.LongDesc != updatedLongDesc || resp.MaxDNSAnswers != updatedMaxDNSAnswers {
		failed = true
		t.Errorf("results do not match actual: %s, expected: %s\n", resp.LongDesc, updatedLongDesc)
		t.Errorf("results do not match actual: %v, expected: %v\n", resp.MaxDNSAnswers, updatedMaxDNSAnswers)
	}
	if !failed {
		log.Debugln("UpdatedTestDeliveryServices() PASSED: ")
	}
}

func DeleteTestDeliveryServices(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServices()
	failed := false
	if err != nil {
		failed = true
		t.Fatalf("cannot GET Servers: %v\n", err)
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
			failed = true
			t.Fatalf("DeliveryService not found in Traffic Ops: %v\n", ds.XMLID)
		}

		delResp, err := TOSession.DeleteDeliveryService(strconv.Itoa(ds.ID))
		if err != nil {
			failed = true
			t.Errorf("cannot DELETE DeliveryService by ID: %v - %v\n", err, delResp)
		}

		// Retrieve the Server to see if it got deleted
		foundDS, err := TOSession.DeliveryService(strconv.Itoa(ds.ID))
		if err == nil && foundDS != nil {
			failed = true
			t.Errorf("expected Delivery Service: %s to be deleted\n", ds.XMLID)
		}
	}

	// clean up parameter created in CreateTestDeliveryServices()
	params, _, err := TOSession.GetParameterByNameAndConfigFile("location", "remap.config")
	for _, param := range params {
		deleted, _, err := TOSession.DeleteParameterByID(param.ID)
		if err != nil {
			failed = true
			t.Errorf("cannot DELETE parameter by ID (%d): %v - %v\n", param.ID, err, deleted)
		}
	}

	if !failed {
		log.Debugln("DeleteTestDeliveryServices() PASSED: ")
	}
}
