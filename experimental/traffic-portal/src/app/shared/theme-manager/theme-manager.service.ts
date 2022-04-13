import {EventEmitter, Injectable} from "@angular/core";

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
	private static readonly STORAGE_KEY = "current-theme-name";
	private static readonly LINK_KEY = "themer";

	public themeChanged = new EventEmitter<Theme>();

	/**
	 * Initialize the theme service
	 */
	public initTheme(): void {
		const themeName = ThemeManagerService.loadStoredTheme();
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
		ThemeManagerService.getThemeLinkElement().setAttribute("href", theme.fileName);
		ThemeManagerService.storeTheme(theme);
		this.themeChanged.emit(theme);
	}

	/**
	 * Revert to the default theme
	 */
	public clearTheme(): void {
		const linkEl = ThemeManagerService.getExistingThemeLinkElement();
		if(linkEl) {
			document.head.removeChild(linkEl);
			ThemeManagerService.clearStoredTheme();
			this.themeChanged.emit(this.themes[0]);
		}
	}

	/**
	 * Stores theme in localStorage
	 *
	 * @param theme Theme to be stored
	 * @private
	 */
	private static storeTheme(theme: Theme): void {
		try {
			window.localStorage.setItem(this.STORAGE_KEY, JSON.stringify(theme));
		} catch (e) {
			console.error(`Unable to store theme into local storage: ${e}`);
		}
	}

	/**
	 * Retrieves theme saved in localStorage
	 *
	 * @private
	 * @returns The stored theme name or null
	 */
	private static loadStoredTheme(): Theme | null {
		try {
			return JSON.parse(window.localStorage.getItem(this.STORAGE_KEY) ?? "");
		} catch (e) {
			console.error(`Unable to load theme from local storage: ${e}`);
		}
		return null;
	}

	/**
	 * Clears theme saved in local storage
	 *
	 * @private
	 */
	private static clearStoredTheme(): void {
		window.localStorage.removeItem(this.STORAGE_KEY);
	}

	/**
	 * Gets or creates the link element used for loading non-default themes
	 *
	 * @private
	 * @returns The html element
	 */
	private static getThemeLinkElement(): Element {
		return this.getExistingThemeLinkElement() || this.createThemeLinkElement();
	}

	/**
	 * Creates the link element used for loading themes
	 *
	 * @private
	 * @returns The html element
	 */
	private static createThemeLinkElement(): Element {
		const linkEl = document.createElement("link");
		linkEl.setAttribute("rel", "stylesheet");
		linkEl.classList.add(this.LINK_KEY);
		document.head.appendChild(linkEl);
		return linkEl;
	}

	/**
	 * Gets the link element used for loading themes
	 *
	 * @private
	 * @returns The html element or null
	 */
	private static getExistingThemeLinkElement(): Element | null {
		return document.head.querySelector(`link[rel="stylesheet"].${this.LINK_KEY}`);
	}
}
