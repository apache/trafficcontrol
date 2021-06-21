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
import { browser, by, element } from 'protractor';
import { randomize } from "../config";
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';

interface CDN {
  description?: string;
  DNSSEC: string;
  Domain: string;
  Name: string;
  validationMessage?: string;
}

interface UpdateCDN {
  description: string;
  Name: string;
  NewName: string;
  validationMessage?: string;
}

interface DeleteCDN {
  Name: string;
  validationMessage?: string;
}

export class CDNPage extends BasePage {

  private btnNewCDN = element(by.name('createCdnButton'));
  private txtCDNName = element(by.name('name'));
  private txtDomain = element(by.name('domainName'));
  private selectDNSSEC = element(by.name('dnssecEnabled'));
  private txtSearch = element(by.id('cdnsTable_filter')).element(by.css('label input'));
  private btnDelete = element(by.buttonText('Delete'));
  private txtConfirmName = element(by.name('confirmWithNameInput'));
  private btnDiffSnapshot = element(by.xpath("//button[@title='Diff CDN Snapshot']"));
  private btnYes = element((by.xpath("//button[text()='Yes']")));
  private btnQueueUpdates = element((by.xpath("//button[contains(text(),'Queue Updates')]")));
  private randomize = randomize;

  public async OpenCDNsPage(): Promise<void> {
    let snp = new SideNavigationPage();
    await snp.NavigateToCDNPage();
  }

  public async CreateCDN(cdn: CDN): Promise<boolean> {
    let snp = new SideNavigationPage();
    let basePage = new BasePage();
    await snp.NavigateToCDNPage();
    await this.btnNewCDN.click();
    await this.txtCDNName.sendKeys(cdn.Name + this.randomize);
    await this.txtDomain.sendKeys(cdn.Domain);
    await this.selectDNSSEC.sendKeys(cdn.DNSSEC);
    await basePage.ClickCreate();
    return await basePage.GetOutputMessage().then(value => cdn.validationMessage === value);
  }

  async SearchCDN(nameCDN: string) {
    let snp = new SideNavigationPage();
    let name = nameCDN + this.randomize;
    await snp.NavigateToCDNPage();
    await this.txtSearch.clear();
    await this.txtSearch.sendKeys(name);
    await element.all(by.repeater('cdn in ::cdns')).filter(function (row) {
      return row.element(by.name('name')).getText().then(function (val) {
        return val === name;
      });
    }).first().click();
  }

  public async UpdateCDN(cdn: UpdateCDN): Promise<boolean | undefined> {
    let result: boolean | undefined = false;
    let basePage = new BasePage();
    switch (cdn.description) {
      case 'perform snapshot':
        await this.btnDiffSnapshot.click();
        if (await browser.isElementPresent(element(by.xpath('//button[@title="Perform ' + cdn.Name + this.randomize + ' Snapshot"]')))) {
          await element(by.xpath('//button[@title="Perform ' + cdn.Name + this.randomize + ' Snapshot"]')).click();
        } else {
          throw new Error("Cannot find Perform Snapshot button")
        }
        await this.btnYes.click();
        break;
      case 'queue CDN updates':
        await this.btnQueueUpdates.click();
        if (await browser.isElementPresent(element(by.linkText('Queue ' + cdn.Name + this.randomize + ' Server Updates')))) {
          await element(by.linkText('Queue ' + cdn.Name + this.randomize + ' Server Updates')).click();
        } else {
          throw new Error("Cannot find Queue CDN updates button")
        }
        await this.btnYes.click();
        break;
      case 'clear CDN updates':
        await this.btnQueueUpdates.click();
        if (await browser.isElementPresent(element(by.linkText('Clear ' + cdn.Name + this.randomize + ' Server Updates')))) {
          await element(by.linkText('Clear ' + cdn.Name + this.randomize + ' Server Updates')).click();
        } else {
          throw new Error("Cannot find Clear CDN updates button")
        }
        await this.btnYes.click();
        break;
      case 'update cdn name':
        await this.txtCDNName.clear();
        await this.txtCDNName.sendKeys(cdn.NewName + this.randomize);
        await this.ClickUpdate();
      default:
        result = undefined;
    }
    result = await basePage.GetOutputMessage().then(function (value) {
      if (cdn.validationMessage == value) {
        return true;
      } else {
        return false;
      }
    })
    return result;
  }

  public async DeleteCDN(cdn: DeleteCDN): Promise<boolean> {
    let name = cdn.Name + this.randomize;
    let basePage = new BasePage();
    await this.btnDelete.click();
    await this.txtConfirmName.sendKeys(name);
    await basePage.ClickDeletePermanently();
    return await basePage.GetOutputMessage().then(value => cdn.validationMessage === value);
  }
  public async CheckCSV(name:string): Promise<boolean> {
    let result = false;
    let linkName = name;
    if (await browser.isElementPresent(element(by.xpath("//span[text()='" + linkName + "']"))) == true) {
      result = true;
    }
    return result;
  }
}
