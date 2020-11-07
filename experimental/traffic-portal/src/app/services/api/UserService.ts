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

import { APIService } from "./apiservice";

import { Role, User, Capability } from "../../models/user";

/**
 * UserService exposes API functionality related to Users, Roles and Capabilities.
 */
@Injectable({providedIn: "root"})
export class UserService extends APIService {

	/**
	 * Performs authentication with the Traffic Ops server.
	 * @param u The username to be used for authentication
	 * @param p The password of user `u`
	 * @returns An observable that will emit the entire HTTP response
	 */
	public login (u: string, p: string): Observable<HttpResponse<any>> {
		const path = `/api/${this.API_VERSION}/user/login`;
		return this.post(path, {u, p});
	}

	/**
	 * Fetches the current user from Traffic Ops.
	 * @returns An observable that will emit a `User` object representing the current user.
	 */
	public getCurrentUser (): Observable<User> {
		const path = `/api/${this.API_VERSION}/user/current`;
		return this.get(path).pipe(map(
			r => {
				return r.body.response as User;
			}
		));
	}

	/**
	 * Gets an array of all users in Traffic Ops
	 * @returns An Observable that will emit an Array of User objects.
	 */
	public getUsers (nameOrID?: string | number): Observable<Array<User> | User | null> {
		const path = `/api/${this.API_VERSION}/users`;
		if (nameOrID) {
			switch (typeof nameOrID) {
				case "string":
					return this.get(`${path}?username=${encodeURIComponent(nameOrID)}`).pipe(map(
						r => {
							for (const u of (r.body.response as Array<User>)) {
								if (u.username === nameOrID) {
									return u;
								}
							}
							return null;
						}
					));
				case "number":
					return this.get(`${path}?id=${nameOrID}`).pipe(map(
						r => {
							for (const u of (r.body.response as Array<User>)) {
								if (u.id === nameOrID) {
									return u;
								}
							}
							return null;
						}
					));
				default:
					throw new TypeError(`expected a username or ID, got '${typeof (nameOrID)}'`);
			}
		}
		return this.get(path).pipe(map(
			r => {
				return r.body.response as Array<User>;
			}
		));
	}

	/** Fetches the Role with the given ID */
	public getRoles (nameOrID: number | string): Observable<Role | null>;
	/** Fetches all Roles */
	public getRoles (): Observable<Array<Role>>;
	/**
	 * Fetches one or all Roles from Traffic Ops
	 * @param nameOrID Optionally, the name or integral, unique identifier of a single Role which will be fetched
	 * @throws {TypeError} When called with an improper argument.
	 * @returns an Observable that will emit either an Array of Roles, or a single Role, depending on whether
	 * `name`/`id` was passed
	 * (In the event that `name`/`id` is given but does not match any Role, `null` will be emitted)
	 */
	public getRoles(nameOrID?: string | number): Observable<Array<Role> | Role | null> {
		const path = `/api/${this.API_VERSION}/roles`;
		if (nameOrID) {
			switch (typeof nameOrID) {
				case "string":
					return this.get(`${path}?name=${nameOrID}`).pipe(map(
						r => {
							for (const role of (r.body.response as Array<Role>)) {
								if (role.name === nameOrID) {
									return role;
								}
							}
							return null;
						}
					));
					break;
				case "number":
					return this.get(`${path}?id=${nameOrID}`).pipe(map(
						r => {
							for (const role of (r.body.response as Array<Role>)) {
								if (role.id === nameOrID) {
									return role;
								}
							}
							return null;
						}
					));
					break;
				default:
					throw new TypeError(`expected a name or ID, got '${typeof (nameOrID)}'`);
					break;
			}
		}
		return this.get(path).pipe(map(
			r => {
				return r.body.response as Array<Role>;
			}
		));
	}

	/** Fetches the User Capability (Permission) with the given name. */
	public getCapabilities (name: string): Observable<Capability | null>;
	/** Fetches all User Capabilities (Permissions). */
	public getCapabilities (): Observable<Array<Capability>>;
	/**
	 * Fetches one or all Capabilities from Traffic Ops
	 * @param name Optionally, the name of a single Capability which will be fetched
	 * @throws {TypeError} When called with an improper argument.
	 * @returns an Observable that will emit either an Array of Capabilities, or a single Capability,
	 * depending on whether `name`/`id` was passed
	 * (In the event that `name`/`id` is given but does not match any Capability, `null` will be emitted)
	 */
	public getCapabilities(name?: string): Observable<Array<Capability> | Capability | null> {
		const path = `/api/${this.API_VERSION}/capabilities`;
		if (name) {
			return this.get(`${path}?name=${encodeURIComponent(name)}`).pipe(map(
				r => {
					for (const cap of (r.body.response as Array<Capability>)) {
						if (cap.name === name) {
							return cap;
						}
					}
					return null;
				}
			));
		}
		return this.get(path).pipe(map(
			r => {
				return r.body.response as Array<Capability>;
			}
		));
	}

	constructor(http: HttpClient) {
		super(http);
	}

}
