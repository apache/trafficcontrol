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

describe("DS Content Invalidation Jobs page", () => {
	let id: number;
	before(() => {
		cy.fixture("test.data").then(
			(data: CreatedData) => {
				({id} = data.ds);
			}
		);
	});
	beforeEach(() => {
		cy.login();
		cy.visit(`/core/deliveryservice/${id}/invalidation-jobs`);
	});

	it("Manage Job", () => {
		cy.get("tp-header").find("h1").should("contain.text", "Content Invalidation Jobs").and("not.contain.text", "Loading");
		cy.get("button#new").focus().trigger("click");

		const startDate = new Date();
		startDate.setDate(startDate.getDate() + 1);
		cy.get("input[name='startDate']").should("be.visible").and("have.value", startDate.toLocaleDateString());
		cy.get("input[name='regexp']").click().focus().type("invalidateMe");
		cy.get("button[type=submit]").click();
		cy.get("simple-snack-bar").contains("created");
		cy.get("simple-snack-bar").find("button").first().click();

		cy.get("div.invalidation-job").find("button").should("be.visible").should("have.length", 2).should("not.be.disabled");
		cy.get("div.invalidation-job button").first().click();
		cy.get("tp-new-invalidation-job-dialog")
			.find("input[name='startDate']")
			.should("be.visible")
			.and("have.value", startDate.toLocaleDateString());
		cy.get("input[name='regexp']").should("have.value", "invalidateMe").click().focus().clear().type("/invalidateMe2");
		cy.get("button[type=submit]").click();

		cy.get("simple-snack-bar").contains("created");
		cy.get("simple-snack-bar").find("button").first().click();

		cy.get("div.invalidation-job button").last().click();
		cy.get("simple-snack-bar").contains("was deleted");
	});
});
