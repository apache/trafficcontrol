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

func (r *TCData) CreateTestStatuses(t *testing.T) {
	for _, status := range r.TestData.Statuses {
		resp, _, err := TOSession.CreateStatusNullable(status)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE status: %v", err)
		}
	}
}

func (r *TCData) DeleteTestStatuses(t *testing.T) {

	for _, status := range r.TestData.Statuses {
		if status.Name == nil {
			t.Fatal("cannot get test statuses: test data statuses must have names")
		}

		// Retrieve the Status by name so we can get the id for the Update
		resp, _, err := TOSession.GetStatusByName(*status.Name)
		if err != nil {
			t.Errorf("cannot GET Status by name: %s - %v", *status.Name, err)
		}
		respStatus := resp[0]

		delResp, _, err := TOSession.DeleteStatusByID(respStatus.ID)
		if err != nil {
			t.Errorf("cannot DELETE Status by ID: %v - %v", err, delResp)
		}

		// Retrieve the Status to see if it got deleted
		statuses, _, err := TOSession.GetStatusByName(*status.Name)
		if err != nil {
			t.Errorf("error getting Status name: %v", err)
		}
		if len(statuses) > 0 {
			t.Errorf("expected Status name: %s to be deleted", *status.Name)
		}
	}
}
