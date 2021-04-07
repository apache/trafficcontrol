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
import { RegionsPage } from '../PageObjects/RegionsPage.po';

let setupFile = 'Data/Regions/Setup.json';
let cleanupFile = 'Data/Regions/Cleanup.json';
let filename = 'Data/Regions/TestCases.json';
let testData = JSON.parse(readFileSync(filename, "utf8"));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let regionsPage = new RegionsPage();

describe('Setup Divisions for Regions Test', function(){
    it('Setup', async function(){
        let setupData = JSON.parse(readFileSync(setupFile, "utf8"));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})

using(testData.Regions, async function(regionsData){
    using(regionsData.Login, function(login){
        describe('Traffic Portal - Regions - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open regions page', async function(){
                await regionsPage.OpenTopologyMenu();
                await regionsPage.OpenRegionsPage();
            })

            using(regionsData.Add, function (add) {
                it(add.description, async function () {
                    expect(await regionsPage.CreateRegions(add)).toBeTruthy();
                    await regionsPage.OpenRegionsPage();
                })
            })
            using(regionsData.Update, function (update) {
                it(update.description, async function () {
                    await regionsPage.SearchRegions(update.Name);
                    expect(await regionsPage.UpdateRegions(update)).toBeTruthy();
                    await regionsPage.OpenRegionsPage();
                })
            })
            using(regionsData.Remove, function (remove) {
                it(remove.description, async function () {
                    await regionsPage.SearchRegions(remove.Name);
                    expect(await regionsPage.DeleteRegions(remove)).toBeTruthy();
                    await regionsPage.OpenRegionsPage();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})

describe('Clean Up Divisions for Regions Test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(readFileSync(cleanupFile, "utf8"));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})
