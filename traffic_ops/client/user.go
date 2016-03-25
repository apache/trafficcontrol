/*
   Copyright 2015 Comcast Cable Communications Management, LLC

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

	"github.com/cihub/seelog"
)

// UserResponse ...
type UserResponse struct {
	Version  string `json:"version"`
	Response []User `json:"response"`
}

// User contains information about a given user in Traffic Ops.
type User struct {
	Username     string `json:"username"`
	PublicSSHKey string `json:"public_ssh_key"`
	Role         string `json:"role"`
	UID          string `json:"uid"`
	GID          string `json:"gid"`
	Company      string `json:"company"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	NewUser      bool   `json:"new_user"`
	LastUpdated  string `json:"lastUpdated"`
}

// Users gets an array of Users.
func (to *Session) Users() ([]User, error) {
	url := "/api/1.2/users.json"
	resp, err := to.request(url, nil)
	if err != nil {
		seelog.Error(err)
		return nil, err
	}

	var data UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		seelog.Error(err)
		return nil, err
	}

	return data.Response, nil
}
