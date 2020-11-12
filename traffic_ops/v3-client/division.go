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
	API_DIVISIONS = apiBase + "/divisions"
)

// Create a Division
func (to *Session) CreateDivision(division tc.Division) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(division)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_DIVISIONS, reqBody, nil)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

func (to *Session) UpdateDivisionByIDWithHdr(id int, division tc.Division, header http.Header) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(division)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s/%d", API_DIVISIONS, id)
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

// Update a Division by ID
// Deprecated: UpdateDivisionByID will be removed in 6.0. Use UpdateDivisionByIDWithHdr.
func (to *Session) UpdateDivisionByID(id int, division tc.Division) (tc.Alerts, ReqInf, error) {
	return to.UpdateDivisionByIDWithHdr(id, division, nil)
}

func (to *Session) GetDivisionsWithHdr(header http.Header) ([]tc.Division, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_DIVISIONS, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Division{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.DivisionsResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// Returns a list of Divisions
// Deprecated: GetDivisions will be removed in 6.0. Use GetDivisionsWithHdr.
func (to *Session) GetDivisions() ([]tc.Division, ReqInf, error) {
	return to.GetDivisionsWithHdr(nil)
}

func (to *Session) GetDivisionByIDWithHdr(id int, header http.Header) ([]tc.Division, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_DIVISIONS, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Division{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.DivisionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a Division by the Division id
// Deprecated: GetDivisionByID will be removed in 6.0. Use GetDivisionByIDWithHdr.
func (to *Session) GetDivisionByID(id int) ([]tc.Division, ReqInf, error) {
	return to.GetDivisionByIDWithHdr(id, nil)
}

func (to *Session) GetDivisionByNameWithHdr(name string, header http.Header) ([]tc.Division, ReqInf, error) {
	url := fmt.Sprintf("%s?name=%s", API_DIVISIONS, name)
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Division{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.DivisionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a Division by the Division name
// Deprecated: GetDivisionByName will be removed in 6.0. Use GetDivisionByNameWithHdr.
func (to *Session) GetDivisionByName(name string) ([]tc.Division, ReqInf, error) {
	return to.GetDivisionByNameWithHdr(name, nil)
}

// DELETE a Division by Division id
func (to *Session) DeleteDivisionByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_DIVISIONS, id)
	resp, remoteAddr, err := to.request(http.MethodDelete, route, nil, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}
