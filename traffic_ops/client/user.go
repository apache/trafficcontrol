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
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-rfc"
)

// Users gets an array of Users.
// Deprecated: use GetUsers
func (to *Session) Users() ([]tc.User, error) {
	us, _, err := to.GetUsers()
	return us, err
}

// GetUsers returns all users accessible from current user
func (to *Session) GetUsers() ([]tc.User, ReqInf, error) {
	route := apiBase + "/users.json"
	resp, remoteAddr, err := to.request("GET", route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.UsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

func (to *Session) GetUserByID(id int) ([]tc.User, ReqInf, error) {
	data := tc.UsersResponse{}
	route := fmt.Sprintf("%s/users/%d", apiBase, id)
	inf, err := get(to, route, &data)
	return data.Response, inf, err
}

func (to *Session) GetUserByUsername(username string) ([]tc.User, ReqInf, error) {
	data := tc.UsersResponse{}
	route := fmt.Sprintf("%s/users?username=%s", apiBase, username)
	inf, err := get(to, route, &data)
	return data.Response, inf, err
}

// GetUserCurrent gets information about the current user
func (to *Session) GetUserCurrent() (*tc.UserCurrent, ReqInf, error) {
	route := apiBase + `/user/current`
	resp := tc.UserCurrentResponse{}
	reqInf, err := get(to, route, &resp)
	if err != nil {
		return nil, reqInf, err
	}
	return &resp.Response, reqInf, nil
}

// CreateUser creates a user
func (to *Session) CreateUser(user *tc.User) (*tc.CreateUserResponse, ReqInf, error) {
	if user.TenantID == nil && user.Tenant != nil {
		tenant, _, err := to.TenantByName(*user.Tenant)
		if err != nil {
			return nil, ReqInf{}, err
		}
		if tenant == nil {
			return nil, ReqInf{}, errors.New("no tenant with name " + *user.Tenant)
		}
		if err != nil {
			return nil, ReqInf{}, err
		}
		user.TenantID = &tenant.ID
	}

	if user.RoleName != nil && *user.RoleName != "" {
		roles, _, _, err := to.GetRoleByName(*user.RoleName)
		if err != nil {
			return nil, ReqInf{}, err
		}
		if len(roles) == 0 || roles[0].ID == nil {
			return nil, ReqInf{}, errors.New("no role with name " + *user.RoleName)
		}
		if err != nil {
			return nil, ReqInf{}, err
		}
		user.Role = roles[0].ID
	}

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(user)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	route := apiBase + "/users.json"
	resp, remoteAddr, err := to.request(http.MethodPost, route, reqBody)
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()
	var clientResp tc.CreateUserResponse
	err = json.NewDecoder(resp.Body).Decode(&clientResp)
	return &clientResp, reqInf, nil
}

// UpdateUserByID updates user with the given id
func (to *Session) UpdateUserByID(id int, u *tc.User) (*tc.UpdateUserResponse, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(u)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	route := apiBase + "/users/" + strconv.Itoa(id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()
	var clientResp tc.UpdateUserResponse
	err = json.NewDecoder(resp.Body).Decode(&clientResp)
	return &clientResp, reqInf, nil
}

// DeleteUserByID updates user with the given id
func (to *Session) DeleteUserByID(id int) (tc.Alerts, ReqInf, error) {
	route := apiBase + "/users/" + strconv.Itoa(id)
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

// RegisterNewUser requests the registration of a new user with the given tenant ID and role ID,
// through their email.
func (to *Session) RegisterNewUser(tenantID uint, roleID uint, email rfc.EmailAddress) (tc.Alerts, ReqInf, error) {
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss}
	var alerts tc.Alerts

	reqBody, err := json.Marshal(tc.UserRegistrationRequest{
		Email: email,
		TenantID: tenantID,
		Role: roleID,
	})
	if err != nil {
		return alerts, reqInf, err
	}

	resp, remoteAddr, err := to.request(http.MethodPost, apiBase+"/users/register", reqBody)
	reqInf.RemoteAddr = remoteAddr
	if err != nil {
		return alerts, reqInf, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, err
}
