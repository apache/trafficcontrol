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

import { TableSectionCommands, TABLE_COMMANDS } from "../globals/tables";

/**
 * Defines the commands for the servers table section.
 */
type ServersTableSectionCommands = TableSectionCommands;

const serversPageObject = {
	api: {} as NightwatchAPI,
	sections: {
		serversTable: {
			commands: {
				...TABLE_COMMANDS
			} as ServersTableSectionCommands,
			elements: {
			},
			selector: "servers-table main"
		}
	},
	url(): string {
		return `${this.api.launchUrl}/core/servers`;
	}
};

/**
 * Defines the servers table section.
 */
type ServersTableSection = EnhancedSectionInstance<ServersTableSectionCommands, typeof serversPageObject.sections.serversTable.elements>;

/**
 * The type of the servers table page object as provided by the Nightwatch API at
 * runtime.
 */
export type ServersPageObject = EnhancedPageObject<{}, {}, { serversTable: ServersTableSection }>;

export default serversPageObject;
