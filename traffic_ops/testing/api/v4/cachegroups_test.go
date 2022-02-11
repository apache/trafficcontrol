package v4

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
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestCacheGroups(t *testing.T) {
	WithObjs(t, []TCObj{Types, Parameters, CacheGroups, CDNs, Profiles, Statuses, Divisions, Regions, PhysLocations, Servers, Topologies}, func() {

		methodTests := map[string]map[string]struct {
			endpointId     func() int
			clientSession  *client.Session
			requestParams  map[string]string
			requestBody    map[string]interface{}
			requestHeaders http.Header
			expectations   []utils.CkReqFunc
		}{
			"GET": {
				"OK when VALID name parameter AND Lat/Long are 0": {
					nil, TOSession, map[string]string{"name": "nullLatLongCG"}, nil, nil,
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1), ValidateResponseFields()),
				},
				"NOT MODIFIED when NO CHANGES made": {
					nil, TOSession, nil, nil,
					http.Header{rfc.IfModifiedSince: {time.Now().AddDate(0, 0, 1).Format(time.RFC1123)}},
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"NOT MODIFIED when VALID name parameter when NO CHANGES made": {
					nil, TOSession, map[string]string{"name": "originCachegroup"}, nil,
					http.Header{rfc.IfModifiedSince: {time.Now().AddDate(0, 0, 1).Format(time.RFC1123)}},
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"NOT MODIFIED when VALID shortName parameter when NO CHANGES made": {
					nil, TOSession, map[string]string{"shortName": "mog1"}, nil,
					http.Header{rfc.IfModifiedSince: {time.Now().AddDate(0, 0, 1).Format(time.RFC1123)}},
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					nil, TOSession, nil, nil, nil,
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID name parameter": {
					nil, TOSession, map[string]string{"name": "parentCachegroup"}, nil, nil,
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						ValidateExpectedField("Name", "parentCachegroup")),
				},
				"OK when VALID shortName parameter": {
					nil, TOSession, map[string]string{"shortName": "pg2"}, nil, nil,
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						ValidateExpectedField("ShortName", "pg2")),
				},
				"OK when VALID topology parameter": {
					nil, TOSession, map[string]string{"topology": "mso-topology"}, nil, nil,
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID type parameter": {
					nil, TOSession, map[string]string{"type": "ORG_LOC"},
					nil, nil,
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						ValidateExpectedField("TypeName", "ORG_LOC")),
				},
				"EMPTY RESPONSE when INVALID id parameter": {
					nil, TOSession, map[string]string{"id": "10000"}, nil, nil,
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID type parameter": {
					nil, TOSession, map[string]string{"type": "10000"}, nil, nil,
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"FIRST RESULT when limit=1": {
					nil, TOSession, map[string]string{"orderby": "id", "limit": "1"}, nil, nil,
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), ValidatePagination("limit")),
				},
				"SECOND RESULT when offset=1": {
					nil, TOSession, map[string]string{"orderby": "id", "limit": "1", "offset": "1"}, nil, nil,
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), ValidatePagination("offset")),
				},
				"SECOND RESULT when page=2": {
					nil, TOSession, map[string]string{"orderby": "id", "limit": "1", "page": "2"}, nil, nil,
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), ValidatePagination("page")),
				},
				"BAD REQUEST when INVALID limit parameter": {
					nil, TOSession, map[string]string{"limit": "-2"}, nil, nil,
					utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID offset parameter": {
					nil, TOSession, map[string]string{"limit": "1", "offset": "0"}, nil, nil,
					utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID page parameter": {
					nil, TOSession, map[string]string{"limit": "1", "page": "0"}, nil, nil,
					utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"UNAUTHORIZED when not logged in": {
					nil, NoAuthTOSession, nil, nil, nil,
					utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusUnauthorized)),
				},
			},
			"POST": {
				"UNAUTHORIZED when not logged in": {
					nil, NoAuthTOSession, nil, nil, nil,
					utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusUnauthorized)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					GetCacheGroupId(t, "cachegroup1"), TOSession, nil,
					map[string]interface{}{
						"latitude":            17.5,
						"longitude":           17.5,
						"name":                "cachegroup1",
						"shortName":           "newShortName",
						"localizationMethods": []string{"CZ"},
						"fallbacks":           []string{"fallback1"},
						"typeName":            "EDGE_LOC",
						"typeId":              -1,
					},
					nil,
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when updating Cache Group with null Lat/Long": {
					GetCacheGroupId(t, "nullLatLongCG"), TOSession, nil,
					map[string]interface{}{
						"name":      "nullLatLongCG",
						"shortName": "null-ll",
						"typeName":  "EDGE_LOC",
						"fallbacks": []string{"fallback1"},
						"typeId":    -1,
					},
					nil,
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when updating type of cache group in topology": {
					GetCacheGroupId(t, "topology-edge-cg-01"), TOSession, nil,
					map[string]interface{}{
						"id":        -1,
						"latitude":  0,
						"longitude": 0,
						"name":      "topology-edge-cg-01",
						"shortName": "te1",
						"typeName":  "MID_LOC",
						"typeId":    -1,
					},
					nil,
					utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					GetCacheGroupId(t, "parentCachegroup"), TOSession, nil,
					map[string]interface{}{
						"latitude":  0,
						"longitude": 0,
						"name":      "parentCachegroup",
						"shortName": "pg1",
						"typeName":  "MID_LOC",
						"typeId":    -1,
					},
					http.Header{
						rfc.IfModifiedSince:   {time.Now().UTC().Add(-5 * time.Second).Format(time.RFC1123)},
						rfc.IfUnmodifiedSince: {time.Now().UTC().Add(-5 * time.Second).Format(time.RFC1123)},
					},
					utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IfMatch ETag Header": {
					GetCacheGroupId(t, "parentCachegroup2"), TOSession, nil,
					map[string]interface{}{
						"latitude":  0,
						"longitude": 0,
						"name":      "parentCachegroup2",
						"shortName": "pg2",
						"typeName":  "MID_LOC",
						"typeId":    -1,
					},
					http.Header{
						rfc.IfMatch: {rfc.ETag(time.Now().UTC().Add(-5 * time.Second))},
					},
					utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"UNAUTHORIZED when not logged in": {
					GetCacheGroupId(t, "cachegroup1"), NoAuthTOSession, nil, nil, nil,
					utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusUnauthorized)),
				},
			},
			"DELETE": {
				"NOT FOUND when INVALID id parameter": {
					func() int { return 111111 }, TOSession, nil, nil, nil,
					utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"UNAUTHORIZED when not logged in": {
					GetCacheGroupId(t, "cachegroup1"), NoAuthTOSession, nil, nil, nil,
					utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusUnauthorized)),
				},
			},
			"GET AFTER CHANGES": {
				"OK when changes were made ": {
					nil, TOSession, nil, nil,
					http.Header{
						rfc.IfModifiedSince:   {time.Now().UTC().Add(-5 * time.Second).Format(time.RFC1123)},
						rfc.IfUnmodifiedSince: {time.Now().UTC().Add(-5 * time.Second).Format(time.RFC1123)},
					},
					utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"CDNLOCK": {
				"FORBIDDEN when updating cache group when CDN LOCK exists": {
					nil, nil, nil, nil, nil, nil,
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					cg := tc.CacheGroupNullable{}
					opts := client.NewRequestOptions()

					if testCase.requestParams != nil {
						for k, v := range testCase.requestParams {
							if k == "type" {
								if _, err := strconv.Atoi(v); err != nil {
									v = strconv.Itoa(GetTypeId(t, v))
								}
							}
							opts.QueryParameters.Add(k, v)
						}
					}
					if testCase.requestHeaders != nil {
						opts.Header = testCase.requestHeaders
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
						dat, _ := json.Marshal(testCase.requestBody)
						err := json.Unmarshal(dat, &cg)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.clientSession.GetCacheGroups(opts)
							for _, check := range testCase.expectations {
								check(t, reqInf, resp, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.clientSession.CreateCacheGroup(cg, opts)
							for _, check := range testCase.expectations {
								check(t, reqInf, resp, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.clientSession.UpdateCacheGroup(testCase.endpointId(), cg, opts)
							for _, check := range testCase.expectations {
								check(t, reqInf, resp, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.clientSession.DeleteCacheGroup(testCase.endpointId(), opts)
							for _, check := range testCase.expectations {
								check(t, reqInf, alerts, err)
							}
						})
					case "CDNLOCK":
						t.Run(name, func(t *testing.T) {
							UpdateCachegroupWithLocks(t)
						})
					}
				}
			})
		}
	})
}

func ValidateExpectedField(field string, expected string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ error) {
		cgResp := resp.(tc.CacheGroupsNullableResponse)
		cg := cgResp.Response[0]
		switch field {
		case "Name":
			assert.Equal(t, expected, *cg.Name, "Expected name to be %v, but got %v", expected, *cg.Name)
		case "ShortName":
			assert.Equal(t, expected, *cg.ShortName, "Expected shortName to be %v, but got %v", expected, *cg.ShortName)
		case "TypeName":
			assert.Equal(t, expected, *cg.Type, "Expected type to be %v, but got %v", expected, *cg.Type)
		default:
			assert.Fail(t, "Expected field: %v, does not exist in response", field)
		}
	}
}

