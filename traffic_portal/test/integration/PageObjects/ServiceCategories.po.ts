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
import { by, element } from "protractor";

import { randomize } from "../config";
import { SideNavigationPage } from "./SideNavigationPage.po";

interface CreateServiceCategory {
    Name: string;
}

interface UpdateServiceCategory {
    description: string;
    NewName: string;
}

interface DeleteServiceCategory {
    Name: string;
}

export class ServiceCategoriesPage extends SideNavigationPage {
    private txtName = element(by.name("name"));
    async OpenServicesMenu() {
        await this.ClickServicesMenu();
    }

    /**
     * Navigates the browser to the Service Categories table page.
     */
    async OpenServiceCategoriesPage() {
        await this.NavigateToServiceCategoriesPage();
    }

    public async CreateServiceCategories(
        serviceCategories: CreateServiceCategory,
        outputMessage: string
    ): Promise<boolean> {
        await this.OpenServiceCategoriesPage();
        await element(by.buttonText("More")).click();
        await element(by.linkText("Create New Service Category")).click();
        this.txtName.sendKeys(serviceCategories.Name + randomize);
        await this.ClickCreate();
        return this.GetOutputMessage().then(
            (v) => v.indexOf(outputMessage ?? "") > -1
        );
    }

    public async SearchServiceCategories(
        nameServiceCategories: string
    ): Promise<void> {
        nameServiceCategories += randomize;
        await this.OpenServiceCategoriesPage();
        const searchInput = element(by.id("quickSearch"));
        await searchInput.clear();
        await searchInput.sendKeys(nameServiceCategories);
        await element(
            by.cssContainingText("span", nameServiceCategories)
        ).click();
    }

    public async UpdateServiceCategories(
        serviceCategories: UpdateServiceCategory,
        outputMessage: string
    ): Promise<boolean | undefined> {
        switch (serviceCategories.description) {
            case "update service categories name":
                await this.txtName.clear();
                await this.txtName.sendKeys(
                    serviceCategories.NewName + randomize
                );
                await this.ClickUpdate();
                break;
            default:
                return undefined;
        }
        return await this.GetOutputMessage().then(
            (v) =>
                outputMessage === v ||
                (outputMessage !== undefined && v.includes(outputMessage))
        );
    }

    public async DeleteServiceCategories(
        serviceCategories: DeleteServiceCategory,
        outputMessage: string
    ): Promise<boolean> {
        const name = serviceCategories.Name + randomize;
        await element(by.buttonText("Delete")).click();
        await element(by.name("confirmWithNameInput")).sendKeys(name);
        await this.ClickDeletePermanently();
        return this.GetOutputMessage().then(
            (v) => v.indexOf(outputMessage ?? "") > -1
        );
    }
}
