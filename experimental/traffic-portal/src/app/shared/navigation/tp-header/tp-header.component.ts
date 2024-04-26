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

import { LoggingService } from "src/app/shared/logging.service";
import { HeaderNavigation, HeaderNavType, NavigationService } from "src/app/shared/navigation/navigation.service";
import { ThemeManagerService } from "src/app/shared/theme-manager/theme-manager.service";

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
		this.navSvc.headerTitle.subscribe(title => {
			this.title = title;
		});
		this.navSvc.headerHidden.subscribe(hidden => {
			this.hidden = hidden;
		});
		this.navSvc.horizontalNavsUpdated.subscribe(navs => {
			this.horizNavs = navs;
		});
		this.navSvc.verticalNavsUpdated.subscribe(navs => {
			this.vertNavs = navs;
		});
	}

	constructor(
		public readonly themeSvc: ThemeManagerService,
		private readonly navSvc: NavigationService,
		private readonly log: LoggingService,
	) { }

	/**
	 * Calls a navs click function, throws an error if null
	 *
	 * @param nav nav to process
	 */
	public navClick(nav: HeaderNavigation): void {
		if(nav.click === undefined) {
			throw new Error(`nav ${nav.text} does not have a click function`);
		} else {
			nav?.click();
		}
	}

	/**
	 * Gets a navs routerLink, logs an error if null
	 *
	 * @param nav nav to process
	 * @returns routerLink
	 */
	public navRouterLink(nav: HeaderNavigation): string {
		if(nav.routerLink === undefined) {
			this.log.error(`nav ${nav.text} does not have a routerLink`);
			return "";
		}
		return nav.routerLink;

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
}
