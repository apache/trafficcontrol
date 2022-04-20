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
import { browser, by, element } from 'protractor';

import { randomize } from "../config";
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';

interface Origin {
    deliveryServiceId: string;
    Name: string;
    Tenant: string;
    FQDN: string;
    Protocol: string;
    validationMessage?: string;
}

interface UpdateOrigin {
    NewDeliveryService: string;
    validationMessage?: string;
}

interface DeleteOrigin {
    Name: string;
    validationMessage?: string;
}

export class OriginsPage extends BasePage {
    private btnCreateNewOrigins = element(by.xpath("//button[@title='Create Origin']"));
    private txtSearch = element(by.id('originsTable_filter')).element(by.css('label input'));
    private txtName = element(by.name('name'));
    private selTenant = element(by.name('tenantId'));
    private txtFQDN = element(by.name('fqdn'));
    private txtProtocol = element(by.name('protocol'));
    private txtDeliveryService = element(by.name("deliveryServiceId"));
    private btnDelete = element(by.xpath("//button[text()='Delete']"));
    private txtConfirmName = element(by.name('confirmWithNameInput'));
    private randomize = randomize;

    async OpenOriginsPage() {
        let snp = new SideNavigationPage();
        await snp.NavigateToOriginsPage();
    }

    async OpenConfigureMenu() {
        let snp = new SideNavigationPage();
        await snp.ClickConfigureMenu();
    }

    public async CreateOrigins(origins: Origin): Promise<boolean> {
        let result = false;
        let basePage = new BasePage();
        let snp = new SideNavigationPage();
        await snp.NavigateToOriginsPage();
        await this.btnCreateNewOrigins.click();
        await this.txtName.sendKeys(origins.Name + this.randomize);
        await this.selTenant.click();
        await element(by.name(origins.Tenant + this.randomize)).click();
        await this.txtFQDN.sendKeys(origins.FQDN);
        await this.txtProtocol.sendKeys(origins.Protocol);
        await this.txtDeliveryService.click();
        await element(by.xpath("//option[@label='" + origins.deliveryServiceId + this.randomize + "']")).click();
        await basePage.ClickCreate();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (origins.validationMessage === value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }

    public async SearchOrigins(nameOrigins: string): Promise<boolean> {
        let name = nameOrigins + this.randomize;
        let snp = new SideNavigationPage();
        await snp.NavigateToOriginsPage();
        await this.txtSearch.clear();
        await this.txtSearch.sendKeys(name);
        if (await browser.isElementPresent(element(by.xpath("//td[@data-search='^" + name + "$']"))) == true) {
            await element(by.xpath("//td[@data-search='^" + name + "$']")).click();
            return true;
        }
        return false;
    }

    async UpdateOrigins(origins: UpdateOrigin): Promise<boolean | undefined> {
        let result: boolean | undefined = false;
        let basePage = new BasePage();
        if (origins.NewDeliveryService != null || origins.NewDeliveryService != undefined) {
            if (await browser.isElementPresent(element(by.xpath(`//select[@name="deliveryServiceId"]//option[@label="` + origins.NewDeliveryService + this.randomize + `"]`)))) {
                await element(by.xpath(`//select[@name="deliveryServiceId"]//option[@label="` + origins.NewDeliveryService + this.randomize + `"]`)).click();
            } else {
                result = undefined;
            }
        }
        if (result != undefined) {
            await basePage.ClickUpdate();
            result = await basePage.GetOutputMessage().then(function (value) {
                if (origins.validationMessage == value) {
                    return true;
                } else {
                    return false;
                }
            })
        }
        return result;
    }

    public async DeleteOrigins(origins: DeleteOrigin): Promise<boolean> {
        let name = origins.Name + this.randomize;
        let result = false;
        let basePage = new BasePage();
        await this.btnDelete.click();
        await this.txtConfirmName.sendKeys(name);
        await basePage.ClickDeletePermanently();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (origins.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }
}
