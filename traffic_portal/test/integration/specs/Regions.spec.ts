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

import { LoginPage } from '../PageObjects/LoginPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { api } from "../config";
import { RegionsPage } from '../PageObjects/RegionsPage.po';
import { regions } from "../Data";

const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const regionsPage = new RegionsPage();

describe('Setup Divisions for Regions Test', () => {
    it('Setup', async () => {
        await api.UseAPI(regions.setup);
    });
});

regions.tests.forEach(async regionsData => {
    regionsData.logins.forEach(login => {
        describe(`Traffic Portal - Regions - ${login.description}`, () => {
            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            });
            it('can open regions page', async () => {
                await regionsPage.OpenTopologyMenu();
                await regionsPage.OpenRegionsPage();
            });
            regionsData.check.forEach(check => {
                it(check.description, async () => {
                    expect(await regionsPage.CheckCSV(check.Name)).toBe(true);
                    await regionsPage.OpenRegionsPage();
                });
            });
            regionsData.add.forEach(add => {
                it(add.description, async () => {
                    expect(await regionsPage.CreateRegions(add)).toBeTruthy();
                    await regionsPage.OpenRegionsPage();
                });
            });
            regionsData.update.forEach(update => {
                it(update.description, async () => {
                    await regionsPage.SearchRegions(update.Name);
                    expect(await regionsPage.UpdateRegions(update)).toBeTruthy();
                    await regionsPage.OpenRegionsPage();
                });
            });
            regionsData.remove.forEach(remove => {
                it(remove.description, async () => {
                    await regionsPage.SearchRegions(remove.Name);
                    expect(await regionsPage.DeleteRegions(remove)).toBeTruthy();
                    await regionsPage.OpenRegionsPage();
                });
            });
            it('can logout', async () => {
                expect(await topNavigation.Logout()).toBeTruthy();
            });
        });
    });
});

describe('Clean Up Divisions for Regions Test', () => {
    it('Cleanup', async () => {
        await api.UseAPI(regions.cleanup);
    });
});
