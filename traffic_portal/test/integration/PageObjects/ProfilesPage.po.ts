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

import { randomize } from "../config";
import { SideNavigationPage } from "./SideNavigationPage.po";

/** A definition of a Profile to create. */
interface CreateProfile {
	cdn: string;
	name: string;
	profileDescription: string;
	routingDisabled: string;
	type: string;
	validationMessage?: string;
}

/** An interface for interacting with various Profile pages. */
export class ProfilesPage extends SideNavigationPage {

	private static readonly txtType = element(by.name("type"));

	/**
	 * Creates the given Profile.
	 *
	 * @param profile The Profile to create.
	 * @returns Whether or not the creation was successful.
	 */
	public async createProfile(profile: CreateProfile): Promise<boolean> {
		await this.NavigateToProfilesPage();
		await element(by.buttonText("More")).click();
		await element(by.linkText("Create New Profile")).click();
		await Promise.all([
			element(by.name("name")).sendKeys(`${profile.name}${randomize}`),
			element(by.name("cdn")).sendKeys(profile.cdn),
			ProfilesPage.txtType.sendKeys(profile.type),
			element(by.name("routingDisabled")).sendKeys(profile.routingDisabled),
			element(by.id("description")).sendKeys(profile.profileDescription)
		]);
		await this.ClickCreate();
		const result = await this.GetOutputMessage();
		return result === profile.validationMessage;
	}

	/**
	 * Searches for a Profile in the table by Name, and clicks on the first
	 * matching entry if any exist.
	 *
	 * @param name The Name of the Profile for which to search.
	 * @returns Whether or not at least one matching entry is present after
	 * searching. If `false`, no navigation was performed (because it couldn"t
	 * be).
	 */
	public async searchProfile(name: string): Promise<boolean> {
		name += randomize;
		await this.NavigateToProfilesPage();

		const quickSearch = element(by.id("quickSearch"));
		await quickSearch.clear();
		await quickSearch.sendKeys(name);

		const match = element(by.cssContainingText("span", name));
		if (!await browser.isElementPresent(match)) {
			return false;
		}

		await match.click();
		return true;
	}

	/**
	 * Opens the "compare" page for the given two Profiles.
	 *
	 * @param profile1 The Name of the first Profile to compare.
	 * @param profile2 The Name of the second Profile to compare.
	 * @returns Whether or not all steps successfully completed and the browser
	 * is sitting
	 */
	public async compareProfiles(profile1: string, profile2: string): Promise<boolean> {
		await this.NavigateToProfilesPage();
		await element(by.buttonText("More")).click();
		await element(by.buttonText("Compare Profiles")).click();
		await Promise.all([
			element(by.name("compareDropdown1")).sendKeys(profile1),
			element(by.name("compareDropdown1")).sendKeys(profile2)
		]);
		await element(by.name("compareSubmit")).click();
		return element(by.id("profilesParamsCompareTable_wrapper")).isDisplayed();
	}

	/**
	 * Updates the Profile for which the browser is currently view details
	 * (navigation to the details page **must** be done *before* calling this
	 * method).
	 *
	 * @param type The new Type to give a Profile
	 * @param validationMessage A literal Alert text that will be checked for to
	 * determine success. If omitted, no Alert is expected.
	 * @returns Whether or not the update was successful.
	 */
	public async updateProfile(type: string, validationMessage: string): Promise<boolean> {
		await ProfilesPage.txtType.sendKeys(type);
		await this.ClickUpdate();
		return (await this.GetOutputMessage()) === validationMessage;
	}

	/**
	 * Deletes the Profile for which the browser is currently view details
	 * (navigation to the details page **must** be done *before* calling this
	 * method).
	 *
	 * @param name The Name of the Profile to be deleted.
	 * @param validationMessage A literal Alert text that will be checked for to
	 * determine success. If omitted, no Alert is expected.
	 * @returns Whether or not the update was successful.
	 */
	public async deleteProfile(name: string, validationMessage: string): Promise<boolean> {
		await element(by.buttonText("Delete")).click();
		await element(by.name("confirmWithNameInput")).sendKeys(`${name}${randomize}`);
		await this.ClickDeletePermanently();
		return await this.GetOutputMessage() === validationMessage;
	}

	/**
	 * Toggles the visibility of a column on the Profiles table page.
	 *
	 * This does not perform navigation to that page!
	 *
	 * @param name The name of the column to toggle.
	 * @returns Whether or not the column is visible after being toggled.
	 */
	public async toggleTableColumn(name: string): Promise<boolean> {
		const btnTableColumn = element(by.className("fa-columns"));
		await btnTableColumn.click();
		try {
			await element(by.cssContainingText("label", name)).click();
		} finally {
			await btnTableColumn.click();
		}
		return element(by.cssContainingText('span[role="columnheader"]', name)).isPresent();
	}
}
