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
	"testing"
)

func (r *TCData) CreateTestPhysLocations(t *testing.T) {
	for _, pl := range r.TestData.PhysLocations {
		resp, _, err := TOSession.CreatePhysLocation(pl)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE physlocations: %v", err)
		}
	}

}

func (r *TCData) DeleteTestPhysLocations(t *testing.T) {

	for _, cdn := range r.TestData.PhysLocations {
		// Retrieve the PhysLocation by name so we can get the id for the Update
		resp, _, err := TOSession.GetPhysLocationByName(cdn.Name)
		if err != nil {
			t.Errorf("cannot GET PhysLocation by name: %s - %v", cdn.Name, err)
		}
		if len(resp) > 0 {
			respPhysLocation := resp[0]

			_, _, err := TOSession.DeletePhysLocationByID(respPhysLocation.ID)
			if err != nil {
				t.Errorf("cannot DELETE PhysLocation by name: '%s' %v", respPhysLocation.Name, err)
			}

			// Retrieve the PhysLocation to see if it got deleted
			cdns, _, err := TOSession.GetPhysLocationByName(cdn.Name)
			if err != nil {
				t.Errorf("error deleting PhysLocation name: %v", err)
			}
			if len(cdns) > 0 {
				t.Errorf("expected PhysLocation name: %s to be deleted", cdn.Name)
			}
		}
	}
}
