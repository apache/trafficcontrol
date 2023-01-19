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
import { TABLE_COMMANDS, TableSectionCommands } from "nightwatch/globals/tables";

/**
 * Defines the Divisions table commands
 */
type DivisionsTableCommands = TableSectionCommands;

/**
 * Defines the Page Object for the Divisions page.
 */
export type DivisionsPageObject = EnhancedPageObject<{}, {},
EnhancedSectionInstance<DivisionsTableCommands>>;

const divisionsPageObject = {
	api: {} as NightwatchAPI,
	sections: {
		divisionsTable: {
			commands: {
				...TABLE_COMMANDS
			},
			elements: {},
			selector: "mat-card"
		}
	},
	url(): string {
		return `${this.api.launchUrl}/core/divisions`;
	}
};

export default divisionsPageObject;
