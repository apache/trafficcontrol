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

describe('Traffic Portal Delivery Service Requests', function() {

	const pageData = new pd();
	const commonFunctions = new cfunc();
	const mockVals = {
		dsType: ["ANY MAP", "DNS", "HTTP", "STEERING"],
		active: "Active",
		xmlId: "xml-id-" + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz'),
		displayName: "dsTest",
		orgServerFqdn: "http://dstest.com",
		longDesc: "This is only a test that should be disposed of by Automated UI Testing.",
		commentInput: "This is the second comment"
	};

	it('should open ds services page and click button to create a new one', function() {
		console.log('Opening delivery service requests page');
		browser.setLocation("delivery-services");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services");
	});

	it('should create and select type of ds from the dropdown and confirm', function() {
		console.log('Clicked Create New and selecting a type');
		browser.driver.findElement(by.name('createDeliveryServiceButton')).click();
		browser.sleep(1000);
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(false);
		browser.driver.findElement(by.name('selectFormDropdown')).sendKeys(mockVals.dsType[1]);
		browser.sleep(250);
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(true);
		pageData.selectFormSubmitButton.click();
	});

	it('should populate and submit the ds form', function() {
		console.log('Filling out form for ' + mockVals.xmlId);
		browser.sleep(250);
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services/new?type=" + mockVals.dsType[1]);
		expect(pageData.createButton.isEnabled()).toBe(false);
		pageData.active.click();
		pageData.active.sendKeys(mockVals.active);
		commonFunctions.selectDropdownbyNum(pageData.type, 1);
		pageData.xmlId.sendKeys(mockVals.xmlId);
		pageData.displayName.sendKeys(mockVals.displayName);
		commonFunctions.selectDropdownbyNum(pageData.tenantId, 1);
		commonFunctions.selectDropdownbyNum(pageData.cdn, 1);
		pageData.orgServerFqdn.sendKeys(mockVals.orgServerFqdn);
		commonFunctions.selectDropdownbyNum(pageData.protocol, 1);
		pageData.longDesc.sendKeys(mockVals.longDesc);
		expect(pageData.createButton.isEnabled()).toBe(true);
		pageData.createButton.click();
		browser.sleep(250);
	});

	it('should select a status from the dropdown and add comment before submitting', function() {
		browser.sleep(250);
		commonFunctions.selectDropdownbyNum(pageData.requestStatus, 2);
		pageData.dialogComment.sendKeys('This is comment one');
		browser.sleep(250);
		expect(pageData.dialogSubmit.isEnabled()).toBe(true);
		pageData.dialogSubmit.click();
		browser.sleep(250);
	});

	it('should redirect to delivery-service-requests page', function() {
		console.log('Backing out and verifying ' + mockVals.xmlId + ' exists');
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-service-requests");
	});

	it('should open up and update the ds', function() {
		console.log('Updating the form for ' + mockVals.xmlId);
		browser.sleep(250);
		element.all(by.repeater('request in ::dsRequests')).filter(function(row){
			return row.element(by.name('xmlId')).getText().then(function(val){
				console.log(val + " this is my val " + mockVals.xmlId);
				return val.toString() === mockVals.xmlId.toString();
			});
		}).get(0).click();
		browser.sleep(250);
		expect(pageData.updateButton.isEnabled()).toBe(false);
		pageData.displayName.sendKeys(mockVals.displayName + "updated");
		expect(pageData.updateButton.isEnabled()).toBe(true);
		pageData.updateButton.click();
		browser.sleep(250);
		expect(pageData.displayName.getText() === mockVals.displayName + "updated");
	});

	it('should select a status from the dropdown and add comment before submitting', function() {
		browser.sleep(250);
		commonFunctions.selectDropdownbyNum(pageData.requestStatus, 2);
		pageData.dialogComment.sendKeys('This is comment two');
		browser.sleep(250);
		expect(pageData.dialogSubmit.isEnabled()).toBe(true);
		pageData.dialogSubmit.click();
		browser.sleep(250);
	});

	it('should add a comment', function () {
		console.log('Adding Comment');
		pageData.newCommentButton.click();
		browser.sleep(250);
		pageData.commentInput.sendKeys(mockVals.commentInput);
		pageData.createCommentButton.click();
	});

	it('should edit a comment', function () {
		console.log('Editing Comment');
		browser.sleep(250);
		element.all(by.css('.link.action-link')).first().click();
		browser.sleep(250);
		pageData.commentInput.sendKeys(mockVals.commentInput);
		pageData.updateCommentButton.click();
	});

	it('should delete a comment', function () {
		console.log('Deleting Comment');
		browser.sleep(250);
		element.all(by.css('.link.action-link')).get(1).click();
		browser.sleep(250);
		pageData.yesButton.click();
	});

	it('should delete the ds request', function() {
		console.log('Deleting ' + mockVals.xmlId);
		pageData.deleteButton.click();
		pageData.confirmWithNameInput.sendKeys(mockVals.xmlId + ' request');
		pageData.deletePermanentlyButton.click();
		browser.sleep(250);
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-service-requests");
	});
});
