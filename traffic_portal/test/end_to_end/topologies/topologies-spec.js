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

describe('Traffic Portal Topologies Test Suite', function() {
	const pageData = new pd();
	const commonFunctions = new cfunc();
	const ec = protractor.ExpectedConditions;
	const myNewTopology = {
		name: 'topology-' + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz0123456789'),
		desc: 'topology-' + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz0123456789')
	};

	it('should go to the topologies page', function() {
		console.log("Go to the topologies page");
		browser.setLocation("topologies");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/topologies");
	});

	it('should verify CSV link exists ', function() {
		console.log("Verify CSV button exists");
		expect(element(by.css('.dt-button.buttons-csv')).isPresent()).toBe(true);
	});

	it('should open new topology form page', function() {
		console.log("Open new topology form page");
		browser.driver.findElement(by.name('createTopologyBtn')).click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/topologies/new");
	});

	it('should build a new topology', function () {
		console.log("Building a new topology");
		pageData.addChildCacheGroupBtn.click();
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(false);
		browser.driver.findElement(by.name('selectFormDropdown')).sendKeys('EDGE_LOC');
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(true);
		pageData.selectFormSubmitButton.click();
		browser.wait(ec.presenceOf(pageData.selectAllCB), 5000);
		pageData.selectAllCB.click();
		pageData.selectFormSubmitButton.click();
	});

	it('should fill out the rest of the topology form, create button is enabled and submit', function () {
		console.log("Filling out topology form, check create button is enabled and submit");
		expect(pageData.createButton.isEnabled()).toBe(false);
		pageData.name.sendKeys(myNewTopology.name);
		pageData.description.sendKeys(myNewTopology.desc);
		expect(pageData.createButton.isEnabled()).toBe(true);
		pageData.createButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/topologies");
	});

	it('should update the topology', function() {
		console.log('Updating the topology: ' + myNewTopology.name);
		pageData.searchFilter.sendKeys(myNewTopology.name);
		element.all(by.repeater('t in ::topologies')).filter(function(row){
			return row.element(by.name('name')).getText().then(function(val){
				return val === myNewTopology.name;
			});
		}).get(0).click();
		expect(pageData.updateButton.isEnabled()).toBe(false);
		expect(pageData.name.getAttribute('disabled')).toBe('true');
		pageData.description.clear().then(function() {
			pageData.description.sendKeys("Updated description");
		});
		expect(pageData.updateButton.isEnabled()).toBe(true);
		pageData.updateButton.click();
		expect(pageData.description.getText() === "Updated description");
	});

	it('should view all delivery services that utilize the topology', function() {
		console.log('Viewing all delivery services that utilize: ' + myNewTopology.name);
		pageData.moreBtn.click();
		pageData.viewDeliveryServicesMenuItem.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl)+"#!/topologies/delivery-services");
	});

	it('should navigate back to the topology and view all cache groups utilized by the topology', function() {
		console.log('Viewing all cache groups utilized by ' + myNewTopology.name);
		pageData.topLink.click();
		pageData.moreBtn.click();
		pageData.viewCacheGroupsMenuItem.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl)+"#!/topologies/cache-groups");
	});

	it('should navigate back to the topology and view all servers utilized by the topology', function() {
		console.log('Viewing all servers utilized by ' + myNewTopology.name);
		pageData.topLink.click();
		pageData.moreBtn.click();
		pageData.viewServersMenuItem.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl)+"#!/topologies/servers");
	});
});
