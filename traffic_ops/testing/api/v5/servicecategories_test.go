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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestServiceCategories(t *testing.T) {
	WithObjs(t, []TCObj{ServiceCategories}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.ServiceCategoryV5]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateServiceCategoriesSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"serviceCategory1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateServiceCategoriesFields(map[string]interface{}{"Name": "serviceCategory1"})),
				},
				"VALID when SORTORDER param is DESC": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"sortOrder": {"desc"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateServiceCategoriesDescSort()),
				},
				"EMPTY RESPONSE when SERVICE CATEGORY DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"invalid"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateServiceCategoriesPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateServiceCategoriesPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateServiceCategoriesPagination("page")),
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
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"BAD REQUEST when ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody:   tc.ServiceCategoryV5{Name: "serviceCategory1"},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when NAME FIELD is BLANK": {
					ClientSession: TOSession,
					RequestBody:   tc.ServiceCategoryV5{Name: ""},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"barServiceCategory2"}}},
					RequestBody:   tc.ServiceCategoryV5{Name: "newName"},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServiceCategoriesUpdateCreateFields("newName", map[string]interface{}{"Name": "newName"})),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					ClientSession: TOSession,
					RequestOpts: client.RequestOptions{
						QueryParameters: url.Values{"name": {"serviceCategory1"}},
						Header:          http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}},
					},
					RequestBody:  tc.ServiceCategoryV5{Name: "newName"},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					ClientSession: TOSession,
					RequestBody:   tc.ServiceCategoryV5{Name: "newName"},
					RequestOpts: client.RequestOptions{
						QueryParameters: url.Values{"name": {"serviceCategory1"}},
						Header:          http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"NOT FOUND when DOESNT EXIST": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"invalid"}}},
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
							resp, reqInf, err := testCase.ClientSession.GetServiceCategories(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateServiceCategory(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							if _, ok := testCase.RequestOpts.QueryParameters["name"]; !ok {
								t.Fatalf("Query Parameter: \"name\" is required for PUT method tests.")
							}
							alerts, reqInf, err := testCase.ClientSession.UpdateServiceCategory(testCase.RequestOpts.QueryParameters["name"][0], testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							if _, ok := testCase.RequestOpts.QueryParameters["name"]; !ok {
								t.Fatalf("Query Parameter: \"name\" is required for DELETE method tests.")
							}
							alerts, reqInf, err := testCase.ClientSession.DeleteServiceCategory(testCase.RequestOpts.QueryParameters["name"][0], testCase.RequestOpts)
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
		serviceCategoryResp := resp.([]tc.ServiceCategoryV5)
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
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		serviceCategories, _, err := TOSession.GetServiceCategories(opts)
		assert.RequireNoError(t, err, "Error getting Service Categories: %v - alerts: %+v", err, serviceCategories.Alerts)
		assert.RequireEqual(t, 1, len(serviceCategories.Response), "Expected one Service Category returned Got: %d", len(serviceCategories.Response))
		validateServiceCategoriesFields(expectedResp)(t, toclientlib.ReqInf{}, serviceCategories.Response, tc.Alerts{}, nil)
	}
}

func validateServiceCategoriesPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.ServiceCategoryV5)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetServiceCategories(opts)
		assert.RequireNoError(t, err, "Cannot get Service Categories: %v - alerts: %+v", err, respBase.Alerts)

		serviceCategories := respBase.Response
		assert.RequireGreaterOrEqual(t, len(serviceCategories), 2, "Need at least 2 Service Categories in Traffic Ops to test pagination support, found: %d", len(serviceCategories))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, serviceCategories[:1], paginationResp, "expected GET Service Categories with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, serviceCategories[1:2], paginationResp, "expected GET Service Categories with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, serviceCategories[1:2], paginationResp, "expected GET Service Categories with limit = 1, page = 2 to return second result")
		}
	}
}

func validateServiceCategoriesSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Service Categories response to not be nil.")
		var serviceCategoryNames []string
		serviceCategoryResp := resp.([]tc.ServiceCategoryV5)
		for _, serviceCategory := range serviceCategoryResp {
			serviceCategoryNames = append(serviceCategoryNames, serviceCategory.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(serviceCategoryNames), "List is not sorted by their names: %v", serviceCategoryNames)
	}
}

func validateServiceCategoriesDescSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Service Categories response to not be nil.")
		serviceCategoriesDescResp := resp.([]tc.ServiceCategoryV5)
		var descSortedList []string
		var ascSortedList []string
		assert.RequireGreaterOrEqual(t, len(serviceCategoriesDescResp), 2, "Need at least 2 Service Categories in Traffic Ops to test desc sort, found: %d", len(serviceCategoriesDescResp))
		// Get Service Categories in the default ascending order for comparison.
		serviceCategoriesAscResp, _, err := TOSession.GetServiceCategories(client.RequestOptions{})
		assert.RequireNoError(t, err, "Unexpected error getting Service Categories with default sort order: %v - alerts: %+v", err, serviceCategoriesAscResp.Alerts)
		// Verify the response match in length, i.e. equal amount of Service Categories.
		assert.RequireEqual(t, len(serviceCategoriesAscResp.Response), len(serviceCategoriesDescResp), "Expected descending order response length: %v, to match ascending order response length %v", len(serviceCategoriesAscResp.Response), len(serviceCategoriesDescResp))
		// Insert Service Category names to the front of a new list, so they are now reversed to be in ascending order.
		for _, serviceCategory := range serviceCategoriesDescResp {
			descSortedList = append([]string{serviceCategory.Name}, descSortedList...)
		}
		// Insert Service Category names by appending to a new list, so they stay in ascending order.
		for _, serviceCategory := range serviceCategoriesAscResp.Response {
			ascSortedList = append(ascSortedList, serviceCategory.Name)
		}
		assert.Exactly(t, ascSortedList, descSortedList, "Service Categories responses are not equal after reversal: %v - %v", ascSortedList, descSortedList)
	}
}

func CreateTestServiceCategories(t *testing.T) {
	for _, serviceCategory := range testData.ServiceCategories {
		resp, _, err := TOSession.CreateServiceCategory(serviceCategory, client.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Service Category: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestServiceCategories(t *testing.T) {
	serviceCategories, _, err := TOSession.GetServiceCategories(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get Service Categories: %v - alerts: %+v", err, serviceCategories.Alerts)

	for _, serviceCategory := range serviceCategories.Response {
		alerts, _, err := TOSession.DeleteServiceCategory(serviceCategory.Name, client.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Service Category '%s': %v - alerts: %+v", serviceCategory.Name, err, alerts.Alerts)
		// Retrieve the Service Category to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", serviceCategory.Name)
		getServiceCategory, _, err := TOSession.GetServiceCategories(opts)
		assert.NoError(t, err, "Error getting Service Category '%s' after deletion: %v - alerts: %+v", serviceCategory.Name, err, getServiceCategory.Alerts)
		assert.Equal(t, 0, len(getServiceCategory.Response), "Expected Service Category '%s' to be deleted, but it was found in Traffic Ops", serviceCategory.Name)
	}
}
