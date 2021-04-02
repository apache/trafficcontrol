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
import { ServersPage } from '../PageObjects/ServersPage.po';
import { API } from '../CommonUtils/API';

let setupFile = 'Data/ServerServerCapabilities/Setup.json';
let cleanupFile = 'Data/ServerServerCapabilities/Cleanup.json';
let filename = 'Data/ServerServerCapabilities/TestCases.json';
let testData = JSON.parse(readFileSync(filename, "utf8"));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let serverCapabilitiesPage = new ServerCapabilitiesPage();
let serverPage = new  ServersPage();

describe("Setup Server Capabilities and Server for prereq", function(){
    it('Setup', async function(){
        let setupData = JSON.parse(readFileSync(setupFile, "utf8"));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})
using(testData.ServerServerCapabilities, async function(serverServerCapData){
    using(serverServerCapData.Login, function(login){
        describe('Traffic Portal - Server Server Capabilities - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open server page', async function(){
                await serverPage.OpenConfigureMenu();
                await serverPage.OpenServerPage();
            })
            using(serverServerCapData.Link, function(link){
                if(link.description.includes("cannot")){
                    it(link.description, async function(){
                        await serverPage.SearchServer(link.Server);
                        expect(await serverPage.AddServerCapabilitiesToServer(link)).toBeUndefined();
                        await serverPage.OpenServerPage();
                    })
                }else{
                    it(link.description, async function(){
                        await serverPage.SearchServer(link.Server);
                        expect(await serverPage.AddServerCapabilitiesToServer(link)).toBeTruthy();
                        await serverPage.OpenServerPage();
                    })
                }
            })
            using(serverServerCapData.Remove, function(remove){
                it(remove.description, async function(){
                    await serverPage.SearchServer(remove.Server);
                    expect(await serverPage.RemoveServerCapabilitiesFromServer(remove.ServerCapability, remove.validationMessage)).toBeTruthy();
                    await serverPage.OpenServerPage();
                })
            })
            it('can open server capabilities page', async function(){
                await serverCapabilitiesPage.OpenServerCapabilityPage();
            })
            using(serverServerCapData.DeleteServerCapability, function(deleteSC){
                it(deleteSC.description, async function(){
                    await serverCapabilitiesPage.SearchServerCapabilities(deleteSC.ServerCapability);
                    expect(await serverCapabilitiesPage.DeleteServerCapabilities(deleteSC.ServerCapability, deleteSC.validationMessage)).toBeTruthy();
                    await serverCapabilitiesPage.OpenServerCapabilityPage();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})
describe("Clean up prereq", function(){
    it('Clean up', async function(){
        let cleanupData = JSON.parse(readFileSync(cleanupFile, "utf8"));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})
