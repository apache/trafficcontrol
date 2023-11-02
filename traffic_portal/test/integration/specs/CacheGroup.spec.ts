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

import { LoginPage } from "../PageObjects/LoginPage.po";
import { CacheGroupPage } from "../PageObjects/CacheGroup.po";
import { TopNavigationPage } from "../PageObjects/TopNavigationPage.po";
import { cachegroups } from "../Data";

let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let cacheGroupPage = new CacheGroupPage();

cachegroups.tests.forEach((cacheGroupData) => {
    for (const login of cacheGroupData.logins) {
        describe(`Traffic Portal - CacheGroup - ${cacheGroupData.testName}`, () => {
            beforeAll(async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
                await cacheGroupPage.OpenTopologyMenu();
                await cacheGroupPage.OpenCacheGroupsPage();
            });
            afterAll(async () => {
                expect(await topNavigation.Logout()).toBeTruthy();
            });
            afterEach(async () => {
                await cacheGroupPage.OpenCacheGroupsPage();
            });
            for (const create of cacheGroupData.create) {
                it(create.Description, async () => {
                    expect(
                        await cacheGroupPage.CreateCacheGroups(
                            create,
                            create.validationMessage
                        )
                    ).toBeTruthy();
                });
            }
            for (const update of cacheGroupData.update) {
                if (update.Description.includes("cannot")) {
                    it(update.Description, async () => {
                        await cacheGroupPage.SearchCacheGroups(update.Name);
                        expect(
                            await cacheGroupPage.UpdateCacheGroups(
                                update,
                                update.validationMessage
                            )
                        ).toBeUndefined();
                    });
                } else {
                    it(update.Description, async () => {
                        await cacheGroupPage.SearchCacheGroups(update.Name);
                        expect(
                            await cacheGroupPage.UpdateCacheGroups(
                                update,
                                update.validationMessage
                            )
                        ).toBeTruthy();
                    });
                }
            }
            for (const remove of cacheGroupData.remove) {
                it(remove.Description, async () => {
                    await cacheGroupPage.SearchCacheGroups(remove.Name);
                    expect(
                        await cacheGroupPage.DeleteCacheGroups(
                            remove.Name,
                            remove.validationMessage
                        )
                    ).toBeTruthy();
                });
            }
        });
    }
});
