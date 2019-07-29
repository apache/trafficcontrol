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

describe('Traffic Portal Cache Groups Test Suite', function() {
	const pageData = new pd();
	const commonFunctions = new cfunc();
	const myNewCG = {
		name: 'cache-group-' + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz0123456789'),
		latitude: 45,
		longitude: 45,
		type: 'EDGE_LOC'
	};

	it('should go to the cache groups page', async () => {
		console.log("Go to the cache groups page");
		await browser.setLocation("cache-groups");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/cache-groups");
	});

	it('should open new cache group form page', async () => {
		console.log("Open new cache groups form page");
		await pageData.createCacheGroupButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/cache-groups/new");
	});

	it('should fill out form, create button is enabled and submit', async () => {
		console.log("Filling out form, check create button is enabled and submit");
		expect(pageData.createButton.isEnabled()).toBe(false);
		await pageData.name.sendKeys(myNewCG.name);
		await pageData.shortName.sendKeys(myNewCG.name);
		await commonFunctions.selectDropdownByLabel(pageData.type, myNewCG.type);
		await pageData.latitude.sendKeys(myNewCG.latitude);
		await pageData.longitude.sendKeys(myNewCG.longitude);
		expect(pageData.createButton.isEnabled()).toBe(true);
		await pageData.createButton.click();
		expect(pageData.successMsg.isPresent()).toBe(true);
        expect(pageData.cacheGroupCreatedText.isPresent()).toBe(true, 'Actual message does not match expected message');
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/cache-groups");
	});

});
