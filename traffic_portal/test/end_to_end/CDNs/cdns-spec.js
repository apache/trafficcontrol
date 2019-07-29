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

describe('Traffic Portal CDNs Test Suite', function() {
	const pageData = new pd();
	const commonFunctions = new cfunc();
	const shuffledText = commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz0123456789');
	const myNewCDN = {
		name : 'cdn-' + shuffledText,
		domainName : 'cdn-' + shuffledText + '.com',
		dnssecEnabled: false,
		numKskDays: commonFunctions.random(365)
	};
	const repeater = 'cdn in ::cdns';

	it('should go to the CDNs page', async () => {
		console.log("Go to the CDNs page");
		await browser.setLocation("cdns");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/cdns");
	});

	it('should open new CDN form page', async () => {
		console.log("Open new CDN form page");
		await pageData.createCdnButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/cdns/new");
	});

	it('should fill out form, create button is enabled and submit', async () => {
		console.log("Filling out form, check create button is enabled and submit");
		expect(pageData.createButton.isEnabled()).toBe(false);
		await pageData.dnssecEnabled.sendKeys(myNewCDN.dnssecEnabled.toString());
		await pageData.name.sendKeys(myNewCDN.name);
		await pageData.domainName.sendKeys(myNewCDN.domainName);
		expect(pageData.createButton.isEnabled()).toBe(true);
		await pageData.createButton.click();
		expect(pageData.successMsg.isPresent()).toBe(true);
		expect(pageData.cdnCreatedText.isPresent()).toBe(true, 'Actual message does not match expected message');
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/cdns");
	});

	it('should verify the new CDN and then update CDN', async () => {
		console.log("Verifying the new CDN and then updating CDN");
		await commonFunctions.clickTableEntry(pageData.searchFilter, myNewCDN.name, repeater);
		await pageData.domainName.clear();
		await pageData.domainName.sendKeys(myNewCDN.domainName + 'updated.com');
		await pageData.dnssecEnabled.sendKeys((!myNewCDN.dnssecEnabled).toString());
		await pageData.updateButton.click();
		expect(pageData.successMsg.isPresent()).toBe(true);
		expect(pageData.cdnUpdatedText.isPresent()).toBe(true, 'Actual message does not match expected message');
		expect(pageData.domainName.getAttribute('value')).toEqual(myNewCDN.domainName + 'updated.com');
	});

	it('should generate DNSSEC keys', async () => {
		console.log("Generating DNSSEC keys for the new CDN and and verifying their expiration date");
		await pageData.moreButton.click();
		await pageData.manageDnssecKeysButton.click();
		expect(pageData.expirationDate.getAttribute('value')).toEqual('');
		await pageData.generateDnssecKeysButton.click();
		await pageData.regenerateButton.click();
		expect(pageData.confirmButton.isEnabled()).toBe(false);
		await pageData.confirmInput.sendKeys(myNewCDN.name);
		expect(pageData.confirmButton.isEnabled()).toBe(true);
		await pageData.confirmButton.click();
		const expirationDate = pageData.expirationDate.getAttribute('value').then((expir) => {return Date.parse(expir + ' UTC');});
		const calculatedExpirationDate = Date.now() + 365*24*60*60*1000;
		expect(expirationDate).toBeCloseTo(calculatedExpirationDate, -4);
	});

	it('should regenerate DNSSEC keys', async () => {
		console.log("Renerating DNSSEC keys and verifying their expiration date");
		await pageData.regenerateDnssecKeysButton.click();
		await pageData.kskExpirationDays.clear().sendKeys(myNewCDN.numKskDays.toString());
		await pageData.regenerateButton.click();
		expect(pageData.confirmButton.isEnabled()).toBe(false);
		await pageData.confirmInput.sendKeys(myNewCDN.name);
		expect(pageData.confirmButton.isEnabled()).toBe(true);
		await pageData.confirmButton.click();
		const expirationDate = pageData.expirationDate.getAttribute('value').then((expir) => {return Date.parse(expir + ' UTC');});
		const calculatedExpirationDate = Date.now() + myNewCDN.numKskDays*24*60*60*1000;
		expect(expirationDate).toBeCloseTo(calculatedExpirationDate, -4);
	});

	it('should regenerate KSK keys', async () => {
		console.log("Regenerating KSK keys and verifying their expiration");
		await pageData.regenerateKskButton.click();
		await pageData.kskExpirationDays.clear().sendKeys(myNewCDN.numKskDays.toString());
		await pageData.generateButton.click();
		expect(pageData.confirmButton.isEnabled()).toBe(false);
		await pageData.confirmInput.sendKeys(myNewCDN.name);
		expect(pageData.confirmButton.isEnabled()).toBe(true);
		await pageData.confirmButton.click();
		const expirationDate = pageData.expirationDate.getAttribute('value').then((expir) => {return Date.parse(expir + ' UTC');});
		const calculatedExpirationDate = Date.now() + myNewCDN.numKskDays*24*60*60*1000;
		expect(expirationDate).toBeCloseTo(calculatedExpirationDate, -4);
	});

});
