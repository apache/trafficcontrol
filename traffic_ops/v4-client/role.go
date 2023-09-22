package client

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

import (
	"net/url"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiRoles is the full path to the /roles API endpoint.
const apiRoles = "/roles"

// CreateRole creates the given Role.
func (to *Session) CreateRole(role tc.RoleV4, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(apiRoles, opts, role, &alerts)
	return alerts, reqInf, err
}

// UpdateRole replaces the Role identified by 'id' with the one provided.
func (to *Session) UpdateRole(name string, role tc.RoleV4, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("name", name)
	var alerts tc.Alerts
	reqInf, err := to.put(apiRoles, opts, role, &alerts)
	return alerts, reqInf, err
}

// GetRoles retrieves Roles from Traffic Ops.
func (to *Session) GetRoles(opts RequestOptions) (tc.RolesResponseV4, toclientlib.ReqInf, error) {
	var data tc.RolesResponseV4
	reqInf, err := to.get(apiRoles, opts, &data)
	return data, reqInf, err
}

// DeleteRole deletes the Role with the given ID.
func (to *Session) DeleteRole(name string, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("name", name)
	var alerts tc.Alerts
	reqInf, err := to.del(apiRoles, opts, &alerts)
	return alerts, reqInf, err
}
