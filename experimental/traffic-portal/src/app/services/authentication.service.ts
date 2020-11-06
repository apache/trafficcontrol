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

import { BehaviorSubject, Observable } from "rxjs";
import { first, map } from "rxjs/operators";

import { UserService } from "./api";

import { Role, User } from "../models";

/**
 * AuthenticationService handles authentication with the Traffic Ops server and
 * providing properties of the current user to service consumers.
 */
@Injectable({ providedIn: "root" })
export class AuthenticationService {
	private readonly currentUserSubject: BehaviorSubject<User | null>;

	/** An observable that emits the current user, or 'null' if they are not logged in. */
	public currentUser: Observable<User | null>;
	private readonly loggedInSubject: BehaviorSubject<boolean>;

	/** An Observable that emits whether or not the current user is logged in. */
	public loggedIn: Observable<boolean>;
	private readonly currentUserCapabilitiesSubject: BehaviorSubject<Set<string>>;

	/** An Observable that emits the current user's capabilities. */
	public currentUserCapabilities: Observable<Set<string>>;

	constructor (private readonly api: UserService) {
		this.currentUserSubject = new BehaviorSubject<User | null>(null);
		this.loggedInSubject = new BehaviorSubject<boolean>(false);
		this.currentUserCapabilitiesSubject = new BehaviorSubject<Set<string>>(new Set<string>());
		this.currentUser = this.currentUserSubject.asObservable();
		this.loggedIn = this.loggedInSubject.asObservable();
		this.currentUserCapabilities = this.currentUserCapabilitiesSubject.asObservable();
	}

	/** The current user's User, or 'null' if they are not logged in. */
	public get currentUserValue (): User | null {
		return this.currentUserSubject.value;
	}

	/** Whether or not the current user is logged in. */
	public get loggedInValue (): boolean {
		return this.loggedInSubject.value;
	}

	/** The Capabilities of the current user. */
	public get currentUserCapabilitiesValue (): Set<string> {
		return this.currentUserCapabilitiesSubject.value;
	}

	/**
	 * Updates the current user, and provides a way for callers to check if the update was succesful
	 * @returns An `Observable` which will emit a boolean value indicating the success of the update
	 */
	public updateCurrentUser (): Observable<boolean> {
		return this.api.getCurrentUser().pipe(first()).pipe(map(
			(u: User) => {
				this.currentUserSubject.next(u);
				if (u.role) {
					this.api.getRoles(u.role).pipe(first()).pipe(map(
						(r: Role) => {
							this.currentUserCapabilitiesSubject.next(new Set<string>(r.capabilities));
						}
					));
				}
				return true;
			},
			e => {
				console.error("Failed to update current user");
				console.debug("User update error: ", e);
				return false;
			}
		));
	}

	/**
	 * Logs in a user and, on successful login, updates the current user.
	 *
	 * @param u The user's username.
	 * @param p The user's password.
	 */
	public login (u: string, p: string): Observable<boolean> {
		return this.api.login(u, p).pipe(map(
			(resp) => {
				if (resp && resp.status === 200) {
					this.loggedInSubject.next(true);
					this.updateCurrentUser();
					return true;
				}
				return false;
			}
		));
	}

	/** Logs the currently logged-in user out. */
	public logout(): void {
		this.currentUserSubject.next(null);
		this.loggedInSubject.next(false);
	}
}
