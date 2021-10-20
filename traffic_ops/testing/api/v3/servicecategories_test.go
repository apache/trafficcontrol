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

	"github.com/apache/trafficcontrol/v6/lib/go-rfc"
	"github.com/apache/trafficcontrol/v6/lib/go-tc"
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
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	params := url.Values{}
	for _, sc := range testData.ServiceCategories {
		params.Add("name", sc.Name)
		_, reqInf, err := TOSession.GetServiceCategoriesWithHdr(&params, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestServiceCategoriesIMSAfterChange(t *testing.T, header http.Header) {
	params := url.Values{}
	for _, sc := range testData.ServiceCategories {
		params.Add("name", sc.Name)
		_, reqInf, err := TOSession.GetServiceCategoriesWithHdr(&params, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	params = url.Values{}
	for _, sc := range testData.ServiceCategories {
		params.Add("name", sc.Name)
		_, reqInf, err := TOSession.GetServiceCategoriesWithHdr(&params, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestServiceCategories(t *testing.T) {
	// loop through service categories, assign FKs and create
	for _, sc := range testData.ServiceCategories {
		resp, _, err := TOSession.CreateServiceCategory(sc)
		if err != nil {
			t.Errorf("could not CREATE service category: %v", err)
		}
		t.Log("Response: ", resp.Alerts)
	}
}

func GetTestServiceCategories(t *testing.T) {
	params := url.Values{}
	for _, sc := range testData.ServiceCategories {
		params.Add("name", sc.Name)
		resp, _, err := TOSession.GetServiceCategories(&params)
		if err != nil {
			t.Errorf("cannot GET Service Category by name: %v - %v", err, resp)
		}
	}
}

func SortTestServiceCategories(t *testing.T) {
	var header http.Header
	params := url.Values{}
	var sortedList []string
	resp, _, err := TOSession.GetServiceCategoriesWithHdr(&params, header)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	for i, _ := range resp {
		sortedList = append(sortedList, resp[i].Name)
	}

	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
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
	_, reqInf, err := TOSession.UpdateServiceCategoryByName(firstServiceCategory.Name, firstServiceCategory, h)
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
	params := url.Values{}
	params.Add("name", firstServiceCategory.Name)
	resp, _, err := TOSession.GetServiceCategories(&params)
	if err != nil {
		t.Errorf("cannot GET Service Category by name: %v - %v", firstServiceCategory.Name, err)
	}
	if len(resp) > 0 {
		remoteServiceCategory := resp[0]
		remoteServiceCategory.Name = "ServiceCategory2"

		var alert tc.Alerts
		alert, _, err = TOSession.UpdateServiceCategoryByName(firstServiceCategory.Name, remoteServiceCategory, nil)
		if err != nil {
			t.Errorf("cannot UPDATE Service Category by name: %v - %v", err, alert)
		}
		t.Logf("alerts: %v", alert)

		// Retrieve the Service Category to check service category got updated
		params := url.Values{}
		params.Add("name", remoteServiceCategory.Name)
		resp, _, err = TOSession.GetServiceCategories(&params)
		if err != nil {
			t.Errorf("cannot GET Service Category by service category: %v - %v", remoteServiceCategory.Name, err)
		}
		if len(resp) < 1 {
			t.Fatal("empty response getting Service Category after update")
		} else if len(resp) > 1 {
			t.Errorf("expected a name to uniquely identify exactly one Service Category, got: %d", len(resp))
		}

		// revert back to original name
		alert, _, err = TOSession.UpdateServiceCategoryByName(remoteServiceCategory.Name, firstServiceCategory, nil)
		if err != nil {
			t.Errorf("cannot UPDATE Service Category by name: %v - %v", err, alert)
		}
		t.Logf("alerts: %v", alert)

		// Retrieve the Service Category to check service category got updated
		params = url.Values{}
		params.Add("name", firstServiceCategory.Name)
		resp, _, err = TOSession.GetServiceCategories(&params)
		if err != nil {
			t.Errorf("cannot GET Service Category by service category: %v - %v", firstServiceCategory.Name, err)
		}
		if len(resp) < 1 {
			t.Fatal("empty response getting Service Category after update")
		} else if len(resp) > 1 {
			t.Errorf("expected a name to uniquely identify exactly one Service Category, got: %d", len(resp))
		}
	}
}

func DeleteTestServiceCategories(t *testing.T) {
	for _, sc := range testData.ServiceCategories {
		// Retrieve the Service Category by name so we can get the id
		params := url.Values{}
		params.Add("name", sc.Name)
		resp, _, err := TOSession.GetServiceCategories(&params)
		if err != nil {
			t.Errorf("cannot GET Service Category by name: %v - %v", sc.Name, err)
		}
		if len(resp) > 0 {
			respServiceCategory := resp[0]

			delResp, _, err := TOSession.DeleteServiceCategoryByName(respServiceCategory.Name)
			if err != nil {
				t.Errorf("cannot DELETE Service Category by service category: %v - %v", err, delResp)
			}

			// Retrieve the Service Category to see if it got deleted
			respDelServiceCategory, _, err := TOSession.GetServiceCategories(&params)
			if err != nil {
				t.Errorf("error deleting Service Category: %s", err.Error())
			}
			if len(respDelServiceCategory) > 0 {
				t.Errorf("expected Service Category : %s to be deleted", sc.Name)
			}
		}
	}
}
