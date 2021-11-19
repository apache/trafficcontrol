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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

// CreateBadProfiles ensures that profiles can't be created with bad values
func (r *TCData) CreateBadProfiles(t *testing.T) {

	// blank profile
	prs := []tc.Profile{
		tc.Profile{Type: "", Name: "", Description: "", CDNID: 0},
		tc.Profile{Type: "ATS_PROFILE", Name: "badprofile", Description: "description", CDNID: 0},
		tc.Profile{Type: "ATS_PROFILE", Name: "badprofile", Description: "", CDNID: 1},
		tc.Profile{Type: "ATS_PROFILE", Name: "", Description: "description", CDNID: 1},
		tc.Profile{Type: "", Name: "badprofile", Description: "description", CDNID: 1},
	}

	for _, pr := range prs {
		resp, _, err := TOSession.CreateProfile(pr)

		if err == nil {
			t.Errorf("Creating bad profile succeeded: %+v\nResponse is %+v", pr, resp)
		}
	}
}

func (r *TCData) CreateTestProfiles(t *testing.T) {

	for _, pr := range r.TestData.Profiles {
		resp, _, err := TOSession.CreateProfile(pr)

		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE profiles with name: %s %v", pr.Name, err)
		}
		profiles, _, err := TOSession.GetProfileByName(pr.Name)
		if err != nil {
			t.Errorf("could not GET profile with name: %s %v", pr.Name, err)
		}
		if len(profiles) == 0 {
			t.Errorf("could not GET profile %+v: not found", pr)
		}
		profileID := profiles[0].ID

		for _, param := range pr.Parameters {
			if param.Name == nil || param.Value == nil || param.ConfigFile == nil {
				t.Errorf("invalid parameter specification: %+v", param)
				continue
			}
			_, _, err := TOSession.CreateParameter(tc.Parameter{Name: *param.Name, Value: *param.Value, ConfigFile: *param.ConfigFile})
			if err != nil {
				// ok if already exists
				if !strings.Contains(err.Error(), "already exists") {
					t.Errorf("could not CREATE parameter %+v: %s", param, err.Error())
					continue
				}
			}
			p, _, err := TOSession.GetParameterByNameAndConfigFileAndValue(*param.Name, *param.ConfigFile, *param.Value)
			if err != nil {
				t.Errorf("could not GET parameter %+v: %s", param, err.Error())
			}
			if len(p) == 0 {
				t.Errorf("could not GET parameter %+v: not found", param)
			}
			_, _, err = TOSession.CreateProfileParameter(tc.ProfileParameter{ProfileID: profileID, ParameterID: p[0].ID})
			if err != nil {
				t.Errorf("could not CREATE profile_parameter %+v: %s", param, err.Error())
			}
		}

	}
}

func (r *TCData) DeleteTestProfiles(t *testing.T) {

	for _, pr := range r.TestData.Profiles {
		// Retrieve the Profile by name so we can get the id for the Update
		resp, _, err := TOSession.GetProfileByName(pr.Name)
		if err != nil {
			t.Errorf("cannot GET Profile by name: %s - %v", pr.Name, err)
			continue
		}
		if len(resp) == 0 {
			t.Errorf("cannot GET Profile by name: not found - %s", pr.Name)
			continue
		}

		profileID := resp[0].ID
		// query by name does not retrieve associated parameters.  But query by id does.
		resp, _, err = TOSession.GetProfileByID(profileID)
		if err != nil {
			t.Errorf("cannot GET Profile by id: %v - %v", err, resp)
		}
		// delete any profile_parameter associations first
		// the parameter is what's being deleted, but the delete is cascaded to profile_parameter
		for _, param := range resp[0].Parameters {
			_, _, err := TOSession.DeleteParameterByID(*param.ID)
			if err != nil {
				t.Errorf("cannot DELETE parameter with parameterID %d: %s", *param.ID, err.Error())
			}
		}
		delResp, _, err := TOSession.DeleteProfileByID(profileID)
		if err != nil {
			t.Errorf("cannot DELETE Profile by name: %v - %v", err, delResp)
		}
		//time.Sleep(1 * time.Second)

		// Retrieve the Profile to see if it got deleted
		prs, _, err := TOSession.GetProfileByName(pr.Name)
		if err != nil {
			t.Errorf("error deleting Profile name: %s", err.Error())
		}
		if len(prs) > 0 {
			t.Errorf("expected Profile Name: %s to be deleted", pr.Name)
		}

		// Attempt to export Profile
		_, _, err = TOSession.ExportProfile(profileID)
		if err == nil {
			t.Errorf("expected Profile: %s to be nil on export", pr.Name)
		}
	}
}
