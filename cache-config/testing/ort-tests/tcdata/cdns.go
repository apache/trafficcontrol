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

func (r *TCData) CreateTestCDNs(t *testing.T) {

	for _, cdn := range r.TestData.CDNs {
		resp, _, err := TOSession.CreateCDN(cdn)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE cdns: %v", err)
		}
	}

}

func (r *TCData) DeleteTestCDNs(t *testing.T) {

	for _, cdn := range r.TestData.CDNs {
		// Retrieve the CDN by name so we can get the id for the Update
		resp, _, err := TOSession.GetCDNByName(cdn.Name)
		if err != nil {
			t.Errorf("cannot GET CDN by name: %v - %v", cdn.Name, err)
		}
		if len(resp) > 0 {
			respCDN := resp[0]

			_, _, err := TOSession.DeleteCDNByID(respCDN.ID)
			if err != nil {
				t.Errorf("cannot DELETE CDN by name: '%s' %v", respCDN.Name, err)
			}

			// Retrieve the CDN to see if it got deleted
			cdns, _, err := TOSession.GetCDNByName(cdn.Name)
			if err != nil {
				t.Errorf("error deleting CDN name: %s", err.Error())
			}
			if len(cdns) > 0 {
				t.Errorf("expected CDN name: %s to be deleted", cdn.Name)
			}
		}
	}
}