func ValidateResponseFields() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ error) {
		cgResp := resp.(tc.CacheGroupsNullableResponse)
		cg := cgResp.Response[0]
		require.NotNil(t, cg.ID, "Expected response id to not be nil")
		require.NotNil(t, cg.Latitude, "Expected response id to not be nil")
		require.NotNil(t, cg.Longitude, "Expected response id to not be nil")
		require.Equal(t, 0.0, *cg.Longitude, "Expected Longitude to be 0, but got %v", cg.Longitude)
		require.Equal(t, 0.0, *cg.Latitude, "Expected Latitude to be 0, but got %v", cg.Latitude)
	}
}

func ValidatePagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ error) {
		assert := assert.New(t)
		require := require.New(t)
		paginationResp := resp.(tc.CacheGroupsNullableResponse)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetCacheGroups(opts)
		require.NoError(err, "cannot get Cache Groups: %v - alerts: %+v", err, respBase.Alerts)

		cachegroup := respBase.Response
		require.GreaterOrEqual(len(cachegroup), 3, "Need at least 3 Cache Groups in Traffic Ops to test pagination support, found: %d", len(cachegroup))
		switch paginationParam {
		case "limit:":
			assert.Exactly(cachegroup[:1], paginationResp.Response, "expected GET Cachegroups with limit = 1 to return first result")
		case "offset":
			assert.Exactly(cachegroup[1:2], paginationResp.Response, "expected GET cachegroup with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(cachegroup[1:2], paginationResp.Response, "expected GET cachegroup with limit = 1, page = 2 to return second result")
		}
	}
}

