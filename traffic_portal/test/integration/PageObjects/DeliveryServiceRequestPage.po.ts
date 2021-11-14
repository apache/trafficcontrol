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

import {  browser, by, element} from 'protractor';
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from '../PageObjects/SideNavigationPage.po';
import { randomize } from '../config';


interface CreateDeliveryServiceRequest{
  description: string;
  XmlId: string;
  DisplayName: string;
  Active: string;
  ContentRoutingType: string;
  Tenant: string;
  CDN: string;
  RawText: string;
  validationMessage:string;
}
interface SearchDeliveryServiceRequest{
  XmlId: string;
}
interface FullfillDeliveryServiceRequest{
  XmlId: string;
  FullfillMessage: string;
}
interface CompleteDeliveryServiceRequest{
  XmlId: string;
  CompleteMessage: string;
}
interface DeleteDeliveryServiceRequest{
  XmlId: string;
  DeleteMessage: string;
}
interface UpdateDeliveryServiceRequest{
  XmlId: string;
  UpdateMessage: string;
}

export class DeliveryServicesRequestPage extends BasePage {
  private btnFullfillRequest = element(by.buttonText("Fulfill Request"))
  private btnCompleteRequest = element(by.buttonText("Complete Request"))
  private btnDeleteRequest = element(by.buttonText("Delete Request"))
  private btnUpdateRequest = element(by.buttonText("Update Request"))
  private btnYes = element(by.buttonText("Yes"))
  private btnMore = element(by.name("moreBtn"))
  private btnCreateDS = element(by.linkText("Create Delivery Service"));
  private formDropDown = element(by.name("selectFormDropdown"))
  private txtXmlId = element(by.id("xmlId"));
  private txtDisplayName = element(by.id("displayName"));
  private txtActive = element(by.id("active"));
  private txtContentRoutingType = element(by.id("type"));
  private txtTenant = element(by.id("tenantId"));
  private txtCDN = element(by.id("cdn"));
  private txtRawRemapText = element(by.id("remapText"));
  private txtRequestStatus = element(by.name("requestStatus"))
  private txtComment = element(by.name("comment"))
  private txtQuickSearch = element(by.id("quickSearch"))
  private txtConfirmInput = element(by.name("confirmWithNameInput"))
  private randomize = randomize;
  public async OpenDeliveryServiceRequestPage(){
    const snp = new SideNavigationPage();
    await snp.NavigateToDeliveryServicesRequestsPage();
  }
  public async OpenServicesMenu(){
    const snp = new SideNavigationPage();
    await snp.ClickServicesMenu();
  }
  public async OpenDeliveryServicePage(){
    const snp = new SideNavigationPage();
    await snp.NavigateToDeliveryServicesPage();
  }
  public async CreateDeliveryServiceRequest(deliveryservicerequest: CreateDeliveryServiceRequest){
    const basePage = new BasePage();
    const outPutMessage = deliveryservicerequest.validationMessage.replace(deliveryservicerequest.XmlId,deliveryservicerequest.XmlId+this.randomize)
    await this.btnMore.click();
    await this.btnCreateDS.click();
    await this.formDropDown.sendKeys("ANY_MAP");
    await basePage.ClickSubmit();
    await this.txtXmlId.sendKeys(deliveryservicerequest.XmlId + this.randomize);
    await this.txtDisplayName.sendKeys(deliveryservicerequest.DisplayName + this.randomize);
    await this.txtActive.sendKeys(deliveryservicerequest.Active);
    await this.txtContentRoutingType.sendKeys(deliveryservicerequest.ContentRoutingType);
    await this.txtTenant.sendKeys(deliveryservicerequest.Tenant);
    await this.txtCDN.sendKeys(deliveryservicerequest.CDN);
    await this.txtRawRemapText.sendKeys(deliveryservicerequest.RawText);
    await basePage.ClickCreate();
    await this.txtRequestStatus.sendKeys("Submit Request for Review and Deployment");
    await this.txtComment.sendKeys("test");
    await basePage.ClickSubmit();
    return await basePage.GetOutputMessage().then(value => outPutMessage === value);
  }
  public async SearchDeliveryServiceRequest(deliveryservicerequest: SearchDeliveryServiceRequest){
    const name = deliveryservicerequest.XmlId+this.randomize;
    await this.txtQuickSearch.clear();
    await this.txtQuickSearch.sendKeys(name)
    await browser.actions().click(element(by.cssContainingText("span", name))).perform();
  }
  public async FullFillDeliveryServiceRequest(deliveryservicerequest: FullfillDeliveryServiceRequest): Promise<boolean>{
    const basePage = new BasePage();
    const outPutMessage = deliveryservicerequest.FullfillMessage.replace(deliveryservicerequest.XmlId,deliveryservicerequest.XmlId+this.randomize)
    await this.btnFullfillRequest.click();
    await this.btnYes.click();
    return await basePage.GetOutputMessage().then(value => outPutMessage === value);
  }
  public async CompleteDeliveryServiceRequest(deliveryservicerequest: CompleteDeliveryServiceRequest): Promise<boolean>{
    const basePage = new BasePage();
    await this.btnCompleteRequest.click();
    await this.btnYes.click();
    return await basePage.GetOutputMessage().then(value => deliveryservicerequest.CompleteMessage === value);

  }
  public async DeleteDeliveryServiceRequest(deliveryservicerequest: DeleteDeliveryServiceRequest): Promise<boolean>{
    const basePage = new BasePage();
    await this.btnDeleteRequest.click();
    await this.txtConfirmInput.sendKeys(deliveryservicerequest.XmlId+this.randomize+" request");
    await basePage.ClickDeletePermanently();
    return await basePage.GetOutputMessage().then(value=>deliveryservicerequest.DeleteMessage === value);
  }
  public async UpdateDeliveryServiceRequest(deliveryservicerequest: UpdateDeliveryServiceRequest): Promise<boolean>{
    const basePage = new BasePage();
    const outPutMessage = deliveryservicerequest.UpdateMessage.replace(deliveryservicerequest.XmlId,deliveryservicerequest.XmlId+this.randomize)
    await this.txtRawRemapText.clear();
    await this.txtRawRemapText.sendKeys("change");
    await this.btnUpdateRequest.click();
    await this.txtRequestStatus.sendKeys("Submit for Review / Deployment")
    await this.txtComment.sendKeys("test");
    await basePage.ClickSubmit();
    return await basePage.GetOutputMessage().then(value => outPutMessage === value);
  }
}

  