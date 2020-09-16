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

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_v13_Parameters = "/api/1.3/parameters"
)

// Create a Parameter
func (to *Session) CreateParameter(pl tc.Parameter) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(pl)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_v13_Parameters, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Update a Parameter by ID
func (to *Session) UpdateParameterByID(id int, pl tc.Parameter) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(pl)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s/%d", API_v13_Parameters, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Returns a list of Parameters
func (to *Session) GetParameters() ([]tc.Parameter, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_v13_Parameters, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ParametersResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// Parameters gets an array of parameter structs for the profile given
// Deprecated: use GetParameters
func (to *Session) Parameters(profileName string) ([]tc.Parameter, error) {
	ps, _, err := to.GetParametersByProfileName(profileName)
	return ps, err
}

func (to *Session) GetParametersByProfileName(profileName string) ([]tc.Parameter, ReqInf, error) {
	url := fmt.Sprintf(API_v13_Parameters+"/profile/%s.json", profileName)
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil)
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

// GET a Parameter by the Parameter ID
func (to *Session) GetParameterByID(id int) ([]tc.Parameter, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_v13_Parameters, id)
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

// GET a Parameter by the Parameter name
func (to *Session) GetParameterByName(name string) ([]tc.Parameter, ReqInf, error) {
	URI := API_v13_Parameters + "?name=" + url.QueryEscape(name)
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

// GET a Parameter by the Parameter ConfigFile
func (to *Session) GetParameterByConfigFile(configFile string) ([]tc.Parameter, ReqInf, error) {
	URI := API_v13_Parameters + "?configFile=" + url.QueryEscape(configFile)
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

// GET a Parameter by the Parameter Name and ConfigFile
func (to *Session) GetParameterByNameAndConfigFile(name string, configFile string) ([]tc.Parameter, ReqInf, error) {
	URI := fmt.Sprintf("%s?name=%s&configFile=%s", API_v13_Parameters, url.QueryEscape(name), url.QueryEscape(configFile))
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

// GET a Parameter by the Parameter Name and ConfigFile and Value
// TODO: API should support all 3,  but does not support filter by value
// currently.  Until then, loop thru hits until you find one with that value
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

// DELETE a Parameter by ID
func (to *Session) DeleteParameterByID(id int) (tc.Alerts, ReqInf, error) {
	URI := fmt.Sprintf("%s/%d", API_v13_Parameters, id)
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
