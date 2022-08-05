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

describe("Delivery Services", () => {
	beforeAll(async () => {
		await api.UseAPI(deliveryservices.setup);
	});

	afterAll(async () => {
		await api.UseAPI(deliveryservices.cleanup);
	});

	for (const data of deliveryservices.tests) {
		describe(`Traffic Portal - Delivery Service - ${data.description}`, () =>{
			beforeAll(async () => {
				browser.get(browser.params.baseUrl);
				await loginPage.Login(data.login);
				expect(await loginPage.CheckUserName(data.login)).toBe(true);
				await deliveryservicesPage.OpenServicesMenu();
				await deliveryservicesPage.OpenDeliveryServicePage();
			});
			afterEach(async () => {
				await deliveryservicesPage.OpenDeliveryServicePage();
				expect((await browser.getCurrentUrl()).split("#").slice(-1).join().replace(/\/$/, "")).toBe("!/delivery-services");
			});
			afterAll(async () => {
				await deliveryservicesPage.OpenServicesMenu();
				return expectAsync(topNavigation.Logout()).toBeResolvedTo(true);
			});

			for (const {description, name, type, tenant, validationMessage} of data.add) {
				it(description, async () => {
					return expectAsync(deliveryservicesPage.CreateDeliveryService(name, type, tenant)).toBeResolvedTo(validationMessage);
				});
			}
			for (const {name, newName, validationMessage} of data.update) {
				it("updates Delivery Service Display Name", async () => {
					await deliveryservicesPage.SearchDeliveryService(name);
					return expectAsync(deliveryservicesPage.UpdateDeliveryServiceDisplayName(newName)).toBeResolved(validationMessage);
				});
			}

			for (const {serverHostname, xmlID, validationMessage} of data.assignServer){
				it("assigns servers to a Delivery Service", async () => {
					await deliveryservicesPage.SearchDeliveryService(xmlID);
					return expectAsync(deliveryservicesPage.AssignServerToDeliveryService(serverHostname)).toBeResolvedTo(validationMessage);
				});
			}

			for (const {rcName, validationMessage, xmlID} of data.assignRequiredCapabilities) {
				it("assign required capabilities to delivery service", async () => {
					await deliveryservicesPage.SearchDeliveryService(xmlID);
					return expectAsync(deliveryservicesPage.AssignRequiredCapabilitiesToDS(rcName)).toBeResolvedTo(validationMessage);
				});
			}

			for (const {name, validationMessage} of data.remove) {
				it("deletes a Delivery Service", async () => {
					await deliveryservicesPage.SearchDeliveryService(name);
					return expectAsync(deliveryservicesPage.DeleteDeliveryService(name)).toBeResolvedTo(validationMessage);
				});
			}
		});
	}
});
