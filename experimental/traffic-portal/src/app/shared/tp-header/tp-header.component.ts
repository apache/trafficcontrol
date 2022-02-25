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
import { Component, Input } from "@angular/core";

import { UserService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";

/**
 * TpHeaderComponent is the controller for the standard Traffic Portal header.
 */
@Component({
	selector: "tp-header",
	styleUrls: ["./tp-header.component.scss"],
	templateUrl: "./tp-header.component.html"
})
export class TpHeaderComponent {

	/**
	 * The title to be used in the header.
	 *
	 * If not given, defaults to "Traffic Portal".
	 */
	@Input() public title?: string;

	constructor(private readonly auth: CurrentUserService, private readonly api: UserService) {
	}

	/**
	 * Checks for a Permission afforded to the currently authenticated user.
	 *
	 * @param perm The Permission for which to check.
	 * @returns Whether the currently authenticated user has the Permission
	 * `perm`.
	 */
	public hasPermission(perm: string): boolean {
		return this.auth.hasPermission(perm);
	}

	/**
	 * Handles when the user clicks the "Logout" button by using the API to
	 * invalidate their session before redirecting them to the login page.
	 */
	public async logout(): Promise<void> {
		if (!(await this.api.logout())) {
			console.warn("Failed to log out - clearing user data anyway!");
		}
		this.auth.logout();
	}
}
