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
import type {
	PostRequestUser,
	PutRequestUser,
	PostResponseRole,
	PutResponseRole,
	RequestRole,
	RequestTenant,
	ResponseCurrentUser,
	ResponseRole,
	ResponseTenant,
	ResponseUser
} from "trafficops-types";

import { LoggingService } from "src/app/shared/logging.service";

/**
 * Represents a request to register a user via email using the `/users/register`
 * API endpoint.
 */
interface UserRegistrationRequest {
	email: string;
	role: number;
	tenantId: number;
}

/**
 * UserService exposes API functionality related to Users, Roles and Capabilities.
 */
@Injectable()
export class UserService {

	private lastID = 0;

	private testAdminUsername = "test-admin";
	private readonly testAdminPassword = "twelve12!";
	private readonly users: Array<ResponseUser> = [
		{
			addressLine1: null,
			addressLine2: null,
			changeLogCount: 0,
			city: null,
			company: null,
			country: null,
			email: "test@adm.in",
			fullName: "Test Admin",
			gid: null,
			id: ++this.lastID,
			lastAuthenticated: new Date(),
			lastUpdated: new Date(),
			newUser: false,
			phoneNumber: null,
			postalCode: null,
			publicSshKey: null,
			registrationSent: null,
			role: "admin",
			stateOrProvince: null,
			tenant: "root",
			tenantId: 1,
			ucdn: "",
			uid: null,
			username: "test-admin"
		}
	];
	private readonly roleDetail: Array<ResponseRole> = [{
		description: "Has access to everything - cannot be modified or deleted",
		lastUpdated: new Date(),
		name: "admin",
		permissions: [
			"ALL"
		],
	}];

	private readonly tenants: Array<ResponseTenant> = [
		{
			active: true,
			id: 1,
			lastUpdated: new Date(),
			name: "root",
			parentId: null
		},
		{
			active: true,
			id: 2,
			lastUpdated: new Date(),
			name: "test",
			parentId: 1
		}
	];

	private readonly tokens = new Map<string, string>();

