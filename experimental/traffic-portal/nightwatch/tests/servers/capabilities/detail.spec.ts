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

describe("Capability Details Spec", () => {
	it("Test capability", () => {
		const page = browser.page.servers.capabilities.capabilityDetail();
		browser.url(`${page.api.launchUrl}/core/capabilities/${browser.globals.testData.capability.name}`, res => {
			browser.assert.ok(res.status === 0);
			page.waitForElementVisible("form")
				.assert.enabled("@name")
				.assert.enabled("@saveBtn")
				.assert.not.enabled("@lastUpdated")
				.assert.valueEquals("@name", browser.globals.testData.capability.name);
		});
	});

	it("New capability", () => {
		const page = browser.page.servers.capabilities.capabilityDetail();
		browser.url(`${page.api.launchUrl}/core/new-capability`, res => {
			browser.assert.ok(res.status === 0);
			page.waitForElementVisible("mat-card")
				.assert.enabled("@name")
				.assert.enabled("@saveBtn")
				.assert.not.elementPresent("@lastUpdated")
				.assert.valueEquals("@name", "");
		});
	});
});
