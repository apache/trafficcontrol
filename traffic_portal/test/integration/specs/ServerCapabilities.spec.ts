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
import { ServerCapabilitiesPage } from '../PageObjects/ServerCapabilitiesPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { serverCapabilities } from "../Data";

const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const serverCapabilitiesPage = new ServerCapabilitiesPage();

serverCapabilities.tests.forEach(serverCapabilitiesData => {
    describe(`Traffic Portal - Server Capabilities - ${serverCapabilitiesData.testName}`,  () => {
        serverCapabilitiesData.logins.forEach(login => {
            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            });
            it('can open server capability page', async () => {
                await serverCapabilitiesPage.OpenConfigureMenu();
                await serverCapabilitiesPage.OpenServerCapabilityPage();
            });
            serverCapabilitiesData.check.forEach(check => {
                it(check.description, async () => {
                    expect(await serverCapabilitiesPage.CheckCSV(check.Name)).toBe(true);
                    await serverCapabilitiesPage.OpenServerCapabilityPage();
                });
            });
            serverCapabilitiesData.add.forEach(add => {
                it(add.description, async () => {
                    expect(await serverCapabilitiesPage.CreateServerCapabilities(add.name, add.capabilityDescription, add.validationMessage)).toBe(true);
                    await serverCapabilitiesPage.OpenServerCapabilityPage();
                });
            });
            serverCapabilitiesData.remove.forEach(remove => {
                if (remove.invalid) {
                    it(remove.description, async () => {
                        await serverCapabilitiesPage.SearchServerCapabilities(remove.name)
                        expect(await serverCapabilitiesPage.DeleteServerCapabilities(remove.invalidName, remove.validationMessage)).toBe(false);
                        await serverCapabilitiesPage.OpenServerCapabilityPage();
                    });
                } else {
                    it(remove.description, async () => {
                        await serverCapabilitiesPage.SearchServerCapabilities(remove.name)
                        expect(await serverCapabilitiesPage.DeleteServerCapabilities(remove.name, remove.validationMessage)).toBe(true);
                        await serverCapabilitiesPage.OpenServerCapabilityPage();
                    });
                }
            });
            it('can logout', async () => {
                expect(await topNavigation.Logout()).toBeTruthy();
            });
        });
    });
});
