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

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_PROFILE_PARAMETERS = apiBase + "/profileparameters"
	ProfileIdQueryParam    = "profileId"
	ParameterIdQueryParam  = "parameterId"
)

// Create a ProfileParameter
func (to *Session) CreateProfileParameter(pp tc.ProfileParameter) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(pp)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_PROFILE_PARAMETERS, reqBody, nil)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// CreateMultipleProfileParameters creates multiple ProfileParameters at once.
func (to *Session) CreateMultipleProfileParameters(pps []tc.ProfileParameter) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(pps)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_PROFILE_PARAMETERS, reqBody, nil)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

func (to *Session) GetProfileParametersWithHdr(header http.Header) ([]tc.ProfileParameter, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_PROFILE_PARAMETERS, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.ProfileParameter{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfileParametersResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// Returns a list of Profile Parameters
// Deprecated: GetProfileParameters will be removed in 6.0. Use GetProfileParametersWithHdr.
func (to *Session) GetProfileParameters() ([]tc.ProfileParameter, ReqInf, error) {
	return to.GetProfileParametersWithHdr(nil)
}

func (to *Session) GetProfileParameterByQueryParamsWithHdr(queryParams string, header http.Header) ([]tc.ProfileParameter, ReqInf, error) {
	URI := API_PROFILE_PARAMETERS + queryParams
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.ProfileParameter{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfileParametersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a Profile Parameter by the Parameter
// Deprecated: GetProfileParameterByQueryParams will be removed in 6.0. Use GetProfileParameterByQueryParamsWithHdr.
func (to *Session) GetProfileParameterByQueryParams(queryParams string) ([]tc.ProfileParameter, ReqInf, error) {
	return to.GetProfileParameterByQueryParamsWithHdr(queryParams, nil)
}

// DELETE a Parameter by Parameter
func (to *Session) DeleteParameterByProfileParameter(profile int, parameter int) (tc.Alerts, ReqInf, error) {
	URI := fmt.Sprintf("%s/%d/%d", API_PROFILE_PARAMETERS, profile, parameter)
	resp, remoteAddr, err := to.request(http.MethodDelete, URI, nil, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}
