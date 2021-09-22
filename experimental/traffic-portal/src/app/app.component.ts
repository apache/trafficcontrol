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
import {Router, RouterOutlet} from "@angular/router";

import { User } from "./models";
import {AuthenticationService} from "./shared/authentication/authentication.service";
import {animate, group, query, style, transition, trigger} from "@angular/animations";

const swipeLeft = [
	query(":enter, :leave", style({height: "100%", position: "fixed", padding: "0 5px 0 5px", width: "calc(100%)"}),
		{optional: true}
	),
	group([  // block executes in parallel
		query(":enter", [
			style({transform: "translateX(-100%)", opacity: "0"}),
			animate("0.5s ease-in-out", style({transform: "translateX(0%)", opacity: ".9"}))
		], {optional: true}),
		query(":leave", [
			style({transform: "translateX(0%)", opacity: ".9"}),
			animate("0.5s ease-in-out", style({transform: "translateX(100%)", opacity: "0"})),
		], {optional: true}),
	])
];

const swipeRight = [
	query(":enter, :leave", style({height: "100%", position: "fixed", width: "calc(100%)"}),
		{optional: true}
	),
	group([  // block executes in parallel
		query(":enter", [
			style({transform: "translateX(100%)"}),
			animate("0.5s ease-in-out", style({transform: "translateX(0%)"}))
		], {optional: true}),
		query(":leave", [
			style({transform: "translateX(0%)"}),
			animate("0.5s ease-in-out", style({transform: "translateX(-100%)"}))
		], {optional: true}),
	])
];

/**
 * The most basic component that contains everything else. This should be kept pretty simple.
 */
@Component({
	animations: [
		trigger("routerAnimations", [
			transition("users => servers", swipeRight),
			transition("servers => users", swipeLeft)
		])
	],
	selector: "app-root",
	styleUrls: ["./app.component.scss"],
	templateUrl: "./app.component.html",
})
export class AppComponent implements OnInit {
	/** The app"s title */
	public title = "Traffic Portal";

	/** The currently logged-in user */
	public currentUser: User | null = null;

	/**
	 * Constructor.
	 */
	constructor(private readonly router: Router, private readonly auth: AuthenticationService) {
	}

	/**
	 * Logs the currently logged-in user out.
	 */
	public logout(): void {
		this.auth.logout();
		this.router.navigate(["/login"]);
	}

	public getState(route: RouterOutlet) {
		return route && route.activatedRouteData && route.activatedRouteData.animation;
	}

	/**
	 * Sets up the current user.
	 */
	public ngOnInit(): void {
		this.auth.updateCurrentUser().then(
			success =>  {
				if (success) {
					this.currentUser = this.auth.currentUser;
				}
			}
		);
	}

}
