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
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestProfiles(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Profiles, Parameters}, func() {
		CreateBadProfiles(t)
		UpdateTestProfiles(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfUnmodifiedSince, time)
		UpdateTestProfilesWithHeaders(t, header)
		GetTestProfilesIMS(t)
		GetTestProfiles(t)
		GetTestProfilesWithParameters(t)
		ImportProfile(t)
		CopyProfile(t)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestProfilesWithHeaders(t, header)
		GetTestPaginationSupportProfiles(t)
	})
}

func UpdateTestProfilesWithHeaders(t *testing.T, header http.Header) {
	if len(testData.Profiles) < 1 {
		t.Fatal("Need at least one Profile to test updating a Profile with HTTP headers")
	}
	firstProfile := testData.Profiles[0]

	// Retrieve the Profile by name so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.Header = header
	opts.QueryParameters.Set("name", firstProfile.Name)
	resp, _, err := TOSession.GetProfiles(opts)
	if err != nil {
		t.Errorf("cannot get Profile '%s' by name: %v - alerts: %+v", firstProfile.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Profile to exist with name '%s', found: %d", firstProfile.Name, len(resp.Response))
	}

	remoteProfile := resp.Response[0]
	opts.QueryParameters.Del("name")
	_, reqInf, err := TOSession.UpdateProfile(remoteProfile.ID, remoteProfile, opts)
	if err == nil {
		t.Errorf("Expected error about precondition failed, but got none")
	}
	if reqInf.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
	}
}

func GetTestProfilesIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)
	for _, pr := range testData.Profiles {
		opts.QueryParameters.Set("name", pr.Name)
		resp, reqInf, err := TOSession.GetProfiles(opts)
		if err != nil {
			t.Errorf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Errorf("Expected 304 status code, got %d", reqInf.StatusCode)
		}
	}
}

// CreateBadProfiles ensures that profiles can't be created with bad values
func CreateBadProfiles(t *testing.T) {

	// blank profile
	prs := []tc.Profile{
		{Type: "", Name: "", Description: "", CDNID: 0},
		{Type: "ATS_PROFILE", Name: "badprofile", Description: "description", CDNID: 0},
		{Type: "ATS_PROFILE", Name: "badprofile", Description: "", CDNID: 1},
		{Type: "ATS_PROFILE", Name: "", Description: "description", CDNID: 1},
		{Type: "", Name: "badprofile", Description: "description", CDNID: 1},
	}

	for _, pr := range prs {
		resp, _, err := TOSession.CreateProfile(pr, client.RequestOptions{})

		if err == nil {
			t.Errorf("Creating bad profile %+v succeeded, response is: %+v", pr, resp)
		}
	}
}

func CopyProfile(t *testing.T) {
	testCases := []struct {
		description  string
		profile      tc.ProfileCopy
		expectedResp string
		err          string
	}{
		{
			description: "copy profile",
			profile: tc.ProfileCopy{
				Name:         "profile-2",
				ExistingName: "EDGE1",
			},
			expectedResp: "created new profile [profile-2] from existing profile [EDGE1]",
		},
		{
			description: "existing profile does not exist",
			profile: tc.ProfileCopy{
				Name:         "profile-3",
				ExistingName: "bogus",
			},
			err: "profile with name bogus does not exist",
		},
		{
			description: "new profile already exists",
			profile: tc.ProfileCopy{
				Name:         "EDGE2",
				ExistingName: "EDGE1",
			},
			err: "profile with name EDGE2 already exists",
		},
	}

	var newProfileNames []string
	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			resp, _, err := TOSession.CopyProfile(c.profile, client.RequestOptions{})
			if c.err != "" {
				if err == nil {
					t.Fatalf("Expected an error like '%s', but got none", c.err)
				}
				found := false
				for _, alert := range resp.Alerts.Alerts {
					if alert.Level == tc.ErrorLevel.String() && strings.Contains(alert.Text, c.err) {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("Didn't find expected error-level alert '%s': %v - alerts: %+v", c.err, err, resp.Alerts)
				}
			} else if err != nil {
				t.Fatalf("Unexpected error: %v - alerts: %+v", err, resp.Alerts)
			}

			if err == nil {
				if got, want := resp.Alerts.ToStrings()[0], c.expectedResp; got != want {
					t.Fatalf("got= %s; expected= %s", got, want)
				}

				newProfileNames = append(newProfileNames, c.profile.Name)
			}
		})
	}

	// Cleanup profiles
	opts := client.NewRequestOptions()
	for _, name := range newProfileNames {
		opts.QueryParameters.Set("name", name)
		profiles, _, err := TOSession.GetProfiles(opts)
		if err != nil {
			t.Errorf("Unexpected error getting Profile '%s': %v - alerts: %+v", name, err, profiles.Alerts)
		}
		if len(profiles.Response) != 1 {
			t.Errorf("Expected exactly one Profile to exist with name '%s', found: %d", name, len(profiles.Response))
			continue
		}
		alerts, _, err := TOSession.DeleteProfile(profiles.Response[0].ID, client.RequestOptions{})
		if err != nil {
			t.Errorf("Unexpected error deleting Profile '%s' (#%d): %v - alerts: %+v", name, profiles.Response[0].ID, err, alerts.Alerts)
		}
	}
}

