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

import { config, randomize } from '../config';
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';

export class ProfilesPage extends BasePage {

    private btnCreateNewProfile = element(by.name('createProfileButton'));
    private txtName = element(by.name('name'));
    private txtCDN = element(by.name('cdn'));
    private txtType = element(by.name('type'));
    private txtRoutingDisable = element(by.name('routingDisabled'));
    private txtDescription = element(by.id('description'));
    private txtSearch = element(by.id('profilesTable_filter')).element(by.css('label input'));
    private mnuProfilesTable = element(by.id('profilesTable'));
    private btnDelete = element(by.buttonText('Delete'));
    private txtConfirmName = element(by.name('confirmWithNameInput'));
    private btnMore = element(by.name('moreBtn'));
    private btnCompareProfile = element(by.name('compareProfilesBtn'));
    private txtCompareDropdown1 = element(by.name('compareDropdown1'));
    private txtCompareDropdown2 = element(by.name('compareDropdown2'));
    private btnCompareSubmit = element(by.name('compareSubmit'));
    private mnuCompareTable = element(by.id('profilesParamsCompareTable_wrapper'));
    private readonly config = config;
    private randomize = randomize;
    async OpenProfilesPage() {
        let snp = new SideNavigationPage();
        await snp.NavigateToProfilesPage();
    }
    async OpenConfigureMenu() {
        let snp = new SideNavigationPage();
        await snp.ClickConfigureMenu();
    }
    async CreateProfile(profile) {
        let result = false;
        let basePage = new BasePage();
        let snp = new SideNavigationPage();
        await snp.NavigateToProfilesPage();
        await this.btnCreateNewProfile.click();
        await this.txtName.sendKeys(profile.Name + this.randomize);
        await this.txtCDN.sendKeys(profile.CDN);
        await this.txtType.sendKeys(profile.Type);
        await this.txtRoutingDisable.sendKeys(profile.RoutingDisable);
        await this.txtDescription.sendKeys(profile.Description);
        await basePage.ClickCreate();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (profile.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }
    async SearchProfile(nameProfiles: string) {
        let result = false;
        let snp = new SideNavigationPage();
        let name = nameProfiles + this.randomize;
        await snp.NavigateToProfilesPage();
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
    async CompareProfile(profile1: string, profile2: string) {
        let result = false;
        let snp = new SideNavigationPage();
        await snp.NavigateToProfilesPage();
        await this.btnMore.click();
        await this.btnCompareProfile.click();
        await this.txtCompareDropdown1.sendKeys(profile1);
        await this.txtCompareDropdown2.sendKeys(profile2);
        await this.btnCompareSubmit.click();
        if (await this.mnuCompareTable.isDisplayed() == true) {
            result = true;
            return result;
        }
    }
    async UpdateProfile(profile) {
        let result = false;
        let basePage = new BasePage();
        switch (profile.description) {
            case "update profile type":
                await this.txtType.sendKeys(profile.Type);
                await basePage.ClickUpdate();
                break;
            default:
                result = undefined;
        }
        if (result =! undefined) {
            result = await basePage.GetOutputMessage().then(function (value) {
                if (profile.validationMessage == value) {
                    return true;
                } else {
                    return false;
                }
            })

        }
        return result;
    }
    async DeleteProfile(profile) {
        let result = false;
        let basePage = new BasePage();
        await this.btnDelete.click();
        await this.txtConfirmName.sendKeys(profile.Name + this.randomize);
        await basePage.ClickDeletePermanently();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (profile.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }
}
