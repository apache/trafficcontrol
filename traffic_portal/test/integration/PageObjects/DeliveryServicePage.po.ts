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

import { BasePage } from './BasePage.po';
import { randomize } from "../config";
import { SideNavigationPage } from './SideNavigationPage.po';
import {browser, by, element} from 'protractor';

interface DeliveryServices {
  Type: string;
  Name: string;
  Tenant: string;
  validationMessage: string;
}
interface UpdateDeliveryService {
  description: string;
  Name: string;
  NewName: string;
  validationMessage: string;
}
interface DeleteDeliveryService {
  Name: string;
  validationMessage: string;
}
interface AssignServer {
  DSName: string;
  ServerName: string;
  validationMessage: string;
}
interface AssignRC {
  RCName: string;
  DSName: string;
  validationMessage: string;
}
export class DeliveryServicePage extends BasePage {
  private btnCreateNewDeliveryServices = element(by.buttonText("Create Delivery Service"));
  private mnuFormDropDown = element(by.name('selectFormDropdown'));
  private btnSubmitFormDropDown = element(by.buttonText('Submit'));
  private txtSearch = element(by.id("quickSearch"))
  private txtConfirmName = element(by.name('confirmWithNameInput'));
  private btnDelete = element(by.buttonText('Delete'));
  private btnMore = element(by.name('moreBtn'));
  private mnuManageRequiredServerCapabilities = element(by.linkText('Manage Required Server Capabilities'));
  private btnAddRequiredServerCapabilities = element(by.name('addCapabilityBtn'));
  private txtInputRC = element(by.name("selectFormDropdown"));
  private mnuManageServers = element(by.buttonText('Manage Servers'));
  private btnAssignServer = element(by.name("selectServersMenuItem"));
  private txtXmlId = element(by.name('xmlId'));
  private txtDisplayName = element(by.name('displayName'));
  private selectActive = element(by.name('active'));
  private selectType = element(by.id('type'));
  private selectTenant = element(by.name('tenantId'));
  private selectCDN = element(by.name('cdn'));
  private txtOrgServerURL = element(by.name('orgServerFqdn'));
  private txtProtocol = element(by.name('protocol'));
  private txtRemapText = element(by.name('remapText'));
  private btnCreateDeliveryServices = element(by.buttonText('Create'));
  private randomize = randomize;

  public async OpenDeliveryServicePage() {
    const snp = new SideNavigationPage();
    await snp.NavigateToDeliveryServicesPage();
  }

  public async OpenServicesMenu() {
    const snp = new SideNavigationPage();
    await snp.ClickServicesMenu();
  }

  public async CreateDeliveryService(deliveryservice: DeliveryServices): Promise<boolean> {
    let result = false;
    let type: string = deliveryservice.Type;
    const basePage = new BasePage();
    await this.btnMore.click();
    await this.btnCreateNewDeliveryServices.click();
    await this.mnuFormDropDown.sendKeys(type);
    await this.btnSubmitFormDropDown.click();
    switch (type) {
      case "ANY_MAP": {
        await this.txtXmlId.sendKeys(deliveryservice.Name + this.randomize);
        await this.txtDisplayName.sendKeys(deliveryservice.Name + this.randomize);
        await this.selectActive.sendKeys('Active')
        await this.selectType.sendKeys('ANY_MAP')
        await this.selectTenant.click();
        await element(by.name(deliveryservice.Tenant + this.randomize)).click();
        await this.selectCDN.sendKeys('dummycdn')
        await this.txtRemapText.sendKeys('test')
        break;
      }
      case "DNS": {
        await this.txtXmlId.sendKeys(deliveryservice.Name + this.randomize);
        await this.txtDisplayName.sendKeys(deliveryservice.Name + this.randomize);
        await this.selectActive.sendKeys('Active')
        await this.selectType.sendKeys('DNS')
        await this.selectTenant.click();
        await element(by.name(deliveryservice.Tenant + this.randomize)).click();
        await this.selectCDN.sendKeys('dummycdn')
        await this.txtOrgServerURL.sendKeys('http://origin.infra.ciab.test');
        await this.txtProtocol.sendKeys('HTTP')
        break;
      }
      case "HTTP": {
        await this.txtXmlId.sendKeys(deliveryservice.Name + this.randomize);
        await this.txtDisplayName.sendKeys(deliveryservice.Name + this.randomize);
        await this.selectActive.sendKeys('Active')
        await this.selectType.sendKeys('HTTP')
        await this.selectTenant.click();
        await element(by.name(deliveryservice.Tenant + this.randomize)).click();
        await this.selectCDN.sendKeys('dummycdn')
        await this.txtOrgServerURL.sendKeys('http://origin.infra.ciab.test');
        await this.txtProtocol.sendKeys('HTTP')
        break;
      }
      case "STEERING": {
        await this.txtXmlId.sendKeys(deliveryservice.Name + this.randomize);
        await this.txtDisplayName.sendKeys(deliveryservice.Name + this.randomize);
        await this.selectActive.sendKeys('Active')
        await this.selectType.sendKeys('STEERING')
        await this.selectTenant.click();
        await element(by.name(deliveryservice.Tenant + this.randomize)).click();
        await this.selectCDN.sendKeys('dummycdn')
        await this.txtProtocol.sendKeys('HTTP')
        break;
      }
      default:
        {
          console.log('Wrong Type name');
          break;
        }
    }
    await this.btnCreateDeliveryServices.click();
    result = await basePage.GetOutputMessage().then(value => value === deliveryservice.validationMessage);
    return result;
  }

