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
		// Assign this server to a created cache-group, needed for Topologies tests
		pageData.cachegroup.all(by.tagName("option")).then(function(options) {
		    options[options.length-1].click();
		});
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

	it('should clear column filter when column is hidden', function() {
		// Confirm we have rows
		let rows = element.all(by.css("div.ag-row"));
		expect(rows.count()).not.toBe(0);

		// Filter one of our columns
		let firstHeaderCell = element.all(by.css('.ag-header-cell')).first();
		firstHeaderCell.all(by.css('span.ag-header-cell-menu-button')).first().click();
		let filterContainer = element(by.css("div.ag-filter"));
		let filterCell = filterContainer.all(by.css('.ag-input-field-input')).first();
		filterCell.sendKeys("nothingshouldmatchthis", protractor.Key.ENTER);

		// Wait for ag-grid to process changes
		browser.sleep(1000);
		rows = element.all(by.css("div.ag-row"));
		expect(rows.count()).toBe(0);

		// Hide filtered column
		let columnToggle = element(by.id('toggleColumns')).click();
		columnToggle.all(by.css('input[type=checkbox]:checked')).first().click();

		// Wait for ag-grid again
		rows = element.all(by.css("div.ag-row"));
		expect(rows.count()).not.toBe(0);
	});

	it('should verify the new Server and then update Server', function() {
		console.log('Verifying new server added and updating ' + mockVals.hostName);
		browser.sleep(1000);
		let row = element(by.cssContainingText('.ag-cell', mockVals.hostName));
		browser.actions().click(row).perform();
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
		pageData.submitButton.click().then(function() {
			element.all(by.css('tbody tr')).then(function(totalRows) {
				expect(totalRows.length).toBe(1);
			});
		});
	});

	it('should navigate back to the new server and view the delivery services assigned to the server', function() {
		console.log('Managing the delivery services of ' + mockVals.hostName);
		browser.navigate().back();
		pageData.moreBtn.click();
		pageData.viewDeliveryServicesMenuItem.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl)+"#!/servers/[0-9]+/delivery-services");
	});
});
