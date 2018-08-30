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
	API_StaticDNSEntries = apiBase + "/staticdnsentries"
)

// Create a StaticDNSEntry
func (to *Session) CreateStaticDNSEntry(cdn StaticDNSEntry) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(cdn)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_StaticDNSEntries, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Update a StaticDNSEntry by ID
func (to *Session) UpdateStaticDNSEntryByID(id int, cdn StaticDNSEntry) (tc.Alerts, ReqInf, int, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(cdn)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, 0, err
	}
	route := fmt.Sprintf("%s?id=%d", API_StaticDNSEntries, id)
	resp, remoteAddr, errClient := to.rawRequest(http.MethodPut, route, reqBody)
	if resp != nil {
		defer resp.Body.Close()
		var alerts tc.Alerts
		if err = json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
			return alerts, reqInf, resp.StatusCode, err
		}
		return alerts, reqInf, resp.StatusCode, errClient
	}
	return tc.Alerts{}, reqInf, 0, errClient
}

// Returns a list of StaticDNSEntrys
func (to *Session) GetStaticDNSEntries() ([]StaticDNSEntry, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_StaticDNSEntries, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data StaticDNSEntriesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// GET a StaticDNSEntry by the StaticDNSEntry ID
func (to *Session) GetStaticDNSEntryByID(id int) ([]StaticDNSEntry, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_StaticDNSEntries, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data StaticDNSEntriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a StaticDNSEntry by the StaticDNSEntry hsot
func (to *Session) GetStaticDNSEntriesByHost(host string) ([]StaticDNSEntry, ReqInf, error) {
	url := fmt.Sprintf("%s?host=%s", API_StaticDNSEntries, host)
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data StaticDNSEntriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// DELETE a StaticDNSEntry by ID
func (to *Session) DeleteStaticDNSEntryByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_StaticDNSEntries, id)
	resp, remoteAddr, err := to.request(http.MethodDelete, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}
