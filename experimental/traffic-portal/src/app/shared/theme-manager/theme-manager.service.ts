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

import { DOCUMENT, isPlatformServer } from "@angular/common";
import { EventEmitter, Inject, Injectable, PLATFORM_ID } from "@angular/core";

import { LoggingService } from "../logging.service";

/**
 * Defines a theme. If fileName is null, it is the default theme
 */
export interface Theme {
	fileName?: string;
	name: string;
}

/**
 * The ThemeManagerService manages the user's theming settings, to be applied
 * throughout the UI.
 */
@Injectable({
	providedIn: "root"
})
export class ThemeManagerService {
	private readonly storageKey = "current-theme-name";
	private readonly linkClass = "themer";
	private readonly isServer: boolean = false;

	public themeChanged = new EventEmitter<Theme>();

	/**
	 * Provides a "safe" accessor for the local session storage. According to
	 * typings, `Document.defaultView` may be `null`, but if it isn't then
	 * `Document.defaultView.localStorage` definitely *isn't* `null`. That's
	 * simply untrue. So this provides that check for you.
	 */
	private get localStorage(): Storage | null {
		if (this.document.defaultView && this.document.defaultView.localStorage) {
			return this.document.defaultView.localStorage;
		}
		return null;
	}

	constructor(@Inject(DOCUMENT) private readonly document: Document, private readonly log: LoggingService,
		@Inject(PLATFORM_ID) private readonly platformId: object) {
		this.isServer = isPlatformServer(this.platformId);
		this.initTheme();
	}

	/**
	 * Initialize the theme service
	 */
	public initTheme(): void {
		if (this.isServer) {
			return;
		}
		const themeName = this.loadStoredTheme();
		if (themeName) {
			this.loadTheme(themeName);
			return;
		}
		// If there is no stored theme and the user has a set preference for dark theme
		// load the dark theme.
		if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
			this.loadTheme(this.themes[1]);
		}
	}

	public readonly themes: Array<Theme> = [{
		name: "Default"
	},
	{
		fileName: "dark-default-theme.css",
		name: "Dark"
	}
	];

	/**
	 * Given a themes bundle name, load the theme and cache the value
	 *
	 * @param theme Theme to load
	 */
	public loadTheme(theme: Theme): void {
		if (!this.isServer) {
			if (theme.fileName === undefined) {
				this.clearTheme();
				return;
			}
			this.getThemeLinkElement().setAttribute("href", theme.fileName);
			this.storeTheme(theme);
			this.themeChanged.emit(theme);
		}
	}

	/**
	 * Revert to the default theme
	 */
	public clearTheme(): void {
		if (!this.isServer) {
			const linkEl = this.getExistingThemeLinkElement();
			if (linkEl) {
				this.document.head.removeChild(linkEl);
				this.clearStoredTheme();
				this.themeChanged.emit(this.themes[0]);
			}
		}
	}

	/**
	 * Stores theme in localStorage
	 *
	 * @param theme Theme to be stored
	 */
	private storeTheme(theme: Theme): void {
		try {
			this.localStorage?.setItem(this.storageKey, JSON.stringify(theme));
		} catch (e) {
			this.log.error(`Unable to store theme into local storage: ${e}`);
		}
	}

	/**
	 * Retrieves theme saved in localStorage
	 *
	 * @returns The stored theme name or null
	 */
	private loadStoredTheme(): Theme | null {
		if (!this.isServer) {
			try {
				return JSON.parse(this.localStorage?.getItem(this.storageKey) ?? "null");
			} catch (e) {
				this.log.error(`Unable to load theme from local storage: ${e}`);
			}
		}
		return null;
	}

	/**
	 * Clears theme saved in local storage
	 */
	private clearStoredTheme(): void {
		this.localStorage?.removeItem(this.storageKey);
	}

	/**
	 * Gets or creates the link element used for loading non-default themes
	 *
	 * @private
	 * @returns The html element
	 */
	private getThemeLinkElement(): Element {
		return this.getExistingThemeLinkElement() || this.createThemeLinkElement();
	}

	/**
	 * Creates the link element used for loading themes
	 *
	 * @returns The html element
	 */
	private createThemeLinkElement(): Element {
		const linkEl = this.document.createElement("link");
		linkEl.setAttribute("rel", "stylesheet");
		linkEl.classList.add(this.linkClass);
		this.document.head.appendChild(linkEl);
		return linkEl;
	}

	/**
	 * Gets the link element used for loading themes
	 *
	 * @returns The html element or null
	 */
	private getExistingThemeLinkElement(): Element | null {
		return this.document.head.querySelector(`link[rel="stylesheet"].${this.linkClass}`);
	}
}
