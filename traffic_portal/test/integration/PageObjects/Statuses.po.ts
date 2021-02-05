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
import { ElementFinder, browser, by, element } from 'protractor'
import { async, delay } from 'q';
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';

export class StatusesPage extends BasePage {
    private btnCreateNewStatus = element(by.xpath("//button[@title='Create Status']"))
    private txtName = element(by.id('name'));
    private txtDescription = element(by.xpath("//textarea[@name='description']"))
    private txtSearch = element(by.id('statusesTable_filter')).element(by.css('label input'));
    private btnDelete = element(by.buttonText('Delete'));
    private txtConfirmName = element(by.name('confirmWithNameInput'));
    private config = require('../config');
    private randomize = this.config.randomize;

    async OpenStatusesPage() {
        let snp = new SideNavigationPage();
        await snp.NavigateToStatusesPage();
    }
    async OpenConfigureMenu() {
        let snp = new SideNavigationPage();
        await snp.ClickConfigureMenu();
    }
    async CreateStatus(status) {
        let result = false;
        let basePage = new BasePage();
        await this.btnCreateNewStatus.click();
        await this.txtName.sendKeys(status.Name + this.randomize)
        await this.txtDescription.sendKeys(status.DescriptionData)
        await basePage.ClickCreate();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (value.includes(status.validationMessage)) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }
    async SearchStatus(nameStatus: string) {
        let result = false;
        let snp = new SideNavigationPage();
        let name = nameStatus + this.randomize;
        await snp.NavigateToStatusesPage();
        await this.txtSearch.clear();
        await this.txtSearch.sendKeys(name);
        if (await browser.isElementPresent(element(by.xpath("//td[@data-search='^" + name + "$']"))) == true) {
            await element(by.xpath("//td[@data-search='^" + name + "$']")).click();
            result = true;
        } else {
            result = undefined;
        }
        return result;
    }
    async UpdateStatus(status) {
        let result = false;
        let basePage = new BasePage();
        switch (status.description) {
            case "update Status description":
                await this.txtDescription.clear();
                await this.txtDescription.sendKeys(status.DescriptionData);
                await basePage.ClickUpdate();
                break;
            default:
                result = undefined;
        }
        result = await basePage.GetOutputMessage().then(function (value) {
            if (status.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }
    async DeleteStatus(status) {
        let name = status.Name + this.randomize;
        let result = false;
        let basePage = new BasePage();
        await this.btnDelete.click();
        await this.txtConfirmName.sendKeys(name);
        await basePage.ClickDeletePermanently();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (status.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;

    }

}