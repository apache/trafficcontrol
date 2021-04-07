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
import { DivisionsPage } from '../PageObjects/Divisions.po';

let setupFile = 'Data/Divisions/Setup.json';
let cleanupFile = 'Data/Divisions/Cleanup.json';
let filename = 'Data/Divisions/TestCases.json';
let testData = JSON.parse(readFileSync(filename, "utf8"));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let divisionsPage = new DivisionsPage();

describe('Setup API for Divisions Test', function(){
    it('Setup', async function(){
        let setupData = JSON.parse(readFileSync(setupFile, "utf8"));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})

using(testData.Divisions, async function(divisionsData){
    using(divisionsData.Login, function(login){
        describe('Traffic Portal - Divisions - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open divisions page', async function(){
                await divisionsPage.OpenTopologyMenu();
                await divisionsPage.OpenDivisionsPage();
            })

            using(divisionsData.Add, function (add) {
                it(add.description, async function () {
                    expect(await divisionsPage.CreateDivisions(add)).toBeTruthy();
                    await divisionsPage.OpenDivisionsPage();
                })
            })
            using(divisionsData.Update, function (update) {
                it(update.description, async function () {
                    await divisionsPage.SearchDivisions(update.Name);
                    expect(await divisionsPage.UpdateDivisions(update)).toBeTruthy();
                    await divisionsPage.OpenDivisionsPage();
                })
            })
            using(divisionsData.Remove, function (remove) {
                it(remove.description, async function () {
                    await divisionsPage.SearchDivisions(remove.Name);
                    expect(await divisionsPage.DeleteDivisions(remove)).toBeTruthy();
                    await divisionsPage.OpenDivisionsPage();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})

describe('Clean Up API for Divisions Test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(readFileSync(cleanupFile, "utf8"));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})
