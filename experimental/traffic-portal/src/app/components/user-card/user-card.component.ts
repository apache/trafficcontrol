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
import { Component, Input, OnInit } from "@angular/core";

import { Observable } from "rxjs";
import { first } from "rxjs/operators";

import { Role, User } from "../../models";
import { UserService } from "../../services/api";

/**
 * UserCardComponent is the controller for collapsible cards containing
 * information about a user.
 */
@Component({
	selector: "user-card",
	styleUrls: ["./user-card.component.scss"],
	templateUrl: "./user-card.component.html"
})
export class UserCardComponent implements OnInit {

	/**
	 * The User being modeled.
	 */
	@Input() public user: User;

	/**
	 * An optional map of role IDs to role names, since often the API will omit role names in user objects
	 * If this isn't given, this component will make its own HTTP request to find out its user's role name
	 */
	@Input() public roleMap?: Observable<Map<number, string>>;

	constructor (private readonly api: UserService) { }

	/** Initializes data like mapping role IDs to role names. */
	public ngOnInit(): void {
		this.user = this.user;
		if (this.user.role === undefined) {
			console.error("User appears to have no role");
			console.debug(this.user);
			return;
		}
		if (!this.user.roleName) {
			if (!this.roleMap) {
				this.api.getRoles(this.user.role).pipe(first()).subscribe(
					(role: Role) => {
						this.user.roleName = role.name;
					}
				);
			} else {
				this.roleMap.subscribe(
					m => {
						// initial value will be null
						if (!m) {
							return;
						}
						if (m.has(this.user.role as number)) {
							this.user.roleName = m.get(this.user.role as number);
						} else {
							console.error("Role ID %d not found in Traffic Ops!");
						}
					}
				);
			}
		}
		// Go emits marshaled JSON date/time structs in a format only Chrome can parse. Because, you know, Google is web standard.
		if (typeof(this.user.lastUpdated) === "string") {
			this.user.lastUpdated = new Date((this.user.lastUpdated as string).replace("-", "/").replace("-", "/").replace("+", " GMT+"));
		}
	}

	/**
	 * Checks if the user has at least on populated "location" field (city,
	 * stateOrProvince etc.).
	 */
	public userHasLocation(): boolean {
		return this.user.city !== null || this.user.stateOrProvince !== null || this.user.country !== null || this.user.postalCode !== null;
	}

	/**
	 * Gets a string representing a user's address, or 'null' if one cannot be
	 * constructed because no relevant information exists.
	 */
	public userLocationString(): string | null {
		let ret = "";
		if (this.user.city) {
			ret += this.user.city;
		}
		if (this.user.stateOrProvince) {
			if (ret.length !== 0) {
				ret += ", ";
			}
			ret += this.user.stateOrProvince;
		}
		if (this.user.country) {
			if (ret.length !== 0) {
				ret += ", ";
			}
			ret += this.user.country;
		}
		if (this.user.postalCode) {
			if (ret.length !== 0) {
				ret += ", ";
			}
			ret += this.user.postalCode;
		}

		return ret.length === 0 ? null : ret;
	}
}
