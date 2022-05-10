/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import {AugmentedBrowser} from "nightwatch/globals/globals";

describe("DS Invalidation Jobs Spec", () => {
	let augBrowser: AugmentedBrowser;
	before(() => {
		augBrowser = browser as AugmentedBrowser;
	});

	beforeEach(() => {
		augBrowser.page.login()
			.navigate().section.loginForm
			.loginAndWait(augBrowser.globals.adminUser, augBrowser.globals.adminPass);
		augBrowser.page.deliveryServiceCard()
			.section.cards
			.viewDetails(`testDS${augBrowser.globals.uniqueString}`);
		augBrowser.page.deliveryServiceDetail()
			.click("@invalidateJobs")
			.assert.urlContains("invalidation-jobs");
	});

	it("Verify page", () => {
		augBrowser.page.deliveryServiceInvalidationJobs()
			.assert.enabled("@addButton")
			.end();
	});

	it("Manage Job", async () => {
		const page = augBrowser.page.deliveryServiceInvalidationJobs();
		const common = augBrowser.page.common();
		page
			.click("@addButton");
		const startDate = new Date();
		startDate.setDate(startDate.getDate() + 1);
		augBrowser.waitForElementVisible("tp-new-invalidation-job-dialog")
			.assert.value("input[name='startDate']", startDate.toLocaleDateString())
			.setValue("input[name='regexp']", "/invalidateMe")
			.click("button#submit");
		common
			.assert.containsText("@snackbarEle", "created")
			.click("simple-snack-bar button");
		page.assert.visible({index: 0, selector: "li.invalidation-job"})
			.assert.enabled({index: 0, selector: "li.invalidation-job button"})
			.assert.enabled({index: 1, selector: "li.invalidation-job button"});
		page
			.click({index: 0, selector: "li.invalidation-job button"});
		augBrowser.waitForElementVisible("tp-new-invalidation-job-dialog")
			.assert.value("input[name='startDate']", startDate.toLocaleDateString())
			.assert.value("input[name='regexp']", "invalidateMe")
			.setValue("input[name='regexp']", "/invalidateMe2")
			.click("button#submit");
		common
			.assert.containsText("@snackbarEle", "created")
			.click("simple-snack-bar button");
		page
			.click({index: 1, selector: "li.invalidation-job button"});
		common
			.assert.containsText("@snackbarEle", "was deleted")
			.end();
	});
});
