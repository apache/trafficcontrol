import { by, element, browser } from 'protractor'
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';
export class CoordinatesPage extends BasePage {

    private btnCreateNewCoordinates = element(by.xpath("//button[@title='Create Coordinate']"));
    private txtName = element(by.name('name'));
    private txtLatitude = element(by.name('latitude'))
    private txtLongitude = element(by.name('longitude'))
    private txtSearch = element(by.id('coordinatesTable_filter')).element(by.css('label input'));
    private btnDelete = element(by.buttonText('Delete'));
    private btnYes = element(by.buttonText('Yes'));
    private txtConfirmName = element(by.name('confirmWithNameInput'));
    private config = require('../config');
    private randomize = this.config.randomize;

    async OpenCoordinatesPage() {
        let snp = new SideNavigationPage();
        await snp.NavigateToCoordinatesPage();
    }
    async OpenTopologyMenu() {
        let snp = new SideNavigationPage();
        await snp.ClickTopologyMenu();
    }

    async CreateCoordinates(coordinates) {
        let result = false;
        let basePage = new BasePage();
        await this.btnCreateNewCoordinates.click();
        await this.txtName.sendKeys(coordinates.Name + this.randomize)
        await this.txtLatitude.sendKeys(coordinates.Latitude);
        await this.txtLongitude.sendKeys(coordinates.Longitude)
        await basePage.ClickCreate();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (coordinates.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;

    }

    async SearchCoordinates(nameCoordinates: string) {
        let result = false;
        let snp = new SideNavigationPage();
        let name = nameCoordinates + this.randomize;
        await snp.NavigateToCoordinatesPage();
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
    async UpdateCoordinates(coordinates) {
        let result = false;
        let basePage = new BasePage();
        switch (coordinates.description) {
            case "update coordinates latitude":
                await this.txtLatitude.clear();
                await this.txtLatitude.sendKeys(coordinates.Latitude);
                await basePage.ClickUpdate();
                break;
            default:
                result = undefined;
        }
        if (result = !undefined) {
            result = await basePage.GetOutputMessage().then(function (value) {
                if (coordinates.validationMessage == value) {
                    return true;
                } else {
                    return false;
                }
            })

        }
        return result;
    }
    async DeleteCoordinates(coordinates) {
        let result = false;
        let basePage = new BasePage();
        await this.btnDelete.click();
        await this.txtConfirmName.sendKeys(coordinates.Name + this.randomize);
        await basePage.ClickDeletePermanently();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (coordinates.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }

}