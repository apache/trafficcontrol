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

describe('Traffic Portal Servers Test Suite', function() {

	const pageData = new pd();
	const commonFunctions = new cfunc();
	const mockVals = {
		status: "OFFLINE",
		hostName: "testHost-" + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz0123456789'),
		domainName: "servertest.com",
		interfaceName: "testInterfaceName",
		ipAddress: "10.42.80.118",
		ipNetmask: "255.255.255.252",
		ipGateway: "10.42.80.117",
		interfaceMtu: "9000",
		types: ['EDGE', 'MID']
	};
	const repeater = 's in ::servers';

	it('should go to the Servers page', async () => {
		console.log('Looading Configure/Servers');
		await browser.setLocation("servers");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/servers");
	});

	it('should open new Servers form page', async () => {
		console.log('Clicking on Create new server ' + mockVals.hostName);
		await pageData.createServerButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/servers/new");
	});

	it('should fill out form, create button is enabled and submit', async () => {
		console.log('Filling out Server form');
		expect(pageData.createButton.isEnabled()).toBe(false);
		await commonFunctions.selectDropdownByLabel(pageData.status, mockVals.status);
		await pageData.hostName.sendKeys(mockVals.hostName);
		await pageData.domainName.sendKeys(mockVals.domainName);
		await commonFunctions.selectDropdownByNum(pageData.cdn, 1);
		await commonFunctions.selectDropdownByNum(pageData.cachegroup, 1);
		await commonFunctions.selectDropdownByLabel(pageData.type, mockVals.types[0]);
		await commonFunctions.selectDropdownByNum(pageData.profile, 1);
		await pageData.interfaceName.sendKeys(mockVals.interfaceName);
		await pageData.ipAddress.sendKeys(mockVals.ipAddress);
		await pageData.ipNetmask.sendKeys(mockVals.ipNetmask);
		await pageData.ipGateway.sendKeys(mockVals.ipGateway);
		await pageData.interfaceMtu.sendKeys(mockVals.interfaceMtu);
		await commonFunctions.selectDropdownByNum(pageData.physLocation, 1);
		expect(pageData.createButton.isEnabled()).toBe(true);
		await pageData.createButton.click();
		expect(pageData.successMsg.isPresent()).toBe(true);
        expect(pageData.serverCreatedText.isPresent()).toBe(true, 'Actual message does not match expected message');
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/servers");
	});

	it('should toggle the visibility of the first table column ', async () => {
		await element(by.id('toggleColumns')).click();
		let first = element.all(by.css('input[type=checkbox]')).first();
		expect(first.isSelected()).toBe(true);
		await first.click();
		expect(first.isSelected()).toBe(false);
		let tableColumns = element.all(by.css('#serversTable tr:first-child td'));
		expect(tableColumns.count()).toBe(11);
	});

	it('should verify the new Server and then update Server', async () => {
		console.log('Verifying new server added and updating ' + mockVals.hostName);
		await pageData.searchFilter.sendKeys(mockVals.hostName);
		await commonFunctions.clickTableEntry(pageData.searchFilter, mockVals.hostName, repeater);
		await pageData.domainName.clear().sendKeys('updated.' + mockVals.domainName);
		await pageData.type.click().sendKeys(mockVals.types[1]);
		await pageData.updateButton.click();
		expect(pageData.domainName.getText() === 'updated.' + mockVals.domainName);
		expect(pageData.type.getText() === mockVals.types[1]);
		expect(pageData.successMsg.isPresent()).toBe(true);
        expect(pageData.serverUpdatedText.isPresent()).toBe(true, 'Actual message does not match expected message');
	});

	it('should delete the new Server', async () => {
		console.log('Deleting the server ' + mockVals.hostName);
		await pageData.deleteButton.click();
		await pageData.confirmWithNameInput.sendKeys(mockVals.hostName);
		await pageData.deletePermanentlyButton.click();
		expect(pageData.successMsg.isPresent()).toBe(true);
        expect(pageData.serverDeletedText.isPresent()).toBe(true, 'Actual message does not match expected message');
	});
});
