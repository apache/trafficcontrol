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


export class DivisionsPage extends BasePage {
    private btnCreateNewDivisions = element(by.name('createDivisionButton'));
    private txtSearch = element(by.id('divisionsTable_filter')).element(by.css('label input'));
    private txtName = element(by.id('name'));
    private btnDelete = element(by.xpath("//button[text()='Delete']"));
    private txtConfirmName = element(by.name('confirmWithNameInput'));
    private readonly config = config;
    private randomize = randomize;

    async OpenDivisionsPage(){
        let snp = new SideNavigationPage();
        await snp.NavigateToDivisionsPage();
    }
    async OpenTopologyMenu(){
        let snp = new SideNavigationPage();
        await snp.ClickTopologyMenu();
    }
    async CreateDivisions(divisions){
        let result = false;
        let basePage = new BasePage();
        let snp = new SideNavigationPage();
        await snp.NavigateToDivisionsPage();
        await this.btnCreateNewDivisions.click();
        await this.txtName.sendKeys(divisions.Name + this.randomize);
        await basePage.ClickCreate();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (divisions.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }
    async SearchDivisions(nameDivisions: string){
        let name = nameDivisions + this.randomize;
        let result = false;
        let snp = new SideNavigationPage();
        await snp.NavigateToDivisionsPage();
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
    async UpdateDivisions(divisions){
        let result = false;
        let basePage = new BasePage();
        await this.txtName.clear();
        await this.txtName.sendKeys(divisions.NewName + this.randomize);
        await basePage.ClickUpdate();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (divisions.validationMessage == value) {
              return true;
            } else {
              return false;
            }
          })
          return result;
    }
    async DeleteDivisions(divisions){
        let name = divisions.Name + this.randomize;
        let result = false;
        let basePage = new BasePage();
        await this.btnDelete.click();
        await this.txtConfirmName.sendKeys(name);
        await basePage.ClickDeletePermanently();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (divisions.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }

}
