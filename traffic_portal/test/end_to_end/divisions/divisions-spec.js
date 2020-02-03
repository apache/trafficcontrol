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

describe('Traffic Portal Divisions Test Suite', function() {
	const pageData = new pd();
	const commonFunctions = new cfunc();
	const myNewDiv = {
		name: 'division-' + Math.random().toString(36).substring(2, 15)
	};

	it('should go to the divisions page', function() {
		console.log("Go to the divisions page");
		browser.setLocation("divisions");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/divisions");
	});

	it('should verify CSV link exists ', function() {
		console.log("Verify CSV button exists");
		expect(element(by.css('.dt-button.buttons-csv')).isPresent()).toBe(true);
	});

	it('should open new division form page', function() {
		console.log("Open new division form page");
		browser.driver.findElement(by.name('createDivisionButton')).click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/divisions/new");
	});

	it('should create a new division', function () {
		console.log("Creating a new division");
		expect(pageData.createButton.isEnabled()).toBe(false);
		pageData.name.sendKeys(myNewDiv.name);
		expect(pageData.createButton.isEnabled()).toBe(true);
		pageData.createButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/divisions");
	});

	it('should update a division', function() {
		console.log('Updating the new division: ' + myNewDiv.name);
		pageData.searchFilter.sendKeys(myNewDiv.name);
		element.all(by.repeater('d in ::divisions')).filter(function(row){
			return row.element(by.name('name')).getText().then(function(val){
				return val === myNewDiv.name;
			});
		}).get(0).click();
		pageData.name.clear();
		pageData.name.sendKeys(myNewDiv.name + ' updated');
		pageData.updateButton.click();
		expect(pageData.name.getText() === myNewDiv.name + ' updated');
	});

});
