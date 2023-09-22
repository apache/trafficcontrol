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
	"net/http"
	"net/url"
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestRoles(t *testing.T) {
	WithObjs(t, []TCObj{Roles}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.RoleV4]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateRoleSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"new_admin"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateRoleFields(map[string]interface{}{"Name": "new_admin"})),
				},
				"VALID when SORTORDER param is DESC": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"name"}, "sortOrder": {"desc"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateRoleDescSort()),
				},
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"BAD REQUEST when MISSING NAME": {
					ClientSession: TOSession,
					RequestBody: tc.RoleV4{
						Description: "missing name",
						Permissions: []string{
							"all-read",
							"all-write",
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING DESCRIPTION": {
					ClientSession: TOSession,
					RequestBody: tc.RoleV4{
						Name: "noDescription",
						Permissions: []string{
							"all-read",
							"all-write",
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ROLE NAME ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody: tc.RoleV4{
						Name:        "new_admin",
						Description: "description",
						Permissions: []string{
							"all-read",
							"all-write",
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"update_role"}}},
					RequestBody: tc.RoleV4{
						Name:        "new_name",
						Description: "new updated description",
						Permissions: []string{
							"all-read",
							"all-write",
						},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateRoleUpdateCreateFields("new_name", map[string]interface{}{"Name": "new_name", "Description": "new updated description"})),
				},
				"BAD REQUEST when MISSING NAME": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"another_role"}}},
					RequestBody: tc.RoleV4{
						Description: "missing name",
						Permissions: []string{
							"all-read",
							"all-write",
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING DESCRIPTION": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"another_role"}}},
					RequestBody: tc.RoleV4{
						Name: "noDescription",
						Permissions: []string{
							"all-read",
							"all-write",
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ADMIN ROLE": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"admin"}}},
					RequestBody: tc.RoleV4{
						Name:        "adminUpdated",
						Description: "description",
						Permissions: []string{
							"all-read",
							"all-write",
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"NOT FOUND when ROLE DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"doesntexist"}}},
					RequestBody: tc.RoleV4{
						Name:        "doesntexist",
						Description: "description",
						Permissions: []string{
							"all-read",
							"all-write",
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when ROLE NAME ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"another_role"}}},
					RequestBody: tc.RoleV4{
						Name:        "new_admin",
						Description: "description",
						Permissions: []string{
							"all-read",
							"all-write",
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					ClientSession: TOSession,
					RequestOpts: client.RequestOptions{
						QueryParameters: url.Values{"name": {"another_role"}},
						Header:          http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					},
					RequestBody: tc.RoleV4{
						Name:        "another_role",
						Description: "super-user 3",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					ClientSession: TOSession,
					RequestOpts: client.RequestOptions{
						QueryParameters: url.Values{"name": {"another_role"}},
						Header:          http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					},
					RequestBody: tc.RoleV4{
						Name:        "another_role",
						Description: "super-user 3",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"BAD REQUEST when SPECIAL ADMIN ROLE": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {tc.AdminRoleName}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetRoles(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateRole(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateRole(testCase.RequestOpts.QueryParameters["name"][0], testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteRole(testCase.RequestOpts.QueryParameters["name"][0], testCase.RequestOpts)
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
		roleResp := resp.([]tc.RoleV4)
		for field, expected := range expectedResp {
			for _, role := range roleResp {
				switch field {
				case "Name":
					assert.Equal(t, expected, role.Name, "Expected Name to be %v, but got %s", expected, role.Name)
				case "Description":
					assert.Equal(t, expected, role.Description, "Expected Description to be %v, but got %s", expected, role.Description)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateRoleUpdateCreateFields(name string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		role, _, err := TOSession.GetRoles(opts)
		assert.RequireNoError(t, err, "Error getting Role: %v - alerts: %+v", err, role.Alerts)
		assert.RequireEqual(t, 1, len(role.Response), "Expected one Role returned Got: %d", len(role.Response))
		validateRoleFields(expectedResp)(t, toclientlib.ReqInf{}, role.Response, tc.Alerts{}, nil)
	}
}

func validateRoleSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Role response to not be nil.")
		var roleNames []string
		roleResp := resp.([]tc.RoleV4)
		for _, role := range roleResp {
			roleNames = append(roleNames, role.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(roleNames), "List is not sorted by their names: %v", roleNames)
	}
}

func validateRoleDescSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Role response to not be nil.")
		roleDescResp := resp.([]tc.RoleV4)
		var descSortedList []string
		var ascSortedList []string
		assert.RequireGreaterOrEqual(t, len(roleDescResp), 2, "Need at least 2 Roles in Traffic Ops to test desc sort, found: %d", len(roleDescResp))
		// Get Roles in the default ascending order for comparison.
		roleAscResp, _, err := TOSession.GetRoles(client.RequestOptions{})
		assert.RequireNoError(t, err, "Unexpected error getting Roles with default sort order: %v - alerts: %+v", err, roleAscResp.Alerts)
		// Verify the response match in length, i.e. equal amount of Roles.
		assert.RequireEqual(t, len(roleAscResp.Response), len(roleDescResp), "Expected descending order response length: %d, to match ascending order response length %d", len(roleAscResp.Response), len(roleDescResp))
		// Insert Role names to the front of a new list, so they are now reversed to be in ascending order.
		for _, role := range roleDescResp {
			descSortedList = append([]string{role.Name}, descSortedList...)
		}
		// Insert Role names by appending to a new list, so they stay in ascending order.
		for _, role := range roleAscResp.Response {
			ascSortedList = append(ascSortedList, role.Name)
		}
		assert.Exactly(t, ascSortedList, descSortedList, "Role responses are not equal after reversal: %v - %v", ascSortedList, descSortedList)
	}
}

func CreateTestRoles(t *testing.T) {
	for _, role := range testData.Roles {
		_, _, err := TOSession.CreateRole(role, client.RequestOptions{})
		assert.NoError(t, err, "No error expected, but got %v", err)
	}
}

func DeleteTestRoles(t *testing.T) {
	roles, _, err := TOSession.GetRoles(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get Roles: %v - alerts: %+v", err, roles.Alerts)
	for _, role := range roles.Response {
		// Don't delete active roles created by test setup
		if role.Name == "admin" || role.Name == "disallowed" || role.Name == "operations" || role.Name == "portal" || role.Name == "read-only" || role.Name == "steering" || role.Name == "federation" {
			continue
		}
		_, _, err := TOSession.DeleteRole(role.Name, client.NewRequestOptions())
		assert.NoError(t, err, "Expected no error while deleting role %s, but got %v", role.Name, err)
		// Retrieve the Role to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", role.Name)
		getRole, _, err := TOSession.GetRoles(opts)
		assert.NoError(t, err, "Error getting Role '%s' after deletion: %v - alerts: %+v", role.Name, err, getRole.Alerts)
		assert.Equal(t, 0, len(getRole.Response), "Expected Role '%s' to be deleted, but it was found in Traffic Ops", role.Name)
	}
}
