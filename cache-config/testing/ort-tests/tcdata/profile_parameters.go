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
	"fmt"
	"sync"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const queryParamFormat = "?profileId=%d&parameterId=%d"

func (r *TCData) CreateTestProfileParameters(t *testing.T) {

	firstProfile := r.TestData.Profiles[0]
	profileResp, _, err := TOSession.GetProfileByName(firstProfile.Name)
	if err != nil {
		t.Errorf("cannot GET Profile by name: %s - %v", firstProfile.Name, err)
	}

	firstParameter := r.TestData.Parameters[0]
	paramResp, _, err := TOSession.GetParameterByName(firstParameter.Name)
	if err != nil {
		t.Errorf("cannot GET Parameter by name: %s - %v", firstParameter.Name, err)
	}

	profileID := profileResp[0].ID
	parameterID := paramResp[0].ID

	pp := tc.ProfileParameter{
		ProfileID:   profileID,
		ParameterID: parameterID,
	}
	resp, _, err := TOSession.CreateProfileParameter(pp)
	t.Log("Response: ", resp)
	if err != nil {
		t.Errorf("could not CREATE profile parameters: %v", err)
	}

}

func (r *TCData) DeleteTestProfileParametersParallel(t *testing.T) {

	var wg sync.WaitGroup
	for _, pp := range r.TestData.ProfileParameters {

		wg.Add(1)
		go func(p tc.ProfileParameter) {
			defer wg.Done()
			DeleteTestProfileParameter(t, p)
		}(pp)

	}
	wg.Wait()
}

func (r *TCData) DeleteTestProfileParameters(t *testing.T) {

	for _, pp := range r.TestData.ProfileParameters {
		DeleteTestProfileParameter(t, pp)
	}
}

func DeleteTestProfileParameter(t *testing.T, pp tc.ProfileParameter) {

	queryParams := fmt.Sprintf(queryParamFormat, pp.ProfileID, pp.ParameterID)
	// Retrieve the PtofileParameter by profile so we can get the id for the Update
	resp, _, err := TOSession.GetProfileParameterByQueryParams(queryParams)
	if err != nil {
		t.Errorf("cannot GET Parameter by profile: %s - %v", pp.Profile, err)
	}
	if len(resp) > 0 {
		respPP := resp[0]

		delResp, _, err := TOSession.DeleteParameterByProfileParameter(respPP.ProfileID, respPP.ParameterID)
		if err != nil {
			t.Errorf("cannot DELETE Parameter by profile: %v - %v", err, delResp)
		}

		// Retrieve the Parameter to see if it got deleted
		pps, _, err := TOSession.GetProfileParameterByQueryParams(queryParams)
		if err != nil {
			t.Errorf("error deleting Parameter name: %s", err.Error())
		}
		if len(pps) > 0 {
			t.Errorf("expected Parameter Name: %s and ConfigFile: %s to be deleted", pp.Profile, pp.Parameter)
		}
	}
}
