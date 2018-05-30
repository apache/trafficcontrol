package v13

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
	"fmt"
	"testing"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

func TestStaticDNSEntries(t *testing.T) {

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
	CreateTestStaticDNSEntries(t)
	GetTestStaticDNSEntries(t)
	UpdateTestStaticDNSEntries(t)
	DeleteTestStaticDNSEntries(t)
	DeleteTestTypes(t)
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

func CreateTestStaticDNSEntries(t *testing.T) {

	for _, staticDNSEntry := range testData.StaticDNSEntries {

		// GET StaticDNSEntry type
		respTypes, _, err := TOSession.GetTypeByName(staticDNSEntry.Type)
		if err != nil {
			t.Errorf("cannot GET Type by name: %v\n", err)
		}
		respType := respTypes[0]
		staticDNSEntry.TypeID = respType.ID

		// GET DeliveryService to associate
		respDSes, _, err := TOSession.GetDeliveryServiceByXMLID(staticDNSEntry.DeliveryService)
		if err != nil {
			t.Errorf("cannot GET DeliveryService by XMLID: %v\n", err)
		}
		respDS := respDSes[0]
		staticDNSEntry.DeliveryServiceID = respDS.ID

		// GET Cachegroup to associate
		respGroups, _, err := TOSession.GetCacheGroupByName(staticDNSEntry.CacheGroupName)
		if err != nil {
			t.Errorf("cannot GET CacheGroup by Name: %v\n", err)
		}
		respGroup := respGroups[0]
		staticDNSEntry.CacheGroupID = respGroup.ID

		resp, _, err := TOSession.CreateStaticDNSEntry(staticDNSEntry)
		log.Debugln("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE staticDNSEntry: %v\n", err)
		}
	}

}

func UpdateTestStaticDNSEntries(t *testing.T) {

	firstStaticDNSEntry := testData.StaticDNSEntries[0]
	// Retrieve the StaticDNSEntries by name so we can get the id for the Update
	resp, _, err := TOSession.GetStaticDNSEntriesByHost(firstStaticDNSEntry.Host)
	if err != nil {
		t.Errorf("cannot GET StaticDNSEntries by name: '%s', %v\n", firstStaticDNSEntry.Host, err)
	}
	remoteStaticDNSEntry := resp[0]
	fmt.Printf("remoteStaticDNSEntry ---> %v\n", remoteStaticDNSEntry)
	fmt.Printf("remoteStaticDNSEntry.ID ---> %v\n", remoteStaticDNSEntry.ID)
	expectedAddress := "address99"
	remoteStaticDNSEntry.Address = expectedAddress
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateStaticDNSEntryByID(remoteStaticDNSEntry.ID, remoteStaticDNSEntry)
	if err != nil {
		t.Errorf("cannot UPDATE StaticDNSEntries by id: %v - %v\n", err, alert)
	}

	// Retrieve the StaticDNSEntries to check StaticDNSEntries name got updated
	resp, _, err = TOSession.GetStaticDNSEntryByID(remoteStaticDNSEntry.ID)
	if err != nil {
		t.Errorf("cannot GET StaticDNSEntries by name: '$%s', %v\n", firstStaticDNSEntry.Host, err)
	}
	respStaticDNSEntry := resp[0]
	if respStaticDNSEntry.Address != expectedAddress {
		t.Errorf("results do not match actual: %s, expected: %s\n", respStaticDNSEntry.Address, expectedAddress)
	}

}

func GetTestStaticDNSEntries(t *testing.T) {

	for _, staticDNSEntry := range testData.StaticDNSEntries {
		resp, _, err := TOSession.GetStaticDNSEntriesByHost(staticDNSEntry.Host)
		if err != nil {
			t.Errorf("cannot GET StaticDNSEntries by name: %v - %v\n", err, resp)
		}
	}
}

func DeleteTestStaticDNSEntries(t *testing.T) {

	for _, staticDNSEntry := range testData.StaticDNSEntries {
		// Retrieve the StaticDNSEntries by name so we can get the id for the Update
		resp, _, err := TOSession.GetStaticDNSEntriesByHost(staticDNSEntry.Host)
		if err != nil {
			t.Errorf("cannot GET StaticDNSEntries by name: %v - %v\n", staticDNSEntry.Host, err)
		}
		if len(resp) > 0 {
			respStaticDNSEntry := resp[0]

			_, _, err := TOSession.DeleteStaticDNSEntryByID(respStaticDNSEntry.ID)
			if err != nil {
				t.Errorf("cannot DELETE StaticDNSEntrie by name: '%s' %v\n", respStaticDNSEntry.Host, err)
			}

			// Retrieve the StaticDNSEntrie to see if it got deleted
			staticDNSEntries, _, err := TOSession.GetStaticDNSEntriesByHost(staticDNSEntry.Host)
			if err != nil {
				t.Errorf("error deleting StaticDNSEntrie name: %s\n", err.Error())
			}
			if len(staticDNSEntries) > 0 {
				t.Errorf("expected StaticDNSEntry name: %s to be deleted\n", staticDNSEntry.Host)
			}
		}
	}
}
