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
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
)

func TestServiceCategories(t *testing.T) {
	WithObjs(t, []TCObj{Tenants, ServiceCategories, Users}, func() {
		UpdateTestServiceCategories(t)
		GetTestServiceCategories(t)
		ServiceCategoryTenancyTest(t)
	})
}

func CreateTestServiceCategories(t *testing.T) {
	tenant1 := "tenant1"
	// loop through service categories, assign FKs and create
	for _, sc := range testData.ServiceCategories {
		tenant, _, err := TOSession.TenantByName(tenant1, nil)
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

func UpdateTestServiceCategories(t *testing.T) {
	firstServiceCategory := tc.ServiceCategory{}
	if len(testData.ServiceCategories) > 0 {
		firstServiceCategory = testData.ServiceCategories[0]
	} else {
		t.Errorf("cannot UPDATE Service Category, test data does not have service categories")
	}
	// Retrieve the Service Category by service category so we can get the id for the Update
	paramsFirst := url.Values{}
	paramsFirst.Add("name", firstServiceCategory.Name)
	resp, _, err := TOSession.GetServiceCategories(&paramsFirst)
	if err != nil {
		t.Errorf("cannot GET Service Category by service category: %v - %v", firstServiceCategory.Name, err)
	}
	if len(resp) > 0 {
		remoteServiceCategory := resp[0]
		expectedServiceCategory := firstServiceCategory.Name + "-test-update"
		remoteServiceCategory.Name = expectedServiceCategory
		var alert tc.Alerts
		alert, _, err = TOSession.UpdateServiceCategoryByID(remoteServiceCategory.ID, remoteServiceCategory)
		if err != nil {
			t.Errorf("cannot UPDATE Service Category by id: %v - %v", err, alert)
		}

		// Retrieve the Service Category to check service category got updated
		paramsRemote := url.Values{}
		paramsRemote.Add("id", strconv.Itoa(remoteServiceCategory.ID))
		resp, _, err = TOSession.GetServiceCategories(&paramsRemote)
		if err != nil {
			t.Errorf("cannot GET Service Category by service category: %v - %v", firstServiceCategory.Name, err)
		}

		respServiceCategory := resp[0]
		if respServiceCategory.Name != expectedServiceCategory {
			t.Errorf("results do not match actual: %s, expected: %s", respServiceCategory.Name, expectedServiceCategory)
		}

		// Set the name back to the fixture value so we can delete it after
		remoteServiceCategory.Name = firstServiceCategory.Name
		alert, _, err = TOSession.UpdateServiceCategoryByID(remoteServiceCategory.ID, remoteServiceCategory)
		if err != nil {
			t.Errorf("cannot UPDATE Service Category by id: %v - %v", err, alert)
		}
	}
}

func ServiceCategoryTenancyTest(t *testing.T) {
	var alert tc.Alerts
	tenant3, _, err := TOSession.TenantByName("tenant3", nil)
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
			alert, _, err = TOSession.UpdateServiceCategoryByID(sc.ID, sc)
			if err != nil {
				t.Errorf("cannot UPDATE Service Category by id: %v - %v", err, alert)
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

			delResp, _, err := TOSession.DeleteServiceCategoryByID(respServiceCategory.ID)
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
