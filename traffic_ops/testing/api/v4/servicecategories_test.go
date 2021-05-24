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
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
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
