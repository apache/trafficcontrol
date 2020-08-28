/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

/**
 * Reprents a Traffic Ops user (`tm_user` in the database)
*/
export interface User {
	addressLine1?:    string;
	addressLine2?:    string;
	city?:            string;
	company?:         string;
	country?:         string;
	email?:           string;
	fullName?:        string;
	gid?:             number;
	id:               number;
	lastUpdated?:     Date;
	localUser?:       boolean;
	newUser:          boolean;
	phoneNumber?:     string;
	postalCode?:      string;
	publicSshKey?:    string;
	role?:            number;
	roleName?:        string;
	stateOrProvince?: string;
	tenant?:          string;
	tenantId?:        number;
	uid?:             number;
	username:         string;
}


/**
 * Represents a role that a user may have
*/
export interface Role {
	capabilities: Array<string>;
	description?: string;
	id:           number;
	name:         string;
	privLevel:    number;
}

/**
 * Represents a user's ability to perform some action
*/
export interface Capability {
	name:         string;
	description:  string;
	lastUpdated?: Date;
}
