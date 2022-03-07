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
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/traffic_ops/v3-client"
)

func TestCacheGroups(t *testing.T) {
	WithObjs(t, []TCObj{Types, Parameters, CacheGroups, CDNs, Profiles, Statuses, Divisions, Regions, PhysLocations, Servers, Topologies}, func() {

		tomorrow := time.Now().AddDate(0, 0, 1).Format(time.RFC1123)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)

		methodTests := map[string]map[string]struct {
			endpointId     func() int
			clientSession  *client.Session
			requestParams  url.Values
			requestHeaders http.Header
			requestBody    map[string]interface{}
			expectations   []utils.CkReqFunc
		}{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					clientSession: TOSession, requestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"NOT MODIFIED when VALID NAME parameter when NO CHANGES made": {
					clientSession: TOSession, requestParams: url.Values{"name": {"originCachegroup"}},
					requestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"NOT MODIFIED when VALID SHORTNAME parameter when NO CHANGES made": {
					clientSession: TOSession, requestParams: url.Values{"shortName": {"mog1"}},
					requestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					clientSession: TOSession, expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID NAME parameter": {
					clientSession: TOSession, requestParams: url.Values{"name": {"parentCachegroup"}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						ValidateExpectedField("Name", "parentCachegroup")),
				},
				"OK when VALID SHORTNAME parameter": {
					clientSession: TOSession, requestParams: url.Values{"shortName": {"pg2"}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						ValidateExpectedField("ShortName", "pg2")),
				},
				"OK when VALID TOPOLOGY parameter": {
					clientSession: TOSession, requestParams: url.Values{"topology": {"mso-topology"}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"UNAUTHORIZED when NOT LOGGED IN": {
					clientSession: NoAuthTOSession, expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusUnauthorized)),
				},
			},
			"POST": {
				"UNAUTHORIZED when NOT LOGGED IN": {
					clientSession: NoAuthTOSession, expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusUnauthorized)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					endpointId: GetCacheGroupId(t, "cachegroup1"), clientSession: TOSession,
					requestBody: map[string]interface{}{
						"latitude":            17.5,
						"longitude":           17.5,
						"name":                "cachegroup1",
						"shortName":           "newShortName",
						"localizationMethods": []string{"CZ"},
						"fallbacks":           []string{"fallback1"},
						"typeName":            "EDGE_LOC",
						"typeId":              -1,
					},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					endpointId: GetCacheGroupId(t, "parentCachegroup"), clientSession: TOSession,
					requestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}, rfc.IfUnmodifiedSince: {currentTimeRFC}},
					requestBody: map[string]interface{}{
						"latitude":  0,
						"longitude": 0,
						"name":      "parentCachegroup",
						"shortName": "pg1",
						"typeName":  "MID_LOC",
						"typeId":    -1,
					},
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					endpointId: GetCacheGroupId(t, "parentCachegroup2"), clientSession: TOSession,
					requestBody: map[string]interface{}{
						"latitude":  0,
						"longitude": 0,
						"name":      "parentCachegroup2",
						"shortName": "pg2",
						"typeName":  "MID_LOC",
						"typeId":    -1,
					},
					requestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"UNAUTHORIZED when NOT LOGGED IN": {
					endpointId: GetCacheGroupId(t, "cachegroup1"), clientSession: NoAuthTOSession,
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusUnauthorized)),
				},
			},
			"DELETE": {
				"UNAUTHORIZED when NOT LOGGED IN": {
					endpointId: GetCacheGroupId(t, "cachegroup1"), clientSession: NoAuthTOSession,
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusUnauthorized)),
				},
			},
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					clientSession:  TOSession,
					requestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}, rfc.IfUnmodifiedSince: {currentTimeRFC}},
					expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					cg := tc.CacheGroupNullable{}

					if testCase.requestParams.Has("type") {
						val := testCase.requestParams.Get("type")
						if _, err := strconv.Atoi(val); err != nil {
							testCase.requestParams.Set("type", strconv.Itoa(GetTypeId(t, val)))
						}
					}

					if testCase.requestBody != nil {
						if _, ok := testCase.requestBody["id"]; ok {
							testCase.requestBody["id"] = testCase.endpointId()
						}
						if typeId, ok := testCase.requestBody["typeId"]; ok {
							if typeId == -1 {
								if typeName, ok := testCase.requestBody["typeName"]; ok {
									testCase.requestBody["typeId"] = GetTypeId(t, typeName.(string))
								}
							}
						}
						dat, err := json.Marshal(testCase.requestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &cg)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.clientSession.GetCacheGroupsByQueryParamsWithHdr(testCase.requestParams, testCase.requestHeaders)
							for _, check := range testCase.expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.clientSession.CreateCacheGroupNullable(cg)
							for _, check := range testCase.expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.clientSession.UpdateCacheGroupNullableByIDWithHdr(testCase.endpointId(), cg, testCase.requestHeaders)
							for _, check := range testCase.expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.clientSession.DeleteCacheGroupByID(testCase.endpointId())
							for _, check := range testCase.expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func ValidateExpectedField(field string, expected string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		cgResp := resp.([]tc.CacheGroupNullable)
		cg := cgResp[0]
		switch field {
		case "Name":
			assert.Equal(t, expected, *cg.Name, "Expected name to be %v, but got %v", expected, *cg.Name)
		case "ShortName":
			assert.Equal(t, expected, *cg.ShortName, "Expected shortName to be %v, but got %v", expected, *cg.ShortName)
		case "TypeName":
			assert.Equal(t, expected, *cg.Type, "Expected type to be %v, but got %v", expected, *cg.Type)
		default:
			t.Errorf("Expected field: %v, does not exist in response", field)
		}
	}
}

