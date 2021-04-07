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
import { ServerCapabilitiesPage } from '../PageObjects/ServerCapabilitiesPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';

let filename = 'Data/ServerCapabilities/TestCases.json';
let testData = JSON.parse(readFileSync(filename, "utf8"));

let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let serverCapabilitiesPage = new ServerCapabilitiesPage();

using(testData.ServerCapabilities, function(serverCapabilitiesData) {
    describe('Traffic Portal - Server Capabilities - '+ serverCapabilitiesData.TestName,  function(){
        using(serverCapabilitiesData.Login, function(login) {
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open server capability page', async function() {
                await serverCapabilitiesPage.OpenConfigureMenu();
                await serverCapabilitiesPage.OpenServerCapabilityPage();
            })
            using(serverCapabilitiesData.Add, function(add) {
                it(add.description, async function(){
                    expect(await serverCapabilitiesPage.CreateServerCapabilities(add.Name, add.validationMessage)).toBeTruthy();
                    await serverCapabilitiesPage.OpenServerCapabilityPage();
                })
            })
            using(serverCapabilitiesData.Delete, function(remove) {
                if(remove.description.includes("invalid")){
                    it(remove.description, async function(){
                        await serverCapabilitiesPage.SearchServerCapabilities(remove.Name)
                        expect(await serverCapabilitiesPage.DeleteServerCapabilities(remove.InvalidName, remove.validationMessage)).toBeFalsy();
                        await serverCapabilitiesPage.OpenServerCapabilityPage();
                    })
                } else {
                    it(remove.description, async function(){
                        await serverCapabilitiesPage.SearchServerCapabilities(remove.Name)
                        expect(await serverCapabilitiesPage.DeleteServerCapabilities(remove.Name, remove.validationMessage)).toBeTruthy();
                        await serverCapabilitiesPage.OpenServerCapabilityPage();
                    })
                }
            })
            it('can logout', async function(){
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})
