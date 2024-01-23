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
 * Defines the Parameters table commands
 */
type ParametersTableCommands = TableSectionCommands;

/**
 * Defines the Page Object for the Parameters page.
 */
export type ParametersPageObject = EnhancedPageObject<{}, {}, EnhancedSectionInstance<ParametersTableCommands>>;

const parametersPageObject = {
	api: {} as NightwatchAPI,
	sections: {
		parametersTable: {
			commands: {
				...TABLE_COMMANDS
			},
			elements: {},
			selector: "mat-card"
		}
	},
	url(): string {
		return `${this.api.launchUrl}/core/parameters`;
	}
};

export default parametersPageObject;
