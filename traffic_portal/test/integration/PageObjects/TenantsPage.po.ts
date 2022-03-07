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
import { SideNavigationPage } from './SideNavigationPage.po';

interface Tenant {
  ParentTenant: string;
  Name: string;
  Active: string;
  validationMessage?: string;
}
interface DeleteTenant {
  Name: string;
  validationMessage: string;
}
interface UpdateTenant {
  Active: string;
  validationMessage: string;
}

export class TenantsPage extends BasePage {
    private btnCreateNewTenant = element(by.xpath("//button[@title='Create New Tenant']"));
    private txtName = element(by.name('name'));
    private txtActive = element(by.name('active'));
    private selParentTenant = element(by.name('parentId'));
    private btnDelete = element(by.buttonText('Delete'));
    private txtConfirmTenantName = element(by.name('confirmWithNameInput'));
    private randomize = randomize;

    public async OpenUserAdminMenu() {
      let snp = new SideNavigationPage();
      await snp.ClickUserAdminMenu();
    }

    public async OpenTenantsPage() {
      let snp = new SideNavigationPage();
      await snp.NavigateToTenantsPage();
    }

    public async CreateTenant(tenant: Tenant): Promise<boolean> {
      const basePage = new BasePage();
      await this.btnCreateNewTenant.click();
      await this.txtName.sendKeys(tenant.Name + this.randomize);
      await this.txtActive.sendKeys(tenant.Active);
      await this.selParentTenant.click();
      await element(by.name(tenant.ParentTenant + this.randomize)).click();
      await basePage.ClickCreate();
      return basePage.GetOutputMessage().then(value => tenant.validationMessage === value);
    }

    public async SearchTenant(name: string) {
      let snp = new SideNavigationPage();
      await snp.NavigateToTenantsPage();
      await element(by.linkText(name + this.randomize)).click();
    }

    public async DeleteTenant(tenant: DeleteTenant) {
      let basePage = new BasePage();
      await this.btnDelete.click();
      await this.txtConfirmTenantName.sendKeys(tenant.Name + this.randomize);
      await basePage.ClickDeletePermanently();
      return basePage.GetOutputMessage().then(value => tenant.validationMessage === value);
    }

    public async UpdateTenant(tenant: UpdateTenant) {
      const basePage = new BasePage();
      await this.txtActive.sendKeys(tenant.Active);
      await basePage.ClickUpdate();
      return basePage.GetOutputMessage().then(value => tenant.validationMessage === value);
    }
}
