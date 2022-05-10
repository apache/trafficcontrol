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

var statusNameMap = map[string]bool{
	"ADMIN_DOWN": true,
	"ONLINE":     true,
	"OFFLINE":    true,
	"REPORTED":   true,
	"PRE_PROD":   true,
}

func (r *TCData) CreateTestStatuses(t *testing.T) {

	for _, status := range r.TestData.Statuses {
		resp, _, err := TOSession.CreateStatusNullable(status)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE types: %v", err)
		}
	}

}

func (r *TCData) DeleteTestStatuses(t *testing.T) {

	for _, status := range r.TestData.Statuses {
		if status.Name == nil {
			t.Fatal("cannot get ftest statuses: test data statuses must have names")
		}

		// Retrieve the Status by name so we can get the id for the Update
		resp, _, err := TOSession.GetStatusByName(*status.Name)
		if err != nil {
			t.Errorf("cannot GET Status by name: %s - %v", *status.Name, err)
		}
		respStatus := resp[0]

		delResp, _, err := TOSession.DeleteStatusByID(respStatus.ID)
		if _, ok := statusNameMap[*status.Name]; !ok {
			if err != nil {
				t.Errorf("cannot DELETE Status by name: %v - %v", err, delResp)
			}

			// Retrieve the Status to see if it got deleted
			types, _, err := TOSession.GetStatusByName(*status.Name)
			if err != nil {
				t.Errorf("error deleting Status name: %v", err)
			}
			if len(types) > 0 {
				t.Errorf("expected Status name: %s to be deleted", *status.Name)
			}
		} else {
			if err == nil {
				t.Errorf("expected an error while trying to delete a reserved status, but got nothing")
			}
		}
	}
	r.ForceDeleteStatuses(t)
}

// ForceDeleteStatuses forcibly deletes the statuses from the db.
func (r *TCData) ForceDeleteStatuses(t *testing.T) {

	// NOTE: Special circumstances!  This should *NOT* be done without a really good reason!
	//  Connects directly to the DB to remove statuses rather than going thru the client.
	//  This is required to delte the special statuses.
	db, err := r.OpenConnection()
	if err != nil {
		t.Error("cannot open db")
	}
	defer db.Close()

	q := `DELETE FROM status`
	err = execSQL(db, q)
	if err != nil {
		t.Errorf("cannot execute SQL: %s; SQL is %s", err.Error(), q)
	}
}
