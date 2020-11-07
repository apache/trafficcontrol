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
import { Component, OnInit } from "@angular/core";
import { FormControl } from "@angular/forms";

import { BehaviorSubject, Observable } from "rxjs";
import { first } from "rxjs/operators";

import { Role, User } from "../../models";
import { AuthenticationService } from "../../services";
import { UserService } from "../../services/api";
import { orderBy } from "../../utils";

/**
 * UsersComponent is the controller for the "users" page.
 */
@Component({
	selector: "tp-users",
	styleUrls: ["./users.component.scss"],
	templateUrl: "./users.component.html"
})
export class UsersComponent implements OnInit {

	/** All (visible) users. */
	public users: Array<User>;

	/** Fuzzy search control. */
	public fuzzControl = new FormControl("");

	/** Whether or not user data is still loading. */
	public loading: boolean;

	/** The ID of the currently logged-in user. */
	public myId: number;

	private readonly rolesMapSubject: BehaviorSubject<Map<number, string>>;

	/** Maps role IDs to role Names. */
	public rolesMap: Observable<Map<number, string>>;

	constructor (private readonly api: UserService, private readonly auth: AuthenticationService) {
		this.rolesMapSubject = new BehaviorSubject<Map<number, string>>(new Map<number, string>());
		this.rolesMap = this.rolesMapSubject.asObservable();
		this.users = new Array<User>();
		this.loading = true;
		this.myId = -1;
	}

	/**
	 * Initializes data like a map of role ids to their names.
	 */
	public ngOnInit(): void {
		// User may have navigated directly with a valid cookie - in which case current user is null
		if (this.auth.currentUserValue === null) {
			this.auth.updateCurrentUser().subscribe(
				v => {
					if (v && this.auth.currentUserValue) {
						this.myId = this.auth.currentUserValue.id;
					}
				}
			);
		} else {
			this.myId = this.auth.currentUserValue.id;
		}

		this.api.getUsers().pipe(first()).subscribe(
			(r: Array<User>) => {
				this.users = orderBy(r, "fullName") as Array<User>;
				this.loading = false;
			}
		);

		this.api.getRoles().pipe(first()).subscribe(
			(roles: Array<Role>) => {
				const roleMap = new Map<number, string>();
				for (const r of roles) {
					roleMap.set(r.id, r.name);
				}
				this.rolesMapSubject.next(roleMap);
			}
		);
	}

	/**
	 * Implements a fuzzy search over usernames
	 * @param u The user being checked for a fuzzy match (currently uses the username)
	 * @returns `true` if `u` is a fuzzy match for the `fuzzControl` value, `false` otherwise
	 */
	public fuzzy (u: User): boolean {
		if (!this.fuzzControl.value) {
			return true;
		}
		const testVal = u.username.toLocaleLowerCase();
		let n = -1;
		for (const l of this.fuzzControl.value.toLocaleLowerCase()) {
			/* tslint:disable */
			if (!~(n = testVal.indexOf(l, n + 1))) {
			/* tslint:enable */
				return false;
			}
		}
		return true;
	}

}
