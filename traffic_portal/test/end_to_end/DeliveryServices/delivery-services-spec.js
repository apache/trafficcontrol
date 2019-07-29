/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

var pd = require('./pageData.js');
var cfunc = require('../common/commonFunctions.js');

describe('Traffic Portal Delivery Services Suite', function() {

	const pageData = new pd();
	const commonFunctions = new cfunc();
	const mockVals = {
		dsTypes: {
			anyMap: [
				"ANY_MAP"
			],
			dns: [
				"DNS",
				"DNS_LIVE_NATNL",
				"DNS_LIVE"
			],
			http: [
				"HTTP",
				"HTTP_NO_CACHE",
				"HTTP_LIVE",
				"HTTP_LIVE_NATNL"
			],
			steering: [
				"STEERING",
				"CLIENT_STEERING"
			]
		},
		xmlIds: {
			anyMap: "any-map-xml-id-" + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz'),
			dns: "dns-xml-id-" + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz'),
			http: "http-xml-id-" + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz'),
			steering: "steering-xml-id-" + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz')
		},
		displayName: "display name",
		longDesc: "This is only a test delivery service that should be disposed of by Automated UI Testing.",
		tenantId: "- root",
		active: "Active",
		orgServerFqdn: "https://example.com",
		protocol: "HTTP to HTTPS"
	};
	const repeater = "ds in ::deliveryServices";

	it('should open delivery services page', async () => {
		console.log('Opening delivery services page');
		await browser.setLocation("delivery-services");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services");
	});

	for(const [dsTypeKey, dsTypeArray] of Object.entries(mockVals.dsTypes)) {
		const dsType = dsTypeArray[0];
		for (const specificDsType of dsTypeArray) {
			it('should click new delivery service and select ' + dsType + ' category from the dropdown', async () => {
				console.log('Clicked Create New and selecting ' + dsType);
				await pageData.createDeliveryServiceButton.click();
				expect(pageData.selectFormSubmitButton.isEnabled()).toBe(false);
				await commonFunctions.selectDropdownByLabel(pageData.selectFormDropdown, dsType);
				expect(pageData.selectFormSubmitButton.isEnabled()).toBe(true);
				await pageData.selectFormSubmitButton.click();
			});
		
			it('should populate and submit the delivery service form', async () => {
				console.log('Creating a DS for ' + mockVals.xmlIds[dsTypeKey]);
				expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services/new?type=" + dsType);
				expect(pageData.createButton.isEnabled()).toBe(false);
				await pageData.xmlId.sendKeys(mockVals.xmlIds[dsTypeKey]);
				await pageData.displayName.sendKeys(mockVals.displayName);
				// no label for this option, so commonFunctions is not used
				await pageData.active.element(by.cssContainingText('option', mockVals.active)).click();
				await commonFunctions.selectDropdownByLabel(pageData.type, specificDsType);
				await commonFunctions.selectDropdownByLabel(pageData.tenantId, mockVals.tenantId);
				await commonFunctions.selectDropdownByNum(pageData.cdn, 1);
				await pageData.orgServerFqdn.isPresent().then(async (present) => {
					if (present)
						await pageData.orgServerFqdn.sendKeys(mockVals.orgServerFqdn);
				});
				await pageData.protocol.isPresent().then(async (present) => {
					if (present)
						// no label for this option, so commonFunctions is not used
						await pageData.protocol.element(by.cssContainingText('option', mockVals.protocol)).click();
				});
				// all required fields have been set, create button should be enabled
				expect(pageData.createButton.isEnabled()).toBe(true);
				await pageData.createButton.click();
				expect(pageData.successMsg.isPresent()).toBe(true);
        		expect(element(by.cssContainingText('div', 'Delivery Service [ '+mockVals.xmlIds[dsTypeKey]+' ] created')).isPresent()).toBe(true, 'Actual message does not match expected message');
				expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services/[0-9]+.+" + dsType);
			});

			it('should toggle the visibility of the first table column ', function() {
				browser.driver.findElement(by.id('toggleColumns')).click();
				let first = element.all(by.css('input[type=checkbox]')).first();
				expect(first.isSelected()).toBe(true);
				first.click();
				expect(first.isSelected()).toBe(false);
				let tableColumns = element.all(by.css('#deliveryServicesTable tr:first-child td'));
				expect(tableColumns.count()).toBe(10);
			});
		
			it('should update the ' + specificDsType + ' delivery service', async () => {
				console.log('Updating the ' + specificDsType + ' delivery service for ' + mockVals.xmlIds[dsTypeKey]);
				await browser.setLocation("delivery-services");
				await commonFunctions.clickTableEntry(pageData.searchFilter, mockVals.xmlIds[dsTypeKey], repeater);
				expect(pageData.updateButton.isEnabled()).toBe(false);
				expect(pageData.xmlId.getAttribute('readonly')).toBe('true');
				await pageData.displayName.clear().sendKeys("Updated " + mockVals.displayName);
				expect(pageData.updateButton.isEnabled()).toBe(true);
				await pageData.updateButton.click();
				expect(pageData.successMsg.isPresent()).toBe(true);
        		expect(element(by.cssContainingText('div', 'Delivery Service [ '+mockVals.xmlIds[dsTypeKey]+' ] updated')).isPresent()).toBe(true, 'Actual message does not match expected message');
				expect(pageData.displayName.getText() === "Updated " + mockVals.displayName);
			});
		
			it('should delete the ' + specificDsType + ' delivery service', async () => {
				console.log('Deleting ' + mockVals.xmlIds[dsTypeKey]);
				await pageData.deleteButton.click();
				await pageData.confirmWithNameInput.sendKeys(mockVals.xmlIds[dsTypeKey]);
				await pageData.deletePermanentlyButton.click();
				expect(pageData.successMsg.isPresent()).toBe(true);
        		expect(element(by.cssContainingText('div', 'Delivery service [ '+mockVals.xmlIds[dsTypeKey]+' ] deleted')).isPresent()).toBe(true, 'Actual message does not match expected message');
				expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services");
			});
		}
	}

});
