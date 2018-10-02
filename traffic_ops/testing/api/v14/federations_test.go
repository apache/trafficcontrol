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
	"testing"

	"github.com/apache/trafficcontrol/lib/go-log"
)

func TestFederations(t *testing.T) {
	CreateTestCDNs(t)
	CreateTestTypes(t)
	CreateTestProfiles(t)
	CreateTestStatuses(t)
	CreateTestDivisions(t)
	CreateTestRegions(t)
	CreateTestPhysLocations(t)
	CreateTestCacheGroups(t)
	CreateTestDeliveryServices(t)
	CreateTestUsersDeliveryServices(t)
	CreateTestCDNFederations(t)

	PostTestFederationsDeliveryServices(t)
	GetTestFederations(t)

	DeleteTestCDNFederations(t)
	DeleteTestUsersDeliveryServices(t)
	DeleteTestDeliveryServices(t)
	DeleteTestCacheGroups(t)
	DeleteTestPhysLocations(t)
	DeleteTestRegions(t)
	DeleteTestDivisions(t)
	DeleteTestStatuses(t)
	DeleteTestProfiles(t)
	DeleteTestTypes(t)
	DeleteTestCDNs(t)
}

func GetTestFederations(t *testing.T) {
	log.Debugln("GetTestFederations")

	if len(testData.Federations) == 0 {
		t.Fatalf("no federations test data")
	}

	feds, _, err := TOSession.AllFederations()
	if err != nil {
		t.Fatalf("getting federations: " + err.Error())
	}

	if len(feds) != 1 {
		t.Fatalf("federations expected 1, actual: %+v", len(feds))
	}
	fed := feds[0]

	if len(fed.Mappings) < 1 {
		t.Fatalf("federation mappings expected <0, actual: 0")
	}

	mapping := fed.Mappings[0]
	if mapping.CName == nil {
		t.Fatalf("federation mapping expected cname, actual: nil")
	}
	if mapping.TTL == nil {
		t.Fatalf("federation mapping expected ttl, actual: nil")
	}

	matched := false
	for _, testFed := range testData.Federations {
		if testFed.CName == nil {
			t.Fatalf("test federation missing cname!")
		}
		if testFed.TTL == nil {
			t.Fatalf("test federation missing ttl!")
		}

		if *mapping.CName != *testFed.CName {
			continue
		}
		matched = true

		if *mapping.TTL != *testFed.TTL {
			t.Fatalf("federation mapping ttl expected: %v, actual: %v", *testFed.TTL, *mapping.TTL)
		}
	}
	if !matched {
		t.Fatalf("federation mapping expected to match test data, actual: cname %v not in test data", *mapping.CName)
	}

	log.Debugln("GetTestFederations PASSED")
}

func PostTestFederationsDeliveryServices(t *testing.T) {
	log.Debugln("PostTestFederationsDeliveryServices")
	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Fatalf("cannot GET DeliveryServices: %v - %v\n", err, dses)
	}
	if len(dses) == 0 {
		t.Fatalf("no delivery services, must have at least 1 ds to test federations deliveryservices\n")
	}
	ds := dses[0]
	if len(fedIDs) == 0 {
		t.Fatalf("no federations, must have at least 1 federation to test federations deliveryservices\n")
	}
	fedID := fedIDs[0]

	if _, err = TOSession.CreateFederationDeliveryServices(fedID, []int{ds.ID}, true); err != nil {
		t.Fatalf("creating federations delivery services: %v\n", err)
	}
	log.Debugln("PostTestFederationsDeliveryServices PASSED")
}
