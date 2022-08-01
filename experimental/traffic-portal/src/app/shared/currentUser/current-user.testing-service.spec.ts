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
import { EventEmitter, Injectable } from "@angular/core";
import { BehaviorSubject } from "rxjs";

import { ADMIN_ROLE, Capability, CurrentUser } from "src/app/models";

/**
 * This is a mock for the {@link CurrentUserService} service for testing.
 *
 * The authenticated user it "manages" is perpetually authenticated, but logging
 * in can fail - it expects the credentials to match the existing currentUser's
 * username with the password 'twelve12!' (determined by the static 'PASSWORD'
 * property of the service).
 */
@Injectable()
export class CurrentUserTestingService {
	public static readonly PASSWORD = "twelve12!";
	public userChanged = new EventEmitter<CurrentUser>();
	public currentUser: CurrentUser = {
		addressLine1: null,
		addressLine2: null,
		city: null,
		company: null,
		country: null,
		email: "a@b.c",
		fullName: "admin",
		gid: null,
		id: 1,
		lastUpdated: new Date(0),
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
		username: "admin"
	};
	public capabilities: BehaviorSubject<Set<string>> = new BehaviorSubject(new Set(["ALL"]));
	public readonly loggedIn = true;

	/**
	 * Gets the current user if currentuser is not already set
	 *
	 * @returns A promise containing the value indicating the success of the update
	 */
	public async fetchCurrentUser(): Promise<boolean> {
		return true;
	}

	/**
	 * Updates the current user, and provides a way for callers to check if the update was successful.
	 *
	 * @returns A boolean value indicating the success of the update
	 */
	public async updateCurrentUser(): Promise<boolean> {
		return true;
	}

	/**
	 * Saves the user
	 *
	 * @param user User to e saved
	 * @returns A promise returning the status of the update.
	 */
	public async saveCurrentUser(user: CurrentUser): Promise<boolean> {
		this.currentUser = user;
		return true;
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
		return uOrT === this.currentUser.username && p === CurrentUserTestingService.PASSWORD;
	}

	/**
	 * Sets the currently authenticated user.
	 *
	 * @param u The new user who has been authenticated.
	 * @param caps The newly authenticated user's Permissions.
	 */
	public setUser(u: CurrentUser, caps: Set<string> | Array<Capability>): void {
		this.currentUser = u;
		const capabilities = Array.isArray(caps) ? new Set(caps.map(c=>c.name)) : caps;
		this.userChanged.emit(this.currentUser);
		this.capabilities.next(capabilities);
	}

	/**
	 * Checks if the user has a given Permission.
	 *
	 * @param perm The Permission in question.
	 * @returns `true` if the user has the Permission `perm`, `false` otherwise.
	 */
	public hasPermission(perm: string): boolean {
		return this.currentUser.roleName === ADMIN_ROLE || this.capabilities.getValue().has(perm);
	}

	/**
	 * Mocks the {@link CurrentUserService}'s logout method. Note that
	 * regardless of what's passed in as an argument, no navigation is performed
	 * by the testing service.
	 *
	 * @param withRedirect If given and `true`, prints a warning that the
	 * testing service doesn't navigate.
	 */
	public logout(withRedirect?: boolean): void {
		if (withRedirect) {
			console.warn("testing service does not navigate!");
		}
	}
}
