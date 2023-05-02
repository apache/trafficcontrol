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

import { EnhancedElementInstance, EnhancedPageObject, EnhancedSectionInstance } from "nightwatch";

/**
 * Defines the commands for the sidebar
 */
interface SidebarCommands  extends EnhancedSectionInstance, EnhancedElementInstance<EnhancedPageObject> {
	navigateToNode(node: string, path: Array<string>): Promise<void>;
}

/**
 * Defines the sidebar section
 */
type SidebarSection = EnhancedSectionInstance<SidebarCommands, typeof commonPageObject.sections.sidebar.elements>;

/**
 * Defines the type for the common PO
 */
export type CommonPageObject = EnhancedPageObject<{}, typeof commonPageObject.elements, { sidebar: SidebarSection }>;

const commonPageObject = {
	elements: {
		snackbarEle: {
			selector: "simple-snack-bar"
		}
	},
	sections: {
		sidebar: {
			commands: {
				async navigateToNode(node: string, path: Array<string>): Promise<void> {
					for (const pathNode of path) {
						await this.click(`@${pathNode}`);
					}
					await this.click(`@${node}`);
				}
			} as SidebarCommands,
			elements: {
				asns: "[aria-label='Navigate to ASNs']",
				cacheGroups: "[aria-label='Navigate to Cache Groups']",
				cacheGroupsContainer: "[aria-label='Toggle Cache Groups']",
				changeLogs: "[aria-label='Navigate to Change Logs']",
				configurationContainer: "[aria-label='Toggle Configuration']",
				coordinates: "[aria-label='Navigate to Coordinates']",
				dashboard: "[aria-label='Navigate to Dashboard']",
				divisions: "[aria-label='Navigate to Divisions']",
				otherContainer: "[aria-label='Toggle Other']",
				physicalLocations: "[aria-label='Navigate to Physical Locations']",
				profile: "[aria-label='Navigate to My Profile']",
				profiles: "[aria-label='Navigate to Profiles']",
				regions: "[aria-label='Navigate to Regions']",
				servers: "[aria-label='Navigate to Servers']",
				serversContainer: "[aria-label='Toggle Servers']",
				statuses: "[aria-label='Navigate to Statuses']",
				tenants: "[aria-label='Navigate to Tenants']",
				types: "[aria-label='Navigate to Types']",
				users: "[aria-label='Navigate to Users']",
				usersContainer: "[aria-label='Toggle Users']"
			},
			selector: "#sidebar-nav-tree"
		}
	}
};

export default commonPageObject;
