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
import { SideNavigationPage } from '../PageObjects/SideNavigationPage.po';

interface CreateCacheGroup {
    Type: string;
    Name: string;
    ShortName: string;
    Latitude: string;
    Longitude: string;
    ParentCacheGroup: string;
    SecondaryParentCG: string;
    FailoverCG?: string;
}

interface UpdateCacheGroup {
    Type: string;
    FailoverCG?: string;
}
export class CacheGroupPage extends BasePage {
    private btnCreateCacheGroups = element(by.name('createCacheGroupButton'));
    private txtName = element(by.name("name"))
    private txtShortName = element(by.name("shortName"));
    private txtType = element(by.name("type"));
    private txtLatitude = element(by.name("latitude"));
    private txtLongtitude = element(by.name("longitude"));
    private txtParentCacheGroup = element(by.name("parentCacheGroup"));
    private txtSecondaryParentCG = element(by.name("secondaryParentCacheGroup"));
    private txtFailoverCG = element(by.name("fallbackOptions"));
    private txtSearch = element(by.id('cacheGroupsTable_filter')).element(by.css('label input'));
    private txtConfirmCacheGroupName = element(by.name("confirmWithNameInput"));
    private btnDelete = element(by.buttonText('Delete'));
    private randomize = randomize;
    private btnTableColumn = element(by.className("caret"))
    async OpenTopologyMenu() {
        let snp = new SideNavigationPage();
        await snp.ClickTopologyMenu();
    }
    async OpenCacheGroupsPage() {
        let snp = new SideNavigationPage();
        await snp.NavigateToCacheGroupsPage();
    }

    public async CreateCacheGroups(cachegroup: CreateCacheGroup, outputMessage: string): Promise<boolean> {
        let result = false
        let basePage = new BasePage();
        if (cachegroup.Type == "EDGE_LOC") {
            if (cachegroup.FailoverCG === undefined) {
                throw new Error(`cachegroups with Type 'EDGE_LOC' must have FailoverCG`);
            }
            await this.btnCreateCacheGroups.click();
            await this.txtName.sendKeys(cachegroup.Name + this.randomize);
            await this.txtShortName.sendKeys(cachegroup.ShortName + this.randomize);
            await this.txtType.sendKeys(cachegroup.Type);
            await this.txtLatitude.sendKeys(cachegroup.Latitude);
            await this.txtLongtitude.sendKeys(cachegroup.Longitude);
            await this.txtParentCacheGroup.sendKeys(cachegroup.ParentCacheGroup);
            await this.txtSecondaryParentCG.sendKeys(cachegroup.SecondaryParentCG);
            await this.txtFailoverCG.sendKeys(cachegroup.FailoverCG);
        } else {
            await this.btnCreateCacheGroups.click();
            await this.txtName.sendKeys(cachegroup.Name + this.randomize);
            await this.txtShortName.sendKeys(cachegroup.ShortName + this.randomize);
            await this.txtType.sendKeys(cachegroup.Type);
            await this.txtLatitude.sendKeys(cachegroup.Latitude);
            await this.txtLongtitude.sendKeys(cachegroup.Longitude);
            await this.txtParentCacheGroup.sendKeys(cachegroup.ParentCacheGroup);
            await this.txtSecondaryParentCG.sendKeys(cachegroup.SecondaryParentCG);
        }
        await basePage.ClickCreate();
        await basePage.GetOutputMessage().then(function (value) {
            if (outputMessage == value) {
                result = true;
            } else {
                result = false;
            }
        })
        return result;
    }

    public async SearchCacheGroups(nameCG: string): Promise<boolean> {
        let name = nameCG + this.randomize;
        await this.txtSearch.clear();
        await this.txtSearch.sendKeys(name);
        if (await browser.isElementPresent(element(by.xpath("//td[@data-search='^" + name + "$']"))) === true) {
            await element(by.xpath("//td[@data-search='^" + name + "$']")).click();
            return true;
        }
        return false;
    }

    public async UpdateCacheGroups(cachegroup: UpdateCacheGroup, outputMessage: string | undefined): Promise<boolean | undefined> {
        let result: boolean | undefined = false;
        let basePage = new BasePage();
        let snp = new SideNavigationPage();
        if (cachegroup.Type == "EDGE_LOC") {
            const name = cachegroup.FailoverCG + this.randomize;
            await this.txtFailoverCG.click();
            if (await browser.isElementPresent(element(by.css(`select[name="fallbackOptions"] > option[label="${name}"]`)))) {
                await element(by.css(`select[name="fallbackOptions"] > option[label="${name}"]`)).click();
            } else {
                result = undefined;
            }
        }
        await this.txtType.sendKeys(cachegroup.Type);
        await snp.ClickUpdate();
        if (result !== undefined) {
            await basePage.GetOutputMessage().then(function (value) {
                if (outputMessage === value) {
                    result = true;
                } else {
                    result = false;
                }
            })
        }
        return result;
    }
    public async DeleteCacheGroups(nameCG: string, outputMessage: string) {
        let result = false;
        let basePage = new BasePage();
        let snp = new SideNavigationPage();
        let name = nameCG + this.randomize;
        await this.btnDelete.click();
        await this.txtConfirmCacheGroupName.sendKeys(name);
        if (await basePage.ClickDeletePermanently() == true) {
            result = await basePage.GetOutputMessage().then(function (value) {
                if (outputMessage == value) {
                    return true
                } else {
                    return false;
                }
            })
        } else {
            await basePage.ClickCancel();
        }
        await snp.NavigateToCacheGroupsPage();
        return result;
    }
    public async CheckCSV(name: string): Promise<boolean> {
        return element(by.cssContainingText("span", name)).isPresent();
    }
    public async ToggleTableColumn(name: string): Promise<boolean> {
        await this.btnTableColumn.click();
        const result = await element(by.cssContainingText("th", name)).isPresent();
        await element(by.cssContainingText("label", name)).click();
        await this.btnTableColumn.click();
        return !result;
    }
}
