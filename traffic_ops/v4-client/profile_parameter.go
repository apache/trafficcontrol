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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	APIProfileParameters  = "/profileparameters"
	ProfileIdQueryParam   = "profileId"
	ParameterIdQueryParam = "parameterId"
)

// Create a ProfileParameter
func (to *Session) CreateProfileParameter(pp tc.ProfileParameter) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIProfileParameters, pp, nil, &alerts)
	return alerts, reqInf, err
}

// CreateMultipleProfileParameters creates multiple ProfileParameters at once.
func (to *Session) CreateMultipleProfileParameters(pps []tc.ProfileParameter) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIProfileParameters, pps, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) GetProfileParametersWithHdr(header http.Header) ([]tc.ProfileParameter, toclientlib.ReqInf, error) {
	return to.GetProfileParameterByQueryParamsWithHdr("", header)
}

// Returns a list of Profile Parameters
// Deprecated: GetProfileParameters will be removed in 6.0. Use GetProfileParametersWithHdr.
func (to *Session) GetProfileParameters() ([]tc.ProfileParameter, toclientlib.ReqInf, error) {
	return to.GetProfileParametersWithHdr(nil)
}

func (to *Session) GetProfileParameterByQueryParamsWithHdr(queryParams string, header http.Header) ([]tc.ProfileParameter, toclientlib.ReqInf, error) {
	URI := APIProfileParameters + queryParams
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

// GET a Profile Parameter by the Parameter
// Deprecated: GetProfileParameterByQueryParams will be removed in 6.0. Use GetProfileParameterByQueryParamsWithHdr.
func (to *Session) GetProfileParameterByQueryParams(queryParams string) ([]tc.ProfileParameter, toclientlib.ReqInf, error) {
	return to.GetProfileParameterByQueryParamsWithHdr(queryParams, nil)
}

// DELETE a Parameter by Parameter
func (to *Session) DeleteParameterByProfileParameter(profile int, parameter int) (tc.Alerts, toclientlib.ReqInf, error) {
	URI := fmt.Sprintf("%s/%d/%d", APIProfileParameters, profile, parameter)
	var alerts tc.Alerts
	reqInf, err := to.del(URI, nil, &alerts)
	return alerts, reqInf, err
}
