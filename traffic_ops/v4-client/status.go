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
	// APIStatuses is the API version-relative path to the /statuses API endpoint.
	APIStatuses = "/statuses"
)

// CreateStatus creates the given Status.
func (to *Session) CreateStatus(status tc.StatusNullable) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIStatuses, status, nil, &alerts)
	return alerts, reqInf, err
}

// UpdateStatus replaces the Status identified by 'id' with the one provided.
func (to *Session) UpdateStatus(id int, status tc.Status, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APIStatuses, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, status, header, &alerts)
	return alerts, reqInf, err
}

// GetStatuses retrieves all Statuses stored in Traffic Ops.
func (to *Session) GetStatuses(header http.Header) ([]tc.Status, toclientlib.ReqInf, error) {
	var data tc.StatusesResponse
	reqInf, err := to.get(APIStatuses, header, &data)
	return data.Response, reqInf, err
}

// GetStatusByID returns the Status with the given ID.
func (to *Session) GetStatusByID(id int, header http.Header) ([]tc.Status, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIStatuses, id)
	var data tc.StatusesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetStatusByName returns the Status with the given Name.
func (to *Session) GetStatusByName(name string, header http.Header) ([]tc.Status, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", APIStatuses, url.QueryEscape(name))
	var data tc.StatusesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// DeleteStatus deletes the Status wtih the given ID.
func (to *Session) DeleteStatus(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APIStatuses, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
