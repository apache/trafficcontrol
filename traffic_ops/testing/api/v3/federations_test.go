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
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestFederations(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices, CDNFederations}, func() {
		PostDeleteTestFederationsDeliveryServices(t)
		GetTestFederations(t)
		GetTestFederationsIMS(t)
		AddFederationResolversForCurrentUserTest(t)
		RemoveFederationResolversForCurrentUserTest(t)
	})
}

func GetTestFederationsIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	if len(testData.Federations) == 0 {
		t.Error("no federations test data")
	}

	_, reqInf, err := TOSession.AllFederationsWithHdr(header)
	if err != nil {
		t.Fatalf("No error expected, but got: %v", err)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestFederations(t *testing.T) {
	if len(testData.Federations) == 0 {
		t.Error("no federations test data")
	}

	feds, _, err := TOSession.AllFederations()
	if err != nil {
		t.Errorf("getting federations: " + err.Error())
	}

	if len(feds) != 2 {
		t.Errorf("federations expected 1, actual: %+v", len(feds))
	}
	fed := feds[0]

	if len(fed.Mappings) < 1 {
		t.Error("federation mappings expected <0, actual: 0")
	}

	mapping := fed.Mappings[0]
	if mapping.CName == nil {
		t.Error("federation mapping expected cname, actual: nil")
	}
	if mapping.TTL == nil {
		t.Error("federation mapping expected ttl, actual: nil")
	}

	matched := false
	for _, testFed := range testData.Federations {
		if testFed.CName == nil {
			t.Error("test federation missing cname!")
		}
		if testFed.TTL == nil {
			t.Error("test federation missing ttl!")
		}

		if *mapping.CName != *testFed.CName {
			continue
		}
		matched = true

		if *mapping.TTL != *testFed.TTL {
			t.Errorf("federation mapping ttl expected: %v, actual: %v", *testFed.TTL, *mapping.TTL)
		}
	}
	if !matched {
		t.Errorf("federation mapping expected to match test data, actual: cname %v not in test data", *mapping.CName)
	}
}

func createFederationToDeliveryServiceAssociation() (int, tc.DeliveryServiceNullable, tc.DeliveryServiceNullable, error) {
	dses, _, err := TOSession.GetDeliveryServicesNullable()
	if err != nil {
		return -1, tc.DeliveryServiceNullable{}, tc.DeliveryServiceNullable{}, fmt.Errorf("cannot GET DeliveryServices: %v - %v", err, dses)
	}
	if len(dses) == 0 {
		return -1, tc.DeliveryServiceNullable{}, tc.DeliveryServiceNullable{}, errors.New("no delivery services, must have at least 1 ds to test federations deliveryservices")
	}
	ds := dses[0]
	ds1 := dses[1]

	if len(fedIDs) == 0 {
		return -1, ds, ds1, errors.New("no federations, must have at least 1 federation to test federations deliveryservices")
	}
	fedID := fedIDs[0]

	_, err = TOSession.CreateFederationDeliveryServices(fedID, []int{*ds.ID, *ds1.ID}, true)
	if err != nil {
		err = fmt.Errorf("creating federations delivery services: %v", err)
	}

	return fedID, ds, ds1, err

}

func PostDeleteTestFederationsDeliveryServices(t *testing.T) {
	fedID, ds, ds1, err := createFederationToDeliveryServiceAssociation()
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Test get created Federation Delivery Services
	fedDSes, _, err := TOSession.GetFederationDeliveryServices(fedID)
	if err != nil {
		t.Fatalf("cannot GET Federation DeliveryServices: %v", err)
	}
	if len(fedDSes) != 2 {
		t.Fatalf("two Federation DeliveryService expected for Federation %v, %v was returned", fedID, len(fedDSes))
	}

	// Delete one of the Delivery Services from the Federation
	_, _, err = TOSession.DeleteFederationDeliveryService(fedID, *ds.ID)
	if err != nil {
		t.Fatalf("cannot Delete Federation %v DeliveryService %v: %v", fedID, ds.ID, err)
	}

	// Make sure it is deleted

	// Test get created Federation Delivery Services
	fedDSes, _, err = TOSession.GetFederationDeliveryServices(fedID)
	if err != nil {
		t.Fatalf("cannot GET Federation DeliveryServices: %v", err)
	}
	if len(fedDSes) != 1 {
		t.Fatalf("one Federation DeliveryService expected for Federation %v, %v was returned", fedID, len(fedDSes))
	}

	// Attempt to delete the last one which should fail as you cannot remove the last
	_, _, err = TOSession.DeleteFederationDeliveryService(fedID, *ds1.ID)
	if err == nil {
		t.Fatal("expected to receive error from attempting to delete last Delivery Service from a Federation")
	}
}

func RemoveFederationResolversForCurrentUserTest(t *testing.T) {
	if len(testData.Federations) < 1 {
		t.Fatal("No test Federations, deleting resolvers cannot be tested!")
	}

	alerts, _, err := TOSession.DeleteFederationResolverMappingsForCurrentUser()
	if err != nil {
		t.Fatalf("Unexpected error deleting Federation Resolvers for current user: %v", err)
	}
	for _, a := range alerts.Alerts {
		if a.Level == tc.SuccessLevel.String() {
			t.Logf("Success message from current user Federation Resolver deletion: %s", a.Text)
		} else {
			t.Errorf("Unexpected %s from deleting Federation Resolvers for current user: %s", a.Level, a.Text)
		}
	}

	// Now try deleting Federation Resolvers when there are none.
	_, _, err = TOSession.DeleteFederationResolverMappingsForCurrentUser()
	if err != nil {
		t.Logf("Received expected error deleting Federation Resolvers for current user: %v", err)
	} else {
		t.Error("Expected an error deleting zero Federation Resolvers, but didn't get one.")
	}
}

func AddFederationResolversForCurrentUserTest(t *testing.T) {
	fedID, ds, ds1, err := createFederationToDeliveryServiceAssociation()
	if err != nil {
		t.Fatalf("%v", err)
	}

	// need to assign myself the federation to set its mappings
	me, _, err := TOSession.GetUserCurrent()
	if err != nil {
		t.Fatalf("Couldn't figure out who I am: %v", err)
	}
	if me.ID == nil {
		t.Fatal("Current user has no ID, cannot continue.")
	}

	_, _, err = TOSession.CreateFederationUsers(fedID, []int{*me.ID}, false)
	if err != nil {
		t.Fatalf("Failed to assign federation to current user: %v", err)
	}

	mappings := tc.DeliveryServiceFederationResolverMappingRequest{
		tc.DeliveryServiceFederationResolverMapping{
			DeliveryService: *ds.XMLID,
			Mappings: tc.ResolverMapping{
				Resolve4: []string{"0.0.0.0"},
				Resolve6: []string{"::1"},
			},
		},
		tc.DeliveryServiceFederationResolverMapping{
			DeliveryService: *ds1.XMLID,
			Mappings: tc.ResolverMapping{
				Resolve4: []string{"1.2.3.4/28"},
				Resolve6: []string{"1234::/110"},
			},
		},
	}

	alerts, _, err := TOSession.AddFederationResolverMappingsForCurrentUser(mappings)
	if err != nil {
		t.Fatalf("Unexpected error adding Federation Resolver mappings for the current user: %v", err)
	}
	for _, a := range alerts.Alerts {
		if a.Level == tc.SuccessLevel.String() {
			t.Logf("Received expected success alert from adding Federation Resolver mappings for the current user: %s", a.Text)
		} else {
			t.Errorf("Unexpected %s from adding Federation Resolver mappings for the current user: %s", a.Level, a.Text)
		}
	}

	mappings = tc.DeliveryServiceFederationResolverMappingRequest{
		tc.DeliveryServiceFederationResolverMapping{
			DeliveryService: "aoeuhtns",
			Mappings: tc.ResolverMapping{
				Resolve4: []string{},
				Resolve6: []string{"dead::beef", "f1d0::f00d/82"},
			},
		},
	}

	alerts, _, err = TOSession.AddFederationResolverMappingsForCurrentUser(mappings)
	if err == nil {
		t.Fatal("Expected error adding Federation Resolver mappings for the current user, but didn't get one")
	}
	for _, a := range alerts.Alerts {
		if a.Level != tc.SuccessLevel.String() {
			t.Logf("Received expected %s from adding Federation Resolver mappings for the current user: %s", a.Level, a.Text)
		} else {
			t.Errorf("Unexpected success from adding Federation Resolver mappings for the current user: %s", a.Text)
		}
	}
}
