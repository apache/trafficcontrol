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

describe('Traffic Portal Types Test Suite', function() {
	const commonFunctions = new cfunc();

	it('should go to the types page', function() {
		console.log("Go to the types page");
		browser.setLocation("types");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/types");
	});

	it('should verify CSV link exists ', function() {
		console.log("Verify CSV button exists");
		expect(element(by.css('.dt-button.buttons-csv')).isPresent()).toBe(true);
	});

	it('should toggle the visibility of the table columns leaving only one visible', function() {
		browser.driver.findElement(by.id('toggleColumns')).click();
		element.all(by.tagName('input[type=checkbox]')).each(function(item) {
			item.click();
		});

		let rowColumns = element.all(by.css('#typesTable tr:first-child td'));
		expect(rowColumns.count()).toBe(1);
	});

});
