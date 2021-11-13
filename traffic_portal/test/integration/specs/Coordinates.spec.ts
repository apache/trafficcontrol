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
import { CoordinatesPage } from '../PageObjects/CoordinatesPage.po';
import { api } from "../config";
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { coordinates } from "../Data";

const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const coordinatesPage = new CoordinatesPage();

describe('Setup API for coordinates test', function () {
    it('Setup', async () => {
        await api.UseAPI(coordinates.setup);
    });
});

coordinates.tests.forEach(async coordinatesData => {
    coordinatesData.logins.forEach(login => {
        describe(`Traffic Portal - Coordinates - ${login.description}`, () => {

            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            });
            it('can open coordinates page', async () => {
                await coordinatesPage.OpenTopologyMenu();
                await coordinatesPage.OpenCoordinatesPage();
            });
            coordinatesData.add.forEach(add => {
                it(add.description, async function () {
                    expect(await coordinatesPage.CreateCoordinates(add)).toBeTruthy();
                    await coordinatesPage.OpenCoordinatesPage();
                });
            });
            coordinatesData.update.forEach(update => {
                it(update.description, async () => {
                    await coordinatesPage.SearchCoordinates(update.Name);
                    expect(await coordinatesPage.UpdateCoordinates(update)).toBeTruthy();
                    await coordinatesPage.OpenCoordinatesPage();
                });
            });

            coordinatesData.remove.forEach(remove => {
                it(remove.description, async () => {
                    await coordinatesPage.SearchCoordinates(remove.Name);
                    expect(await coordinatesPage.DeleteCoordinates(remove)).toBeTruthy();
                    await coordinatesPage.OpenCoordinatesPage();
                });
            });

            it('can logout', async () => {
                expect(await topNavigation.Logout()).toBeTruthy();
            });
        });
    });
});


describe('Clean up API for coordinates test', () => {
    it('Cleanup', async () => {
        await api.UseAPI(coordinates.cleanup);
    });
});
