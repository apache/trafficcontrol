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
import { ServersPage } from '../PageObjects/ServersPage.po';
import { api } from "../config";
import { serverServerCapabilities } from "../Data";

const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const serverCapabilitiesPage = new ServerCapabilitiesPage();
const serverPage = new  ServersPage();

describe("Setup Server Capabilities and Server for prereq", () => {
    it('Setup', async () => {
        await api.UseAPI(serverServerCapabilities.setup);
    });
});
serverServerCapabilities.tests.forEach(async serverServerCapData => {
    serverServerCapData.logins.forEach(login => {
        describe(`Traffic Portal - Server Server Capabilities - ${login.description}`, () => {
            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBe(true);
            });
            it('can open server page', async () => {
                await serverPage.OpenConfigureMenu();
                await serverPage.OpenServerPage();
            });
            serverServerCapData.link.forEach(link => {
                if(link.description.includes("cannot")){
                    it(link.description, async () => {
                        await serverPage.SearchServer(link.Server);
                        expect(await serverPage.AddServerCapabilitiesToServer(link)).toBe(false);
                        await serverPage.OpenServerPage();
                    });
                } else {
                    it(link.description, async () => {
                        await serverPage.SearchServer(link.Server);
                        expect(await serverPage.AddServerCapabilitiesToServer(link)).toBe(true);
                        await serverPage.OpenServerPage();
                    });
                }
            });
            serverServerCapData.remove.forEach(remove => {
                it(remove.description, async () => {
                    await serverPage.SearchServer(remove.Server);
                    expect(await serverPage.RemoveServerCapabilitiesFromServer(remove.ServerCapability, remove.validationMessage)).toBe(true);
                    await serverPage.OpenServerPage();
                });
            });
            it('can open server capabilities page', async () => {
                await serverCapabilitiesPage.OpenServerCapabilityPage();
            });
            serverServerCapData.deleteServerCapability.forEach(deleteSC => {
                it(deleteSC.description, async () => {
                    await serverCapabilitiesPage.SearchServerCapabilities(deleteSC.ServerCapability);
                    expect(await serverCapabilitiesPage.DeleteServerCapabilities(deleteSC.ServerCapability, deleteSC.validationMessage)).toBe(true);
                    await serverCapabilitiesPage.OpenServerCapabilityPage();
                });
            });
            it('can logout', async () => {
                expect(await topNavigation.Logout()).toBe(true);
            });
        });
    });
});
describe("Clean up prereq", () => {
    it('Clean up', async () => {
        await api.UseAPI(serverServerCapabilities.cleanup);
    });
});
