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
	// APIDivisions is the API version-relative path to the /divisions API route.
	APIDivisions = "/divisions"
)

// CreateDivision creates the given Division.
func (to *Session) CreateDivision(division tc.Division) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIDivisions, division, nil, &alerts)
	return alerts, reqInf, err
}

// UpdateDivision replaces the Division identified by 'id' with the one
// provided.
func (to *Session) UpdateDivision(id int, division tc.Division, header http.Header) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APIDivisions, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, division, header, &alerts)
	return alerts, reqInf, err
}

// GetDivisions returns all Divisions in Traffic Ops.
func (to *Session) GetDivisions(params url.Values, header http.Header) ([]tc.Division, toclientlib.ReqInf, error) {
	uri := APIDivisions
	if params != nil {
		uri += "?" + params.Encode()
	}
	var data tc.DivisionsResponse
	reqInf, err := to.get(uri, header, &data)
	return data.Response, reqInf, err
}

// GetDivisionByID retrieves the Division with the given ID.
func (to *Session) GetDivisionByID(id int, header http.Header) ([]tc.Division, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIDivisions, id)
	var data tc.DivisionsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetDivisionByName retrieves the Division with the given Name.
func (to *Session) GetDivisionByName(name string, header http.Header) ([]tc.Division, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", APIDivisions, url.QueryEscape(name))
	var data tc.DivisionsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// DeleteDivision deletes the Division identified by 'id'.
func (to *Session) DeleteDivision(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s/%d", APIDivisions, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
