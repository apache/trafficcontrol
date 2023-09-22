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

package v5

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
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestProfiles(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, ProfileParameters}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.ProfileV5]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"RASCAL1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateProfilesFields(map[string]interface{}{"Name": "RASCAL1"})),
				},
				"OK when VALID PARAM parameter": {
					ClientSession: TOSession,
					RequestOpts: client.RequestOptions{QueryParameters: url.Values{
						"id":    {strconv.Itoa(GetProfileID(t, "EDGE1")())},
						"param": {strconv.Itoa(GetParameterID(t, "health.threshold.loadavg", "rascal.properties", "25.0")())},
					}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateProfilesFields(map[string]interface{}{"Parameter": "health.threshold.loadavg"})),
				},
				"OK when VALID CDN parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {strconv.Itoa(GetCDNID(t, "cdn1")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateProfilesFields(map[string]interface{}{"CDNName": "cdn1"})),
				},
				"OK when VALID ID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"id": {strconv.Itoa(GetProfileID(t, "EDGEInCDN2")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateProfilesFields(map[string]interface{}{"ID": GetProfileID(t, "EDGEInCDN2")()})),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateProfilesPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateProfilesPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateProfilesPagination("page")),
				},
				"BAD REQUEST when INVALID LIMIT parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"-2"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID OFFSET parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "offset": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PAGE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "page": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"POST": {
				"BAD REQUEST when NAME has SPACES": {
					ClientSession: TOSession,
					RequestBody: tc.ProfileV5{
						CDNID:       GetCDNID(t, "cdn1")(),
						Description: "name has spaces test",
						Name:        "name has space",
						Type:        tc.CacheServerProfileType,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING ALL FIELDS": {
					ClientSession: TOSession,
					RequestBody:   tc.ProfileV5{},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID CDN ID": {
					ClientSession: TOSession,
					RequestBody: tc.ProfileV5{
						CDNID:       0,
						Description: "description",
						Name:        "badprofile",
						Type:        tc.CacheServerProfileType,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING DESCRIPTION FIELD": {
					ClientSession: TOSession,
					RequestBody: tc.ProfileV5{
						CDNID: GetCDNID(t, "cdn1")(),
						Name:  "missing_description",
						Type:  tc.CacheServerProfileType,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING NAME FIELD": {
					ClientSession: TOSession,
					RequestBody: tc.ProfileV5{
						CDNID:       GetCDNID(t, "cdn1")(),
						Description: "missing name",
						Type:        tc.CacheServerProfileType,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING TYPE FIELD": {
					ClientSession: TOSession,
					RequestBody: tc.ProfileV5{
						CDNID:       GetCDNID(t, "cdn1")(),
						Description: "missing type",
						Name:        "missing_type",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"OK when VALID REQUEST": {
					EndpointID:    GetProfileID(t, "EDGE2"),
					ClientSession: TOSession,
					RequestBody: tc.ProfileV5{
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
					EndpointID:    GetProfileID(t, "CCR1"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: tc.ProfileV5{
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
					RequestBody: tc.ProfileV5{
						CDNID:           GetCDNID(t, "cdn1")(),
						Description:     "cdn1 description",
						Name:            "CCR1",
						RoutingDisabled: false,
						Type:            "TR_PROFILE",
					},
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetProfiles(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateProfile(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateProfile(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteProfile(testCase.EndpointID(), testCase.RequestOpts)
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
		profileResp := resp.([]tc.ProfileV5)
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
					assert.Equal(t, true, validateProfileContainsParameter(t, expected.(string), profile.Parameters), "Expected to find Parameter in Profiles Parameters list.")
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
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		profiles, _, err := TOSession.GetProfiles(opts)
		assert.RequireNoError(t, err, "Error getting Profile: %v - alerts: %+v", err, profiles.Alerts)
		assert.RequireEqual(t, 1, len(profiles.Response), "Expected one Profile returned Got: %d", len(profiles.Response))
		validateProfilesFields(expectedResp)(t, toclientlib.ReqInf{}, profiles.Response, tc.Alerts{}, nil)
	}
}

func validateProfilesPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.ProfileV5)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetProfiles(opts)
		assert.RequireNoError(t, err, "Cannot get Profiles: %v - alerts: %+v", err, respBase.Alerts)

		profiles := respBase.Response
		assert.RequireGreaterOrEqual(t, len(profiles), 3, "Need at least 3 Profiles in Traffic Ops to test pagination support, found: %d", len(profiles))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, profiles[:1], paginationResp, "expected GET Profiles with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, profiles[1:2], paginationResp, "expected GET Profiles with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, profiles[1:2], paginationResp, "expected GET Profiles with limit = 1, page = 2 to return second result")
		}
	}
}

func validateProfileContainsParameter(t *testing.T, expectedParameter string, actualParameters []tc.ParameterNullable) bool {
	for _, parameter := range actualParameters {
		assert.RequireNotNil(t, parameter.Name, "Expected Parameter Name to not be nil.")
		if expectedParameter == *parameter.Name {
			return true
		}
	}
	return false
}

func GetProfileID(t *testing.T, profileName string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", profileName)
		resp, _, err := TOSession.GetProfiles(opts)
		assert.RequireNoError(t, err, "Get Profiles Request failed with error: %v", err)
		assert.RequireEqual(t, 1, len(resp.Response), "Expected response object length 1, but got %d", len(resp.Response))
		return resp.Response[0].ID
	}
}

func CreateTestProfiles(t *testing.T) {
	for _, profile := range testData.Profiles {
		resp, _, err := TOSession.CreateProfile(profile, client.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Profile '%s': %v - alerts: %+v", profile.Name, err, resp.Alerts)
	}
}

func DeleteTestProfiles(t *testing.T) {
	profiles, _, err := TOSession.GetProfiles(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get Profiles: %v - alerts: %+v", err, profiles.Alerts)
	for _, profile := range profiles.Response {
		alerts, _, err := TOSession.DeleteProfile(profile.ID, client.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Profile: %v - alerts: %+v", err, alerts.Alerts)
		// Retrieve the Profile to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(profile.ID))
		getProfiles, _, err := TOSession.GetProfiles(opts)
		assert.NoError(t, err, "Error getting Profile '%s' after deletion: %v - alerts: %+v", profile.Name, err, getProfiles.Alerts)
		assert.Equal(t, 0, len(getProfiles.Response), "Expected Profile '%s' to be deleted, but it was found in Traffic Ops", profile.Name)
	}
}
