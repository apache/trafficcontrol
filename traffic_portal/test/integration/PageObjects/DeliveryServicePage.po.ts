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

import { BasePage } from "./BasePage.po";
import { randomize } from "../config";
import { SideNavigationPage } from "./SideNavigationPage.po";
import { browser, by, element, ExpectedConditions } from "protractor";

/**
 * The DeliveryServicePage is a page object modelling of the Delivery Service
 * editing/creation view. For simplicity"s sake, it also provides functionality
 * that relates to the Delivery Services table view.
 */
export class DeliveryServicePage extends BasePage {

	/** The search box in the DS table view. */
	private readonly txtSearch = element(by.id("quickSearch"));

	/** The "Display Name" text input in the editing/creation view(s). */
	private readonly txtDisplayName = element(by.name("displayName"));
	/** The "More" dropdown menu button in the editing/creation view(s). */
	private readonly  btnMore = element(by.name("moreBtn"));

	/**
	 * Navigates to the Delivery Services table view.
	 */
	public async OpenDeliveryServicePage(): Promise<void> {
		const snp = new SideNavigationPage();
		return snp.NavigateToDeliveryServicesPage();
	}

	/**
	 * Toggles the open/close state of the "Services" sub-menu in the left-side
	 * navigation pane.
	 */
	public async OpenServicesMenu(): Promise<void> {
		const snp = new SideNavigationPage();
		return snp.ClickServicesMenu();
	}

	/**
	 * Creates a new Delivery Service.
	 *
	 * @param deliveryservice Details for the Delivery Service to be created.
	 * @returns The text shown in the first Alert pane found after creation.
	 */
	public async CreateDeliveryService(name: string, type: string, tenant: string): Promise<string> {
		await this.btnMore.click();
		await element(by.buttonText("Create Delivery Service")).click();
		await element(by.name("selectFormDropdown")).sendKeys(type);
		await element(by.buttonText("Submit")).click();

		name += randomize;
		tenant += randomize;

		const ps = [];
		switch (type) {
			case "ANY_MAP":
				ps.push(element(by.name("remapText")).sendKeys("test"));
			break;

			case "DNS":
				ps.push(element(by.name("capability-0")).click());
			case "HTTP":
				ps.push(
					element(by.name("orgServerFqdn")).sendKeys("http://origin.infra.ciab.test"),
					element(by.name("capability-0")).click()
				);
			case "STEERING":
				ps.push(element(by.name("protocol")).sendKeys("HTTP"));
			break;

			default:
				throw new Error(`invalid Delivery Service routing type: ${type}`);
		}
		ps.push(
			element(by.name("xmlId")).sendKeys(name),
			this.txtDisplayName.sendKeys(name),
			element(by.name("active")).sendKeys("Active"),
			element(by.id("type")).sendKeys(type),
			element(by.name("tenantId")).click().then(() => element(by.name(tenant)).click()),
			element(by.name("cdn")).sendKeys("dummycdn"+randomize)
		);

		await Promise.all(ps);
		await element(by.buttonText("Create")).click();

		return this.GetOutputMessage();
	}

	/**
	 * Searches the table for a Delivery Service in the table.
	 *
	 * (Note this neither checks nor enforces that the sought-after DS is
	 * actually found.)
	 *
	 * @param name The name for which to search.
	 */
	public async SearchDeliveryService(name: string): Promise<void> {
		name += randomize;

		await this.txtSearch.clear();
		await this.txtSearch.sendKeys(name);
		const nameSpan = element(by.cssContainingText("span", name));
		await nameSpan.click();
	}

	/**
	 * Changes a Delivery Service's Display Name to the provided value (after
	 * randomization).
	 *
	 * @param newName The new Display Name to be given to the Delivery Service.
	 * @returns The text shown in the first Alert pane found after attempting to
	 * submit the update.
	 */
	public async UpdateDeliveryServiceDisplayName(newName: string): Promise<string> {
		await this.txtDisplayName.clear();
		await this.txtDisplayName.sendKeys(newName + randomize);
		await this.ClickUpdate();
		return this.GetOutputMessage();
	}

	/**
	 * Attempts to delete a Delivery Service.
	 *
	 * @param name The XMLID of the Delivery Service to be deleted.
	 * @returns The text shown in the first Alert pane found after attempting
	 * the deletion.
	 */
	public async DeleteDeliveryService(name: string): Promise<string> {
		name += randomize;
		await element(by.buttonText("Delete")).click();
		await element(by.name("confirmWithNameInput")).sendKeys(name);
		await this.ClickDeletePermanently();
		return this.GetOutputMessage();
	}

	/**
	 * Assigns the server with the given hostname to the Delivery Service. Note
	 * that the browser must already be on a Delivery Service edit view for this
	 * to work, as this method neither navigates to it nor back to the table
	 * view afterward!
	 *
	 * @param serverName The name of the server being assigned.
	 * @returns The text shown in the first Alert pane found after attempting
	 * the assignment.
	 */
	public async AssignServerToDeliveryService(serverName: string): Promise<string>{
		await this.btnMore.click();
		await element(by.linkText("Manage Servers")).click();
		await this.btnMore.click();
		await element(by.partialButtonText("Assign")).click();
		const serverCell = element(by.cssContainingText(".ag-cell-value", serverName));
		await browser.wait(ExpectedConditions.elementToBeClickable(serverCell), 3000);
		await serverCell.click();
		await this.ClickSubmit();
		return this.GetOutputMessage();
	}
}
