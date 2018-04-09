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
	"testing"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

func TestLocations(t *testing.T) {
	CreateTestLocations(t)
	GetTestLocations(t)
	UpdateTestLocations(t)
	DeleteTestLocations(t)
}

func CreateTestLocations(t *testing.T) {
	failed := false

	for _, loc := range testData.Locations {

		_, _, err := TOSession.CreateLocation(loc)
		if err != nil {
			t.Errorf("could not CREATE locations: %v\n", err)
			failed = true
		}
	}
	if !failed {
		log.Debugln("CreateTestLocations() PASSED: ")
	}
}

func GetTestLocations(t *testing.T) {
	failed := false
	for _, loc := range testData.Locations {
		resp, _, err := TOSession.GetLocationByName(loc.Name)
		if err != nil {
			t.Errorf("cannot GET Location: %v - %v\n", err, resp)
			failed = true
		}
	}
	if !failed {
		log.Debugln("GetTestLocations() PASSED: ")
	}
}

func UpdateTestLocations(t *testing.T) {
	failed := false
	firstLoc := testData.Locations[0]
	resp, _, err := TOSession.GetLocationByName(firstLoc.Name)
	if err != nil {
		t.Errorf("cannot GET Location by name: %v - %v\n", firstLoc.Name, err)
		failed = true
	}
	loc := resp[0]
	expectedName := "blah"
	loc.Name = expectedName

	var alert tc.Alerts
	alert, _, err = TOSession.UpdateLocationByID(loc.ID, loc)
	if err != nil {
		t.Errorf("cannot UPDATE Location by id: %v - %v\n", err, alert)
		failed = true
	}

	// Retrieve the Location to check Location name got updated
	resp, _, err = TOSession.GetLocationByID(loc.ID)
	if err != nil {
		t.Errorf("cannot GET Location by name: '$%s', %v\n", firstLoc.Name, err)
		failed = true
	}
	loc = resp[0]
	if loc.Name != expectedName {
		t.Errorf("results do not match actual: %s, expected: %s\n", loc.Name, expectedName)
	}
	if !failed {
		log.Debugln("UpdateTestLocations() PASSED: ")
	}
}

func DeleteTestLocations(t *testing.T) {
	failed := false

	for _, loc := range testData.Locations {
		// Retrieve the Location by name so we can get the id for the Update
		resp, _, err := TOSession.GetLocationByName(loc.Name)
		if err != nil {
			t.Errorf("cannot GET Location by name: %v - %v\n", loc.Name, err)
			failed = true
		}
		if len(resp) > 0 {
			respLoc := resp[0]
			_, _, err := TOSession.DeleteLocationByID(respLoc.ID)
			if err != nil {
				t.Errorf("cannot DELETE Location by name: '%s' %v\n", respLoc.Name, err)
				failed = true
			}
			// Retrieve the Location to see if it got deleted
			locs, _, err := TOSession.GetLocationByName(loc.Name)
			if err != nil {
				t.Errorf("error deleting Location name: %s\n", err.Error())
				failed = true
			}
			if len(locs) > 0 {
				t.Errorf("expected Location name: %s to be deleted\n", loc.Name)
				failed = true
			}
		}
	}

	if !failed {
		log.Debugln("DeleteTestLocations() PASSED: ")
	}
}
