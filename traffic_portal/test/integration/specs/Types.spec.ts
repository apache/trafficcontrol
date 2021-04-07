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
import { TypesPage } from '../PageObjects/Types.po'

let setupFile = 'Data/Types/Setup.json';
let cleanupFile = 'Data/Types/Cleanup.json';
let filename = 'Data/Types/TestCases.json';
let testData = JSON.parse(readFileSync(filename, "utf8"));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let typesPage = new TypesPage();

describe('Setup API for Types Test', function(){
    it('Setup', async function(){
        let setupData = JSON.parse(readFileSync(setupFile, "utf8"));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})
using(testData.Types, async function(typesData){
    using(typesData.Login, function(login){
        describe('Traffic Portal - Types - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open types page', async function(){
                await typesPage.OpenConfigureMenu();
                await typesPage.OpenTypesPage();
            })

            using(typesData.Add, function (add) {
                it(add.description, async function () {
                    expect(await typesPage.CreateType(add)).toBeTruthy();
                    await typesPage.OpenTypesPage();
                })
            })
            using(typesData.Update, function (update) {
                it(update.description, async function () {
                    await typesPage.SearchType(update.Name);
                    expect(await typesPage.UpdateType(update)).toBeTruthy();
                    await typesPage.OpenTypesPage();
                })
            })
            using(typesData.Remove, function (remove) {
                it(remove.description, async function () {
                    await typesPage.SearchType(remove.Name);
                    expect(await typesPage.DeleteTypes(remove)).toBeTruthy();
                    await typesPage.OpenTypesPage();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})
describe('Clean Up API for Types Test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(readFileSync(cleanupFile, "utf8"));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})
