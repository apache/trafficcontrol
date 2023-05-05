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
import { ResponseServer } from "trafficops-types";

import { TableSectionCommands, TABLE_COMMANDS } from "../../globals/tables";

/**
 * Defines the commands for the servers table section.
 */
interface ServersTableSectionCommands extends TableSectionCommands {
	createNew(): Promise<void>;
	openDetails(s: ResponseServer): Promise<void>;
	open(): Promise<void>;
}

const serversPageObject = {
	api: {} as NightwatchAPI,
	sections: {
		serversTable: {
			commands: {
				async createNew(): Promise<void> {
					await this.open();
					await browser.click("a.page-fab[routerLink='new']");
				},
				async open(): Promise<void> {
					await browser.page.common()
						.section.sidebar
						.navigateToNode("servers", ["serversContainer"]);
				},
				async openDetails(server: ResponseServer): Promise<void> {
					await this.open();
					const table = browser.page.servers.serversTable().section.serversTable;
					await table
						.filterTableByColumn("Host", server.hostName);
					await table.doubleClickRow(1);
				},
				...TABLE_COMMANDS
			} as ServersTableSectionCommands,
			elements: {
			},
			selector: "mat-card"
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
export type ServersTablePageObject = EnhancedPageObject<{}, {}, { serversTable: ServersTableSection }>;

export default serversPageObject;
