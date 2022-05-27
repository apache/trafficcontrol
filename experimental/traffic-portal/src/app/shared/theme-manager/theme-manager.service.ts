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

import {DOCUMENT} from "@angular/common";
import {EventEmitter, Inject, Injectable} from "@angular/core";

/**
 * Defines a theme. If fileName is null, it is the default theme
 */
export interface Theme {
	fileName?: string;
	name: string;
}

/**
 *
 */
@Injectable({
	providedIn: "root"
})
export class ThemeManagerService {
	private readonly storageKey = "current-theme-name";
	private readonly linkClass = "themer";

	public themeChanged = new EventEmitter<Theme>();

	constructor(@Inject(DOCUMENT) private readonly document: Document) {
		this.initTheme();
	}

	/**
	 * Initialize the theme service
	 */
	public initTheme(): void {
		const themeName = this.loadStoredTheme();
		if(themeName) {
			this.loadTheme(themeName);
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
		if(theme.fileName === undefined) {
			this.clearTheme();
			return;
		}
		this.getThemeLinkElement().setAttribute("href", theme.fileName);
		this.storeTheme(theme);
		this.themeChanged.emit(theme);
	}

	/**
	 * Revert to the default theme
	 */
	public clearTheme(): void {
		const linkEl = this.getExistingThemeLinkElement();
		if(linkEl) {
			this.document.head.removeChild(linkEl);
			this.clearStoredTheme();
			this.themeChanged.emit(this.themes[0]);
		}
	}

	/**
	 * Stores theme in localStorage
	 *
	 * @param theme Theme to be stored
	 */
	private storeTheme(theme: Theme): void {
		try {
			window.localStorage.setItem(this.storageKey, JSON.stringify(theme));
		} catch (e) {
			console.error(`Unable to store theme into local storage: ${e}`);
		}
	}

	/**
	 * Retrieves theme saved in localStorage
	 *
	 * @returns The stored theme name or null
	 */
	private loadStoredTheme(): Theme | null {
		try {
			return JSON.parse(window.localStorage.getItem(this.storageKey) ?? "");
		} catch (e) {
			console.error(`Unable to load theme from local storage: ${e}`);
		}
		return null;
	}

	/**
	 * Clears theme saved in local storage
	 */
	private clearStoredTheme(): void {
		window.localStorage.removeItem(this.storageKey);
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
