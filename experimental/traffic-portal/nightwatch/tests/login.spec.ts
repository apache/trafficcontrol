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
import {type NightwatchBrowser} from "nightwatch";

import {type GlobalConfig} from "../globals/globals";
import {type LoginPageObject} from "../page_objects/login";

module.exports = {
	"Clear form test": (browser: NightwatchBrowser): void => {
		const page: LoginPageObject = browser.page.login();
		page.navigate()
			.section.loginForm
			.fillOut("test", "asdf")
			.click("@clearBtn")
			.assert.containsText("@usernameTxt", "")
			.assert.containsText("@passwordTxt", "")
			.end();
	},
	"Incorrect password test":  (browser: NightwatchBrowser): void => {
		const page: LoginPageObject = browser.page.login();
		page.navigate()
			.section.loginForm
			.login("test", "asdf")
			.assert.value("@usernameTxt", "test")
			.assert.value("@passwordTxt", "asdf");
		page
			.assert.containsText("@snackbarEle", "Invalid")
			.end();
	},
	"Login test": (browser: NightwatchBrowser): void => {
		const page: LoginPageObject = browser.page.login();
		const globals = browser.globals as GlobalConfig;
		page.navigate()
			.section.loginForm
			.login(globals.adminUser, globals.adminPass)
			.parent
			.assert.containsText("@snackbarEle", "Success")
			.end();
	}
};
