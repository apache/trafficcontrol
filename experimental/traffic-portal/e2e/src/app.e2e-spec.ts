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

import { browser, until } from "protractor";

import { LoginPage } from "./app.po";

describe("workspace-project App", () => {
	let page: LoginPage;

	beforeEach(() => {
		page = new LoginPage();
		page.navigateTo();
	});

	it("should allow login", async done => {
		// TODO: put this in a config somewhere
		await Promise.all([
			page.usernameInput.sendKeys("admin"),
			page.passwordInput.sendKeys("twelve12")
		]);
		await page.loginButton.click();
		const loggedInURL = `${browser.baseUrl}core`;
		try {
			await browser.wait(until.urlIs(loggedInURL), 1000);
		} catch(e) {
			done.fail(`page did not navigate after login: ${e}`);
			return;
		}
		expect(browser.getCurrentUrl()).toBe(loggedInURL);
		done();
	});

	it("should clear the form", async () => {
		await page.usernameInput.sendKeys("foo");
		await page.passwordInput.sendKeys("bar");
		await page.clearButton.click();
		expect(page.usernameInput.getText()).toBe("");
		expect(page.passwordInput.getText()).toBe("");
	});
});
