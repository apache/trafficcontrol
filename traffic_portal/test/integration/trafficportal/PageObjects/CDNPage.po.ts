import { ElementFinder, browser, by, element, ExpectedConditions, protractor } from 'protractor';
import { async, delay } from 'q';
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';

export class CDNPage extends BasePage {

  private btnNewCDN = element(by.name('createCdnButton'));
  private txtCDNName = element(by.name('name'));
  private txtDomain = element(by.name('domainName'));
  private selectDNSSEC = element(by.name('dnssecEnabled'));
  private txtSearch = element(by.id('cdnsTable_filter')).element(by.css('label input'));
  private mnuCDNTable = element(by.id('cdnsTable'));
  private btnDelete = element(by.buttonText('Delete'));
  private txtConfirmName = element(by.name('confirmWithNameInput'));
  private btnDiffSnapshot = element(by.xpath("//button[@title='Diff CDN Snapshot']"));
  private btnYes = element((by.xpath("//button[text()='Yes']")));
  private btnQueueUpdates = element((by.xpath("//button[contains(text(),'Queue Updates')]")));
  private config = require('../config');
  private randomize = this.config.randomize;
  async OpenCDNsPage() {
    let snp = new SideNavigationPage();
    await snp.NavigateToCDNPage();
  }
  async CreateCDN(cdn) {
    let result = false;
    let snp = new SideNavigationPage();
    let basePage = new BasePage();
    await snp.NavigateToCDNPage();
    await this.btnNewCDN.click();
    await this.txtCDNName.sendKeys(cdn.Name + this.randomize);
    await this.txtDomain.sendKeys(cdn.Domain);
    await this.selectDNSSEC.sendKeys(cdn.DNSSEC);
    await basePage.ClickCreate();
    result = await basePage.GetOutputMessage().then(function (value) {
      if (cdn.validationMessage == value) {
        return true;
      } else {
        return false;
      }
    })
    return result;
  }

  async SearchCDN(nameCDN: string) {
    let result = false;
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

  async UpdateCDN(cdn) {
    let result = false;
    let snp = new SideNavigationPage();
    let basePage = new BasePage();
    switch (cdn.description) {
      case 'perform snapshot':
        await this.btnDiffSnapshot.click();
        await element(by.xpath(`//button[@title="Perform ` + cdn.Name + this.randomize + ` Snapshot"]`)).click();
        await this.btnYes.click();
        break;
      case 'queue CDN updates':
        await this.btnQueueUpdates.click();
        await element(by.linkText(`Queue ` + cdn.Name + this.randomize + ` Server Updates`)).click();
        await this.btnYes.click();
        break;
      case 'clear CDN updates':
        await this.btnQueueUpdates.click();
        await element(by.linkText(`Clear ` + cdn.Name + this.randomize + ` Server Updates`)).click();
        await this.btnYes.click();
        break;
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
  async DeleteCDN(cdn) {
    let name = cdn.Name + this.randomize;
    let result = false;
    let basePage = new BasePage();
    await this.btnDelete.click();
    await this.txtConfirmName.sendKeys(name);
    await basePage.ClickDeletePermanently();
    result = await basePage.GetOutputMessage().then(function (value) {
      if (cdn.validationMessage == value) {
        return true;
      } else {
        return false;
      }
    })
    return result;
  }
}