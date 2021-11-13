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
import { OriginsPage } from '../PageObjects/OriginsPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { api } from "../config";
import { origins } from "../Data";

const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const originsPage = new OriginsPage();

describe('Setup Origin Delivery Service', function () {
    it('Setup', async function () {
        await api.UseAPI(origins.setup);
    });
});
origins.tests.forEach(async originsData => {
    originsData.logins.forEach(login => {
        describe(`Traffic Portal - Origins - ${login.description}`, () => {
            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            });
            it('can open origins page', async () => {
                await originsPage.OpenConfigureMenu();
                await originsPage.OpenOriginsPage();
            });
            originsData.add.forEach(add => {
                it(add.description, async function () {
                    expect(await originsPage.CreateOrigins(add)).toBeTruthy();
                    await originsPage.OpenOriginsPage();
                });
            });
            originsData.update.forEach(update => {
                if (!update.validationMessage) {
                    it(update.description, async () => {
                        await originsPage.SearchOrigins(update.Name);
                        expect(await originsPage.UpdateOrigins(update)).toBeUndefined();
                        await originsPage.OpenOriginsPage();
                    });
                } else {
                    it(update.description, async () => {
                        await originsPage.SearchOrigins(update.Name);
                        expect(await originsPage.UpdateOrigins(update)).toBeTruthy();
                        await originsPage.OpenOriginsPage();
                    });
                }
            });
            originsData.remove.forEach(remove => {
                it(remove.description, async () => {
                    await originsPage.SearchOrigins(remove.Name);
                    expect(await originsPage.DeleteOrigins(remove)).toBeTruthy();
                    await originsPage.OpenOriginsPage();
                });
            });
            it('can logout', async () => {
                expect(await topNavigation.Logout()).toBeTruthy();
            });
        });
    });
});

describe('Clean up Origin Delivery Service', () => {
    it('Cleanup', async () => {
        await api.UseAPI(origins.cleanup);
    })
})
