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
 * Represents a Traffic Ops user (`tm_user` in the database)
 */
export interface User {
	/** Line one of the user's address. */
	addressLine1?:    string | null;
	/** Line two of the user's address. */
	addressLine2?:    string | null;
	/** The city in which the user lives/is based. */
	city?:            string | null;
	/** The company for which the user works. */
	company?:         string | null;
	/** A confirmation field for the user's password - this has no known effect, but we set it anyway on password update. */
	confirmLocalPasswd?: string | null;
	/** The country in which the user lives/is based. */
	country?:         string | null;
	/** The user's email address. */
	email?:           string | null;
	/** The user's full name. */
	fullName?:        string | null;
	/** legacy field with no purpose. */
	gid?:             number | null;
	/** An integral, unique identifier for the user. */
	id:               number;
	/** The date/time at which the user was last updated. */
	lastUpdated?:     Date | null;
	/** The user's password - this should only be populated on update, and only if updating the password. */
	localPasswd?:     string | null;
	/**
	 * Whether (false) or not (true) the user has reset their password after
	 * registration.
	 */
	newUser:          boolean;
	/** The user's phone number. */
	phoneNumber?:     string | null;
	/** The postal code where the user lives/is based. */
	postalCode?:      string | null;
	/** The user's public SSH key. */
	publicSshKey?:    string | null;
	/** The integral, unique identifier of the Role the user has. */
	role?:            number;
	/** The user's Role. */
	rolename?:        string | null;
	/** The state or province within which the user lives/is based. */
	stateOrProvince?: string | null;
	/** The Tenant to which the user belongs. */
	tenant?:          string | null;
	/** An integral, unique identifier for the Tenant to which the user belongs. */
	tenantId?:        number;
	/** legacy field with no purpose. */
	uid?:             number | null;
	/** The user's username. */
	username:         string;
}

/**
 * CurrentUser represents a "current user" (mostly as seen in API *responses* -
 * request structures are subtly different in a few ways). This differs from a
 * `User` in a few key ways, most notably `rolename` vs `roleName`.
 */
export interface CurrentUser {
	addressLine1: string | null;
	addressLine2: string | null;
	city: string | null;
	confirmLocalPasswd?: string | null;
	company: string | null;
	country: string | null;
	email: string;
	fullName: string | null;
	gid: number | null;
	id: number;
	lastUpdated: Date;
	localPasswd?: string | null;
	localUser: boolean;
	newUser: boolean;
	phoneNumber: string | null;
	postalCode: string | null;
	publicSshKey: string | null;
	role: number;
	roleName: string;
	stateOrProvince: string | null;
	tenant: string;
	tenantId: number;
	uid: number | null;
	username: string;
}

/**
 * Gets a new `CurrentUser` to use as a default structure.
 *
 * @returns A valid `CurrentUser` - but one that will absolutely fail validation
 * server-side for several reasons. Should not be used directly.
 */
export function newCurrentUser(): CurrentUser {
	return {
		addressLine1: null,
		addressLine2: null,
		city: null,
		company: null,
		country: null,
		email: "",
		fullName: "",
		gid: null,
		id: -1,
		lastUpdated: new Date(),
		localUser: true,
		newUser: false,
		phoneNumber: null,
		postalCode: null,
		publicSshKey: null,
		role: -1,
		roleName: "",
		stateOrProvince: null,
		tenant: "",
		tenantId: -1,
		uid: null,
		username: "",
	};
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
 * Represents a group of Users that can own certain resources.
 */
export interface Tenant {
	active: boolean;
	readonly id: number;
	readonly lastUpdated: Date;
	name: string;
	parentId: number | null;
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

/** The name of a special Role that is always allowed to do whatever it wants. */
export const ADMIN_ROLE = "admin";