	constructor(private readonly log: LoggingService) {}

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
		if (p !== undefined) {
			if (uOrT !== this.testAdminUsername || p !== this.testAdminPassword) {
				this.log.error("Invalid username or password.");
				return null;
			}
			return new HttpResponse({body: {alerts: [{level: "success", text: "Successfully logged in."}]}, status: 200});
		}
		const email = this.tokens.get(uOrT);
		if (email === undefined) {
			this.log.error(`token '${uOrT}' did not match any set token for any user`);
			return null;
		}
		const user = this.users.find(u=>u.email === email);
		if (!user) {
			this.log.error(`email '${email}' associated with token '${uOrT}' did not belong to any User`);
			return null;
		}
		this.tokens.delete(uOrT);
		return new HttpResponse({body: {alerts: [{level: "success", text: "Successfully logged in."}]}, status: 200});
	}

	/**
	 * Ends the current user's session - but does *not* affect the
	 * CurrentUserService's user data, which must be separately cleared.
	 *
	 * Note that in the testing environment this has no affect on the value of
	 * the "current user".
	 *
	 * @returns The entire HTTP response on success, or `null` on failure.
	 */
	public async logout(): Promise<HttpResponse<object> | null> {
		return new HttpResponse({body: {alerts: [{level: "success", text: "You are logged out."}]}});
	}

	/**
	 * Fetches the current user from Traffic Ops.
	 *
	 * @returns A `User` object representing the current user.
	 */
	public async getCurrentUser(): Promise<ResponseCurrentUser> {
		let user = this.users.filter(u=>u.username === this.testAdminUsername)[0];
		const transformUser = (u: ResponseUser): ResponseCurrentUser => ({
			...u,
			addressLine1: u.addressLine1,
			addressLine2: u.addressLine2,
			city: u.city,
			company: u.company,
			country: u.country,
			email: u.email,
			fullName: u.fullName,
			gid: u.gid,
			id: u.id,
			lastUpdated: u.lastUpdated,
			localUser: true,
			newUser: u.newUser ?? false,
			phoneNumber: u.phoneNumber,
			postalCode: u.postalCode,
			publicSshKey: u.publicSshKey,
			registrationSent: u.registrationSent ?? null,
			role: u.role,
			stateOrProvince: u.stateOrProvince,
			tenant: u.tenant,
			tenantId: u.tenantId,
			uid: u.uid,
			username: u.username
		});
		if (user) {
			return transformUser(user);
		}
		this.log.warn("stored admin username not found in stored users: from now on the current user will be (more or less) random");
		user = this.users[0];
		if (!user) {
			throw new Error("no users exist");
		}
		return transformUser(user);
	}

	/**
	 * Updates the current user to match the one passed in.
	 *
	 * @param user Unused. This method does nothing in the testing environment yet.
	 * @returns whether or not the request was successful.
	 */
	public async updateCurrentUser(user: ResponseCurrentUser): Promise<boolean> {
		const storedUser = this.users.findIndex(u=>u.id === user.id);
		if (storedUser < 0) {
			this.log.error(`no such User: #${user.id}`);
			return false;
		}
		this.testAdminUsername = user.username;
		this.users[storedUser] = {
			...user,
			lastUpdated: new Date(),
		};
		return true;
	}

	/**
	 * Gets a specific User.
	 *
	 * @param nameOrID The username (string) or ID (number) of the User to
	 * fetch.
	 * @returns The requested User.
	 */
	public async getUsers(nameOrID: string | number): Promise<ResponseUser>;
	/**
	 * Gets all stored Users.
	 *
	 * @returns All Users that are visible to the current user's Tenant.
	 */
	public async getUsers(): Promise<Array<ResponseUser>>;
	/**
	 * Gets one or all Users.
	 *
	 * @param nameOrID If given, returns only the User with the given username
	 * (string) or ID (number).
	 * @returns The requested User(s).
	 */
	public async getUsers(nameOrID?: string | number): Promise<Array<ResponseUser> | ResponseUser> {
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
		if (typeof(user) !== "number") {
			const idx = this.users.findIndex(u=>u.id === user.id);
			if (idx < 0) {
				throw new Error(`no such User: ${user.id}`);
			}
			const response = {
				...user,
				lastUpdated: new Date(),
			};
			this.users[idx] = response;
			return response;
		}
		if (!payload) {
			throw new Error("must supply a request body along with ID to update a User");
		}
		const index = this.users.findIndex(u => u.id === user);
		if (index < 0) {
			throw new Error(`no such User: ${user}`);
		}
		const tenant = await this.getTenants(payload.tenantId);
		const updated = {
			...payload,
			addressLine1: payload.addressLine1 ?? null,
			addressLine2: payload.addressLine2 ?? null,
			changeLogCount: 0,
			city: payload.city ?? null,
			company: payload.company ?? null,
			country: payload.country ?? null,
			gid: payload.gid ?? null,
			id: user,
			lastAuthenticated: null,
			lastUpdated: new Date(),
			newUser: payload.newUser ?? null,
			phoneNumber: payload.phoneNumber ?? null,
			postalCode: payload.postalCode ?? null,
			publicSshKey: payload.publicSshKey ?? null,
			registrationSent: null,
			stateOrProvince: payload.stateOrProvince ?? null,
			tenant: tenant.name,
			ucdn: payload.ucdn ?? "",
			uid: payload.uid ?? null,
		};
		this.users[index] = updated;
		return updated;
	}

	/**
	 * Creates a new user.
	 *
	 * @param user The user to create.
	 * @returns The created user.
	 */
	public async createUser(user: PostRequestUser): Promise<ResponseUser> {
		const role = this.roleDetail.find(r=>r.name === user.role);
		if (!role) {
			throw new Error(`no such Role: #${user.role}`);
		}
		const tenant = this.tenants.find(t=>t.id === user.tenantId);
		if (!tenant) {
			throw new Error(`no such Tenant: #${user.tenantId}`);
		}
		const response = {
			...user,
			addressLine1: user.addressLine1 ?? null,
			addressLine2: user.addressLine2 ?? null,
			changeLogCount: 0,
			city: user.city ?? null,
			company: user.company ?? null,
			confirmLocalPasswd: undefined,
			country: user.country ?? null,
			gid: user.gid ?? null,
			id: ++this.lastID,
			lastAuthenticated: null,
			lastUpdated: new Date(),
			newUser: user.newUser ?? null,
			phoneNumber: user.phoneNumber ?? null,
			postalCode: user.postalCode ?? null,
			publicSshKey: user.publicSshKey ?? null,
			registrationSent: null,
			stateOrProvince: user.stateOrProvince ?? null,
			tenant: tenant.name,
			tenantId: user.tenantId,
			ucdn: user.ucdn ?? "",
			uid: user.uid ?? null,
		};
		this.users.push(response);
		return response;
	}

	/**
	 * Registers a new user via email.
	 *
	 * Note that in testing this has no real effect.
	 *
	 * @param request The full registration request.
	 */
	public async registerUser(request: UserRegistrationRequest): Promise<void>;
	/**
	 * Registers a new user via email.
	 *
	 * Note that in testing this has no real effect.
	 *
	 * @param email The email address to use for registration.
	 * @param role The new user's Role (or just its ID).
	 * @param tenant The new user's Tenant (or just its ID).
	 */
	public async registerUser(email: string, role: number | ResponseRole, tenant: number | ResponseTenant): Promise<void>;
	/**
	 * Registers a new user via email.
	 *
	 * Note that in testing this has no real effect.
	 *
	 * @param userOrEmail Either the full registration request, or just the
	 * email address to use for registration.
	 * @param role The new user's Role (or just its ID). This is required if
	 * `userOrEmail` is given as an email address, and is ignored otherwise.
	 * @param tenant The new user's Tenant (or just its ID). This is required if
	 * `userOrEmail` is given as an email address, and is ignored otherwise.
	 */
	public async registerUser(
		userOrEmail: UserRegistrationRequest | string,
		role?: number | ResponseRole,
		tenant?: number | ResponseTenant
	): Promise<void> {
		if (typeof(userOrEmail) === "string") {
			if (role === undefined || tenant === undefined) {
				throw new Error("arguments 'role' and 'tenant' must be supplied when 'userOrEmail' is an email address");
			}
		}
	}

	/** Fetches the Role with the given name. */
	public async getRoles (name: string): Promise<ResponseRole>;
	/** Fetches all Roles. */
	public async getRoles (): Promise<Array<ResponseRole>>;
	/**
	 * Fetches one or all Roles from Traffic Ops.
	 *
	 * @param name unique identifier (name) of a single Role which will be fetched
	 * @throws {TypeError} When called with an improper argument.
	 * @returns Either an Array of Roles, or a single Role, depending on whether
	 * name was passed
	 */
	public async getRoles(name?: string): Promise<Array<ResponseRole> | ResponseRole> {
		if (name !== undefined) {
			const role = this.roleDetail.find(r=>r.name === name);
			if (!role) {
				throw new Error(`no such Role: ${name}`);
			}
			return role;
		}
		return this.roleDetail;
	}

	/**
	 * Creates a new role.
	 *
	 * @param role The role to create.
	 * @returns The created role along with lastUpdated field.
	 */
	public async createRole(role: RequestRole): Promise<PostResponseRole> {
		const resp = {
			lastUpdated: new Date(),
			...role
		};
		this.roleDetail.push(resp);
		return resp;
	}

	/**
	 * Updates an existing role.
	 *
	 * @param name The original role name
	 * @param role The role to update.
	 * @returns The updated role without lastUpdated field.
	 */
	public async updateRole(name: string, role: ResponseRole): Promise<PutResponseRole> {
		const roleName = this.roleDetail.findIndex(r => r.name === name);
		if (roleName < 0 ) {
			throw new Error(`no such Role: ${name}`);
		}
		this.roleDetail[roleName] = role;
		return role;
	}

	/**
	 * Deletes an existing role.
	 *
	 * @param role The role to be deleted.
	 * @returns The deleted role.
	 */
	public async deleteRole(role: string | ResponseRole): Promise<void> {
		const roleName = typeof(role) === "string" ? role : role.name;
		const index = this.roleDetail.findIndex(r => r.name === roleName);
		if (index === -1) {
			throw new Error(`no such role: ${role}`);
		}
		this.roleDetail.splice(index, 1);
	}

	/**
	 * Retrieves all (visible) Tenants from Traffic Ops.
	 *
	 * @returns All Tenants visible to the requesting user's Tenant.
	 */
	public async getTenants(): Promise<Array<ResponseTenant>>;
	/**
	 * Retrieves a specific Tenant from Traffic Ops.
	 *
	 * @param nameOrID Either the name or ID of a single desired Tenant.
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
		if (nameOrID !== undefined) {
			let tenant;
			switch (typeof nameOrID) {
				case "string":
					tenant = this.tenants.find(t=>t.name === nameOrID);
					break;
				case "number":
					tenant = this.tenants.find(t=>t.id === nameOrID);
			}
			if (!tenant) {
				throw new Error(`no such Tenant: ${nameOrID}`);
			}
			return tenant;
		}
		return this.tenants;
	}
	/**
	 * Creates a new tenant.
	 *
	 * @param tenant The Tenant to create.
	 * @returns The created tenant.
	 */
	public async createTenant(tenant: RequestTenant): Promise<ResponseTenant> {
		const resp = {
			...tenant,
			id: ++this.lastID,
			lastUpdated: new Date()
		};
		this.tenants.push(resp);
		return resp;
	}

	/**
	 * Updates an existing tenant.
	 *
	 * @param tenant The tenant to update.
	 * @returns The updated tenant.
	 */
	public async updateTenant(tenant: ResponseTenant): Promise<ResponseTenant> {
		const id = this.tenants.findIndex(t => t.id === tenant.id);
		if (id < 0) {
			throw new Error(`no such Tenant: ${tenant.id}`);
		}
		this.tenants[id] = tenant;
		return tenant;
	}

	/**
	 * Deletes an existing tenant.
	 *
	 * @param tenant The Tenant to be deleted, or just its ID.
	 * @returns The deleted Tenant.
	 */
	public async deleteTenant(tenant: number | ResponseTenant): Promise<ResponseTenant> {
		const id = typeof(tenant) === "number" ? tenant : tenant.id;
		const index = this.tenants.findIndex(t => t.id === id);
		if (index < 0) {
			throw new Error(`no such Tenant: ${id}`);
		}
		return this.tenants.splice(index, 1)[0];
	}

	/**
	 * Requests a password reset for a user.
	 *
	 * @param email The email of the user for whom to reset a password.
	 */
	public async resetPassword(email: string): Promise<void> {
		if (!this.users.some(u=>u.email === email)) {
			this.log.error(`no User exists with email '${email}' - TO doesn't expose that information with an error, so neither will we`);
			return;
		}
		const token = (Math.random() + 1).toString(36).substring(2);
		this.log.debug("setting token", token, "for email", email);
		this.tokens.set(token, email);
	}

}
