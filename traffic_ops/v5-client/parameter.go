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

// apiParameters is the full path to the /parameters API endpoint.
const apiParameters = "/parameters"

// CreateParameter performs a POST to create a Parameter.
func (to *Session) CreateParameter(pl tc.ParameterV5, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiParameters, opts, pl, &alerts)
	return alerts, reqInf, err
}

// CreateMultipleParameters performs a POST to create multiple Parameters at once.
func (to *Session) CreateMultipleParameters(pls []tc.ParameterV5, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiParameters, opts, pls, &alerts)
	return alerts, reqInf, err
}

// UpdateParameter replaces the Parameter identified by 'id' with the one
// provided.
func (to *Session) UpdateParameter(id int, pl tc.ParameterV5, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", apiParameters, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, opts, pl, &alerts)
	return alerts, reqInf, err
}

// GetParameters returns all Parameters in Traffic Ops.
func (to *Session) GetParameters(opts RequestOptions) (tc.ParametersResponseV5, toclientlib.ReqInf, error) {
	var data tc.ParametersResponseV5
	reqInf, err := to.get(apiParameters, opts, &data)
	return data, reqInf, err
}

// DeleteParameter deletes the Parameter with the given ID.
func (to *Session) DeleteParameter(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	URI := fmt.Sprintf("%s/%d", apiParameters, id)
	var alerts tc.Alerts
	reqInf, err := to.del(URI, opts, &alerts)
	return alerts, reqInf, err
}
