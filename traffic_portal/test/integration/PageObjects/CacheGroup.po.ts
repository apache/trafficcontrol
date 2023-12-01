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
import { browser, by, element } from "protractor";

import { SideNavigationPage } from "./SideNavigationPage.po";
import { randomize } from "../config";

interface CreateCacheGroup {
    Type: string;
    Name: string;
    ShortName: string;
    Latitude: string;
    Longitude: string;
    ParentCacheGroup: string;
    SecondaryParentCG: string;
    FailoverCG?: string;
}

interface UpdateCacheGroup {
    Type: string;
    FailoverCG?: string;
}

export class CacheGroupPage extends SideNavigationPage {
    private txtName = element(by.name("name"));

    async OpenTopologyMenu() {
        await this.ClickTopologyMenu();
    }

    /**
     * Navigates the browser to the Cache Groups table page.
     */
    async OpenCacheGroupsPage() {
        await this.NavigateToCacheGroupsPage();
    }

    /**
     * Creates a given Cache Group.
     *
     * @param cachegroup The CacheGroup to create.
     * @param outputMessage The expected output message
     * @returns Whether or not creation succeeded, which is judged by comparing
     * the displayed Alert message to outputMessage - they must match
     * *exactly*!
     */
    public async CreateCacheGroups(
        cachegroup: CreateCacheGroup,
        outputMessage: string
    ): Promise<boolean> {
        await this.OpenCacheGroupsPage();
        await element(by.buttonText("More")).click();
        await element(by.linkText("Create New Cache Group")).click();

        if (
            cachegroup.Type == "EDGE_LOC" &&
            cachegroup.FailoverCG === undefined
        ) {
            throw new Error(
                `cachegroups with Type 'EDGE_LOC' must have FailoverCG`
            );
        }

        const actions = [
            this.txtName.sendKeys(cachegroup.Name + randomize),
            element(by.name("shortName")).sendKeys(cachegroup.ShortName),
            element(by.name("type")).sendKeys(cachegroup.Type),
            element(by.name("latitude")).sendKeys(cachegroup.Latitude),
            element(by.name("longitude")).sendKeys(cachegroup.Longitude),
            element(by.name("parentCacheGroup")).sendKeys(
                cachegroup.ParentCacheGroup
            ),
            element(by.name("secondaryParentCacheGroup")).sendKeys(
                cachegroup.SecondaryParentCG
            ),
        ];

        if (cachegroup.Type == "EDGE_LOC" && cachegroup.FailoverCG) {
            actions.push(
                element(by.name("fallbackOptions")).sendKeys(
                    cachegroup.FailoverCG
                )
            );
        }

        await Promise.all(actions);
        await this.ClickCreate();
        return this.GetOutputMessage().then((v) => outputMessage === v);
    }

    /**
     * Searches the Cache Groups table for a specific Cache Group, and navigates to its details
     * page.
     *
     * @param nameCG The Name of the Cache Group for which to search.
     */
    public async SearchCacheGroups(nameCG: string): Promise<void> {
        nameCG += randomize;
        await this.OpenCacheGroupsPage();
        const searchInput = element(by.id("quickSearch"));
        await searchInput.clear();
        await searchInput.sendKeys(nameCG);
        await element(by.cssContainingText("span", nameCG)).click();
    }

    /**
     * Updates a Cache Group's Name.
     *
     * @param cachegroup A definition of the CacheGroup renaming.
     * @param outputMessage The expected output message
     * @returns Whether or not renaming succeeded.
     */
    public async UpdateCacheGroups(
        cachegroup: UpdateCacheGroup,
        outputMessage: string | undefined
    ): Promise<boolean | undefined> {
        let result: boolean | undefined = false;
        if (cachegroup.Type == "EDGE_LOC") {
            const name = cachegroup.FailoverCG + randomize;
            await element(by.name("fallbackOptions")).click();
            if (
                await browser.isElementPresent(
                    element(
                        by.css(
                            `select[name="fallbackOptions"] > option[label="${name}"]`
                        )
                    )
                )
            ) {
                await element(
                    by.css(
                        `select[name="fallbackOptions"] > option[label="${name}"]`
                    )
                ).click();
            } else {
                result = undefined;
            }
        }
        await element(by.name("type")).sendKeys(cachegroup.Type);
        await this.ClickUpdate();
        if (result !== undefined) {
            return (await this.GetOutputMessage()) === outputMessage;
        }
    }

    /**
     * Deletes a Cache Group.
     *
     * @param nameCG The Name of the Cache Group to be deleted.
     * @param outputMessage The expected output message
     * @returns Whether or not the deletion succeeded.
     */
    public async DeleteCacheGroups(
        nameCG: string,
        outputMessage: string
    ): Promise<boolean> {
        nameCG += randomize;
        await element(by.buttonText("Delete")).click();
        await element(by.name("confirmWithNameInput")).sendKeys(nameCG);
        await this.ClickDeletePermanently();
        return (await this.GetOutputMessage()) === outputMessage;
    }
}
