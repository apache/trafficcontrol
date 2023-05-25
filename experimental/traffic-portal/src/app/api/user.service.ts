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

import { HttpClient, type HttpResponse } from "@angular/common/http";
import { Injectable } from "@angular/core";
import {
	type ResponseUser,
	type PostRequestUser,
	type RequestRole,
	type RequestTenant,
	type ResponseCurrentUser,
	type ResponseRole,
	type ResponseTenant,
	type PutRequestUser,
	type RegistrationRequest,
	userEmailIsValid
} from "trafficops-types";

import { APIService } from "./base-api.service";

/**
 * UserService exposes API functionality related to Users, Roles and Tenants.
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
			return this.http.post(path, {p, u: uOrT}, this.defaultOptions).toPromise();
		}
		path += "/token";
		return this.http.post(path, {t: uOrT}, this.defaultOptions).toPromise();
	}

	/**
	 * Ends the current user's session - but does *not* affect the
	 * CurrentUserService's user data, which must be separately cleared.
	 *
	 * @returns The entire HTTP response on success, or `null` on failure.
	 */
	public async logout(): Promise<HttpResponse<object> | null> {
		const path = `/api/${this.apiVersion}/user/logout`;
		return this.http.post(path, undefined, this.defaultOptions).toPromise();
	}

	/**
	 * Fetches the current user from Traffic Ops.
	 *
	 * @returns A `User` object representing the current user.
	 */
	public async getCurrentUser(): Promise<ResponseCurrentUser> {
		const path = "user/current";
		return this.get<ResponseCurrentUser>(path).toPromise();
	}

	/**
	 * Updates the current user to match the one passed in.
	 *
	 * @param user The new form of the user.
	 * @returns whether or not the request was successful.
	 */
	public async updateCurrentUser(user: ResponseCurrentUser): Promise<boolean> {
		const path = "user/current";
		return this.put<ResponseCurrentUser>(path, user).toPromise().then(
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
	public async getUsers(nameOrID: string | number): Promise<ResponseUser>;
	/**
	 * Gets an array of all users in Traffic Ops visible to the current user's
	 * Tenant.
	 *
	 * @param nameOrID If given, returns only the User with the given username
	 * (string) or ID (number).
	 * @returns An Array of User objects - or a single User object if 'nameOrID'
	 * was given.
	 */
	public async getUsers(): Promise<Array<ResponseUser>>;
	/**
	 * Gets an array of users from Traffic Ops.
	 *
	 * @param nameOrID If given, returns only the User with the given username
	 * (string) or ID (number).
	 * @returns An Array of User objects - or a single User object if 'nameOrID'
	 * was given.
	 */
	public async getUsers(nameOrID?: string | number): Promise<Array<ResponseUser> | ResponseUser> {
		const path = "users";
		if (nameOrID) {
			let params;
			switch (typeof nameOrID) {
				case "string":
					params = {username: nameOrID};
					break;
				case "number":
					params = {id: nameOrID};
			}
			const r = await this.get<[ResponseUser]>(path, undefined, params).toPromise();
			return r[0];
		}
		return this.get<Array<ResponseUser>>(path).toPromise();
	}

	/**
	 * Replaces the current definition of a user with the one given.
	 *
	 * @param user The full new definition of the User.
	 * @returns The user as updated.
	 */
	public async updateUser(user: ResponseUser): Promise<ResponseUser>;
	/**
	 * Replaces the current definition of a user with the one given.
	 *
	 * @param user The ID of the User being updated.
	 * @param payload The new definition of the User.
	 * @returns The user as updated.
	 */
	public async updateUser(user: number, payload: PutRequestUser): Promise<ResponseUser>;
	/**
	 * Replaces the current definition of a user with the one given.
	 *
	 * @param user The new definition of the User, or just its ID.
	 * @param payload The new definition of the User. This is required if `user`
	 * is an ID, and ignored otherwise.
	 * @returns The user as updated.
	 */
	public async updateUser(user: ResponseUser | number, payload?: PutRequestUser): Promise<ResponseUser> {
		let id;
		let body;
		if (typeof(user) === "number") {
			id = user;
			body = payload;
			if (!body) {
				throw new Error("must supply a request body along with ID to update a User");
			}
		} else {
			body = user;
			id = user.id;
		}
		const path = `users/${id}`;
		return this.put<ResponseUser>(path, body).toPromise();
	}

	/**
	 * Creates a new user through the API.
	 *
	 * @param user The user to create.
	 * @returns The created user.
	 */
	public async createUser(user: PostRequestUser): Promise<ResponseUser> {
		return this.post<ResponseUser>("users", user).toPromise();
	}

	/**
	 * Registers a new user via email.
	 *
	 * @param request The full registration request.
	 */
	public async registerUser(request: RegistrationRequest): Promise<void>;
	/**
	 * Registers a new user via email.
	 *
	 * @param email The email address to use for registration.
	 * @param role The new user's Role (or just its ID).
	 * @param tenant The new user's Tenant (or just its ID).
	 */
	public async registerUser(email: string, role: string | ResponseRole, tenant: number | ResponseTenant): Promise<void>;
	/**
	 * Registers a new user via email.
	 *
	 * @param userOrEmail Either the full registration request, or just the
	 * email address to use for registration.
	 * @param role The new user's Role (or just its ID). This is required if
	 * `userOrEmail` is given as an email address, and is ignored otherwise.
	 * @param tenant The new user's Tenant (or just its ID). This is required if
	 * `userOrEmail` is given as an email address, and is ignored otherwise.
	 */
	public async registerUser(
		userOrEmail: RegistrationRequest | string,
		role?: string | ResponseRole,
		tenant?: number | ResponseTenant
	): Promise<void> {
		let request: RegistrationRequest;
		if (typeof(userOrEmail) === "string") {
			if (!userEmailIsValid(userOrEmail)) {
				throw new Error(`invalid email address: '${userOrEmail}'`);
			}
			if (role === undefined || tenant === undefined) {
				throw new Error("arguments 'role' and 'tenant' must be supplied when 'userOrEmail' is an email address");
			}
			request = {
				email: userOrEmail,
				role: typeof(role) === "string" ? role : role.name,
				tenantId: typeof(tenant) === "number" ? tenant : tenant.id
			};
		} else {
			request = userOrEmail;
		}

		await this.post("users/register", request).toPromise();
	}

	/**
	 * Fetches one or all Roles from Traffic Ops.
	 *
	 * @param nameOrID The name or integral, unique identifier of the single
	 * Role which will be fetched.
	 * @returns The requested Role.
	 */
	public async getRoles (nameOrID: string): Promise<ResponseRole>;
	/**
	 * Fetches all Roles from Traffic Ops.
	 *
	 * @returns An Array of Roles.
	 */
	public async getRoles (): Promise<Array<ResponseRole>>;
	/**
	 * Fetches one or all Roles from Traffic Ops.
	 *
	 * @param name Optionally, the name of a single Role which will be fetched.
	 * @throws {TypeError} When called with an improper argument.
	 * @returns Either an Array of Roles, or a single Role, depending on whether
	 * `name`/`id` was passed.
	 */
	public async getRoles(name?: string): Promise<Array<ResponseRole> | ResponseRole> {
		const path = "roles";
		if (name !== undefined) {
			const resp = await this.get<[ResponseRole]>(path, undefined, {name}).toPromise();
			if (resp.length !== 1) {
				throw new Error(`Traffic Ops responded with ${resp.length} Roles by identifier ${name}`);
			}
			return resp[0];
		}
		return this.get<Array<ResponseRole>>(path).toPromise();
	}

	/**
	 * Creates a new Role.
	 *
	 * @param role The role to create.
	 * @returns The created role along with lastUpdated field.
	 */
	public async createRole(role: RequestRole): Promise<ResponseRole> {
		return this.post<ResponseRole>("roles", role).toPromise();
	}

	/**
	 * Updates an existing Role.
	 *
	 * @param name The original role name
	 * @param role The role to update.
	 * @returns The updated role without lastUpdated field.
	 */
	public async updateRole(name: string, role: ResponseRole): Promise<ResponseRole> {
		return this.put<ResponseRole>("roles?", role, {name}).toPromise();
	}

	/**
	 * Deletes an existing role.
	 *
	 * @param role The role to be deleted.
	 * @returns The deleted role.
	 */
	public async deleteRole(role: string | ResponseRole): Promise<void> {
		const name = typeof(role) === "string" ? role : role.name;
		return this.delete("roles", undefined, {name}).toPromise();
	}

	/**
	 * Retrieves Tenants from Traffic Ops.
	 *
	 * @returns All Tenants visible to the requesting user's Tenant.
	 */
	public async getTenants(): Promise<Array<ResponseTenant>>;
	/**
	 * Retrieves a Tenant from Traffic Ops.
	 *
	 * @param nameOrID Either the name or ID of the desired Tenant.
	 * @returns The Tenant identified by `nameOrID`.
	 */
	public async getTenants(nameOrID: string | number): Promise<ResponseTenant>;
	/**
	 * Retrieves one or all Tenants from Traffic Ops.
	 *
	 * @param nameOrID Either the name or ID of a single desired Tenant.
	 * @returns The Tenant identified by `nameOrID` if given, otherwise all
	 * Tenants visible to the requesting user's Tenant.
	 */
	public async getTenants(nameOrID?: string | number): Promise<Array<ResponseTenant> | ResponseTenant> {
		const path = "tenants";
		if (nameOrID !== undefined) {
			let params;
			switch (typeof nameOrID) {
				case "string":
					params = {name: nameOrID};
					break;
				case "number":
					params = {id: nameOrID};
			}
			const resp = await this.get<[ResponseTenant]>(path, undefined, params).toPromise();
			return resp[0];
		}
		return this.get<Array<ResponseTenant>>(path).toPromise();
	}

	/**
	 * Creates a new tenant.
	 *
	 * @param tenant The Tenant to create.
	 * @returns The created tenant.
	 */
	public async createTenant(tenant: RequestTenant): Promise<ResponseTenant> {
		return this.post<ResponseTenant>("tenants", tenant).toPromise();
	}

	/**
	 * Updates an existing tenant.
	 *
	 * @param tenant The tenant to update.
	 * @returns The updated tenant.
	 */
	public async updateTenant(tenant: ResponseTenant): Promise<ResponseTenant> {
		return this.put<ResponseTenant>(`tenants/${tenant.id}`, tenant).toPromise();
	}

	/**
	 * Deletes an existing tenant.
	 *
	 * @param tenant The Tenant to be deleted, or just its ID.
	 * @returns The deleted Tenant.
	 */
	public async deleteTenant(tenant: number | ResponseTenant): Promise<ResponseTenant> {
		const id = typeof(tenant) === "number" ? tenant : tenant.id;
		return this.delete<ResponseTenant>(`tenants/${id}`).toPromise();
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
