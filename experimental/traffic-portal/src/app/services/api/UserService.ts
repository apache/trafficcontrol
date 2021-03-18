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

import { Observable } from "rxjs";
import { map } from "rxjs/operators";

import { Role, User, Capability } from "../../models/user";

import { APIService } from "./apiservice";

/**
 * UserService exposes API functionality related to Users, Roles and Capabilities.
 */
@Injectable({providedIn: "root"})
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
	 * @param u The username to be used for authentication
	 * @param p The password of user `u`
	 * @returns An observable that will emit the entire HTTP response
	 */
	public login(u: string, p: string): Observable<HttpResponse<object>> {
		const path = `/api/${this.apiVersion}/user/login`;
		return this.http.post(path, {p, u}, this.defaultOptions);
	}

	/**
	 * Fetches the current user from Traffic Ops.
	 *
	 * @returns An observable that will emit a `User` object representing the current user.
	 */
	public getCurrentUser(): Observable<User> {
		const path = "user/current";
		return this.get<User>(path).pipe(map(
			r => {
				r.lastUpdated = new Date((r.lastUpdated as unknown as string).replace("+00", "Z"));
				return r;
			}
		));
	}

	public getUsers(nameOrID: string | number): Observable<User>;
	public getUsers(): Observable<Array<User>>;
	/**
	 * Gets an array of all users in Traffic Ops.
	 *
	 * @param nameOrID If given, returns only the User with the given username (string) or ID (number).
	 * @returns An Observable that will emit an Array of User objects - or a single User object if 'nameOrID' was given.
	 */
	public getUsers(nameOrID?: string | number): Observable<Array<User> | User> {
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
			return this.get<[User]>(path, undefined, params).pipe(map(
				r => {
					r[0].lastUpdated = new Date((r[0].lastUpdated as unknown as string).replace("+00", "Z"));
					return r[0];
				}
			));
		}
		return this.get<Array<User>>(path).pipe(map(r => r.map(
			u => {
				u.lastUpdated = new Date((u.lastUpdated as unknown as string).replace("+00", "Z"));
				return u;
			}
		)));
	}

	/** Fetches the Role with the given ID. */
	public getRoles (nameOrID: number | string): Observable<Role>;
	/** Fetches all Roles. */
	public getRoles (): Observable<Array<Role>>;
	/**
	 * Fetches one or all Roles from Traffic Ops.
	 *
	 * @param nameOrID Optionally, the name or integral, unique identifier of a single Role which will be fetched
	 * @throws {TypeError} When called with an improper argument.
	 * @returns an Observable that will emit either an Array of Roles, or a single Role, depending on whether
	 * `name`/`id` was passed
	 * (In the event that `name`/`id` is given but does not match any Role, `null` will be emitted)
	 */
	public getRoles(nameOrID?: string | number): Observable<Array<Role> | Role> {
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
			return this.get<[Role]>(path, undefined, params).pipe(map(r => r[0]));
		}
		return this.get<Array<Role>>(path);
	}

	/** Fetches the User Capability (Permission) with the given name. */
	public getCapabilities (name: string): Observable<Capability>;
	/** Fetches all User Capabilities (Permissions). */
	public getCapabilities (): Observable<Array<Capability>>;
	/**
	 * Fetches one or all Capabilities from Traffic Ops.
	 *
	 * @param name Optionally, the name of a single Capability which will be fetched
	 * @throws {TypeError} When called with an improper argument.
	 * @returns an Observable that will emit either an Array of Capabilities, or a single Capability,
	 * depending on whether `name`/`id` was passed
	 */
	public getCapabilities(name?: string): Observable<Array<Capability> | Capability> {
		const path = "capabilities";
		if (name) {
			return this.get<[Capability]>(path, undefined, {name}).pipe(map(
				r => r[0]
			));
		}
		return this.get<Array<Capability>>(path);
	}

}
