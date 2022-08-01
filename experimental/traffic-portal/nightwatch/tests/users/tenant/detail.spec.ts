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

describe("Tenant Detail Spec", () => {
	it("Root tenant", () => {
		const page = browser.page.tenantDetail();
		browser.url(`${page.api.launchUrl}/core/tenants/1`, res => {
			browser.assert.ok(res.status === 0);
			page.waitForElementVisible("mat-card")
				.assert.not.enabled("@active")
				.assert.not.enabled("@name")
				.assert.not.enabled("@parent")
				.assert.not.enabled("@saveBtn")
				.assert.value("@name", "root")
				.assert.value("@active", "on");
		});
	});

	it("Test tenant", () => {
		const page = browser.page.tenantDetail();
		browser.url(`${page.api.launchUrl}/core/tenants/2`, res => {
			browser.assert.ok(res.status === 0);
			page.waitForElementVisible("mat-card")
				.assert.enabled("@active")
				.assert.enabled("@name")
				.assert.enabled("@parent")
				.assert.enabled("@saveBtn");
		});
	});

	it("New tenant", () => {
		const page = browser.page.tenantDetail();
		browser.url(`${page.api.launchUrl}/core/tenants/new`, res => {
			browser.assert.ok(res.status === 0);
			page.waitForElementVisible("mat-card")
				.assert.enabled("@active")
				.assert.enabled("@name")
				.assert.enabled("@parent")
				.assert.enabled("@saveBtn")
				.assert.containsText("@name", "")
				.assert.value("@active", "on")
				.assert.value("@parent", "");
		});
	});
});
