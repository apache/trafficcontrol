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
import { ServersPage } from '../PageObjects/ServersPage.po';
import { api } from "../config";
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { servers } from "../Data";

const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const serversPage = new ServersPage();

describe('Setup API call for Servers Test', () => {
    it('Setup', async () => {
        await api.UseAPI(servers.setup);
    });
});
servers.tests.forEach(async serversData => {
    serversData.logins.forEach(login => {
        describe(`Traffic Portal - Servers - ${login.description}`, () => {
            afterEach(async function () {
                await serversPage.OpenServerPage();
            });
            afterAll(async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBe(true);
                await serversPage.OpenConfigureMenu();
            });
            serversData.toggle.forEach(toggle => {
                it(toggle.description, async () => {
                    if (toggle.description.includes('hide')) {
                        expect(await serversPage.ToggleTableColumn(toggle.Name)).toBe(true);
                    } else {
                        expect(await serversPage.ToggleTableColumn(toggle.Name)).toBe(true);
                    }
                });
            })
            serversData.add.forEach(add => {
                it(add.description, async () => {
                    expect(await serversPage.CreateServer(add)).toBe(true);
                });
            });
            serversData.update.forEach(update => {
                it(update.description, async () => {
                    await serversPage.SearchServer(update.Name);
                    expect(await serversPage.UpdateServer(update)).toBe(true);
                });
            });
            serversData.remove.forEach(remove => {
                it(remove.description, async () => {
                    await serversPage.SearchServer(remove.Name);
                    expect(await serversPage.DeleteServer(remove)).toBe(true);
                });
            });
        })
    })
})
describe('API Clean Up for Servers Test', () => {
    it('Cleanup', async () => {
        await api.UseAPI(servers.cleanup);
    });
});