  public async SearchDeliveryService(nameDS: string): Promise<boolean> {
    const name = nameDS + this.randomize;
    await this.txtSearch.clear();
    await this.txtSearch.sendKeys(name);
    const result = await element(by.cssContainingText("span", name)).isPresent();
    await element(by.cssContainingText("span", name)).click();
    return !result;
  }

  public async UpdateDeliveryService(deliveryservice: UpdateDeliveryService): Promise<boolean | undefined> {
    let result: boolean | undefined = false;
    const basePage = new BasePage();
    switch (deliveryservice.description) {
      case "update delivery service display name":
        await this.txtDisplayName.clear();
        await this.txtDisplayName.sendKeys(deliveryservice.NewName + this.randomize);
        await basePage.ClickUpdate();
        break;
      default:
        result = undefined;
    }
    if (result = !undefined) {
      result = await basePage.GetOutputMessage().then(value => value === deliveryservice.validationMessage);
    }
    return result;
  }

  public async DeleteDeliveryService(deliveryservice: DeleteDeliveryService): Promise<boolean> {
    let result = false;
    const basePage = new BasePage();
    if (deliveryservice.validationMessage.includes("deleted")) {
      deliveryservice.validationMessage = deliveryservice.validationMessage.replace(deliveryservice.Name, deliveryservice.Name + this.randomize);
    }
    await this.btnDelete.click();
    await this.txtConfirmName.sendKeys(deliveryservice.Name + this.randomize);
    await basePage.ClickDeletePermanently();
    result = await basePage.GetOutputMessage().then(value => value === deliveryservice.validationMessage);
    return result;
  }

  public async AssignServerToDeliveryService(deliveryservice: AssignServer): Promise<boolean>{
    let result = false;
    const basePage = new BasePage();
    await this.btnMore.click();
    await this.mnuManageServers.click();
    await this.btnMore.click();
    await this.btnAssignServer.click();
    await browser.sleep(3000);
    await element(by.cssContainingText(".ag-cell-value", deliveryservice.ServerName)).click();
    await this.ClickSubmit();
    result = await basePage.GetOutputMessage().then(value => value === deliveryservice.validationMessage);
    return result;
  }

  public async AssignRequiredCapabilitiesToDS(deliveryservice: AssignRC): Promise<boolean>{
    let result = false;
    const basePage = new BasePage();
    await this.btnMore.click();
    await this.mnuManageRequiredServerCapabilities.click();
    await this.btnAddRequiredServerCapabilities.click();
    await this.txtInputRC.sendKeys(deliveryservice.RCName);
    await this.ClickSubmit();
    result = await basePage.GetOutputMessage().then(value => value === deliveryservice.validationMessage);
    return result;
  }


}
