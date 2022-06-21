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
import {Component, OnInit} from "@angular/core";

import { UserService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import {ThemeManagerService} from "src/app/shared/theme-manager/theme-manager.service";
import {HeaderNavigation, HeaderNavType, TpHeaderService} from "src/app/shared/tp-header/tp-header.service";

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
	 * The title to be used in the header.
	 *
	 * If not given, defaults to "Traffic Portal".
	 */
	public title = "";

	public hidden = false;

	// Will try to display each of these navs on the header, space allowing.
	public horizNavs: Array<HeaderNavigation> = new Array<HeaderNavigation>();
	// Navs that are not directly displayed on the header.
	public vertNavs: Array<HeaderNavigation> = new Array<HeaderNavigation>();

	/**
	 * Angular lifecycle hook
	 */
	public ngOnInit(): void {
		this.headerSvc.addHorizontalNav({
			routerLink: "/core",
			text: "Home",
			type: "anchor",
		}, "home");
		this.headerSvc.addHorizontalNav({
			routerLink: "/core/users",
			text: "Users",
			type: "anchor",
			visible: () => this.hasPermission("USER:READ"),
		}, "Users");
		this.headerSvc.addHorizontalNav({
			routerLink: "/core/servers",
			text: "Servers",
			type: "anchor",
			visible: () => this.hasPermission("SERVER:READ"),
		}, "Servers");
		this.headerSvc.addVerticalNav({
			routerLink: "/core/me",
			text: "Profile",
			type: "anchor"
		}, "Profile");
		this.headerSvc.addVerticalNav({
			click: async () => this.logout(),
			text: "Logout",
			type: "anchor"
		}, "Logout");

		this.headerSvc.headerTitle.subscribe(title => {
			this.title = title;
		});
		this.headerSvc.headerHidden.subscribe(hidden => {
			this.hidden = hidden;
		});
		this.headerSvc.horizontalNavsUpdated.subscribe(navs => {
			this.horizNavs = navs;
		});
		this.headerSvc.verticalNavsUpdated.subscribe(navs => {
			this.vertNavs = navs;
		});
	}

	constructor(private readonly auth: CurrentUserService, private readonly api: UserService,
		public readonly themeSvc: ThemeManagerService, private readonly headerSvc: TpHeaderService) {
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
	 * Checks if a nav is shown
	 *
	 * @param nav nav to check
	 * @param type which type of nav to check for
	 * @returns If the nav should be rendered
	 */
	public navShown(nav: HeaderNavigation, type: HeaderNavType): boolean {
		return nav.type === type && (nav.visible === undefined || nav.visible());
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
