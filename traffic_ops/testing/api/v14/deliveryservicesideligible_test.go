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
)

func TestDeliveryServicesEligible(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		GetTestDeliveryServicesEligible(t)
	})
}

func GetTestDeliveryServicesEligible(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v\n", err)
	}
	if len(dses) == 0 {
		t.Errorf("GET DeliveryServices returned no delivery services, need at least 1 to test")
	}
	dsID := dses[0].ID
	servers, _, err := TOSession.GetDeliveryServicesEligible(dsID)
	if err != nil {
		t.Errorf("getting delivery services eligible: %v\n", err)
	}
	if len(servers) == 0 {
		t.Errorf("getting delivery services eligible returned no servers\n")
	}
}

/*func TestDeliveryServicesNotEligible(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices, ServerCapabilities, ServerServerCapabilities, DeliveryServicesRequiredCapabilities}, func() {
		GetTestDeliveryServicesNotEligible(t)
	})
}

func GetTestDeliveryServicesNotEligible(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v\n", err)
	}
	if len(dses) == 0 {
		t.Errorf("GET DeliveryServices returned no delivery services, need at least 1 to test")
	}
	dsID := dses[0].ID
	rc := testData.ServerCapabilities[0]

	cap := tc.DeliveryServicesRequiredCapability{
		DeliveryServiceID:  &dsID,
		RequiredCapability: &rc.Name,
	}

	_, _, err = TOSession.CreateDeliveryServicesRequiredCapability(cap)
	if err != nil {
		t.Errorf("error creating required capability: %s", err.Error())
	}

	s := testData.Servers[0]
	fmt.Println("***", s.ID, "***")
	ssc := tc.ServerServerCapability{
		ServerID:         &s.ID,
		ServerCapability: &rc.Name,
	}

	_, _, err = TOSession.CreateServerServerCapability(ssc)
	if err != nil {
		t.Errorf("error creating server server capability: %s", err.Error())
	}
	fmt.Println("*** AM I HERE? ***")

	servers, _, err := TOSession.GetDeliveryServicesEligible(dsID)
	if err != nil {
		t.Errorf("getting delivery services eligible: %v\n", err)
	}
	if len(servers) == 0 {
		t.Errorf("getting delivery services eligible returned no servers\n")
	}

	_, _, err = TOSession.DeleteDeliveryServicesRequiredCapability(dsID, rc.Name)
	if err != nil {
		t.Errorf("error deleting ds  required capability: %s", err.Error())
	}
	_, _, err = TOSession.DeleteServerServerCapability(s.ID, rc.Name)
	if err != nil {
		t.Errorf("error deleting ds  required capability: %s", err.Error())
	}
}*/
