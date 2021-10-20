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

	tc "github.com/apache/trafficcontrol/v6/lib/go-tc"
)

func TestCoordinates(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Coordinates}, func() {
		GetTestCoordinates(t)
		UpdateTestCoordinates(t)
	})
}

func CreateTestCoordinates(t *testing.T) {
	for _, coord := range testData.Coordinates {

		_, _, err := TOSession.CreateCoordinate(coord)
		if err != nil {
			t.Errorf("could not CREATE coordinates: %v", err)
		}
	}
}

func GetTestCoordinates(t *testing.T) {
	for _, coord := range testData.Coordinates {
		resp, _, err := TOSession.GetCoordinateByName(coord.Name)
		if err != nil {
			t.Errorf("cannot GET Coordinate: %v - %v", err, resp)
		}
	}
}

func UpdateTestCoordinates(t *testing.T) {
	firstCoord := testData.Coordinates[0]
	resp, _, err := TOSession.GetCoordinateByName(firstCoord.Name)
	if err != nil {
		t.Errorf("cannot GET Coordinate by name: %v - %v", firstCoord.Name, err)
	}
	coord := resp[0]
	expectedLat := 12.34
	coord.Latitude = expectedLat

	var alert tc.Alerts
	alert, _, err = TOSession.UpdateCoordinateByID(coord.ID, coord)
	if err != nil {
		t.Errorf("cannot UPDATE Coordinate by id: %v - %v", err, alert)
	}

	// Retrieve the Coordinate to check Coordinate name got updated
	resp, _, err = TOSession.GetCoordinateByID(coord.ID)
	if err != nil {
		t.Errorf("cannot GET Coordinate by name: '$%s', %v", firstCoord.Name, err)
	}
	coord = resp[0]
	if coord.Latitude != expectedLat {
		t.Errorf("results do not match actual: %s, expected: %f", coord.Name, expectedLat)
	}
}

func DeleteTestCoordinates(t *testing.T) {
	for _, coord := range testData.Coordinates {
		// Retrieve the Coordinate by name so we can get the id for the Update
		resp, _, err := TOSession.GetCoordinateByName(coord.Name)
		if err != nil {
			t.Errorf("cannot GET Coordinate by name: %v - %v", coord.Name, err)
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
