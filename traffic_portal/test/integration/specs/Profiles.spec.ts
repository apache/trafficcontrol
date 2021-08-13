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
import { browser } from 'protractor'

import { LoginPage } from '../PageObjects/LoginPage.po'
import { ProfilesPage } from '../PageObjects/ProfilesPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { API } from '../CommonUtils/API';
import { profiles } from "../Data";

const api = new API();
const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const profilesPage = new ProfilesPage();

describe('Setup API for Profiles', () => {
    it('Setup', async () => {
        await api.UseAPI(profiles.setup);
    });
    profiles.tests.forEach(async profilesData => {
        profilesData.logins.forEach(login => {
            describe(`Traffic Portal - Profiles - ${login.description}`, () => {
                afterEach(async function () {
                    await profilesPage.OpenProfilesPage();
                });
                afterAll(async function () {
                    expect(await topNavigation.Logout()).toBeTruthy();
                })
                it('can login', async () => {
                    browser.get(browser.params.baseUrl);
                    await loginPage.Login(login);
                    expect(await loginPage.CheckUserName(login)).toBeTruthy();
                    await profilesPage.OpenConfigureMenu();
                });
                profilesData.check.forEach(check => {
                    it(check.description, async () => {
                        expect(await profilesPage.CheckCSV(check.Name)).toBe(true);
                    });
                });
                profilesData.toggle.forEach(toggle => {
                    it(toggle.description, async () => {
                        if (toggle.description.includes('hide')) {
                            expect(await profilesPage.ToggleTableColumn(toggle.Name)).toBe(false);
                        } else {
                            expect(await profilesPage.ToggleTableColumn(toggle.Name)).toBe(true);
                        }
                    });
                })
                profilesData.add.forEach(add => {
                    it(add.description, async () => {
                        expect(await profilesPage.CreateProfile(add)).toBeTruthy();
                    });
                });
                profilesData.update.forEach(update => {
                    it(update.description, async () => {
                        await profilesPage.SearchProfile(update.Name);
                        expect(await profilesPage.UpdateProfile(update)).toBeTruthy();
                    });
                });
                profilesData.remove.forEach(remove => {
                    it(remove.description, async () => {
                        await profilesPage.SearchProfile(remove.Name);
                        expect(await profilesPage.DeleteProfile(remove)).toBeTruthy();
                    });
                });
            });
        });
    });
    it('Cleanup', async () => {
        await api.UseAPI(profiles.cleanup);
    });
});