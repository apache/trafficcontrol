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
import { TypesPage } from '../PageObjects/Types.po'
import { types } from "../Data";

const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const typesPage = new TypesPage();

describe('Setup API for Types Test', () => {
    it('Setup', async () => {
        await api.UseAPI(types.setup);
    });
});
types.tests.forEach(async typesData => {
    typesData.logins.forEach(login => {
        describe(`Traffic Portal - Types - ${login.description}`, () => {
            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            });
            it('can open types page', async () => {
                await typesPage.OpenConfigureMenu();
                await typesPage.OpenTypesPage();
            });

            typesData.add.forEach(add => {
                it(add.description, async () => {
                    expect(await typesPage.CreateType(add)).toBeTruthy();
                    await typesPage.OpenTypesPage();
                });
            });
            typesData.update.forEach(update => {
                it(update.description, async () => {
                    await typesPage.SearchType(update.Name);
                    expect(await typesPage.UpdateType(update)).toBeTruthy();
                    await typesPage.OpenTypesPage();
                });
            });
            typesData.remove.forEach(remove => {
                it(remove.description, async () => {
                    await typesPage.SearchType(remove.Name);
                    expect(await typesPage.DeleteTypes(remove)).toBeTruthy();
                    await typesPage.OpenTypesPage();
                });
            });
            it('can logout', async () => {
                expect(await topNavigation.Logout()).toBeTruthy();
            });
        });
    });
});
describe('Clean Up API for Types Test', () => {
    it('Cleanup', async () => {
        await api.UseAPI(types.cleanup);
    });
});
