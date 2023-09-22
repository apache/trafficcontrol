package client

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
	"fmt"
	"net/url"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiServiceCategories is the API version-relative path to the /service_categories API
// endpoints.
const apiServiceCategories = "/service_categories"

// CreateServiceCategory creates the given Service Category.
func (to *Session) CreateServiceCategory(serviceCategory tc.ServiceCategory, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiServiceCategories, opts, serviceCategory, &alerts)
	return alerts, reqInf, err
}

// UpdateServiceCategory replaces the Service Category with the given Name with
// the one provided.
func (to *Session) UpdateServiceCategory(name string, serviceCategory tc.ServiceCategory, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%s", apiServiceCategories, url.PathEscape(name))
	var alerts tc.Alerts
	reqInf, err := to.put(route, opts, serviceCategory, &alerts)
	return alerts, reqInf, err
}

// GetServiceCategories retrieves Service Categories from Traffic Ops.
func (to *Session) GetServiceCategories(opts RequestOptions) (tc.ServiceCategoriesResponse, toclientlib.ReqInf, error) {
	var data tc.ServiceCategoriesResponse
	reqInf, err := to.get(apiServiceCategories, opts, &data)
	return data, reqInf, err
}

// DeleteServiceCategory deletes the Service Category with the given Name.
func (to *Session) DeleteServiceCategory(name string, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%s", apiServiceCategories, url.PathEscape(name))
	var alerts tc.Alerts
	reqInf, err := to.del(route, opts, &alerts)
	return alerts, reqInf, err
}
