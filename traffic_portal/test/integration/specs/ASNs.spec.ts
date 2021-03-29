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
import { ASNsPage } from '../PageObjects/ASNs.po';

let setupFile = 'Data/ASNs/Setup.json';
let cleanupFile = 'Data/ASNs/Cleanup.json';
let filename = 'Data/ASNs/TestCases.json';
let testData = JSON.parse(readFileSync(filename, "utf8"));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let asnsPage = new ASNsPage();

describe('Setup API for ASNs Test', function(){
    it('Setup', async function(){
        let setupData = JSON.parse(readFileSync(setupFile, "utf8"));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})

using(testData.ASNs, async function(asnsData){
    using(asnsData.Login, function(login){
        describe('Traffic Portal - ASNs - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open asns page', async function(){
                await asnsPage.OpenTopologyMenu();
                await asnsPage.OpenASNsPage();
            })

            using(asnsData.Add, function (add) {
                it(add.description, async function () {
                    expect(await asnsPage.CreateASNs(add)).toBeTruthy();
                    await asnsPage.OpenASNsPage();
                })
            })
            using(asnsData.Update, function (update) {
                it(update.description, async function () {
                    await asnsPage.SearchASNs(update.ASNs);
                    expect(await asnsPage.UpdateASNs(update)).toBeTruthy();
                    await asnsPage.OpenASNsPage();
                })
            })
            using(asnsData.Remove, function (remove) {
                it(remove.description, async function () {
                    await asnsPage.SearchASNs(remove.ASNs);
                    expect(await asnsPage.DeleteASNs(remove)).toBeTruthy();
                    await asnsPage.OpenASNsPage();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})

describe('Clean Up API for ASNs Test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(readFileSync(cleanupFile, "utf8"));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})
