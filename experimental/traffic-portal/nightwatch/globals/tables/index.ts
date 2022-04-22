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

import type { EnhancedElementInstance, EnhancedPageObject, EnhancedSectionInstance } from "nightwatch";

/**
 * TableSectionCommands is the base type for page object sections representing
 * pages containing AG-Grid generic tables.
 */
export interface TableSectionCommands extends EnhancedSectionInstance, EnhancedElementInstance<EnhancedPageObject> {
	getColumnState(column: string): Promise<boolean>;
	searchText<T extends this>(text: string): T;
	toggleColumn<T extends this>(column: string): T;
}

/**
 * A CSS selector for an AG-Grid generic table's column visibility dropdown
 * menu.
 */
export const columnMenuSelector = "button.dropdown-toggle";

/**
 * A CSS selector for an AG-Grid generic table's "Fuzzy Search" input text box.
 */
export const searchboxSelector = "input[name='fuzzControl']";

/**
 * Gets the state of an AG-Grid column by checking whether or not it's checked
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
	return new Promise((resolve, reject) => {
		this.click(columnMenuSelector).getElementProperty(`input[name='column-${column}']`, "checked",
			result => {
				if (typeof result.value !== "boolean") {
					console.error("incorrect type for 'checked' DOM property:", result.value);
					reject(new Error(`incorrect type for 'checked' DOM property: ${typeof result.value}`));
					return;
				}
				resolve(result.value);
			}
		).click(columnMenuSelector);
	});
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
export function searchText<T extends TableSectionCommands>(this: T, text: string): T  {
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
export function toggleColumn<T extends TableSectionCommands>(this: T, column: string): T {
	return this.click(columnMenuSelector).click(`input[name='${column}']`).click(columnMenuSelector);
}

/**
 * This is meant to be mixed-in to generic table page object command sections,
 * to most easily provide all the functionality of a table.
 */
export const TABLE_COMMANDS = {
	getColumnState,
	searchText,
	toggleColumn
};
