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

describe("Status edit/creation page", () => {
	beforeEach(() => {
		cy.login();
	});
	it("Edits an existing Status", () => {
		cy.fixture("test.data").then(
			(data: CreatedData) => {
				const {status} = data;
				cy.visit(`/core/statuses/${status.id}`);
				cy.get("input[name=name]").should("be.enabled").and("have.value", status.name);
				cy.get("input[name=description]").should("be.enabled").and("have.value", status.description);
				cy.get("input[name=id]").should("be.disabled").and("have.value", status.id);
				cy.get("input[name=lastUpdated]").should("be.disabled");
				cy.get("button").contains("Save").should("not.be.disabled");
				cy.get("button").contains("Create").should("not.exist");
			}
		);
	});

	it("Creates a new Status", () => {
		cy.visit("/core/statuses/new");
		cy.get("input[name=name]").should("be.enabled").and("be.empty");
		cy.get("input[name=description]").should("be.enabled").and("be.empty");
		cy.get("input[name=id]").should("not.exist");
		cy.get("input[name=lastUpdated]").should("not.exist");
		cy.get("button").contains("Save").should("not.exist");
		cy.get("button").contains("Create").should("not.be.disabled");
	});
});