func CreateTestProfiles(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, pr := range testData.Profiles {
		resp, _, err := TOSession.CreateProfile(pr, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create Profile '%s': %v - alerts: %+v", pr.Name, err, resp.Alerts)
		}

		opts.QueryParameters.Set("name", pr.Name)
		profiles, _, err := TOSession.GetProfiles(opts)
		if err != nil {
			t.Errorf("could not GET profile with name: %s %v", pr.Name, err)
		}
		if len(profiles.Response) != 1 {
			t.Errorf("Expected exactly one Profile to exist with name '%s', found: %d", pr.Name, len(profiles.Response))
			continue
		}
		profileID := profiles.Response[0].ID

		paramOpts := client.NewRequestOptions()
		for _, param := range pr.Parameters {
			if param.Name == nil || param.Value == nil || param.ConfigFile == nil {
				t.Errorf("invalid parameter specification: %+v", param)
				continue
			}
			alerts, _, err := TOSession.CreateParameter(tc.Parameter{Name: *param.Name, Value: *param.Value, ConfigFile: *param.ConfigFile}, client.RequestOptions{})
			if err != nil {
				found := false
				for _, alert := range alerts.Alerts {
					if alert.Level == tc.ErrorLevel.String() && strings.Contains(alert.Text, "already exists") {
						found = true
						break
					}
				}
				// ok if already exists
				if !found {
					t.Errorf("Unexpected error creating Parameter %+v: %v - alerts: %+v", param, err, alerts.Alerts)
					continue
				}
			}
			paramOpts.QueryParameters.Set("name", *param.Name)
			paramOpts.QueryParameters.Set("configFile", *param.ConfigFile)
			paramOpts.QueryParameters.Set("value", *param.Value)
			p, _, err := TOSession.GetParameters(paramOpts)
			if err != nil {
				t.Errorf("could not get Parameter %+v: %v - alerts: %+v", param, err, p.Alerts)
			}
			if len(p.Response) == 0 {
				t.Fatalf("could not get parameter %+v: not found", param)
			}
			req := tc.ProfileParameterCreationRequest{ProfileID: profileID, ParameterID: p.Response[0].ID}
			alerts, _, err = TOSession.CreateProfileParameter(req, client.RequestOptions{})
			if err != nil {
				t.Errorf("could not associate Parameter %+v with Profile #%d: %v - alerts: %+v", param, profileID, err, alerts.Alerts)
			}
		}

	}
}

