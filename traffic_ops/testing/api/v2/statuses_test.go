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

	for _, status := range testData.Statuses {
		resp, _, err := TOSession.CreateStatusNullable(status)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE types: %v", err)
		}
	}

}

func UpdateTestStatuses(t *testing.T) {

	firstStatus := testData.Statuses[0]
	if firstStatus.Name == nil {
		t.Fatal("cannot update test statuses: first test data status must have a name")
	}

	// Retrieve the Status by name so we can get the id for the Update
	resp, _, err := TOSession.GetStatusByName(*firstStatus.Name)
	if err != nil {
		t.Errorf("cannot GET Status by name: %v - %v", firstStatus.Name, err)
	}
	remoteStatus := resp[0]
	expectedStatusDesc := "new description"
	remoteStatus.Description = expectedStatusDesc
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateStatusByID(remoteStatus.ID, remoteStatus)
	if err != nil {
		t.Errorf("cannot UPDATE Status by id: %v - %v", err, alert)
	}

	// Retrieve the Status to check Status name got updated
	resp, _, err = TOSession.GetStatusByID(remoteStatus.ID)
	if err != nil {
		t.Errorf("cannot GET Status by ID: %v - %v", firstStatus.Description, err)
	}
	respStatus := resp[0]
	if respStatus.Description != expectedStatusDesc {
		t.Errorf("results do not match actual: %s, expected: %s", respStatus.Name, expectedStatusDesc)
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
			t.Fatal("cannot get ftest statuses: test data statuses must have names")
		}

		// Retrieve the Status by name so we can get the id for the Update
		resp, _, err := TOSession.GetStatusByName(*status.Name)
		if err != nil {
			t.Errorf("cannot GET Status by name: %v - %v", status.Name, err)
		}
		respStatus := resp[0]

		delResp, _, err := TOSession.DeleteStatusByID(respStatus.ID)
		if err != nil {
			t.Errorf("cannot DELETE Status by name: %v - %v", err, delResp)
		}

		// Retrieve the Status to see if it got deleted
		types, _, err := TOSession.GetStatusByName(*status.Name)
		if err != nil {
			t.Errorf("error deleting Status name: %s", err.Error())
		}
		if len(types) > 0 {
			t.Errorf("expected Status name: %s to be deleted", *status.Name)
		}
	}
}
