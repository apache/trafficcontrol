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
		commentInput: "This is the second comment",
		tenantId: "- root",
		protocol: "HTTP to HTTPS"
	};
	const repeater = 'request in ::dsRequests';

	it('should create and select type of ds from the dropdown and confirm', async () => {
		console.log('Opening delivery service requests page');
		await browser.setLocation("delivery-services");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services");
		console.log('Clicked Create New and selecting a type');
		await pageData.createDeliveryServiceButton.click();
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(false);
		await pageData.selectFormDropdown.click().sendKeys(mockVals.dsType[1]);
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(true);
		await pageData.selectFormSubmitButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services/new?type=" + mockVals.dsType[1]);
	});

	it('should populate and submit the ds form', async () => {
		console.log('Filling out form for ' + mockVals.xmlId);
		expect(pageData.createButton.isEnabled()).toBe(false);
		await pageData.active.click().sendKeys(mockVals.active);
		await commonFunctions.selectDropdownByLabel(pageData.type, mockVals.dsType[1]);
		await pageData.xmlId.sendKeys(mockVals.xmlId);
		await pageData.displayName.sendKeys(mockVals.displayName);
		await commonFunctions.selectDropdownByLabel(pageData.tenantId, mockVals.tenantId);
		await commonFunctions.selectDropdownByNum(pageData.cdn, 1);
		await pageData.orgServerFqdn.sendKeys(mockVals.orgServerFqdn);
		await pageData.protocol.click().sendKeys(mockVals.protocol);
		await pageData.longDesc.sendKeys(mockVals.longDesc);
		expect(pageData.createButton.isEnabled()).toBe(true);
		await pageData.createButton.click();
	});

	it('should select a status from the dropdown and add comment before submitting', async () => {
		await commonFunctions.selectDropdownByNum(pageData.requestStatus, 2);
		await pageData.dialogComment.sendKeys('This is comment one');
		expect(pageData.dialogSubmit.isEnabled()).toBe(true);
		await pageData.dialogSubmit.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-service-requests");
	});

	it('should open up and update the ds', async () => {
		console.log('Updating the form for ' + mockVals.xmlId);
		await commonFunctions.clickTableEntry(pageData.searchFilter, mockVals.xmlId, repeater);
		expect(pageData.updateButton.isEnabled()).toBe(false);
		await pageData.displayName.sendKeys(mockVals.displayName + "updated");
		expect(pageData.updateButton.isEnabled()).toBe(true);
		await pageData.updateButton.click();
		expect(pageData.displayName.getText() === mockVals.displayName + "updated");
	});

	it('should select a status from the dropdown and add comment before submitting', async () => {
		await commonFunctions.selectDropdownByNum(pageData.requestStatus, 2);
		await pageData.dialogComment.sendKeys('This is comment two');
		expect(pageData.dialogSubmit.isEnabled()).toBe(true);
		await pageData.dialogSubmit.click();
	});

	it('should add a comment', async () => {
		console.log('Adding Comment');
		await pageData.newCommentButton.click();
		await pageData.commentInput.sendKeys(mockVals.commentInput);
		await pageData.createCommentButton.click();
	});

	it('should edit a comment', async () => {
		console.log('Editing Comment');
		await element.all(by.css('.link.action-link')).first().click();
		await pageData.commentInput.sendKeys(mockVals.commentInput);
		await pageData.updateCommentButton.click();
	});

	it('should delete a comment', async () => {
		console.log('Deleting Comment');
		await element.all(by.css('.link.action-link')).get(1).click();
		await pageData.yesButton.click();
	});

	it('should delete the ds request', async () => {
		console.log('Deleting ' + mockVals.xmlId);
		await pageData.deleteButton.click();
		await pageData.confirmWithNameInput.sendKeys(mockVals.xmlId + ' request');
		await pageData.deletePermanentlyButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-service-requests");
	});
});
