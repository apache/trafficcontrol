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

import { browser } from 'protractor';
import using from "jasmine-data-provider";

import { LoginPage } from '../PageObjects/LoginPage.po'
import { OriginsPage } from '../PageObjects/OriginsPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { API } from '../CommonUtils/API';

let setupFile = 'Data/Origins/Setup.json';
let cleanupFile = 'Data/Origins/Cleanup.json';
let filename = 'Data/Origins/TestCases.json';
let testData = JSON.parse(readFileSync(filename, "utf8"));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let originsPage = new OriginsPage();

describe('Setup Origin Delivery Service', function () {
    it('Setup', async function () {
        let setupData = JSON.parse(readFileSync(setupFile, "utf8"));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})
using(testData.Origins, async function (originsData) {
    using(originsData.Login, function (login) {
        describe('Traffic Portal - Origins - ' + login.description, function () {
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open origins page', async function () {
                await originsPage.OpenConfigureMenu();
                await originsPage.OpenOriginsPage();
            })
            using(originsData.Add, function (add) {
                it(add.description, async function () {
                    expect(await originsPage.CreateOrigins(add)).toBeTruthy();
                    await originsPage.OpenOriginsPage();
                })
            })
            using(originsData.Update, function (update) {
                if (update.validationMessage == undefined) {
                    it(update.description, async function () {
                        await originsPage.SearchOrigins(update.Name);
                        expect(await originsPage.UpdateOrigins(update)).toBeUndefined();
                        await originsPage.OpenOriginsPage();
                    })
                } else {
                    it(update.description, async function () {
                        await originsPage.SearchOrigins(update.Name);
                        expect(await originsPage.UpdateOrigins(update)).toBeTruthy();
                        await originsPage.OpenOriginsPage();
                    })
                }
            })
            using(originsData.Remove, function (remove) {
                it(remove.description, async function () {
                    await originsPage.SearchOrigins(remove.Name);
                    expect(await originsPage.DeleteOrigins(remove)).toBeTruthy();
                    await originsPage.OpenOriginsPage();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})

describe('Clean up Origin Delivery Service', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(readFileSync(cleanupFile, "utf8"));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})
