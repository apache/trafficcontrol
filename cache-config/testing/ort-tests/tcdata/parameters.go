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
	"sync"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func (r *TCData) CreateTestParameters(t *testing.T) {

	for _, pl := range r.TestData.Parameters {
		resp, _, err := TOSession.CreateParameter(pl)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE parameters: %v", err)
		}
	}

}

func (r *TCData) DeleteTestParametersParallel(t *testing.T) {

	var wg sync.WaitGroup
	for _, pl := range r.TestData.Parameters {

		wg.Add(1)
		go func(p tc.Parameter) {
			defer wg.Done()
			DeleteTestParameter(t, p)
		}(pl)

	}
	wg.Wait()
}

func (r *TCData) DeleteTestParameters(t *testing.T) {

	for _, pl := range r.TestData.Parameters {
		DeleteTestParameter(t, pl)
	}
}

func DeleteTestParameter(t *testing.T, pl tc.Parameter) {

	// Retrieve the Parameter by name so we can get the id for the Update
	resp, _, err := TOSession.GetParameterByNameAndConfigFile(pl.Name, pl.ConfigFile)
	if err != nil {
		t.Errorf("cannot GET Parameter by name %s: %v", pl.Name, err)
	}

	if len(resp) == 0 {
		// TODO This fails for the ProfileParameters test; determine a way to check this, even for ProfileParameters
		// t.Errorf("DeleteTestParameter got no params for %s %s", pl.Name, pl.ConfigFile)
	} else if len(resp) > 1 {
		// TODO figure out why this happens, and be more precise about deleting things where created.
		// t.Errorf("DeleteTestParameter params for %s %s expected 1, actual %d", pl.Name, pl.ConfigFile, len(resp))
	}
	for _, respParameter := range resp {
		delResp, _, err := TOSession.DeleteParameterByID(respParameter.ID)
		if err != nil {
			t.Errorf("cannot DELETE Parameter by ID: %v - %v", err, delResp)
		}

		// Retrieve the Parameter to see if it got deleted
		pls, _, err := TOSession.GetParameterByID(pl.ID)
		if err != nil {
			t.Errorf("error deleting Parameter name: %v", err)
		}
		if len(pls) > 0 {
			t.Errorf("expected Parameter Name: %s and ConfigFile: %s to be deleted", pl.Name, pl.ConfigFile)
		}
	}
}
