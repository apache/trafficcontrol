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

describe("DS Invalidation Jobs Spec", () => {
	beforeEach(() => {
		browser.page.deliveryServices.deliveryServiceCard()
			.navigate()
			.section.cards
			.viewDetails(`testDS${browser.globals.uniqueString}`);
		browser.page.deliveryServices.deliveryServiceDetail()
			.click("@invalidateJobs")
			.assert.urlContains("invalidation-jobs");
	});

	it("Verify page", () => {
		browser.page.deliveryServices.deliveryServiceInvalidationJobs()
			.assert.enabled("@addButton");
	});

	it("Manage Job", async () => {
		const page = browser.page.deliveryServices.deliveryServiceInvalidationJobs();
		const common = browser.page.common();
		page
			.click("@addButton");
		const startDate = new Date();
		startDate.setDate(startDate.getDate() + 1);
		browser.waitForElementVisible("tp-new-invalidation-job-dialog")
			.assert.valueEquals("input[name='startDate']", startDate.toLocaleDateString())
			.setValue("input[name='regexp']", "/invalidateMe")
			.click("button#submit");
		common
			.assert.textContains("@snackbarEle", "created")
			.click("simple-snack-bar button");
		page.assert.visible({index: 0, selector: "div.invalidation-job"})
			.assert.enabled({index: 0, selector: "div.invalidation-job button"})
			.assert.enabled({index: 1, selector: "div.invalidation-job button"});
		page
			.click({index: 0, selector: "div.invalidation-job button"});
		browser.waitForElementVisible("tp-new-invalidation-job-dialog")
			.assert.valueEquals("input[name='startDate']", startDate.toLocaleDateString())
			.assert.valueEquals("input[name='regexp']", "invalidateMe")
			.setValue("input[name='regexp']", "/invalidateMe2")
			.click("button#submit");
		common
			.assert.textContains("@snackbarEle", "created")
			.click("simple-snack-bar button");
		page
			.click({index: 1, selector: "div.invalidation-job button"});
		common
			.assert.textContains("@snackbarEle", "was deleted");
	});
});
