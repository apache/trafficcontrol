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
import { browser } from 'protractor';

import { LoginPage } from '../PageObjects/LoginPage.po'
import { CDNPage } from '../PageObjects/CDNPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { cdns } from "../Data";


const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const cdnsPage = new CDNPage();

cdns.tests.forEach(async cdnsData =>{
    cdnsData.logins.forEach(login => {
        describe('Traffic Portal - CDN - ' + login.description, function(){
            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            });
            it('can open CDN page', async () => {
                await cdnsPage.OpenCDNsPage();
            });
             cdnsData.check.forEach(check => {
                 it(check.description, async () => {
                     expect(await cdnsPage.CheckCSV(check.Name)).toBe(true);
                     await cdnsPage.OpenCDNsPage();
                 });
            });
            cdnsData.add.forEach(add => {
                it(add.description, async () => {
                    expect(await cdnsPage.CreateCDN(add)).toBeTruthy();
                    await cdnsPage.OpenCDNsPage();
                });
            });
            cdnsData.update.forEach(update => {
                it(update.description, async () => {
                    await cdnsPage.SearchCDN(update.Name);
                    expect(await cdnsPage.UpdateCDN(update)).toBeTruthy();
                });
            });
            cdnsData.remove.forEach(remove => {
                it(remove.description, async () => {
                    await cdnsPage.SearchCDN(remove.Name);
                    expect(await cdnsPage.DeleteCDN(remove)).toBeTruthy();

                });
            });
            it('can logout', async () => {
                expect(await topNavigation.Logout()).toBeTruthy();
            });
        });
    });
});
