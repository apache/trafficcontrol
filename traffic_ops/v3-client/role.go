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
	API_ROLES = apiBase + "/roles"
)

// CreateRole creates a Role.
func (to *Session) CreateRole(region tc.Role) (tc.Alerts, ReqInf, int, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(region)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, 0, err
	}
	resp, remoteAddr, errClient := to.RawRequest(http.MethodPost, API_ROLES, reqBody)
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

// UpdateRoleByID updates a Role by ID.
func (to *Session) UpdateRoleByID(id int, region tc.Role) (tc.Alerts, ReqInf, int, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(region)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, 0, err
	}
	route := fmt.Sprintf("%s/?id=%d", API_ROLES, id)
	resp, remoteAddr, errClient := to.RawRequest(http.MethodPut, route, reqBody)
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

func (to *Session) GetRolesWithHdr(header http.Header) ([]tc.Role, ReqInf, int, error) {
	resp, remoteAddr, errClient := to.RawRequestWithHdr(http.MethodGet, API_ROLES, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Role{}, reqInf, http.StatusNotModified, nil
		}
	}
	if resp != nil {
		defer resp.Body.Close()

		var data tc.RolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return data.Response, reqInf, resp.StatusCode, err
		}
		return data.Response, reqInf, resp.StatusCode, errClient
	}
	return []tc.Role{}, reqInf, 0, errClient
}

// GetRoles returns a list of roles.
// Deprecated: GetRoles will be removed in 6.0. Use GetRolesWithHdr.
func (to *Session) GetRoles() ([]tc.Role, ReqInf, int, error) {
	return to.GetRolesWithHdr(nil)
}

func (to *Session) GetRoleByIDWithHdr(id int, header http.Header) ([]tc.Role, ReqInf, int, error) {
	route := fmt.Sprintf("%s/?id=%d", API_ROLES, id)
	resp, remoteAddr, errClient := to.RawRequestWithHdr(http.MethodGet, route, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Role{}, reqInf, http.StatusNotModified, nil
		}
	}
	if resp != nil {
		defer resp.Body.Close()

		var data tc.RolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return data.Response, reqInf, resp.StatusCode, err
		}
		return data.Response, reqInf, resp.StatusCode, errClient
	}
	return []tc.Role{}, reqInf, 0, errClient
}

// GetRoleByID GETs a Role by the Role ID.
// Deprecated: GetRoleByID will be removed in 6.0. Use GetRoleByIDWithHdr.
func (to *Session) GetRoleByID(id int) ([]tc.Role, ReqInf, int, error) {
	return to.GetRoleByIDWithHdr(id, nil)
}

func (to *Session) GetRoleByNameWithHdr(name string, header http.Header) ([]tc.Role, ReqInf, int, error) {
	route := fmt.Sprintf("%s?name=%s", API_ROLES, url.QueryEscape(name))
	resp, remoteAddr, errClient := to.RawRequestWithHdr(http.MethodGet, route, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Role{}, reqInf, http.StatusNotModified, nil
		}
	}
	if resp != nil {
		defer resp.Body.Close()

		var data tc.RolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return data.Response, reqInf, resp.StatusCode, err
		}
		return data.Response, reqInf, resp.StatusCode, errClient
	}
	return []tc.Role{}, reqInf, 0, errClient
}

// GetRoleByName GETs a Role by the Role name.
// Deprecated: GetRoleByName will be removed in 6.0. Use GetRoleByNameWithHdr.
func (to *Session) GetRoleByName(name string) ([]tc.Role, ReqInf, int, error) {
	return to.GetRoleByNameWithHdr(name, nil)
}

func (to *Session) GetRoleByQueryParamsWithHdr(queryParams map[string]string, header http.Header) ([]tc.Role, ReqInf, int, error) {
	route := fmt.Sprintf("%s?", API_ROLES)
	for param, val := range queryParams {
		route += fmt.Sprintf("%s=%s&", url.QueryEscape(param), url.QueryEscape(val))
	}
	resp, remoteAddr, errClient := to.RawRequestWithHdr(http.MethodGet, route, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
		if reqInf.StatusCode == http.StatusNotModified {
			return []tc.Role{}, reqInf, http.StatusNotModified, nil
		}
	}
	if resp != nil {
		defer resp.Body.Close()

		var data tc.RolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return data.Response, reqInf, resp.StatusCode, err
		}
		return data.Response, reqInf, resp.StatusCode, errClient
	}
	return []tc.Role{}, reqInf, 0, errClient
}

// GetRoleByQueryParams gets a Role by the Role query parameters.
// Deprecated: GetRoleByQueryParams will be removed in 6.0. Use GetRoleByQueryParamsWithHdr.
func (to *Session) GetRoleByQueryParams(queryParams map[string]string) ([]tc.Role, ReqInf, int, error) {
	return to.GetRoleByQueryParamsWithHdr(queryParams, nil)
}

// DeleteRoleByID DELETEs a Role by ID.
func (to *Session) DeleteRoleByID(id int) (tc.Alerts, ReqInf, int, error) {
	route := fmt.Sprintf("%s/?id=%d", API_ROLES, id)
	resp, remoteAddr, errClient := to.RawRequest(http.MethodDelete, route, nil)
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
