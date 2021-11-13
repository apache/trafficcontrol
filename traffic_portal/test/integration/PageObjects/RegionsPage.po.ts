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

import { randomize } from '../config';
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';

interface CreateRegion {
    Division: string;
    Name: string;
    validationMessage?: string;
}

interface UpdateRegion {
    description: string;
    Division: string;
    validationMessage?: string;
}

interface DeleteRegion {
    Name: string;
    validationMessage?: string;
}

export class RegionsPage extends BasePage {
    private btnCreateNewRegions = element(by.name('createRegionButton'));
    private txtSearch = element(by.id('regionsTable_filter')).element(by.css('label input'));
    private txtName = element(by.id('name'));
    private txtDivision = element(by.name('division'));
    private btnDelete = element(by.buttonText('Delete'));
    private txtConfirmName = element(by.name('confirmWithNameInput'));
    private randomize = randomize;

    public async OpenRegionsPage(){
        let snp = new SideNavigationPage();
        await snp.NavigateToRegionsPage();
    }
    public async OpenTopologyMenu(){
        let snp = new SideNavigationPage();
        await snp.ClickTopologyMenu();
    }

    public async CreateRegions(regions: CreateRegion): Promise<boolean> {
        let result = false;
        let basePage = new BasePage();
        let snp = new SideNavigationPage();
        await snp.NavigateToRegionsPage();
        await this.btnCreateNewRegions.click();
        await this.txtName.sendKeys(regions.Name + this.randomize);
        await this.txtDivision.sendKeys(regions.Division);
        await basePage.ClickCreate();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (regions.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }

    public async SearchRegions(nameRegions:string): Promise<boolean> {
        let name = nameRegions + this.randomize;
        let snp = new SideNavigationPage();
        await snp.NavigateToRegionsPage();
        await this.txtSearch.clear();
        await this.txtSearch.sendKeys(name);
        if (await browser.isElementPresent(element(by.xpath("//td[@data-search='^" + name + "$']"))) == true) {
            await element(by.xpath("//td[@data-search='^" + name + "$']")).click();
            return true;
        }
        return false;
    }

    public async UpdateRegions(regions: UpdateRegion): Promise<boolean> {
        let basePage = new BasePage();
        switch(regions.description){
            case "update Region's Division":
                await this.txtDivision.sendKeys(regions.Division);
                await basePage.ClickUpdate();
                break;
            default:
                return false;
        }
        return await basePage.GetOutputMessage().then(value => regions.validationMessage === value);
    }

    public async DeleteRegions(regions: DeleteRegion): Promise<boolean> {
        let name = regions.Name + this.randomize;
        let result = false;
        let basePage = new BasePage();
        await this.btnDelete.click();
        await this.txtConfirmName.sendKeys(name);
        await basePage.ClickDeletePermanently();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (regions.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }

    public async CheckCSV(name: string): Promise<boolean> {
        return element(by.cssContainingText("span", name)).isPresent();
      }
}
