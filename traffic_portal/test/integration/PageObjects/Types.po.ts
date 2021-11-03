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

interface CreateType {
    Name: string;
    DescriptionData: string;
    validationMessage?: string;
}

interface UpdateType {
    description: string;
    DescriptionData: string;
    validationMessage?: string;
}

interface DeleteType {
    Name: string;
    validationMessage?: string;
}

export class TypesPage extends BasePage {
    private btnCreateNewType = element(by.xpath("//button[@title='Create Type']//i[1]"));
    private txtName = element(by.id('name'));
    private txtDescription = element(by.xpath("//textarea[@name='description']"))
    private txtSearch = element(by.id('typesTable_filter')).element(by.css('label input'));
    private btnDelete = element(by.buttonText('Delete'));
    private txtConfirmName = element(by.name('confirmWithNameInput'));
    private btnTableColumn = element(by.className("caret"))
    private randomize = randomize;

    async OpenTypesPage() {
        let snp = new SideNavigationPage();
        await snp.NavigateToTypesPage();
    }
    async OpenConfigureMenu() {
        let snp = new SideNavigationPage();
        await snp.ClickConfigureMenu();
    }

    public async CreateType(type: CreateType): Promise<boolean> {
        let result = false;
        let basePage = new BasePage();
        await this.btnCreateNewType.click();
        await this.txtName.sendKeys(type.Name + this.randomize)
        await this.txtDescription.sendKeys(type.DescriptionData)
        await basePage.ClickCreate();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (type.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }

    public async SearchType(nameTypes: string): Promise<boolean> {
        let snp = new SideNavigationPage();
        let name = nameTypes + this.randomize;
        await snp.NavigateToTypesPage();
        await this.txtSearch.clear();
        await this.txtSearch.sendKeys(name);
        if (await browser.isElementPresent(element(by.xpath("//td[@data-search='^" + name + "$']"))) == true) {
            await element(by.xpath("//td[@data-search='^" + name + "$']")).click();
            return true;
        }
        return false;
    }

    public async UpdateType(type: UpdateType): Promise<boolean | undefined> {
        let basePage = new BasePage();
        switch (type.description) {
            case "update description type":
                await this.txtDescription.clear();
                await this.txtDescription.sendKeys(type.DescriptionData);
                await basePage.ClickUpdate();
                break;
            default:
                return undefined;
        }
        return await basePage.GetOutputMessage().then(value => type.validationMessage === value);
    }

    public async DeleteTypes(type: DeleteType): Promise<boolean> {
        let name = type.Name + this.randomize;
        let result = false;
        let basePage = new BasePage();
        await this.btnDelete.click();
        await this.txtConfirmName.sendKeys(name);
        await basePage.ClickDeletePermanently();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (type.validationMessage == value) {
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
    public async ToggleTableColumn(name: string): Promise<boolean> {
        await this.btnTableColumn.click();
        const result = await element(by.cssContainingText("th", name.charAt(0).toUpperCase() + name.slice(1))).isPresent();
        await element(by.cssContainingText("label", name)).click();
        await this.btnTableColumn.click();
        return !result;
    }


}
