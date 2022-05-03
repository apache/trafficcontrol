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
	EnhancedElementInstance,
	EnhancedPageObject, EnhancedSectionInstance, NightwatchAPI
} from "nightwatch";

/**
 * Defines the commands for the loginForm section
 */
interface DeliveryServiceCardCommands extends EnhancedSectionInstance, EnhancedElementInstance<EnhancedPageObject> {
	expandDS(xmlId: string): Promise<boolean>;
	viewDetails(xmlId: string): Promise<boolean>;
}

/**
 * Defines the loginForm section
 */
type DeliveryServiceCardSection = EnhancedSectionInstance<DeliveryServiceCardCommands,
	typeof deliveryServiceCardPageObject.sections.cards.elements>;

/**
 * Define the type for our PO
 */
export type DeliveryServiceCardPageObject = EnhancedPageObject<{}, {}, { cards: DeliveryServiceCardSection }>;

const deliveryServiceCardPageObject = {
	api: {} as NightwatchAPI,
	sections: {
		cards: {
			commands: {
				async expandDS(xmlId: string): Promise<boolean> {
					return new Promise((resolve, reject) => {
						this.click("css selector", `mat-card#${xmlId}`, result => {
							if (result.status === 1) {
								reject(new Error(`Unable to find by css mat-card#${xmlId}`));
								return;
							}
							this.waitForElementVisible(`mat-card#${xmlId} mat-card-content > div`,
								undefined, undefined, undefined, () => {
									resolve(true);
								});
						});
					});
				},
				async viewDetails(xmlId: string): Promise<boolean> {
					await this.expandDS(xmlId);
					return new Promise((resolve) => {
						this.click("css selector", `mat-card#${xmlId} mat-card-actions > a`, () => {
							browser.assert.urlContains("deliveryservice");
							resolve(true);
						});
					});
				}
			} as DeliveryServiceCardCommands,
			elements: {
			},
			selector: "article#deliveryservices"
		}
	},
	url(): string {
		return `${this.api.launchUrl}/core`;
	}
};

export default deliveryServiceCardPageObject;
