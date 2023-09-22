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
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	tc "github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestProfiles(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, ProfileParameters}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.Profile]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"RASCAL1"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateProfilesFields(map[string]interface{}{"Name": "RASCAL1"})),
				},
				"OK when VALID CDN parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"cdn": {strconv.Itoa(GetCDNID(t, "cdn1")())}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateProfilesFields(map[string]interface{}{"CDNName": "cdn1"})),
				},
				"OK when VALID ID parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"id": {strconv.Itoa(GetProfileID(t, "EDGEInCDN2")())}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateProfilesFields(map[string]interface{}{"ID": GetProfileID(t, "EDGEInCDN2")()})),
				},
			},
			"POST": {
				"BAD REQUEST when NAME has SPACES": {
					ClientSession: TOSession,
					RequestBody: tc.Profile{
						CDNID:       GetCDNID(t, "cdn1")(),
						Description: "name has spaces test",
						Name:        "name has space",
						Type:        tc.CacheServerProfileType,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING ALL FIELDS": {
					ClientSession: TOSession,
					RequestBody: tc.Profile{
						CDNID:       0,
						Description: "",
						Name:        "",
						Type:        "",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID CDN ID": {
					ClientSession: TOSession,
					RequestBody: tc.Profile{
						CDNID:       0,
						Description: "description",
						Name:        "badprofile",
						Type:        tc.CacheServerProfileType,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING DESCRIPTION FIELD": {
					ClientSession: TOSession,
					RequestBody: tc.Profile{
						CDNID:       GetCDNID(t, "cdn1")(),
						Description: "",
						Name:        "missing_description",
						Type:        tc.CacheServerProfileType,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING NAME FIELD": {
					ClientSession: TOSession,
					RequestBody: tc.Profile{
						CDNID:       GetCDNID(t, "cdn1")(),
						Description: "missing name",
						Name:        "",
						Type:        tc.CacheServerProfileType,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING TYPE FIELD": {
					ClientSession: TOSession,
					RequestBody: tc.Profile{
						CDNID:       GetCDNID(t, "cdn1")(),
						Description: "missing type",
						Name:        "missing_type",
						Type:        "",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"OK when VALID REQUEST": {
					EndpointID:    GetProfileID(t, "EDGE2"),
					ClientSession: TOSession,
					RequestBody: tc.Profile{
						CDNID:           GetCDNID(t, "cdn2")(),
						Description:     "edge2 description updated",
						Name:            "EDGE2UPDATED",
						RoutingDisabled: false,
						Type:            "TR_PROFILE",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateProfilesUpdateCreateFields("EDGE2UPDATED",
							map[string]interface{}{"CDNName": "cdn2", "Description": "edge2 description updated",
								"Name": "EDGE2UPDATED", "RoutingDisabled": false, "Type": "TR_PROFILE"})),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:     GetProfileID(t, "CCR1"),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody: tc.Profile{
						CDNID:           GetCDNID(t, "cdn1")(),
						Description:     "cdn1 description",
						Name:            "CCR1",
						RoutingDisabled: false,
						Type:            "TR_PROFILE",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetProfileID(t, "CCR1"),
					ClientSession: TOSession,
					RequestBody: tc.Profile{
						CDNID:           GetCDNID(t, "cdn1")(),
						Description:     "cdn1 description",
						Name:            "CCR1",
						RoutingDisabled: false,
						Type:            "TR_PROFILE",
					},
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							if name == "OK when VALID NAME parameter" {
								resp, reqInf, err := testCase.ClientSession.GetProfileByNameWithHdr(testCase.RequestParams["name"][0], testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else if name == "OK when VALID CDN parameter" {
								cdnID, err := strconv.Atoi(testCase.RequestParams["cdn"][0])
								resp, reqInf, err := testCase.ClientSession.GetProfileByCDNIDWithHdr(cdnID, testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else if name == "OK when VALID ID parameter" {
								id, err := strconv.Atoi(testCase.RequestParams["id"][0])
								resp, reqInf, err := testCase.ClientSession.GetProfileByIDWithHdr(id, testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else if name == "OK when VALID PARAM parameter" {
								paramID, err := strconv.Atoi(testCase.RequestParams["param"][0])
								resp, reqInf, err := testCase.ClientSession.GetProfileByParameterIdWithHdr(paramID, testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else {
								resp, reqInf, err := testCase.ClientSession.GetProfilesWithHdr(testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateProfile(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateProfileByIDWithHdr(testCase.EndpointID(), testCase.RequestBody, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteProfileByID(testCase.EndpointID())
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func validateProfilesFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Profiles response to not be nil.")
		profileResp := resp.([]tc.Profile)
		for field, expected := range expectedResp {
			for _, profile := range profileResp {
				switch field {
				case "CDNID":
					assert.Equal(t, expected, profile.CDNID, "Expected CDNID to be %v, but got %d", expected, profile.CDNID)
				case "CDNName":
					assert.Equal(t, expected, profile.CDNName, "Expected CDNName to be %v, but got %s", expected, profile.CDNName)
				case "Description":
					assert.Equal(t, expected, profile.Description, "Expected Description to be %v, but got %s", expected, profile.Description)
				case "ID":
					assert.Equal(t, expected, profile.ID, "Expected ID to be %v, but got %d", expected, profile.ID)
				case "Name":
					assert.Equal(t, expected, profile.Name, "Expected Name to be %v, but got %s", expected, profile.Name)
				case "Parameter":
					assert.Equal(t, expected, profile.Parameter, "Expected Parameter to be %v, but got %s", expected, profile.Parameter)
				case "Parameters":
					assert.Exactly(t, expected, profile.Parameters, "Expected Parameters to be %v, but got %s", expected, profile.Parameters)
				case "RoutingDisabled":
					assert.Equal(t, expected, profile.RoutingDisabled, "Expected RoutingDisabled to be %v, but got %v", expected, profile.RoutingDisabled)
				case "Type":
					assert.Equal(t, expected, profile.Type, "Expected Type to be %v, but got %s", expected, profile.Type)
				default:
					t.Fatalf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateProfilesUpdateCreateFields(name string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		profiles, _, err := TOSession.GetProfileByNameWithHdr(name, nil)
		assert.RequireNoError(t, err, "Error getting Profile: %v", err)
		assert.RequireEqual(t, 1, len(profiles), "Expected one Profile returned Got: %d", len(profiles))
		validateProfilesFields(expectedResp)(t, toclientlib.ReqInf{}, profiles, tc.Alerts{}, nil)
	}
}

func GetProfileID(t *testing.T, profileName string) func() int {
	return func() int {
		resp, _, err := TOSession.GetProfileByNameWithHdr(profileName, nil)
		assert.RequireNoError(t, err, "Get Profiles Request failed with error: %v", err)
		assert.RequireEqual(t, 1, len(resp), "Expected response object length 1, but got %d", len(resp))
		return resp[0].ID
	}
}

func CreateTestProfiles(t *testing.T) {
	for _, profile := range testData.Profiles {
		resp, _, err := TOSession.CreateProfile(profile)
		assert.RequireNoError(t, err, "Could not create Profile '%s': %v - alerts: %+v", profile.Name, err, resp.Alerts)
	}
}

func DeleteTestProfiles(t *testing.T) {
	profiles, _, err := TOSession.GetProfilesWithHdr(nil)
	assert.NoError(t, err, "Cannot get Profiles: %v", err)
	for _, profile := range profiles {
		alerts, _, err := TOSession.DeleteProfileByID(profile.ID)
		assert.NoError(t, err, "Cannot delete Profile: %v - alerts: %+v", err, alerts.Alerts)
		// Retrieve the Profile to see if it got deleted
		getProfiles, _, err := TOSession.GetProfileByIDWithHdr(profile.ID, nil)
		assert.NoError(t, err, "Error getting Profile '%s' after deletion: %v", profile.Name, err)
		assert.Equal(t, 0, len(getProfiles), "Expected Profile '%s' to be deleted, but it was found in Traffic Ops", profile.Name)
	}
}
