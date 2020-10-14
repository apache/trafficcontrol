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

describe('Traffic Portal Jobs Test Suite', function() {
	const pageData = new pd();
	const commonFunctions = new cfunc();
	const newJob = {
		regex: '/foo.png',
		ttl: 24
	};

	it('should go to the jobs page', function() {
		console.log("Go to the jobs page");
		browser.setLocation("jobs");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/jobs");
	});

	it('should open new job form page', function() {
		console.log("Open new job form page");
		pageData.moreBtn.click();
		pageData.createJobMenuItem.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/jobs/new");
	});

	it('should build a new job', function () {
		console.log("Building a new job");
		expect(pageData.createButton.isEnabled()).toBe(false);
		commonFunctions.selectDropdownbyNum(pageData.deliveryservice, 1);
		pageData.regex.sendKeys(newJob.regex);
		pageData.ttl.sendKeys(newJob.ttl);
		expect(pageData.createButton.isEnabled()).toBe(true);
		pageData.createButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/jobs");
	});

});
