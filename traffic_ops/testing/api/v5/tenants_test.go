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
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestTenants(t *testing.T) {
	WithObjs(t, []TCObj{Tenants}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.TenantV5]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						utils.ResponseLengthGreaterOrEqual(1), validateTenantSort()),
				},
				"OK when VALID ACTIVE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"active": {"true"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateTenantFields(map[string]interface{}{"Active": true})),
				},
				"VALID when SORTORDER param is DESC": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"sortOrder": {"desc"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateTenantDescSort()),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateTenantPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateTenantPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateTenantPagination("page")),
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
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestBody: tc.TenantV5{
						Active:     util.Ptr(true),
						Name:       util.Ptr("tenant5"),
						ParentName: util.Ptr("root"),
						ParentID:   util.Ptr(GetTenantID(t, "root")()),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateTenantCreateUpdateFields(map[string]interface{}{"Name": "tenant5"})),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetTenantID(t, "tenant4"),
					ClientSession: TOSession,
					RequestBody: tc.TenantV5{
						Active:     util.Ptr(false),
						Name:       util.Ptr("newname"),
						ParentName: util.Ptr("root"),
						ParentID:   util.Ptr(GetTenantID(t, "root")()),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateTenantCreateUpdateFields(map[string]interface{}{"Name": "newname", "Active": false})),
				},
				"BAD REQUEST when ROOT TENANT": {
					EndpointID:    GetTenantID(t, "root"),
					ClientSession: TOSession,
					RequestBody: tc.TenantV5{
						Active:     util.Ptr(false),
						Name:       util.Ptr("tenant1"),
						ParentName: util.Ptr("root"),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:    GetTenantID(t, "tenant2"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: tc.TenantV5{
						Active:     util.Ptr(false),
						Name:       util.Ptr("tenant2"),
						ParentName: util.Ptr("root"),
						ParentID:   util.Ptr(GetTenantID(t, "root")()),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetTenantID(t, "tenant2"),
					ClientSession: TOSession,
					RequestBody: tc.TenantV5{
						Active:     util.Ptr(false),
						Name:       util.Ptr("tenant2"),
						ParentName: util.Ptr("root"),
						ParentID:   util.Ptr(GetTenantID(t, "root")()),
					},
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"BAD REQUEST when TENANT HAS CHILDREN": {
					EndpointID:    GetTenantID(t, "tenant1"),
					ClientSession: TOSession,
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
							resp, reqInf, err := testCase.ClientSession.GetTenants(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateTenant(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateTenant(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteTenant(testCase.EndpointID(), testCase.RequestOpts)
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

func validateTenantFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Tenant response to not be nil.")
		tenantResp := resp.([]tc.TenantV5)
		for field, expected := range expectedResp {
			for _, tenant := range tenantResp {
				switch field {
				case "Active":
					assert.Equal(t, expected, *tenant.Active, "Expected Active to be %v, but got %b", expected, tenant.Active)
				case "Name":
					assert.Equal(t, expected, *tenant.Name, "Expected Name to be %v, but got %s", expected, tenant.Name)
				case "ParentName":
					assert.Equal(t, expected, *tenant.ParentName, "Expected ParentName to be %v, but got %s", expected, tenant.ParentName)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateTenantCreateUpdateFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Tenant response to not be nil.")
		tenantResp := resp.(tc.TenantV5)
		tenants := []tc.TenantV5{tenantResp}
		validateTenantFields(expectedResp)(t, toclientlib.ReqInf{}, tenants, tc.Alerts{}, nil)
	}
}

func validateTenantSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Tenant response to not be nil.")
		var tenants []string
		tenantResp := resp.([]tc.TenantV5)
		for _, tenant := range tenantResp {
			tenants = append(tenants, *tenant.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(tenants), "List is not sorted by their names: %v", tenants)
	}
}

func validateTenantDescSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Tenant response to not be nil.")
		tenantDescResp := resp.([]tc.TenantV5)
		var descSortedList []string
		var ascSortedList []string
		assert.RequireGreaterOrEqual(t, len(tenantDescResp), 2, "Need at least 2 Tenants in Traffic Ops to test desc sort, found: %d", len(tenantDescResp))
		// Get Tenants in the default ascending order for comparison.
		tenantsAscResp, _, err := TOSession.GetTenants(client.RequestOptions{})
		assert.RequireNoError(t, err, "Unexpected error getting Tenants with default sort order: %v - alerts: %+v", err, tenantsAscResp.Alerts)
		// Verify the response match in length, i.e. equal amount of Tenants.
		assert.RequireEqual(t, len(tenantsAscResp.Response), len(tenantDescResp), "Expected descending order response length: %v, to match ascending order response length %v", len(tenantsAscResp.Response), len(tenantDescResp))
		// Insert Tenant names to the front of a new list, so they are now reversed to be in ascending order.
		for _, tenant := range tenantDescResp {
			descSortedList = append([]string{*tenant.Name}, descSortedList...)
		}
		// Insert Tenant names by appending to a new list, so they stay in ascending order.
		for _, tenant := range tenantsAscResp.Response {
			ascSortedList = append(ascSortedList, *tenant.Name)
		}
		assert.Exactly(t, ascSortedList, descSortedList, "Tenant responses are not equal after reversal: %v - %v", ascSortedList, descSortedList)
	}
}

func validateTenantPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.TenantV5)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetTenants(opts)
		assert.RequireNoError(t, err, "Cannot get Tenants: %v - alerts: %+v", err, respBase.Alerts)

		tenants := respBase.Response
		assert.RequireGreaterOrEqual(t, len(tenants), 3, "Need at least 3 Tenants in Traffic Ops to test pagination support, found: %d", len(tenants))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, tenants[:1], paginationResp, "expected GET Tenants with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, tenants[1:2], paginationResp, "expected GET Tenants with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, tenants[1:2], paginationResp, "expected GET Tenants with limit = 1, page = 2 to return second result")
		}
	}
}

func GetTenantID(t *testing.T, name string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		tenants, _, err := TOSession.GetTenants(opts)
		assert.RequireNoError(t, err, "Get Tenants Request failed with error:", err)
		assert.RequireEqual(t, 1, len(tenants.Response), "Expected response object length 1, but got %d", len(tenants.Response))
		return *tenants.Response[0].ID
	}
}

func CreateTestTenants(t *testing.T) {
	for _, tenant := range testData.Tenants {
		resp, _, err := TOSession.CreateTenant(tenant, client.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Tenant '%s': %v - alerts: %+v", tenant.Name, err, resp.Alerts)
	}
}

func DeleteTestTenants(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("sortOrder", "desc")
	tenants, _, err := TOSession.GetTenants(opts)
	assert.NoError(t, err, "Cannot get Tenants: %v - alerts: %+v", err, tenants.Alerts)

	for _, tenant := range tenants.Response {
		if *tenant.Name == "root" {
			continue
		}
		alerts, _, err := TOSession.DeleteTenant(*tenant.ID, client.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Tenant '%s' (#%d): %v - alerts: %+v", tenant.Name, tenant.ID, err, alerts.Alerts)
		// Retrieve the Tenant to see if it got deleted
		opts.QueryParameters.Set("id", strconv.Itoa(*tenant.ID))
		getTenants, _, err := TOSession.GetTenants(opts)
		assert.NoError(t, err, "Error getting Tenant '%s' after deletion: %v - alerts: %+v", tenant.Name, err, getTenants.Alerts)
		assert.Equal(t, 0, len(getTenants.Response), "Expected Tenant '%s' to be deleted, but it was found in Traffic Ops", tenant.Name)
	}
}
