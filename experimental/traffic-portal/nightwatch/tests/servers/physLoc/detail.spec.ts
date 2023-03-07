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

describe("Phys Location Detail Spec", () => {
	it("Test physLoc", () => {
		const page = browser.page.servers.physLocDetail();
		browser.url(`${page.api.launchUrl}/core/phys-locs/${browser.globals.testData.physLoc.id}`, res => {
			browser.assert.ok(res.status === 0);
			page.waitForElementVisible("form")
				.assert.enabled("@name")
				.assert.enabled("@saveBtn")
				.assert.not.enabled("@id")
				.assert.not.enabled("@lastUpdated")
				.assert.valueEquals("@name", browser.globals.testData.physLoc.name)
				.assert.valueEquals("@id", String(browser.globals.testData.physLoc.id))
				.assert.valueEquals("@address", browser.globals.testData.physLoc.address)
				.assert.valueEquals("@city", browser.globals.testData.physLoc.city)
				.assert.valueEquals("@comments", browser.globals.testData.physLoc.comments ?? "")
				.assert.valueEquals("@email", browser.globals.testData.physLoc.email ?? "")
				.assert.valueEquals("@poc", browser.globals.testData.physLoc.poc ?? "")
				.assert.valueEquals("@shortName", browser.globals.testData.physLoc.shortName)
				.assert.valueEquals("@state", browser.globals.testData.physLoc.state)
				.assert.valueEquals("@zip", browser.globals.testData.physLoc.zip)
				.assert.valueEquals("@phone", browser.globals.testData.physLoc.phone ?? "")
				.assert.textEquals("@region", String(browser.globals.testData.physLoc.region));
		});
	});

	it("New physLoc", () => {
		const page = browser.page.servers.physLocDetail();
		browser.url(`${page.api.launchUrl}/core/phys-locs/new`, res => {
			browser.assert.ok(res.status === 0);
			page.waitForElementVisible("mat-card")
				.assert.enabled("@name")
				.assert.enabled("@saveBtn")
				.assert.not.elementPresent("@id")
				.assert.not.elementPresent("@lastUpdated")
				.assert.valueEquals("@name", "");
		});
	});
});
