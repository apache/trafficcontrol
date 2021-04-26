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
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';

interface Topologies {
    description: string;
    Name:string;
    DescriptionData: string;
    Type: string;
    CacheGroup: string;
    validationMessage: string;
}
interface DeleteTopologies{
    Name: string;
    validationMessage?: string;
}
export class TopologiesPage extends BasePage {
    private btnCreateNewTopologies = element(by.xpath("//button[@title='Create Topology']"));
    private txtCacheGroupType = element(by.name('selectFormDropdown'));
    private txtSearch = element(by.id('topologiesTable_filter')).element(by.css('label input'));
    private txtName = element(by.name('name'));
    private txtDescription = element(by.id('description'));
    private btnAddCacheGroup = element(by.xpath("//a[@title='Add child cache groups to TOPOLOGY']"));
    private txtSearchCacheGroup = element(by.id('availableCacheGroupsTable_filter')).element(by.css('label input'));
    private btnDelete = element(by.xpath("//button[text()='Delete']"));
    private txtConfirmName = element(by.name('confirmWithNameInput'));
    private config = require('../config');
    private randomize = this.config.randomize;
    async OpenTopologiesPage(){
        let snp = new SideNavigationPage();
        await snp.NavigateToTopologiesPage();
    }
    async OpenTopologyMenu(){
        let snp = new SideNavigationPage();
        await snp.ClickTopologyMenu();
    }
    public async CreateTopologies(topologies: Topologies): Promise<boolean> {
        let basePage = new BasePage();
        let snp = new SideNavigationPage();
        await snp.NavigateToTopologiesPage();
        //click '+'
        await this.btnCreateNewTopologies.click();
        await this.txtName.sendKeys(topologies.Name + this.randomize)
        await this.txtDescription.sendKeys(topologies.DescriptionData + this.randomize)
        //click add cache group +
        await this.btnAddCacheGroup.click();
        //choose type
        await this.txtCacheGroupType.sendKeys(topologies.Type);
        await basePage.ClickSubmit();
        //choose Cachegroup
        await this.txtSearchCacheGroup.sendKeys(topologies.CacheGroup + this.randomize)
        if(await browser.isElementPresent(by.xpath("//td[@data-search='^" + topologies.CacheGroup + this.randomize + "$']")) === true){
            await element(by.xpath("//td[@data-search='^" + topologies.CacheGroup + this.randomize + "$']")).click();
        }
        await basePage.ClickSubmit();
        await basePage.ClickCreate();
        if(topologies.description === "create a Topologies with empty cachegroup (no server)"){
            topologies.validationMessage = topologies.validationMessage + this.randomize;
        }
        return await basePage.GetOutputMessage().then(value => value === topologies.validationMessage);
    }
       
    async SearchTopologies(nameTopologies:string): Promise<boolean>{
        let name = nameTopologies + this.randomize;
        let snp = new SideNavigationPage();
        await snp.NavigateToTopologiesPage();
        await this.txtSearch.clear();
        await this.txtSearch.sendKeys(name);
        if (await browser.isElementPresent(element(by.xpath("//td[@data-search='^" + name + "$']"))) == true) {
            await element(by.xpath("//td[@data-search='^" + name + "$']")).click();
            return true;
        } 
        return false;
    }
    async DeleteTopologies(topologies: DeleteTopologies):Promise<boolean>{
        let name = topologies.Name + this.randomize;
        let basePage = new BasePage();
        await this.btnDelete.click();
        await this.txtConfirmName.sendKeys(name);
        await this.ClickDeletePermanently();
        return await basePage.GetOutputMessage().then(value => value === topologies.validationMessage);
    }
}
