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

const (
	// APIParameters is the full path to the /parameters API endpoint.
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

// UpdateParameter replaces the Parameter identified by 'id' with the one
// provided.
func (to *Session) UpdateParameter(id int, pl tc.Parameter, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APIParameters, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, pl, header, &alerts)
	return alerts, reqInf, err
}

// GetParameters returns all Parameters in Traffic Ops.
func (to *Session) GetParameters(header http.Header, params url.Values) ([]tc.Parameter, toclientlib.ReqInf, error) {
	route := APIParameters
	if len(params) > 0 {
		route += "?" + params.Encode()
	}

	var data tc.ParametersResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetParameterByID returns the Parameter with the given ID.
func (to *Session) GetParameterByID(id int, header http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIParameters, id)
	var data tc.ParametersResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetParametersByName retrieves all Parameters with the given Name.
func (to *Session) GetParametersByName(name string, header http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
	URI := APIParameters + "?name=" + url.QueryEscape(name)
	var data tc.ParametersResponse
	reqInf, err := to.get(URI, header, &data)
	return data.Response, reqInf, err
}

// GetParametersByConfigFile retrieves all Parameters that have the given
// Config File.
func (to *Session) GetParametersByConfigFile(configFile string, header http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
	URI := APIParameters + "?configFile=" + url.QueryEscape(configFile)
	var data tc.ParametersResponse
	reqInf, err := to.get(URI, header, &data)
	return data.Response, reqInf, err
}

// GetParametersByValue retrieves all Parameters that have the given Value.
func (to *Session) GetParametersByValue(value string, header http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
	URI := fmt.Sprintf("%s?value=%s", APIParameters, url.QueryEscape(value))
	var data tc.ParametersResponse
	reqInf, err := to.get(URI, header, &data)
	return data.Response, reqInf, err
}

// DeleteParameter deletes the Parameter with the given ID.
func (to *Session) DeleteParameter(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	URI := fmt.Sprintf("%s/%d", APIParameters, id)
	var alerts tc.Alerts
	reqInf, err := to.del(URI, nil, &alerts)
	return alerts, reqInf, err
}
