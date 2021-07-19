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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestDivisions(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Divisions, Regions}, func() {
		TryToDeleteDivision(t)
		UpdateTestDivisions(t)
		GetTestDivisions(t)
	})
}

func TryToDeleteDivision(t *testing.T) {
	division := testData.Divisions[0]

	resp, _, err := TOSession.GetDivisionByName(division.Name)
	if err != nil {
		t.Errorf("cannot GET Division by name: %v - %v", division.Name, err)
	}
	division = resp[0]
	_, _, err = TOSession.DeleteDivisionByID(division.ID)

	if err == nil {
		t.Error("should not be able to delete a division prematurely")
		return
	}

	if strings.Contains(err.Error(), "Resource not found.") {
		t.Errorf("division with name %v does not exist", division.Name)
		return
	}

	if strings.Contains(err.Error(), "cannot delete division because it is being used by a region") {
		return
	}

	t.Errorf("unexpected error occured: %v", err)
}

func CreateTestDivisions(t *testing.T) {
	for _, division := range testData.Divisions {
		resp, _, err := TOSession.CreateDivision(division)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE division: %v", err)
		}
	}
}

func UpdateTestDivisions(t *testing.T) {

	firstDivision := testData.Divisions[0]
	// Retrieve the Division by division so we can get the id for the Update
	resp, _, err := TOSession.GetDivisionByName(firstDivision.Name)
	if err != nil {
		t.Errorf("cannot GET Division by division: %v - %v", firstDivision.Name, err)
	}
	remoteDivision := resp[0]
	expectedDivision := "division-test"
	remoteDivision.Name = expectedDivision
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateDivisionByID(remoteDivision.ID, remoteDivision)
	if err != nil {
		t.Errorf("cannot UPDATE Division by id: %v - %v", err, alert)
	}

	// Retrieve the Division to check division got updated
	resp, _, err = TOSession.GetDivisionByID(remoteDivision.ID)
	if err != nil {
		t.Errorf("cannot GET Division by division: %v - %v", firstDivision.Name, err)
	}
	respDivision := resp[0]
	if respDivision.Name != expectedDivision {
		t.Errorf("results do not match actual: %s, expected: %s", respDivision.Name, expectedDivision)
	}

	// Set the name back to the fixture value so we can delete it after
	remoteDivision.Name = firstDivision.Name
	alert, _, err = TOSession.UpdateDivisionByID(remoteDivision.ID, remoteDivision)
	if err != nil {
		t.Errorf("cannot UPDATE Division by id: %v - %v", err, alert)
	}

}

func GetTestDivisions(t *testing.T) {
	for _, division := range testData.Divisions {
		resp, _, err := TOSession.GetDivisionByName(division.Name)
		if err != nil {
			t.Errorf("cannot GET Division by division: %v - %v", err, resp)
		}
	}
}

func DeleteTestDivisions(t *testing.T) {

	for _, division := range testData.Divisions {
		// Retrieve the Division by name so we can get the id
		resp, _, err := TOSession.GetDivisionByName(division.Name)
		if err != nil {
			t.Errorf("cannot GET Division by name: %v - %v", division.Name, err)
		}
		if len(resp) != 1 {
			t.Errorf("Expected exactly one Division to exist with name '%s', found: %d", division.Name, len(resp))
			continue
		}
		respDivision := resp[0]

		delResp, _, err := TOSession.DeleteDivisionByID(respDivision.ID)
		if err != nil {
			t.Errorf("cannot DELETE Division by division: %v - %v", err, delResp)
		}

		// Retrieve the Division to see if it got deleted
		divisionResp, _, err := TOSession.GetDivisionByName(division.Name)
		if err != nil {
			t.Errorf("error deleting Division division: %s", err.Error())
		}
		if len(divisionResp) > 0 {
			t.Errorf("expected Division : %s to be deleted", division.Name)
		}
	}
}
