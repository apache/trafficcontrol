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
	};

	it('should go to the Servers page', function() {
		console.log('Looading Configure/Servers');
		browser.setLocation("servers");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/servers");
	});

	it('should verify CSV link exists ', function() {
		console.log("Verify CSV button exists");
		expect(element(by.css('.dt-button.buttons-csv')).isPresent()).toBe(true);
	});

	it('should open new Servers form page', function() {
		console.log('Clicking on Create new server ' + mockVals.hostName);
		browser.driver.findElement(by.name('createServersButton')).click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/servers/new");
	});

	it('should fill out form, create button is enabled and submit', function () {
		console.log('Filling out Server form');
		expect(pageData.createButton.isEnabled()).toBe(false);
		pageData.status.click();
		pageData.status.sendKeys(mockVals.status);
		pageData.hostName.sendKeys(mockVals.hostName);
		pageData.domainName.sendKeys(mockVals.domainName);
		commonFunctions.selectDropdownbyNum(pageData.cdn, 1);
		commonFunctions.selectDropdownbyNum(pageData.cachegroup, 1);
		commonFunctions.selectDropdownbyNum(pageData.type, 1);
		commonFunctions.selectDropdownbyNum(pageData.profile, 1);
		pageData.interfaceName.sendKeys(mockVals.interfaceName);
		pageData.ipAddress.sendKeys(mockVals.ipAddress);
		pageData.ipNetmask.sendKeys(mockVals.ipNetmask);
		pageData.ipGateway.sendKeys(mockVals.ipGateway);
		pageData.interfaceMtu.sendKeys(mockVals.interfaceMtu);
		commonFunctions.selectDropdownbyNum(pageData.physLocation, 1);
		expect(pageData.createButton.isEnabled()).toBe(true);
		pageData.createButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/servers");
	});

	it('should toggle the visibility of the first table column ', function() {
		browser.driver.findElement(by.id('toggleColumns')).click();
		let first = element.all(by.css('input[type=checkbox]')).first();
		expect(first.isSelected()).toBe(true);
		first.click();
		expect(first.isSelected()).toBe(false);
		let tableColumns = element.all(by.css('#serversTable tr:first-child td'));
		expect(tableColumns.count()).toBe(11);
	});

	it('should verify the new Server and then update Server', function() {
		console.log('Verifying new server added and updating ' + mockVals.hostName);
		browser.sleep(1000);
		pageData.searchFilter.sendKeys(mockVals.hostName);
		browser.sleep(250);
		element.all(by.repeater('s in ::servers')).filter(function(row){
			return row.element(by.name('hostName')).getText().then(function(val){
				return val === mockVals.hostName;
			});
		}).get(0).click();
		browser.sleep(1000);
		pageData.domainName.clear();
		pageData.domainName.sendKeys('testupdated.com');
		pageData.type.click();
		pageData.updateButton.click();
		expect(pageData.domainName.getText() === 'testupdated.com');
	});

	it('should add a server capability to the server', function() {
		console.log('Adding new server capability to ' + mockVals.hostName);
		pageData.moreBtn.click();
		pageData.viewCapabilitiesMenuItem.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl)+"#!/servers/[0-9]+/capabilities");
		pageData.addCapabilityBtn.click();
		expect(pageData.submitButton.isEnabled()).toBe(false);
		commonFunctions.selectDropdownbyNum(pageData.selectFormDropdown, 1);
		expect(pageData.submitButton.isEnabled()).toBe(true);
		pageData.submitButton.click();
		element.all(by.css('tbody tr')).then(function(totalRows) {
			expect(totalRows.length).toBe(1);
		});
	});

	it('should navigate back to the new server and delete it', function() {
		console.log('Deleting the server ' + mockVals.hostName);
		browser.navigate().back();
		pageData.deleteButton.click();
		pageData.confirmWithNameInput.sendKeys(mockVals.hostName);
		pageData.deletePermanentlyButton.click();
	});
});
