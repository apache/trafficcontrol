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

describe('Traffic Portal Delivery Services Suite', function() {

	var pageData = new pd();
	var mockVals = {
		dsType: ["ANY MAP", "DNS", "HTTP", "STEERING"],
		active: "true",
		type: "DNS",
		xmlId: 4343432432424,
		displayName: "dsTest",
		tenantId: "- root",
		cdn: "CDN-in-a-Box",
		orgServerFqdn: "http://dstest.com",
		protocol: "0 - HTTP",
		longDesc: "This is only a test that should be disposed of by Automated UI Testing."
	};

	it('should open ds page and click button to create a new one', function() {
		browser.get(browser.baseUrl + "/#!/delivery-services");
		expect(browser.getCurrentUrl()).toEqual(browser.baseUrl+"/#!/delivery-services");
	});

	it('should create and select type of ds from the dropdown and confirm', function() {
		browser.driver.findElement(by.name('createDeliveryServiceButton')).click();
		browser.sleep(1000);
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(false);
		browser.driver.findElement(by.name('selectFormDropdown')).sendKeys(mockVals.dsType[1]);
		browser.sleep(250);
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(true);
		pageData.selectFormSubmitButton.click();
	});

	it('should populate and submit the ds form', function() {
		browser.sleep(250);
		expect(browser.getCurrentUrl()).toEqual(browser.baseUrl+"/#!/delivery-services/new?type=" + mockVals.dsType[1]);
		expect(pageData.createButton.isEnabled()).toBe(false);
		pageData.active.click();
		pageData.active.sendKeys(mockVals.active);
		pageData.type.click();
		pageData.type.sendKeys(mockVals.type);
		pageData.xmlId.sendKeys(mockVals.xmlId);
		pageData.displayName.sendKeys(mockVals.displayName);
		pageData.tenantId.click();
		pageData.tenantId.sendKeys(mockVals.tenantId);
		pageData.cdn.click();
		pageData.cdn.sendKeys(mockVals.cdn);
		pageData.orgServerFqdn.sendKeys(mockVals.orgServerFqdn);
		pageData.protocol.click();
		pageData.protocol.sendKeys(mockVals.protocol);
		pageData.longDesc.sendKeys(mockVals.longDesc);
		expect(pageData.createButton.isEnabled()).toBe(true);
		pageData.createButton.click();
		browser.sleep(250);
	});

	it('should back out to ds page and verify new ds and update it', function() {
		browser.get(browser.baseUrl + "/#!/delivery-services");
		expect(browser.getCurrentUrl()).toEqual(browser.baseUrl+"/#!/delivery-services");
		browser.sleep(250);
		element(by.repeater('ds in ::deliveryServices').row(0)).click();
		// Need to get the below code to function. Currently throwing an index out of bounds error
		// And then remove the single line click above
		//
		// element.all(by.repeater('ds in ::deliveryServices')).filter(function(row){
		// 	return row.element(by.name('xmlId')).getText().then(function(val){
		// 		console.log(val + " this is my val " + row);
		// 		return val === mockVals.xmlId;
		// 	});
		// })
		// 	.get(0).click();
		browser.sleep(250);
		expect(pageData.updateButton.isEnabled()).toBe(false);
		pageData.displayName.sendKeys(mockVals.displayName + "updated");
		expect(pageData.updateButton.isEnabled()).toBe(true);
		pageData.updateButton.click();
		browser.sleep(250);
		expect(pageData.displayName.getText() === mockVals.displayName + "updated");
	});

	it('should delete the ds', function() {
		pageData.deleteButton.click();
		pageData.confirmWithNameInput.sendKeys(mockVals.xmlId);
		pageData.deletePermanentlyButton.click();
		expect(browser.getCurrentUrl()).toEqual(browser.baseUrl+"/#!/delivery-services");
	});
});
