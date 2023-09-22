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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestCacheGroups(t *testing.T) {
	WithObjs(t, []TCObj{Types, Parameters, CacheGroups, CDNs, Profiles, Statuses, Divisions, Regions, PhysLocations, Servers, Topologies}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.CacheGroupNullable]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"NOT MODIFIED when VALID NAME parameter when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestParams:  url.Values{"name": {"originCachegroup"}},
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"NOT MODIFIED when VALID SHORTNAME parameter when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestParams:  url.Values{"shortName": {"mog1"}},
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"parentCachegroup"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						ValidateExpectedField("Name", "parentCachegroup")),
				},
				"OK when VALID SHORTNAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"shortName": {"pg2"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						ValidateExpectedField("ShortName", "pg2")),
				},
				"OK when VALID TOPOLOGY parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"topology": {"mso-topology"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"UNAUTHORIZED when NOT LOGGED IN": {
					ClientSession: NoAuthTOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusUnauthorized)),
				},
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"UNAUTHORIZED when NOT LOGGED IN": {
					ClientSession: NoAuthTOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusUnauthorized)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetCacheGroupId(t, "cachegroup1"),
					ClientSession: TOSession,
					RequestBody: tc.CacheGroupNullable{
						Latitude:            util.Ptr(17.5),
						Longitude:           util.Ptr(17.5),
						Name:                util.Ptr("cachegroup1"),
						ShortName:           util.Ptr("newShortName"),
						LocalizationMethods: util.Ptr([]tc.LocalizationMethod{tc.LocalizationMethodCZ}),
						Fallbacks:           util.Ptr([]string{"fallback1"}),
						Type:                util.Ptr("EDGE_LOC"),
						TypeID:              util.Ptr(GetTypeId(t, "EDGE_LOC")),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID: GetCacheGroupId(t, "cachegroup1"), ClientSession: TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody: tc.CacheGroupNullable{
						Name:      util.Ptr("cachegroup1"),
						ShortName: util.Ptr("changeName"),
						Type:      util.Ptr("EDGE_LOC"),
						TypeID:    util.Ptr(GetTypeId(t, "EDGE_LOC")),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetCacheGroupId(t, "cachegroup1"),
					ClientSession: TOSession,
					RequestBody: tc.CacheGroupNullable{
						Name:      util.Ptr("cachegroup1"),
						ShortName: util.Ptr("changeName"),
						Type:      util.Ptr("EDGE_LOC"),
						TypeID:    util.Ptr(GetTypeId(t, "EDGE_LOC")),
					},
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"UNAUTHORIZED when NOT LOGGED IN": {
					EndpointID:    GetCacheGroupId(t, "cachegroup1"),
					ClientSession: NoAuthTOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusUnauthorized)),
				},
			},
			"DELETE": {
				"UNAUTHORIZED when NOT LOGGED IN": {
					EndpointID:    GetCacheGroupId(t, "cachegroup1"),
					ClientSession: NoAuthTOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusUnauthorized)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetCacheGroupsByQueryParamsWithHdr(testCase.RequestParams, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateCacheGroupNullable(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateCacheGroupNullableByIDWithHdr(testCase.EndpointID(), testCase.RequestBody, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteCacheGroupByID(testCase.EndpointID())
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
