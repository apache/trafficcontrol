package v3

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
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

var fedIDs = make(map[string]int)

// All prerequisite Federations are associated to this cdn and this xmlID
var cdnName = "cdn1"
var xmlId = "ds1"

func TestCDNFederations(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Tenants, CacheGroups, Statuses, Divisions, Regions, PhysLocations, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.CDNFederation]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID ID parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"id": {strconv.Itoa(GetFederationID(t, "the.cname.com.")())}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"SORTED by CNAME when ORDERBY=CNAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"orderby": {"cname"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCDNFederationCNameSort()),
				},
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetFederationID(t, "google.com."),
					ClientSession: TOSession,
					RequestBody: tc.CDNFederation{
						CName:       util.Ptr("new.cname."),
						TTL:         util.Ptr(34),
						Description: util.Ptr("updated"),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCDNFederationUpdateFields(map[string]interface{}{"CName": "new.cname."})),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:     GetFederationID(t, "booya.com."),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody: tc.CDNFederation{
						CName:       util.Ptr("booya.com."),
						TTL:         util.Ptr(34),
						Description: util.Ptr("fooya"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetFederationID(t, "booya.com."),
					ClientSession: TOSession,
					RequestBody: tc.CDNFederation{
						CName:       util.Ptr("new.cname."),
						TTL:         util.Ptr(34),
						Description: util.Ptr("updated"),
					},
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var fedID int
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							if name == "OK when VALID ID parameter" {
								if val, ok := testCase.RequestParams["id"]; ok {
									id, err := strconv.Atoi(val[0])
									assert.RequireNoError(t, err, "Failed to convert ID to an integer.")
									fedID = id
								}
								resp, reqInf, err := testCase.ClientSession.GetCDNFederationsByIDWithHdr(cdnName, fedID, testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else {
								resp, reqInf, err := testCase.ClientSession.GetCDNFederationsByNameWithHdr(cdnName, testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp.Response, resp.Alerts, err)
								}
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateCDNFederationByName(testCase.RequestBody, cdnName)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateCDNFederationsByIDWithHdr(testCase.RequestBody, cdnName, testCase.EndpointID(), testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteCDNFederationByID(cdnName, testCase.EndpointID())
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts.Alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func validateCDNFederationUpdateFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected CDN Federation response to not be nil.")
		CDNFederationResp := resp.(tc.CDNFederation)
		for field, expected := range expectedResp {
			switch field {
			case "CName":
				assert.RequireNotNil(t, CDNFederationResp.CName, "Expected CName to not be nil.")
				assert.Equal(t, expected, *CDNFederationResp.CName, "Expected CName to be %v, but got %s", expected, *CDNFederationResp.CName)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}

func validateCDNFederationCNameSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected CDN Federation response to not be nil.")
		var federationCNames []string
		CDNFederationResp := resp.([]tc.CDNFederation)
		for _, CDNFederation := range CDNFederationResp {
			assert.RequireNotNil(t, CDNFederation.CName, "Expected CDN Federation CName to not be nil.")
			federationCNames = append(federationCNames, *CDNFederation.CName)
		}
		assert.Equal(t, true, sort.StringsAreSorted(federationCNames), "List is not sorted by their names: %v", federationCNames)
	}
}

func GetFederationID(t *testing.T, cname string) func() int {
	return func() int {
		ID, ok := fedIDs[cname]
		assert.RequireEqual(t, true, ok, "Expected to find Federation CName: %s to have associated ID", cname)
		return ID
	}
}

func setFederationID(t *testing.T, cdnFederation tc.CDNFederation) {
	assert.RequireNotNil(t, cdnFederation.CName, "Federation CName was nil after posting.")
	assert.RequireNotNil(t, cdnFederation.ID, "Federation ID was nil after posting.")
	fedIDs[*cdnFederation.CName] = *cdnFederation.ID
}

func CreateTestCDNFederations(t *testing.T) {
	for _, federation := range testData.Federations {
		dsResp, _, err := TOSession.GetDeliveryServiceByXMLIDNullableWithHdr(*federation.DeliveryServiceIDs.XmlId, nil)
		assert.RequireNoError(t, err, "Could not get Delivery Service by XML ID: %v", err)
		assert.RequireEqual(t, 1, len(dsResp), "Expected one Delivery Service, but got %d", len(dsResp))
		assert.RequireNotNil(t, dsResp[0].CDNName, "Expected Delivery Service CDN Name to not be nil.")

		resp, _, err := TOSession.CreateCDNFederationByName(federation, *dsResp[0].CDNName)
		assert.NoError(t, err, "Could not create CDN Federations: %v - alerts: %+v", err, resp.Alerts)

		// Need to save the ids, otherwise the other tests won't be able to reference the federations
		setFederationID(t, resp.Response)
		assert.RequireNotNil(t, resp.Response.ID, "Federation ID was nil after posting.")
		assert.RequireNotNil(t, dsResp[0].ID, "Delivery Service ID was nil.")
		_, err = TOSession.CreateFederationDeliveryServices(*resp.Response.ID, []int{*dsResp[0].ID}, false)
		assert.NoError(t, err, "Could not create Federation Delivery Service: %v", err)
	}
}

func DeleteTestCDNFederations(t *testing.T) {
	for _, id := range fedIDs {
		resp, _, err := TOSession.DeleteCDNFederationByID(cdnName, id)
		assert.NoError(t, err, "Cannot delete federation #%d: %v - alerts: %+v", id, err, resp.Alerts)
	}
	data, _, _ := TOSession.GetCDNFederationsByNameWithHdr(cdnName, nil)
	assert.Equal(t, 0, len(data.Response), "expected federation to be deleted")
	fedIDs = make(map[string]int) // reset the global variable for the next test
}
