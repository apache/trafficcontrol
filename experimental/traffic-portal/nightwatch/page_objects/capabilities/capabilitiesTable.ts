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

import { TABLE_COMMANDS, TableSectionCommands } from "../../../globals/tables";

/**
 * Defines the Capabilities table commands.
 */
type CapabilitiesTableCommands = TableSectionCommands;

/**
 * Defines the Page Object for the Capabilities page.
 */
export type CapabilitiesPageObject = EnhancedPageObject<{}, {}, EnhancedSectionInstance<CapabilitiesTableCommands>>;

const capabilitiesPageObject = {
	api: {} as NightwatchAPI,
	sections: {
		capabilitiesTable: {
			commands: {
				...TABLE_COMMANDS
			},
			elements: {},
			selector: "mat-card"
		}
	},
	url(): string {
		return `${this.api.launchUrl}/core/capabilities`;
	}
};

export default capabilitiesPageObject;
