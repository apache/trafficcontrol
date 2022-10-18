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
import { randomize, twoNumberRandomize } from "../config";
import { SideNavigationPage } from "./SideNavigationPage.po";

interface CDN {
	description?: string;
	DNSSEC: string;
	Domain: string;
	Name: string;
	TTLOverride?: string;
	validationMessage?: string;
}

interface UpdateCDN {
	description: string;
	Name: string;
	NewName: string;
	validationMessage?: string;
}

export class CDNPage extends SideNavigationPage {

	private readonly txtCDNName = element(by.name("name"));
	private readonly txtCDNTTLOverride = element(by.name("ttlOverride"));

	/**
	 * Navigates the browser to the CDNs table page.
	 */
	public async openCDNsPage(): Promise<void> {
		await this.NavigateToCDNPage();
	}

	/**
	 * Creates a given CDN.
	 *
	 * @param cdn The CDN to create.
	 * @returns Whether or not creation succeeded, which is judged by comparing
	 * the displayed Alert messge to `cdn`'s `validtionMessge` - they must mtch
	 * *exactly*!
	 */
	public async createCDN(cdn: CDN): Promise<boolean> {
		await this.openCDNsPage();
		await element(by.buttonText("More")).click();
		await element(by.linkText("Create New CDN")).click();
		const actions = [
			this.txtCDNName.sendKeys(cdn.Name + randomize),
			element(by.name("domainName")).sendKeys(cdn.Domain),
			element(by.name("dnssecEnabled")).sendKeys(cdn.DNSSEC),
		];
		if (cdn.TTLOverride !== undefined) {
			actions.push(element(this.txtCDNTTLOverride).sendKeys(cdn.TTLOverride));
		}
		await Promise.all(actions);
		await this.ClickCreate();
		return this.GetOutputMessage().then(v => cdn.validationMessage === v);
	}

	/**
	 * Searches the CDNs table for a specific CDN, and navigates to its details
	 * page.
	 *
	 * @param cdnName The Name of the CDN for which to search.
	 */
	public async searchCDN(cdnName: string): Promise<void> {
		cdnName += randomize;
		await this.NavigateToCDNPage();
		const searchInput = element(by.id("quickSearch"));
		await searchInput.clear();
		await searchInput.sendKeys(cdnName);
		await element(by.cssContainingText("span", cdnName)).click();
	}

	/**
	 * Updates a CDN's Name.
	 *
	 * @param cdn A definition of the CDN renaming.
	 * @returns Whether or not renaming succeeded.
	 */
	public async updateCDN(cdn: UpdateCDN): Promise<boolean> {
		const yesButton = element(by.buttonText("Yes"));
		const queueButton = element(by.buttonText("Queue Updates"));

		switch (cdn.description) {
			case "perform snapshot":
				await element(by.name("diffCDNbtn")).click();
				const snapshotButton = by.partialButtonText("Perform Snapshot");
				if (!await browser.isElementPresent(snapshotButton)) {
					throw new Error("cannot find 'Perform Snapshot' button")
				}
				await element(snapshotButton).click();
				await yesButton.click();
				break;
			case "queue CDN updates":
				await queueButton.click();
				const queueSelection = by.linkText(`Queue ${cdn.Name}${randomize} Server Updates`);
				if (!await browser.isElementPresent(queueSelection)) {
					throw new Error("cannot find 'Queue CDN updates' button")
				}
				await element(queueSelection).click();
				await yesButton.click();
				break;
			case "clear CDN updates":
				await queueButton.click();
				const clearSelection = by.linkText(`Clear ${cdn.Name}${randomize} Server Updates`);
				if (!await browser.isElementPresent(clearSelection)) {
					throw new Error("Cannot find Clear CDN updates button")
				}
				await element(clearSelection).click();
				await yesButton.click();
				break;
			case "update cdn name":
				await this.txtCDNName.clear();
				await this.txtCDNName.sendKeys(cdn.NewName + randomize);
				await this.ClickUpdate();
				break;
			case "update cdn ttl override":
				await this.txtCDNTTLOverride.clear();
				await this.txtCDNTTLOverride.sendKeys(twoNumberRandomize);
				await this.ClickUpdate();
				break;
			default:
				throw new Error(`unhandleable description: '${cdn.description}'`);
		}

		return (await this.GetOutputMessage()) === cdn.validationMessage;
	}

	/**
	 * Deletes a CDN.
	 *
	 * @param cdn The Name of the CDN to be deleted.
	 * @param validationMessage A literal Alert text that will be checked for to
	 * determine success. If omitted, no Alert is expected.
	 * @returns Whether or not the deletion succeeded.
	 */
	public async deleteCDN(cdn: string, validationMessage?: string): Promise<boolean> {
		cdn += randomize;
		await element(by.buttonText("Delete")).click();
		await element(by.name("confirmWithNameInput")).sendKeys(cdn);
		await this.ClickDeletePermanently();
		return (await this.GetOutputMessage()) === validationMessage;
	}
}
