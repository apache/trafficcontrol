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

func (r *TCData) CreateTestDivisions(t *testing.T) {
	for _, division := range r.TestData.Divisions {
		resp, _, err := TOSession.CreateDivision(division)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE division: %v", err)
		}
	}
}

func (r *TCData) DeleteTestDivisions(t *testing.T) {

	for _, division := range r.TestData.Divisions {
		// Retrieve the Division by name so we can get the id
		resp, _, err := TOSession.GetDivisionByName(division.Name)
		if err != nil {
			t.Errorf("cannot GET Division by name: %v - %v", division.Name, err)
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
