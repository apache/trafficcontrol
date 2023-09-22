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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestRoles(t *testing.T) {
	WithObjs(t, []TCObj{Roles}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.Role]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateRoleSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"new_admin"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateRoleFields(map[string]interface{}{"Name": "new_admin"})),
				},
				"VALID when SORTORDER param is DESC": {
					ClientSession: TOSession,
					RequestParams: url.Values{"orderby": {"name"}, "sortOrder": {"desc"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateRoleDescSort()),
				},
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"BAD REQUEST when INVALID CAPABILITY": {
					ClientSession: TOSession,
					RequestBody: tc.Role{
						RoleV11: tc.RoleV11{
							Name:        util.Ptr("bad_admin"),
							Description: util.Ptr("super-user 3"),
							PrivLevel:   util.Ptr(30),
						},
						Capabilities: util.Ptr([]string{
							"all-read",
							"all-write",
							"invalid-capability",
						}),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING NAME": {
					ClientSession: TOSession,
					RequestBody: tc.Role{
						RoleV11: tc.RoleV11{
							Description: util.Ptr("missing name"),
							PrivLevel:   util.Ptr(30),
						},
						Capabilities: util.Ptr([]string{
							"all-read",
							"all-write",
						}),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING DESCRIPTION": {
					ClientSession: TOSession,
					RequestBody: tc.Role{
						RoleV11: tc.RoleV11{
							Name:      util.Ptr("nodescription"),
							PrivLevel: util.Ptr(30),
						},
						Capabilities: util.Ptr([]string{
							"all-read",
							"all-write",
						}),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ROLE NAME ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody: tc.Role{
						RoleV11: tc.RoleV11{
							Name:        util.Ptr("new_admin"),
							Description: util.Ptr("description"),
							PrivLevel:   util.Ptr(30),
						},
						Capabilities: util.Ptr([]string{
							"all-read",
							"all-write",
						}),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetRoleID(t, "update_role"),
					ClientSession: TOSession,
					RequestBody: tc.Role{
						RoleV11: tc.RoleV11{
							Name:        util.Ptr("new_name"),
							Description: util.Ptr("new updated description"),
							PrivLevel:   util.Ptr(30),
						},
						Capabilities: util.Ptr([]string{
							"all-read",
							"all-write",
						}),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateRoleUpdateCreateFields("new_name", map[string]interface{}{"Name": "new_name", "Description": "new updated description"})),
				},
				"BAD REQUEST when MISSING NAME": {
					EndpointID:    GetRoleID(t, "another_role"),
					ClientSession: TOSession,
					RequestBody: tc.Role{
						RoleV11: tc.RoleV11{
							Description: util.Ptr("missing name"),
							PrivLevel:   util.Ptr(30),
						},
						Capabilities: util.Ptr([]string{
							"all-read",
							"all-write",
						}),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING DESCRIPTION": {
					EndpointID:    GetRoleID(t, "another_role"),
					ClientSession: TOSession,
					RequestBody: tc.Role{
						RoleV11: tc.RoleV11{
							Name:      util.Ptr("noDescription"),
							PrivLevel: util.Ptr(30),
						},
						Capabilities: util.Ptr([]string{
							"all-read",
							"all-write",
						}),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ADMIN ROLE": {
					EndpointID:    GetRoleID(t, "admin"),
					ClientSession: TOSession,
					RequestBody: tc.Role{
						RoleV11: tc.RoleV11{
							Name:        util.Ptr("adminUpdated"),
							PrivLevel:   util.Ptr(30),
							Description: util.Ptr("description"),
						},
						Capabilities: util.Ptr([]string{
							"all-read",
							"all-write",
						}),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"NOT FOUND when ROLE DOESNT EXIST": {
					EndpointID:    func() int { return 9999999 },
					ClientSession: TOSession,
					RequestBody: tc.Role{
						RoleV11: tc.RoleV11{
							Name:        util.Ptr("doesntexist"),
							PrivLevel:   util.Ptr(30),
							Description: util.Ptr("description"),
						},
						Capabilities: util.Ptr([]string{
							"all-read",
							"all-write",
						}),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when ROLE NAME ALREADY EXISTS": {
					EndpointID:    GetRoleID(t, "another_role"),
					ClientSession: TOSession,
					RequestBody: tc.Role{
						RoleV11: tc.RoleV11{
							Name:        util.Ptr("new_admin"),
							PrivLevel:   util.Ptr(30),
							Description: util.Ptr("description"),
						},
						Capabilities: util.Ptr([]string{
							"all-read",
							"all-write",
						}),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:     GetRoleID(t, "another_role"),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody: tc.Role{
						RoleV11: tc.RoleV11{
							Name:        util.Ptr("another_role"),
							Description: util.Ptr("super-user 3"),
							PrivLevel:   util.Ptr(30),
						},
						Capabilities: util.Ptr([]string{
							"all-read",
						}),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:     GetRoleID(t, "another_role"),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					RequestBody: tc.Role{
						RoleV11: tc.RoleV11{
							Name:        util.Ptr("another_role"),
							Description: util.Ptr("super-user 3"),
							PrivLevel:   util.Ptr(30),
						},
						Capabilities: util.Ptr([]string{
							"all-read",
						}),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {

					params := make(map[string]string)
					if testCase.RequestParams != nil {
						for k, v := range testCase.RequestParams {
							params[k] = v[0]
						}
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, _, err := testCase.ClientSession.GetRoleByQueryParamsWithHdr(params, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, _, err := testCase.ClientSession.CreateRole(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, _, err := testCase.ClientSession.UpdateRoleByIDWithHdr(testCase.EndpointID(), testCase.RequestBody, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, _, err := testCase.ClientSession.DeleteRoleByID(testCase.EndpointID())
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

func validateRoleFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Role response to not be nil.")
		roleResp := resp.([]tc.Role)
		for field, expected := range expectedResp {
			for _, role := range roleResp {
				switch field {
				case "Name":
					assert.RequireNotNil(t, role.Name, "Expected Name to not be nil.")
					assert.Equal(t, expected, *role.Name, "Expected Name to be %v, but got %s", expected, *role.Name)
				case "Description":
					assert.RequireNotNil(t, role.Description, "Expected Description to not be nil.")
					assert.Equal(t, expected, *role.Description, "Expected Description to be %v, but got %s", expected, *role.Description)
				case "PrivLevel":
					assert.RequireNotNil(t, role.PrivLevel, "Expected PrivLevel to not be nil.")
					assert.Equal(t, expected, *role.PrivLevel, "Expected PrivLevel to be %v, but got %d", expected, *role.PrivLevel)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateRoleUpdateCreateFields(name string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		role, _, _, err := TOSession.GetRoleByNameWithHdr(name, nil)
		assert.RequireNoError(t, err, "Error getting Role: %v", err)
		assert.RequireEqual(t, 1, len(role), "Expected one Role returned Got: %d", len(role))
		validateRoleFields(expectedResp)(t, toclientlib.ReqInf{}, role, tc.Alerts{}, nil)
	}
}

func validateRoleSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Role response to not be nil.")
		var roleNames []string
		roleResp := resp.([]tc.Role)
		for _, role := range roleResp {
			assert.RequireNotNil(t, role.Name, "Expected Name to not be nil.")
			roleNames = append(roleNames, *role.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(roleNames), "List is not sorted by their names: %v", roleNames)
	}
}

func validateRoleDescSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Role response to not be nil.")
		roleDescResp := resp.([]tc.Role)
		var descSortedList []string
		var ascSortedList []string
		assert.RequireGreaterOrEqual(t, len(roleDescResp), 2, "Need at least 2 Roles in Traffic Ops to test desc sort, found: %d", len(roleDescResp))
		// Get Roles in the default ascending order for comparison.
		roleAscResp, _, _, err := TOSession.GetRolesWithHdr(nil)
		assert.RequireNoError(t, err, "Unexpected error getting Roles with default sort order: %v", err)
		// Verify the response match in length, i.e. equal amount of Roles.
		assert.RequireEqual(t, len(roleAscResp), len(roleDescResp), "Expected descending order response length: %d, to match ascending order response length %d", len(roleAscResp), len(roleDescResp))
		// Insert Role names to the front of a new list, so they are now reversed to be in ascending order.
		for _, role := range roleDescResp {
			assert.RequireNotNil(t, role.Name, "Expected Name to not be nil.")
			descSortedList = append([]string{*role.Name}, descSortedList...)
		}
		// Insert Role names by appending to a new list, so they stay in ascending order.
		for _, role := range roleAscResp {
			assert.RequireNotNil(t, role.Name, "Expected Name to not be nil.")
			ascSortedList = append(ascSortedList, *role.Name)
		}
		assert.Exactly(t, ascSortedList, descSortedList, "Role responses are not equal after reversal: %v - %v", ascSortedList, descSortedList)
	}
}

func GetRoleID(t *testing.T, name string) func() int {
	return func() int {
		role, _, _, err := TOSession.GetRoleByNameWithHdr(name, nil)
		assert.RequireNoError(t, err, "Get Roles Request failed with error:", err)
		assert.RequireEqual(t, 1, len(role), "Expected response object length 1, but got %d", len(role))
		assert.RequireNotNil(t, role, "Expected role to not be nil.")
		assert.RequireNotNil(t, role[0].ID, "Expected ID to not be nil.")
		return *role[0].ID
	}
}

func CreateTestRoles(t *testing.T) {
	for _, role := range testData.Roles {
		_, _, _, err := TOSession.CreateRole(role)
		assert.NoError(t, err, "No error expected, but got %v", err)
	}
}

func DeleteTestRoles(t *testing.T) {
	roles, _, _, err := TOSession.GetRolesWithHdr(nil)
	assert.NoError(t, err, "Cannot get Roles: %v", err)
	for _, role := range roles {
		// Don't delete active roles created by test setup
		assert.RequireNotNil(t, role.Name, "Expected Name to not be nil.")
		assert.RequireNotNil(t, role.ID, "Expected ID to not be nil.")
		if *role.Name == "admin" || *role.Name == "disallowed" || *role.Name == "operations" || *role.Name == "portal" || *role.Name == "read-only" || *role.Name == "steering" || *role.Name == "federation" {
			continue
		}
		_, _, _, err := TOSession.DeleteRoleByID(*role.ID)
		assert.NoError(t, err, "Expected no error while deleting role %s, but got %v", *role.Name, err)
		// Retrieve the Role to see if it got deleted
		getRole, _, _, err := TOSession.GetRoleByIDWithHdr(*role.ID, nil)
		assert.NoError(t, err, "Error getting Role '%s' after deletion: %v", *role.Name, err)
		assert.Equal(t, 0, len(getRole), "Expected Role '%s' to be deleted, but it was found in Traffic Ops", *role.Name)
	}
}
