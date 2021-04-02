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
import { CDNPage } from '../PageObjects/CDNPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';

let filename = 'Data/CDN/TestCases.json';
let testData = JSON.parse(readFileSync(filename, "utf8"));

let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let cdnsPage = new CDNPage();

using(testData.CDN, async function(cdnsData){
    using(cdnsData.Login, function(login){
        describe('Traffic Portal - CDN - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open CDN page', async function(){
                await cdnsPage.OpenCDNsPage();
            })

            using(cdnsData.Add, function (add){
                it(add.description, async function(){
                    expect(await cdnsPage.CreateCDN(add)).toBeTruthy();
                    await cdnsPage.OpenCDNsPage();
                })
            })
            using(cdnsData.Update, function(update){
                it(update.description, async function(){
                    await cdnsPage.SearchCDN(update.Name);
                    expect(await cdnsPage.UpdateCDN(update)).toBeTruthy();
                })

            })
            using(cdnsData.Remove, function(remove){
                it(remove.description, async function(){
                    await cdnsPage.SearchCDN(remove.Name);
                    expect(await cdnsPage.DeleteCDN(remove)).toBeTruthy();

                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })

        })
    })
})
