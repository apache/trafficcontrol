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
import { browser } from "protractor";

import { LoginPage } from "../PageObjects/LoginPage.po"
import { CDNPage } from "../PageObjects/CDNPage.po";
import { TopNavigationPage } from "../PageObjects/TopNavigationPage.po";
import { cdns } from "../Data";


const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const cdnsPage = new CDNPage();

cdns.tests.forEach(async cdnsData =>{
	for (const login of cdnsData.logins) {
		describe(`Traffic Portal - CDN - ${login.description}`, () =>{
			beforeAll(async () => {
				browser.get(browser.params.baseUrl);
				await loginPage.Login(login);
				expect(await loginPage.CheckUserName(login)).toBeTruthy();
				await cdnsPage.openCDNsPage();
			});
			afterAll(async ()=>{
				expect(await topNavigation.Logout()).toBeTruthy();
			});
			afterEach(async ()=> {
				await cdnsPage.openCDNsPage();
			});
			for (const add of cdnsData.add) {
				it(add.description, async () => {
					expect(await cdnsPage.createCDN(add)).toBeTruthy();
				});
			}
			for (const update of cdnsData.update) {
				it(update.description, async () => {
					await cdnsPage.searchCDN(update.Name);
					expect(await cdnsPage.updateCDN(update)).toBeTruthy();
				});
			}
			for (const remove of cdnsData.remove) {
				it(remove.description, async () => {
					await cdnsPage.searchCDN(remove.Name);
					expect(await cdnsPage.deleteCDN(remove.Name, remove.validationMessage)).toBeTruthy();
				});
			}
		});
	}
});
