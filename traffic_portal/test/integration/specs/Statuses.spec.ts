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

import { LoginPage } from '../PageObjects/LoginPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { API } from '../CommonUtils/API';
import { StatusesPage } from '../PageObjects/Statuses.po'

let setupFile = 'Data/Statuses/Setup.json';
let cleanupFile = 'Data/Statuses/Cleanup.json';
let filename = 'Data/Statuses/TestCases.json';
let testData = JSON.parse(readFileSync(filename, "utf8"));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let statusesPage = new StatusesPage();

describe('Setup API for Statuses Test', function(){
    it('Setup', async function(){
        let setupData = JSON.parse(readFileSync(setupFile, "utf8"));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})
using(testData.Statuses, async function(statusesData){
    using(statusesData.Login, function(login){
        describe('Traffic Portal - Statuses - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open statuses page', async function(){
                await statusesPage.OpenConfigureMenu();
                await statusesPage.OpenStatusesPage();
            })

            using(statusesData.Add, function (add) {
                it(add.description, async function () {
                    expect(await statusesPage.CreateStatus(add)).toBeTruthy();
                    await statusesPage.OpenStatusesPage();
                })
            })
            using(statusesData.Update, function (update) {
                it(update.description, async function () {
                    await statusesPage.SearchStatus(update.Name);
                    expect(await statusesPage.UpdateStatus(update)).toBeTruthy();
                    await statusesPage.OpenStatusesPage();
                })
            })
            using(statusesData.Remove, function (remove) {
                it(remove.description, async function () {
                    await statusesPage.SearchStatus(remove.Name);
                    expect(await statusesPage.DeleteStatus(remove)).toBeTruthy();
                    await statusesPage.OpenStatusesPage();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})
describe('Clean Up API for Statuses Test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(readFileSync(cleanupFile, "utf8"));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})
