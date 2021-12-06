package tcdata

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
	"testing"
)

func (r *TCData) CreateTestServiceCategories(t *testing.T) {
	// loop through service categories, assign FKs and create
	for _, sc := range r.TestData.ServiceCategories {
		resp, _, err := TOSession.CreateServiceCategory(sc)
		if err != nil {
			t.Errorf("could not CREATE service category: %v", err)
		}
		t.Log("Response: ", resp.Alerts)
	}
}

func (r *TCData) DeleteTestServiceCategories(t *testing.T) {
	for _, sc := range r.TestData.ServiceCategories {
		// Retrieve the Service Category by name so we can get the id
		params := url.Values{}
		params.Add("name", sc.Name)
		resp, _, err := TOSession.GetServiceCategories(&params)
		if err != nil {
			t.Errorf("cannot GET Service Category by name: %s - %v", sc.Name, err)
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
				t.Errorf("error deleting Service Category: %v", err)
			}
			if len(respDelServiceCategory) > 0 {
				t.Errorf("expected Service Category %s to be deleted", sc.Name)
			}
		}
	}
}
