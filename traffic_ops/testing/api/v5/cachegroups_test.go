package v5

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
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestCacheGroups(t *testing.T) {
	WithObjs(t, []TCObj{Types, Parameters, CacheGroups, CDNs, Profiles, Statuses, Divisions, Regions, PhysLocations, Servers, Topologies}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.CacheGroupNullableV5]{
			"GET": {
				"OK when VALID NAME parameter AND Lat/Long are 0": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"nullLatLongCG"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1), ValidateResponseFields()),
				},
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"NOT MODIFIED when VALID NAME parameter when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts: client.RequestOptions{
						Header:          http.Header{rfc.IfModifiedSince: {tomorrow}},
						QueryParameters: url.Values{"name": {"originCachegroup"}},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"NOT MODIFIED when VALID SHORTNAME parameter when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts: client.RequestOptions{
						Header:          http.Header{rfc.IfModifiedSince: {tomorrow}},
						QueryParameters: url.Values{"shortName": {"mog1"}},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"parentCachegroup"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						ValidateExpectedField("Name", "parentCachegroup")),
				},
				"OK when VALID SHORTNAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"shortName": {"pg2"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						ValidateExpectedField("ShortName", "pg2")),
				},
				"OK when VALID TOPOLOGY parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"topology": {"mso-topology"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID TYPE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"type": {strconv.Itoa(GetTypeId(t, "ORG_LOC"))}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						ValidateExpectedField("TypeName", "ORG_LOC")),
				},
				"EMPTY RESPONSE when INVALID ID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"id": {"10000"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID TYPE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"type": {"10000"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), ValidatePagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), ValidatePagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), ValidatePagination("page")),
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
				"UNAUTHORIZED when NOT LOGGED IN": {
					ClientSession: NoAuthTOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusUnauthorized)),
				},
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
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
					EndpointID: GetCacheGroupId(t, "cachegroup1"), ClientSession: TOSession,
					RequestBody: tc.CacheGroupNullableV5{
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
				"OK when updating CG with null Lat/Long": {
					EndpointID: GetCacheGroupId(t, "nullLatLongCG"), ClientSession: TOSession,
					RequestBody: tc.CacheGroupNullableV5{
						Name:      util.Ptr("nullLatLongCG"),
						ShortName: util.Ptr("null-ll"),
						Type:      util.Ptr("EDGE_LOC"),
						Fallbacks: util.Ptr([]string{"fallback1"}),
						TypeID:    util.Ptr(GetTypeId(t, "EDGE_LOC")),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when updating TYPE of CG in TOPOLOGY": {
					EndpointID: GetCacheGroupId(t, "topology-edge-cg-01"), ClientSession: TOSession,
					RequestBody: tc.CacheGroupNullableV5{
						Latitude:  util.Ptr(0.0),
						Longitude: util.Ptr(0.0),
						Name:      util.Ptr("topology-edge-cg-01"),
						ShortName: util.Ptr("te1"),
						Type:      util.Ptr("MID_LOC"),
						TypeID:    util.Ptr(GetTypeId(t, "MID_LOC")),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID: GetCacheGroupId(t, "cachegroup1"), ClientSession: TOSession,
					RequestOpts: client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: tc.CacheGroupNullableV5{
						Name:      util.Ptr("cachegroup1"),
						ShortName: util.Ptr("changeName"),
						Type:      util.Ptr("EDGE_LOC"),
						TypeID:    util.Ptr(GetTypeId(t, "EDGE_LOC")),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID: GetCacheGroupId(t, "cachegroup1"), ClientSession: TOSession,
					RequestOpts: client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					RequestBody: tc.CacheGroupNullableV5{
						Name:      util.Ptr("cachegroup1"),
						ShortName: util.Ptr("changeName"),
						Type:      util.Ptr("EDGE_LOC"),
						TypeID:    util.Ptr(GetTypeId(t, "EDGE_LOC")),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"UNAUTHORIZED when NOT LOGGED IN": {
					EndpointID:    GetCacheGroupId(t, "cachegroup1"),
					ClientSession: NoAuthTOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusUnauthorized)),
				},
			},
			"DELETE": {
				"NOT FOUND when INVALID ID parameter": {
					EndpointID:    func() int { return 111111 },
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
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
							resp, reqInf, err := testCase.ClientSession.GetCacheGroups(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateCacheGroup(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateCacheGroup(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteCacheGroup(testCase.EndpointID(), testCase.RequestOpts)
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
		cgResp := resp.([]tc.CacheGroupNullableV5)
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

func ValidateResponseFields() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		cgResp := resp.([]tc.CacheGroupNullableV5)
		cg := cgResp[0]
		assert.NotNil(t, cg.ID, "Expected response id to not be nil")
		assert.NotNil(t, cg.Latitude, "Expected latitude to not be nil")
		assert.NotNil(t, cg.Longitude, "Expected longitude to not be nil")
		assert.Equal(t, 0.0, *cg.Longitude, "Expected Longitude to be 0, but got %v", cg.Longitude)
		assert.Equal(t, 0.0, *cg.Latitude, "Expected Latitude to be 0, but got %v", cg.Latitude)
	}
}

func ValidatePagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.CacheGroupNullableV5)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetCacheGroups(opts)
		assert.RequireNoError(t, err, "cannot get Cache Groups: %v - alerts: %+v", err, respBase.Alerts)

		cachegroup := respBase.Response
		assert.RequireGreaterOrEqual(t, len(cachegroup), 3, "Need at least 3 Cache Groups in Traffic Ops to test pagination support, found: %d", len(cachegroup))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, cachegroup[:1], paginationResp, "expected GET Cachegroups with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, cachegroup[1:2], paginationResp, "expected GET cachegroup with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, cachegroup[1:2], paginationResp, "expected GET cachegroup with limit = 1, page = 2 to return second result")
		}
	}
}

func GetTypeId(t *testing.T, typeName string) int {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", typeName)
	resp, _, err := TOSession.GetTypes(opts)

	assert.RequireNoError(t, err, "Get Types Request failed with error: %v", err)
	assert.RequireEqual(t, 1, len(resp.Response), "Expected response object length 1, but got %d", len(resp.Response))
	assert.RequireNotNil(t, &resp.Response[0].ID, "Expected id to not be nil")

	return resp.Response[0].ID
}

func GetCacheGroupId(t *testing.T, cacheGroupName string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", cacheGroupName)

		resp, _, err := TOSession.GetCacheGroups(opts)
		assert.RequireNoError(t, err, "Get Cache Groups Request failed with error: %v", err)
		assert.RequireEqual(t, len(resp.Response), 1, "Expected response object length 1, but got %d", len(resp.Response))
		assert.RequireNotNil(t, resp.Response[0].ID, "Expected id to not be nil")

		return *resp.Response[0].ID
	}
}

