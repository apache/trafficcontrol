import { ElementFinder, browser, by, element } from 'protractor'
import { async, delay } from 'q';
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';
export class RegionsPage extends BasePage {
    private btnCreateNewRegions = element(by.name('createRegionButton'));
    private txtSearch = element(by.id('regionsTable_filter')).element(by.css('label input'));
    private txtName = element(by.id('name'));
    private txtDivision = element(by.name('division'));
    private btnDelete = element(by.xpath("//button[text()='Delete']"));
    private txtConfirmName = element(by.name('confirmWithNameInput'));
    private config = require('../config');
    private randomize = this.config.randomize;
    async OpenRegionsPage(){
        let snp = new SideNavigationPage();
        await snp.NavigateToRegionsPage();
    }
    async OpenTopologyMenu(){
        let snp = new SideNavigationPage();
        await snp.ClickTopologyMenu();
    }

    async CreateRegions(regions){
        let result = false;
        let basePage = new BasePage();
        let snp = new SideNavigationPage();
        await snp.NavigateToRegionsPage();
        await this.btnCreateNewRegions.click();
        await this.txtName.sendKeys(regions.Name + this.randomize);
        await this.txtDivision.sendKeys(regions.Division);
        await basePage.ClickCreate();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (regions.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }
    async SearchRegions(nameRegions:string){
        let name = nameRegions + this.randomize;
        let result = false;
        let snp = new SideNavigationPage();
        await snp.NavigateToRegionsPage();
        await this.txtSearch.clear();
        await this.txtSearch.sendKeys(name);
        if (await browser.isElementPresent(element(by.xpath("//td[@data-search='^" + name + "$']"))) == true) {
            await element(by.xpath("//td[@data-search='^" + name + "$']")).click();
            result = true;
        } else {
            result = undefined;
        }
        return result;
    }
    async UpdateRegions(regions){
        let result = false;
        let basePage = new BasePage();
        switch(regions.description){
            case "update Region's Division":
                await this.txtDivision.sendKeys(regions.Division);
                await basePage.ClickUpdate();
                break;
            default:
                result = undefined;
        }
        result = await basePage.GetOutputMessage().then(function (value) {
            if (regions.validationMessage == value) {
              return true;
            } else {
              return false;
            }
          })
          return result;

    }
    async DeleteRegions(regions){
        let name = regions.Name + this.randomize;
        let result = false;
        let basePage = new BasePage();
        await this.btnDelete.click();
        await this.txtConfirmName.sendKeys(name);
        await basePage.ClickDeletePermanently();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (regions.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }
}