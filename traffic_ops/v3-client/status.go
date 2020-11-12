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
	API_STATUSES = apiBase + "/statuses"
)

// CreateStatusNullable creates a new status, using the tc.StatusNullable structure.
func (to *Session) CreateStatusNullable(status tc.StatusNullable) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(status)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_STATUSES, reqBody, nil)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

func (to *Session) UpdateStatusByIDWithHdr(id int, status tc.Status, header http.Header) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(status)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s/%d", API_STATUSES, id)
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

// UpdateStatusByID updates a Status by ID.
// Deprecated: UpdateStatusByID will be removed in 6.0. Use UpdateStatusByIDWithHdr.
func (to *Session) UpdateStatusByID(id int, status tc.Status) (tc.Alerts, ReqInf, error) {
	return to.UpdateStatusByIDWithHdr(id, status, nil)
}

func (to *Session) GetStatusesWithHdr(header http.Header) ([]tc.Status, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_STATUSES, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Status{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.StatusesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// GetStatuses returns a list of Statuses.
// Deprecated: GetStatuses will be removed in 6.0. Use GetStatusesWithHdr.
func (to *Session) GetStatuses() ([]tc.Status, ReqInf, error) {
	return to.GetStatusesWithHdr(nil)
}

func (to *Session) GetStatusByIDWithHdr(id int, header http.Header) ([]tc.Status, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_STATUSES, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Status{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.StatusesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetStatusByID GETs a Status by the Status ID.
// Deprecated: GetStatusByID will be removed in 6.0. Use GetStatusByIDWithHdr.
func (to *Session) GetStatusByID(id int) ([]tc.Status, ReqInf, error) {
	return to.GetStatusByIDWithHdr(id, nil)
}

func (to *Session) GetStatusByNameWithHdr(name string, header http.Header) ([]tc.Status, ReqInf, error) {
	url := fmt.Sprintf("%s?name=%s", API_STATUSES, name)
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Status{}, reqInf, nil
		}
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.StatusesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetStatusByName GETs a Status by the Status name.
// Deprecated: GetStatusByName will be removed in 6.0. Use GetStatusByNameWithHdr.
func (to *Session) GetStatusByName(name string) ([]tc.Status, ReqInf, error) {
	return to.GetStatusByNameWithHdr(name, nil)
}

// DeleteStatusByID DELETEs a Status by ID.
func (to *Session) DeleteStatusByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_STATUSES, id)
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
