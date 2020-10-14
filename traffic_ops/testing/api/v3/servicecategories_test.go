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

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
)

func TestServiceCategories(t *testing.T) {
	WithObjs(t, []TCObj{Tenants, ServiceCategories, Users}, func() {
		GetTestServiceCategoriesIMS(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		SortTestServiceCategories(t)
		UpdateTestServiceCategories(t)
		GetTestServiceCategories(t)
		ServiceCategoryTenancyTest(t)
		GetTestServiceCategoriesIMSAfterChange(t, header)
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
		tenant, _, err := TOSession.TenantByName(sc.TenantName)
		sc.TenantID = tenant.ID
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

func UpdateTestServiceCategories(t *testing.T) {
	firstServiceCategory := tc.ServiceCategory{}
	if len(testData.ServiceCategories) > 0 {
		firstServiceCategory = testData.ServiceCategories[0]
	} else {
		t.Fatalf("cannot UPDATE Service Category, test data does not have service categories")
	}

	tenants, _, err := TOSession.Tenants()
	if err != nil {
		t.Fatalf("Failed to get tenants: %v", err)
	}
	if len(tenants) < 2 {
		t.Fatalf("Need at least two tenants to test changing tenant; got: %d", len(tenants))
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

		originalTenant := remoteServiceCategory.TenantID
		found := false
		for _, tenant := range tenants {
			if tenant.ID != originalTenant {
				remoteServiceCategory.TenantID = tenant.ID
				found = true
				break
			}
		}
		if !found {
			t.Fatal("Could not find tenant that isn't the same as the remote service category's tenant")
		}

		var alert tc.Alerts
		alert, _, err = TOSession.UpdateServiceCategoryByName(remoteServiceCategory.Name, remoteServiceCategory)
		if err != nil {
			t.Errorf("cannot UPDATE Service Category by name: %v - %v", err, alert)
		}
		t.Logf("alerts: %v", alert)

		// Retrieve the Service Category to check service category got updated
		resp, _, err = TOSession.GetServiceCategories(&params)
		if err != nil {
			t.Errorf("cannot GET Service Category by service category: %v - %v", firstServiceCategory.Name, err)
		}
		if len(resp) < 1 {
			t.Fatal("empty response getting Service Category after update")
		} else if len(resp) > 1 {
			t.Errorf("expected a name to uniquely identify exactly one Service Category, got: %d", len(resp))
		}

		respServiceCategory := resp[0]
		if respServiceCategory.TenantID != remoteServiceCategory.TenantID {
			t.Errorf("results do not match; want: %d, got: %d", remoteServiceCategory.TenantID, respServiceCategory.TenantID)
		}

		// Set the name back to the fixture value so we can delete it after
		remoteServiceCategory.TenantID = originalTenant
		alert, _, err = TOSession.UpdateServiceCategoryByName(remoteServiceCategory.Name, remoteServiceCategory)
		if err != nil {
			t.Errorf("cannot UPDATE Service Category by name: %v - %v", err, alert)
		}
	}
}

func ServiceCategoryTenancyTest(t *testing.T) {
	var alert tc.Alerts
	tenant3, _, err := TOSession.TenantByName("tenant3")
	if err != nil {
		t.Errorf("cannot GET Tenant3: %v", err)
	}

	params := url.Values{}
	serviceCategories, _, err := TOSession.GetServiceCategories(&params)
	if err != nil {
		t.Errorf("cannot GET Service Categories: %v", err)
	}
	for _, sc := range serviceCategories {
		if sc.Name == "serviceCategory1" {
			alert, _, err = TOSession.UpdateServiceCategoryByName(sc.Name, sc)
			if err != nil {
				t.Errorf("cannot UPDATE Service Category by name: %v - %v", err, alert)
			}
			sc.TenantID = tenant3.ID
		}
	}

	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	tenant4TOClient, _, err := toclient.LoginWithAgent(TOSession.URL, "tenant4user", "pa$$word", true, "to-api-v3-client-tests/tenant4user", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to log in with tenant4user: %v", err.Error())
	}

	serviceCategoriesReadableByTenant4, _, err := tenant4TOClient.GetServiceCategories(&params)
	if err != nil {
		t.Error("tenant4user cannot GET service categories")
	}

	// assert that tenant4user cannot read service categories outside of its tenant
	for _, sc := range serviceCategoriesReadableByTenant4 {
		if sc.Name == "serviceCategory1" {
			t.Error("expected tenant4 to be unable to read service categories from tenant 3")
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
