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

import { UserService } from "src/app/api";
import type { Role, User } from "src/app/models";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import { orderBy } from "src/app/utils";

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

	/** An observation subject for the map of a user's Roles. */
	private readonly rolesMapSubject: BehaviorSubject<Map<number, string>>;

	/** Maps role IDs to role Names. */
	public rolesMap: Observable<Map<number, string>>;

	/**
	 * Constructor.
	 */
	constructor(private readonly api: UserService, private readonly auth: CurrentUserService) {
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
		if (!this.auth.currentUser) {
			this.auth.updateCurrentUser().then(
				v => {
					if (v && this.auth.currentUser) {
						this.myId = this.auth.currentUser.id;
					}
				}
			);
		} else {
			this.myId = this.auth.currentUser.id;
		}

		this.api.getUsers().then(
			r => {
				this.users = orderBy(r, "fullName");
				this.loading = false;
			}
		);

		this.api.getRoles().then(
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
	 *
	 * @param u The user being checked for a fuzzy match (currently uses the username)
	 * @returns `true` if `u` is a fuzzy match for the `fuzzControl` value, `false` otherwise
	 */
	public fuzzy(u: User): boolean {
		if (!this.fuzzControl.value) {
			return true;
		}
		const testVal = u.username.toLocaleLowerCase();
		let n = -1;
		for (const l of this.fuzzControl.value.toLocaleLowerCase()) {
			/* eslint-disable */
			if (!~(n = testVal.indexOf(l, n + 1))) {
			/* eslint-enable */
				return false;
			}
		}
		return true;
	}

	/**
	 * Checks if the user has any render-able address piece(s).
	 *
	 * @param user The user to check.
	 * @returns 'true' if the user has at least one populated "location" field (city,
	 * stateOrProvince etc.), 'false' otherwise.
	 */
	public userHasLocation(user: User): boolean {
		return user.city !== null || user.stateOrProvince !== null || user.country !== null || user.postalCode !== null;
	}

	/**
	 * Gets a string representing a user's address.
	 *
	 * @param user The user for whom to fetch a location string.
	 * @returns The user's address, or 'null' if one cannot be
	 * constructed because no relevant information exists.
	 */
	public userLocationString(user: User): string | null {
		let ret = "";
		if (user.city) {
			ret += user.city;
		}
		if (user.stateOrProvince) {
			if (ret.length !== 0) {
				ret += ", ";
			}
			ret += user.stateOrProvince;
		}
		if (user.country) {
			if (ret.length !== 0) {
				ret += ", ";
			}
			ret += user.country;
		}
		if (user.postalCode) {
			if (ret.length !== 0) {
				ret += ", ";
			}
			ret += user.postalCode;
		}

		return ret.length === 0 ? null : ret;
	}
}
