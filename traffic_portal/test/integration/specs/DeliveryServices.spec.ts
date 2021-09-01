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
import { DeliveryServicePage } from '../PageObjects/DeliveryServicePage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { api } from "../config";
import { deliveryservices } from "../Data/deliveryservices";

const topNavigation = new TopNavigationPage();
const loginPage = new LoginPage();
const deliveryservicesPage = new DeliveryServicePage();

describe('Setup API for delivery service test', function () {
    it('Setup', async () => {
        await api.UseAPI(deliveryservices.setup);
    });
});

deliveryservices.tests.forEach(async deliveryservicesData => {
    deliveryservicesData.logins.forEach(login =>{
        describe(`Traffic Portal - Delivery Service - ${login.description}`, () =>{
            afterEach(async function () {
                await deliveryservicesPage.OpenDeliveryServicePage();
            });
            afterAll(async function () {
                await deliveryservicesPage.OpenServicesMenu();
                expect(await topNavigation.Logout()).toBe(true);
            })
            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBe(true);
                await deliveryservicesPage.OpenServicesMenu();
            });
            it('can perform add test suits', async () => {
                deliveryservicesData.add.forEach(add => {
                    it(add.description, async function () {
                        expect(await deliveryservicesPage.CreateDeliveryService(add)).toBe(true);
                    });
                });
            })
            it('can perform update suits', async () => {
                deliveryservicesData.update.forEach(update => {
                    it(update.description, async function () {
                        await deliveryservicesPage.SearchDeliveryService(update.Name);
                        expect(await deliveryservicesPage.UpdateDeliveryService(update)).toBe(true);
                    });
                })
            })
            it('can perform assignserver suits', async () => {
                deliveryservicesData.assignserver.forEach(assignserver => {
                    it(assignserver.description, async function(){
                        await deliveryservicesPage.SearchDeliveryService(assignserver.DSName);
                        expect(await deliveryservicesPage.AssignServerToDeliveryService(assignserver)).toBe(true);
                    })
                })
            })
            it('can perform assignrequirecapabilities suits', async () => {
                deliveryservicesData.assignrequiredcapabilities.forEach(assignrc => {
                    it(assignrc.description, async function(){
                        await deliveryservicesPage.SearchDeliveryService(assignrc.DSName);
                        expect(await deliveryservicesPage.AssignRequiredCapabilitiesToDS(assignrc)).toBe(true);
                    })
                })  
            })
            it('can perform remove test suits', async () => {
                deliveryservicesData.remove.forEach(remove => {
                    it(remove.description, async () => {
                        await deliveryservicesPage.SearchDeliveryService(remove.Name);
                        expect(await deliveryservicesPage.DeleteDeliveryService(remove)).toBe(true);
                    });
                });  
            })
        })
    })
})

describe('Clean up API for delivery service test', () => {
    afterAll(async function () {
        it('Cleanup', async () => {
            await api.UseAPI(deliveryservices.cleanup);
        });
    })
});
