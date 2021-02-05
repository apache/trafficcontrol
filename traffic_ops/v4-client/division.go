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
)

const (
	API_DIVISIONS = apiBase + "/divisions"
)

// Create a Division
func (to *Session) CreateDivision(division tc.Division) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(API_DIVISIONS, division, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdateDivisionByIDWithHdr(id int, division tc.Division, header http.Header) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_DIVISIONS, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, division, header, &alerts)
	return alerts, reqInf, err
}

// Update a Division by ID
// Deprecated: UpdateDivisionByID will be removed in 6.0. Use UpdateDivisionByIDWithHdr.
func (to *Session) UpdateDivisionByID(id int, division tc.Division) (tc.Alerts, ReqInf, error) {
	return to.UpdateDivisionByIDWithHdr(id, division, nil)
}

func (to *Session) GetDivisionsWithHdr(header http.Header) ([]tc.Division, ReqInf, error) {
	var data tc.DivisionsResponse
	reqInf, err := to.get(API_DIVISIONS, header, &data)
	return data.Response, reqInf, err
}

// Returns a list of Divisions
// Deprecated: GetDivisions will be removed in 6.0. Use GetDivisionsWithHdr.
func (to *Session) GetDivisions() ([]tc.Division, ReqInf, error) {
	return to.GetDivisionsWithHdr(nil)
}

func (to *Session) GetDivisionByIDWithHdr(id int, header http.Header) ([]tc.Division, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_DIVISIONS, id)
	var data tc.DivisionsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GET a Division by the Division id
// Deprecated: GetDivisionByID will be removed in 6.0. Use GetDivisionByIDWithHdr.
func (to *Session) GetDivisionByID(id int) ([]tc.Division, ReqInf, error) {
	return to.GetDivisionByIDWithHdr(id, nil)
}

func (to *Session) GetDivisionByNameWithHdr(name string, header http.Header) ([]tc.Division, ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", API_DIVISIONS, url.QueryEscape(name))
	var data tc.DivisionsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GET a Division by the Division name
// Deprecated: GetDivisionByName will be removed in 6.0. Use GetDivisionByNameWithHdr.
func (to *Session) GetDivisionByName(name string) ([]tc.Division, ReqInf, error) {
	return to.GetDivisionByNameWithHdr(name, nil)
}

// DELETE a Division by Division id
func (to *Session) DeleteDivisionByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_DIVISIONS, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
