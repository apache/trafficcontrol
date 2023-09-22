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
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestTenants(t *testing.T) {
	WithObjs(t, []TCObj{Tenants}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.Tenant]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						utils.ResponseLengthGreaterOrEqual(1), validateTenantSort()),
				},
			},
			"POST": {
				"NO ERROR when VALID request": {
					ClientSession: TOSession,
					RequestBody: tc.Tenant{
						Active:     true,
						Name:       "tenant5",
						ParentName: "root",
						ParentID:   GetTenantID(t, "root")(),
					},
					Expectations: utils.CkRequest(utils.NoError(), validateTenantCreateUpdateFields(map[string]interface{}{"Name": "tenant5"})),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetTenantID(t, "tenant4"),
					ClientSession: TOSession,
					RequestBody: tc.Tenant{
						Active:     false,
						Name:       "newname",
						ParentName: "root",
						ParentID:   GetTenantID(t, "root")(),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateTenantCreateUpdateFields(map[string]interface{}{"Name": "newname", "Active": false})),
				},
				"BAD REQUEST when ROOT TENANT": {
					EndpointID:    GetTenantID(t, "root"),
					ClientSession: TOSession,
					RequestBody: tc.Tenant{
						Active:     false,
						Name:       "tenant1",
						ParentName: "root",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:     GetTenantID(t, "tenant2"),
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody: tc.Tenant{
						Active:     false,
						Name:       "tenant2",
						ParentName: "root",
						ParentID:   GetTenantID(t, "root")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetTenantID(t, "tenant2"),
					ClientSession: TOSession,
					RequestBody: tc.Tenant{
						Active:     false,
						Name:       "tenant2",
						ParentName: "root",
						ParentID:   GetTenantID(t, "root")(),
					},
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"ERROR when TENANT HAS CHILDREN": {
					EndpointID:    GetTenantID(t, "tenant1"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError()),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.TenantsWithHdr(testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, err := testCase.ClientSession.CreateTenant(&testCase.RequestBody)
							for _, check := range testCase.Expectations {
								if resp != nil {
									check(t, toclientlib.ReqInf{}, resp.Response, resp.Alerts, err)
								}
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateTenantWithHdr(strconv.Itoa(testCase.EndpointID()), &testCase.RequestBody, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								if resp != nil {
									check(t, reqInf, resp.Response, resp.Alerts, err)
								}
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							_, err := testCase.ClientSession.DeleteTenant(strconv.Itoa(testCase.EndpointID()))
							for _, check := range testCase.Expectations {
								check(t, toclientlib.ReqInf{}, nil, tc.Alerts{}, err)
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
		tenantResp := resp.([]tc.Tenant)
		for field, expected := range expectedResp {
			for _, tenant := range tenantResp {
				switch field {
				case "Active":
					assert.Equal(t, expected, tenant.Active, "Expected Active to be %v, but got %b", expected, tenant.Active)
				case "Name":
					assert.Equal(t, expected, tenant.Name, "Expected Name to be %v, but got %s", expected, tenant.Name)
				case "ParentName":
					assert.Equal(t, expected, tenant.ParentName, "Expected ParentName to be %v, but got %s", expected, tenant.ParentName)
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
		tenantResp := resp.(tc.Tenant)
		tenants := []tc.Tenant{tenantResp}
		validateTenantFields(expectedResp)(t, toclientlib.ReqInf{}, tenants, tc.Alerts{}, nil)
	}
}

func validateTenantSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Tenant response to not be nil.")
		var tenants []string
		tenantResp := resp.([]tc.Tenant)
		for _, tenant := range tenantResp {
			tenants = append(tenants, tenant.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(tenants), "List is not sorted by their names: %v", tenants)
	}
}

func GetTenantID(t *testing.T, name string) func() int {
	return func() int {
		tenant, _, err := TOSession.TenantByNameWithHdr(name, nil)
		assert.RequireNoError(t, err, "Get Tenants Request failed with error:", err)
		assert.RequireNotNil(t, tenant, "Expected tenant to not be nil")
		return tenant.ID
	}
}

func CreateTestTenants(t *testing.T) {
	for _, tenant := range testData.Tenants {
		resp, err := TOSession.CreateTenant(&tenant)
		assert.RequireNoError(t, err, "Could not create Tenant '%s': %v - alerts: %+v", tenant.Name, err, resp.Alerts)
	}
}

func DeleteTestTenants(t *testing.T) {
	tenants, _, err := TOSession.TenantsWithHdr(nil)
	assert.NoError(t, err, "Cannot get Tenants: %v", err)

	for i := len(tenants) - 1; i >= 0; i-- {
		if tenants[i].Name == "root" {
			continue
		}
		alerts, err := TOSession.DeleteTenant(strconv.Itoa(tenants[i].ID))
		assert.NoError(t, err, "Unexpected error deleting Tenant '%s' (#%d): %v - alerts: %+v", tenants[i].Name, tenants[i].ID, err, alerts.Alerts)
	}
	// Retrieve Tenants to see if they got deleted, only root should exist
	tenants, _, err = TOSession.TenantsWithHdr(nil)
	assert.NoError(t, err, "Error getting Tenants after deletion: %v", err)
	assert.Equal(t, len(tenants), 1, "Expected only 1 Tenant, but found %d.", len(tenants))
}
