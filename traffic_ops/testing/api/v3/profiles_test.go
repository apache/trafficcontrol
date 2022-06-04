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
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
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
	})
}

func UpdateTestProfilesWithHeaders(t *testing.T, header http.Header) {
	if len(testData.Profiles) > 0 {
		firstProfile := testData.Profiles[0]
		// Retrieve the Profile by name so we can get the id for the Update
		resp, _, err := TOSession.GetProfileByNameWithHdr(firstProfile.Name, header)
		if err != nil {
			t.Errorf("cannot GET Profile by name: %v - %v", firstProfile.Name, err)
		}
		if len(resp) > 0 {
			remoteProfile := resp[0]
			_, reqInf, err := TOSession.UpdateProfileByIDWithHdr(remoteProfile.ID, remoteProfile, header)
			if err == nil {
				t.Errorf("Expected error about precondition failed, but got none")
			}
			if reqInf.StatusCode != http.StatusPreconditionFailed {
				t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
			}
		}
	}
}

func GetTestProfilesIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	for _, pr := range testData.Profiles {
		_, reqInf, err := TOSession.GetProfileByNameWithHdr(pr.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
		if len(pr.Parameters) > 0 {
			parameter := pr.Parameters[0]
			respParameter, _, err := TOSession.GetParameterByName(*parameter.Name)
			if err != nil {
				t.Errorf("Cannot GET Parameter by name: %v", err)
			}
			if len(respParameter) > 0 {
				parameterID := respParameter[0].ID
				if parameterID > 0 {
					resp, _, err := TOSession.GetProfileByParameterIdWithHdr(parameterID, nil)
					if err != nil {
						t.Fatalf("Expected no error, but got %v", err.Error())
					}
					if reqInf.StatusCode != http.StatusNotModified {
						t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
					}
					if len(resp) < 1 {
						t.Errorf("Expected atleast one response for Get Profile by Parameters, but found %d", len(resp))
					}
				} else {
					t.Errorf("Invalid parameter ID %d", parameterID)
				}
			} else {
				t.Errorf("No response found for GET Parameters by name")
			}
		}
	}
}

