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
import { Injectable } from "@angular/core";
import {ReplaySubject} from "rxjs";

/**
 * Defines the type of the header nav
 */
export declare type HeaderNavType = "anchor" | "button";

/**
 * Specifies the setting for the nav
 */
export interface HeaderNavigation {
	type: HeaderNavType;
	visible?: () => boolean;
	routerLink?: string;
	click?: () => Promise<void>;
	text: string;
}
/**
 *
 */
@Injectable({
	providedIn: "root"
})
export class TpHeaderService {
	public readonly headerTitle: ReplaySubject<string>;
	public readonly headerHidden: ReplaySubject<boolean>;
	public readonly horizontalNavsUpdated: ReplaySubject<Array<HeaderNavigation>>;
	public readonly verticalNavsUpdated: ReplaySubject<Array<HeaderNavigation>>;

	private readonly horizontalNavs: Map<string, HeaderNavigation>;
	private readonly verticalNavs: Map<string, HeaderNavigation>;

	constructor() {
		this.horizontalNavs = new Map<string, HeaderNavigation>();
		this.verticalNavs = new Map<string, HeaderNavigation>();
		this.horizontalNavsUpdated = new ReplaySubject();
		this.verticalNavsUpdated = new ReplaySubject();
		this.headerTitle = new ReplaySubject(1);
		this.headerHidden = new ReplaySubject(1);
		this.headerHidden.next(false);
	}

	/**
	 * Builds the horizontal header navigation array for consumption.
	 *
	 * @returns Header Navs
	 */
	private buildHorizontalNavs(): Array<HeaderNavigation> {
		return Array.from(this.horizontalNavs.values());
	}

	/**
	 * Builds the vertical header navigation array for consumption.
	 *
	 * @returns Header Navs
	 */
	private buildVerticalNavs(): Array<HeaderNavigation> {
		return Array.from(this.verticalNavs.values());
	}

	/**
	 * Removes a nav from the list. Does not throw an exception if not exists
	 *
	 * @param key key to delete by
	 * @returns boolean indicating if a nav was deleted
	 */
	public removeHorizontalNav(key: string): boolean {
		return this.horizontalNavs.delete(key);
	}

	/**
	 * Removes a nav from the list. Does not throw an exception if not exists
	 *
	 * @param key key to delete by
	 * @returns boolean indicating if a nav was deleted
	 */
	public removeVerticalNav(key: string): boolean {
		return this.verticalNavs.delete(key);
	}

	/**
	 * Adds to the horizontal nav list
	 *
	 * @param hn nav element to add
	 * @param key key to use for determining uniqueness
	 * @returns boolean indicating if a nav was replaced.
	 */
	public addHorizontalNav(hn: HeaderNavigation, key: string): boolean {
		const present = this.horizontalNavs.has(key);
		this.horizontalNavs.set(key, hn);
		this.horizontalNavsUpdated.next(this.buildHorizontalNavs());
		return present;
	}

	/**
	 * Adds to the vertical nav list
	 *
	 * @param hn nav element to add
	 * @param key key to use for determining uniqueness
	 * @returns boolean indicating if a nav was replaced.
	 */
	public addVerticalNav(hn: HeaderNavigation, key: string): boolean {
		const present = this.verticalNavs.has(key);
		this.verticalNavs.set(key, hn);
		this.verticalNavsUpdated.next(this.buildVerticalNavs());
		return present;
	}
}
