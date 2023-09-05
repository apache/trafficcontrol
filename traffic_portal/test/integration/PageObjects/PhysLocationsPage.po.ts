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

interface UpdatePhysicalLocation {
  description: string;
  Region: string;
  validationMessage?: string;
}

interface CreatePhysicalLocation {
  Address: string;
  City: string;
  Comments: string;
  Email: string;
  Name: string;
  Phone: string;
  Poc: string;
  Region: string;
  ShortName: string;
  State: string;
  Zip: string;
  validationMessage?: string;
}

interface DeletePhysicalLocation {
  Name: string;
  validationMessage?: string;
}

export class PhysLocationsPage extends BasePage {

  private btnCreateNewPhysLocation = element(by.name('createPhysLocationButton'));
  private txtName = element(by.name('name'));
  private txtShortName = element(by.name('shortName'));
  private txtAddress = element(by.name('address'));
  private txtCity = element(by.name('city'));
  private txtState = element(by.name('state'));
  private txtZip = element(by.name('zip'));
  private txtPoc = element(by.name('poc'));
  private txtPhone = element(by.name('phone'));
  private txtEmail = element(by.name('email'));
  private txtRegion = element(by.name('region'));
  private txtComments = element(by.name('comments'));
  private txtSearch = element(by.id('physLocationsTable_filter')).element(by.css('label input'));
  private btnDelete = element(by.buttonText('Delete'));
  private txtConfirmName = element(by.name('confirmWithNameInput'));
  private randomize = randomize;

  public async OpenPhysLocationPage() {
    const snp = new SideNavigationPage();
    await snp.NavigateToPhysLocation();
  }
  public async OpenConfigureMenu() {
    const snp = new SideNavigationPage();
    await snp.ClickTopologyMenu();
  }

  public async CreatePhysLocation(physlocation: CreatePhysicalLocation): Promise<boolean> {
    let result = false;
    const basePage = new BasePage();
    const snp = new SideNavigationPage();
    await snp.NavigateToPhysLocation();
    await this.btnCreateNewPhysLocation.click();
    await this.txtName.sendKeys(physlocation.Name + this.randomize);
    await this.txtShortName.sendKeys(physlocation.ShortName);
    await this.txtAddress.sendKeys(physlocation.Address);
    await this.txtCity.sendKeys(physlocation.City);
    await this.txtState.sendKeys(physlocation.State);
    await this.txtZip.sendKeys(physlocation.Zip);
    await this.txtPoc.sendKeys(physlocation.Poc);
    await this.txtPhone.sendKeys(physlocation.Phone);
    await this.txtEmail.sendKeys(physlocation.Email);
    await this.txtRegion.sendKeys(physlocation.Region + this.randomize);
    await this.txtComments.sendKeys(physlocation.Comments);
    await basePage.ClickCreate();
    result = await basePage.GetOutputMessage().then(function (value) {
      if (physlocation.validationMessage == value) {
        return true;
      } else {
        return false;
      }
    })
    return result;
  }

  public async SearchPhysLocation(physlocationName: string): Promise<void> {
    const snp = new SideNavigationPage();
    let name = physlocationName + this.randomize;
    await snp.NavigateToPhysLocation();
    await this.txtSearch.clear();
    await this.txtSearch.sendKeys(name);
    await element.all(by.repeater('pl in ::physLocations')).filter(function (row) {
      return row.element(by.name('name')).getText().then(function (val) {
        return val === name;
      });
    }).first().click();
  }

  public async UpdatePhysLocation(physlocation: UpdatePhysicalLocation): Promise<boolean> {
    const basePage = new BasePage();

    switch (physlocation.description) {
      case "update physlocation region":
        await this.txtRegion.sendKeys(physlocation.Region + this.randomize);
        await basePage.ClickUpdate();
        break;
      default:
        return false;
    }
    return await basePage.GetOutputMessage().then(value => physlocation.validationMessage === value);
  }

  public async DeletePhysLocation(physlocation: DeletePhysicalLocation): Promise<boolean> {
    let result = false;
    const basePage = new BasePage();
    await this.btnDelete.click();
    await this.txtConfirmName.sendKeys(physlocation.Name + this.randomize);
    await basePage.ClickDeletePermanently();
    result = await basePage.GetOutputMessage().then(function (value) {
      if (value.indexOf(physlocation.validationMessage ?? "") > -1) {
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