// Note this test will break if certain changes are made to the content and/or
// structure of the testing CDN and/or Profile data collections.
func UpdateTestProfiles(t *testing.T) {
	if len(testData.Profiles) < 1 {
		t.Fatal("Need at least one Profile to test updating Profiles")
	}
	firstProfile := testData.Profiles[0]

	// Retrieve the Profile by name so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", firstProfile.Name)
	resp, _, err := TOSession.GetProfiles(opts)
	if err != nil {
		t.Errorf("cannot get Profile '%s' by name: %v - alerts %+v", firstProfile.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Profile to exist with name '%s', found: %d", firstProfile.Name, len(resp.Response))
	}
	remoteProfile := resp.Response[0]

	opts.QueryParameters.Set("name", "cdn2")
	cdns, _, err := TOSession.GetCDNs(opts)
	if err != nil {
		t.Errorf("Unexpected error getting CDNs filtered by name 'cdn2': %v - alerts: %+v", err, cdns.Alerts)
	}
	if len(cdns.Response) != 1 {
		t.Fatalf("Expected exactly one CDN to exist with name 'cdn2', found: %d", len(cdns.Response))
	}
	oldName := remoteProfile.Name

	expectedProfileDesc := "UPDATED"
	expectedCDNId := cdns.Response[0].ID
	expectedName := "testing"
	expectedRoutingDisabled := true
	expectedType := "TR_PROFILE"

	remoteProfile.Description = expectedProfileDesc
	remoteProfile.Type = expectedType
	remoteProfile.CDNID = expectedCDNId
	remoteProfile.Name = expectedName
	remoteProfile.RoutingDisabled = expectedRoutingDisabled

	alert, _, err := TOSession.UpdateProfile(remoteProfile.ID, remoteProfile, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Profile: %v - alerts: %+v", err, alert.Alerts)
	}

	// Retrieve the Profile to check Profile name got updated
	opts.QueryParameters.Del("name")
	opts.QueryParameters.Set("id", strconv.Itoa(remoteProfile.ID))
	resp, _, err = TOSession.GetProfiles(opts)
	if err != nil {
		t.Errorf("cannot get Profile '%s' by ID (%d): %v - alerts: %+v", firstProfile.Name, remoteProfile.ID, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Profile to exist with ID %d, found: %d", remoteProfile.ID, len(resp.Response))
	}
	respProfile := resp.Response[0]
	if respProfile.Description != expectedProfileDesc {
		t.Errorf("results do not match actual: %s, expected: %s", respProfile.Description, expectedProfileDesc)
	}
	if respProfile.Type != expectedType {
		t.Errorf("results do not match actual: %s, expected: %s", respProfile.Type, expectedType)
	}
	if respProfile.CDNID != expectedCDNId {
		t.Errorf("results do not match actual: %d, expected: %d", respProfile.CDNID, expectedCDNId)
	}
	if respProfile.RoutingDisabled != expectedRoutingDisabled {
		t.Errorf("results do not match actual: %t, expected: %t", respProfile.RoutingDisabled, expectedRoutingDisabled)
	}
	if respProfile.Name != expectedName {
		t.Errorf("results do not match actual: %v, expected: %v", respProfile.Name, expectedName)
	}

	respProfile.Name = oldName
	alert, _, err = TOSession.UpdateProfile(respProfile.ID, respProfile, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error restoring Profile name: %v - alerts: %+v", err, alert.Alerts)
	}
}

func GetTestProfiles(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, pr := range testData.Profiles {
		opts.QueryParameters.Set("name", pr.Name)
		resp, _, err := TOSession.GetProfiles(opts)
		opts.QueryParameters.Del("name")
		if err != nil {
			t.Errorf("cannot get Profile '%s' by name: %v - alerts: %+v", pr.Name, err, resp.Alerts)
		}
		if len(resp.Response) != 1 {
			t.Errorf("Expected exactly one Profile to exist with name '%s', found: %d", pr.Name, len(resp.Response))
			continue
		}
		profileID := resp.Response[0].ID

		// TODO: figure out what the 'Parameter' field of a Profile is and how
		// passing it to this is supposed to work.
		// resp, _, err = TOSession.GetProfileByParameter(pr.Parameter)
		// if err != nil {
		// 	t.Errorf("cannot GET Profile by param: %v - %v", err, resp)
		// }

		opts.QueryParameters.Set("cdn", strconv.Itoa(pr.CDNID))
		resp, _, err = TOSession.GetProfiles(opts)
		opts.QueryParameters.Del("cdn")
		if err != nil {
			t.Errorf("cannot get Profiles by CDN ID %d: %v - alerts: %+v", pr.CDNID, err, resp.Alerts)
		}

		// Export Profile
		exportResp, _, err := TOSession.ExportProfile(profileID, client.RequestOptions{})
		if err != nil {
			t.Errorf("error exporting Profile #%d: %v - alerts: %+v", profileID, err, exportResp.Alerts)
		}
	}
}

func ImportProfile(t *testing.T) {
	if len(testData.Profiles) < 1 {
		t.Fatal("Need at least one Profile to test importing Profiles")
	}

	// Get ID of Profile to export
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", testData.Profiles[0].Name)
	resp, _, err := TOSession.GetProfiles(opts)
	if err != nil {
		t.Errorf("cannot get Profile '%s' by name: %v - alerts: %+v", testData.Profiles[0].Name, err, resp)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Profiles expected 1, actual %v", len(resp.Response))
	}
	profileID := resp.Response[0].ID

	// Export Profile to import
	exportResp, _, err := TOSession.ExportProfile(profileID, client.RequestOptions{})
	if err != nil {
		t.Fatalf("error exporting Profile #%d: %v - alerts: %+v", profileID, err, exportResp.Alerts)
	}

	// Modify Profile and import

	// Add parameter and change name
	profile := exportResp.Profile
	profile.Name = util.StrPtr("TestProfileImport")

	newParam := tc.ProfileExportImportParameterNullable{
		ConfigFile: util.StrPtr("config_file_import_test"),
		Name:       util.StrPtr("param_import_test"),
		Value:      util.StrPtr("import_test"),
	}
	parameters := append(exportResp.Parameters, newParam)
	// Import Profile
	importReq := tc.ProfileImportRequest{
		Profile:    profile,
		Parameters: parameters,
	}
	importResp, _, err := TOSession.ImportProfile(importReq, client.RequestOptions{})
	if err != nil {
		t.Fatalf("error importing Profile #%d: %v - alerts: %+v", profileID, err, importResp.Alerts)
	}

	// TODO: just delete it now?
	// Add newly create profile and parameter to testData so it gets deleted
	testData.Profiles = append(testData.Profiles, tc.Profile{
		Name:        *profile.Name,
		CDNName:     *profile.CDNName,
		Description: *profile.Description,
		Type:        *profile.Type,
	})

	testData.Parameters = append(testData.Parameters, tc.Parameter{
		ConfigFile: *newParam.ConfigFile,
		Name:       *newParam.Name,
		Value:      *newParam.Value,
	})
}

func GetTestProfilesWithParameters(t *testing.T) {
	if len(testData.Profiles) < 1 {
		t.Fatal("Need at least one Profile to test updating Profiles")
	}
	firstProfile := testData.Profiles[0]

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", firstProfile.Name)
	resp, _, err := TOSession.GetProfiles(opts)
	if err != nil {
		t.Errorf("cannot get Profile '%s' by name: %v - alerts: %+v", firstProfile.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Profile to exist with name '%s', found: %d", firstProfile.Name, len(resp.Response))
	}
	respProfile := resp.Response[0]

	// query by name does not retrieve associated parameters. But query by id does.
	// TODO ... what??
	opts.QueryParameters.Del("name")
	opts.QueryParameters.Set("id", strconv.Itoa(respProfile.ID))
	resp, _, err = TOSession.GetProfiles(opts)
	if err != nil {
		t.Errorf("cannot get Profile %s (#%d) by ID: %v - alerts: %+v", firstProfile.Name, respProfile.ID, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Profile to exist with ID %d, found: %d", respProfile.ID, len(resp.Response))
	}
	respProfile = resp.Response[0]
	respParameters := respProfile.Parameters
	if len(respParameters) == 0 {
		t.Error("expected a profile with parameters to be retrieved, recieved one without any parameters")
	}
}

func DeleteTestProfiles(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, pr := range testData.Profiles {
		// Retrieve the Profile by name so we can get the id for the Update
		opts.QueryParameters.Set("name", pr.Name)
		resp, _, err := TOSession.GetProfiles(opts)
		if err != nil {
			t.Errorf("cannot get Profile '%s' by name': %v - alerts: %+v", pr.Name, err, resp.Alerts)
		}
		if len(resp.Response) != 1 {
			t.Errorf("Expected exactly one Profile to exist with name '%s' found: %d", pr.Name, len(resp.Response))
			continue
		}
		profileID := resp.Response[0].ID

		// query by name does not retrieve associated parameters.  But query by id does.
		opts.QueryParameters.Del("name")
		opts.QueryParameters.Set("id", strconv.Itoa(profileID))
		resp, _, err = TOSession.GetProfiles(opts)
		opts.QueryParameters.Del("id")
		if err != nil {
			t.Errorf("cannot get Profile '%s' (#%d) by ID: %v - alerts: %+v", pr.Name, profileID, err, resp.Alerts)
		}
		if len(resp.Response) != 1 {
			t.Errorf("Expected exactly one Profile to exist with ID %d, found: %d", profileID, len(resp.Response))
			continue
		}
		// delete any profile_parameter associations first
		// the parameter is what's being deleted, but the delete is cascaded to profile_parameter
		for _, param := range resp.Response[0].Parameters {
			if param.ID == nil {
				t.Error("Traffic Ops responded with a representation of a Parameter with null or undefined ID")
				continue
			}
			alerts, _, err := TOSession.DeleteParameter(*param.ID, client.RequestOptions{})
			if err != nil {
				t.Errorf("cannot delete Parameter #%d: %v - alerts: %+v", *param.ID, err, alerts.Alerts)
			}
		}
		delResp, _, err := TOSession.DeleteProfile(profileID, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot delete Profile: %v - alerts: %+v", err, delResp.Alerts)
		}

		// Retrieve the Profile to see if it got deleted
		opts.QueryParameters.Set("name", pr.Name)
		prs, _, err := TOSession.GetProfiles(opts)
		if err != nil {
			t.Errorf("error fetching Profile after deletion: %v - alerts: %+v", err, prs.Alerts)
		}
		if len(prs.Response) > 0 {
			t.Errorf("expected Profile '%s' to be deleted, but it was found in Traffic Ops", pr.Name)
		}

		// Attempt to export Profile
		_, _, err = TOSession.ExportProfile(profileID, client.RequestOptions{})
		if err == nil {
			t.Errorf("export deleted profile %s - expected: error, actual: nil", pr.Name)
		}
	}
}

func GetTestPaginationSupportProfiles(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "id")
	resp, _, err := TOSession.GetProfiles(opts)
	if err != nil {
		t.Errorf("Unexpected error getting Profiles: %v - alerts: %+v", err, resp.Alerts)
	}
	profiles := resp.Response

	if len(profiles) > 0 {
		opts.QueryParameters = url.Values{}
		opts.QueryParameters.Set("orderby", "id")
		opts.QueryParameters.Set("limit", "1")
		profilesWithLimit, _, err := TOSession.GetProfiles(opts)
		if err == nil {
			if !reflect.DeepEqual(profiles[:1], profilesWithLimit.Response) {
				t.Error("expected GET Profiles with limit = 1 to return first result")
			}
		} else {
			t.Error("Error in getting Profiles by limit")
		}
		if len(profiles) > 1 {
			opts.QueryParameters = url.Values{}
			opts.QueryParameters.Set("orderby", "id")
			opts.QueryParameters.Set("limit", "1")
			opts.QueryParameters.Set("offset", "1")
			profilesWithOffset, _, err := TOSession.GetProfiles(opts)
			if err == nil {
				if !reflect.DeepEqual(profiles[1:2], profilesWithOffset.Response) {
					t.Error("expected GET Profiles with limit = 1, offset = 1 to return second result")
				}
			} else {
				t.Error("Error in getting Profiles by limit and offset")
			}

			opts.QueryParameters = url.Values{}
			opts.QueryParameters.Set("orderby", "id")
			opts.QueryParameters.Set("limit", "1")
			opts.QueryParameters.Set("page", "2")
			profilesWithPage, _, err := TOSession.GetProfiles(opts)
			if err == nil {
				if !reflect.DeepEqual(profiles[1:2], profilesWithPage.Response) {
					t.Error("expected GET Profiles with limit = 1, page = 2 to return second result")
				}
			} else {
				t.Error("Error in getting Profiles by limit and page")
			}
		} else {
			t.Errorf("only one Profiles found, so offset functionality can't test")
		}
	} else {
		t.Errorf("No Profiles found to check pagination")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "-2")
	resp, _, err = TOSession.GetProfiles(opts)
	if err == nil {
		t.Error("expected GET Profiles to return an error when limit is not bigger than -1")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be bigger than -1") {
		t.Errorf("expected getting Profiles where limit is not bigger than -1 to return an error stating so, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "0")
	resp, _, err = TOSession.GetProfiles(opts)
	if err == nil {
		t.Error("expected GET Profiles to return an error when offset is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected getting Profiles where offset is not a positive integer to return an error stating so, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "0")
	resp, _, err = TOSession.GetProfiles(opts)
	if err == nil {
		t.Error("expected GET Profiles to return an error when page is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected getting Profiles where page is not a positive integer to return an error stating so, actual error: %v - alerts: %+v", err, resp.Alerts)
	}
}
