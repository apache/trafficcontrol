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
import { PhysLocationsPage } from '../PageObjects/PhysLocationsPage.po';
import { api } from "../config";
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { physLocations } from "../Data";

const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const physlocationsPage = new PhysLocationsPage();

describe('Setup API for physlocation test', () => {
    it('Setup', async () => {
        await api.UseAPI(physLocations.setup);
    })
})

physLocations.tests.forEach(async physlocationsData => {
    physlocationsData.logins.forEach(login => {
        describe(`Traffic Portal - PhysLocation - ${login.description}`, () => {

            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            });
            it('can open physical locations page', async () => {
                await physlocationsPage.OpenConfigureMenu();
                await physlocationsPage.OpenPhysLocationPage();
            });
            physlocationsData.check.forEach(check => {
                it(check.description, async () => {
                    expect(await physlocationsPage.CheckCSV(check.Name)).toBe(true);
                    await physlocationsPage.OpenPhysLocationPage();
                });
            });
            physlocationsData.add.forEach(add => {
                it(add.description, async () => {
                    expect(await physlocationsPage.CreatePhysLocation(add)).toBeTruthy();
                    await physlocationsPage.OpenPhysLocationPage();
                });
            });
            physlocationsData.update.forEach(update => {
                it(update.description, async () => {
                    await physlocationsPage.SearchPhysLocation(update.Name);
                    expect(await physlocationsPage.UpdatePhysLocation(update)).toBeTruthy();
                    await physlocationsPage.OpenPhysLocationPage();
                });
            });
            physlocationsData.remove.forEach(remove => {
                it(remove.description, async () => {
                    await physlocationsPage.SearchPhysLocation(remove.Name);
                    expect(await physlocationsPage.DeletePhysLocation(remove)).toBeTruthy();
                    await physlocationsPage.OpenPhysLocationPage();
                });
            });
            it('can logout', async () => {
                expect(await topNavigation.Logout()).toBeTruthy();
            });
        });
    });
});

describe('Clean up API for physlocation test', () => {
    it('Cleanup', async () => {
        await api.UseAPI(physLocations.cleanup);
    });
});
