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
	resp, remoteAddr, err := to.request(http.MethodPost, API_PARAMETERS, reqBody, nil)
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
	resp, remoteAddr, err := to.request(http.MethodPost, API_PARAMETERS, reqBody, nil)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

func (to *Session) UpdateParameterByIDWithHdr(id int, pl tc.Parameter, header http.Header) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(pl)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s/%d", API_PARAMETERS, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody, header)
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
	}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// UpdateParameterByID performs a PUT to update a Parameter by ID.
// Deprecated: UpdateParameterByID will be removed in 6.0. Use UpdateParameterByIDWithHdr.
func (to *Session) UpdateParameterByID(id int, pl tc.Parameter) (tc.Alerts, ReqInf, error) {
	return to.UpdateParameterByIDWithHdr(id, pl, nil)
}

func (to *Session) GetParametersWithHdr(header http.Header) ([]tc.Parameter, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_PARAMETERS, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Parameter{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ParametersResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// GetParameters returns a list of Parameters.
// Deprecated: GetParameters will be removed in 6.0. Use GetParametersWithHdr.
func (to *Session) GetParameters() ([]tc.Parameter, ReqInf, error) {
	return to.GetParametersWithHdr(nil)
}

func (to *Session) GetParameterByIDWithHdr(id int, header http.Header) ([]tc.Parameter, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_PARAMETERS, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Parameter{}, reqInf, nil
		}
	}
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

// GetParameterByID GETs a Parameter by the Parameter ID.
// Deprecated: GetParameterByID will be removed in 6.0. Use GetParameterByIDWithHdr.
func (to *Session) GetParameterByID(id int) ([]tc.Parameter, ReqInf, error) {
	return to.GetParameterByIDWithHdr(id, nil)
}

func (to *Session) GetParameterByNameWithHdr(name string, header http.Header) ([]tc.Parameter, ReqInf, error) {
	URI := API_PARAMETERS + "?name=" + url.QueryEscape(name)
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Parameter{}, reqInf, nil
		}
	}
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
// Deprecated: GetParameterByName will be removed in 6.0. Use GetParameterByNameWithHdr.
func (to *Session) GetParameterByName(name string) ([]tc.Parameter, ReqInf, error) {
	return to.GetParameterByNameWithHdr(name, nil)
}

func (to *Session) GetParameterByConfigFileWithHdr(configFile string, header http.Header) ([]tc.Parameter, ReqInf, error) {
	URI := API_PARAMETERS + "?configFile=" + url.QueryEscape(configFile)
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Parameter{}, reqInf, nil
		}
	}
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
// Deprecated: GetParameterByConfigFile will be removed in 6.0. Use GetParameterByConfigFileWithHdr.
func (to *Session) GetParameterByConfigFile(configFile string) ([]tc.Parameter, ReqInf, error) {
	return to.GetParameterByConfigFileWithHdr(configFile, nil)
}

func (to *Session) GetParameterByNameAndConfigFileWithHdr(name string, configFile string, header http.Header) ([]tc.Parameter, ReqInf, error) {
	URI := fmt.Sprintf("%s?name=%s&configFile=%s", API_PARAMETERS, url.QueryEscape(name), url.QueryEscape(configFile))
	resp, remoteAddr, err := to.request(http.MethodGet, URI, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Parameter{}, reqInf, nil
		}
	}
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
// Deprecated: GetParameterByNameAndConfigFile will be removed in 6.0. Use GetParameterByNameAndConfigFileWithHdr.
func (to *Session) GetParameterByNameAndConfigFile(name string, configFile string) ([]tc.Parameter, ReqInf, error) {
	return to.GetParameterByNameAndConfigFileWithHdr(name, configFile, nil)
}

func (to *Session) GetParameterByNameAndConfigFileAndValueWithHdr(name, configFile, value string, header http.Header) ([]tc.Parameter, ReqInf, error) {
	params, reqInf, err := to.GetParameterByNameAndConfigFileWithHdr(name, configFile, header)
	if reqInf.StatusCode == http.StatusNotModified {
		return []tc.Parameter{}, reqInf, nil
	}
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

// GetParameterByNameAndConfigFileAndValue GETs a Parameter by the Parameter Name and ConfigFile and Value.
// TODO: API should support all 3, but does not support filter by value
// currently. Until then, loop through hits until you find one with that value.
// Deprecated: GetParameterByNameAndConfigFileAndValue will be removed in 6.0. Use GetParameterByNameAndConfigFileAndValueWithHdr.
func (to *Session) GetParameterByNameAndConfigFileAndValue(name, configFile, value string) ([]tc.Parameter, ReqInf, error) {
	return to.GetParameterByNameAndConfigFileAndValueWithHdr(name, configFile, value, nil)
}

// DeleteParameterByID DELETEs a Parameter by ID.
func (to *Session) DeleteParameterByID(id int) (tc.Alerts, ReqInf, error) {
	URI := fmt.Sprintf("%s/%d", API_PARAMETERS, id)
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
