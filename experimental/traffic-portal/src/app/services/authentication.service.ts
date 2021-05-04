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
import { Injectable } from "@angular/core";
import { User } from "../models";

import { UserService } from "./api";
import { CurrentUserService } from "./current-user.service";

/**
 * AuthenticationService handles authentication with the Traffic Ops server and
 * providing properties of the current user to service consumers.
 */
@Injectable({ providedIn: "root" })
export class AuthenticationService {

	/**
	 * The currently authenticated user - or `null` if not authenticated.
	 */
	public get currentUser(): User | null {
		return this.currentUserService.currentUser;
	}

	/**
	 * All of the Permissions afforded to the currently authenticated user.
	 */
	public get capabilities(): Set<string> {
		return this.currentUserService.capabilities;
	}

	/**
	 * Constructs the service with its required dependencies injected.
	 *
	 * @param api A reference to the UserService.
	 */
	constructor(private readonly api: UserService, private readonly currentUserService: CurrentUserService) {
		this.updateCurrentUser();
	}

	/**
	 * Updates the current user, and provides a way for callers to check if the update was succesful.
	 *
	 * @param token If given, the service will first attempt to login using this token.
	 * @returns A boolean value indicating the success of the update
	 */
	public async updateCurrentUser(token?: string): Promise<boolean> {
		if (token) {
			if (!(await this.api.login(token))) {
				console.error("invalid token");
				return false;
			}
		}
		return this.api.getCurrentUser().then(
			async u => {
				if (u.role === undefined) {
					throw new Error("current user had no Role");
				}
				const role = await this.api.getRoles(u.role);
				this.currentUserService.setUser(u, new Set(role.capabilities));
				return true;
			}
		).catch(
			e => {
				console.error("Failed to update current user:", e);
				return false;
			}
		);
	}

	/**
	 * Logs in a user and, on successful login, updates the current user.
	 *
	 * @param uOrT The user's username, if `p` is given. If `p` is *not* given,
	 * this is treated as a login token.
	 * @param p The user's password.
	 * @returns An observable that emits whether or not login succeeded.
	 */
	public async login(uOrT: string, p?: string): Promise<boolean> {
		return this.api.login(uOrT, p).then(
			async resp => {
				if (resp && resp.status === 200) {
					return this.updateCurrentUser();
				}
				return false;
			}
		);
	}

	/** Logs the currently logged-in user out. */
	public logout(): void {
		this.currentUserService.logout();
	}
}
