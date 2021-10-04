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
import { TenantsPage } from '../PageObjects/TenantsPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { api } from "../config";
import { tenants } from '../Data/tenant';

const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const tenantsPage = new TenantsPage();

describe('Setup API for Tenants', () => {
    it('Setup', async () => {
        await api.UseAPI(tenants.setup);
    });
});

tenants.tests.forEach(async tenantsData => {
    tenantsData.logins.forEach(login => {
        describe(`Traffic Portal - tenants - ${login.description}`, () => {
            afterEach(async function () {
                await tenantsPage.OpenTenantsPage();
            });
            afterAll(async function () {
                expect(await topNavigation.Logout()).toBe(true);
            })
            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBe(true);
                await tenantsPage.OpenUserAdminMenu();
            });
            tenantsData.add.forEach(add => {
                it(add.description, async () => {
                    expect(await tenantsPage.CreateTenant(add)).toBe(true);
                });
            });
            tenantsData.update.forEach(update => {
                it(update.description, async () => {
                    await tenantsPage.SearchTenant(update.Name);
                    expect(await tenantsPage.UpdateTenant(update)).toBe(true);
                });
            });
            tenantsData.remove.forEach(remove => {
                it(remove.description, async () => {
                    await tenantsPage.SearchTenant(remove.Name);
                    expect(await tenantsPage.DeleteTenant(remove)).toBe(true);
                });
            });
        });
    });
});

describe('Clean Up API for Tenants Test', () => {
    afterAll(async () => {
        await api.UseAPI(tenants.cleanup);
    });
    it('Cleanup', async() => {
      expect(true).toBeTruthy();
    });
});
