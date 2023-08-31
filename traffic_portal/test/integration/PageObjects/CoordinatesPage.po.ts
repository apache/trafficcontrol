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
import { by, element, browser } from 'protractor';

import { randomize } from "../config";
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';

interface Coordinates {
    description?: string;
    Latitude: number;
    Longitude: number;
    Name: string;
    validationMessage?: string;
}

interface UpdateCoordinate {
    description: string;
    Latitude: number;
    validationMessage?: string;
}

interface DeleteCoordinate {
    Name: string;
    validationMessage?: string;
}

export class CoordinatesPage extends BasePage {

    private btnCreateNewCoordinates = element(by.xpath("//button[@title='Create Coordinate']"));
    private txtName = element(by.name('name'));
    private txtLatitude = element(by.name('latitude'))
    private txtLongitude = element(by.name('longitude'))
    private txtSearch = element(by.id('coordinatesTable_filter')).element(by.css('label input'));
    private btnDelete = element(by.buttonText('Delete'));
    private txtConfirmName = element(by.name('confirmWithNameInput'));
    private randomize = randomize;

    async OpenCoordinatesPage() {
        let snp = new SideNavigationPage();
        await snp.NavigateToCoordinatesPage();
    }
    async OpenTopologyMenu() {
        let snp = new SideNavigationPage();
        await snp.ClickTopologyMenu();
    }

    public async CreateCoordinates(coordinates: Coordinates): Promise<boolean> {
        let result = false;
        let basePage = new BasePage();
        await this.btnCreateNewCoordinates.click();
        await this.txtName.sendKeys(coordinates.Name + this.randomize)
        await this.txtLatitude.sendKeys(coordinates.Latitude);
        await this.txtLongitude.sendKeys(coordinates.Longitude)
        await basePage.ClickCreate();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (value.indexOf(coordinates.validationMessage ?? "") > -1) {
                return true;
            } else {
                return false;
            }
        })
        return result;

    }

    public async SearchCoordinates(nameCoordinates: string): Promise<boolean> {
        let snp = new SideNavigationPage();
        let name = nameCoordinates + this.randomize;
        await snp.NavigateToCoordinatesPage();
        await this.txtSearch.clear();
        await this.txtSearch.sendKeys(name);
        if (await browser.isElementPresent(element(by.xpath("//td[@data-search='^" + name + "$']"))) == true) {
            await element(by.xpath("//td[@data-search='^" + name + "$']")).click();
            return true;
        }
        return false;
    }

    async UpdateCoordinates(coordinates: UpdateCoordinate): Promise<boolean | undefined> {
        let result: boolean | undefined = false;
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
                if (value.indexOf(coordinates.validationMessage ?? "") > -1) {
                    return true;
                } else {
                    return false;
                }
            })

        }
        return result;
    }

    public async DeleteCoordinates(coordinates: DeleteCoordinate): Promise<boolean> {
        let result = false;
        let basePage = new BasePage();
        await this.btnDelete.click();
        await this.txtConfirmName.sendKeys(coordinates.Name + this.randomize);
        await basePage.ClickDeletePermanently();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (value.indexOf(coordinates.validationMessage ?? "") > -1) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }

}
