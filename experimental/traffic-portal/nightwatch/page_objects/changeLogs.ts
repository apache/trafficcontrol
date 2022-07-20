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

import { TABLE_COMMANDS, TableSectionCommands } from "../globals/tables";

/**
 * Defines Change Logs Commands
 */
type ChangeLogsTableSectionCommand = TableSectionCommands;

const changeLogsPageObject = {
	api: {} as NightwatchAPI,
	sections: {
		changeLogsTable: {
			commands: {
				...TABLE_COMMANDS
			} as ChangeLogsTableSectionCommand,
			elements: {
			},
			selector: "tp-change-logs mat-card"
		}
	},
	url(): string {
		return `${this.api.launchUrl}/core/change-logs`;
	}
};

/**
 * Defines the changeLogs table section.
 */
type ChangeLogsTableSection = EnhancedSectionInstance<ChangeLogsTableSectionCommand,
	typeof changeLogsPageObject.sections.changeLogsTable.elements>;

/**
 * The type of the changeLogs table page object as provided by the Nightwatch API at
 * runtime.
 */
export type ChangeLogsPageObject = EnhancedPageObject<{}, {}, { changeLogsTable: ChangeLogsTableSection }>;

export default changeLogsPageObject;
