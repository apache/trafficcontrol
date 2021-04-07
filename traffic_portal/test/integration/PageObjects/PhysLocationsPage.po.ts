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

import { config, randomize } from '../config';
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';

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
  private mnuPhysLocationsTable = element(by.id('physLocationsTable'));
  private btnDelete = element(by.buttonText('Delete'));
  private txtConfirmName = element(by.name('confirmWithNameInput'));
  private readonly config = config;
  private randomize = randomize;

  async OpenPhysLocationPage() {
    let snp = new SideNavigationPage();
    await snp.NavigateToPhysLocation();
  }
  async OpenConfigureMenu() {
    let snp = new SideNavigationPage();
    await snp.ClickTopologyMenu();
  }
  async CreatePhysLocation(physlocation) {
    let result = false;
    let basePage = new BasePage();
    let snp = new SideNavigationPage();
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
  async SearchPhysLocation(physlocationName) {
    let result = false;
    let snp = new SideNavigationPage();
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
  async UpdatePhysLocation(physlocation) {
    let result = false;
    let basePage = new BasePage();

    switch (physlocation.description) {
      case "update physlocation region":
        await this.txtRegion.sendKeys(physlocation.Region + this.randomize);
        await basePage.ClickUpdate();
        break;
      default:
        result = undefined;
    }
    if (result = !undefined) {
        result = await basePage.GetOutputMessage().then(function (value) {
          if (physlocation.validationMessage == value) {
            return true;
          } else {
            return false;
          }
        })

      }
    return result;
  }
  async DeletePhysLocation(physlocation) {
    let result = false;
    let basePage = new BasePage();
    await this.btnDelete.click();
    await this.txtConfirmName.sendKeys(physlocation.Name + this.randomize);
    await basePage.ClickDeletePermanently();
    result = await basePage.GetOutputMessage().then(function (value) {
      if (physlocation.validationMessage == value) {
        return true;
      } else {
        return false;
      }
    })
    return result;
  }
}
