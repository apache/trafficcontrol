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
import { API } from '../CommonUtils/API';
import { deliveryservices } from "../Data/deliveryservices";

const api = new API();
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
            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBe(true);
            });
            it('can open delivery service page', async () => {
                await deliveryservicesPage.OpenServicesMenu();
                await deliveryservicesPage.OpenDeliveryServicePage();
            });
            deliveryservicesData.add.forEach(add => {
                it(add.description, async function () {
                    expect(await deliveryservicesPage.CreateDeliveryService(add)).toBe(true);
                    await deliveryservicesPage.OpenDeliveryServicePage();
                });
            });
            deliveryservicesData.update.forEach(update => {
                it(update.description, async function () {
                    await deliveryservicesPage.SearchDeliveryService(update.Name);
                    expect(await deliveryservicesPage.UpdateDeliveryService(update)).toBe(true);
                    await deliveryservicesPage.OpenDeliveryServicePage();
                });
            })
            deliveryservicesData.assignserver.forEach(assignserver => {
                it(assignserver.description, async function(){
                    await deliveryservicesPage.SearchDeliveryService(assignserver.DSName);
                    expect(await deliveryservicesPage.AssignServerToDeliveryService(assignserver)).toBe(true);
                    await deliveryservicesPage.OpenDeliveryServicePage();
                }
                
                )
            })
            deliveryservicesData.assignrequiredcapabilities.forEach(assignrc => {
                it(assignrc.description, async function(){
                    await deliveryservicesPage.SearchDeliveryService(assignrc.DSName);
                    expect(await deliveryservicesPage.AssignRequiredCapabilitiesToDS(assignrc)).toBe(true);
                    await deliveryservicesPage.OpenDeliveryServicePage();
                })
            })
            deliveryservicesData.remove.forEach(remove => {
                it(remove.description, async () => {
                    await deliveryservicesPage.SearchDeliveryService(remove.Name);
                    expect(await deliveryservicesPage.DeleteDeliveryService(remove)).toBe(true);
                    await deliveryservicesPage.OpenDeliveryServicePage();
                });
            });
            it('can close service menu tab', async () => {
                await deliveryservicesPage.OpenServicesMenu();
            });
            it('can logout', async () => {
                expect(await topNavigation.Logout()).toBe(true);
            });
        })
    })

})
describe('Clean up API for delivery service test', () => {
    it('Cleanup', async () => {
        await api.UseAPI(deliveryservices.cleanup);
    });
});
