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

	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

func TestStatuses(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Statuses}, func() {
		UpdateTestStatuses(t)
		GetTestStatuses(t)
	})
}

func CreateTestStatuses(t *testing.T) {
	response, _, err := TOSession.GetStatuses()
	if err != nil {
		t.Errorf("could not get statuses: %v", err)
	}
	statusNameMap := make(map[string]bool, 0)
	for _, r := range response {
		statusNameMap[r.Name] = true
	}
	for _, status := range testData.Statuses {
		if status.Name != nil {
			if _, ok := statusNameMap[*status.Name]; !ok {
				resp, _, err := TOSession.CreateStatusNullable(status)
				t.Log("Response: ", resp)
				if err != nil {
					t.Errorf("could not CREATE status: %v", err)
				}
			}
		}
	}

}

func UpdateTestStatuses(t *testing.T) {

	if len(testData.Statuses) < 1 {
		t.Fatal("Need at least one Status to test updating a Status")
	}
	for _, status := range testData.Statuses {
		if status.Name == nil {
			t.Fatal("cannot update test statuses: test data status must have a name")
		}
		// Retrieve the Status by name so we can get the id for the Update
		resp, _, err := TOSession.GetStatusByName(*status.Name)
		if err != nil {
			t.Errorf("cannot GET Status by name: %s - %v", *status.Name, err)
		}
		remoteStatus := resp[0]
		expectedStatusDesc := "new description"
		remoteStatus.Description = expectedStatusDesc
		var alert tc.Alerts
		alert, _, err = TOSession.UpdateStatusByID(remoteStatus.ID, remoteStatus)

		if tc.IsReservedStatus(*status.Name) {
			if err == nil {
				t.Errorf("expected an error about while updating a reserved status, but got nothing")
			}
		} else {
			if err != nil {
				t.Errorf("cannot UPDATE Status by id: %d, err: %v - %v", remoteStatus.ID, err, alert)
			}

			// Retrieve the Status to check Status name got updated
			resp, _, err = TOSession.GetStatusByID(remoteStatus.ID)
			if err != nil {
				t.Errorf("cannot GET Status by ID: %d - %v", remoteStatus.ID, err)
			}
			respStatus := resp[0]
			if respStatus.Description != expectedStatusDesc {
				t.Errorf("results do not match actual: %s, expected: %s", respStatus.Name, expectedStatusDesc)
			}
		}
	}
}

func GetTestStatuses(t *testing.T) {

	for _, status := range testData.Statuses {
		if status.Name == nil {
			t.Fatal("cannot get ftest statuses: test data statuses must have names")
		}
		resp, _, err := TOSession.GetStatusByName(*status.Name)
		if err != nil {
			t.Errorf("cannot GET Status by name: %v - %v", err, resp)
		}
	}
}

func DeleteTestStatuses(t *testing.T) {

	for _, status := range testData.Statuses {
		if status.Name == nil {
			t.Error("Found status in testing data with null or undefined Name")
			continue
		}

		// Retrieve the Status by name so we can get the id for the Update
		resp, _, err := TOSession.GetStatusByName(*status.Name)
		if err != nil {
			t.Errorf("cannot GET Status by name: %s - %v", *status.Name, err)
		}
		if len(resp) != 1 {
			t.Errorf("Expected exactly one Status to exist with name '%s', found: %d", *status.Name, len(resp))
			continue
		}
		respStatus := resp[0]

		delResp, _, err := TOSession.DeleteStatusByID(respStatus.ID)
		if !tc.IsReservedStatus(*status.Name) {
			if err != nil {
				t.Errorf("cannot DELETE Status by name: %v - %v", err, delResp)
			}

			// Retrieve the Status to see if it got deleted
			types, _, err := TOSession.GetStatusByName(*status.Name)
			if err != nil {
				t.Errorf("error deleting status name: %s, err: %v", *status.Name, err)
			}
			if len(types) > 0 {
				t.Errorf("expected Status name: %s to be deleted", *status.Name)
			}
		} else if err == nil {
			t.Errorf("expected an error while trying to delete a reserved status, but got nothing")
		}
	}
}
