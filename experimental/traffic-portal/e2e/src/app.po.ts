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

import { browser, by, element, ElementFinder } from "protractor";

/**
 * LoginPage is the main page of the app, for end-to-end testing purposes.
 */
export class LoginPage {
	/** The base URL for this page (i.e. without any fragment or query string). */
	private readonly baseURL = `${browser.baseUrl}login`;

	/**
	 * Navigates to the base URL.
	 */
	public async navigateTo(): Promise<unknown> {
		return browser.get(this.baseURL);
	}

	/**
	 * Gets the text from the title element.
	 */
	public async getTitleText(): Promise<string> {
		return element(by.css("app-root .content span")).getText() as Promise<string>;
	}

	/**
	 * The input text box for the username.
	 */
	public get usernameInput(): ElementFinder {
		return element(by.id("u"));
	}

	/**
	 * The input text box for the password.
	 */
	public get passwordInput(): ElementFinder {
		return element(by.id("p"));
	}

	/**
	 * The "Login" button.
	 */
	public get loginButton(): ElementFinder {
		return element(by.partialButtonText("Login"));
	}

	/**
	 * The "Clear" button.
	 */
	public get clearButton(): ElementFinder {
		return element(by.partialButtonText("Clear"));
	}

	/**
	 * Uses the Login form to authenticate a user - or at least attempt to.
	 *
	 * @param username The username of the user as whom to authenticate.
	 * @param password The user's password.
	 * @returns whether or not the login succeeded.
	 */
	public async login(username: string, password: string): Promise<boolean> {
		this.usernameInput.sendKeys(username);
		this.passwordInput.sendKeys(password);
		return this.loginButton.click().then(
			async () => (await browser.getCurrentUrl()) !== this.baseURL
		);
	}
}
