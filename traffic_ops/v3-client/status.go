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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	// API_STATUSES is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_STATUSES = apiBase + "/statuses"

	APIStatuses = "/statuses"
)

// CreateStatusNullable creates a new status, using the tc.StatusNullable structure.
func (to *Session) CreateStatusNullable(status tc.StatusNullable) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIStatuses, status, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdateStatusByIDWithHdr(id int, status tc.Status, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APIStatuses, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, status, header, &alerts)
	return alerts, reqInf, err
}

// UpdateStatusByID updates a Status by ID.
// Deprecated: UpdateStatusByID will be removed in 6.0. Use UpdateStatusByIDWithHdr.
func (to *Session) UpdateStatusByID(id int, status tc.Status) (tc.Alerts, toclientlib.ReqInf, error) {
	return to.UpdateStatusByIDWithHdr(id, status, nil)
}

func (to *Session) GetStatusesWithHdr(header http.Header) ([]tc.Status, toclientlib.ReqInf, error) {
	var data tc.StatusesResponse
	reqInf, err := to.get(APIStatuses, header, &data)
	return data.Response, reqInf, err
}

// GetStatuses returns a list of Statuses.
// Deprecated: GetStatuses will be removed in 6.0. Use GetStatusesWithHdr.
func (to *Session) GetStatuses() ([]tc.Status, toclientlib.ReqInf, error) {
	return to.GetStatusesWithHdr(nil)
}

func (to *Session) GetStatusByIDWithHdr(id int, header http.Header) ([]tc.Status, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIStatuses, id)
	var data tc.StatusesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetStatusByID GETs a Status by the Status ID.
// Deprecated: GetStatusByID will be removed in 6.0. Use GetStatusByIDWithHdr.
func (to *Session) GetStatusByID(id int) ([]tc.Status, toclientlib.ReqInf, error) {
	return to.GetStatusByIDWithHdr(id, nil)
}

func (to *Session) GetStatusByNameWithHdr(name string, header http.Header) ([]tc.Status, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", APIStatuses, url.QueryEscape(name))
	var data tc.StatusesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetStatusByName GETs a Status by the Status name.
// Deprecated: GetStatusByName will be removed in 6.0. Use GetStatusByNameWithHdr.
func (to *Session) GetStatusByName(name string) ([]tc.Status, toclientlib.ReqInf, error) {
	return to.GetStatusByNameWithHdr(name, nil)
}

// DeleteStatusByID DELETEs a Status by ID.
func (to *Session) DeleteStatusByID(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APIStatuses, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
