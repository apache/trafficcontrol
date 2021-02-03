package tcdata

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
	"strings"
	"testing"
)

func (r *TCData) CreateTestRegions(t *testing.T) {

	for _, region := range r.TestData.Regions {
		resp, _, err := TOSession.CreateRegion(region)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE region: %v", err)
		}
	}
}

func (r *TCData) DeleteTestRegionsByName(t *testing.T) {
	for _, region := range r.TestData.Regions {
		delResp, _, err := TOSession.DeleteRegion(nil, &region.Name)
		if err != nil {
			t.Errorf("cannot DELETE Region by name: %v - %v", err, delResp)
		}

		deleteLog, _, err := TOSession.GetLogsByLimit(1)
		if err != nil {
			t.Fatalf("unable to get latest audit log entry")
		}
		if len(deleteLog) != 1 {
			t.Fatalf("log entry length - expected: 1, actual: %d", len(deleteLog))
		}
		if !strings.Contains(*deleteLog[0].Message, region.Name) {
			t.Errorf("region deletion audit log entry - expected: message containing region name '%s', actual: %s", region.Name, *deleteLog[0].Message)
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

	r.CreateTestRegions(t)
}

func (r *TCData) DeleteTestRegions(t *testing.T) {

	for _, region := range r.TestData.Regions {
		// Retrieve the Region by name so we can get the id
		resp, _, err := TOSession.GetRegionByName(region.Name)
		if err != nil {
			t.Errorf("cannot GET Region by name: %v - %v", region.Name, err)
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
