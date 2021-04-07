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
import { ServiceCategoriesPage } from '../PageObjects/ServiceCategories.po';

let setupFile = 'Data/ServiceCategories/Setup.json';
let cleanupFile = 'Data/ServiceCategories/Cleanup.json';
let filename = 'Data/ServiceCategories/TestCases.json';
let testData = JSON.parse(readFileSync(filename, "utf8"));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let serviceCategoriesPage = new ServiceCategoriesPage();

describe('Setup API for Service Categories Test', function(){
    it('Setup', async function(){
        let setupData = JSON.parse(readFileSync(setupFile, "utf8"));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})
using(testData.ServiceCategories, async function(serviceCategoriesData){
    using(serviceCategoriesData.Login, function(login){
        describe('Traffic Portal - ServiceCategories - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open service categories page', async function(){
                await serviceCategoriesPage.OpenServicesMenu();
                await serviceCategoriesPage.OpenServiceCategoriesPage();
            })

            using(serviceCategoriesData.Add, function (add) {
                it(add.description, async function () {
                    expect(await serviceCategoriesPage.CreateServiceCategories(add)).toBeTruthy();
                    await serviceCategoriesPage.OpenServiceCategoriesPage();
                })
            })
            using(serviceCategoriesData.Update, function (update) {
                it(update.description, async function () {
                    await serviceCategoriesPage.SearchServiceCategories(update.Name);
                    expect(await serviceCategoriesPage.UpdateServiceCategories(update)).toBeTruthy();
                    await serviceCategoriesPage.OpenServiceCategoriesPage();
                })
            })
            using(serviceCategoriesData.Remove, function (remove) {
                it(remove.description, async function () {
                    await serviceCategoriesPage.SearchServiceCategories(remove.Name);
                    expect(await serviceCategoriesPage.DeleteServiceCategories(remove)).toBeTruthy();
                    await serviceCategoriesPage.OpenServiceCategoriesPage();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})

describe('Clean Up API for Service Categories Test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(readFileSync(cleanupFile, "utf8"));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})