func GetTypeId(t *testing.T, typeName string) int {
	require := require.New(t)
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", typeName)
	resp, _, err := TOSession.GetTypes(opts)

	require.NoError(err, "Get Types Request failed with error: %v", err)
	require.Equal(1, len(resp.Response), "Expected response object length 1, but got %d", len(resp.Response))
	require.NotNil(&resp.Response[0].ID, "Expected id to not be nil")

	return resp.Response[0].ID
}

func GetCacheGroupId(t *testing.T, cacheGroupName string) func() int {
	return func() int {
		require := require.New(t)
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", cacheGroupName)

		resp, _, err := TOSession.GetCacheGroups(opts)
		require.NoError(err, "Get Cache Groups Request failed with error: %v", err)
		require.Equal(len(resp.Response), 1, "Expected response object length 1, but got %d", len(resp.Response))
		require.NotNil(resp.Response[0].ID, "Expected id to not be nil")

		return *resp.Response[0].ID
	}
}

func UpdateCachegroupWithLocks(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var cdnName string
	servers := make([]tc.ServerV40, 0)
	opts := client.NewRequestOptions()
	opts.QueryParameters.Add("name", "cachegroup1")
	cgResp, _, err := TOSession.GetCacheGroups(opts)
	require.NoError(err, "couldn't get cachegroup: %v", err)
	require.Equal(1, len(cgResp.Response), "expected only one cachegroup in response, but got %d, quitting", len(cgResp.Response))

	opts.QueryParameters.Del("name")
	opts.QueryParameters.Add("cachegroupName", "cachegroup1")
	serversResp, _, err := TOSession.GetServers(opts)
	require.NoError(err, "couldn't get servers for cachegroup: %v", err)

	servers = serversResp.Response
	require.GreaterOrEqual(len(servers), 1, "couldn't get a cachegroup with 1 or more servers assigned, quitting")

	server := servers[0]
	if server.CDNName != nil {
		cdnName = *server.CDNName
	} else if server.CDNID != nil {
		opts = client.RequestOptions{}
		opts.QueryParameters.Add("id", strconv.Itoa(*server.CDNID))
		cdnResp, _, err := TOSession.GetCDNs(opts)
		require.NoError(err, "couldn't get CDN: %v", err)
		require.Equal(1, len(cdnResp.Response), "expected only one CDN in response, but got %d", len(cdnResp.Response))

		cdnName = cdnResp.Response[0].Name
	}

	// Create a new user with operations level privileges
	user1 := tc.UserV4{
		Username:             "lock_user1",
		RegistrationSent:     new(time.Time),
		LocalPassword:        util.StrPtr("test_pa$$word"),
		ConfirmLocalPassword: util.StrPtr("test_pa$$word"),
		Role:                 "operations",
	}
	user1.Email = util.StrPtr("lockuseremail@domain.com")
	user1.TenantID = 1
	user1.FullName = util.StrPtr("firstName LastName")
	_, _, err = TOSession.CreateUser(user1, client.RequestOptions{})
	require.NoError(err, "could not create test user with username: %s. err: %v", user1.Username, err)

	defer ForceDeleteTestUsersByUsernames(t, []string{"lock_user1"})

	// Establish a session with the newly created non admin level user
	userSession, _, err := client.LoginWithAgent(Config.TrafficOps.URL, user1.Username, *user1.LocalPassword, true, "to-api-v4-client-tests", false, toReqTimeout)
	require.NoError(err, "could not login with user lock_user1: %v", err)

	// Create a lock for this user
	_, _, err = userSession.CreateCDNLock(tc.CDNLock{
		CDN:     cdnName,
		Message: util.StrPtr("test lock"),
		Soft:    util.BoolPtr(false),
	}, client.RequestOptions{})
	require.NoError(err, "couldn't create cdn lock: %v", err)

	cg := cgResp.Response[0]

	// Try to update a cachegroup on a CDN that another user has a hard lock on -> this should fail
	cg.ShortName = util.StrPtr("changedShortName")
	_, reqInf, err := TOSession.UpdateCacheGroup(*cg.ID, cg, client.RequestOptions{})
	assert.Error(err, "expected an error while updating a cachegroup for a CDN for which a hard lock is held by another user, but got nothing")
	assert.Equal(http.StatusForbidden, reqInf.StatusCode, "expected a 403 forbidden status while updating a cachegroup for a CDN for which a hard lock is held by another user, but got %d", reqInf.StatusCode)

	// Try to update a cachegroup on a CDN that the same user has a hard lock on -> this should succeed
	_, reqInf, err = userSession.UpdateCacheGroup(*cg.ID, cg, client.RequestOptions{})
	assert.NoError(err, "expected no error while updating a cachegroup for a CDN for which a hard lock is held by the same user, but got %v", err)

	// Delete the lock
	_, _, err = userSession.DeleteCDNLocks(client.RequestOptions{QueryParameters: url.Values{"cdn": []string{cdnName}}})
	assert.NoError(err, "expected no error while deleting other user's lock using admin endpoint, but got %v", err)
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
		if resp.Response.LocalizationMethods == nil {
			t.Error("Localization methods are null")
		}
		if resp.Response.Fallbacks == nil {
			t.Error("Fallbacks are null")
		}
	}
}

