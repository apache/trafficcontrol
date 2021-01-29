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
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_SERVICE_CATEGORIES = apiBase + "/service_categories"
)

// CreateServiceCategory performs a post to create a service category.
func (to *Session) CreateServiceCategory(serviceCategory tc.ServiceCategory) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(serviceCategory)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_SERVICE_CATEGORIES, reqBody, nil)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return tc.Alerts{}, reqInf, err
	}
	return alerts, reqInf, nil
}

// UpdateServiceCategoryByName updates a service category by its unique name.
func (to *Session) UpdateServiceCategoryByName(name string, serviceCategory tc.ServiceCategory, header http.Header) (tc.Alerts, ReqInf, error) {
	var remoteAddr net.Addr
	reqBody, err := json.Marshal(serviceCategory)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s/%s", API_SERVICE_CATEGORIES, name)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody, header)
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		defer resp.Body.Close()
	}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	var alerts tc.Alerts
	if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return tc.Alerts{}, reqInf, err
	}
	return alerts, reqInf, nil
}

// GetServiceCategoriesWithHdr gets a list of service categories by the passed in url values and http headers.
func (to *Session) GetServiceCategoriesWithHdr(values *url.Values, header http.Header) ([]tc.ServiceCategory, ReqInf, error) {
	path := fmt.Sprintf("%s?%s", API_SERVICE_CATEGORIES, values.Encode())
	resp, remoteAddr, err := to.request(http.MethodGet, path, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.ServiceCategory{}, reqInf, nil
		}
	}
	defer resp.Body.Close()

	var data tc.ServiceCategoriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetServiceCategories gets a list of service categories by the passed in url values.
func (to *Session) GetServiceCategories(values *url.Values) ([]tc.ServiceCategory, ReqInf, error) {
	return to.GetServiceCategoriesWithHdr(values, nil)
}

// DeleteServiceCategoryByName deletes a service category by the service
// category's unique name.
func (to *Session) DeleteServiceCategoryByName(name string) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%s", API_SERVICE_CATEGORIES, name)
	resp, remoteAddr, err := to.request(http.MethodDelete, route, nil, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return tc.Alerts{}, reqInf, err
	}
	return alerts, reqInf, nil
}
