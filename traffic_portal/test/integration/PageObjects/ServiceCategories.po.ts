import { ElementFinder, browser, by, element } from 'protractor'
import { async, delay } from 'q';
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';

export class ServiceCategoriesPage extends BasePage {

    private btnCreateServiceCategories = element(by.name("createServiceCategoryButton"));
    private txtSearch = element(by.id('serviceCategoriesTable_filter')).element(by.css('label input'));
    private txtName = element(by.id('name'));
    private txtTenant = element(by.name("tenantId"))

    private btnDelete = element(by.buttonText('Delete'));
    private txtConfirmName = element(by.name('confirmWithNameInput'));
    private config = require('../config');
    private randomize = this.config.randomize;

    async OpenServicesMenu() {
        let snp = new SideNavigationPage();
        await snp.ClickServicesMenu();
    }

    async OpenServiceCategoriesPage() {
        let snp = new SideNavigationPage();
        await snp.NavigateToServiceCategoriesPage();
    }
    async CreateServiceCategories(serviceCategories) {
        let result = false;
        let basePage = new BasePage();
        await this.btnCreateServiceCategories.click();
        await this.txtName.sendKeys(serviceCategories.Name + this.randomize);
        await basePage.ClickCreate();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (serviceCategories.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }
    async SearchServiceCategories(nameServiceCategories: string) {
        let name = nameServiceCategories + this.randomize;
        let result = false;
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
    async UpdateServiceCategories(serviceCategories) {
        let result = false;
        let basePage = new BasePage();
        switch (serviceCategories.description) {
            case "update service categories name":
                await this.txtName.clear();
                await this.txtName.sendKeys(serviceCategories.NewName + this.randomize);
                await basePage.ClickUpdate();
                break;
            default:
                result = undefined;
        }
        result = await basePage.GetOutputMessage().then(function (value) {
            if (serviceCategories.validationMessage == value || value.includes(serviceCategories.validationMessage)) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }
    async DeleteServiceCategories(serviceCategories) {
        let name = serviceCategories.Name + this.randomize;
        let result = false;
        let basePage = new BasePage();
        await this.btnDelete.click();
        await this.txtConfirmName.sendKeys(name);
        await basePage.ClickDeletePermanently();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (serviceCategories.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;

    }
}