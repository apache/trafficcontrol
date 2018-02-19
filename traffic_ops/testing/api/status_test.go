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

package api

import (
	"fmt"
	"testing"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

func TestStatuses(t *testing.T) {

	CreateTestStatuses(t)
	UpdateTestStatuses(t)
	GetTestStatuses(t)
	DeleteTestStatuses(t)

}

func CreateTestStatuses(t *testing.T) {
	for _, status := range testData.Statuses {
		fmt.Printf("status ---> %v\n", status)
		resp, _, err := TOSession.CreateStatus(status)
		log.Debugln("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE statuss: %v\n", err)
		}
	}

}

func UpdateTestStatuses(t *testing.T) {

	firstStatus := testData.Statuses[0]
	// Retrieve the Status by status so we can get the id for the Update
	resp, _, err := TOSession.GetStatusByName(firstStatus.Name)
	if err != nil {
		t.Errorf("cannot GET Status by status: %v - %v\n", firstStatus.Name, err)
	}
	remoteStatus := resp[0]
	expectedStatus := "OFFLINE-TEST"
	remoteStatus.Name = expectedStatus
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateStatusByID(remoteStatus.ID, remoteStatus)
	if err != nil {
		t.Errorf("cannot UPDATE Status by id: %v - %v\n", err, alert)
	}

	// Retrieve the Status to check status got updated
	resp, _, err = TOSession.GetStatusByID(remoteStatus.ID)
	if err != nil {
		t.Errorf("cannot GET Status by status: %v - %v\n", firstStatus.Name, err)
	}
	respStatus := resp[0]
	if respStatus.Name != expectedStatus {
		t.Errorf("results do not match actual: %s, expected: %s\n", respStatus.Name, expectedStatus)
	}

}

func GetTestStatuses(t *testing.T) {
	for _, status := range testData.Statuses {
		resp, _, err := TOSession.GetStatusByName(status.Name)
		if err != nil {
			t.Errorf("cannot GET Status by status: %v - %v\n", err, resp)
		}
	}
}

func DeleteTestStatuses(t *testing.T) {

	status := testData.Statuses[1]
	// Retrieve the Status by name so we can get the id
	resp, _, err := TOSession.GetStatusByName(status.Name)
	if err != nil {
		t.Errorf("cannot GET Status by name: %v - %v\n", status.Name, err)
	}
	respStatus := resp[0]

	delResp, _, err := TOSession.DeleteStatusByID(respStatus.ID)
	if err != nil {
		t.Errorf("cannot DELETE Status by status: %v - %v\n", err, delResp)
	}

	// Retrieve the Status to see if it got deleted
	statusResp, _, err := TOSession.GetStatusByName(status.Name)
	if err != nil {
		t.Errorf("error deleting Status status: %s\n", err.Error())
	}
	if len(statusResp) > 0 {
		t.Errorf("expected Status status: %s to be deleted\n", status.Name)
	}
}
