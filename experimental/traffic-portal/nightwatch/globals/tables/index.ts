/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

import type {
	Awaitable,
	EnhancedElementInstance,
	EnhancedPageObject,
	EnhancedSectionInstance,
	WebDriverProtocolUserActions
} from "nightwatch";

/**
 * TableSectionCommands is the base type for page object sections representing
 * pages containing AG-Grid generic tables.
 */
export interface TableSectionCommands extends EnhancedSectionInstance,
	EnhancedElementInstance<EnhancedPageObject>, WebDriverProtocolUserActions {
	doubleClickRow<T extends this>(row: number): Promise<T>;

	filterTableByColumn<T extends this>(column: string, search: string): Promise<T>;

	gotoRowByColumn<T extends this>(column: string, search: string): Promise<T>;

	getColumnState(column: string): Promise<boolean>;

	searchText<T extends this>(text: string): T;

	toggleColumn(column: string): Promise<this>;
}

/**
 * A CSS selector for an AG-Grid generic table's column visibility dropdown
 * menu.
 */
export const columnMenuBtnSelector = "div.toggle-columns > button.mat-mdc-menu-trigger";

/**
 * A CSS selector for an AG-Grid generic table's column visibility dropdown
 * menu.
 */
export const columnMenuCloseSelector = ".cdk-overlay-backdrop";

/**
 * A CSS selector for an AG-Grid generic table's "Fuzzy Search" input text box.
 */
export const searchboxSelector = "input[name='fuzzControl']";

/**
 * CSS selector for the AG-Grid row(s).
 */
export const tableRowsSelector = ".ag-center-cols-clipper .ag-row";

/**
 * Gets the state of an AG-Grid column by checking whether it's checked
 * in the column visibility menu (doesn't actually verify that this means the
 * column is visible).
 *
 * @param this Special parameter that tells the compiler what `this` is in a
 * valid context for this function.
 * @param column The name of the column being retrieved.
 * @returns The state of the column named `column`. Behavior is undefined if
 * multiple columns exist with the same given name.
 */
export async function getColumnState(this: TableSectionCommands, column: string): Promise<boolean> {
	const selector = `input[type='checkbox'][name='${column}']`;
	await this.click(columnMenuBtnSelector);
	const selected = await browser.isSelected(selector);
	return Promise.resolve(selected);
}

/**
 * Filters a table by a column
 *
 * @param this Special parameter that tells the compiler what `this` is in a
 * valid context for this function.
 * @param column Which column to filter
 * @param text Text to filter by
 */
export async function filterTableByColumn<T extends TableSectionCommands>(
	this: TableSectionCommands,
	column: string,
	text: string): Promise<T> {
	if (!await this.getColumnState(column)) {
		await this.toggleColumn(column);
	}
	this.searchText(text);
	return Promise.resolve(this) as Promise<T>;
}

/**
 * Double-clicks the nth row on a table
 *
 * @param this Special parameter that tells the compiler what `this` is in a
 * valid context for this function.
 * @param rowNumber Which row to click
 * @returns The calling command section for call-chaining the way Nightwatch
 * likes to do.
 */
export function doubleClickRow<T extends TableSectionCommands>(this: TableSectionCommands, rowNumber: number): Awaitable<T, null> {
	return this.doubleClick("css selector", `${tableRowsSelector}:nth-of-type(${rowNumber})`) as Awaitable<T, null>;
}

/**
 * Filters a table by a column, then double-clicks the first row resulting from the filtering.
 *
 * @param this Special parameter that tells the compiler what `this` is in a
 * valid context for this function.
 * @param column Which column to filter
 * @param text Text to filter by
 */
export async function gotoRowByColumn<T extends TableSectionCommands>(
	this: TableSectionCommands, column: string, text: string): Promise<T> {
	await this.filterTableByColumn(column, text);
	return this.doubleClickRow(1);
}

/**
 * Sets the text of the table's "Fuzzy Search" searchbox.
 *
 * @param this Special parameter that tells the compiler what `this` is in a
 * valid context for this function.
 * @param text The text to set in the search input.
 * @returns The calling command section for call-chaining the way Nightwatch
 * likes to do.
 */
export function searchText<T extends TableSectionCommands>(this: T, text: string): Awaitable<T, null> {
	return this.setValue(searchboxSelector, text);
}

/**
 * Toggles the presence of a given column.
 *
 * @param this Special parameter that tells the compiler what `this` is in a
 * valid context for this function.
 * @param column The name of the column to be toggled.
 * @returns The calling command section for call-chaining the way Nightwatch
 * likes to do.
 */
export async function toggleColumn<T extends TableSectionCommands>(this: T, column: string): Promise<T> {
	const selector = `input[type='checkbox'][name='${column}']`;
	await browser.findElement(".mat-mdc-menu-panel")
		.getLocationInView(selector)
		.click(selector);
	await browser.click(columnMenuCloseSelector);
	return Promise.resolve(this);
}

/**
 * This is meant to be mixed-in to generic table page object command sections,
 * to most easily provide all the functionality of a table.
 */
export const TABLE_COMMANDS = {
	doubleClickRow,
	filterTableByColumn,
	getColumnState,
	gotoRowByColumn,
	searchText,
	toggleColumn
};
