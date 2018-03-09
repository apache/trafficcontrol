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

package v13

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

const (
	API_v13_Profile_Parameters = "/api/1.3/profile_parameters"
)

// Create a ProfileParameter
func (to *Session) CreateProfileParameter(pp tc.ProfileParameter) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(pp)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_v13_Profile_Parameters, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Update a Profile Parameter by Profile
func (to *Session) UpdateParameterByProfile(id int, pp tc.ProfileParameter) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(pp)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	URI := fmt.Sprintf("%s/%d", API_v13_Profile_Parameters, id)
	resp, remoteAddr, err := to.request(http.MethodPut, URI, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Returns a list of Profile Parameters
func (to *Session) GetProfileParameters() ([]tc.ProfileParameter, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_v13_Profile_Parameters, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ProfileParametersResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// GET a Profile Parameter by the Profile
func (to *Session) GetProfileParameterByIDs(profile int, parameter int) ([]tc.ProfileParameter, ReqInf, error) {
	URI := fmt.Sprintf("%s?profile=%d&parameter=%d", API_v13_Profile_Parameters, profile)
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
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
func (to *Session) GetProfileParameterByParameter(parameter int) ([]tc.ProfileParameter, ReqInf, error) {
	URI := API_v13_Profile_Parameters + "?parameter=" + strconv.Itoa(parameter)
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
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

// DELETE a Parameter by Parameter
func (to *Session) DeleteParameterByParameter(parameter int) (tc.Alerts, ReqInf, error) {
	URI := fmt.Sprintf("%s/%d", API_v13_Profile_Parameters, parameter)
	resp, remoteAddr, err := to.request(http.MethodDelete, URI, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}
