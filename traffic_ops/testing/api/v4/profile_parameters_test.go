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
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestProfileParameters(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, ProfileParameters}, func() {
		GetTestProfileParametersIMS(t)
		GetTestProfileParameters(t)
		InvalidCreateTestProfileParameters(t)
	})
}

func GetTestProfileParametersIMS(t *testing.T) {
	opts := client.NewRequestOptions()

	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	opts.Header.Set(rfc.IfModifiedSince, time)

	for _, pp := range testData.ProfileParameters {
		opts.QueryParameters.Set("profileId", strconv.Itoa(pp.ProfileID))
		opts.QueryParameters.Set("parameterId", strconv.Itoa(pp.ParameterID))
		resp, reqInf, err := TOSession.GetProfileParameters(opts)
		if err != nil {
			t.Errorf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Errorf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestProfileParameters(t *testing.T) {
	if len(testData.Parameters) < 1 || len(testData.Profiles) < 1 {
		t.Fatal("Need at least one Profile and one Parameter to test associating a Parameter with a Profile")
	}

	firstProfile := testData.Profiles[0]
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", firstProfile.Name)
	profileResp, _, err := TOSession.GetProfiles(opts)
	if err != nil {
		t.Errorf("cannot get Profile '%s' by name: %v - alerts: %+v", firstProfile.Name, err, profileResp.Alerts)
	}
	if len(profileResp.Response) != 1 {
		t.Fatalf("Expected exactly one Profile to exist with name '%s', found: %d", firstProfile.Name, len(profileResp.Response))
	}

	firstParameter := testData.Parameters[0]
	opts.QueryParameters.Set("name", firstParameter.Name)
	paramResp, _, err := TOSession.GetParameters(opts)
	if err != nil {
		t.Errorf("cannot get Parameter by name '%s': %v - alerts: %+v", firstParameter.Name, err, paramResp.Alerts)
	}
	if len(paramResp.Response) < 1 {
		t.Fatalf("Expected at least one Parameter to exist with name '%s'", firstParameter.Name)
	}

	profileID := profileResp.Response[0].ID
	parameterID := paramResp.Response[0].ID

	pp := tc.ProfileParameterCreationRequest{
		ProfileID:   profileID,
		ParameterID: parameterID,
	}
	resp, _, err := TOSession.CreateProfileParameter(pp, client.RequestOptions{})
	if err != nil {
		t.Errorf("could not associate parameters to profile: %v - alerts: %+v", err, resp.Alerts)
	}

}

func InvalidCreateTestProfileParameters(t *testing.T) {
	pp := tc.ProfileParameterCreationRequest{
		ProfileID:   0,
		ParameterID: 0,
	}
	resp, _, err := TOSession.CreateProfileParameter(pp, client.RequestOptions{})
	if err == nil {
		t.Fatalf("creating invalid profile parameter - expected: error, actual: nil")
	}

	foundProfile := false
	foundParam := false
	for _, alert := range resp.Alerts {
		if alert.Level == tc.ErrorLevel.String() {
			if strings.Contains(alert.Text, "profileId") {
				foundProfile = true
			}
			if strings.Contains(alert.Text, "parameterId") {
				foundParam = true
			}
			if foundProfile && foundParam {
				break
			}
		}
	}

	if !foundProfile {
		t.Errorf("expected: error message to contain 'profileId', actual: %v - alerts: %+v", err, resp.Alerts)
	}
	if !foundParam {
		t.Errorf("expected: error message to contain 'parameterId', actual: %v - alerts: %+v", err, resp.Alerts)
	}

}

func GetTestProfileParameters(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, pp := range testData.ProfileParameters {
		opts.QueryParameters.Set("profileId", strconv.Itoa(pp.ProfileID))
		opts.QueryParameters.Set("parameterId", strconv.Itoa(pp.ParameterID))
		resp, _, err := TOSession.GetProfileParameters(opts)
		if err != nil {
			t.Errorf("cannot get Profile #%d/Parameter #%d association: %v - alerts: %+v", pp.ProfileID, pp.ParameterID, err, resp.Alerts)
		}
	}
}

func DeleteTestProfileParameters(t *testing.T) {

	for _, pp := range testData.ProfileParameters {
		DeleteTestProfileParameter(t, pp)
	}
}

func DeleteTestProfileParameter(t *testing.T, pp tc.ProfileParameter) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("profileId", strconv.Itoa(pp.ProfileID))
	opts.QueryParameters.Set("parameterId", strconv.Itoa(pp.ParameterID))
	// Retrieve the PtofileParameter by profile so we can get the id for the Update
	resp, _, err := TOSession.GetProfileParameters(opts)
	if err != nil {
		t.Errorf("cannot get Profile #%d/Parameter #%d association: %v - alerts: %+v", pp.ProfileID, pp.ParameterID, err, resp.Alerts)
	}
	if len(resp.Response) > 0 {
		respPP := resp.Response[0]

		delResp, _, err := TOSession.DeleteProfileParameter(pp.ProfileID, respPP.Parameter, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot delete Profile #%d/Parameter #%d association: %v - alerts: %+v", pp.ProfileID, pp.ParameterID, err, delResp.Alerts)
		}

		// Retrieve the Parameter to see if it got deleted
		pps, _, err := TOSession.GetProfileParameters(opts)
		if err != nil {
			t.Errorf("error getting #%d/Parameter #%d association after deletion: %v - alerts: %+v", pp.ProfileID, pp.ParameterID, err, pps.Alerts)
		}
		if len(pps.Response) > 0 {
			t.Errorf("expected #%d/Parameter #%d association to be deleted, but it was found in Traffic Ops", pp.ProfileID, pp.ParameterID)
		}
	}
}
