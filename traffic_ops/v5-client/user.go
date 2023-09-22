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
	"fmt"
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// UserCurrentResponseV4 is an alias to avoid client breaking changes. In-case of a minor or major version change, we replace the below alias with a new structure.
type UserCurrentResponseV4 = tc.UserCurrentResponseV4

// GetUsers retrieves all (Tenant-accessible) Users stored in Traffic Ops.
func (to *Session) GetUsers(opts RequestOptions) (tc.UsersResponseV4, toclientlib.ReqInf, error) {
	data := tc.UsersResponseV4{}
	route := "/users"
	inf, err := to.get(route, opts, &data)
	return data, inf, err
}

// GetUserCurrent retrieves the currently authenticated User.
func (to *Session) GetUserCurrent(opts RequestOptions) (UserCurrentResponseV4, toclientlib.ReqInf, error) {
	route := `/user/current`
	resp := UserCurrentResponseV4{}
	reqInf, err := to.get(route, opts, &resp)
	return resp, reqInf, err
}

// UpdateCurrentUser replaces the current user data with the provided tc.UserV4 structure.
func (to *Session) UpdateCurrentUser(u tc.UserV4, opts RequestOptions) (tc.UpdateUserResponseV4, toclientlib.ReqInf, error) {
	var clientResp tc.UpdateUserResponseV4
	reqInf, err := to.put("/user/current", opts, u, &clientResp)
	return clientResp, reqInf, err
}

// CreateUser creates the given user.
func (to *Session) CreateUser(user tc.UserV4, opts RequestOptions) (tc.CreateUserResponseV4, toclientlib.ReqInf, error) {
	if user.Tenant != nil {
		innerOpts := NewRequestOptions()
		innerOpts.QueryParameters.Set("name", *user.Tenant)
		tenant, _, err := to.GetTenants(innerOpts)
		if err != nil {
			return tc.CreateUserResponseV4{Alerts: tenant.Alerts}, toclientlib.ReqInf{}, fmt.Errorf("resolving Tenant name '%s' to an ID: %w", *user.Tenant, err)
		}
		if len(tenant.Response) < 1 {
			return tc.CreateUserResponseV4{Alerts: tenant.Alerts}, toclientlib.ReqInf{}, fmt.Errorf("no such Tenant: '%s'", *user.Tenant)
		}
		user.TenantID = *tenant.Response[0].ID
	}

	if user.Role != "" {
		innerOpts := NewRequestOptions()
		innerOpts.QueryParameters.Set("name", user.Role)
		roles, _, err := to.GetRoles(innerOpts)
		if err != nil {
			return tc.CreateUserResponseV4{Alerts: roles.Alerts}, toclientlib.ReqInf{}, fmt.Errorf("resolving Role name '%s' to an ID: %w", user.Role, err)
		}
		if len(roles.Response) == 0 {
			return tc.CreateUserResponseV4{Alerts: roles.Alerts}, toclientlib.ReqInf{}, fmt.Errorf("no such Role: '%s'", user.Role)
		}
		user.Role = roles.Response[0].Name
	}

	route := "/users"
	var clientResp tc.CreateUserResponseV4
	reqInf, err := to.post(route, opts, user, &clientResp)
	return clientResp, reqInf, err
}

// UpdateUser replaces the User identified by 'id' with the one provided.
func (to *Session) UpdateUser(id int, u tc.UserV4, opts RequestOptions) (tc.UpdateUserResponseV4, toclientlib.ReqInf, error) {
	route := "/users/" + strconv.Itoa(id)
	var clientResp tc.UpdateUserResponseV4
	reqInf, err := to.put(route, opts, u, &clientResp)
	return clientResp, reqInf, err
}

// DeleteUser deletes the User with the given ID.
func (to *Session) DeleteUser(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	route := "/users/" + strconv.Itoa(id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, opts, &alerts)
	return alerts, reqInf, err
}

// RegisterNewUser requests the registration of a new user with the given tenant ID and role ID,
// through their email.
func (to *Session) RegisterNewUser(tenantID uint, role string, email rfc.EmailAddress, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqBody := tc.UserRegistrationRequestV40{
		Email:    email,
		TenantID: tenantID,
		Role:     role,
	}
	reqInf, err := to.post("/users/register", opts, reqBody, &alerts)
	return alerts, reqInf, err
}
