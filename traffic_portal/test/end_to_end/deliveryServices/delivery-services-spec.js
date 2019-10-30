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

describe('Traffic Portal Delivery Services Suite', function() {

	const pageData = new pd();
	const commonFunctions = new cfunc();
	const mockVals = {
		dsTypes: {
			anyMap: [
				"ANY_MAP"
			],
			dns: [
				"DNS",
				"DNS_LIVE_NATNL",
				"DNS_LIVE"
			],
			http: [
				"HTTP",
				"HTTP_NO_CACHE",
				"HTTP_LIVE",
				"HTTP_LIVE_NATNL"
			],
			steering: [
				"STEERING",
				"CLIENT_STEERING"
			]
		},
		anyMapXmlId: "any-map-xml-id-" + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz'),
		dnsXmlId: "dns-xml-id-" + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz'),
		httpXmlId: "http-xml-id-" + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz'),
		steeringXmlId: "http-xml-id-" + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz'),
		longDesc: "This is only a test delivery service that should be disposed of by Automated UI Testing."
	};

	it('should open delivery services page', function() {
		console.log('Opening delivery services page');
		browser.setLocation("delivery-services");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services");
	});

	it('should verify CSV link exists ', function() {
		console.log("Verify CSV button exists");
		expect(element(by.css('.dt-button.buttons-csv')).isPresent()).toBe(true);
	});

	// ANY_MAP delivery service

	it('should click new delivery service and select ANY_MAP category from the dropdown', function() {
		console.log('Clicked Create New and selecting ANY_MAP');
		browser.driver.findElement(by.name('createDeliveryServiceButton')).click();
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(false);
		browser.driver.findElement(by.name('selectFormDropdown')).sendKeys('ANY_MAP');
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(true);
		pageData.selectFormSubmitButton.click();
	});

	it('should populate and submit the delivery service form', function() {
		console.log('Creating a DS for ' + mockVals.anyMapXmlId);
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services/new?type=ANY_MAP");
		expect(pageData.createButton.isEnabled()).toBe(false);
		// set required fields
		// set xml id
		pageData.xmlId.sendKeys(mockVals.anyMapXmlId);
		// set display name
		pageData.displayName.sendKeys(mockVals.anyMapXmlId);
		// set active status
		pageData.active.click();
		pageData.active.sendKeys('Active');
		// set content routing type
		pageData.type.click();
		pageData.type.sendKeys(mockVals.dsTypes.anyMap[0]);
		// set tenant
		commonFunctions.selectDropdownbyNum(pageData.tenantId, 1);
		// set cdn
		commonFunctions.selectDropdownbyNum(pageData.cdn, 1);
		// all required fields have been set, create button should be enabled
		expect(pageData.createButton.isEnabled()).toBe(true);
		pageData.createButton.click();
	});

	it('should back out to delivery services page and verify the new ANY_MAP delivery service and update it', function() {
		console.log('Verifying that ' + mockVals.anyMapXmlId + ' exists');
		browser.setLocation("delivery-services");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services");
	});

	it('should toggle the visibility of the first table column ', function() {
		browser.driver.findElement(by.id('toggleColumns')).click();
		let first = element.all(by.css('input[type=checkbox]')).first();
		expect(first.isSelected()).toBe(true);
		first.click();
		expect(first.isSelected()).toBe(false);
		let tableColumns = element.all(by.css('#deliveryServicesTable tr:first-child td'));
		expect(tableColumns.count()).toBe(10);
	});

	it('should update the ANY_MAP delivery service', function() {
		console.log('Updating the ANY_MAP delivery service for ' + mockVals.anyMapXmlId);
		pageData.searchFilter.clear().then(function() {
			pageData.searchFilter.sendKeys(mockVals.anyMapXmlId);
		});
		element.all(by.repeater('ds in ::deliveryServices')).filter(function(row){
			return row.element(by.name('xmlId')).getText().then(function(val){
				return val.toString() === mockVals.anyMapXmlId.toString();
			});
		}).get(0).click();
		expect(pageData.updateButton.isEnabled()).toBe(false);
		expect(pageData.xmlId.getAttribute('readonly')).toBe('true');
		pageData.displayName.clear().then(function() {
			pageData.displayName.sendKeys("Updated display name");
		});
		expect(pageData.updateButton.isEnabled()).toBe(true);
		pageData.updateButton.click();
		expect(pageData.displayName.getText() === "Updated display name");
	});

	it('should add a required server capability to the ANY_MAP delivery service', function() {
		console.log('Adding required server capability to ' + mockVals.anyMapXmlId);
		pageData.moreBtn.click();
		pageData.viewCapabilitiesMenuItem.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services/[0-9]+/required-server-capabilities");
		pageData.addCapabilityBtn.click();
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(false);
		commonFunctions.selectDropdownbyNum(pageData.selectFormDropdown, 1);
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(true);
		pageData.selectFormSubmitButton.click();
		element.all(by.css('tbody tr')).then(function(totalRows) {
			expect(totalRows.length).toBe(1);
		});
	});

	it('should navigate back to the ANY_MAP delivery service and delete it', function() {
		console.log('Deleting ' + mockVals.anyMapXmlId);
		pageData.dsLink.click();
		pageData.deleteButton.click();
		pageData.confirmWithNameInput.sendKeys(mockVals.anyMapXmlId);
		pageData.deletePermanentlyButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services");
	});

	// DNS delivery service

	it('should click new delivery service and select DNS category from the dropdown', function() {
		console.log('Clicked Create New and selecting DNS');
		browser.driver.findElement(by.name('createDeliveryServiceButton')).click();
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(false);
		browser.driver.findElement(by.name('selectFormDropdown')).sendKeys('DNS');
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(true);
		pageData.selectFormSubmitButton.click();
	});

	it('should populate and submit the ds form', function() {
		console.log('Creating a DS for ' + mockVals.dnsXmlId);
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services/new?type=DNS");
		expect(pageData.createButton.isEnabled()).toBe(false);
		// set required fields
		// set xml id
		pageData.xmlId.sendKeys(mockVals.dnsXmlId);
		// set display name
		pageData.displayName.sendKeys(mockVals.dnsXmlId);
		// set active status
		pageData.active.click();
		pageData.active.sendKeys('Active');
		// set content routing type
		pageData.type.click();
		pageData.type.sendKeys(mockVals.dsTypes.dns[0]);
		// set tenant
		commonFunctions.selectDropdownbyNum(pageData.tenantId, 1);
		// set cdn
		commonFunctions.selectDropdownbyNum(pageData.cdn, 1);
		// set origin server
		pageData.orgServerFqdn.sendKeys('http://' + mockVals.dnsXmlId + '.com');
		// set protocol
		commonFunctions.selectDropdownbyNum(pageData.protocol, 1);
		// all required fields have been set, create button should be enabled
		expect(pageData.createButton.isEnabled()).toBe(true);
		pageData.createButton.click();
	});

	it('should back out to delivery services page and verify the new DNS delivery service and update it', function() {
		console.log('Verifying that ' + mockVals.dnsXmlId + ' exists');
		browser.setLocation("delivery-services");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services");
	});

	it('should update the DNS delivery service', function() {
		console.log('Updating the DNS delivery service for ' + mockVals.dnsXmlId);
		pageData.searchFilter.clear().then(function() {
			pageData.searchFilter.sendKeys(mockVals.dnsXmlId);
		});
		element.all(by.repeater('ds in ::deliveryServices')).filter(function(row){
			return row.element(by.name('xmlId')).getText().then(function(val){
				return val.toString() === mockVals.dnsXmlId.toString();
			});
		}).get(0).click();
		expect(pageData.updateButton.isEnabled()).toBe(false);
		expect(pageData.xmlId.getAttribute('readonly')).toBe('true');
		pageData.displayName.clear().then(function() {
			pageData.displayName.sendKeys("Updated display name");
		});
		expect(pageData.updateButton.isEnabled()).toBe(true);
		pageData.updateButton.click();
		expect(pageData.displayName.getText() === "Updated display name");
	});

	it('should add a required server capability to the DNS delivery service', function() {
		console.log('Adding required server capability to ' + mockVals.dnsXmlId);
		pageData.moreBtn.click();
		pageData.viewCapabilitiesMenuItem.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services/[0-9]+/required-server-capabilities");
		pageData.addCapabilityBtn.click();
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(false);
		commonFunctions.selectDropdownbyNum(pageData.selectFormDropdown, 1);
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(true);
		pageData.selectFormSubmitButton.click();
		element.all(by.css('tbody tr')).then(function(totalRows) {
			expect(totalRows.length).toBe(1);
		});
	});

	it('should navigate back to the DNS delivery service and delete it', function() {
		console.log('Deleting ' + mockVals.dnsXmlId);
		pageData.dsLink.click();
		pageData.deleteButton.click();
		pageData.confirmWithNameInput.sendKeys(mockVals.dnsXmlId);
		pageData.deletePermanentlyButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services");
	});

	// HTTP delivery service

	it('should click new delivery service and select HTTP category from the dropdown', function() {
		console.log('Clicked Create New and selecting HTTP');
		browser.driver.findElement(by.name('createDeliveryServiceButton')).click();
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(false);
		browser.driver.findElement(by.name('selectFormDropdown')).sendKeys('HTTP');
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(true);
		pageData.selectFormSubmitButton.click();
	});

	it('should populate and submit the delivery service form', function() {
		console.log('Creating a HTTP DS for ' + mockVals.dnsXmlId);
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services/new?type=HTTP");
		expect(pageData.createButton.isEnabled()).toBe(false);
		// set required fields
		// set xml id
		pageData.xmlId.sendKeys(mockVals.httpXmlId);
		// set display name
		pageData.displayName.sendKeys(mockVals.httpXmlId);
		// set active status
		pageData.active.click();
		pageData.active.sendKeys('Active');
		// set content routing type
		pageData.type.click();
		pageData.type.sendKeys(mockVals.dsTypes.http[0]);
		// set tenant
		commonFunctions.selectDropdownbyNum(pageData.tenantId, 1);
		// set cdn
		commonFunctions.selectDropdownbyNum(pageData.cdn, 1);
		// set origin server
		pageData.orgServerFqdn.sendKeys('http://' + mockVals.httpXmlId + '.com');
		// set protocol
		commonFunctions.selectDropdownbyNum(pageData.protocol, 1);
		// all required fields have been set, create button should be enabled
		expect(pageData.createButton.isEnabled()).toBe(true);
		pageData.createButton.click();
	});

	it('should back out to delivery services page and verify the new HTTP delivery service and update it', function() {
		console.log('Verifying that ' + mockVals.httpXmlId + ' exists');
		browser.setLocation("delivery-services");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services");
	});

	it('should update the HTTP delivery service', function() {
		console.log('Updating the HTTP delivery service for ' + mockVals.httpXmlId);
		pageData.searchFilter.clear().then(function() {
			pageData.searchFilter.sendKeys(mockVals.httpXmlId);
		});
		element.all(by.repeater('ds in ::deliveryServices')).filter(function(row){
			return row.element(by.name('xmlId')).getText().then(function(val){
				return val.toString() === mockVals.httpXmlId.toString();
			});
		}).get(0).click();
		expect(pageData.updateButton.isEnabled()).toBe(false);
		expect(pageData.xmlId.getAttribute('readonly')).toBe('true');
		pageData.displayName.clear().then(function() {
			pageData.displayName.sendKeys("Updated display name");
		});
		expect(pageData.updateButton.isEnabled()).toBe(true);
		pageData.updateButton.click();
		expect(pageData.displayName.getText() === "Updated display name");
	});

	it('should add a required server capability to the HTTP delivery service', function() {
		console.log('Adding required server capability to ' + mockVals.httpXmlId);
		pageData.moreBtn.click();
		pageData.viewCapabilitiesMenuItem.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services/[0-9]+/required-server-capabilities");
		pageData.addCapabilityBtn.click();
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(false);
		commonFunctions.selectDropdownbyNum(pageData.selectFormDropdown, 1);
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(true);
		pageData.selectFormSubmitButton.click();
		element.all(by.css('tbody tr')).then(function(totalRows) {
			expect(totalRows.length).toBe(1);
		});
	});

	it('should navigate back to the HTTP delivery service and delete it', function() {
		console.log('Deleting ' + mockVals.httpXmlId);
		pageData.dsLink.click();
		pageData.deleteButton.click();
		pageData.confirmWithNameInput.sendKeys(mockVals.httpXmlId);
		pageData.deletePermanentlyButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services");
	});

	// Steering delivery service

	it('should click new delivery service and select Steering category from the dropdown', function() {
		console.log('Clicked Create New and selecting Steering');
		browser.driver.findElement(by.name('createDeliveryServiceButton')).click();
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(false);
		browser.driver.findElement(by.name('selectFormDropdown')).sendKeys('STEERING');
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(true);
		pageData.selectFormSubmitButton.click();
	});

	it('should populate and submit the delivery service form', function() {
		console.log('Creating a Steering DS for ' + mockVals.dnsXmlId);
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services/new?type=STEERING");
		expect(pageData.createButton.isEnabled()).toBe(false);
		// set required fields
		// set xml id
		pageData.xmlId.sendKeys(mockVals.steeringXmlId);
		// set display name
		pageData.displayName.sendKeys(mockVals.steeringXmlId);
		// set active status
		pageData.active.click();
		pageData.active.sendKeys('Active');
		// set content routing type
		pageData.type.click();
		pageData.type.sendKeys(mockVals.dsTypes.steering[0]);
		// set tenant
		commonFunctions.selectDropdownbyNum(pageData.tenantId, 1);
		// set cdn
		commonFunctions.selectDropdownbyNum(pageData.cdn, 1);
		// set protocol
		commonFunctions.selectDropdownbyNum(pageData.protocol, 1);
		// all required fields have been set, create button should be enabled
		expect(pageData.createButton.isEnabled()).toBe(true);
		pageData.createButton.click();
	});

	it('should back out to delivery services page and verify the new Steering delivery service and update it', function() {
		console.log('Verifying that ' + mockVals.steeringXmlId + ' exists');
		browser.setLocation("delivery-services");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services");
	});

	it('should update the Steering delivery service', function() {
		console.log('Updating the Steering delivery service for ' + mockVals.steeringXmlId);
		pageData.searchFilter.clear().then(function() {
			pageData.searchFilter.sendKeys(mockVals.steeringXmlId);
		});
		element.all(by.repeater('ds in ::deliveryServices')).filter(function(row){
			return row.element(by.name('xmlId')).getText().then(function(val){
				return val.toString() === mockVals.steeringXmlId.toString();
			});
		}).get(0).click();
		expect(pageData.updateButton.isEnabled()).toBe(false);
		expect(pageData.xmlId.getAttribute('readonly')).toBe('true');
		pageData.displayName.clear().then(function() {
			pageData.displayName.sendKeys("Updated display name");
		});
		expect(pageData.updateButton.isEnabled()).toBe(true);
		pageData.updateButton.click();
		expect(pageData.displayName.getText() === "Updated display name");
	});

	it('should add a required server capability to the STEERING delivery service', function() {
		console.log('Adding required server capability to ' + mockVals.steeringXmlId);
		pageData.moreBtn.click();
		pageData.viewCapabilitiesMenuItem.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services/[0-9]+/required-server-capabilities");
		pageData.addCapabilityBtn.click();
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(false);
		commonFunctions.selectDropdownbyNum(pageData.selectFormDropdown, 1);
		expect(pageData.selectFormSubmitButton.isEnabled()).toBe(true);
		pageData.selectFormSubmitButton.click();
		element.all(by.css('tbody tr')).then(function(totalRows) {
			expect(totalRows.length).toBe(1);
		});
	});

	it('should navigate back to the STEERING delivery service and delete it', function() {
		console.log('Deleting ' + mockVals.steeringXmlId);
		pageData.dsLink.click();
		pageData.deleteButton.click();
		pageData.confirmWithNameInput.sendKeys(mockVals.steeringXmlId);
		pageData.deletePermanentlyButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/delivery-services");
	});

});
