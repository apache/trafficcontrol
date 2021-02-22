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

import { browser, by, element } from "protractor";

/**
 * AppPage is the main page of the app, for end-to-end testing purposes.
 */
export class AppPage {
	/**
	 * Navigates to the base URL.
	 */
	public async navigateTo(): Promise<unknown> {
		return browser.get(browser.baseUrl) as Promise<unknown>;
	}

	/**
	 * Gets the text from the title element.
	 */
	public async getTitleText(): Promise<string> {
		return element(by.css("app-root .content span")).getText() as Promise<string>;
	}
}
