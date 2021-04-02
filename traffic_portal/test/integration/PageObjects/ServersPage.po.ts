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

import randomIpv6 from "random-ipv6";

import { BasePage } from './BasePage.po';
import {SideNavigationPage} from '../PageObjects/SideNavigationPage.po';
import { config, randomize } from '../config';

export class ServersPage extends BasePage {

  private btnSelectFormSubmit = element(by.buttonText('Submit'));
  private lnkExportAsCSV = element(by.xpath("//button[@title='Export as CSV']"));
  private btnrefresh = element(by.xpath("//button[@title='Refresh']"));
  private selectTableColumns = element(by.xpath("//div[@id='toggleColumns']"));
  private btnMore = element(by.xpath("//button[contains(text(),'More')]"));
  private lnkQueueCDN = element(by.xpath("//a[text()='Queue CDN Server Updates']"));
  private lnkClearCDN = element(by.xpath("//a[text()='Clear CDN Server Updates']"));
  private btnShowEntries = element(by.name('serversTable_length'));

  //private btnCreate = element(by.buttonText('Create'));
  private btnDelete = element(by.buttonText('Delete'));
  private txtUpdateStatus = element(by.xpath("//a[text()='Update Status']"));
  private lnkQueueServerUpdates = element(by.xpath("//a[text()='Queue Server Updates']"));
  private lnkClearServerUpdates = element(by.xpath("//a[text()='Clear Server Updates']"));
  private ViewConfigFiles = element(by.xpath("//a[text()='View Config Files']"));

  private lnkcacheGroupHeading = element(by.xpath("//th[text()='Cache Group']"));
  private lnkCDNHeading = element(by.xpath("//th[text()='CDN']"));
  private lnkDomainHeading = element(by.xpath("//th[text()='Domain']"));
  private lnkHostHeading = element(by.xpath("//th[text()='Host']"));
  private lnkILOIPHeading = element(by.xpath("//th[text()='ILO IP Address']"));
  private lnkIPv6AddressHeading = element(by.xpath("//th[text()='IPv6 Address']"));
  private lnknetworkIPHeading = element(by.xpath("//th[text()='Network IP']"));
  private lnkphysLocationHeading = element(by.xpath("//th[text()='Phys Location']"));
  private lnkstatusHeading = element(by.xpath("//th[text()='Status']"));
  private lnktypeHeading = element(by.xpath("//th[text()='Type']"));
  private lnkupdatePendingHeading = element(by.xpath("//th[text()='Update Pending']"));

  private txtStatus = element(by.name('status'));
  private txtUpdatePending = element(by.name('updPending'));
  private txtHostName = element(by.xpath("//ol[@class='breadcrumb pull-left']//li[@class='active ng-binding']"))
  private txtDomainName = element(by.name('domainName'));
  private txtProfile = element(by.name('profile'));
  private txtNetworkGateway = element(by.id("--gateway"));
  private txtNetworkMTU = element(by.id("-mtu"));
  private txtPhysLocation = element(by.name('physLocation'));
  private txtIp6Address = element(by.name('ip6Address'));
  private txtIp6Gateway = element(by.name('ip6Gateway'));
  private txtTcpPort = element(by.name('tcpPort'));
  private txtHttpsPort = element(by.name('httpsPort'));
  private txtRack = element(by.name('rack'));
  private txtMgmtIpAddress = element(by.name('mgmtIpAddress'));
  private txtMgmtIpNetmask = element(by.name('mgmtIpNetmask'));
  private txtMgmtIpGateway = element(by.name('mgmtIpGateway'));
  private txtIloIpAddress = element(by.name('iloIpAddress'));
  private txtIloIpNetmask = element(by.name('iloIpNetmask'));
  private txtIloIpGateway = element(by.name('iloIpGateway'));
  private txtIloUsername = element(by.name('iloUsername'));
  private txtIloPassword = element(by.name('iloPassword'));
  private txtRouterHostName = element(by.name('routerHostName'));
  private txtRouterPortName = element(by.name('routerPortName'));
  private lblInputError = element(by.className("input-error"));

  private txtHostname = element(by.name('hostName'));
  private txtCDN = element(by.name('cdn'));
  private txtCacheGroup = element(by.name('cachegroup'));
  private txtType = element(by.name('type'));
  private txtNetworkIP= element(by.id("-"));
  private txtNetworkSubnet = element(by.name('ipNetmask'));
  private txtSearch = element(by.id('serversTable_filter')).element(by.css('label input'));
  private mnuServerTable = element(by.id('serversTable'));
  private txtConfirmServerName = element(by.name('confirmWithNameInput'));

