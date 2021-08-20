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
import { StatusesPage } from '../PageObjects/Statuses.po'
import { statuses } from "../Data";

const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const statusesPage = new StatusesPage();

describe('Setup API for Statuses Test', () => {
    it('Setup', async () => {
        await api.UseAPI(statuses.setup);
    });
});
statuses.tests.forEach(async statusesData => {
    statusesData.logins.forEach(login => {
        describe(`Traffic Portal - Statuses - ${login.description}`, () => {
            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            });
            it('can open statuses page', async () => {
                await statusesPage.OpenConfigureMenu();
                await statusesPage.OpenStatusesPage();
            });

            statusesData.add.forEach(add => {
                it(add.description, async () => {
                    expect(await statusesPage.CreateStatus(add)).toBeTruthy();
                    await statusesPage.OpenStatusesPage();
                });
            });
            statusesData.update.forEach(update => {
                it(update.description, async () => {
                    await statusesPage.SearchStatus(update.Name);
                    expect(await statusesPage.UpdateStatus(update)).toBeTruthy();
                    await statusesPage.OpenStatusesPage();
                });
            });
            statusesData.remove.forEach(remove => {
                it(remove.description, async () => {
                    await statusesPage.SearchStatus(remove.Name);
                    expect(await statusesPage.DeleteStatus(remove)).toBeTruthy();
                    await statusesPage.OpenStatusesPage();
                });
            });
            it('can logout', async () => {
                expect(await topNavigation.Logout()).toBeTruthy();
            });
        });
    });
});
describe('Clean Up API for Statuses Test', () => {
    it('Cleanup', async () => {
        await api.UseAPI(statuses.cleanup);
    });
});
