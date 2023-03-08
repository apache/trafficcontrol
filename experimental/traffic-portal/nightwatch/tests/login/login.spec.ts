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

describe("Login Spec", () => {

	it("Clear form test", () => {
		browser.page.login()
			.navigate().section.loginForm
			.fillOut("test", "asdf")
			.click("@clearBtn")
			.assert.textContains("@usernameTxt", "")
			.assert.textContains("@passwordTxt", "");
	});
	it("Incorrect password test", () => {
		browser.page.login()
			.navigate().section.loginForm
			.login("test", "asdf")
			.assert.valueEquals("@usernameTxt", "test")
			.assert.valueEquals("@passwordTxt", "asdf");
		browser.page.common()
			.assert.textContains("@snackbarEle", "Invalid");
	});
	it("Login test", () => {
		browser.page.login()
			.navigate()
			.section.loginForm
			.loginAndWait(browser.globals.adminUser, browser.globals.adminPass);
	});
});
