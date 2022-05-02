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
import type { TestSuite } from "../globals";
import type { LoginPageObject } from "../page_objects/login";
import type { ServersPageObject } from "../page_objects/servers";

const suite: TestSuite = {
	"Filter by hostname": async browser => {
		const username = browser.globals.adminUser;
		const password = browser.globals.adminPass;

		const loginPage: LoginPageObject = browser.page.login();
		loginPage.navigate().section.loginForm.login(username, password);

		const page: ServersPageObject = browser.waitForElementPresent("main").page.servers().navigate();
		page.pause(4000);
		let tbl = page.waitForElementPresent("input[name=fuzzControl]").section.serversTable;
		tbl = tbl.searchText("edge");
		tbl.parent.assert.urlContains("search=edge").end();
	}
};

export default suite;
