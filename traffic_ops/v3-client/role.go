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
	// API_ROLES is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_ROLES = apiBase + "/roles"

	APIRoles = "/roles"
)

// CreateRole creates a Role.
func (to *Session) CreateRole(role tc.Role) (tc.Alerts, toclientlib.ReqInf, int, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIRoles, role, nil, &alerts)
	return alerts, reqInf, reqInf.StatusCode, err
}

func (to *Session) UpdateRoleByIDWithHdr(id int, role tc.Role, header http.Header) (tc.Alerts, toclientlib.ReqInf, int, error) {
	route := fmt.Sprintf("%s/?id=%d", APIRoles, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, role, header, &alerts)
	return alerts, reqInf, reqInf.StatusCode, err
}

// UpdateRoleByID updates a Role by ID.
// Deprecated: UpdateRoleByID will be removed in 6.0. Use UpdateRoleByIDWithHdr.
func (to *Session) UpdateRoleByID(id int, role tc.Role) (tc.Alerts, toclientlib.ReqInf, int, error) {

	return to.UpdateRoleByIDWithHdr(id, role, nil)
}

func (to *Session) GetRolesWithHdr(header http.Header) ([]tc.Role, toclientlib.ReqInf, int, error) {
	var data tc.RolesResponse
	reqInf, err := to.get(APIRoles, header, &data)
	return data.Response, reqInf, reqInf.StatusCode, err
}

// GetRoles returns a list of roles.
// Deprecated: GetRoles will be removed in 6.0. Use GetRolesWithHdr.
func (to *Session) GetRoles() ([]tc.Role, toclientlib.ReqInf, int, error) {
	return to.GetRolesWithHdr(nil)
}

func (to *Session) GetRoleByIDWithHdr(id int, header http.Header) ([]tc.Role, toclientlib.ReqInf, int, error) {
	route := fmt.Sprintf("%s/?id=%d", APIRoles, id)
	var data tc.RolesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, reqInf.StatusCode, err
}

// GetRoleByID GETs a Role by the Role ID.
// Deprecated: GetRoleByID will be removed in 6.0. Use GetRoleByIDWithHdr.
func (to *Session) GetRoleByID(id int) ([]tc.Role, toclientlib.ReqInf, int, error) {
	return to.GetRoleByIDWithHdr(id, nil)
}

func (to *Session) GetRoleByNameWithHdr(name string, header http.Header) ([]tc.Role, toclientlib.ReqInf, int, error) {
	route := fmt.Sprintf("%s?name=%s", APIRoles, url.QueryEscape(name))
	var data tc.RolesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, reqInf.StatusCode, err
}

// GetRoleByName GETs a Role by the Role name.
// Deprecated: GetRoleByName will be removed in 6.0. Use GetRoleByNameWithHdr.
func (to *Session) GetRoleByName(name string) ([]tc.Role, toclientlib.ReqInf, int, error) {
	return to.GetRoleByNameWithHdr(name, nil)
}

func (to *Session) GetRoleByQueryParamsWithHdr(queryParams map[string]string, header http.Header) ([]tc.Role, toclientlib.ReqInf, int, error) {
	route := fmt.Sprintf("%s%s", APIRoles, mapToQueryParameters(queryParams))
	var data tc.RolesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, reqInf.StatusCode, err
}

// GetRoleByQueryParams gets a Role by the Role query parameters.
// Deprecated: GetRoleByQueryParams will be removed in 6.0. Use GetRoleByQueryParamsWithHdr.
func (to *Session) GetRoleByQueryParams(queryParams map[string]string) ([]tc.Role, toclientlib.ReqInf, int, error) {
	return to.GetRoleByQueryParamsWithHdr(queryParams, nil)
}

// DeleteRoleByID DELETEs a Role by ID.
func (to *Session) DeleteRoleByID(id int) (tc.Alerts, toclientlib.ReqInf, int, error) {
	route := fmt.Sprintf("%s/?id=%d", APIRoles, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, reqInf.StatusCode, err
}
