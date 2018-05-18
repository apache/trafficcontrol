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

package v13

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc/v13"
)

const (
	API_v13_IPs = "/api/1.3/ips"
)

// Create an IP
func (to *Session) CreateIP(serverID int, ip v13.IP) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(ip)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := API_v13_IPs
	resp, remoteAddr, err := to.request(http.MethodPost, route, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)

	return alerts, reqInf, nil
}

// Update an IP by ID
func (to *Session) UpdateIPs(id int, ip v13.IP) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(ip)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := API_v13_IPs + fmt.Sprintf("?id=%d", id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Returns a list of IPs
func (to *Session) GetIPsByServer(serverID int) ([]v13.IP, ReqInf, error) {
	route := API_v13_IPs + fmt.Sprintf("?serverId=%d", serverID)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data v13.IPsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// Returns an IP
func (to *Session) GetIPsById(id int) ([]v13.IP, ReqInf, error) {
	route := API_v13_IPs + fmt.Sprintf("?id=%d", id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data v13.IPsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// DELETE an IP
func (to *Session) DeleteIPs(id int) (tc.Alerts, ReqInf, error) {
	route := API_v13_IPs + fmt.Sprintf("?id=%d", id)
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
