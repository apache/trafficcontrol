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

import type { CreatedData } from "cypress/support/testing.data";
import type { ResponseDeliveryService } from "trafficops-types";

describe("Delivery Service info card", () => {
	let ds: ResponseDeliveryService;
	beforeEach(() => {
		cy.login();
		cy.fixture("test.data").then(
			(data: CreatedData) => {
				({ds} = data);
			}
		);
	});
	it("Expands and collapses", () => {
		cy.visit("/core");
		const selector = `#${ds.xmlId}`;
		cy.get(selector).find("mat-card-content > div").should("not.exist");
		cy.get(selector).click().find("mat-card-content > div").should("exist");
	});

	it("Has a button that takes you to a details view", async (): Promise<void> => {
		cy.visit("/core");
		const selector = `#${ds.xmlId}`;
		cy.get(selector).click();
		cy.get(selector).find("a").contains("View Details").click();
		cy.window().its("location.pathname").should("eq", `/core/deliveryservice/${ds.id}`);
	});
});
