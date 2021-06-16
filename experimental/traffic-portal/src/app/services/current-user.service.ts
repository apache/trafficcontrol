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
import { Router } from "@angular/router";
import { Capability, User } from "../models";

/**
 * This service keeps track of the currently authenticated user.
 *
 * This needs to be done separately from the AuthenticationService's
 * methods, because those depend on the API services and the API services use
 * an implicitly injected ErrorInterceptor which clears the authenticated user
 * value when it hits a 401 error - so that would be a circular dependency.
 */
@Injectable({
	providedIn: "root"
})
export class CurrentUserService {
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

	constructor(private readonly router: Router) {}

	/**
	 * Sets the currently authenticated user.
	 *
	 * @param u The new user who has been authenticated.
	 * @param caps The newly authenticated user's Permissions.
	 */
	public setUser(u: User, caps: Set<string> | Array<Capability>): void {
		this.user = u;
		this.caps = caps instanceof Array ? new Set(caps.map(c=>c.name)) : caps;
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
		}
		console.log("query params:", queryParams);
		this.router.navigate(["/login"], {queryParams});
	}
}
