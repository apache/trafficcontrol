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

package v3

import (
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"net/http"
	"sync"
	"testing"
	"time"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

func TestParameters(t *testing.T) {

	//toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	//SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Portal, Config.TrafficOps.UserPassword)

	WithObjs(t, []TCObj{Parameters}, func() {
		GetTestParametersIMS(t)
		currentTime := time.Now().Add(-1 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		UpdateTestParameters(t)
		GetTestParameters(t)
		GetTestParametersIMSAfterChange(t, header)
	})
}

func GetTestParametersIMSAfterChange(t *testing.T, header http.Header) {
	for _, pl := range testData.Parameters {
		_, reqInf, err := TOSession.GetParameterByName(pl.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}
	currentTime := time.Now()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	for _, pl := range testData.Parameters {
		_, reqInf, err := TOSession.GetParameterByName(pl.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestParameters(t *testing.T) {

	for _, pl := range testData.Parameters {
		resp, _, err := TOSession.CreateParameter(pl)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE parameters: %v", err)
		}
	}

}

func UpdateTestParameters(t *testing.T) {

	firstParameter := testData.Parameters[0]
	// Retrieve the Parameter by name so we can get the id for the Update
	resp, _, err := TOSession.GetParameterByName(firstParameter.Name, nil)
	if err != nil {
		t.Errorf("cannot GET Parameter by name: %v - %v", firstParameter.Name, err)
	}
	remoteParameter := resp[0]
	expectedParameterValue := "UPDATED"
	remoteParameter.Value = expectedParameterValue
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateParameterByID(remoteParameter.ID, remoteParameter)
	if err != nil {
		t.Errorf("cannot UPDATE Parameter by id: %v - %v", err, alert)
	}

	// Retrieve the Parameter to check Parameter name got updated
	resp, _, err = TOSession.GetParameterByID(remoteParameter.ID, nil)
	if err != nil {
		t.Errorf("cannot GET Parameter by name: %v - %v", firstParameter.Name, err)
	}
	respParameter := resp[0]
	if respParameter.Value != expectedParameterValue {
		t.Errorf("results do not match actual: %s, expected: %s", respParameter.Value, expectedParameterValue)
	}

}

func GetTestParametersIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0,0,1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	for _, pl := range testData.Parameters {
		_, reqInf, err := TOSession.GetParameterByName(pl.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestParameters(t *testing.T) {

	for _, pl := range testData.Parameters {
		resp, _, err := TOSession.GetParameterByName(pl.Name, nil)
		if err != nil {
			t.Errorf("cannot GET Parameter by name: %v - %v", err, resp)
		}
	}
}

func DeleteTestParametersParallel(t *testing.T) {

	var wg sync.WaitGroup
	for _, pl := range testData.Parameters {

		wg.Add(1)
		go func() {
			defer wg.Done()
			DeleteTestParameter(t, pl)
		}()

	}
	wg.Wait()
}

func DeleteTestParameters(t *testing.T) {

	for _, pl := range testData.Parameters {
		DeleteTestParameter(t, pl)
	}
}

func DeleteTestParameter(t *testing.T, pl tc.Parameter) {

	// Retrieve the Parameter by name so we can get the id for the Update
	resp, _, err := TOSession.GetParameterByNameAndConfigFile(pl.Name, pl.ConfigFile, nil)
	if err != nil {
		t.Errorf("cannot GET Parameter by name: %v - %v", pl.Name, err)
	}

	if len(resp) == 0 {
		// TODO This fails for the ProfileParameters test; determine a way to check this, even for ProfileParameters
		// t.Errorf("DeleteTestParameter got no params for %+v %+v", pl.Name, pl.ConfigFile)
	} else if len(resp) > 1 {
		// TODO figure out why this happens, and be more precise about deleting things where created.
		// t.Errorf("DeleteTestParameter params for %+v %+v expected 1, actual %+v", pl.Name, pl.ConfigFile, len(resp))
	}
	for _, respParameter := range resp {
		delResp, _, err := TOSession.DeleteParameterByID(respParameter.ID)
		if err != nil {
			t.Errorf("cannot DELETE Parameter by name: %v - %v", err, delResp)
		}

		// Retrieve the Parameter to see if it got deleted
		pls, _, err := TOSession.GetParameterByID(pl.ID, nil)
		if err != nil {
			t.Errorf("error deleting Parameter name: %s", err.Error())
		}
		if len(pls) > 0 {
			t.Errorf("expected Parameter Name: %s and ConfigFile: %s to be deleted", pl.Name, pl.ConfigFile)
		}
	}
}
