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
		fullName: 'test',
		email: 'test@viper.com',
		roleName: 'admin',
		tenantId: ' - root',
		localPasswd: 'test@123',
		confirmLocalPasswd: 'test@123'
	};
	const myNewRegisteredUser = {
		email: 'test1@viper.com',
		roleName: 'operations',
		tenantId: ' -- tenant01'
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
		pageData.roleName.sendKeys(myNewUser.roleName);
		pageData.tenantId.sendKeys(myNewUser.tenantId);
		pageData.localPasswd.sendKeys(myNewUser.localPasswd);
		pageData.confirmLocalPasswd.sendKeys(myNewUser.confirmLocalPasswd);
		expect(pageData.createButton.isEnabled()).toBe(true);
		pageData.createButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/users");
	});

	it('should update a user', function() {
		console.log('Updating the new user: ' + myNewUser.username);
		pageData.searchFilter.sendKeys(myNewUser.fullName);
		element.all(by.repeater('u in ::users')).filter(function(row){
			return row.element(by.name('fullName')).getText().then(function(val){
				return val === myNewUser.fullName;
			});
		}).get(0).click();
		pageData.fullName.clear();
		pageData.username.sendKeys(myNewUser.username + ' updated');
		pageData.updateButton.click();
		expect(pageData.username.getText() === myNewUser.username + ' updated');
	});

	it('should open new register users form page', function() {
		console.log("Open new register users form page");
		browser.driver.findElement(by.name('createRegisterUserButton')).click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/users/register");
	});

	it('should create a new registered user', function () {
		console.log("Creating a new registered user");
		expect(pageData.registerButton.isEnabled()).toBe(false);
		pageData.email.sendKeys(myNewRegisteredUser.email);
		pageData.roleName.sendKeys(myNewRegisteredUser.roleName);
		pageData.tenantId.sendKeys(myNewRegisteredUser.tenantId);
		expect(pageData.registerButton.isEnabled()).toBe(true);
		pageData.registerButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/users/register");
	});

	// it('should update a new registered user', function() {
	// 	expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/users");
	// 	console.log('Updating the new registered user: ' + myNewRegisteredUser.username);
	// 	pageData.searchFilter.sendKeys(myNewRegisteredUser.fullName);
	// 	element.all(by.repeater('u in ::users')).filter(function(row){
	// 		return row.element(by.name('fullName')).getText().then(function(val){
	// 			return val === myNewRegisteredUser.fullName;
	// 		});
	// 	}).get(0).click();
	// 	pageData.fullName.clear();
	// 	pageData.username.sendKeys(myNewRegisteredUser.username + ' updated');
	// 	pageData.updateButton.click();
	// 	expect(pageData.username.getText() === myNewRegisteredUser.username + ' updated');
	// });

});
