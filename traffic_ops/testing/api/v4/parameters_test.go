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

package v4

import (
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestParameters(t *testing.T) {

	//toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	//SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Portal, Config.TrafficOps.UserPassword)

	WithObjs(t, []TCObj{Parameters}, func() {
		GetTestParametersIMS(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		header.Set(rfc.IfUnmodifiedSince, time)
		UpdateTestParameters(t)
		UpdateTestParametersWithHeaders(t, header)
		GetTestParameters(t)
		GetTestParametersIMSAfterChange(t, header)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestParametersWithHeaders(t, header)
	})
}

func UpdateTestParametersWithHeaders(t *testing.T, header http.Header) {
	if len(testData.Parameters) < 1 {
		t.Fatal("Need at least one Parameter to test updating Parameters with HTTP headers")
	}
	firstParameter := testData.Parameters[0]

	// Retrieve the Parameter by name so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.Header = header
	resp, _, err := TOSession.GetParametersByProfileName(firstParameter.Name, opts)
	if err != nil {
		t.Errorf("cannot get Parameter by name '%s': %v - alerts: %+v", firstParameter.Name, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one Parameter to exist with name '%s'", firstParameter.Name)
	}
	remoteParameter := resp.Response[0]
	expectedParameterValue := "UPDATED"
	remoteParameter.Value = expectedParameterValue
	_, reqInf, err := TOSession.UpdateParameter(remoteParameter.ID, remoteParameter, opts)
	if err == nil {
		t.Error("Expected error about precondition failed, but got none")
	}
	if reqInf.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("Expected status code 412, got %d", reqInf.StatusCode)
	}
}

func GetTestParametersIMSAfterChange(t *testing.T, header http.Header) {
	opts := client.NewRequestOptions()
	opts.Header = header
	for _, pl := range testData.Parameters {
		opts.QueryParameters.Set("name", pl.Name)
		resp, reqInf, err := TOSession.GetParameters(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %d", reqInf.StatusCode)
		}
	}

	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)

	opts.Header.Set(rfc.IfModifiedSince, timeStr)
	for _, pl := range testData.Parameters {
		opts.QueryParameters.Set("name", pl.Name)
		resp, reqInf, err := TOSession.GetParameters(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %d", reqInf.StatusCode)
		}
	}
}

func CreateTestParameters(t *testing.T) {
	resp, _, err := TOSession.CreateMultipleParameters(testData.Parameters, client.RequestOptions{})
	if err != nil {
		t.Errorf("could not create Parameters: %v - alerts: %+v", err, resp)
	}
}

func UpdateTestParameters(t *testing.T) {
	if len(testData.Parameters) < 1 {
		t.Fatal("Need at least one Parameter to test updating Parameters")
	}
	firstParameter := testData.Parameters[0]

	// Retrieve the Parameter by name so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", firstParameter.Name)
	resp, _, err := TOSession.GetParameters(opts)
	if err != nil {
		t.Errorf("cannot get Parameter by name '%s': %v - alerts: %+v", firstParameter.Name, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one Parameter to exist with name '%s'", firstParameter.Name)
	}
	remoteParameter := resp.Response[0]

	expectedParameterValue := "UPDATED"
	remoteParameter.Value = expectedParameterValue
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateParameter(remoteParameter.ID, remoteParameter, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Parameter: %v - alerts: %+v", err, alert.Alerts)
	}

	// Retrieve the Parameter to check Parameter name got updated
	opts.QueryParameters.Del("name")
	opts.QueryParameters.Set("id", strconv.Itoa(remoteParameter.ID))
	resp, _, err = TOSession.GetParameters(opts)
	if err != nil {
		t.Errorf("cannot get Parameter by ID %d: %v - alerts: %+v", remoteParameter.ID, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Parameter to exist with ID %d, found: %d", remoteParameter.ID, len(resp.Response))
	}
	respParameter := resp.Response[0]
	if respParameter.Value != expectedParameterValue {
		t.Errorf("results do not match actual: %s, expected: %s", respParameter.Value, expectedParameterValue)
	}

}

func GetTestParametersIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)
	for _, pl := range testData.Parameters {
		opts.QueryParameters.Set("name", pl.Name)
		resp, reqInf, err := TOSession.GetParameters(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %d", reqInf.StatusCode)
		}
	}
}

func GetTestParameters(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, pl := range testData.Parameters {
		opts.QueryParameters.Set("name", pl.Name)
		resp, _, err := TOSession.GetParameters(opts)
		if err != nil {
			t.Errorf("cannot GET Parameter by name '%s': %v - alerts: %+v", pl.Name, err, resp.Alerts)
		}
	}
}

func DeleteTestParametersParallel(t *testing.T) {
	var wg sync.WaitGroup
	for _, pl := range testData.Parameters {

		wg.Add(1)
		go func(p tc.Parameter) {
			defer wg.Done()
			DeleteTestParameter(t, p)
		}(pl)

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
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", pl.Name)
	opts.QueryParameters.Set("configFile", pl.ConfigFile)
	resp, _, err := TOSession.GetParameters(opts)
	if err != nil {
		t.Errorf("cannot get Parameter by name '%s' and configFile '%s': %v - alerts: %+v", pl.Name, pl.ConfigFile, err, resp.Alerts)
	}

	// TODO This fails for the ProfileParameters test; determine a way to check this, even for ProfileParameters
	// if len(resp.Response) == 0 {
	// t.Errorf("DeleteTestParameter got no params for %+v %+v", pl.Name, pl.ConfigFile)
	// TODO figure out why this happens, and be more precise about deleting things where created.
	// } else if len(resp.Response) > 1 {
	// t.Errorf("DeleteTestParameter params for %+v %+v expected 1, actual %+v", pl.Name, pl.ConfigFile, len(resp))
	// }

	opts.QueryParameters.Del("name")
	opts.QueryParameters.Del("configFile")
	for _, respParameter := range resp.Response {
		delResp, _, err := TOSession.DeleteParameter(respParameter.ID, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot delete Parameter #%d: %v - alerts: %+v", respParameter.ID, err, delResp.Alerts)
		}

		// Retrieve the Parameter to see if it got deleted
		opts.QueryParameters.Set("id", strconv.Itoa(pl.ID))
		pls, _, err := TOSession.GetParameters(opts)
		if err != nil {
			t.Errorf("Unexpected error fetching Parameter #%d after deletion: %v - alerts: %+v", pl.ID, err, pls.Alerts)
		}
		if len(pls.Response) > 0 {
			t.Errorf("expected Parameter with name '%s' and configFile '%s' to be deleted, but it was found in a Traffic Ops response", pl.Name, pl.ConfigFile)
		}
	}
}
