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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiStatuses is the API version-relative path to the /statuses API endpoint.
const apiStatuses = "/statuses"

// CreateStatus creates the given Status.
func (to *Session) CreateStatus(status tc.StatusNullable, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiStatuses, opts, status, &alerts)
	return alerts, reqInf, err
}

// UpdateStatus replaces the Status identified by 'id' with the one provided.
func (to *Session) UpdateStatus(id int, status tc.Status, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiStatuses, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, opts, status, &alerts)
	return alerts, reqInf, err
}

// GetStatuses retrieves all Statuses stored in Traffic Ops.
func (to *Session) GetStatuses(opts RequestOptions) (tc.StatusesResponse, toclientlib.ReqInf, error) {
	var data tc.StatusesResponse
	reqInf, err := to.get(apiStatuses, opts, &data)
	return data, reqInf, err
}

// DeleteStatus deletes the Status with the given ID.
func (to *Session) DeleteStatus(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiStatuses, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, opts, &alerts)
	return alerts, reqInf, err
}
