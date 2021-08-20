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
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { api } from "../config";
import { TopologiesPage } from '../PageObjects/TopologiesPage.po';
import { topologies } from "../Data/topologies";

const loginPage = new LoginPage();
const topologiesPage = new TopologiesPage();
const topNavigation = new TopNavigationPage();

describe('Setup API for Topologies Test', () => {
    it('Setup', async () => {
        await api.UseAPI(topologies.setup);
    });
});

topologies.tests.forEach(async  topologiesData =>{
    topologiesData.logins.forEach(login =>{
        describe(`Traffic Portal - Topologies - ${login.description}`, () => {
            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            });
            it('can open topologies page', async () => {
                await topologiesPage.OpenTopologyMenu();
                await topologiesPage.OpenTopologiesPage();
            });
            topologiesData.add.forEach(add => {
                it(add.description, async () => {
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

describe('Clean Up API for Topologies Test', () => {
    it('Cleanup', async () => {
        await api.UseAPI(topologies.cleanup);
    });
});
