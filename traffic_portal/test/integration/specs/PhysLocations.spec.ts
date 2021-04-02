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
import { PhysLocationsPage } from '../PageObjects/PhysLocationsPage.po';
import { API } from '../CommonUtils/API';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let physlocationsPage = new PhysLocationsPage();

let setupFile = 'Data/PhysLocations/Setup.json';
let cleanupFile = 'Data/PhysLocations/Cleanup.json';
let filename = 'Data/PhysLocations/TestCases.json';
let testData = JSON.parse(readFileSync(filename, "utf8"));

describe('Setup API for physlocation test', function () {
    it('Setup', async function () {
        let setupData = JSON.parse(readFileSync(setupFile, "utf8"));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})

using(testData.PhysLocations, async function(physlocationsData){
    using(physlocationsData.Login, function(login){
        describe('Traffic Portal - PhysLocation - ' + login.description, function(){

            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open parameters page', async function(){
                await physlocationsPage.OpenConfigureMenu();
                await physlocationsPage.OpenPhysLocationPage();
            })
            using(physlocationsData.Add, function (add) {
                it(add.description, async function () {
                    expect(await physlocationsPage.CreatePhysLocation(add)).toBeTruthy();
                    await physlocationsPage.OpenPhysLocationPage();
                })
            })
            using(physlocationsData.Update, function (update) {
                it(update.description, async function () {
                    await physlocationsPage.SearchPhysLocation(update.Name);
                    expect(await physlocationsPage.UpdatePhysLocation(update)).toBeTruthy();
                    await physlocationsPage.OpenPhysLocationPage();
                })
            })

            using(physlocationsData.Remove, function (remove) {
                it(remove.description, async function () {
                    await physlocationsPage.SearchPhysLocation(remove.Name);
                    expect(await physlocationsPage.DeletePhysLocation(remove)).toBeTruthy();
                    await physlocationsPage.OpenPhysLocationPage();
                })
            })
            it('can logout', async function(){
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})

describe('Clean up API for physlocation test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(readFileSync(cleanupFile, "utf8"));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})
