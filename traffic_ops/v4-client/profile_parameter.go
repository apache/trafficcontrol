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

// APIProfileParameters is the full path to the /profileparameters API endpoint.
const APIProfileParameters = "/profileparameters"

// Supported query string parameter names.
const (
	ProfileIDQueryParam   = "profileId"
	ParameterIDQueryParam = "parameterId"
)

// CreateProfileParameter assigns a Parameter to a Profile.
func (to *Session) CreateProfileParameter(pp tc.ProfileParameter) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIProfileParameters, pp, nil, &alerts)
	return alerts, reqInf, err
}

// CreateMultipleProfileParameters assigns multip Parameters to one or more
// Profiles at once.
func (to *Session) CreateMultipleProfileParameters(pps []tc.ProfileParameter) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIProfileParameters, pps, nil, &alerts)
	return alerts, reqInf, err
}

// GetProfileParameters retrieves associations between Profiles and Parameters.
func (to *Session) GetProfileParameters(queryParams url.Values, header http.Header) ([]tc.ProfileParameter, toclientlib.ReqInf, error) {
	URI := APIProfileParameters
	if len(queryParams) > 0 {
		URI += "?" + queryParams.Encode()
	}
	var data tc.ProfileParametersNullableResponse
	reqInf, err := to.get(URI, header, &data)
	if err != nil {
		return nil, reqInf, err
	}
	ret := make([]tc.ProfileParameter, len(data.Response))
	for i, pp := range data.Response {
		ret[i] = tc.ProfileParameter{}
		if pp.Profile != nil {
			ret[i].Profile = *pp.Profile
		}
		if pp.Parameter != nil {
			ret[i].ParameterID = *pp.Parameter
		}
	}
	return ret, reqInf, nil
}

// DeleteProfileParameter removes the Parameter with the ID 'parameter' from
// the Profile identified by the ID 'profile'.
func (to *Session) DeleteProfileParameter(profile int, parameter int) (tc.Alerts, toclientlib.ReqInf, error) {
	URI := fmt.Sprintf("%s/%d/%d", APIProfileParameters, profile, parameter)
	var alerts tc.Alerts
	reqInf, err := to.del(URI, nil, &alerts)
	return alerts, reqInf, err
}
