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

func (r *TCData) CreateTestOrigins(t *testing.T) {
	// loop through origins, assign FKs and create
	for _, origin := range r.TestData.Origins {
		_, _, err := TOSession.CreateOrigin(origin)
		if err != nil {
			t.Errorf("could not CREATE origins: %v", err)
		}
	}
}

func (r *TCData) DeleteTestOrigins(t *testing.T) {
	for _, origin := range r.TestData.Origins {
		resp, _, err := TOSession.GetOriginByName(*origin.Name)
		if err != nil {
			t.Errorf("cannot GET Origin by name: %v - %v", *origin.Name, err)
		}
		if len(resp) > 0 {
			respOrigin := resp[0]

			delResp, _, err := TOSession.DeleteOriginByID(*respOrigin.ID)
			if err != nil {
				t.Errorf("cannot DELETE Origin by ID: %v - %v", err, delResp)
			}

			// Retrieve the Origin to see if it got deleted
			org, _, err := TOSession.GetOriginByName(*origin.Name)
			if err != nil {
				t.Errorf("error deleting Origin name: %s", err.Error())
			}
			if len(org) > 0 {
				t.Errorf("expected Origin name: %s to be deleted", *origin.Name)
			}
		}
	}
}
