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

/** The name of a special Role that is always allowed to do whatever it wants. */
export const ADMIN_ROLE = "admin";
