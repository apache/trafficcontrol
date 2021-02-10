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
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// GetUsers retrieves all (Tenant-accessible) Users stored in Traffic Ops.
func (to *Session) GetUsers(header http.Header) ([]tc.User, toclientlib.ReqInf, error) {
	data := tc.UsersResponse{}
	route := "/users"
	inf, err := to.get(route, header, &data)
	return data.Response, inf, err
}

// GetUsersByRole retrieves all (Tenant-accessible) with the Role that has the
// given Name.
func (to *Session) GetUsersByRole(roleName string, header http.Header) ([]tc.User, toclientlib.ReqInf, error) {
	data := tc.UsersResponse{}
	route := "/users?role=" + url.QueryEscape(roleName)
	inf, err := to.get(route, header, &data)
	return data.Response, inf, err
}

// GetUserByID retrieves the User with the given ID.
func (to *Session) GetUserByID(id int, header http.Header) ([]tc.User, toclientlib.ReqInf, error) {
	data := tc.UsersResponse{}
	route := fmt.Sprintf("/users/%d", id)
	inf, err := to.get(route, header, &data)
	return data.Response, inf, err
}

// GetUserByUsername retrieves the User with the given Username.
func (to *Session) GetUserByUsername(username string, header http.Header) ([]tc.User, toclientlib.ReqInf, error) {
	data := tc.UsersResponse{}
	route := "/users?username=" + url.QueryEscape(username)
	inf, err := to.get(route, header, &data)
	return data.Response, inf, err
}

// GetUserCurrent retrieves the currently authenticated User.
func (to *Session) GetUserCurrent(header http.Header) (tc.UserCurrent, toclientlib.ReqInf, error) {
	route := `/user/current`
	resp := tc.UserCurrentResponse{}
	reqInf, err := to.get(route, header, &resp)
	return resp.Response, reqInf, err
}

// UpdateCurrentUser replaces the current user data with the provided tc.User structure.
func (to *Session) UpdateCurrentUser(u tc.User) (tc.UpdateUserResponse, toclientlib.ReqInf, error) {
	user := struct {
		User tc.User `json:"user"`
	}{u}
	var clientResp tc.UpdateUserResponse
	reqInf, err := to.put("/user/current", user, nil, &clientResp)
	return clientResp, reqInf, err
}

// CreateUser creates the given user
func (to *Session) CreateUser(user tc.User) (tc.CreateUserResponse, toclientlib.ReqInf, error) {
	if user.TenantID == nil && user.Tenant != nil {
		tenant, _, err := to.GetTenantByName(*user.Tenant, nil)
		if err != nil {
			return tc.CreateUserResponse{}, toclientlib.ReqInf{}, err
		}
		user.TenantID = &tenant.ID
	}

	if user.RoleName != nil && *user.RoleName != "" {
		roles, _, _, err := to.GetRoleByName(*user.RoleName, nil)
		if err != nil {
			return tc.CreateUserResponse{}, toclientlib.ReqInf{}, err
		}
		if len(roles) == 0 || roles[0].ID == nil {
			return tc.CreateUserResponse{}, toclientlib.ReqInf{}, errors.New("no role with name " + *user.RoleName)
		}
		user.Role = roles[0].ID
	}

	route := "/users"
	var clientResp tc.CreateUserResponse
	reqInf, err := to.post(route, user, nil, &clientResp)
	return clientResp, reqInf, err
}

// UpdateUser replaces the User identified by 'id' with the one provided.
func (to *Session) UpdateUser(id int, u tc.User) (tc.UpdateUserResponse, toclientlib.ReqInf, error) {
	route := "/users/" + strconv.Itoa(id)
	var clientResp tc.UpdateUserResponse
	reqInf, err := to.put(route, u, nil, &clientResp)
	return clientResp, reqInf, err
}

// DeleteUser deletes the User with the given ID.
func (to *Session) DeleteUser(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := "/users/" + strconv.Itoa(id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// RegisterNewUser requests the registration of a new user with the given tenant ID and role ID,
// through their email.
func (to *Session) RegisterNewUser(tenantID uint, roleID uint, email rfc.EmailAddress) (tc.Alerts, toclientlib.ReqInf, error) {
	var alerts tc.Alerts
	reqBody := tc.UserRegistrationRequest{
		Email:    email,
		TenantID: tenantID,
		Role:     roleID,
	}
	reqInf, err := to.post("/users/register", reqBody, nil, &alerts)
	return alerts, reqInf, err
}
