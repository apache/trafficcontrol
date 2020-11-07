/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/
import { Component, OnInit, OnDestroy } from "@angular/core";

import { Subscription } from "rxjs";
import { first } from "rxjs/operators";

import { User } from "../../models";
import { AuthenticationService } from "../../services";

/**
 * CurrentuserComponent is the controller for the current user's profile page.
 */
@Component({
	selector: "tp-currentuser",
	styleUrls: ["./currentuser.component.scss"],
	templateUrl: "./currentuser.component.html"
})
export class CurrentuserComponent implements OnInit, OnDestroy {

	/** The currently logged-in user - or 'null' if not logged-in. */
	public currentUser: User | null;
	private subscription: Subscription;

	constructor (private readonly auth: AuthenticationService) {
	}

	/**
	 * Runs initialization, setting the currently logged-in user from the
	 * authentication service.
	 */
	public ngOnInit(): void {
		this.currentUser = this.auth.currentUserValue;
		if (this.currentUser === null) {
			this.auth.updateCurrentUser().pipe(first()).subscribe(
				(r: boolean) => {
					if (r) {
						this.currentUser = this.auth.currentUserValue;
					}
				}
			);
		}
		this.subscription = this.auth.currentUser.subscribe(
			(u: User) => {
				this.currentUser = u;
			}
		);
	}

	/**
	 * Runs when the component is destroyed - cleans up active subscriptions.
	 */
	public ngOnDestroy(): void {
		this.subscription.unsubscribe();
	}
}
