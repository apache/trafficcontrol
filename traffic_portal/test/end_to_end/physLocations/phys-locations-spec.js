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

describe('Traffic Portal Phys Locations Test Suite', function() {
	const pageData = new pd();
	const  commonFunctions = new cfunc();
	const myNewPhysLoc = {
		name: 'phys-loc-' + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz0123456789'),
		address: '1200 sycamore lane',
		city: 'Pottersville',
		state: 'AK',
		zip: '12345'
	};

	it('should go to the phys locations page', async () => {
		console.log("Go to the phys locations page");
		await browser.setLocation("phys-locations");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/phys-locations");
	});

	it('should open new phys locations form page', async () => {
		console.log("Open new phys location form page");
		await pageData.createPhysLocationButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/phys-locations/new");
	});

	it('should fill out form, create button is enabled and submit', async () => {
		console.log("Filling out form, check create button is enabled and submit");
		expect(pageData.createButton.isEnabled()).toBe(false);
		await pageData.name.sendKeys(myNewPhysLoc.name);
		await pageData.shortName.sendKeys(myNewPhysLoc.name);
		await pageData.address.sendKeys(myNewPhysLoc.address);
		await pageData.city.sendKeys(myNewPhysLoc.city);
		await pageData.state.sendKeys(myNewPhysLoc.state);
		await pageData.zip.sendKeys(myNewPhysLoc.zip);
		await commonFunctions.selectDropdownByNum(pageData.region, 1);
		expect(pageData.createButton.isEnabled()).toBe(true);
		await pageData.createButton.click();
		expect(pageData.successMsg.isPresent()).toBe(true);
        expect(pageData.physLocationCreatedText.isPresent()).toBe(true, 'Actual message does not match expected message');
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/phys-locations");
	});

});
