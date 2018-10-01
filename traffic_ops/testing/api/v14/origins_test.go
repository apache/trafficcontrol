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

func TestOrigins(t *testing.T) {
	CreateTestCDNs(t)
	defer DeleteTestCDNs(t)
	CreateTestTypes(t)
	defer DeleteTestTypes(t)
	CreateTestProfiles(t)
	defer DeleteTestProfiles(t)
	CreateTestStatuses(t)
	defer DeleteTestStatuses(t)
	CreateTestDivisions(t)
	defer DeleteTestDivisions(t)
	CreateTestRegions(t)
	defer DeleteTestRegions(t)
	CreateTestPhysLocations(t)
	defer DeleteTestPhysLocations(t)
	CreateTestCacheGroups(t)
	defer DeleteTestCacheGroups(t)
	CreateTestServers(t)
	defer DeleteTestServers(t)
	CreateTestDeliveryServices(t)
	defer DeleteTestDeliveryServices(t)
	CreateTestCoordinates(t)
	defer DeleteTestCoordinates(t)
	// TODO: add tenants once their API integration tests are implemented

	CreateTestOrigins(t)
	defer DeleteTestOrigins(t)
	UpdateTestOrigins(t)
	GetTestOrigins(t)
}

func CreateTestOrigins(t *testing.T) {
	failed := false

	// GET ORIGIN1 profile
	respProfiles, _, err := TOSession.GetProfileByName("ORIGIN1")
	if err != nil {
		t.Errorf("cannot GET Profiles - %v\n", err)
		failed = true
	}
	respProfile := respProfiles[0]

	// GET originCachegroup cachegroup
	respCacheGroups, _, err := TOSession.GetCacheGroupNullableByName("originCachegroup")
	if err != nil {
		t.Errorf("cannot GET CacheGroup by name: originCachegroup - %v\n", err)
		failed = true
	}
	respCacheGroup := respCacheGroups[0]

	// GET deliveryservices
	respDeliveryServices, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Errorf("cannot GET Delivery Services - %v\n", err)
		failed = true
	}
	if len(respDeliveryServices) == 0 {
		t.Errorf("no delivery services found")
		failed = true
	}

	// GET coordinate1 coordinate
	respCoordinates, _, err := TOSession.GetCoordinateByName("coordinate1")
	if err != nil {
		t.Errorf("cannot GET Coordinate by name: coordinate1 - %v\n", err)
		failed = true
	}
	respCoordinate := respCoordinates[0]

	// loop through origins, assign FKs and create
	for _, origin := range testData.Origins {
		origin.CachegroupID = respCacheGroup.ID
		origin.CoordinateID = &respCoordinate.ID
		origin.ProfileID = &respProfile.ID
		origin.DeliveryServiceID = &respDeliveryServices[0].ID

		_, _, err = TOSession.CreateOrigin(origin)
		if err != nil {
			t.Errorf("could not CREATE origins: %v\n", err)
			failed = true
		}
	}

	if !failed {
		log.Debugln("CreateTestOrigins() PASSED")
	}

}

func GetTestOrigins(t *testing.T) {
	failed := false

	for _, origin := range testData.Origins {
		resp, _, err := TOSession.GetServerByHostName(*origin.Name)
		if err != nil {
			t.Errorf("cannot GET Origin by name: %v - %v\n", err, resp)
			failed = true
		}
	}

	if !failed {
		log.Debugln("GetTestOrigins() PASSED")
	}
}

func UpdateTestOrigins(t *testing.T) {
	failed := false

	firstOrigin := testData.Origins[0]
	// Retrieve the origin by name so we can get the id for the Update
	resp, _, err := TOSession.GetOriginByName(*firstOrigin.Name)
	if err != nil {
		t.Errorf("cannot GET origin by name: %v - %v\n", *firstOrigin.Name, err)
		failed = true
	}
	remoteOrigin := resp[0]
	updatedPort := 4321
	updatedFQDN := "updated.example.com"

	// update port and FQDN values on origin
	remoteOrigin.Port = &updatedPort
	remoteOrigin.FQDN = &updatedFQDN
	updResp, _, err := TOSession.UpdateOriginByID(*remoteOrigin.ID, remoteOrigin)
	if err != nil {
		t.Errorf("cannot UPDATE Origin by name: %v - %v\n", err, updResp.Alerts)
		failed = true
	}

	// Retrieve the origin to check port and FQDN values were updated
	resp, _, err = TOSession.GetOriginByID(*remoteOrigin.ID)
	if err != nil {
		t.Errorf("cannot GET Origin by ID: %v - %v\n", *remoteOrigin.Name, err)
		failed = true
	}

	respOrigin := resp[0]
	if *respOrigin.Port != updatedPort || *respOrigin.FQDN != updatedFQDN {
		t.Errorf("results do not match actual: %d, expected: %d\n", *respOrigin.Port, updatedPort)
		t.Errorf("results do not match actual: %s, expected: %s\n", *respOrigin.FQDN, updatedFQDN)
		failed = true
	}

	if !failed {
		log.Debugln("UpdateTestOrigins() PASSED")
	}
}

func DeleteTestOrigins(t *testing.T) {
	failed := false

	for _, origin := range testData.Origins {
		resp, _, err := TOSession.GetOriginByName(*origin.Name)
		if err != nil {
			t.Errorf("cannot GET Origin by name: %v - %v\n", *origin.Name, err)
			failed = true
		}
		if len(resp) > 0 {
			respOrigin := resp[0]

			delResp, _, err := TOSession.DeleteOriginByID(*respOrigin.ID)
			if err != nil {
				t.Errorf("cannot DELETE Origin by ID: %v - %v\n", err, delResp)
				failed = true
			}

			// Retrieve the Origin to see if it got deleted
			org, _, err := TOSession.GetOriginByName(*origin.Name)
			if err != nil {
				t.Errorf("error deleting Origin name: %s\n", err.Error())
				failed = true
			}
			if len(org) > 0 {
				t.Errorf("expected Origin name: %s to be deleted\n", *origin.Name)
				failed = true
			}
		}
	}

	if !failed {
		log.Debugln("DeleteTestOrigins() PASSED")
	}
}
