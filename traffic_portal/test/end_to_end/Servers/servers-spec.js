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

describe('Traffic Portal Servers Test Suite', function() {

	var pageData = new pd();
	var mockVals = {
		status: "OFFLINE",
		hostName: "testHost",
		domainName: "servertest.com",
		cdn: "CDN-in-a-Box",
		cachegroup: "CDN_in_a_Box_Edge",
		type: "EDGE",
		profile: "ATS_EDGE_TIER_CACHE",
		interfaceName: "testInterfaceName",
		ipAddress: "10.42.80.118",
		ipNetmask: "255.255.255.252",
		ipGateway: "10.42.80.117",
		interfaceMtu: "9000",
		physLocation: "CDN-in-a-Box",
	};

	it('should go to the Servers page', function() {
		browser.get(browser.baseUrl + "/#!/servers");
		expect(browser.getCurrentUrl()).toEqual(browser.baseUrl+"/#!/servers");
	});

    it('should open new Servers form page', function() {
	    browser.driver.findElement(by.name('createServersButton')).click();
	    expect(browser.getCurrentUrl()).toEqual(browser.baseUrl+"/#!/servers/new");
    });

	it('should fill out form, create button is enabled and submit', function () {
		expect(pageData.createButton.isEnabled()).toBe(false);
		pageData.status.click();
		pageData.status.sendKeys(mockVals.status);
		pageData.hostName.sendKeys(mockVals.hostName);
		pageData.domainName.sendKeys(mockVals.domainName);
		pageData.cdn.click();
		pageData.cdn.sendKeys(mockVals.cdn);
		pageData.cachegroup.click();
		pageData.cachegroup.sendKeys(mockVals.cachegroup);
		pageData.type.click();
		pageData.type.sendKeys(mockVals.type);
		pageData.profile.click();
		pageData.profile.sendKeys(mockVals.profile);
		pageData.interfaceName.sendKeys(mockVals.interfaceName);
		pageData.ipAddress.sendKeys(mockVals.ipAddress);
		pageData.ipNetmask.sendKeys(mockVals.ipNetmask);
		pageData.ipGateway.sendKeys(mockVals.ipGateway);
		pageData.interfaceMtu.sendKeys(mockVals.interfaceMtu);
		pageData.physLocation.click();
		pageData.physLocation.sendKeys(mockVals.physLocation);
		expect(pageData.createButton.isEnabled()).toBe(true);
		pageData.createButton.click();
		expect(browser.getCurrentUrl()).toEqual(browser.baseUrl+"/#!/servers");
	});

	it('should verify the new Server and then update Server', function() {
		browser.sleep(1000);
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
		pageData.deleteButton.click();
		pageData.confirmWithNameInput.sendKeys(mockVals.hostName);
		pageData.deletePermanentlyButton.click();
	});
});