  private btnYesRemoveSC = element(by.buttonText("Yes"))
  private btnManageCapabilities = element(by.linkText('Manage Capabilities'));
  private btnAddCapabilities = element(by.name('addCapabilityBtn'));
  private selectCapabilities = element(by.name('selectFormDropdown'));
  private searchFilter = element(by.id('serverCapabilitiesTable_filter')).element(by.css('label input'));
  private btnManageDeliveryService = element(by.linkText('Manage Delivery Services'));
  private btnLinkDStoServer = element(by.xpath("//button[@title='Link Delivery Services to Server']"));
  private txtDSSearch = element(by.id('assignDSTable_filter')).element(by.css('label input'));
  private btnIpIsService = element(by.xpath("//input[@ng-model='ip.serviceAddress']"))
  private btnHostSearch = element(by.xpath("(//span[@class='ag-header-icon ag-header-cell-menu-button']//span)[4]"));
  private txtFilter= element(by.xpath("(//input[@placeholder='Filter...'])[1]"));
  private btnAddInterfaces = element(by.xpath("//button[@title='add a new interface']//i[1]"));
  private txtInterfaceName = element(by.id("-name"));
  private btnAddIp = element(by.name('addIPBtn'));
  private btnMoreCreateServer = element(by.name("moreBtn"))
  private btnCreateServer = element(by.name("createServerMenuItem"))
  private txtQuickSearch = element(by.id("quickSearch"));
  private readonly config = config;
  private randomize = randomize;

  async OpenServerPage(){
    let snp = new SideNavigationPage();
    await snp.NavigateToServersPage();
   }
   async OpenConfigureMenu(){
    let snp = new SideNavigationPage();
    await snp.ClickConfigureMenu();
   }
  GetInputErrorDisplayed() {
    return this.lblInputError.getText()
  }

  IsServersItemPresent(serversName: string) {
    return element(by.xpath("//table[@id='serversTable']//tr/td[text()='" + "']")).isPresent()
  }

  async ClickAddServer() {
    await this.btnCreateServer.click()
  }

  async CreateServer(server){
    let result = false;
    let basePage = new BasePage();
    let networkIp = Math.round(Math.random() * 100).toString()+ "." + Math.round(Math.random() * 100).toString() + "." + Math.round(Math.random() * 100).toString() +
    "." + Math.round(Math.random() * 100).toString();
    let ipv6 = randomIpv6();
    await this.btnMoreCreateServer.click();
    await this.btnCreateServer.click();
    await this.txtStatus.sendKeys(server.Status);
    await this.txtHostname.sendKeys(server.Hostname+this.randomize);
    await this.txtDomainName.sendKeys(server.Domainname);
    await this.txtCDN.sendKeys("ALL");
    await this.txtCDN.sendKeys(server.CDN + this.randomize);
    await this.txtCacheGroup.sendKeys(server.CacheGroup + this.randomize);
    await this.txtType.sendKeys(server.Type);
    await this.txtProfile.sendKeys(server.Profile + this.randomize);
    await this.txtPhysLocation.sendKeys(server.PhysLocation);
    await this.txtInterfaceName.sendKeys(server.InterfaceName);
    await element(by.id(""+server.InterfaceName+"-")).sendKeys(ipv6.toString());
    if (!await basePage.ClickCreate())
        result = false;
    await basePage.GetOutputMessage().then(function(value){
      if(server.validationMessage == value){
        result = true;
      }else{
        result = false;
      }
    })
    await this.OpenServerPage();
    return result;
  }

  async SearchServer(nameServer:string){
    let result = false;
    let name = nameServer+this.randomize;
    await this.txtQuickSearch.clear();
    await this.txtQuickSearch.sendKeys(name);
    await browser.actions().mouseMove(element(by.xpath("//span[text()='"+name+"']"))).perform();
    await browser.actions().doubleClick(element(by.xpath("//span[text()='"+name+"']"))).perform();
  }

  async SearchDeliveryServiceFromServerPage(name:string){
    let result = false;
    await this.txtDSSearch.clear();
    await this.txtDSSearch.sendKeys(name);
    if(await browser.isElementPresent(element(by.xpath("//td[@data-search='^"+name+"$']"))) == true){
      await element(by.xpath("//td[@data-search='^"+name+"$']")).click();
      result = true;
    }else{
      result = undefined;
    }
    return result;
  }

