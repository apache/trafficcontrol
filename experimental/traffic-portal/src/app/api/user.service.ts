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
import type { GetResponseUser, PostRequestUser, PutOrPostResponseUser } from "trafficops-types";

import {
	type Role,
	type Tenant,
	type Capability,
	type CurrentUser,
	newCurrentUser
} from "src/app/models";

import { APIService } from "./base-api.service";

/**
 * UserService exposes API functionality related to Users, Roles and Capabilities.
 */
@Injectable()
export class UserService extends APIService {

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
	 * @returns The entire HTTP response on success, or `null` on failure.
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

	/**
	 * Gets a specific user from Traffic Ops.
	 *
	 * @param nameOrID The username (string) or ID (number) of the user to
	 * fetch.
	 * @returns An Array of User objects - or a single User object if 'nameOrID'
	 * was given.
	 */
	public async getUsers(nameOrID: string | number): Promise<GetResponseUser>;
	/**
	 * Gets an array of all users in Traffic Ops visible to the current user's
	 * Tenant.
	 *
	 * @param nameOrID If given, returns only the User with the given username
	 * (string) or ID (number).
	 * @returns An Array of User objects - or a single User object if 'nameOrID'
	 * was given.
	 */
	public async getUsers(): Promise<Array<GetResponseUser>>;
	/**
	 * Gets an array of users from Traffic Ops.
	 *
	 * @param nameOrID If given, returns only the User with the given username
	 * (string) or ID (number).
	 * @returns An Array of User objects - or a single User object if 'nameOrID'
	 * was given.
	 */
	public async getUsers(nameOrID?: string | number): Promise<Array<GetResponseUser> | GetResponseUser> {
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
			const r = await this.get<[GetResponseUser]>(path, undefined, params).toPromise();
			return {...r[0], lastUpdated: new Date((r[0].lastUpdated as unknown as string).replace("+00", "Z"))};
		}
		const users = await this.get<Array<GetResponseUser>>(path).toPromise();
		return users.map(
			u => ({...u, lastUpdated: new Date((u.lastUpdated as unknown as string).replace("+00", "Z"))})
		);
	}

	/**
	 * Replaces the current definition of a user with the one given.
	 *
	 * @param user The new definition of the User.
	 * @returns The user as updated.
	 */
	public async updateUser(user: PutOrPostResponseUser | GetResponseUser): Promise<PutOrPostResponseUser> {
		const path = `users/${user.id}`;
		const response = await this.put<PutOrPostResponseUser>(path, user).toPromise();
		if (response.registrationSent) {
			response.registrationSent = new Date((response.registrationSent as unknown as string));
		}
		return {
			...response,
			lastUpdated: new Date((response.lastUpdated as unknown as string).replace(" ", "T").replace("+00", "Z"))
		};
	}

	/**
	 * Creates a new user through the API.
	 *
	 * @param user The user to create.
	 * @returns The created user.
	 */
	public async createUser(user: PostRequestUser): Promise<PutOrPostResponseUser> {
		const response = await  this.post<PutOrPostResponseUser>("users", user).toPromise();
		if (response.registrationSent) {
			response.registrationSent = new Date((response.registrationSent as unknown as string));
		}
		return {
			...response,
			lastUpdated: new Date((response.lastUpdated as unknown as string).replace(" ", "T").replace("+00", "Z"))
		};
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

	/**
	 * Retrieves Tenants from Traffic Ops.
	 *
	 * @returns All Tenants visible to the requesting user's Tenant.
	 */
	public async getTenants(): Promise<Array<Tenant>>;
	/**
	 * Retrieves a Tenant from Traffic Ops.
	 *
	 * @param nameOrID Either the name or ID of the desired Tenant.
	 * @returns The Tenant identified by `nameOrID`.
	 */
	public async getTenants(nameOrID: string | number): Promise<Tenant>;
	/**
	 * Retrieves one or all Tenants from Traffic Ops.
	 *
	 * @param nameOrID Either the name or ID of a single desired Tenant.
	 * @returns The Tenant identified by `nameOrID` if given, otherwise all
	 * Tenants visible to the requesting user's Tenant.
	 */
	public async getTenants(nameOrID?: string | number): Promise<Array<Tenant> | Tenant> {
		const path = "tenants";
		if (nameOrID !== undefined) {
			let params;
			switch (typeof nameOrID) {
				case "string":
					params = {name: nameOrID};
					break;
				case "number":
					params = {id: String(nameOrID)};
			}
			const resp = await this.get<[Tenant]>(path, undefined, params).toPromise();
			return resp[0];
		}
		return this.get<Array<Tenant>>(path).toPromise();
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
