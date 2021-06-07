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
import { browser, by, element} from 'protractor';
interface DeliveryServices {
  Type: string;
  Name: string;
  validationMessage: string;
}
interface DeleteDeliveryService {
  Name: string;
  validationMessage: string;
}
export class DeliveryServicePage extends BasePage {
  private btnMore = element(by.name("moreBtn"));
  private btnCreateNewDeliveryServices = element(by.linkText("Create Delivery Service"));
  private mnuFormDropDown = element(by.name('selectFormDropdown'));
  private btnSubmitFormDropDown = element(by.buttonText('Submit'));
  private txtSearch = element(by.id("quickSearch"))
  // private txtSearchRC = element(by.id('serverCapabilitiesTable_filter')).element(by.css('label input'));
  // private mnuDeliveryServicesTable = element(by.id('deliveryServicesTable'));
  private txtConfirmName = element(by.name('confirmWithNameInput'));
  private btnDelete = element(by.buttonText('Delete'));
  // private btnMore = element(by.name('moreBtn'));
  // private btnManageRequiredServerCapabilities = element(by.linkText('Manage Required Server Capabilities'));
  // private btnAddRequiredServerCapabilities = element(by.name('addCapabilityBtn'));
  // private txtSearchRequiredServerCapabilities = element(by.id('deliveryServiceCapabilitiesTable_filter')).element(by.css('label input'));
  // private btnRemoveRequiredServerCapabilities = element(by.css('a[title="Remove Required Server Capability"]'));
  // private btnYesRemoveRC = element(by.buttonText('Yes'));
  // private txtRequestStatus = element(by.name('requestStatus'));
  // private txtComment = element(by.name('comment'));
  // private mnuManageServers = element(by.linkText('Manage Servers'));
  // private btnLinkServersToDS = element(by.xpath("//button[@title='Link Servers to Delivery Service']"));
  // private txtSearchServersForDS = element(by.xpath("//div[@id='dsServersUnassignedTable_filter']//input[@type='search']"));
  // private txtSearchServerToUnlinkFromDS = element(by.xpath("//div[@id='dsServersTable_filter']//input[@type='search']"));
  // private mnuUnlinkServerFromDS = element(by.linkText("Unlink Server from Delivery Service"));
  // private txtServerAssign = element(by.id('dsServersUnassignedTable'));

  private txtXmlId = element(by.name('xmlId'));
  private txtDisplayName = element(by.name('displayName'));
  private selectActive = element(by.name('active'));
  private selectType = element(by.id('type'));
  private selectTenant = element(by.name('tenantId'));
  private selectCDN = element(by.name('cdn'));
  // private txtLongDesc = element(by.name('longDesc'));
  // private txtlongDesc2 = element(by.name('longDesc1'));
  // private txtlongDesc3 = element(by.name('longDesc2'));
  private txtOrgServerURL = element(by.name('orgServerFqdn'));
  private txtProtocol = element(by.name('protocol'));
  // private txtDscp = element(by.name('dscp'));
  // private txtProfile = element(by.name('profile'));
  // private txtInfoUrl = element(by.name('infoUrl'));
  // private txtCheckPath = element(by.name('checkPath'));
  // private txtOriginShield = element(by.name('originShield'));

  // Cache Configuration Settings 
  // private txtMaxOriginConnections = element(by.name('maxOriginConnections'));
  // private txtSigningAlgorithm = element(by.name('signingAlgorithm'));
  // private txtRangeRequestHandling = element(by.name('rangeRequestHandling'));
  // private txtQueryStringHandling = element(by.name('qstringIgnore'));
  // private txtEdgeHeaderRewrite = element(by.name('edgeHeaderRewrite'));
  // private txtMidHeaderRewrite = element(by.name('midHeaderRewrite'));
  // private txtRegexRemap = element(by.name('regexRemap'));
  private txtRemapText = element(by.name('remapText'));
  // private txtMultiSiteOrigin = element(by.name('multiSiteOrigin'));
  // private txtCacheURL = element(by.name('cacheurl'));

