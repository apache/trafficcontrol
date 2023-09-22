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
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestServiceCategories(t *testing.T) {
	WithObjs(t, []TCObj{ServiceCategories}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.ServiceCategory]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateServiceCategoriesSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"serviceCategory1"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServiceCategoriesFields(map[string]interface{}{"Name": "serviceCategory1"})),
				},
				"EMPTY RESPONSE when SERVICE CATEGORY DOESNT EXIST": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"invalid"}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"BAD REQUEST when ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody:   tc.ServiceCategory{Name: "serviceCategory1"},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when NAME FIELD is BLANK": {
					ClientSession: TOSession,
					RequestBody:   tc.ServiceCategory{Name: ""},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"barServiceCategory2"}},
					RequestBody:   tc.ServiceCategory{Name: "newName"},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServiceCategoriesUpdateCreateFields("newName", map[string]interface{}{"Name": "newName"})),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					ClientSession:  TOSession,
					RequestParams:  url.Values{"name": {"serviceCategory1"}},
					RequestHeaders: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					RequestBody:    tc.ServiceCategory{Name: "newName"},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					ClientSession:  TOSession,
					RequestParams:  url.Values{"name": {"serviceCategory1"}},
					RequestHeaders: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					RequestBody:    tc.ServiceCategory{Name: "newName"},
					Expectations:   utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"NOT FOUND when DOESNT EXIST": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"invalid"}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetServiceCategoriesWithHdr(&testCase.RequestParams, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateServiceCategory(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateServiceCategoryByName(testCase.RequestParams["name"][0], testCase.RequestBody, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteServiceCategoryByName(testCase.RequestParams["name"][0])
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

func validateServiceCategoriesFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Service Categories response to not be nil.")
		serviceCategoryResp := resp.([]tc.ServiceCategory)
		for field, expected := range expectedResp {
			for _, serviceCategory := range serviceCategoryResp {
				switch field {
				case "Name":
					assert.Equal(t, expected, serviceCategory.Name, "Expected Name to be %v, but got %s", expected, serviceCategory.Name)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateServiceCategoriesUpdateCreateFields(name string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		values := url.Values{}
		values.Set("name", name)
		serviceCategories, _, err := TOSession.GetServiceCategoriesWithHdr(&values, nil)
		assert.RequireNoError(t, err, "Error getting Service Categories: %v", err)
		assert.RequireEqual(t, 1, len(serviceCategories), "Expected one Service Category returned Got: %d", len(serviceCategories))
		validateServiceCategoriesFields(expectedResp)(t, toclientlib.ReqInf{}, serviceCategories, tc.Alerts{}, nil)
	}
}

func validateServiceCategoriesSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Service Categories response to not be nil.")
		var serviceCategoryNames []string
		serviceCategoryResp := resp.([]tc.ServiceCategory)
		for _, serviceCategory := range serviceCategoryResp {
			serviceCategoryNames = append(serviceCategoryNames, serviceCategory.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(serviceCategoryNames), "List is not sorted by their names: %v", serviceCategoryNames)
	}
}

func CreateTestServiceCategories(t *testing.T) {
	for _, serviceCategory := range testData.ServiceCategories {
		resp, _, err := TOSession.CreateServiceCategory(serviceCategory)
		assert.RequireNoError(t, err, "Could not create Service Category: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestServiceCategories(t *testing.T) {
	values := url.Values{}
	serviceCategories, _, err := TOSession.GetServiceCategoriesWithHdr(&values, nil)
	assert.NoError(t, err, "Cannot get Service Categories: %v", err)

	for _, serviceCategory := range serviceCategories {
		alerts, _, err := TOSession.DeleteServiceCategoryByName(serviceCategory.Name)
		assert.NoError(t, err, "Unexpected error deleting Service Category '%s': %v - alerts: %+v", serviceCategory.Name, err, alerts.Alerts)
		// Retrieve the Service Category to see if it got deleted
		values.Set("name", serviceCategory.Name)
		getServiceCategory, _, err := TOSession.GetServiceCategoriesWithHdr(&values, nil)
		assert.NoError(t, err, "Error getting Service Category '%s' after deletion: %v", serviceCategory.Name, err)
		assert.Equal(t, 0, len(getServiceCategory), "Expected Service Category '%s' to be deleted, but it was found in Traffic Ops", serviceCategory.Name)
	}
}
