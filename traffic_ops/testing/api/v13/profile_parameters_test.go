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
	"sync"
	"testing"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

func TestProfileParameters(t *testing.T) {

	CreateTestProfileParameters(t)
	UpdateTestProfileParameters(t)
	GetTestProfileParameters(t)
	DeleteTestProfileParameters(t)

}

func CreateTestProfileParameters(t *testing.T) {

	for _, pp := range testData.ProfileParameters {
		resp, _, err := TOSession.CreateProfileParameter(pp)
		log.Debugln("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE profile parameters: %v\n", err)
		}
	}

}

func UpdateTestProfileParameters(t *testing.T) {

	firstPP := testData.ProfileParameters[0]
	// Retrieve the Parameter by profile so we can get the id for the Update
	resp, _, err := TOSession.GetProfileParameterByParameter(firstPP.Profile)
	if err != nil {
		t.Errorf("cannot GET Parameter by name: %v - %v\n", firstPP.Profile, err)
	}
	remotePP := resp[0]
	expectedPP := 1
	remotePP.Profile = expectedPP
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateParameterByProfile(remotePP.Profile, remotePP)
	if err != nil {
		t.Errorf("cannot UPDATE Parameter by profile: %v - %v\n", err, alert)
	}

	// Retrieve the Parameter to check Parameter name got updated
	resp, _, err = TOSession.GetProfileParameterByProfile(remotePP.Profile)
	if err != nil {
		t.Errorf("cannot GET Parameter by profile: %v - %v\n", firstPP.Profile, err)
	}
	respParameter := resp[0]
	if respParameter.Profile != expectedPP {
		t.Errorf("results do not match actual: %s, expected: %s\n", respParameter.Profile, expectedPP)
	}

}

func GetTestProfileParameters(t *testing.T) {

	for _, pp := range testData.ProfileParameters {
		resp, _, err := TOSession.GetProfileParameterByProfile(pp.Profile)
		if err != nil {
			t.Errorf("cannot GET Parameter by name: %v - %v\n", err, resp)
		}
	}
}

func DeleteTestProfileParametersParallel(t *testing.T) {

	var wg sync.WaitGroup
	for _, pp := range testData.ProfileParameters {

		wg.Add(1)
		go func() {
			defer wg.Done()
			DeleteTestProfileParameter(t, pp)
		}()

	}
	wg.Wait()
}

func DeleteTestProfileParameters(t *testing.T) {

	for _, pp := range testData.ProfileParameters {
		DeleteTestProfileParameter(t, pp)
	}
}

func DeleteTestProfileParameter(t *testing.T, pp tc.ProfileParameter) {

	// Retrieve the PtofileParameter by profile so we can get the id for the Update
	resp, _, err := TOSession.GetProfileParameterByIDs(pp.Profile, pp.Parameter)
	if err != nil {
		t.Errorf("cannot GET Parameter by profile: %v - %v\n", pp.Profile, err)
	}
	if len(resp) > 0 {
		respPP := resp[0]

		delResp, _, err := TOSession.DeleteProfileParameterByProfile(respPP.Profile)
		if err != nil {
			t.Errorf("cannot DELETE Parameter by profile: %v - %v\n", err, delResp)
		}

		// Retrieve the Parameter to see if it got deleted
		pls, _, err := TOSession.GetProfileParameterByIDs(pp.Profile)
		if err != nil {
			t.Errorf("error deleting Parameter name: %s\n", err.Error())
		}
		if len(pls) > 0 {
			t.Errorf("expected Parameter Name: %s and ConfigFile: %s to be deleted\n", pp.Profile, pp.Parameter)
		}
	}
}
