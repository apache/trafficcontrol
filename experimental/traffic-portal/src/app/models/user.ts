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
	/** Line one of the user's address. */
	addressLine1?:    string;
	/** Line two of the user's address. */
	addressLine2?:    string;
	/** The city in which the user lives/is based. */
	city?:            string;
	/** The company for which the user works. */
	company?:         string;
	/** A confirmation field for the user's password - this has no known effect, but we set it anyway on password update. */
	confirmLocalPasswd?: string;
	/** The country in which the user lives/is based. */
	country?:         string;
	/** The user's email address. */
	email?:           string;
	/** The user's full name. */
	fullName?:        string;
	/** legacy field with no purpose. */
	gid?:             number;
	/** An integral, unique identifier for the user. */
	id:               number;
	/** The date/time at which the user was last updated. */
	lastUpdated?:     Date;
	/** The user's password - this should only be populated on update, and only if updating the password. */
	localPasswd?:     string;
	/** legacy field with no purpose. */
	localUser?:       boolean;
	/**
	 * Whether (false) or not (true) the user has reset their password after
	 * registration.
	 */
	newUser:          boolean;
	/** The user's phone number. */
	phoneNumber?:     string;
	/** The postal code where the user lives/is based. */
	postalCode?:      string;
	/** The user's public SSH key. */
	publicSshKey?:    string;
	/** The integral, unique identifier of the Role the user has. */
	role?:            number;
	/** The user's Role. */
	roleName?:        string;
	/** The state or province within which the user lives/is based. */
	stateOrProvince?: string;
	/** The Tenant to which the user belongs. */
	tenant?:          string;
	/** An integral, unique identifier for the Tenant to which the user belongs. */
	tenantId?:        number;
	/** legacy field with no purpose. */
	uid?:             number;
	/** The user's username. */
	username:         string;
}


/**
 * Represents a role that a user may have
 */
export interface Role {
	/**
	 * The Capabilities afforded by this Role.
	 */
	capabilities: Array<string>;
	/** A description of the Role. */
	description?: string;
	/** An integral, unique identifier for the Role. */
	id:           number;
	/** The Role's name. */
	name:         string;
	/** The Role's "privilege level". */
	privLevel:    number;
}

/**
 * Represents a user's ability to perform some action
 */
export interface Capability {
	/** The Capability's name. */
	name:         string;
	/** A description of the capability. */
	description:  string;
	/** The date/time at which the Capability was last updated. */
	lastUpdated?: Date;
}
