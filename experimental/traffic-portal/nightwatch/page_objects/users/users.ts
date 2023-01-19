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
	EnhancedPageObject,
	EnhancedSectionInstance,
	NightwatchAPI
} from "nightwatch";
import { TableSectionCommands, TABLE_COMMANDS } from "nightwatch/globals/tables";

/**
 * Defines the commands for the users table section.
 */
type UsersTableSectionCommands = TableSectionCommands;

const usersPageObject = {
	api: {} as NightwatchAPI,
	sections: {
		usersTable: {
			commands: {
				...TABLE_COMMANDS
			} as UsersTableSectionCommands,
			elements: {
			},
			selector: "mat-card"
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
