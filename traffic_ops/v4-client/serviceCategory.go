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
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	// APIServiceCategories is the API version-relative path to the /service_categories API
	// endpoints.
	APIServiceCategories = "/service_categories"
)

// CreateServiceCategory creates the given Service Category.
func (to *Session) CreateServiceCategory(serviceCategory tc.ServiceCategory) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIServiceCategories, serviceCategory, nil, &alerts)
	return alerts, reqInf, err
}

// UpdateServiceCategory replaces the Service Category with the given Name with
// the one provided.
func (to *Session) UpdateServiceCategory(name string, serviceCategory tc.ServiceCategory, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%s", APIServiceCategories, name)
	var alerts tc.Alerts
	reqInf, err := to.put(route, serviceCategory, header, &alerts)
	return alerts, reqInf, err
}

// GetServiceCategories retrieves Service Categories from Traffic Ops.
func (to *Session) GetServiceCategories(values url.Values, header http.Header) ([]tc.ServiceCategory, toclientlib.ReqInf, error) {
	path := fmt.Sprintf("%s?%s", APIServiceCategories, values.Encode())
	var data tc.ServiceCategoriesResponse
	reqInf, err := to.get(path, header, &data)
	return data.Response, reqInf, err
}

// DeleteServiceCategory deletes the Service Category with the given Name.
func (to *Session) DeleteServiceCategory(name string) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%s", APIServiceCategories, name)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
