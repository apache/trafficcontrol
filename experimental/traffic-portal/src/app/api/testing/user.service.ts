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

import { HttpResponse } from "@angular/common/http";
import { Injectable } from "@angular/core";

import type { Role, User, Capability, CurrentUser } from "../../models/user";

/**
 * UserService exposes API functionality related to Users, Roles and Capabilities.
 */
@Injectable()
export class UserService {

	private testAdminUsername = "test-admin";
	private readonly testAdminPassword = "twelve12!";
	private readonly users: Array<CurrentUser> = [
		{
			addressLine1: null,
			addressLine2: null,
			city: null,
			company: null,
			country: null,
			email: "test@adm.in",
			fullName: "Test Admin",
			gid: null,
			id: 1,
			lastUpdated: new Date(),
			localUser: true,
			newUser: false,
			phoneNumber: null,
			postalCode: null,
			publicSshKey: null,
			role: 1,
			roleName: "admin",
			stateOrProvince: null,
			tenant: "root",
			tenantId: 1,
			uid: null,
			username: "test-admin"
		}
	];
	private readonly roles = [
		{
			capabilities: [
				"ALL",
				"PARAMETER-SECURE:READ"
			],
			description: "Has access to everything - cannot be modified or deleted",
			id: 1,
			lastUpdated: new Date(),
			name: "admin",
			privLevel: 30
		}
	];
	private readonly capabilities = [
		{
			description: "unknown - comes from a Permission",
			lastUpdated: new Date(),
			name: "ALL"
		},
		{
			description: "unknown - comes from a Permission",
			lastUpdated: new Date(),
			name: "PARAMETER-SECURE:READ"
		}
	];

	private readonly tokens = new Map<string, string>();

	/**
	 * Performs authentication with the Traffic Ops server.
	 *
	 * Note that in the testing environment, this gives back more information
	 * than the concrete service in the event of a token authentication error.
	 *
	 * Also note that the value of the "current user" is unaffected by any calls
	 * to login (in the testing environment).
	 *
	 * @param uOrT The username to be used for authentication, if `p` is
	 * provided. If `p` is **not** provided, then this is used as a token.
	 * @param p The password of user `u`
	 * @returns The entire HTTP response on success, or `null` on failure.
	 */
	public async login(uOrT: string, p?: string): Promise<HttpResponse<object> | null> {
		if (p !== undefined && (uOrT !== "test-admin" || p !== this.testAdminPassword)) {
			throw new Error("Invalid username or password.");
		}
		const email = this.tokens.get(uOrT);
		if (email === undefined) {
			throw new Error(`token '${uOrT}' did not match any set token for any user`);
		}
		const user = this.users.find(u=>u.email === email);
		if (!user) {
			throw new Error(`email '${email}' associated with token '${uOrT}' did not belong to any User`);
		}
		this.tokens.delete(uOrT);
		return new HttpResponse({body: {alerts: [{level: "success", text: "Successfully logged in."}]}});
	}

	/**
	 * Ends the current user's session - but does *not* affect the
	 * CurrentUserService's user data, which must be separately cleared.
	 *
	 * Note that in the testing environment this has no affect on the value of
	 * the "current user".
	 *
	 * @returns The entire HTTP response on succes, or `null` on failure.
	 */
	public async logout(): Promise<HttpResponse<object> | null> {
		return new HttpResponse({body: {alerts: [{level: "success", text: "You are logged out."}]}});
	}

	/**
	 * Fetches the current user from Traffic Ops.
	 *
	 * @returns A `User` object representing the current user.
	 */
	public async getCurrentUser(): Promise<CurrentUser> {
		let user = this.users.filter(u=>u.username === this.testAdminUsername)[0];
		if (user) {
			return user;
		}
		console.warn("stored admin username not found in stored users: from now on the current user will be (more or less) random");
		user = this.users[0];
		if (!user) {
			throw new Error("no users exist");
		}
		return user;
	}

