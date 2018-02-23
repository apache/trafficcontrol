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

package v13

import (
	"testing"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

func TestPhysLocations(t *testing.T) {

	CreateTestPhysLocations(t)
	UpdateTestPhysLocations(t)
	GetTestPhysLocations(t)
	DeleteTestPhysLocations(t)

}

func CreateTestPhysLocations(t *testing.T) {

	for _, pl := range testData.PhysLocations {
		resp, _, err := TOSession.CreatePhysLocation(pl)
		log.Debugln("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE phys_locations: %v\n", err)
		}
	}

}

func UpdateTestPhysLocations(t *testing.T) {

	firstPhysLocation := testData.PhysLocations[0]
	// Retrieve the PhysLocation by name so we can get the id for the Update
	resp, _, err := TOSession.GetPhysLocationByName(firstPhysLocation.Name)
	if err != nil {
		t.Errorf("cannot GET PhysLocation by name: %v - %v\n", firstPhysLocation.Name, err)
	}
	remotePhysLocation := resp[0]
	expectedPhysLocationName := "testPhysLoc1"
	remotePhysLocation.Name = expectedPhysLocationName
	var alert tc.Alerts
	alert, _, err = TOSession.UpdatePhysLocationByID(remotePhysLocation.ID, remotePhysLocation)
	if err != nil {
		t.Errorf("cannot UPDATE PhysLocation by id: %v - %v\n", err, alert)
	}

	// Retrieve the PhysLocation to check PhysLocation name got updated
	resp, _, err = TOSession.GetPhysLocationByID(remotePhysLocation.ID)
	if err != nil {
		t.Errorf("cannot GET PhysLocation by name: %v - %v\n", firstPhysLocation.Name, err)
	}
	respPhysLocation := resp[0]
	if respPhysLocation.Name != expectedPhysLocationName {
		t.Errorf("results do not match actual: %s, expected: %s\n", respPhysLocation.Name, expectedPhysLocationName)
	}

}

func GetTestPhysLocations(t *testing.T) {

	for _, pl := range testData.PhysLocations {
		resp, _, err := TOSession.GetPhysLocationByName(pl.Name)
		if err != nil {
			t.Errorf("cannot GET PhysLocation by name: %v - %v\n", err, resp)
		}
	}
}

func DeleteTestPhysLocations(t *testing.T) {

	pl := testData.PhysLocations[1]
	// Retrieve the PhysLocation by name so we can get the id for the Update
	resp, _, err := TOSession.GetPhysLocationByName(pl.Name)
	if err != nil {
		t.Errorf("cannot GET PhysLocation by name: %v - %v\n", pl.Name, err)
	}
	respPhysLocation := resp[0]

	delResp, _, err := TOSession.DeletePhysLocationByID(respPhysLocation.ID)
	if err != nil {
		t.Errorf("cannot DELETE PhysLocation by name: %v - %v\n", err, delResp)
	}

	// Retrieve the PhysLocation to see if it got deleted
	pls, _, err := TOSession.GetPhysLocationByName(pl.Name)
	if err != nil {
		t.Errorf("error deleting PhysLocation name: %s\n", err.Error())
	}
	if len(pls) > 0 {
		t.Errorf("expected PhysLocation name: %s to be deleted\n", pl.Name)
	}
}
