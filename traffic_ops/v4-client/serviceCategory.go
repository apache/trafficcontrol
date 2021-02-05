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

package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_SERVICE_CATEGORIES = apiBase + "/service_categories"
)

// CreateServiceCategory performs a post to create a service category.
func (to *Session) CreateServiceCategory(serviceCategory tc.ServiceCategory) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(API_SERVICE_CATEGORIES, serviceCategory, nil, &alerts)
	return alerts, reqInf, err
}

// UpdateServiceCategoryByName updates a service category by its unique name.
func (to *Session) UpdateServiceCategoryByName(name string, serviceCategory tc.ServiceCategory, header http.Header) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%s", API_SERVICE_CATEGORIES, name)
	var alerts tc.Alerts
	reqInf, err := to.put(route, serviceCategory, header, &alerts)
	return alerts, reqInf, err
}

// GetServiceCategoriesWithHdr gets a list of service categories by the passed in url values and http headers.
func (to *Session) GetServiceCategoriesWithHdr(values *url.Values, header http.Header) ([]tc.ServiceCategory, ReqInf, error) {
	path := fmt.Sprintf("%s?%s", API_SERVICE_CATEGORIES, values.Encode())
	var data tc.ServiceCategoriesResponse
	reqInf, err := to.get(path, header, &data)
	return data.Response, reqInf, err
}

// GetServiceCategories gets a list of service categories by the passed in url values.
func (to *Session) GetServiceCategories(values *url.Values) ([]tc.ServiceCategory, ReqInf, error) {
	return to.GetServiceCategoriesWithHdr(values, nil)
}

// DeleteServiceCategoryByName deletes a service category by the service
// category's unique name.
func (to *Session) DeleteServiceCategoryByName(name string) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%s", API_SERVICE_CATEGORIES, name)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
