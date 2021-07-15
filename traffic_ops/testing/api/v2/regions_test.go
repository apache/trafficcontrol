package v2

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

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestRegions(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Divisions, Regions}, func() {
		UpdateTestRegions(t)
		GetTestRegions(t)
	})
}

func CreateTestRegions(t *testing.T) {

	for _, region := range testData.Regions {
		resp, _, err := TOSession.CreateRegion(region)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE region: %v", err)
		}
	}
}

func UpdateTestRegions(t *testing.T) {

	firstRegion := testData.Regions[0]
	// Retrieve the Region by region so we can get the id for the Update
	resp, _, err := TOSession.GetRegionByName(firstRegion.Name)
	if err != nil {
		t.Errorf("cannot GET Region by region: %v - %v", firstRegion.Name, err)
	}
	remoteRegion := resp[0]
	expectedRegion := "OFFLINE-TEST"
	remoteRegion.Name = expectedRegion
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateRegionByID(remoteRegion.ID, remoteRegion)
	if err != nil {
		t.Errorf("cannot UPDATE Region by id: %v - %v", err, alert)
	}

	// Retrieve the Region to check region got updated
	resp, _, err = TOSession.GetRegionByID(remoteRegion.ID)
	if err != nil {
		t.Errorf("cannot GET Region by region: %v - %v", firstRegion.Name, err)
	}
	respRegion := resp[0]
	if respRegion.Name != expectedRegion {
		t.Errorf("results do not match actual: %s, expected: %s", respRegion.Name, expectedRegion)
	}

	// Set the name back to the fixture value so we can delete it after
	remoteRegion.Name = firstRegion.Name
	alert, _, err = TOSession.UpdateRegionByID(remoteRegion.ID, remoteRegion)
	if err != nil {
		t.Errorf("cannot UPDATE Region by id: %v - %v", err, alert)
	}

}

func GetTestRegions(t *testing.T) {
	for _, region := range testData.Regions {
		resp, _, err := TOSession.GetRegionByName(region.Name)
		if err != nil {
			t.Errorf("cannot GET Region by region: %v - %v", err, resp)
		}
	}
}

func DeleteTestRegions(t *testing.T) {

	for _, region := range testData.Regions {
		// Retrieve the Region by name so we can get the id
		resp, _, err := TOSession.GetRegionByName(region.Name)
		if err != nil {
			t.Errorf("cannot GET Region by name: %v - %v", region.Name, err)
		}
		if len(resp) != 1 {
			t.Errorf("Expected exactly one Region to exist with name '%s', found: %d", region.Name, len(resp))
			continue
		}
		respRegion := resp[0]

		delResp, _, err := TOSession.DeleteRegionByID(respRegion.ID)
		if err != nil {
			t.Errorf("cannot DELETE Region by region: %v - %v", err, delResp)
		}

		// Retrieve the Region to see if it got deleted
		regionResp, _, err := TOSession.GetRegionByName(region.Name)
		if err != nil {
			t.Errorf("error deleting Region region: %s", err.Error())
		}
		if len(regionResp) > 0 {
			t.Errorf("expected Region : %s to be deleted", region.Name)
		}
	}
}
