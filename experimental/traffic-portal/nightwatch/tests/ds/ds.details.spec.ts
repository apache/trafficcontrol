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

describe("DS Detail Spec", () => {
	beforeEach(() => {
		browser.page.login()
			.navigate().section.loginForm
			.loginAndWait(browser.globals.adminUser, browser.globals.adminPass);
		browser.page.deliveryServiceCard()
			.section.cards
			.viewDetails(`testDS${browser.globals.uniqueString}`);
	});

	it("Verify page test", (): void => {
		const page = browser.page.deliveryServiceDetail();
		page.assert.visible("@bandwidthChart")
			.assert.visible("@tpsChart")
			.assert.enabled("@invalidateJobs");

		page.section.dateInputForm
			.assert.enabled("@fromDate")
			.assert.enabled("@fromTime")
			.assert.enabled("@toDate")
			.assert.enabled("@toTime")
			.assert.enabled("@refreshBtn")
			.end();
	});

	it("Default values test", (): void => {
		const page = browser.page.deliveryServiceDetail();
		const now = new Date();
		const nowString = now.toISOString();
		const date = nowString.split("T")[0];
		let time = nowString.split("T")[1].substring(0, 5);
		time = `${(+time.split(":")[0] - now.getTimezoneOffset()/60).toString().padStart(2, "0")}:${time.split(":")[1]}`;

		page.section.dateInputForm
			.assert.value("@fromDate", date)
			.assert.value("@fromTime", "00:00")
			.assert.value("@toDate", date)
			.assert.value("@toTime", time)
			.end();
	});
});