  //Routing Configuration Settings 
  // private txtRoutingName = element(by.name('routingName'));
  // private txtIpv6RoutingEnabled = element(by.name('ipv6RoutingEnabled'));
  // private txtSubnetEnabled = element(by.name('ecsEnabled'));
  // private txtGeoProvider = element(by.name('geoProvider'));
  // private txtMissLat = element(by.name('missLat'));
  // private txtMissLong = element(by.name('missLong'));
  // private txtGeoLimit = element(by.name('geoLimit'));
  // private txtGeoLimitRedirectURL = element(by.name('geoLimitRedirectURL'));
  // private txtDnsBypassIp = element(by.name('dnsBypassIp'));
  // private txtDnsBypassIpv6 = element(by.name('dnsBypassIp6'));
  // private txtDnsBypassCname = element(by.name('dnsBypassCname'));
  // private txtDnsBypassTtl = element(by.name('dnsBypassTtl'));
  // private txtMaxDnsAnswers = element(by.name('maxDnsAnswers'));
  // private txtCcrDnsTtl = element(by.name('ccrDnsTtl'));
  // private txtGlobalMaxMbps = element(by.name('globalMaxMbps'));
  // private txtGlobalMaxTps = element(by.name('globalMaxTps'));
  // private txtFqPacingRate = element(by.name('fqPacingRate'));
  // private txtTrResponseHeaders = element(by.name('trResponseHeaders'));
  // private txtTrRequestHeaders = element(by.name('trRequestHeaders'));
  // private txtDeepCachingType = element(by.name('deepCachingType'));
  // private txtInitialDispersion = element(by.name('initialDispersion'));
  // private txtAnonymousBlockingEnabled = element(by.name('anonymousBlockingEnabled'));
  // private txtRegionalGeoBlocking = element(by.name('regionalGeoBlocking'));
  // private txtHttpBypassFqdn = element(by.name('httpBypassFqdn'));
  // private txtConsistentHashRegex = element(by.name('consistentHashRegex'));
  // private txtConsistentHashQueryParam = element(by.name('consistentHashQueryParam'));
  private btnCreateDeliveryServices = element(by.xpath("//div[@class='pull-right']//button[text()='Create']"));
  // private btnFullfillRequest = element(by.xpath("//button[text()='Fulfill Request']"));
  // private txtSearchDSRequest = element(by.xpath("//div[@id='dsRequestsTable_filter']//input[@type='search']"));
  // private lnkCompleteRequest = element(by.css('a[title="Complete Request"]'));
  // private txtCommentCompleteDS = element(by.name("text"));
  // private btnDeleteRequest = element(by.xpath("//button[text()='Delete Request']"));
  // private txtNoMatchingError = element(by.xpath("//td[text()='No data available in table']"));
  private randomize = randomize;

  async OpenDeliveryServicePage(){
    let snp = new SideNavigationPage();
    await snp.NavigateToDeliveryServicesPage();
  }
  async OpenServicesMenu(){
    let snp = new SideNavigationPage();
    await snp.ClickServicesMenu();
  }
  public async CreateDeliveryService(deliveryservice: DeliveryServices): Promise<boolean> {
    let result = false;
    let type: string = deliveryservice.Type;
    let basePage = new BasePage();
    if (deliveryservice.validationMessage.includes("created")) {
      deliveryservice.validationMessage = deliveryservice.validationMessage.replace(deliveryservice.Name, deliveryservice.Name + this.randomize)
    }
    await this.btnMore.click();
    await this.btnCreateNewDeliveryServices.click();
    await this.mnuFormDropDown.sendKeys(type);
    await this.btnSubmitFormDropDown.click();
    switch (type) {
      case "ANY_MAP": {
        await this.txtXmlId.sendKeys(deliveryservice.Name + this.randomize);
        await this.txtDisplayName.sendKeys('cdntest' + this.randomize);
        await this.selectActive.sendKeys('Active')
        await this.selectType.sendKeys('ANY_MAP')
        await this.selectTenant.sendKeys('- root')
        await this.selectCDN.sendKeys('dummycdn')
        await this.txtRemapText.sendKeys('test')
        break;
      }
      case "DNS": {
        await this.txtXmlId.sendKeys(deliveryservice.Name + this.randomize);
        await this.txtDisplayName.sendKeys(deliveryservice.Name + this.randomize);
        await this.selectActive.sendKeys('Active')
        await this.selectType.sendKeys('DNS')
        await this.selectTenant.sendKeys('- root')
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
        await this.selectTenant.sendKeys('- root')
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
        await this.selectTenant.sendKeys('- root')
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
    let name = nameDS + this.randomize;
    await this.txtSearch.clear();
    await this.txtSearch.sendKeys(name);
    if (await browser.isElementPresent(element(by.xpath("//td[@data-search='^" + name + "$']"))) === true) {
      await element(by.xpath("//td[@data-search='^" + name + "$']")).click();
      return true;
    }
    return false;
  }
  public async DeleteDeliveryService(deliveryservice:DeleteDeliveryService): Promise<boolean> {
    let result = false;
    let basePage = new BasePage();
    if (deliveryservice.validationMessage.includes("deleted")) {
      deliveryservice.validationMessage = deliveryservice.validationMessage.replace(deliveryservice.Name, deliveryservice.Name + this.randomize);
    }
    await this.btnDelete.click();
    await this.txtConfirmName.sendKeys(deliveryservice.Name + this.randomize);
    await basePage.ClickDeletePermanently();
    result = await basePage.GetOutputMessage().then(value => value === deliveryservice.validationMessage);
    return result;
  }
}
