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
import { DivisionsPage } from '../PageObjects/Divisions.po';
import { divisions } from "../Data";


const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const divisionsPage = new DivisionsPage();

describe('Setup API for Divisions Test', () => {
    it('Setup', async () => {
        await api.UseAPI(divisions.setup);
    });
});

divisions.tests.forEach(divisionsData => {
    divisionsData.logins.forEach(login => {
        describe(`Traffic Portal - Divisions - ${login.description}`, () => {
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            });
            it('can open divisions page', async () => {
                await divisionsPage.OpenTopologyMenu();
                await divisionsPage.OpenDivisionsPage();
            });
            divisionsData.check.forEach(check => {
                it(check.description, async () => {
                    expect(await divisionsPage.CheckCSV(check.Name)).toBe(true);
                    await divisionsPage.OpenDivisionsPage();
                });
            });
            divisionsData.add.forEach(add => {
                it(add.description, async () => {
                    expect(await divisionsPage.CreateDivisions(add)).toBeTruthy();
                    await divisionsPage.OpenDivisionsPage();
                });
            });
            divisionsData.update.forEach(update => {
                it(update.description, async () => {
                    await divisionsPage.SearchDivisions(update.Name);
                    expect(await divisionsPage.UpdateDivisions(update)).toBeTruthy();
                    await divisionsPage.OpenDivisionsPage();
                });
            });
            divisionsData.remove.forEach(remove => {
                it(remove.description, async () => {
                    await divisionsPage.SearchDivisions(remove.Name);
                    expect(await divisionsPage.DeleteDivisions(remove)).toBeTruthy();
                    await divisionsPage.OpenDivisionsPage();
                });
            });
            it('can logout', async () => {
                expect(await topNavigation.Logout()).toBeTruthy();
            });
        });
    });
});

describe('Clean Up API for Divisions Test', () => {
    it('Cleanup', async () => {
        await api.UseAPI(divisions.cleanup);
    });
});
