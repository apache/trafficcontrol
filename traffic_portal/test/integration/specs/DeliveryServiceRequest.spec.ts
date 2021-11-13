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
import { browser } from 'protractor'
import { LoginPage } from '../PageObjects/LoginPage.po'
import { DeliveryServicesRequestPage } from '../PageObjects/DeliveryServiceRequestPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { deliveryservicerequest } from '../Data/deliveryservicerequest';


const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const deliveryServiceRequestPage = new DeliveryServicesRequestPage();

deliveryservicerequest.tests.forEach(deliveryServiceRequestData => {
    deliveryServiceRequestData.logins.forEach(login => {
        describe(`Traffic Portal - Delivery Service Request - ${login.description}`, () => {
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            })
            it('can open delivery service page', async function () {
                await deliveryServiceRequestPage.OpenServicesMenu();
                await deliveryServiceRequestPage.OpenDeliveryServicePage();
            })
            deliveryServiceRequestData.add.forEach(add => {
                it(add.description, async () => {
                    expect(await deliveryServiceRequestPage.CreateDeliveryService(add)).toBe(true);
                });
            });
            it('can logout', async () => {
                expect(await topNavigation.Logout()).toBeTruthy();
            });
        });
    });
});

