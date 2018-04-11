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

package v13

import (
	"fmt"
	"testing"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

func TestProfiles(t *testing.T) {

	CreateTestCDNs(t)
	CreateTestTypes(t)
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

func CreateTestProfiles(t *testing.T) {

	for _, pr := range testData.Profiles {
		cdns, _, err := TOSession.GetCDNByName(pr.CDNName)
		respCDN := cdns[0]
		cdnName := respCDN.Name
		fmt.Printf("profileName: %s, cdnName %s\n", pr.Name, cdnName)
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
