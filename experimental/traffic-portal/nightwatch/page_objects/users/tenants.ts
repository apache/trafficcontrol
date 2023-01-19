/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { EnhancedPageObject, EnhancedSectionInstance, NightwatchAPI } from "nightwatch";

import { TABLE_COMMANDS, TableSectionCommands } from "../../globals/tables";

/**
 * Defines the Tenants table commands
 */
type TenantsTableCommands = TableSectionCommands;

/**
 * Defines the Page Object for the Tenants page.
 */
export type TenantsPageObject = EnhancedPageObject<{}, {},
EnhancedSectionInstance<TenantsTableCommands>>;

const tenantsPageObject = {
	api: {} as NightwatchAPI,
	sections: {
		tenantsTable: {
			commands: {
				...TABLE_COMMANDS
			},
			elements: {},
			selector: "mat-card"
		}
	},
	url(): string {
		return `${this.api.launchUrl}/core/tenants`;
	}
};

export default tenantsPageObject;
