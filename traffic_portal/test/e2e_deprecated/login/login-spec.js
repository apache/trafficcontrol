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

var cfunc = require('../common/commonFunctions.js');

describe('Traffic Portal Login Test Suite', function() {
	const commonFunctions = new cfunc();

	beforeEach(function() {
		browser.get(browser.baseUrl + '/#!/cdns');
		browser.wait(function() {
			return element(by.name('loginUsername')).isPresent();
		}, 5000);
	});

	it('should not show environment banner in prod mode', function() {
		console.log('Verifying environment banner does not have the prod class');
		expect(element(by.css('.enviro-banner.prod')).isPresent()).toBe(false);
	});

	it('should fail login to Traffic Portal with bad user', function() {
		console.log('Negative login test');
		browser.driver.findElement(by.name('loginUsername')).sendKeys('badUser');
		browser.driver.findElement(by.name('loginPass')).sendKeys('badPassword');
		browser.driver.findElement(by.name('loginSubmit')).click();
		browser.sleep(250);
		expect(browser.driver.findElement(by.css('div.ng-binding')).getText()).toEqual('Invalid username or password.');
	});

	it('should successfully login to Traffic Portal', function() {
		console.log('Logging in to Traffic Portal "' + browser.baseUrl + '" with user "' + browser.params.adminUser + '"');
		browser.driver.findElement(by.name('loginUsername')).sendKeys(browser.params.adminUser);
		browser.driver.findElement(by.name('loginPass')).sendKeys(browser.params.adminPassword);
		browser.driver.findElement(by.name('loginSubmit')).click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/cdns");
	});
});
