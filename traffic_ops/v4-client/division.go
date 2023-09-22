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

// apiDivisions is the API version-relative path to the /divisions API route.
const apiDivisions = "/divisions"

// CreateDivision creates the given Division.
func (to *Session) CreateDivision(division tc.Division, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiDivisions, opts, division, &alerts)
	return alerts, reqInf, err
}

// UpdateDivision replaces the Division identified by 'id' with the one
// provided.
func (to *Session) UpdateDivision(id int, division tc.Division, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiDivisions, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, opts, division, &alerts)
	return alerts, reqInf, err
}

// GetDivisions returns Divisions from Traffic Ops.
func (to *Session) GetDivisions(opts RequestOptions) (tc.DivisionsResponse, toclientlib.ReqInf, error) {
	var data tc.DivisionsResponse
	reqInf, err := to.get(apiDivisions, opts, &data)
	return data, reqInf, err
}

// DeleteDivision deletes the Division identified by 'id'.
func (to *Session) DeleteDivision(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiDivisions, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, opts, &alerts)
	return alerts, reqInf, err
}
