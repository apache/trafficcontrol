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

import { EnhancedPageObject, EnhancedSectionInstance } from "nightwatch";

/**
 * Defines the Section Instance for the Servers Details.
 */
type ServersDetailSection = EnhancedSectionInstance<
EnhancedSectionInstance<typeof serversDetailPageObject.sections.detailCard.commands,
		typeof serversDetailPageObject.sections.detailCard.elements>>;

/**
 * Defines the PageObject for Servers Details.
 */
export type ServersDetailPageObject = EnhancedPageObject<{}, {}, { detailCard: ServersDetailSection}>;

const serversDetailPageObject = {
	sections: {
		detailCard: {
			commands: {
			},
			elements: {
				cacheGroup: "mat-select[name='cachegroup']",
				cdn: "mat-select[name='cdn']",
				deleteBtn: "button[aria-label='Delete Server']",
				domainName: "input[name='domainname']",
				hostName: "input[name='hostname']",
				httpPort: "input[name='httpport']",
				httpsPort: "input[name='httpsport']",
				id: "input[name='serverId']",
				iloGateway: "input[name='iloGateway']",
				iloIP: "input[name='iloIP']",
				iloNetmask: "input[name='iloNetmask']",
				iloPassword: "input[name='iloPassword']",
				iloUsername: "input[name='iloUsername']",
				intfAddBtn: "button[aria-label='Add An Interface']",
				lastUpdated: "input[name='lastUpdated']",
				mgmtGateway: "input[name='mgmtIpGateway']",
				mgmtIP: "input[name='mgmtIP']",
				mgmtNetmask: "input[name='mgmtIpNetmask']",
				offlineReason: "mat-select[name='offlineReason']",
				physLoc: "mat-select[name='physLocation']",
				profileNames: "mat-select[name='profiles']",
				rack: "input[name='rack']",
				status: "mat-select[name='status']",
				statusDisabled: "input[name='status']",
				statusLastUpdated: "input[name='statusLastUpdated']",
				submitBtn: "button[aria-label='Submit Server']",
				type: "mat-select[name='type']"
			},
			selector: "mat-card.page-content"
		}
	}
};

export default serversDetailPageObject;
