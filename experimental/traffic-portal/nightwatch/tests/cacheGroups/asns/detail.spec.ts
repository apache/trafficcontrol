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

describe("Asn Detail Spec", () => {
	it("Test Asn", () => {
		const page = browser.page.cacheGroups.asnDetail();
		console.log(browser.globals.testData.asn.id)
		browser.url(`${page.api.launchUrl}/core/asns/${browser.globals.testData.asn.id}`, res => {
			browser.assert.ok(res.status === 0);
			page.waitForElementVisible("mat-card")
				.assert.enabled("@asn")
				.assert.enabled("@cachegroup")
				.assert.enabled("@saveBtn")
				.assert.not.enabled("@id")
				.assert.not.enabled("@lastUpdated")
				.assert.valueEquals("@asn", String(browser.globals.testData.asn.asn))
				.assert.valueEquals("@id", String(browser.globals.testData.asn.id));
		});
	});

	it("New asn", () => {
		const page = browser.page.cacheGroups.asnDetail();
		browser.url(`${page.api.launchUrl}/core/asns/new`, res => {
			browser.assert.ok(res.status === 0);
			page.waitForElementVisible("mat-card")
				.assert.enabled("@asn")
				.assert.enabled("@cachegroup")
				.assert.enabled("@saveBtn")
				.assert.not.elementPresent("@id")
				.assert.not.elementPresent("@lastUpdated")
				.assert.valueEquals("@asn", "");
		});
	});
});