  async AddDeliveryServiceToServer(deliveryServiceName:string,outputMessage:string){
    let result = false;
    let basePage = new BasePage();
    let deliveryService = deliveryServiceName+this.randomize;
    let serverNameRandomized;
    await this.txtHostName.getText().then(function(value){
      serverNameRandomized = value;
    })
    let serverName = serverNameRandomized.replace(this.randomize,"")
    if(outputMessage.includes("delivery services assigned")){
      outputMessage = outputMessage.replace(serverName,serverNameRandomized)
    }
    if(outputMessage.includes("cannot assign")){
      let dsCapRequired = outputMessage.slice(112,118);
      outputMessage = outputMessage.replace(dsCapRequired,dsCapRequired+this.randomize)
      outputMessage = outputMessage.replace(serverName,serverNameRandomized)
    }
    await this.btnMore.click();
    if(await this.btnManageDeliveryService.isPresent() == true){
      await this.btnManageDeliveryService.click();
      await this.btnLinkDStoServer.click();
      if(await this.SearchDeliveryServiceFromServerPage(deliveryService) == true){
        await basePage.ClickSubmit();
         result = await basePage.GetOutputMessage().then(function(value){
          if(value == outputMessage){
            return true;
          }else{
            return false;
          }
        })
      }
    }else{
      result = undefined;
    }
    return result;

  }

  async AddServerCapabilitiesToServer(serverCapabilities){
    let result = false;
    let serverPage = new ServersPage();
    let basePage = new BasePage();
    let serverCapabilitiesName = serverCapabilities.ServerCapability + this.randomize;
    await this.btnMore.click();
    if((await this.btnManageCapabilities.isPresent()) == true){
      await this.btnManageCapabilities.click();
      await this.btnAddCapabilities.click();
      await this.selectCapabilities.sendKeys(serverCapabilitiesName);
      await basePage.ClickSubmit();
      result = await basePage.GetOutputMessage().then(function(value){
        if(serverCapabilities.validationMessage == value || value.includes(serverCapabilities.validationMessage)){
          result = true;
        }else{
          result = false;
        }
        return result;
      })
    }else{
      result = undefined;
    }
    await this.OpenServerPage();
    return result;
   }

   async SearchServerServerCapabilities(name:string){
    let result = false;
    await this.searchFilter.clear();
    await this.searchFilter.sendKeys(name);
    result = await element.all(by.repeater('sc in ::serverCapabilities')).filter(function(row){
      return row.element(by.name('name')).getText().then(function(val){
        return val === name;
      });
    }).first().getText().then(function(value){
      if(value == name){
        return true;
      }else{
        return false;
      }
    })
    return result;
   }


   async RemoveServerCapabilitiesFromServer(serverCapabilities:string,outputMessage:string){
    let result = false;
    let basePage = new BasePage();
    let serverCapabilitiesname = serverCapabilities+this.randomize;
    let url;
    await browser.getCurrentUrl().then((link) => {
      url = link.toString();
    })
    let serverNumber = url.substring(url.lastIndexOf('/') + 1);
    if(outputMessage.includes("cannot remove")){
      outputMessage = outputMessage.replace(serverCapabilities,serverCapabilitiesname)
      outputMessage = outputMessage.slice(0,56) + serverNumber + " " + outputMessage.slice(56);
    }
    await this.btnMore.click();
    if((await this.btnManageCapabilities.isPresent()) == true){
      await this.btnManageCapabilities.click();
      if(await this.SearchServerServerCapabilities(serverCapabilitiesname) == true){
        await element(by.xpath("//td[text()='"+ serverCapabilitiesname +"']/following-sibling::td/a[@title='Remove Server Capability']")).click();
      }
      await this.btnYesRemoveSC.click();
      result = await basePage.GetOutputMessage().then(function(value){
        if(outputMessage == value){
          return true;
        }else if(value.includes(outputMessage)){
          return true;
        }else{
          return false;
        }
      })
    }else{
      result = undefined;
    }
    await this.OpenServerPage();
    return result;
   }
  async UpdateServer(server){
    let result = false;
    let basePage = new BasePage();
    if(server.description.includes('change the cdn of a Server')){
      await this.txtCDN.sendKeys(server.CDN + this.randomize);
      await this.txtProfile.sendKeys(server.Profile + this.randomize)
      await basePage.ClickUpdate();
      result = await basePage.GetOutputMessage().then(function (value) {
          if (server.validationMessage == value) {
            return true;
          } else {
            return false;
          }
        })
        return result;
    }
  }
  async DeleteServer(server){
    let result = false;
    let basePage = new BasePage();
    let name = server.Name + this.randomize
    await this.btnDelete.click();
    await browser.wait(ExpectedConditions.visibilityOf(this.txtConfirmServerName), 1000);
    await this.txtConfirmServerName.sendKeys(name);
    if(await basePage.ClickDeletePermanently() == true){
      result = await basePage.GetOutputMessage().then(function(value){
        if(server.validationMessage == value){
          return true
        }else{
          return false;
        }
      })
    }else{
      await basePage.ClickCancel();
    }
    await this.OpenServerPage();
    return result;
   }
}
