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
import { by, element } from 'protractor';

import { randomize } from '../config';
import { BasePage } from './BasePage.po';
import {SideNavigationPage} from './SideNavigationPage.po';

interface Tenant {
  ParentTenant: string;
  Name: string;
  Active: string;
  existsMessage?: string;
  validationMessage?: string;
}

export class TenantsPage extends BasePage {

    private btnCreateNewTenant = element(by.xpath("//button[@title='Create New Tenant']"));
    private txtName = element(by.name('name'));
    private txtActive = element(by.name('active'));
    private txtParentTenant = element(by.name('parentId'));
    private txtSearch = element(by.id('tenantsTable_filter')).element(by.css('label input'));
    private btnDelete = element(by.buttonText('Delete'));
    private txtConfirmTenantName = element(by.name('confirmWithNameInput'));
    private randomize = randomize;

    async OpenTenantPage(){
      let snp = new SideNavigationPage();
      await snp.ClickUserAdminMenu();
      await snp.NavigateToTenantsPage();
    }

    public async CreateTenant(tenant: Tenant): Promise<boolean> {
        let result = false;
        let basePage = new BasePage();
        let snp = new SideNavigationPage();
        await this.btnCreateNewTenant.click();
        await this.txtName.sendKeys(tenant.Name+this.randomize);
        await this.txtActive.sendKeys(tenant.Active);
        if(tenant.ParentTenant == '- root'){
          await this.txtParentTenant.sendKeys(tenant.ParentTenant);
        }else{
          await this.txtParentTenant.sendKeys(tenant.ParentTenant+this.randomize);
        }
        await basePage.ClickCreate();
        if(await basePage.GetOutputMessage() == tenant.existsMessage){
          await snp.NavigateToTenantsPage();
          result = true;
        }else if(await basePage.GetOutputMessage() == tenant.validationMessage){
          result = true;
        }else{
          result = false;
        }
        return result;
    }
    async SearchTenant(name:string){
        let snp = new SideNavigationPage();
        await snp.NavigateToTenantsPage();
        await this.txtSearch.clear();
        await this.txtSearch.sendKeys(name);
        await element.all(by.repeater('t in ::tenants')).filter(function(row){
            return row.element(by.name('name')).getText().then(function(val){
              return val === name;
            });
          }).first().click();
    }
    async DeleteTenant(name:string,outputMessage:string){
        let result = false;
        let basePage = new BasePage();
        await this.btnDelete.click();
        await this.txtConfirmTenantName.sendKeys(name);
        await basePage.ClickDeletePermanently();
        result = await basePage.GetOutputMessage().then(function(value){
            if(outputMessage == value){
              return true;
            }else{
              return false;
            }
          })
          return result;
    }
}
