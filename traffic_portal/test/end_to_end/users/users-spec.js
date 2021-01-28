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

describe('Traffic Portal Users Test Suite', function() {
	const pageData = new pd();
	const commonFunctions = new cfunc();
	const myNewUser = {
		username: 'user-' + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz0123456789'),
		fullName: 'test-' + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz0123456789'),
		email: 'test@cdn.' + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz') + '.com',
		localPasswd: commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz'),
	};
	const myNewRegisteredUser = {
		email: 'test1@cdn.' + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz') + '.com'
	};

	it('should go to the users page', function() {
		console.log("Go to the users page");
		browser.setLocation("users");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/users");
	});

	it('should verify CSV link exists ', function() {
		console.log("Verify CSV button exists");
		expect(element(by.css('.dt-button.buttons-csv')).isPresent()).toBe(true);
	});

	it('should toggle the visibility of the first table column ', function() {
		browser.driver.findElement(by.id('toggleColumns')).click();
		let first = element.all(by.css('input[type=checkbox]')).first();
		expect(first.isSelected()).toBe(true);
		first.click();
		expect(first.isSelected()).toBe(false);
		let tableColumns = element.all(by.css('#usersTable tr:first-child td'));
		expect(tableColumns.count()).toBe(4);
		first.click();
	});

	it('should open new users form page', function() {
		console.log("Open new users form page");
		browser.driver.findElement(by.name('createUserButton')).click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/users/new");
	});

	it('should create a new user', function () {
		console.log("Creating a new user");
		expect(pageData.createButton.isEnabled()).toBe(false);
		pageData.username.sendKeys(myNewUser.username);
		pageData.fullName.sendKeys(myNewUser.fullName);
		pageData.email.sendKeys(myNewUser.email);
		commonFunctions.selectDropdownbyNum(pageData.roleName, 1);
		commonFunctions.selectDropdownbyNum(pageData.tenantId, 1);
		pageData.localPasswd.sendKeys(myNewUser.localPasswd);
		pageData.confirmLocalPasswd.sendKeys(myNewUser.localPasswd);
		expect(pageData.createButton.isEnabled()).toBe(true);
		pageData.createButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/users");
	});

	it('should update a user', function() {
		console.log('Updating the new user: ' + myNewUser.username);
		browser.sleep(250);
		pageData.searchFilter.sendKeys(myNewUser.username);
		browser.sleep(250);
		element.all(by.repeater('u in ::users')).filter(function(row){
			return row.element(by.name('username')).getText().then(function(val){
				return val === myNewUser.username;
			});
		}).get(0).click();
		browser.sleep(1000);
		pageData.fullName.clear();
		pageData.fullName.sendKeys(myNewUser.fullName + ' updated');
		pageData.updateButton.click();
		expect(pageData.fullName.getText() === myNewUser.fullName + ' updated');
	});

	it('should open new registered users form page', function() {
		console.log("Open new register users form page");
		browser.setLocation("users");
		pageData.registerNewUserButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/users/register");
	});

	it('should create a new registered user', function () {
		console.log("Creating a new registered user");
		expect(pageData.registerEmailButton.isEnabled()).toBe(false);
		pageData.email.sendKeys(myNewRegisteredUser.email);
		commonFunctions.selectDropdownbyNum(pageData.roleName, 2);
		commonFunctions.selectDropdownbyNum(pageData.tenantId, 2);
		expect(pageData.registerEmailButton.isEnabled()).toBe(true);
		pageData.registerEmailButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/users/register");
	});

	it('should update a new registered user', function() {
		console.log('Updating the new registered user: ' + myNewRegisteredUser.email);
		browser.setLocation("users");
		pageData.searchFilter.clear();
		browser.sleep(250);
		pageData.searchFilter.sendKeys(myNewRegisteredUser.email);
		browser.sleep(250);
		element.all(by.repeater('u in ::users')).filter(function(row){
			return row.element(by.name('email')).getText().then(function(val){
				return val === myNewRegisteredUser.email;
			});
		}).get(0).click();
		browser.sleep(1000);
		pageData.fullName.clear();
		pageData.fullName.sendKeys('test1 updated');
		expect(pageData.registerSent.getAttribute('readOnly')).toBe('true');
		pageData.updateButton.click();
		expect(pageData.fullName.getText() === 'test1 updated');
	});

});
