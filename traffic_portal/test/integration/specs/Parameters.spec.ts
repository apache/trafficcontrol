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

import { LoginPage } from '../PageObjects/LoginPage.po'
import { ParametersPage } from '../PageObjects/ParametersPage.po';
import { api } from "../config";
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { parameters } from "../Data";

const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const parametersPage = new ParametersPage();

describe('Setup API for parameter test', () => {
    it('Setup', async () => {
        await api.UseAPI(parameters.setup);
    });
});

parameters.tests.forEach(async parametersData => {
    parametersData.logins.forEach(login => {
        describe(`Traffic Portal - Parameters - ${login.description}`, () => {

            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            });
            it('can open parameters page', async function(){
                await parametersPage.OpenConfigureMenu();
                await parametersPage.OpenParametersPage();
            });
            parametersData.check.forEach(check => {
                it(check.description, async () => {
                    expect(await parametersPage.CheckCSV(check.Name)).toBe(true);
                    await parametersPage.OpenParametersPage();
                });
            });
            parametersData.toggle.forEach(toggle => {
                it(toggle.description, async () => {
                    if(toggle.description.includes('hide')){
                        expect(await parametersPage.ToggleTableColumn(toggle.Name)).toBe(false);
                        await parametersPage.OpenParametersPage();
                    }else{
                        expect(await parametersPage.ToggleTableColumn(toggle.Name)).toBe(true);
                        await parametersPage.OpenParametersPage();
                    }

                });
            })
            parametersData.add.forEach(add => {
                it(add.description, async () => {
                    expect(await parametersPage.CreateParameter(add)).toBeTruthy();
                    await parametersPage.OpenParametersPage();
                });
            });
            parametersData.update.forEach(update => {
                it(update.description, async () => {
                    await parametersPage.SearchParameter(update.Name);
                    expect(await parametersPage.UpdateParameter(update)).toBeTruthy();
                    await parametersPage.OpenParametersPage();
                });
            });

            parametersData.remove.forEach(remove => {
                it(remove.description, async () => {
                    await parametersPage.SearchParameter(remove.Name);
                    expect(await parametersPage.DeleteParameter(remove)).toBeTruthy();
                    await parametersPage.OpenParametersPage();
                });
            });

            it('can logout', async () => {
                expect(await topNavigation.Logout()).toBeTruthy();
            });
        });
    });
});

describe('Clean up API for parameter test', () => {
    it('Cleanup', async () => {
        await api.UseAPI(parameters.cleanup);
    });
});
