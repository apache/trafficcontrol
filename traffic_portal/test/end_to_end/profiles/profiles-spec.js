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

describe('Traffic Portal Profiles Test Suite', function() {
	const pageData = new pd();
	const commonFunctions = new cfunc();
	const myNewProfile = {
		name: 'profile-' + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz0123456789'),
	};

	it('should go to the profiles page', function() {
		console.log("Go to the profiles page");
		browser.setLocation("profiles");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/profiles");
	});

	it('should verify CSV link exists ', function() {
		console.log("Verify CSV button exists");
		expect(element(by.css('.dt-button.buttons-csv')).isPresent()).toBe(true);
	});

	it('should compare profiles', function() {
		pageData.moreBtn.click();
		pageData.compareProfilesMenuItem.click();
		expect(pageData.compareSubmit.isEnabled()).toBe(false);
		commonFunctions.selectDropdownbyNum(pageData.compareDropdown1, 1);
		commonFunctions.selectDropdownbyNum(pageData.compareDropdown2, 2);
		expect(pageData.compareSubmit.isEnabled()).toBe(true);
		pageData.compareSubmit.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl)+"#!/profiles/[0-9]+/[0-9]+/compare/diff");
	});

	it('should open new profile form page', function() {
		console.log("Open new profile form page");
		browser.setLocation("profiles");
		pageData.createProfileButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/profiles/new");
	});

	it('should fill out form, create button is enabled and submit', function () {
		console.log("Filling out form, check create button is enabled and submit");
		expect(pageData.createButton.isEnabled()).toBe(false);
		pageData.name.sendKeys(myNewProfile.name);
		commonFunctions.selectDropdownbyNum(pageData.cdn, 2); // the ALL CDN is first so let's pick a real CDN
		commonFunctions.selectDropdownbyNum(pageData.type, 1);
		pageData.routingDisabled.click();
		pageData.routingDisabled.sendKeys('false');
		pageData.description.sendKeys(myNewProfile.name);
		expect(pageData.createButton.isEnabled()).toBe(true);
		pageData.createButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl)+"#!/profiles/[0-9]+/parameters");
	});

	it('should go back to the profiles page', function() {
		console.log("Go to the profiles page");
		browser.setLocation("profiles");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/profiles");
	});

	it('should toggle the visibility of the first table column ', function () {
		browser.driver.findElement(by.id('toggleColumns')).click();
		let first = element.all(by.css('input[type=checkbox]')).first();
		expect(first.isSelected()).toBe(true);
		first.click();
		expect(first.isSelected()).toBe(false);
		let tableColumns = element.all(by.css('#profilesTable tr:first-child td'));
		expect(tableColumns.count()).toBe(4);
	});


});
