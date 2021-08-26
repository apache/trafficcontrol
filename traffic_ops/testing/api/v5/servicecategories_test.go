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
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v5-client"
)

func TestServiceCategories(t *testing.T) {
	WithObjs(t, []TCObj{ServiceCategories}, func() {
		GetTestServiceCategoriesIMS(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		header.Set(rfc.IfUnmodifiedSince, time)
		SortTestServiceCategories(t)
		UpdateTestServiceCategories(t)
		GetTestServiceCategories(t)
		GetTestServiceCategoriesIMSAfterChange(t, header)
		UpdateTestServiceCategoriesWithHeaders(t, header)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestServiceCategoriesWithHeaders(t, header)
		GetTestServiceCategoriesByInvalidName(t)
		VerifyPaginationSupportServiceCategories(t)
		SortTestServiceCategoriesDesc(t)
		CreateTestServiceCategoriesAlreadyExist(t)
		CreateTestServiceCategoriesInvalidName(t)
		DeleteTestServiceCategoriesInvalidName(t)
	})
}

func GetTestServiceCategoriesIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)
	for _, sc := range testData.ServiceCategories {
		opts.QueryParameters.Add("name", sc.Name)
		resp, reqInf, err := TOSession.GetServiceCategories(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestServiceCategoriesIMSAfterChange(t *testing.T, header http.Header) {
	opts := client.NewRequestOptions()
	opts.Header = header
	for _, sc := range testData.ServiceCategories {
		opts.QueryParameters.Set("name", sc.Name)
		resp, reqInf, err := TOSession.GetServiceCategories(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	opts.Header.Set(rfc.IfModifiedSince, timeStr)
	for _, sc := range testData.ServiceCategories {
		opts.QueryParameters.Set("name", sc.Name)
		resp, reqInf, err := TOSession.GetServiceCategories(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestServiceCategories(t *testing.T) {
	// loop through service categories, assign FKs and create
	for _, sc := range testData.ServiceCategories {
		resp, _, err := TOSession.CreateServiceCategory(sc, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create Service Category: %v - alerts: %+v", err, resp.Alerts)
		}
	}
}

func CreateTestServiceCategoriesAlreadyExist(t *testing.T) {
	firstServiceCategory := tc.ServiceCategory{}
	if len(testData.ServiceCategories) > 0 {
		firstServiceCategory = testData.ServiceCategories[0]
	} else {
		t.Fatalf("cannot CREATE DUPLICATE Service Category, test data does not have service categories")
	}
	resp, reqInf, err := TOSession.CreateServiceCategory(firstServiceCategory, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected service_category name already exists. but got: %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
	}
}

func CreateTestServiceCategoriesInvalidName(t *testing.T) {
	firstServiceCategory := tc.ServiceCategory{}
	firstServiceCategory.Name = ""
	resp, reqInf, err := TOSession.CreateServiceCategory(firstServiceCategory, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected 'name' cannot be blanks. but got: %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestServiceCategories(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, sc := range testData.ServiceCategories {
		opts.QueryParameters.Set("name", sc.Name)
		resp, _, err := TOSession.GetServiceCategories(opts)
		if err != nil {
			t.Errorf("cannot get Service Category by name: %v - alerts: %+v", err, resp.Alerts)
		}
	}
}

func SortTestServiceCategories(t *testing.T) {
	resp, _, err := TOSession.GetServiceCategories(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
	}
	sortedList := make([]string, 0, len(resp.Response))
	for _, sc := range resp.Response {
		sortedList = append(sortedList, sc.Name)
	}

	if !sort.StringsAreSorted(sortedList) {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func UpdateTestServiceCategoriesWithHeaders(t *testing.T, h http.Header) {
	firstServiceCategory := tc.ServiceCategory{}
	if len(testData.ServiceCategories) > 0 {
		firstServiceCategory = testData.ServiceCategories[0]
	} else {
		t.Fatalf("cannot UPDATE Service Category, test data does not have service categories")
	}
	_, reqInf, err := TOSession.UpdateServiceCategory(firstServiceCategory.Name, firstServiceCategory, client.RequestOptions{Header: h})
	if err == nil {
		t.Errorf("attempting to update service category with headers - expected: error, actual: nil")
	}
	if reqInf.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("expected status code: %d, actual: %d", http.StatusPreconditionFailed, reqInf.StatusCode)
	}
}

func UpdateTestServiceCategories(t *testing.T) {
	firstServiceCategory := tc.ServiceCategory{}
	if len(testData.ServiceCategories) > 0 {
		firstServiceCategory = testData.ServiceCategories[0]
	} else {
		t.Fatalf("cannot UPDATE Service Category, test data does not have service categories")
	}

	// Retrieve the Service Category by service category so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", firstServiceCategory.Name)
	resp, _, err := TOSession.GetServiceCategories(opts)
	if err != nil {
		t.Errorf("cannot get Service Category '%s' by name: %v - alerts: %+v", firstServiceCategory.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Service Category to exist with name '%s', found: %d", firstServiceCategory.Name, len(resp.Response))
	}
	remoteServiceCategory := resp.Response[0]
	remoteServiceCategory.Name = "ServiceCategory2"

	alert, _, err := TOSession.UpdateServiceCategory(firstServiceCategory.Name, remoteServiceCategory, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Service Category: %v - alerts: %+v", err, alert.Alerts)
	}

	// Retrieve the Service Category to check service category got updated
	opts.QueryParameters.Set("name", remoteServiceCategory.Name)
	resp, _, err = TOSession.GetServiceCategories(opts)
	if err != nil {
		t.Errorf("cannot get Service Category '%s' by name: %v - alerts: %+v", remoteServiceCategory.Name, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatal("empty response getting Service Category after update")
	} else if len(resp.Response) > 1 {
		t.Errorf("expected a name to uniquely identify exactly one Service Category, got: %d", len(resp.Response))
	}

	// revert back to original name
	alert, _, err = TOSession.UpdateServiceCategory(remoteServiceCategory.Name, firstServiceCategory, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Service Category: %v - alerts: %+v", err, alert.Alerts)
	}

	// Retrieve the Service Category to check service category got updated
	opts.QueryParameters.Set("name", firstServiceCategory.Name)
	resp, _, err = TOSession.GetServiceCategories(opts)
	if err != nil {
		t.Errorf("cannot get Service Category '%s' by name: %v - alerts: %+v", firstServiceCategory.Name, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatal("empty response getting Service Category after update")
	} else if len(resp.Response) > 1 {
		t.Errorf("expected a name to uniquely identify exactly one Service Category, got: %d", len(resp.Response))
	}
}

func DeleteTestServiceCategories(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, sc := range testData.ServiceCategories {
		// Retrieve the Service Category by name so we can get the id
		opts.QueryParameters.Set("name", sc.Name)
		resp, _, err := TOSession.GetServiceCategories(opts)
		if err != nil {
			t.Errorf("cannot get Service Category '%s' by name: %v - alerts: %+v", sc.Name, err, resp.Alerts)
		}
		if len(resp.Response) != 1 {
			t.Errorf("Expected exactly one Service Category to exist with name '%s', found: %d", sc.Name, len(resp.Response))
			continue
		}
		respServiceCategory := resp.Response[0]

		delResp, _, err := TOSession.DeleteServiceCategory(respServiceCategory.Name, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot delete Service Category: %v - alerts: %+v", err, delResp.Alerts)
		}

		// Retrieve the Service Category to see if it got deleted
		respDelServiceCategory, _, err := TOSession.GetServiceCategories(opts)
		if err != nil {
			t.Errorf("error deleting Service Category: %v - alerts: %+v", err, respDelServiceCategory.Alerts)
		}
		if len(respDelServiceCategory.Response) > 0 {
			t.Errorf("expected Service Category '%s' to be deleted", sc.Name)
		}
	}
}

func DeleteTestServiceCategoriesInvalidName(t *testing.T) {

	delResp, reqInf, err := TOSession.DeleteServiceCategory("invalid", client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected no serviceCategory with that key found but got: %v - alerts: %+v", err, delResp.Alerts)
	}
	if reqInf.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected 404 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestServiceCategoriesByInvalidName(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", "invalid")
	resp, _, _ := TOSession.GetServiceCategories(opts)
	if len(resp.Response) > 0 {
		t.Errorf("Expected 0 response, but got many %v", resp)
	}
}

func VerifyPaginationSupportServiceCategories(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "id")
	resp, _, err := TOSession.GetServiceCategories(opts)
	if err != nil {
		t.Fatalf("cannot get Service Categories: %v - alerts: %+v", err, resp.Alerts)
	}
	serviceCategories := resp.Response
	if len(serviceCategories) < 2 {
		t.Fatalf("Need at least 2 Service Categories in Traffic Ops to test pagination support, found: %d", len(serviceCategories))
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	serviceCategoriesWithLimit, _, err := TOSession.GetServiceCategories(opts)
	if err == nil {
		if !reflect.DeepEqual(serviceCategories[:1], serviceCategoriesWithLimit.Response) {
			t.Error("expected GET Service Categories with limit = 1 to return first result")
		}
	} else {
		t.Errorf("Error in getting Service Categories by limit %v", err)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "1")
	serviceCategoriesWithOffset, _, err := TOSession.GetServiceCategories(opts)
	if err == nil {
		if !reflect.DeepEqual(serviceCategories[1:2], serviceCategoriesWithOffset.Response) {
			t.Error("expected GET Service Categories with limit = 1, offset = 1 to return second result")
		}
	} else {
		t.Errorf("Error in getting Service Categories by limit and offset %v", err)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "2")
	serviceCategoriesWithPage, _, err := TOSession.GetServiceCategories(opts)
	if err == nil {
		if !reflect.DeepEqual(serviceCategories[1:2], serviceCategoriesWithPage.Response) {
			t.Error("expected GET Service Categories with limit = 1, page = 2 to return second result")
		}
	} else {
		t.Errorf("Error in getting Service Categories by limit and page %v", err)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "-2")
	resp, _, err = TOSession.GetServiceCategories(opts)
	if err == nil {
		t.Error("expected GET Service Categories to return an error when limit is not bigger than -1")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be bigger than -1") {
		t.Errorf("expected GET Service Categories to return an error for limit is not bigger than -1, actual error: " + err.Error())
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "0")
	resp, _, err = TOSession.GetServiceCategories(opts)
	if err == nil {
		t.Error("expected GET Service Categories to return an error when offset is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET Service Categories to return an error for offset is not a positive integer, actual error: " + err.Error())
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "0")
	resp, _, err = TOSession.GetServiceCategories(opts)
	if err == nil {
		t.Error("expected GET Service Categories to return an error when page is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET Service Categories to return an error for page is not a positive integer, actual error: " + err.Error())
	}
}

func SortTestServiceCategoriesDesc(t *testing.T) {
	resp, _, err := TOSession.GetServiceCategories(client.RequestOptions{})
	if err != nil {
		t.Errorf("Expected no error, but got error in Service Categories with default ordering: %v - alerts: %+v", err, resp.Alerts)
	}
	respAsc := resp.Response
	if len(respAsc) < 1 {
		t.Fatal("Need at least one Service Categories in Traffic Ops to test Service Categories sort ordering")
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("sortOrder", "desc")
	resp, _, err = TOSession.GetServiceCategories(opts)
	if err != nil {
		t.Errorf("Expected no error, but got error in Service Categories with Descending ordering: %v - alerts: %+v", err, resp.Alerts)
	}
	respDesc := resp.Response
	if len(respDesc) < 1 {
		t.Fatal("Need at least one Service Categories in Traffic Ops to test Service Categories sort ordering")
	}

	if len(respAsc) != len(respDesc) {
		t.Fatalf("Traffic Ops returned %d Service Categories using default sort order, but %d Service Categories when sort order was explicitly set to descending", len(respAsc), len(respDesc))
	}

	// reverse the descending-sorted response and compare it to the ascending-sorted one
	// TODO ensure at least two in each slice? A list of length one is
	// trivially sorted both ascending and descending.
	for start, end := 0, len(respDesc)-1; start < end; start, end = start+1, end-1 {
		respDesc[start], respDesc[end] = respDesc[end], respDesc[start]
	}
	if respDesc[0].Name != respAsc[0].Name {
		t.Errorf("Service Categories responses are not equal after reversal: Asc: %s - Desc: %s", respDesc[0].Name, respAsc[0].Name)
	}
}
