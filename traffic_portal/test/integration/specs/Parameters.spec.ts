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
import { ParametersPage } from '../PageObjects/ParametersPage.po';
import { API } from '../CommonUtils/API';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let parametersPage = new ParametersPage();


let setupFile = 'Data/Parameters/Setup.json';
let cleanupFile = 'Data/Parameters/Cleanup.json';
let filename = 'Data/Parameters/TestCases.json';
let testData = JSON.parse(readFileSync(filename, "utf8"));

describe('Setup API for parameter test', function () {
    it('Setup', async function () {
        let setupData = JSON.parse(readFileSync(setupFile, "utf8"));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})

using(testData.Parameters, async function(parametersData){
    using(parametersData.Login, function(login){
        describe('Traffic Portal - Parameters - ' + login.description, function(){

            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open parameters page', async function(){
                await parametersPage.OpenConfigureMenu();
                await parametersPage.OpenParametersPage();
            })
            using(parametersData.Add, function (add) {
                it(add.description, async function () {
                    expect(await parametersPage.CreateParameter(add)).toBeTruthy();
                    await parametersPage.OpenParametersPage();
                })
            })
            using(parametersData.Update, function (update) {
                it(update.description, async function () {
                    await parametersPage.SearchParameter(update.Name);
                    expect(await parametersPage.UpdateParameter(update)).toBeTruthy();
                    await parametersPage.OpenParametersPage();
                })
            })

            using(parametersData.Remove, function (remove) {
                it(remove.description, async function () {
                    await parametersPage.SearchParameter(remove.Name);
                    expect(await parametersPage.DeleteParameter(remove)).toBeTruthy();
                    await parametersPage.OpenParametersPage();
                })
            })

            it('can logout', async function(){
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})

describe('Clean up API for parameter test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(readFileSync(cleanupFile, "utf8"));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})