// CreateBadProfiles ensures that profiles can't be created with bad values
func CreateBadProfiles(t *testing.T) {

	// blank profile
	prs := []tc.Profile{
		{Type: "", Name: "", Description: "", CDNID: 0},
		{Type: tc.CacheServerProfileType, Name: "badprofile", Description: "description", CDNID: 0},
		{Type: tc.CacheServerProfileType, Name: "badprofile", Description: "", CDNID: 1},
		{Type: tc.CacheServerProfileType, Name: "", Description: "description", CDNID: 1},
		{Type: "", Name: "badprofile", Description: "description", CDNID: 1},
	}

	for _, pr := range prs {
		resp, _, err := TOSession.CreateProfile(pr)

		if err == nil {
			t.Errorf("Creating bad profile succeeded: %+v\nResponse is %+v", pr, resp)
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
			resp, _, err := TOSession.CopyProfile(c.profile)
			if c.err != "" {
				if err != nil && !strings.Contains(err.Error(), c.err) {
					t.Fatalf("got err= %s; expected err= %s", err, c.err)
				}
			} else if err != nil {
				t.Fatalf("got err= %s; expected err= nil", err)
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
	for _, name := range newProfileNames {
		profiles, _, err := TOSession.GetProfileByName(name)
		if err != nil {
			t.Fatalf("got err= %s; expected err= nil", err)
		}
		if len(profiles) == 0 {
			t.Errorf("could not GET profile %+v: not found", name)
		}
		_, _, err = TOSession.DeleteProfileByID(profiles[0].ID)
		if err != nil {
			t.Fatalf("got err= %s; expected err= nil", err)
		}
	}
}

func CreateTestProfiles(t *testing.T) {

	for _, pr := range testData.Profiles {
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

func UpdateTestProfiles(t *testing.T) {

	firstProfile := testData.Profiles[0]
	// Retrieve the Profile by name so we can get the id for the Update
	resp, _, err := TOSession.GetProfileByName(firstProfile.Name)
	if err != nil {
		t.Errorf("cannot GET Profile by name: %v - %v", firstProfile.Name, err)
	}
	remoteProfile := resp[0]
	expectedProfileDesc := "UPDATED"
	remoteProfile.Description = expectedProfileDesc
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateProfileByID(remoteProfile.ID, remoteProfile)
	if err != nil {
		t.Errorf("cannot UPDATE Profile by id: %v - %v", err, alert)
	}

	// Retrieve the Profile to check Profile name got updated
	resp, _, err = TOSession.GetProfileByID(remoteProfile.ID)
	if err != nil {
		t.Errorf("cannot GET Profile by name: %v - %v", firstProfile.Name, err)
	}
	respProfile := resp[0]
	if respProfile.Description != expectedProfileDesc {
		t.Errorf("results do not match actual: %s, expected: %s", respProfile.Description, expectedProfileDesc)
	}

}

func GetTestProfiles(t *testing.T) {

	for _, pr := range testData.Profiles {
		resp, _, err := TOSession.GetProfileByName(pr.Name)
		if err != nil {
			t.Errorf("cannot GET Profile by name: %v - %v", err, resp)
		}
		profileID := resp[0].ID
		if len(pr.Parameters) > 0 {
			parameter := pr.Parameters[0]
			respParameter, _, err := TOSession.GetParameterByName(*parameter.Name)
			if err != nil {
				t.Errorf("Cannot GET Parameter by name: %v", err)
			}
			if len(respParameter) > 0 {
				parameterID := respParameter[0].ID
				if parameterID > 0 {
					resp, _, err = TOSession.GetProfileByParameterId(parameterID)
					if err != nil {
						t.Errorf("cannot GET Profile by param: %v - %v", err, resp)
					}
					if len(resp) < 1 {
						t.Errorf("Expected atleast one response for Get Profile by Parameters, but found %d", len(resp))
					}
				} else {
					t.Errorf("Invalid parameter ID %d", parameterID)
				}
			} else {
				t.Errorf("No response found for GET Parameters by name")
			}
		}

		resp, _, err = TOSession.GetProfileByCDNID(pr.CDNID)
		if err != nil {
			t.Errorf("cannot GET Profile by cdn: %v - %v", err, resp)
		}

		// Export Profile
		exportResp, _, err := TOSession.ExportProfile(profileID)
		if err != nil {
			t.Errorf("error exporting Profile: %v - %v", profileID, err)
		}
		if exportResp == nil {
			t.Error("error exporting Profile: response nil")
		}
	}
}

func ImportProfile(t *testing.T) {
	// Get ID of Profile to export
	resp, _, err := TOSession.GetProfileByName(testData.Profiles[0].Name)
	if err != nil {
		t.Fatalf("cannot GET Profile by name: %v - %v", err, resp)
	}
	if resp == nil {
		t.Fatal("error getting Profile: response nil")
	}
	if len(resp) != 1 {
		t.Fatalf("Profiles expected 1, actual %v", len(resp))
	}
	profileID := resp[0].ID

	// Export Profile to import
	exportResp, _, err := TOSession.ExportProfile(profileID)
	if err != nil {
		t.Fatalf("error exporting Profile: %v - %v", profileID, err)
	}
	if exportResp == nil {
		t.Fatal("error exporting Profile: response nil")
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
	importResp, _, err := TOSession.ImportProfile(&importReq)
	if err != nil {
		t.Fatalf("error importing Profile: %v - %v", profileID, err)
	}
	if importResp == nil {
		t.Error("error importing Profile: response nil")
	}

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
	firstProfile := testData.Profiles[0]
	resp, _, err := TOSession.GetProfileByName(firstProfile.Name)
	if err != nil {
		t.Errorf("cannot GET Profile by name: %v - %v", err, resp)
		return
	}
	if len(resp) == 0 {
		t.Errorf("cannot GET Profile by name: not found - %v", resp)
		return
	}
	respProfile := resp[0]
	// query by name does not retrieve associated parameters.  But query by id does.
	resp, _, err = TOSession.GetProfileByID(respProfile.ID)
	if err != nil {
		t.Errorf("cannot GET Profile by name: %v - %v", err, resp)
	}
	if len(resp) > 0 {
		respProfile = resp[0]
		respParameters := respProfile.Parameters
		if len(respParameters) == 0 {
			t.Errorf("expected a profile with parameters to be retrieved: %v - %v", err, respParameters)
		}
	}
}

func DeleteTestProfiles(t *testing.T) {

	for _, pr := range testData.Profiles {
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
			t.Errorf("export deleted profile %s - expected: error, actual: nil", pr.Name)
		}
	}
}
