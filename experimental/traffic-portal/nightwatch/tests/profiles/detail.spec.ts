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

describe("Profile Detail Spec", () => {
	it("Test Profile", () => {
		const page = browser.page.profiles.profileDetail();
		browser.url(`${page.api.launchUrl}/core/profiles/${browser.globals.testData.profile.id}`, res => {
			browser.assert.ok(res.status === 0);
			page.waitForElementVisible("mat-card")
				.assert.enabled("@name")
				.assert.enabled("@cdn")
				.assert.enabled("@type")
				.assert.enabled("@routingDisabled")
				.assert.enabled("@description")
				.assert.not.enabled("@id")
				.assert.not.enabled("@lastUpdated")
				.assert.enabled("@saveBtn")
				.assert.valueEquals("@name", browser.globals.testData.profile.name)
				.assert.valueEquals("@id", String(browser.globals.testData.profile.id));
		});
	});

	it("New Profile", () => {
		const page = browser.page.profiles.profileDetail();
		browser.url(`${page.api.launchUrl}/core/types/new`, res => {
			browser.assert.ok(res.status === 0);
			page.waitForElementVisible("mat-card")
				.assert.enabled("@cdn")
				.assert.enabled("@description")
				.assert.enabled("@name")
				.assert.enabled("@routingDisabled")
				.assert.enabled("@saveBtn")
				.assert.enabled("@type")
				.assert.not.elementPresent("@id")
				.assert.not.elementPresent("@lastUpdated")
				.assert.valueEquals("@name", "");
		});
	});
});
