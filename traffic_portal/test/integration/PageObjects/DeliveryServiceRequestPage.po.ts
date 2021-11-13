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

import { browser, by, element, ExpectedConditions } from 'protractor';
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from '../PageObjects/SideNavigationPage.po';
import { randomize } from '../config';
interface DeliveryService{
  description: string;
  XmlId: string;
  Type: string;
  DisplayName: string;
  Active: string;
  ContentRoutingType: string;
  Tenant: string;
  CDN: string;
  RawRemapText: string;
  RequestStatus: string;
  validationMessage: string;
}

export class DeliveryServicesRequestPage extends BasePage {

  private btnCreateNewDeliveryServices = element(by.name('createDeliveryServiceButton'));
  private mnuFormDropDown = element(by.name('selectFormDropdown'));
  private btnSubmitFormDropDown = element(by.buttonText('Submit'));
  private txtSearch = element(by.id('deliveryServicesTable_filter')).element(by.css('label input'));
  private txtConfirmName = element(by.name('confirmWithNameInput'));
  private btnDelete = element(by.className('pull-right')).element(by.buttonText('Delete'));

  private btnYesRemoveRC = element(by.buttonText('Yes'));
  private txtRequestStatus = element(by.name('requestStatus'));
  private txtComment = element(by.name('comment'));

  private txtXmlId = element(by.name('xmlId'));
  private txtDisplayName = element(by.name('displayName'));
  private selectActive = element(by.name('active'));
  private selectType = element(by.id('type'));
  private selectTenant = element(by.name('tenantId'));
  private selectCDN = element(by.name('cdn'));

  // Cache Configuration Settings 
  private txtRemapText = element(by.name('remapText'));

  //Routing Configuration Settings 
  private btnCreateDeliveryServices = element(by.xpath("//div[@class='pull-right']//button[text()='Create']"));
  private btnFullfillRequest = element(by.xpath("//button[text()='Fulfill Request']"));
  private txtSearchDSRequest = element(by.xpath("//div[@id='dsRequestsTable_filter']//input[@type='search']"));
  private lnkCompleteRequest = element(by.css('a[title="Complete Request"]'));
  private txtCommentCompleteDS = element(by.name("text"));
  private txtNoMatchingError = element(by.xpath("//td[text()='No data available in table']"));
  private randomize = randomize;
  async OpenDeliveryServicePage() {
    let snp = new SideNavigationPage();
    await snp.NavigateToDeliveryServicesPage();
  }
  async OpenServicesMenu() {
    let snp = new SideNavigationPage();
    await snp.ClickServicesMenu();
  }

  async CreateDeliveryService(deliveryservice:DeliveryService) {
    let result = false;
    let type: string = deliveryservice.Type;
    let basePage = new BasePage();
    let snp = new SideNavigationPage();
    if (deliveryservice.validationMessage.includes("created")) {
      deliveryservice.validationMessage = deliveryservice.validationMessage.replace(deliveryservice.XmlId, deliveryservice.XmlId + this.randomize)
    }
    await snp.NavigateToDeliveryServicesPage();
    await this.btnCreateNewDeliveryServices.click();
    await this.mnuFormDropDown.sendKeys(type);
    await this.btnSubmitFormDropDown.click();
    await this.txtXmlId.sendKeys(deliveryservice.XmlId + this.randomize);
    await this.txtDisplayName.sendKeys(deliveryservice.DisplayName);
    await this.selectActive.sendKeys(deliveryservice.Active)
    await this.selectType.sendKeys(deliveryservice.ContentRoutingType)
    await this.selectTenant.sendKeys(deliveryservice.Tenant + this.randomize)
    await this.selectCDN.sendKeys(deliveryservice.CDN)
    await this.txtRemapText.sendKeys(deliveryservice.RawRemapText)
    await this.btnCreateDeliveryServices.click();
    await this.txtRequestStatus.sendKeys(deliveryservice.RequestStatus);
    await this.txtComment.sendKeys('test');
    await basePage.ClickSubmit();
    if (deliveryservice.RequestStatus.includes('Review and Deployment')) {
      if (await this.SearchDeliveryServiceRequest(deliveryservice.XmlId + this.randomize) == true) {
        await browser.actions().mouseMove(this.btnFullfillRequest).perform();
        await this.btnFullfillRequest.click();
        await this.btnYesRemoveRC.click();
        result = await basePage.GetOutputMessage().then(function (value) {
          if (deliveryservice.validationMessage == value) {
            return true;
          } else {
            return false;
          }
        })
        await snp.NavigateToDeliveryServicesRequestsPage();
        await this.txtSearchDSRequest.clear();
        await this.txtSearchDSRequest.sendKeys(deliveryservice.XmlId + this.randomize);
        await this.lnkCompleteRequest.click();
        await this.btnYesRemoveRC.click();
        await this.txtCommentCompleteDS.sendKeys('test');
        await basePage.ClickSubmit();
      }
    }
    return result;
  }
  async SearchDeliveryServiceRequest(name: string) {
    let result = false;
    await browser.wait(ExpectedConditions.visibilityOf(this.txtSearchDSRequest), 2000);
    await this.txtSearchDSRequest.clear();
    await this.txtSearchDSRequest.sendKeys(name);
    if (await browser.isElementPresent(element(by.xpath("//td[@data-search='^" + name + "$']"))) == true) {
      await element(by.xpath("//td[@data-search='^" + name + "$']")).click();
      result = true;
    }
    return result;
  }

  async SearchDeliveryService(nameDS: string) {
    let name = nameDS + this.randomize;
    await this.txtSearch.clear();
    await this.txtSearch.sendKeys(name);
    if (await this.txtNoMatchingError.isPresent() == true) {
      return undefined;
    } else {
      await element.all(by.repeater('ds in ::deliveryServices')).filter(function (row) {
        return row.element(by.name('xmlId')).getText().then(function (val) {
          return val === name;
        });
      }).first().click();
    }
  }


  async DeleteDeliveryService(nameDS: string, requestStatus: string, outputMessage: string) {
    let result = false;
    let basePage = new BasePage();
    let snp = new SideNavigationPage();
    let name = nameDS + this.randomize;
    if (outputMessage.includes("deleted")) {
      outputMessage = outputMessage.replace(nameDS, name);
    }
    await this.btnDelete.click();
    await this.txtConfirmName.sendKeys(name);
    await basePage.ClickDeletePermanently();
    await this.txtRequestStatus.sendKeys(requestStatus);
    await this.txtComment.sendKeys('test');
    await basePage.ClickSubmit();
    if (requestStatus.includes('Review and Deployment')) {
      if (await this.SearchDeliveryServiceRequest(name) == true) {
        await browser.actions().mouseMove(this.btnFullfillRequest).perform();
        await this.btnFullfillRequest.click();
        await this.btnYesRemoveRC.click();
        await this.txtConfirmName.sendKeys(name);
        await basePage.ClickDeletePermanently();
        result = await basePage.GetOutputMessage().then(function (value) {
          if (outputMessage == value) {
            return true;
          } else {
            return false;
          }
        })
        await snp.NavigateToDeliveryServicesRequestsPage();
        await this.txtSearchDSRequest.clear();
        await this.txtSearchDSRequest.sendKeys(name);
        await this.lnkCompleteRequest.click();
        await this.btnYesRemoveRC.click();
        await this.txtCommentCompleteDS.sendKeys('test');
        await basePage.ClickSubmit();
      }
    } else {
      result = await basePage.GetOutputMessage().then(function (value) {
        if (outputMessage == value) {
          return true;
        } else {
          return false;
        }
      })

    }
    return result;
  }
}
