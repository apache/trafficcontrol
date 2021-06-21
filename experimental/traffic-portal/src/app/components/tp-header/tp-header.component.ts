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
import { AuthenticationService } from "src/app/services";
import { UserService } from "src/app/services/api";

/**
 * TpHeaderComponent is the controller for the standard Traffic Portal header.
 */
@Component({
	selector: "tp-header",
	styleUrls: ["./tp-header.component.scss"],
	templateUrl: "./tp-header.component.html"
})
export class TpHeaderComponent implements OnInit {

	/**
	 * The set of permissions available to the authenticated user.
	 */
	public permissions = new Set<string>();

	/**
	 * Holds a continuous subscription for the current user's permissions, in case they change.
	 */
	// private permissionSubscription: Subscription | undefined;

	/**
	 * The title to be used in the header.
	 *
	 * If not given, defaults to "Traffic Portal".
	 */
	@Input() public title?: string;

	/** Constructor */
	constructor(private readonly auth: AuthenticationService, private readonly api: UserService) {
	}

	/** Sets up data dependencies. */
	public ngOnInit(): void {
		this.permissions = this.auth.capabilities;
	}

	/**
	 * Handles when the user clicks the "Logout" button by using the API to
	 * invalidate their session before redirecting them to the login page.
	 */
	public logout(): void {
		this.api.logout().then(
			r => {
				if (!r) {
					console.warn("Failed to log out - clearing user data anyway!");
				}
				this.auth.logout();
			}
		);
	}
}
