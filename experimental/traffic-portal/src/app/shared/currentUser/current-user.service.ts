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
import { Router } from "@angular/router";
import { BehaviorSubject } from "rxjs";

import { UserService } from "src/app/api";
import { type Capability, type CurrentUser, ADMIN_ROLE } from "src/app/models";

/**
 * This service keeps track of the currently authenticated user.
 *
 * This needs to be done separately from the CurrentUserService's
 * methods, because those depend on the API services and the API services use
 * an implicitly injected ErrorInterceptor which clears the authenticated user
 * value when it hits a 401 error - so that would be a circular dependency.
 */
@Injectable()
export class CurrentUserService {
	/** Makes updateCurrentUser able to be called from multiple places without regard to order */
	private updatingUserPromise: Promise<boolean> | null = null;
	/** To allow downstream code to stay up to date with the current user */
	public userChanged = new EventEmitter<CurrentUser>();
	/** The currently authenticated user - or `null` if not authenticated. */
	private user: CurrentUser | null = null;

	/** The currently authenticated user - or `null` if not authenticated. */
	public get currentUser(): CurrentUser | null {
		return this.user;
	}

	/** The Permissions afforded to the currently authenticated user. */
	public capabilities: BehaviorSubject<Set<string>> = new BehaviorSubject(new Set());

	/** Whether or not the user is authenticated. */
	public get loggedIn(): boolean {
		return this.currentUser !== null;
	}

	constructor(private readonly router: Router, private readonly api: UserService) {
		this.updateCurrentUser();
	}

	/**
	 * Gets the current user if currentuser is not already set
	 *
	 * @returns A promise containing the value indicating the success of the update
	 */
	public async fetchCurrentUser(): Promise<boolean> {
		if (this.currentUser !== null){
			return true;
		}
		return this.updateCurrentUser();
	}

	/**
	 * Updates the current user, and provides a way for callers to check if the update was successful.
	 *
	 * @returns A boolean value indicating the success of the update
	 */
	public async updateCurrentUser(): Promise<boolean> {
		if (this.updatingUserPromise === null) {
			this.updatingUserPromise = this.api.getCurrentUser().then(
				async u => {
					if (u.role === undefined) {
						throw new Error("current user had no Role");
					}
					const role = await this.api.getRoles(u.role);
					this.setUser(u, new Set(role.capabilities));
					return true;
				}
			).catch(
				e => {
					console.error("Failed to update current user:", e);
					return false;
				}
			).finally(() => this.updatingUserPromise = null );
		}
		return this.updatingUserPromise;
	}

	/**
	 * Saves the user
	 *
	 * @param user User to e saved
	 * @returns A promise returning the status of the update.
	 */
	public async saveCurrentUser(user: CurrentUser): Promise<boolean> {
		return this.api.updateCurrentUser(user);
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
		const resp = await this.api.login(uOrT, p);
		if (resp && resp.status === 200) {
			return this.updateCurrentUser();
		}
		return false;
	}

	/**
	 * Sets the currently authenticated user.
	 *
	 * @param u The new user who has been authenticated.
	 * @param caps The newly authenticated user's Permissions.
	 */
	public setUser(u: CurrentUser, caps: Set<string> | Array<Capability>): void {
		this.user = u;
		const capabilities = caps instanceof Array ? new Set(caps.map(c=>c.name)) : caps;
		this.userChanged.emit(this.user);
		this.capabilities.next(capabilities);
	}

	/**
	 * Checks if the user has a given Permission.
	 *
	 * @param perm The Permission in question.
	 * @returns `true` if the user has the Permission `perm`, `false` otherwise.
	 */
	public hasPermission(perm: string): boolean {
		if (!this.user) {
			return false;
		}
		return this.user.roleName === ADMIN_ROLE || this.capabilities.getValue().has(perm);
	}

	/**
	 * Clears authentication data associated with the current user, and
	 * redirects to login.
	 *
	 * @param withRedirect If given and `true`, will redirect with the
	 * `returnUrl` query string parameter set to the current route.
	 */
	public logout(withRedirect?: boolean): void {
		this.user = null;
		this.capabilities.next(new Set());

		const queryParams: Record<string | symbol, string> = {};
		if (withRedirect) {
			queryParams.returnUrl = this.router.url;
			if (queryParams.returnUrl.startsWith("/login")) {
				queryParams.returnUrl = "/core";
			}
		}
		this.router.navigate(["/login"], {queryParams});
	}
}
