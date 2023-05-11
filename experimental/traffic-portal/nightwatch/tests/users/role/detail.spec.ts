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

describe("Role Detail Spec", () => {
	it("Test Role", () => {
		const page = browser.page.users.roleDetail();
		browser.url(`${page.api.launchUrl}/core/roles/${browser.globals.testData.role.name}`, res => {
			browser.assert.ok(res.status === 0);
			page.waitForElementVisible("mat-card")
				.assert.enabled("@description")
				.assert.enabled("@name")
				.assert.attributeEquals("@permissions", "readonly", "true")
				.assert.enabled("@saveBtn")
				.assert.valueEquals("@name", String(browser.globals.testData.role.name))
				.assert.valueEquals("@permissions", String(browser.globals.testData.role.permissions));
		});
	});

	it("New asn", () => {
		const page = browser.page.users.roleDetail();
		browser.url(`${page.api.launchUrl}/core/roles/new`, res => {
			browser.assert.ok(res.status === 0);
			page.waitForElementVisible("mat-card")
				.assert.enabled("@description")
				.assert.enabled("@name")
				.assert.not.hasAttribute("@permissions", "readonly")
				.assert.enabled("@saveBtn")
				.assert.valueEquals("@name", "");
		});
	});
});
