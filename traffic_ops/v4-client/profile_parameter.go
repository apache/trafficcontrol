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

// apiProfileParameters is the full path to the /profileparameters API endpoint.
const apiProfileParameters = "/profileparameters"
const apiProfileParameter = "/profileparameter"

// Supported query string parameter names.
const (
	ProfileIDQueryParam   = "profileId"
	ParameterIDQueryParam = "parameterId"
)

// CreateProfileParameter assigns a Parameter to a Profile.
func (to *Session) CreateProfileParameter(pp tc.ProfileParameterCreationRequest, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiProfileParameters, opts, pp, &alerts)
	return alerts, reqInf, err
}

// CreateMultipleProfileParameters assigns multip Parameters to one or more
// Profiles at once.
func (to *Session) CreateMultipleProfileParameters(pps []tc.ProfileParameterCreationRequest, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiProfileParameters, opts, pps, &alerts)
	return alerts, reqInf, err
}

// CreateProfileWithMultipleParameters assigns multipe Parameters to one or more
// Profiles at once.
func (to *Session) CreateProfileWithMultipleParameters(pps tc.PostProfileParam, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiProfileParameter, opts, pps, &alerts)
	return alerts, reqInf, err
}

// GetProfileParameters retrieves associations between Profiles and Parameters.
func (to *Session) GetProfileParameters(opts RequestOptions) (tc.ProfileParametersAPIResponse, toclientlib.ReqInf, error) {
	var data tc.ProfileParametersAPIResponse
	reqInf, err := to.get(apiProfileParameters, opts, &data)
	return data, reqInf, err
}

// DeleteProfileParameter removes the Parameter with the ID 'parameter' from
// the Profile identified by the ID 'profile'.
func (to *Session) DeleteProfileParameter(profile int, parameter int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	URI := fmt.Sprintf("%s/%d/%d", apiProfileParameters, profile, parameter)
	var alerts tc.Alerts
	reqInf, err := to.del(URI, opts, &alerts)
	return alerts, reqInf, err
}
