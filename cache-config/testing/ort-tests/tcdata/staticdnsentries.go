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

func (r *TCData) CreateTestStaticDNSEntries(t *testing.T) {
	for _, staticDNSEntry := range r.TestData.StaticDNSEntries {
		resp, _, err := TOSession.CreateStaticDNSEntry(staticDNSEntry)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE staticDNSEntry: %v", err)
		}
	}

}

func (r *TCData) DeleteTestStaticDNSEntries(t *testing.T) {

	for _, staticDNSEntry := range r.TestData.StaticDNSEntries {
		// Retrieve the StaticDNSEntries by name so we can get the id for the Update
		resp, _, err := TOSession.GetStaticDNSEntriesByHost(staticDNSEntry.Host)
		if err != nil {
			t.Errorf("cannot GET StaticDNSEntries by name: %s - %v", staticDNSEntry.Host, err)
		}
		if len(resp) > 0 {
			respStaticDNSEntry := resp[0]

			_, _, err := TOSession.DeleteStaticDNSEntryByID(respStaticDNSEntry.ID)
			if err != nil {
				t.Errorf("cannot DELETE StaticDNSEntry by name: '%s' %v", respStaticDNSEntry.Host, err)
			}

			// Retrieve the StaticDNSEntry to see if it got deleted
			staticDNSEntries, _, err := TOSession.GetStaticDNSEntriesByHost(staticDNSEntry.Host)
			if err != nil {
				t.Errorf("error deleting StaticDNSEntrie name: %v", err)
			}
			if len(staticDNSEntries) > 0 {
				t.Errorf("expected StaticDNSEntry name: %s to be deleted", staticDNSEntry.Host)
			}
		}
	}
}
