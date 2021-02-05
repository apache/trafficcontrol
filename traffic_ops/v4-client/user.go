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
)

func (to *Session) GetUsersWithHdr(header http.Header) ([]tc.User, ReqInf, error) {
	data := tc.UsersResponse{}
	route := fmt.Sprintf("%s/users", apiBase)
	inf, err := to.get(route, header, &data)
	return data.Response, inf, err
}

// GetUsers returns all users accessible from current user
// Deprecated: GetUsers will be removed in 6.0. Use GetUsersWithHdr.
func (to *Session) GetUsers() ([]tc.User, ReqInf, error) {
	return to.GetUsersWithHdr(nil)
}

func (to *Session) GetUsersByRoleWithHdr(roleName string, header http.Header) ([]tc.User, ReqInf, error) {
	data := tc.UsersResponse{}
	route := fmt.Sprintf("%s/users?role=%s", apiBase, url.QueryEscape(roleName))
	inf, err := to.get(route, header, &data)
	return data.Response, inf, err
}

// GetUsersByRole returns all users accessible from current user for a given role
// Deprecated: GetUsersByRole will be removed in 6.0. Use GetUsersByRoleWithHdr.
func (to *Session) GetUsersByRole(roleName string) ([]tc.User, ReqInf, error) {
	return to.GetUsersByRoleWithHdr(roleName, nil)
}

func (to *Session) GetUserByIDWithHdr(id int, header http.Header) ([]tc.User, ReqInf, error) {
	data := tc.UsersResponse{}
	route := fmt.Sprintf("%s/users/%d", apiBase, id)
	inf, err := to.get(route, header, &data)
	return data.Response, inf, err
}

// Deprecated: GetUserByID will be removed in 6.0. Use GetUserByIDWithHdr.
func (to *Session) GetUserByID(id int) ([]tc.User, ReqInf, error) {
	return to.GetUserByIDWithHdr(id, nil)
}

func (to *Session) GetUserByUsernameWithHdr(username string, header http.Header) ([]tc.User, ReqInf, error) {
	data := tc.UsersResponse{}
	route := fmt.Sprintf("%s/users?username=%s", apiBase, url.QueryEscape(username))
	inf, err := to.get(route, header, &data)
	return data.Response, inf, err
}

// Deprecated: GetUserByUsername will be removed in 6.0. Use GetUserByUsernameWithHdr.
func (to *Session) GetUserByUsername(username string) ([]tc.User, ReqInf, error) {
	return to.GetUserByUsernameWithHdr(username, nil)
}

func (to *Session) GetUserCurrentWithHdr(header http.Header) (*tc.UserCurrent, ReqInf, error) {
	route := apiBase + `/user/current`
	resp := tc.UserCurrentResponse{}
	reqInf, err := to.get(route, header, &resp)
	if err != nil {
		return nil, reqInf, err
	}
	return &resp.Response, reqInf, nil
}

// GetUserCurrent gets information about the current user
// Deprecated: GetUserCurrent will be removed in 6.0. Use GetUserCurrentWithHdr.
func (to *Session) GetUserCurrent() (*tc.UserCurrent, ReqInf, error) {
	return to.GetUserCurrentWithHdr(nil)
}

// UpdateCurrentUser replaces the current user data with the provided tc.User structure.
func (to *Session) UpdateCurrentUser(u tc.User) (*tc.UpdateUserResponse, ReqInf, error) {
	user := struct {
		User tc.User `json:"user"`
	}{u}
	var clientResp tc.UpdateUserResponse
	reqInf, err := to.put(apiBase+"/user/current", user, nil, &clientResp)
	return &clientResp, reqInf, err
}

// CreateUser creates a user
func (to *Session) CreateUser(user *tc.User) (*tc.CreateUserResponse, ReqInf, error) {
	if user.TenantID == nil && user.Tenant != nil {
		tenant, _, err := to.TenantByNameWithHdr(*user.Tenant, nil)
		if err != nil {
			return nil, ReqInf{}, err
		}
		if tenant == nil {
			return nil, ReqInf{}, errors.New("no tenant with name " + *user.Tenant)
		}
		user.TenantID = &tenant.ID
	}

	if user.RoleName != nil && *user.RoleName != "" {
		roles, _, _, err := to.GetRoleByNameWithHdr(*user.RoleName, nil)
		if err != nil {
			return nil, ReqInf{}, err
		}
		if len(roles) == 0 || roles[0].ID == nil {
			return nil, ReqInf{}, errors.New("no role with name " + *user.RoleName)
		}
		user.Role = roles[0].ID
	}

	route := apiBase + "/users"
	var clientResp tc.CreateUserResponse
	reqInf, err := to.post(route, user, nil, &clientResp)
	return &clientResp, reqInf, err
}

// UpdateUserByID updates user with the given id
func (to *Session) UpdateUserByID(id int, u *tc.User) (*tc.UpdateUserResponse, ReqInf, error) {
	route := apiBase + "/users/" + strconv.Itoa(id)
	var clientResp tc.UpdateUserResponse
	reqInf, err := to.put(route, u, nil, &clientResp)
	return &clientResp, reqInf, err
}

// DeleteUserByID updates user with the given id
func (to *Session) DeleteUserByID(id int) (tc.Alerts, ReqInf, error) {
	route := apiBase + "/users/" + strconv.Itoa(id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// RegisterNewUser requests the registration of a new user with the given tenant ID and role ID,
// through their email.
func (to *Session) RegisterNewUser(tenantID uint, roleID uint, email rfc.EmailAddress) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	reqBody := tc.UserRegistrationRequest{
		Email:    email,
		TenantID: tenantID,
		Role:     roleID,
	}
	reqInf, err := to.post(apiBase+"/users/register", reqBody, nil, &alerts)
	return alerts, reqInf, err
}
