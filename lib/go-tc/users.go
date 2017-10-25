package tc

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// UsersResponse ...
type UsersResponse struct {
	Response []User `json:"response"`
}

// User contains information about a given user in Traffic Ops.
type User struct {
	Username     string `json:"username,omitempty"`
	PublicSSHKey string `json:"publicSshKey,omitempty"`
	Role         int    `json:"role,omitempty"`
	RoleName     string `json:"rolename,omitempty"`
	ID           int    `json:"id,omitempty"`
	UID          int    `json:"uid,omitempty"`
	GID          int    `json:"gid,omitempty"`
	Company      string `json:"company,omitempty"`
	Email        string `json:"email,omitempty"`
	FullName     string `json:"fullName,omitempty"`
	NewUser      bool   `json:"newUser,omitempty"`
	LastUpdated  string `json:"lastUpdated,omitempty"`
}

// Credentials contains Traffic Ops login credentials
type UserCredentials struct {
	Username string `json:"u"`
	Password string `json:"p"`
}