func GetTypeId(t *testing.T, typeName string) int {
	resp, _, err := TOSession.GetTypeByNameWithHdr(typeName, nil)
	assert.RequireNoError(t, err, "Get Types Request failed with error: %v", err)
	assert.RequireEqual(t, 1, len(resp), "Expected response object length 1, but got %d", len(resp))
	assert.RequireNotNil(t, &resp[0].ID, "Expected id to not be nil")

	return resp[0].ID
}

func GetCacheGroupId(t *testing.T, cacheGroupName string) func() int {
	return func() int {
		resp, _, err := TOSession.GetCacheGroupNullableByNameWithHdr(cacheGroupName, nil)
		assert.RequireNoError(t, err, "Get Cache Groups Request failed with error: %v", err)
		assert.RequireEqual(t, len(resp), 1, "Expected response object length 1, but got %d", len(resp))
		assert.RequireNotNil(t, resp[0].ID, "Expected id to not be nil")

		return *resp[0].ID
	}
}

func CreateTestCacheGroups(t *testing.T) {
	var err error
	var resp *tc.CacheGroupDetailResponse

	for _, cg := range testData.CacheGroups {

		resp, _, err = TOSession.CreateCacheGroupNullable(cg)
		if err != nil {
			t.Errorf("could not CREATE cachegroups: %v, request: %v", err, cg)
			continue
		}

		// Testing 'join' fields during create
		if cg.ParentName != nil && resp.Response.ParentName == nil {
			t.Error("Parent cachegroup is null in response when it should have a value")
		}
		if cg.SecondaryParentName != nil && resp.Response.SecondaryParentName == nil {
			t.Error("Secondary parent cachegroup is null in response when it should have a value\n")
		}
		if cg.Type != nil && resp.Response.Type == nil {
			t.Error("Type is null in response when it should have a value\n")
		}
		assert.NotNil(t, resp.Response.LocalizationMethods, "Localization methods are null")
		assert.NotNil(t, resp.Response.Fallbacks, "Fallbacks are null")

	}
}

func DeleteTestCacheGroups(t *testing.T) {
	var parentlessCacheGroups []tc.CacheGroupNullable

	// delete the edge caches.
	for _, cg := range testData.CacheGroups {
		// Retrieve the CacheGroup by name so we can get the id for the Update
		resp, _, err := TOSession.GetCacheGroupNullableByNameWithHdr(*cg.Name, nil)
		assert.NoError(t, err, "Cannot GET CacheGroup by name '%s': %v", *cg.Name, err)
		cg = resp[0]

		// Cachegroups that are parents (usually mids but sometimes edges)
		// need to be deleted only after the children cachegroups are deleted.
		if cg.ParentCachegroupID == nil && cg.SecondaryParentCachegroupID == nil {
			parentlessCacheGroups = append(parentlessCacheGroups, cg)
			continue
		}
		if len(resp) > 0 {
			respCG := resp[0]
			_, _, err := TOSession.DeleteCacheGroupByID(*respCG.ID)
			assert.NoError(t, err, "Cannot delete Cache Group: %v - alerts: %+v", *respCG.Name, err)

			// Retrieve the CacheGroup to see if it got deleted
			cgs, _, err := TOSession.GetCacheGroupNullableByNameWithHdr(*cg.Name, nil)
			assert.NoError(t, err, "Error deleting Cache Group by name: %v", err)
			assert.Equal(t, 0, len(cgs), "Expected CacheGroup name: %s to be deleted", *cg.Name)
		}
	}

	// now delete the parentless cachegroups
	for _, cg := range parentlessCacheGroups {
		// Retrieve the CacheGroup by name so we can get the id for the Update
		resp, _, err := TOSession.GetCacheGroupNullableByNameWithHdr(*cg.Name, nil)
		assert.NoError(t, err, "Cannot GET CacheGroup by name '%s': %v", *cg.Name, err)
		if len(resp) > 0 {
			respCG := resp[0]
			_, _, err := TOSession.DeleteCacheGroupByID(*respCG.ID)
			assert.NoError(t, err, "Cannot delete Cache Group: %v - alerts: %+v", *respCG.Name, err)

			// Retrieve the CacheGroup to see if it got deleted
			cgs, _, err := TOSession.GetCacheGroupNullableByShortNameWithHdr(*cg.Name, nil)
			assert.NoError(t, err, "Error deleting Cache Group by name: %v", err)
			assert.Equal(t, 0, len(cgs), "Expected CacheGroup name: %s to be deleted", *cg.Name)
		}
	}
}
