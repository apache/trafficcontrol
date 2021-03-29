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
import { CacheGroupPage } from '../PageObjects/CacheGroup.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';


let filename = 'Data/CacheGroup/TestCases.json';
let testData = JSON.parse(readFileSync(filename, "utf8"));

let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let cacheGroupPage = new CacheGroupPage();

using(testData.CacheGroup, function (cacheGroupData) {
    describe('Traffic Portal - CacheGroup - ' + cacheGroupData.TestName, function () {
        using(cacheGroupData.Login, function (login) {
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open cache group page', async function () {
                await cacheGroupPage.OpenTopologyMenu();
                await cacheGroupPage.OpenCacheGroupsPage();
            })
            using(cacheGroupData.Create, function (create) {
                it(create.description, async function () {
                    expect(await cacheGroupPage.CreateCacheGroups(create, create.validationMessage)).toBeTruthy();
                    await cacheGroupPage.OpenCacheGroupsPage();
                })
            })
            using(cacheGroupData.Update, function (update) {
                if(update.description.includes("cannot")){
                    it(update.description, async function () {
                        await cacheGroupPage.SearchCacheGroups(update.Name)
                        expect(await cacheGroupPage.UpdateCacheGroups(update, update.validationMessage)).toBeUndefined();
                        await cacheGroupPage.OpenCacheGroupsPage();
                    })
                }else{
                    it(update.description, async function () {
                        await cacheGroupPage.SearchCacheGroups(update.Name)
                        expect(await cacheGroupPage.UpdateCacheGroups(update, update.validationMessage)).toBeTruthy();
                        await cacheGroupPage.OpenCacheGroupsPage();
                    })
                }

            })
            using(cacheGroupData.Remove, function (remove) {
                it(remove.description, async function () {
                    await cacheGroupPage.SearchCacheGroups(remove.Name)
                    expect(await cacheGroupPage.DeleteCacheGroups(remove.Name, remove.validationMessage)).toBeTruthy();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})