	/**
	 * Updates the current user to match the one passed in.
	 *
	 * @param user Unused. This method does nothing in the testing environment yet.
	 * @returns whether or not the request was successful.
	 */
	public async updateCurrentUser(user: CurrentUser): Promise<boolean> {
		const storedUser = this.users.findIndex(u=>u.id === user.id);
		if (storedUser < 0) {
			console.error(`no such User: #${user.id}`);
			return false;
		}
		this.testAdminUsername = user.username;
		this.users[storedUser] = user;
		this.users[storedUser].lastUpdated = new Date();
		return true;
	}

	public async getUsers(nameOrID: string | number): Promise<User>;
	public async getUsers(): Promise<Array<User>>;
	/**
	 * Gets an array of all users in Traffic Ops.
	 *
	 * @param nameOrID If given, returns only the User with the given username (string) or ID (number).
	 * @returns An Array of User objects - or a single User object if 'nameOrID' was given.
	 */
	public async getUsers(nameOrID?: string | number): Promise<Array<User> | User> {
		if (nameOrID) {
			let user;
			switch (typeof nameOrID) {
				case "string":
					user = this.users.filter(u=>u.username === nameOrID)[0];
					break;
				case "number":
					user = this.users.filter(u=>u.id === nameOrID)[0];
			}
			if (!user) {
				throw new Error(`no such User: ${nameOrID}`);
			}
			return user;
		}
		return this.users;
	}

	/** Fetches the Role with the given ID. */
	public async getRoles (nameOrID: number | string): Promise<Role>;
	/** Fetches all Roles. */
	public async getRoles (): Promise<Array<Role>>;
	/**
	 * Fetches one or all Roles from Traffic Ops.
	 *
	 * @param nameOrID Optionally, the name or integral, unique identifier of a single Role which will be fetched
	 * @throws {TypeError} When called with an improper argument.
	 * @returns Either an Array of Roles, or a single Role, depending on whether
	 * `name`/`id` was passed
	 */
	public async getRoles(nameOrID?: string | number): Promise<Array<Role> | Role> {
		if (nameOrID !== undefined) {
			let role;
			switch (typeof nameOrID) {
				case "string":
					role = this.roles.find(r=>r.name === nameOrID);
					break;
				case "number":
					role = this.roles.find(r=>r.id === nameOrID);
			}
			if (!role) {
				throw new Error(`no such Role: ${nameOrID}`);
			}
			return role;
		}
		return this.roles;
	}

	/** Fetches the User Capability (Permission) with the given name. */
	public async getCapabilities (name: string): Promise<Capability>;
	/** Fetches all User Capabilities (Permissions). */
	public async getCapabilities (): Promise<Array<Capability>>;
	/**
	 * Fetches one or all Capabilities from Traffic Ops.
	 *
	 * @deprecated "Capabilities" are deprecated in favor of Permissions.
	 * "Capabilities" are removed from API v4 and later.
	 *
	 * @param name Optionally, the name of a single Capability which will be fetched
	 * @throws {TypeError} When called with an improper argument.
	 * @returns Either an Array of Capabilities, or a single Capability,
	 * depending on whether `name`/`id` was passed
	 */
	public async getCapabilities(name?: string): Promise<Array<Capability> | Capability> {
		if (name) {
			const cap = this.capabilities.find(c=>c.name === name);
			if (!cap) {
				throw new Error(`no such Capability: ${name}`);
			}
			return cap;
		}
		return this.capabilities;
	}

	/**
	 * Requests a password reset for a user.
	 *
	 * @param email The email of the user for whom to reset a password.
	 */
	public async resetPassword(email: string): Promise<void> {
		if (!this.users.some(u=>u.email === email)) {
			console.error(`no User exists with email '${email}' - TO doesn't expose that information with an error, so neither will we`);
			return;
		}
		const token = (Math.random() + 1).toString(36).substring(2);
		console.log("setting token", token, "for email", email);
		this.tokens.set(token, email);
	}

}
