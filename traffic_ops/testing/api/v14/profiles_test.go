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

package v14

import (
	"testing"

	"github.com/apache/trafficcontrol/lib/go-log"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

func TestProfiles(t *testing.T) {

	CreateTestCDNs(t)
	CreateTestTypes(t)

	// attempt to create profiles with missing info
	CreateBadProfiles(t)
	CreateTestProfiles(t)
	CreateTestParameters(t)
	CreateTestProfileParameters(t)
	UpdateTestProfiles(t)
	GetTestProfiles(t)
	GetTestProfilesWithParameters(t)
	DeleteTestProfileParameters(t)
	DeleteTestParameters(t)
	DeleteTestProfiles(t)
	DeleteTestTypes(t)
	DeleteTestCDNs(t)
}

// CreateBadProfiles ensures that profiles can't be created with bad values
func CreateBadProfiles(t *testing.T) {

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
		} else {
			log.Debugf("bad profile creation failed appropriately")
		}
	}
}

func CreateTestProfiles(t *testing.T) {

	for _, pr := range testData.Profiles {
		cdns, _, err := TOSession.GetCDNByName(pr.CDNName)
		respCDN := cdns[0]
		pr.CDNID = respCDN.ID

		resp, _, err := TOSession.CreateProfile(pr)

		log.Debugln("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE profiles with name: %s %v\n", pr.Name, err)
		}
	}
}

func UpdateTestProfiles(t *testing.T) {

	firstProfile := testData.Profiles[0]
	// Retrieve the Profile by name so we can get the id for the Update
	resp, _, err := TOSession.GetProfileByName(firstProfile.Name)
	if err != nil {
		t.Errorf("cannot GET Profile by name: %v - %v\n", firstProfile.Name, err)
	}
	remoteProfile := resp[0]
	expectedProfileDesc := "UPDATED"
	remoteProfile.Description = expectedProfileDesc
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateProfileByID(remoteProfile.ID, remoteProfile)
	if err != nil {
		t.Errorf("cannot UPDATE Profile by id: %v - %v\n", err, alert)
	}

	// Retrieve the Profile to check Profile name got updated
	resp, _, err = TOSession.GetProfileByID(remoteProfile.ID)
	if err != nil {
		t.Errorf("cannot GET Profile by name: %v - %v\n", firstProfile.Name, err)
	}
	respProfile := resp[0]
	if respProfile.Description != expectedProfileDesc {
		t.Errorf("results do not match actual: %s, expected: %s\n", respProfile.Description, expectedProfileDesc)
	}

}

func GetTestProfiles(t *testing.T) {

	for _, pr := range testData.Profiles {
		resp, _, err := TOSession.GetProfileByName(pr.Name)
		if err != nil {
			t.Errorf("cannot GET Profile by name: %v - %v\n", err, resp)
		}

		resp, _, err = TOSession.GetProfileByParameter(pr.Parameter)
		if err != nil {
			t.Errorf("cannot GET Profile by param: %v - %v\n", err, resp)
		}

		resp, _, err = TOSession.GetProfileByCDNID(pr.CDNID)
		if err != nil {
			t.Errorf("cannot GET Profile by cdn: %v - %v\n", err, resp)
		}
	}
}
func GetTestProfilesWithParameters(t *testing.T) {
	firstProfile := testData.Profiles[0]
	resp, _, err := TOSession.GetProfileByName(firstProfile.Name)
	if len(resp) > 0 {
		respProfile := resp[0]
		resp, _, err := TOSession.GetProfileByID(respProfile.ID)
		if err != nil {
			t.Errorf("cannot GET Profile by name: %v - %v\n", err, resp)
		}
		if len(resp) > 0 {
			respProfile = resp[0]
			respParameters := respProfile.Parameters
			if len(respParameters) == 0 {
				t.Errorf("expected a profile with parameters to be retrieved: %v - %v\n", err, respParameters)
			}
		}
	}
	if err != nil {
		t.Errorf("cannot GET Profile by name: %v - %v\n", err, resp)
	}
}

func DeleteTestProfiles(t *testing.T) {

	for _, pr := range testData.Profiles {
		// Retrieve the Profile by name so we can get the id for the Update
		resp, _, err := TOSession.GetProfileByName(pr.Name)
		if err != nil {
			t.Errorf("cannot GET Profile by name: %v - %v\n", pr.Name, err)
		}
		if len(resp) > 0 {
			respProfile := resp[0]

			delResp, _, err := TOSession.DeleteProfileByID(respProfile.ID)
			if err != nil {
				t.Errorf("cannot DELETE Profile by name: %v - %v\n", err, delResp)
			}
			//time.Sleep(1 * time.Second)

			// Retrieve the Profile to see if it got deleted
			prs, _, err := TOSession.GetProfileByName(pr.Name)
			if err != nil {
				t.Errorf("error deleting Profile name: %s\n", err.Error())
			}
			if len(prs) > 0 {
				t.Errorf("expected Profile Name: %s to be deleted\n", pr.Name)
			}
		}
	}
}