func DeleteTestCacheGroups(t *testing.T) {
	var parentlessCacheGroups []tc.CacheGroupNullable
	opts := client.NewRequestOptions()

	// delete the edge caches.
	for _, cg := range testData.CacheGroups {
		if cg.Name == nil {
			t.Error("Found a Cache Group with null or undefined name")
			continue
		}
		// Retrieve the CacheGroup by name so we can get the id for the Update
		opts.QueryParameters.Set("name", *cg.Name)
		resp, _, err := TOSession.GetCacheGroups(opts)
		if err != nil {
			t.Errorf("cannot GET CacheGroup by name '%s': %v - alerts: %+v", *cg.Name, err, resp.Alerts)
		}
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

		// TODO: Typo here? cg is already reassigned to resp.Response[0] - is respCG supposed to be different?
		respCG := resp.Response[0]
		if respCG.ID == nil {
			t.Error("Traffic Ops returned a Cache Group with null or undefined ID")
			continue
		}
		if respCG.Name == nil {
			t.Error("Traffic Ops returned a Cache Group with null or undefined name")
			continue
		}
		alerts, _, err := TOSession.DeleteCacheGroup(*respCG.ID, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot delete Cache Group: %v - alerts: %+v", err, alerts)
		}
		// Retrieve the CacheGroup to see if it got deleted
		opts.QueryParameters.Set("name", *respCG.Name)
		cgs, _, err := TOSession.GetCacheGroups(opts)
		if err != nil {
			t.Errorf("error deleting Cache Group by name: %v - alerts: %+v", err, cgs.Alerts)
		}
		if len(cgs.Response) > 0 {
			t.Errorf("expected CacheGroup name: %s to be deleted", *cg.Name)
		}
	}

	opts = client.NewRequestOptions()
	// now delete the parentless cachegroups
	for _, cg := range parentlessCacheGroups {
		// nil check for cg.Name occurs prior to insertion into parentlessCacheGroups
		opts.QueryParameters.Set("name", *cg.Name)
		// Retrieve the CacheGroup by name so we can get the id for the Update
		resp, _, err := TOSession.GetCacheGroups(opts)
		if err != nil {
			t.Errorf("cannot get Cache Group by name '%s': %v - alerts: %+v", *cg.Name, err, resp.Alerts)
		}
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
		if err != nil {
			t.Errorf("cannot delete Cache Group '%s': %v - alerts: %+v", *respCG.Name, err, delResp.Alerts)
		}

		// Retrieve the CacheGroup to see if it got deleted
		opts.QueryParameters.Set("name", *cg.Name)
		cgs, _, err := TOSession.GetCacheGroups(opts)
		if err != nil {
			t.Errorf("error attempting to fetch Cache Group '%s' after deletion: %v - alerts: %+v", *cg.Name, err, cgs.Alerts)
		}
		if len(cgs.Response) > 0 {
			t.Errorf("expected Cache Group '%s' to be deleted", *cg.Name)
		}
	}
}
