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

func (r *TCData) CreateTestCoordinates(t *testing.T) {
	for _, coord := range r.TestData.Coordinates {

		_, _, err := TOSession.CreateCoordinate(coord)
		if err != nil {
			t.Errorf("could not CREATE coordinates: %v", err)
		}
	}
}

func (r *TCData) DeleteTestCoordinates(t *testing.T) {
	for _, coord := range r.TestData.Coordinates {
		// Retrieve the Coordinate by name so we can get the id for the Update
		resp, _, err := TOSession.GetCoordinateByName(coord.Name)
		if err != nil {
			t.Errorf("cannot GET Coordinate by name: %s - %v", coord.Name, err)
		}
		if len(resp) > 0 {
			respCoord := resp[0]
			_, _, err := TOSession.DeleteCoordinateByID(respCoord.ID)
			if err != nil {
				t.Errorf("cannot DELETE Coordinate by name: '%s' %v", respCoord.Name, err)
			}
			// Retrieve the Coordinate to see if it got deleted
			coords, _, err := TOSession.GetCoordinateByName(coord.Name)
			if err != nil {
				t.Errorf("error deleting Coordinate name: %s", err.Error())
			}
			if len(coords) > 0 {
				t.Errorf("expected Coordinate name: %s to be deleted", coord.Name)
			}
		}
	}
}
