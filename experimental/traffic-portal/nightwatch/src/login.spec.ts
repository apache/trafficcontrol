import {NightwatchBrowser} from "nightwatch";

import {type LoginPageObject} from "../page_objects/login";

module.exports = {
	"Incorrect password test":  (browser: NightwatchBrowser): void => {
		const page: LoginPageObject = browser.page.login();
		page.navigate()
			.section.loginForm
			.login("test", "asdf")
			.assert.value("@usernameTxt", "test")
			.assert.value("@passwordTxt", "asdf")
			.parent
			.assert.containsText("@snackbarEle", "Invalid")
			.end();
	},
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
	"Login test": (browser: NightwatchBrowser): void => {
		const page: LoginPageObject = browser.page.login();
		page.navigate()
			.section.loginForm
			.login("admin", "twelve12")
			.parent
			.assert.containsText("@snackbarEle", "Success")
			.end();
	}
};
