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
import type { GlobalConfig } from "nightwatch/globals/globals";
import { LoginPageObject } from "nightwatch/page_objects/login";
import type { UsersPageObject } from "nightwatch/page_objects/users";

describe("Users Spec", () => {
	beforeEach(() => {
		const globals = browser.globals as GlobalConfig;
		(browser.page.login() as LoginPageObject)
			.navigate().section.loginForm
			.loginAndWait(globals.adminUser, globals.adminPass);
	});

	it("Filter by username", async () => {
		const username = (browser.globals as GlobalConfig).adminUser;

		const page: UsersPageObject = browser.waitForElementPresent("main").page.users();
		page.navigate()
			.waitForElementPresent(".ag-row");
		let tbl = page.section.usersTable;
		if (! await tbl.getColumnState("Username")) {
			tbl = tbl.toggleColumn("Username");
		}

		tbl = tbl.searchText(username);
		page.assert.urlContains(`search=${username}`);

		tbl.api.elements("css selector", ".ag-row:not(.ag-hidden .ag-row)",
			result => {
				if (result.status === 1) {
					browser.assert.equal(true, false, `failed to select ag-grid rows: ${result.value.message}`);
					return;
				}
				browser.assert.equal(result.value.length, 1)
					.end();
			}
		);
	});
	// Uncomment when user details page exists
	// "View user details":  browser => {
	// 	const username = browser.globals.adminUser;
	// 	const password = browser.globals.adminPass;

	// 	const loginPage: LoginPageObject = browser.page.login();
	// 	loginPage.navigate().section.loginForm.login(username, password);

	// 	const page: UsersPageObject = browser.waitForElementPresent("main").page.users();
	// 	const tbl = page.navigate().waitForElementPresent(".ag-row").section.usersTable.searchText(username);

	// 	const userRow = tbl.parent.api.moveToElement(".ag-row:not(.ag-hidden .ag-row)", 2, 2, 100, {pointer: 0, viewport: 0});
	// 	userRow.mouseButtonClick("right").click("button[name=View-User-Details]").assert.urlContains(username);
	// }
});
