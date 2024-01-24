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
import { browser } from "protractor";

import { LoginPage } from "../PageObjects/LoginPage.po";
import { TopNavigationPage } from "../PageObjects/TopNavigationPage.po";
import { api } from "../config";
import { ServiceCategoriesPage } from "../PageObjects/ServiceCategories.po";
import { serviceCategories } from "../Data";

const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const serviceCategoriesPage = new ServiceCategoriesPage();

describe("Setup API for Service Categories Test", () => {
    it("Setup", async () => {
        await api.UseAPI(serviceCategories.setup);
    });
});
serviceCategories.tests.forEach(async (serviceCategoriesData) => {
    serviceCategoriesData.logins.forEach((login) => {
        describe(`Traffic Portal - ServiceCategories - ${login.description}`, () => {
            it("can login", async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            });
            it("can open service categories page", async () => {
                await serviceCategoriesPage.OpenServicesMenu();
                await serviceCategoriesPage.OpenServiceCategoriesPage();
            });

            serviceCategoriesData.add.forEach((add) => {
                it(add.description, async () => {
                    expect(
                        await serviceCategoriesPage.CreateServiceCategories(
                            add,
                            add.validationMessage
                        )
                    ).toBeTruthy();
                    await serviceCategoriesPage.OpenServiceCategoriesPage();
                });
            });
            serviceCategoriesData.update.forEach((update) => {
                it(update.description, async () => {
                    await serviceCategoriesPage.SearchServiceCategories(
                        update.Name
                    );
                    expect(
                        await serviceCategoriesPage.UpdateServiceCategories(
                            update,
                            update.validationMessage
                        )
                    ).toBeTruthy();
                    await serviceCategoriesPage.OpenServiceCategoriesPage();
                });
            });
            serviceCategoriesData.remove.forEach((remove) => {
                it(remove.description, async () => {
                    await serviceCategoriesPage.SearchServiceCategories(
                        remove.Name
                    );
                    expect(
                        await serviceCategoriesPage.DeleteServiceCategories(
                            remove,
                            remove.validationMessage
                        )
                    ).toBeTruthy();
                    await serviceCategoriesPage.OpenServiceCategoriesPage();
                });
            });
            it("can logout", async () => {
                expect(await topNavigation.Logout()).toBeTruthy();
            });
        });
    });
});

describe("Clean Up API for Service Categories Test", () => {
    it("Cleanup", async () => {
        await api.UseAPI(serviceCategories.cleanup);
    });
});
