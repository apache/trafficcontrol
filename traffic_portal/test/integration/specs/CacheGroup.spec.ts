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
import { CacheGroupPage } from '../PageObjects/CacheGroup.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { cachegroups } from "../Data";



let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let cacheGroupPage = new CacheGroupPage();

cachegroups.tests.forEach(cacheGroupData => {
    describe(`Traffic Portal - CacheGroup - ${cacheGroupData.testName}`, () => {
        cacheGroupData.logins.forEach(login => {
            it('can login', async function () {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open cache group page', async function () {
                await cacheGroupPage.OpenTopologyMenu();
                await cacheGroupPage.OpenCacheGroupsPage();
            })
            cacheGroupData.check.forEach(check => {
                it(check.description, async () => {
                    expect(await cacheGroupPage.CheckCSV(check.Name)).toBe(true);
                    await cacheGroupPage.OpenCacheGroupsPage();
                });
            });
            cacheGroupData.toggle.forEach(toggle => {
                it(toggle.description, async () => {
                    if(toggle.description.includes('hide')){
                        expect(await cacheGroupPage.ToggleTableColumn(toggle.Name)).toBe(false);
                        await cacheGroupPage.OpenCacheGroupsPage();
                    }else{
                        expect(await cacheGroupPage.ToggleTableColumn(toggle.Name)).toBe(true);
                        await cacheGroupPage.OpenCacheGroupsPage();
                    }
                    
                });
            })
            cacheGroupData.create.forEach(create => {
                it(create.Description, async function () {
                    expect(await cacheGroupPage.CreateCacheGroups(create, create.validationMessage)).toBeTruthy();
                    await cacheGroupPage.OpenCacheGroupsPage();
                })
            })
            cacheGroupData.update.forEach(update => {
                if (update.Description.includes("cannot")) {
                    it(update.Description, async function () {
                        await cacheGroupPage.SearchCacheGroups(update.Name)
                        expect(await cacheGroupPage.UpdateCacheGroups(update, update.validationMessage)).toBeUndefined();
                        await cacheGroupPage.OpenCacheGroupsPage();
                    })
                } else {
                    it(update.Description, async function () {
                        await cacheGroupPage.SearchCacheGroups(update.Name)
                        expect(await cacheGroupPage.UpdateCacheGroups(update, update.validationMessage)).toBeTruthy();
                        await cacheGroupPage.OpenCacheGroupsPage();
                    })
                }

            })
            cacheGroupData.remove.forEach(remove => {
                it(remove.Description, async function () {
                    await cacheGroupPage.SearchCacheGroups(remove.Name)
                    expect(await cacheGroupPage.DeleteCacheGroups(remove.Name, remove.validationMessage)).toBeTruthy();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})
