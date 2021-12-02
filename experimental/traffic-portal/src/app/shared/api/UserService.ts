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

import { HttpClient, HttpResponse } from "@angular/common/http";
import { Injectable } from "@angular/core";

import { Role, User, Capability, CurrentUser, newCurrentUser } from "../../models/user";

import { APIService } from "./APIService";

/**
 * UserService exposes API functionality related to Users, Roles and Capabilities.
 */
@Injectable()
export class UserService extends APIService {

	/**
	 * Injects the Angular HTTP client service into the parent constructor.
	 *
	 * @param http The Angular HTTP client service.
	 */
	constructor(http: HttpClient) {
		super(http);
	}

	/**
	 * Performs authentication with the Traffic Ops server.
	 *
	 * @param uOrT The username to be used for authentication, if `p` is
	 * provided. If `p` is **not** provided, then this is used as a token.
	 * @param p The password of user `u`
	 * @returns The entire HTTP response on success, or `null` on failure.
	 */
	public async login(uOrT: string, p?: string): Promise<HttpResponse<object> | null> {
		let path = `/api/${this.apiVersion}/user/login`;
		if (p !== undefined) {
			return this.http.post(path, {p, u: uOrT}, this.defaultOptions).toPromise().catch(
				e => {
					console.error("Failed to login:", e);
					return null;
				}
			);
		}
		path += "/token";
		return this.http.post(path, {t: uOrT}, this.defaultOptions).toPromise().catch(
			e => {
				console.error("Failed to login with token:", e);
				return null;
			}
		);
	}

	/**
	 * Ends the current user's session - but does *not* affect the
	 * CurrentUserService's user data, which must be separately cleared.
	 *
	 * @returns The entire HTTP response on succes, or `null` on failure.
	 */
	public async logout(): Promise<HttpResponse<object> | null> {
		const path = `/api/${this.apiVersion}/user/logout`;
		return this.http.post(path, undefined, this.defaultOptions).toPromise().catch(
			e => {
				console.error("Failed to logout:", e);
				return null;
			}
		);
	}

	/**
	 * Fetches the current user from Traffic Ops.
	 *
	 * @returns A `User` object representing the current user.
	 */
	public async getCurrentUser(): Promise<CurrentUser> {
		const path = "user/current";
		return this.get<CurrentUser>(path).toPromise().then(
			r => {
				r.lastUpdated = new Date((r.lastUpdated as unknown as string).replace("+00", "Z"));
				return r;
			}
		).catch(
			e => {
				console.error("Failed to get current user:", e);
				return newCurrentUser();
			}
		);
	}

	/**
	 * Updates the current user to match the one passed in.
	 *
	 * @param user The new form of the user.
	 * @returns whether or not the request was successful.
	 */
	public async updateCurrentUser(user: CurrentUser): Promise<boolean> {
		const path = "user/current";
		return this.put<CurrentUser>(path, {user}).toPromise().then(
			() => true,
			() => false
		);
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
		const path = "users";
		if (nameOrID) {
			let params;
			switch (typeof nameOrID) {
				case "string":
					params = {username: nameOrID};
					break;
				case "number":
					params = {id: String(nameOrID)};
			}
			return this.get<[User]>(path, undefined, params).toPromise().then(
				r => {
					r[0].lastUpdated = new Date((r[0].lastUpdated as unknown as string).replace("+00", "Z"));
					return r[0];
				}
			).catch(
				e => {
					console.error("Failed to get user:", e);
					return {
						id: -1,
						newUser: false,
						username: ""
					};
				}
			);
		}
		return this.get<Array<User>>(path).toPromise().then(r => r.map(
			u => {
				u.lastUpdated = new Date((u.lastUpdated as unknown as string).replace("+00", "Z"));
				return u;
			}
		)).catch(
			e => {
				console.error("Failed to get users:", e);
				return [];
			}
		);
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
		const path = "roles";
		if (nameOrID !== undefined) {
			let params;
			switch (typeof nameOrID) {
				case "string":
					params = {name: nameOrID};
					break;
				case "number":
					params = {id: String(nameOrID)};
			}
			return this.get<[Role]>(path, undefined, params).toPromise().then(r => r[0]).catch(
				e => {
					console.error("Failed to get Role:", e);
					return {
						capabilities: [],
						id: -1,
						name: "",
						privLevel: -1,
					};
				}
			);
		}
		return this.get<Array<Role>>(path).toPromise().catch(
			e => {
				console.error("Failed to get Roles:", e);
				return [];
			}
		);
	}

	/** Fetches the User Capability (Permission) with the given name. */
	public async getCapabilities (name: string): Promise<Capability>;
	/** Fetches all User Capabilities (Permissions). */
	public async getCapabilities (): Promise<Array<Capability>>;
	/**
	 * Fetches one or all Capabilities from Traffic Ops.
	 *
	 * @param name Optionally, the name of a single Capability which will be fetched
	 * @throws {TypeError} When called with an improper argument.
	 * @returns Either an Array of Capabilities, or a single Capability,
	 * depending on whether `name`/`id` was passed
	 */
	public async getCapabilities(name?: string): Promise<Array<Capability> | Capability> {
		const path = "capabilities";
		if (name) {
			return this.get<[Capability]>(path, undefined, {name}).toPromise().then(
				r => r[0]
			).catch(
				e => {
					console.error("Failed to get user Permission:", e);
					return {
						description: "",
						name: ""
					};
				}
			);
		}
		return this.get<Array<Capability>>(path).toPromise().catch(
			e => {
				console.error("Failed to get user Permissions:", e);
				return [];
			}
		);
	}

	/**
	 * Requests a password reset for a user.
	 *
	 * @param email The email of the user for whom to reset a password.
	 */
	public async resetPassword(email: string): Promise<void> {
		await this.post("user/reset_password", {email}).toPromise();
	}

}
