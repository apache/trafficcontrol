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
import { browser, by, element } from 'protractor'

import { twoNumberRandomize } from '../config';
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';

interface CreateASN {
    ASNs: string;
    CacheGroup: string;
    validationMessage?: string;
}

interface DeleteASN {
    ASNs: string;
    validationMessage?: string;
}

interface UpdateASN {
    CacheGroup?: string;
    description: string;
    NewASNs?: string;
    validationMessage?: string;
}

export class ASNsPage extends BasePage {
    private btnCreateNewASNs = element(by.xpath("//button[@title='Create ASN']"));
    private txtSearch = element(by.id('asnsTable_filter')).element(by.css('label input'));
    private txtASN = element(by.name("asn"));
    private txtCacheGroup = element(by.name("cachegroup"))
    private btnDelete = element(by.xpath("//button[text()='Delete']"));
    private txtConfirmName = element(by.name('confirmWithNameInput'));

    async OpenASNsPage() {
        let snp = new SideNavigationPage();
        await snp.NavigateToASNsPage();
    }
    async OpenTopologyMenu() {
        let snp = new SideNavigationPage();
        await snp.ClickTopologyMenu();
    }

    public async CreateASNs(asns: CreateASN): Promise<boolean> {
        let basePage = new BasePage();
        let snp = new SideNavigationPage();
        await snp.NavigateToASNsPage();
        await this.btnCreateNewASNs.click();
        await this.txtASN.sendKeys(asns.ASNs + twoNumberRandomize);
        await this.txtCacheGroup.sendKeys(asns.CacheGroup)
        await basePage.ClickCreate();
        return await basePage.GetOutputMessage().then(v => v.indexOf(asns.validationMessage ?? "") > -1);
    }

    public async SearchASNs(nameASNs: string): Promise<boolean> {
        let name = nameASNs + twoNumberRandomize;
        let snp = new SideNavigationPage();
        await snp.NavigateToASNsPage();
        await this.txtSearch.clear();
        await this.txtSearch.sendKeys(name);
        if (await browser.isElementPresent(element(by.xpath("//td[@data-search='^" + name + "$']"))) == true) {
            await element(by.xpath("//td[@data-search='^" + name + "$']")).click();
            return true;
        }
        return false;
    }

    public async UpdateASNs(asns: UpdateASN): Promise<boolean> {
        let result = false;
        let basePage = new BasePage();
        if(asns.description.includes("update cachegroup")){
            // preserves old behavior, but with a better error message
            if (!asns.CacheGroup) {
                throw new Error("ASN update data indicated in the description that it was for updating cachegroup linking, but data included no CacheGroup");
            }
            await this.txtCacheGroup.sendKeys(asns.CacheGroup);
            await basePage.ClickUpdate();
        }else if(asns.description.includes("update an ASN")){
            // preserves old behavior, but with a better error message
            if (!asns.NewASNs) {
                throw new Error("ASN update data indicated in the description that it was NOT for updating cachegroup linking, but data included no NewASNs");
            }
            await this,this.txtASN.clear();
            await this.txtASN.sendKeys(asns.NewASNs + twoNumberRandomize);
            await basePage.ClickUpdate();
        }else{
            result = false;
        }
        result = await basePage.GetOutputMessage().then(function (value) {
            if (value.indexOf(asns.validationMessage ?? "") > -1) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }

    public async DeleteASNs(asns: DeleteASN): Promise<boolean> {
        let name = asns.ASNs + twoNumberRandomize;
        let result = false;
        let basePage = new BasePage();
        await this.btnDelete.click();
        await this.txtConfirmName.sendKeys(name);
        await basePage.ClickDeletePermanently();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (asns.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }


}
