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

package integration

import (
	"encoding/json"
	"fmt"
	"testing"

	traffic_ops "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

func TestProfiles(t *testing.T) {

	uri := fmt.Sprintf("/api/1.2/profiles.json")
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var apiProfileRes traffic_ops.ProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiProfileRes); err != nil {
		t.Errorf("Could not decode profile json.  Error is: %v\n", err)
		t.FailNow()
	}
	apiProfiles := apiProfileRes.Response

	clientProfiles, err := to.Profiles()
	if err != nil {
		t.Errorf("Could not get profiles from client.  Error is: %v\n", err)
		t.FailNow()
	}

	if len(apiProfiles) != len(clientProfiles) {
		t.Errorf("Profile Response Length -- expected %v, got %v\n", len(apiProfiles), len(clientProfiles))
	}

	for _, apiProfile := range apiProfiles {
		match := false
		for _, clientProfile := range clientProfiles {
			if apiProfile.ID == clientProfile.ID {
				match = true
				if apiProfile.Description != clientProfile.Description {
					t.Errorf("Description -- Expected %v, got %v\n", apiProfile.Description, clientProfile.Description)
				}
				if apiProfile.LastUpdated != clientProfile.LastUpdated {
					t.Errorf("Last Updated -- Expected %v, got %v\n", apiProfile.LastUpdated, clientProfile.LastUpdated)
				}
				if apiProfile.Name != clientProfile.Name {
					t.Errorf("Name -- Expected %v, got %v\n", apiProfile.Name, clientProfile.Name)
				}
			}
		}
		if !match {
			t.Errorf("Did not get a profile matching %v\n", apiProfile.Name)
		}
	}
}
