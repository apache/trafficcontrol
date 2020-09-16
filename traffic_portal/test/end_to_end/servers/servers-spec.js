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
		interfaceMtu: "9000",
	};

	it('should go to the Servers page', function() {
		console.log('Loading Configure/Servers');
		browser.setLocation("servers");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/servers");
	});

	it('should open new Servers form page', function() {
		console.log('Clicking on Create new server ' + mockVals.hostName);
		pageData.moreBtn.click();
		pageData.createServerMenuItem.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/servers/new");
	});

	it('should fill out form, create button is enabled and submit', function () {
		console.log('Filling out Server form');
		expect(pageData.createButton.isEnabled()).toBe(false);
		pageData.status.click();
		pageData.status.sendKeys(mockVals.status);
		pageData.hostName.sendKeys(mockVals.hostName);
		pageData.domainName.sendKeys(mockVals.domainName);
		commonFunctions.selectDropdownbyNum(pageData.cdn, 2); // the ALL CDN is first so let's pick a real CDN
		commonFunctions.selectDropdownbyNum(pageData.cachegroup, 1);
		element(by.css("#type [label='EDGE']")).click();
		commonFunctions.selectDropdownbyNum(pageData.profile, 1);
		commonFunctions.selectDropdownbyNum(pageData.physLocation, 1);
		pageData.interfaceName.sendKeys(mockVals.interfaceName);
		pageData.ipAddress.sendKeys(mockVals.ipAddress);
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
		let tableColumns = element.all(by.css('.ag-header-cell'));
		expect(tableColumns.count()).toBe(9);
	});

	it('should verify the new Server and then update Server', function() {
		console.log('Verifying new server added and updating ' + mockVals.hostName);
		browser.sleep(1000);
		element(by.cssContainingText('.ag-cell', mockVals.hostName)).click();
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

	it('should navigate back to the new server and view the delivery services assigned to the server', function() {
		console.log('Managing the delivery services of ' + mockVals.hostName);
		browser.navigate().back();
		pageData.moreBtn.click();
		pageData.viewDeliveryServicesMenuItem.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl)+"#!/servers/[0-9]+/delivery-services");
	});

	it('should ensure you cannot clone delivery service assignments because there are no delivery services assigned to the server', function() {
		console.log('Ensure you cannot clone delivery service assignments for ' + mockVals.hostName);
		pageData.moreBtn.click();
		expect(element(by.css('.clone-ds-assignments')).isPresent()).toEqual(false);
	});

});
