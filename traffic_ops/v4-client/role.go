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
	// APIRoles is the full path to the /roles API endpoint.
	APIRoles = "/roles"
)

// CreateRole creates the given Role.
func (to *Session) CreateRole(role tc.Role) (tc.Alerts, toclientlib.ReqInf, int, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIRoles, role, nil, &alerts)
	return alerts, reqInf, reqInf.StatusCode, err
}

// UpdateRole replaces the Role identified by 'id' with the one provided.
func (to *Session) UpdateRole(id int, role tc.Role, header http.Header) (tc.Alerts, toclientlib.ReqInf, int, error) {
	route := fmt.Sprintf("%s/?id=%d", APIRoles, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, role, header, &alerts)
	return alerts, reqInf, reqInf.StatusCode, err
}

// GetRoleByID returns the Role with the given ID.
func (to *Session) GetRoleByID(id int, header http.Header) ([]tc.Role, toclientlib.ReqInf, int, error) {
	route := fmt.Sprintf("%s/?id=%d", APIRoles, id)
	var data tc.RolesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, reqInf.StatusCode, err
}

// GetRoleByName returns the Role with the given Name.
func (to *Session) GetRoleByName(name string, header http.Header) ([]tc.Role, toclientlib.ReqInf, int, error) {
	route := fmt.Sprintf("%s?name=%s", APIRoles, url.QueryEscape(name))
	var data tc.RolesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, reqInf.StatusCode, err
}

// GetRoles retrieves Roles from Traffic Ops.
func (to *Session) GetRoles(queryParams url.Values, header http.Header) ([]tc.Role, toclientlib.ReqInf, int, error) {
	route := APIRoles
	if len(queryParams) > 0 {
		route += "?" + queryParams.Encode()
	}
	var data tc.RolesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, reqInf.StatusCode, err
}

// DeleteRole deletes the Role with the given ID.
func (to *Session) DeleteRole(id int) (tc.Alerts, toclientlib.ReqInf, int, error) {
	route := fmt.Sprintf("%s/?id=%d", APIRoles, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, reqInf.StatusCode, err
}
