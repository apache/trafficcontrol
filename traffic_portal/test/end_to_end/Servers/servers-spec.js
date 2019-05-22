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

	var pageData = new pd();
	var commonFunctions = new cfunc();
	var mockVals = {
		status: "OFFLINE",
		hostName: "testHost-" + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz'),
		domainName: "servertest.com",
		interfaceName: "testInterfaceName",
		ipAddress: "10.42.80.118",
		ipNetmask: "255.255.255.252",
		ipGateway: "10.42.80.117",
		interfaceMtu: "9000",
	};

	it('should go to the Servers page', function() {
		console.log('Looading Configure/Servers');
		browser.get(browser.baseUrl + "/#!/servers");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/servers");
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
		pageData.type.sendKeys('MID');
		pageData.updateButton.click();
		expect(pageData.domainName.getText() === 'testupdated.com');
		expect(pageData.type.getText() === 'MID');
	});

	it('should delete the new Server', function() {
		console.log('Deleting the server ' + mockVals.hostName);
		pageData.deleteButton.click();
		pageData.confirmWithNameInput.sendKeys(mockVals.hostName);
		pageData.deletePermanentlyButton.click();
	});
});
