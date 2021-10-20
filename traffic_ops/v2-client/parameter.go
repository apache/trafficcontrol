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
	"net/url"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

const (
	API_PARAMETERS = apiBase + "/parameters"
)

// CreateParameter performs a POST to create a Parameter.
func (to *Session) CreateParameter(pl tc.Parameter) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(pl)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_PARAMETERS, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// CreateMultipleParameters performs a POST to create multiple Parameters at once.
func (to *Session) CreateMultipleParameters(pls []tc.Parameter) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(pls)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_PARAMETERS, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// UpdateParameterByID performs a PUT to update a Parameter by ID.
func (to *Session) UpdateParameterByID(id int, pl tc.Parameter) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(pl)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s/%d", API_PARAMETERS, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// GetParameters returns a list of Parameters.
func (to *Session) GetParameters() ([]tc.Parameter, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_PARAMETERS, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ParametersResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// GetParameterByID GETs a Parameter by the Parameter ID.
func (to *Session) GetParameterByID(id int) ([]tc.Parameter, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_PARAMETERS, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ParametersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetParameterByName GETs a Parameter by the Parameter name.
func (to *Session) GetParameterByName(name string) ([]tc.Parameter, ReqInf, error) {
	URI := API_PARAMETERS + "?name=" + url.QueryEscape(name)
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ParametersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetParameterByConfigFile GETs a Parameter by the Parameter ConfigFile.
func (to *Session) GetParameterByConfigFile(configFile string) ([]tc.Parameter, ReqInf, error) {
	URI := API_PARAMETERS + "?configFile=" + url.QueryEscape(configFile)
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ParametersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetParameterByNameAndConfigFile GETs a Parameter by the Parameter Name and ConfigFile.
func (to *Session) GetParameterByNameAndConfigFile(name string, configFile string) ([]tc.Parameter, ReqInf, error) {
	URI := fmt.Sprintf("%s?name=%s&configFile=%s", API_PARAMETERS, url.QueryEscape(name), url.QueryEscape(configFile))
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ParametersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetParameterByNameAndConfigFileAndValue GETs a Parameter by the Parameter Name and ConfigFile and Value.
// TODO: API should support all 3, but does not support filter by value
// currently. Until then, loop through hits until you find one with that value.
func (to *Session) GetParameterByNameAndConfigFileAndValue(name, configFile, value string) ([]tc.Parameter, ReqInf, error) {
	params, reqInf, err := to.GetParameterByNameAndConfigFile(name, configFile)
	if err != nil {
		return params, reqInf, err
	}
	for _, p := range params {
		if p.Value == value {
			return []tc.Parameter{p}, reqInf, err
		}
	}
	return nil, reqInf, err
}

// DeleteParameterByID DELETEs a Parameter by ID.
func (to *Session) DeleteParameterByID(id int) (tc.Alerts, ReqInf, error) {
	URI := fmt.Sprintf("%s/%d", API_PARAMETERS, id)
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
