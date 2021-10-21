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
import {EventEmitter, Injectable} from "@angular/core";
import { Router } from "@angular/router";
import {UserService} from "src/app/shared/api";
import { Capability, User } from "../../models";

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
	public userChanged: EventEmitter<User> = new EventEmitter<User>();
	/** The currently authenticated user - or `null` if not authenticated. */
	private user: User | null = null;
	/** The Permissions afforded to the currently authenticated user. */
	private caps = new Set<string>();

	/** The currently authenticated user - or `null` if not authenticated. */
	public get currentUser(): User | null {
		return this.user;
	}
	/** The Permissions afforded to the currently authenticated user. */
	public get capabilities(): Set<string> {
		return this.caps;
	}

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
	 * @returns Promise<boolean> a value indicating the success of the update
	 */
	public async fetchCurrentUser(): Promise<boolean> {
		if(this.currentUser !== null){
			return new Promise<boolean>(resolve => resolve(true));
		}
		return this.updateCurrentUser();
	}

	/**
	 * Updates the current user, and provides a way for callers to check if the update was succesful.
	 *
	 * @returns Promise<boolean> a value indicating the success of the update
	 */
	public async updateCurrentUser(): Promise<boolean> {
		if (this.updatingUserPromise == null) {
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
	 * @returns Promise<boolean> promise returning the status of the update.
	 */
	public async saveCurrentUser(user: User): Promise<boolean> {
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
		return this.api.login(uOrT, p).then(
			async resp => {
				if (resp && resp.status === 200) {
					return this.updateCurrentUser();
				}
				return false;
			}
		);
	}

	/**
	 * Sets the currently authenticated user.
	 *
	 * @param u The new user who has been authenticated.
	 * @param caps The newly authenticated user's Permissions.
	 */
	public setUser(u: User, caps: Set<string> | Array<Capability>): void {
		this.user = u;
		this.caps = caps instanceof Array ? new Set(caps.map(c=>c.name)) : caps;
		this.userChanged.emit(this.user);
	}

	/**
	 * Checks if the user has a given Permission.
	 *
	 * @param perm The Permission in question.
	 * @returns `true` if the user has the Permission `perm`, `false` otherwise.
	 */
	public hasPermission(perm: string): boolean {
		return this.user ? this.caps.has(perm) : false;
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
		this.caps.clear();

		const queryParams: Record<string | symbol, string> = {};
		if (withRedirect) {
			queryParams.returnUrl = this.router.url;
			if (queryParams.returnUrl.startsWith("/login")) {
				queryParams.returnUrl = "/core";
			}
		}
		console.log("query params:", queryParams);
		this.router.navigate(["/login"], {queryParams});
	}
}
