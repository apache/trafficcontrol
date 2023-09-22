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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	// API_PARAMETERS is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_PARAMETERS = apiBase + "/parameters"

	APIParameters = "/parameters"
)

// CreateParameter performs a POST to create a Parameter.
func (to *Session) CreateParameter(pl tc.Parameter) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIParameters, pl, nil, &alerts)
	return alerts, reqInf, err
}

// CreateMultipleParameters performs a POST to create multiple Parameters at once.
func (to *Session) CreateMultipleParameters(pls []tc.Parameter) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIParameters, pls, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdateParameterByIDWithHdr(id int, pl tc.Parameter, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APIParameters, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, pl, header, &alerts)
	return alerts, reqInf, err
}

// UpdateParameterByID performs a PUT to update a Parameter by ID.
// Deprecated: UpdateParameterByID will be removed in 6.0. Use UpdateParameterByIDWithHdr.
func (to *Session) UpdateParameterByID(id int, pl tc.Parameter) (tc.Alerts, toclientlib.ReqInf, error) {
	return to.UpdateParameterByIDWithHdr(id, pl, nil)
}

func (to *Session) GetParametersWithHdr(header http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
	var data tc.ParametersResponse
	reqInf, err := to.get(APIParameters, header, &data)
	return data.Response, reqInf, err
}

// GetParameters returns a list of Parameters.
// Deprecated: GetParameters will be removed in 6.0. Use GetParametersWithHdr.
func (to *Session) GetParameters() ([]tc.Parameter, toclientlib.ReqInf, error) {
	return to.GetParametersWithHdr(nil)
}

func (to *Session) GetParameterByIDWithHdr(id int, header http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIParameters, id)
	var data tc.ParametersResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetParameterByID GETs a Parameter by the Parameter ID.
// Deprecated: GetParameterByID will be removed in 6.0. Use GetParameterByIDWithHdr.
func (to *Session) GetParameterByID(id int) ([]tc.Parameter, toclientlib.ReqInf, error) {
	return to.GetParameterByIDWithHdr(id, nil)
}

func (to *Session) GetParameterByNameWithHdr(name string, header http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
	uri := APIParameters + "?name=" + url.QueryEscape(name)
	var data tc.ParametersResponse
	reqInf, err := to.get(uri, header, &data)
	return data.Response, reqInf, err
}

// GetParameterByName GETs a Parameter by the Parameter name.
// Deprecated: GetParameterByName will be removed in 6.0. Use GetParameterByNameWithHdr.
func (to *Session) GetParameterByName(name string) ([]tc.Parameter, toclientlib.ReqInf, error) {
	return to.GetParameterByNameWithHdr(name, nil)
}

func (to *Session) GetParameterByConfigFileWithHdr(configFile string, header http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
	uri := APIParameters + "?configFile=" + url.QueryEscape(configFile)
	var data tc.ParametersResponse
	reqInf, err := to.get(uri, header, &data)
	return data.Response, reqInf, err
}

// GetParameterByConfigFile GETs a Parameter by the Parameter ConfigFile.
// Deprecated: GetParameterByConfigFile will be removed in 6.0. Use GetParameterByConfigFileWithHdr.
func (to *Session) GetParameterByConfigFile(configFile string) ([]tc.Parameter, toclientlib.ReqInf, error) {
	return to.GetParameterByConfigFileWithHdr(configFile, nil)
}

func (to *Session) GetParameterByNameAndConfigFileWithHdr(name string, configFile string, header http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
	uri := fmt.Sprintf("%s?name=%s&configFile=%s", APIParameters, url.QueryEscape(name), url.QueryEscape(configFile))
	var data tc.ParametersResponse
	reqInf, err := to.get(uri, header, &data)
	return data.Response, reqInf, err
}

// GetParameterByNameAndConfigFile GETs a Parameter by the Parameter Name and ConfigFile.
// Deprecated: GetParameterByNameAndConfigFile will be removed in 6.0. Use GetParameterByNameAndConfigFileWithHdr.
func (to *Session) GetParameterByNameAndConfigFile(name string, configFile string) ([]tc.Parameter, toclientlib.ReqInf, error) {
	return to.GetParameterByNameAndConfigFileWithHdr(name, configFile, nil)
}

func (to *Session) GetParameterByNameAndConfigFileAndValueWithHdr(name, configFile, value string, header http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
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
func (to *Session) GetParameterByNameAndConfigFileAndValue(name, configFile, value string) ([]tc.Parameter, toclientlib.ReqInf, error) {
	return to.GetParameterByNameAndConfigFileAndValueWithHdr(name, configFile, value, nil)
}

// DeleteParameterByID DELETEs a Parameter by ID.
func (to *Session) DeleteParameterByID(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	uri := fmt.Sprintf("%s/%d", APIParameters, id)
	var alerts tc.Alerts
	reqInf, err := to.del(uri, nil, &alerts)
	return alerts, reqInf, err
}
