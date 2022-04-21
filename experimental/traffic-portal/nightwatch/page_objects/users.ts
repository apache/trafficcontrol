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
import {
	EnhancedElementInstance,
	EnhancedPageObject,
	EnhancedSectionInstance,
	NightwatchAPI
} from "nightwatch";

/**
 * Defines the commands for the users table section.
 */
interface UsersTableSectionCommands extends EnhancedSectionInstance, EnhancedElementInstance<EnhancedPageObject> {
	getColumnState(column: string): Promise<boolean>;
	searchText(text: string): this;
	toggleColumn(column: string): this;
}

const usersPageObject = {
	api: {} as NightwatchAPI,
	sections: {
		usersTable: {
			commands: {
				async getColumnState(column: string): Promise<boolean> {
					return new Promise((resolve, reject) => {
						this.click("@columnMenu").getElementProperty(`input[name='column-${column}']`, "checked",
							result => {
								if (typeof result.value !== "boolean") {
									console.error("incorrect type for 'checked' DOM property:", result.value);
									reject(new Error(`incorrect type for 'checked' DOM property: ${typeof result.value}`));
									return;
								}
								resolve(result.value);
							}
						).click("@columnMenu");
					});
				},
				searchText(text: string): UsersTableSectionCommands  {
					 return this.setValue("@searchbox", text);
				},
				toggleColumn(column: string): UsersTableSectionCommands {
					return this.click("@columnMenu").click(`input[name='${column}']`).click("@columnMenu");
				},
			} as UsersTableSectionCommands,
			elements: {
				columnMenu: {
					selector: "button.dropdown-toggle"
				},
				searchbox: {
					selector: "input[name='fuzzControl']"
				},
			},
			selector: "main > main"
		}
	},
	url(): string {
		return `${this.api.launchUrl}/core/users`;
	}
};

/**
 * Defines the users table section.
 */
type UsersTableSection = EnhancedSectionInstance<UsersTableSectionCommands, typeof usersPageObject.sections.usersTable.elements>;

/**
 * The type of the users table page object as provided by the Nightwatch API at
 * runtime.
 */
export type UsersPageObject = EnhancedPageObject<{}, {}, { usersTable: UsersTableSection }>;

export default usersPageObject;
