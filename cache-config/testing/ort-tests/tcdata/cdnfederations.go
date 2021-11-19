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
	"encoding/json"
	"testing"
)

var fedIDs []int

func (r *TCData) CreateTestCDNFederations(t *testing.T) {

	// Every federation is associated with a cdn
	for i, f := range r.TestData.Federations {

		// CDNs test data and Federations test data are not naturally parallel
		if i >= len(r.TestData.CDNs) {
			break
		}

		data, _, err := TOSession.CreateCDNFederationByName(f, r.TestData.CDNs[i].Name)
		if err != nil {
			t.Errorf("could not POST federations: " + err.Error())
		}
		bytes, _ := json.Marshal(data)
		t.Logf("POST Response: %s", string(bytes))

		// need to save the ids, otherwise the other tests won't be able to reference the federations
		if data.Response.ID == nil {
			t.Error("Federation id is nil after posting")
		} else {
			fedIDs = append(fedIDs, *data.Response.ID)
		}
	}
}

func (r *TCData) DeleteTestCDNFederations(t *testing.T) {

	for _, id := range fedIDs {
		resp, _, err := TOSession.DeleteCDNFederationByID("foo", id)
		if err != nil {
			t.Errorf("cannot DELETE federation by id: '%d' %v", id, err)
		}
		bytes, err := json.Marshal(resp)
		t.Logf("DELETE Response: %s", string(bytes))

		data, _, err := TOSession.GetCDNFederationsByID("foo", id)
		if len(data.Response) != 0 {
			t.Error("expected federation to be deleted")
		}
	}
	fedIDs = nil // reset the global variable for the next test
}
