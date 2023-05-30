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

describe("Servers Detail Spec", () => {
	it("New server loads correctly", async () => {
		await browser.page.servers.serversTable()
			.section.serversTable
			.createNew();
		await browser.assert.urlContains("core/servers/new");
		const page = browser.page.servers.serversDetail();
		await page.section.detailCard
			.assert.enabled("@hostName")
			.assert.enabled("@cdn")
			.assert.enabled("@cacheGroup")
			.assert.enabled("@physLoc")
			.assert.enabled("@status")
			.assert.not.elementPresent("@offlineReason")
			.assert.enabled("@type")
			.assert.enabled("@httpPort")
			.assert.enabled("@httpsPort")
			.assert.enabled("@rack")
			.assert.not.elementPresent("@id")
			.assert.not.elementPresent("@lastUpdated")
			.assert.not.elementPresent("@statusLastUpdated")
			.assert.enabled("@profileNames")
			.assert.enabled("@intfAddBtn")
			.assert.enabled("@iloIP")
			.assert.enabled("@iloGateway")
			.assert.enabled("@iloNetmask")
			.assert.enabled("@iloUsername")
			.assert.enabled("@iloPassword")
			.assert.enabled("@mgmtIP")
			.assert.enabled("@mgmtGateway")
			.assert.enabled("@mgmtNetmask")
			.assert.not.elementPresent("@deleteBtn")
			.assert.enabled("@submitBtn");
	});

	it("Fields are loaded correctly", async () => {
		await browser.page.servers.serversTable()
			.section.serversTable
			.openDetails(browser.globals.testData.edgeServer);
		await browser.assert.urlContains(`core/servers/${  browser.globals.testData.edgeServer.id}`);
		const page = browser.page.servers.serversDetail();
		await page.section.detailCard
			.assert.enabled("@hostName")
			.assert.valueEquals("@hostName", browser.globals.testData.edgeServer.hostName)
			.assert.enabled("@cdn")
			.assert.textEquals("@cdn", browser.globals.testData.edgeServer.cdnName)
			.assert.enabled("@cacheGroup")
			.assert.textEquals("@cacheGroup", browser.globals.testData.edgeServer.cachegroup)
			.assert.enabled("@physLoc")
			.assert.textEquals("@physLoc", browser.globals.testData.edgeServer.physLocation)
			.assert.not.enabled("@statusDisabled")
			.assert.valueEquals("@statusDisabled", browser.globals.testData.edgeServer.status)
			.assert.not.elementPresent("@offlineReason")
			.assert.enabled("@type")
			.assert.textEquals("@type", browser.globals.testData.edgeServer.type)
			.assert.enabled("@httpPort")
			.assert.enabled("@httpsPort")
			.assert.textEquals("@httpsPort", String(browser.globals.testData.edgeServer.httpsPort ?? ""))
			.assert.enabled("@rack")
			.assert.textEquals("@rack", browser.globals.testData.edgeServer.rack ?? "")
			.assert.not.enabled("@id")
			.assert.valueEquals("@id", String(browser.globals.testData.edgeServer.id))
			.assert.not.enabled("@lastUpdated")
			.assert.not.elementPresent("@statusLastUpdated")
			.assert.enabled("@profileNames")
			.assert.textEquals("@profileNames", browser.globals.testData.edgeServer.profileNames[0])
			.assert.enabled("@intfAddBtn")
			.assert.enabled("@iloIP")
			.assert.valueEquals("@iloIP", browser.globals.testData.edgeServer.iloIpAddress ?? "")
			.assert.enabled("@iloGateway")
			.assert.valueEquals("@iloGateway", browser.globals.testData.edgeServer.iloIpGateway ?? "" )
			.assert.enabled("@iloNetmask")
			.assert.valueEquals("@iloNetmask", browser.globals.testData.edgeServer.iloIpNetmask ?? "")
			.assert.enabled("@iloUsername")
			.assert.valueEquals("@iloUsername", browser.globals.testData.edgeServer.iloUsername ?? "")
			.assert.enabled("@iloPassword")
			.assert.valueEquals("@iloPassword", browser.globals.testData.edgeServer.iloPassword ?? "")
			.assert.enabled("@mgmtIP")
			.assert.valueEquals("@mgmtIP", browser.globals.testData.edgeServer.mgmtIpAddress ?? "")
			.assert.enabled("@mgmtGateway")
			.assert.valueEquals("@mgmtGateway", browser.globals.testData.edgeServer.mgmtIpGateway ?? "")
			.assert.enabled("@mgmtNetmask")
			.assert.valueEquals("@mgmtNetmask", browser.globals.testData.edgeServer.mgmtIpNetmask ?? "")
			.assert.enabled("@deleteBtn")
			.assert.enabled("@submitBtn");
	});
});
