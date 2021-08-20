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
import { ASNsPage } from '../PageObjects/ASNs.po';
import { ASNs } from "../Data";

const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const asnsPage = new ASNsPage();

describe('Setup API for ASNs Test', () => {
    it('Setup', () => {
        api.UseAPI(ASNs.setup);
    });
});

ASNs.tests.forEach(async asnsData => {
    asnsData.logins.forEach( login => {
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

            asnsData.add.forEach( add => {
                it(add.description, async function () {
                    expect(await asnsPage.CreateASNs(add)).toBeTruthy();
                    await asnsPage.OpenASNsPage();
                })
            });
            asnsData.update.forEach( update => {
                it(update.description, async function () {
                    await asnsPage.SearchASNs(update.ASNs);
                    expect(await asnsPage.UpdateASNs(update)).toBeTruthy();
                    await asnsPage.OpenASNsPage();
                })
            });
            asnsData.remove.forEach( remove => {
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

describe('Clean Up API for ASNs Test', () => {
    it('Cleanup', () => {
        api.UseAPI(ASNs.cleanup);
    })
})
