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

import { browser } from 'protractor';
import { LoginPage } from '../PageObjects/LoginPage.po';
import { TopologiesPage } from '../PageObjects/TopologiesPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import  * as using  from "jasmine-data-provider";
import { readFileSync } from "fs"
import { API } from '../CommonUtils/API';

const api = new API();
const testFile = 'Data/Topologies/TestCases.json';
const setupFile = 'Data/Topologies/Setup.json';
const cleanupFile = 'Data/Topologies/Cleanup.json';
const testData = JSON.parse(readFileSync(testFile,'utf-8'));
const loginPage = new LoginPage();
const topologiesPage = new TopologiesPage();
const topNavigation = new TopNavigationPage();

describe('Setup prereq for Topologies Test', function(){
    it('Setup', async function(){
        let setupData = JSON.parse(readFileSync(setupFile,'utf-8'));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})

using(testData.Topologies, async function(topologiesData){
    using(topologiesData.Login, function(login){
        describe('Traffic Portal - Topologies - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open topologies page', async function(){
                await topologiesPage.OpenTopologyMenu();
                await topologiesPage.OpenTopologiesPage();
            })
            using(topologiesData.Add, function (add) {
                it(add.description, async function () {
                    expect(await topologiesPage.CreateTopologies(add)).toBeTruthy();
                    await topologiesPage.OpenTopologiesPage();
                })
            })
            it('can logout', async function(){
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        
        })
    })
})

describe('Clean up prereq and test data for Topologies Test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(readFileSync(cleanupFile,'utf-8'));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})