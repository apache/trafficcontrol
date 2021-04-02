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
import { readFileSync } from "fs";

import { browser } from 'protractor'
import using from "jasmine-data-provider";

import { LoginPage } from '../PageObjects/LoginPage.po'
import { ProfilesPage } from '../PageObjects/ProfilesPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { API } from '../CommonUtils/API';

let setupFile = 'Data/Profiles/Setup.json';
let cleanupFile = 'Data/Profiles/Cleanup.json';
let filename = 'Data/Profiles/TestCases.json';
let testData = JSON.parse(readFileSync(filename, "utf8"));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let profilesPage = new ProfilesPage();

describe('Setup API for Profiles', function () {
    it('Setup', async function () {
        let setupData = JSON.parse(readFileSync(setupFile, "utf8"));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})
using(testData.Profiles, async function(profilesData){
    using(profilesData.Login, function(login){
        describe('Traffic Portal - Profiles - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open profiles page', async function () {
                await profilesPage.OpenConfigureMenu();
                await profilesPage.OpenProfilesPage();
            })
            using(profilesData.Add, function (add) {
                it(add.description, async function () {
                    expect(await profilesPage.CreateProfile(add)).toBeTruthy();
                    await profilesPage.OpenProfilesPage();
                })
            })
            using(profilesData.Update, function (update) {
                it(update.description, async function () {
                    await profilesPage.SearchProfile(update.Name);
                    expect(await profilesPage.UpdateProfile(update)).toBeTruthy();
                    await profilesPage.OpenProfilesPage();
                })
            })
            using(profilesData.Remove, function (remove) {
                it(remove.description, async function () {
                    await profilesPage.SearchProfile(remove.Name);
                    expect(await profilesPage.DeleteProfile(remove)).toBeTruthy();
                    await profilesPage.OpenProfilesPage();
                })
            })
            it('can logout', async function(){
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})
describe('Clean up API for Profiles', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(readFileSync(cleanupFile, "utf8"));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})
