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
	API_v13_ROLES = "/api/1.3/roles"
)

// Create a Role
func (to *Session) CreateRole(region v13.Role) (tc.Alerts, ReqInf, int, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(region)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, 0, err
	}
	resp, remoteAddr, errClient := to.rawRequest(http.MethodPost, API_v13_ROLES, reqBody)
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

// Update a Role by ID
func (to *Session) UpdateRoleByID(id int, region v13.Role) (tc.Alerts, ReqInf, int, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(region)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, 0, err
	}
	route := fmt.Sprintf("%s/?id=%d", API_v13_ROLES, id)
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

// Returns a list of roles
func (to *Session) GetRoles() ([]v13.Role, ReqInf, int, error) {
	resp, remoteAddr, errClient := to.rawRequest(http.MethodGet, API_v13_ROLES, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		defer resp.Body.Close()

		var data v13.RolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return data.Response, reqInf, resp.StatusCode, err
		}
		return data.Response, reqInf, resp.StatusCode, errClient
	}
	return []v13.Role{}, reqInf, 0, errClient
}

// GET a Role by the Role id
func (to *Session) GetRoleByID(id int) ([]v13.Role, ReqInf, int, error) {
	route := fmt.Sprintf("%s/?id=%d", API_v13_ROLES, id)
	resp, remoteAddr, errClient := to.rawRequest(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		defer resp.Body.Close()

		var data v13.RolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return data.Response, reqInf, resp.StatusCode, err
		}
		return data.Response, reqInf, resp.StatusCode, errClient
	}
	return []v13.Role{}, reqInf, 0, errClient
}

// GET a Role by the Role name
func (to *Session) GetRoleByName(name string) ([]v13.Role, ReqInf, int, error) {
	url := fmt.Sprintf("%s?name=%s", API_v13_ROLES, name)
	resp, remoteAddr, errClient := to.rawRequest(http.MethodGet, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		defer resp.Body.Close()

		var data v13.RolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return data.Response, reqInf, resp.StatusCode, err
		}
		return data.Response, reqInf, resp.StatusCode, errClient
	}
	return []v13.Role{}, reqInf, 0, errClient
}

// DELETE a Role by ID
func (to *Session) DeleteRoleByID(id int) (tc.Alerts, ReqInf, int, error) {
	route := fmt.Sprintf("%s/?id=%d", API_v13_ROLES, id)
	resp, remoteAddr, errClient := to.rawRequest(http.MethodDelete, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		defer resp.Body.Close()

		var alerts tc.Alerts
		if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
			return alerts, reqInf, resp.StatusCode, err
		}
		return alerts, reqInf, resp.StatusCode, errClient
	}
	return tc.Alerts{}, reqInf, 0, errClient
}
