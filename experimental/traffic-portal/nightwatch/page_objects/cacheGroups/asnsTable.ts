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
 * Defines the ASNs table commands
 */
type AsnsTableCommands = TableSectionCommands;

/**
 * Defines the Page Object for the ASNs page.
 */
export type AsnsPageObject = EnhancedPageObject<{}, {},
EnhancedSectionInstance<AsnsTableCommands>>;

const asnsPageObject = {
	api: {} as NightwatchAPI,
	sections: {
		asnsTable: {
			commands: {
				...TABLE_COMMANDS
			},
			elements: {},
			selector: "mat-card"
		}
	},
	url(): string {
		return `${this.api.launchUrl}/core/asns`;
	}
};

export default asnsPageObject;
