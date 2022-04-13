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

import {AfterViewInit, Component, OnInit, ViewChild} from "@angular/core";
import { Router } from "@angular/router";

import { CurrentUser } from "src/app/models";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import {TpHeaderComponent} from "src/app/shared/tp-header/tp-header.component";
import {TpHeaderService} from "src/app/shared/tp-header/tp-header.service";

/**
 * The most basic component that contains everything else. This should be kept pretty simple.
 */
@Component({
	selector: "app-root",
	styleUrls: ["./app.component.scss"],
	templateUrl: "./app.component.html",
})
export class AppComponent implements OnInit, AfterViewInit {
	// @ts-ignore
	@ViewChild("header") private readonly headerRef: TpHeaderComponent;
	/** The app's title */
	public title = "Traffic Portal";

	/** The currently logged-in user */
	public currentUser: CurrentUser | null = null;

	constructor(private readonly router: Router, private readonly auth: CurrentUserService, private readonly headerSvc: TpHeaderService) {
	}

	/**
	 * Logs the currently logged-in user out.
	 */
	public logout(): void {
		this.auth.logout();
		this.router.navigate(["/login"]);
	}

	/**
	 * Angular lifecycle hook, sets header text.
	 */
	public ngAfterViewInit(): void {
		this.headerSvc.registerHeader(this.headerRef);
	}

	/**
	 * Sets up the current user.
	 */
	public ngOnInit(): void {
		this.auth.userChanged.subscribe(user => {
			this.currentUser = user;
		});
	}

}