func CreateTestCacheGroups(t *testing.T) {
	for _, cg := range testData.CacheGroups {

		resp, _, err := TOSession.CreateCacheGroup(cg, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create Cache Group: %v - alerts: %+v", err, resp.Alerts)
			continue
		}

		// Testing 'join' fields during create
		if cg.ParentName != nil && resp.Response.ParentName == nil {
			t.Error("Parent cachegroup is null in response when it should have a value")
		}
		if cg.SecondaryParentName != nil && resp.Response.SecondaryParentName == nil {
			t.Error("Secondary parent cachegroup is null in response when it should have a value")
		}
		if cg.Type != nil && resp.Response.Type == nil {
			t.Error("Type is null in response when it should have a value")
		}
		assert.NotNil(t, resp.Response.LocalizationMethods, "Localization methods are null")
		assert.NotNil(t, resp.Response.Fallbacks, "Fallbacks are null")
	}
}

func DeleteTestCacheGroups(t *testing.T) {
	var parentlessCacheGroups []tc.CacheGroupNullableV5
	opts := client.NewRequestOptions()

	// delete the edge caches.
	for _, cg := range testData.CacheGroups {
		if cg.Name == nil {
			t.Error("Found a Cache Group with null or undefined name")
			continue
		}

		// Retrieve the CacheGroup by name so we can get the id for Deletion
		opts.QueryParameters.Set("name", *cg.Name)
		resp, _, err := TOSession.GetCacheGroups(opts)
		assert.NoError(t, err, "Cannot GET CacheGroup by name '%s': %v - alerts: %+v", *cg.Name, err, resp.Alerts)

		if len(resp.Response) < 1 {
			t.Errorf("Could not find test data Cache Group '%s' in Traffic Ops", *cg.Name)
			continue
		}
		cg = resp.Response[0]

		// Cachegroups that are parents (usually mids but sometimes edges)
		// need to be deleted only after the children cachegroups are deleted.
		if cg.ParentCachegroupID == nil && cg.SecondaryParentCachegroupID == nil {
			parentlessCacheGroups = append(parentlessCacheGroups, cg)
			continue
		}

		if cg.ID == nil {
			t.Error("Traffic Ops returned a Cache Group with null or undefined ID")
			continue
		}

		alerts, _, err := TOSession.DeleteCacheGroup(*cg.ID, client.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Cache Group: %v - alerts: %+v", err, alerts)

		// Retrieve the CacheGroup to see if it got deleted
		opts.QueryParameters.Set("name", *cg.Name)
		cgs, _, err := TOSession.GetCacheGroups(opts)
		assert.NoError(t, err, "Error deleting Cache Group by name: %v - alerts: %+v", err, cgs.Alerts)
		assert.Equal(t, 0, len(cgs.Response), "Expected CacheGroup name: %s to be deleted", *cg.Name)
	}

	opts = client.NewRequestOptions()
	// now delete the parentless cachegroups
	for _, cg := range parentlessCacheGroups {
		// nil check for cg.Name occurs prior to insertion into parentlessCacheGroups
		opts.QueryParameters.Set("name", *cg.Name)
		// Retrieve the CacheGroup by name so we can get the id for Deletion
		resp, _, err := TOSession.GetCacheGroups(opts)
		assert.NoError(t, err, "Cannot get Cache Group by name '%s': %v - alerts: %+v", *cg.Name, err, resp.Alerts)

		if len(resp.Response) < 1 {
			t.Errorf("Cache Group '%s' somehow stopped existing since the last time we ask Traffic Ops about it", *cg.Name)
			continue
		}

		respCG := resp.Response[0]
		if respCG.ID == nil {
			t.Errorf("Traffic Ops returned Cache Group '%s' with null or undefined ID", *cg.Name)
			continue
		}
		delResp, _, err := TOSession.DeleteCacheGroup(*respCG.ID, client.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Cache Group '%s': %v - alerts: %+v", *respCG.Name, err, delResp.Alerts)

		// Retrieve the CacheGroup to see if it got deleted
		opts.QueryParameters.Set("name", *cg.Name)
		cgs, _, err := TOSession.GetCacheGroups(opts)
		assert.NoError(t, err, "Error attempting to fetch Cache Group '%s' after deletion: %v - alerts: %+v", *cg.Name, err, cgs.Alerts)
		assert.Equal(t, 0, len(cgs.Response), "Expected Cache Group '%s' to be deleted", *cg.Name)
	}
}
