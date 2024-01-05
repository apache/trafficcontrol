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

describe("Server Capability edit/creation page", () => {
	beforeEach(() => {
		cy.login();
	});
	it("Edits an existing Capability", () => {
		cy.fixture("test.data").then(
			(data: CreatedData) => {
				const {capability} = data;
				cy.visit(`/core/capabilities/${capability.name}`);
				cy.get("input[name=name]").should("be.enabled").and("have.value", capability.name);
				cy.get("input[name=lastUpdated]").should("be.disabled");
				cy.get("button").contains("Save").should("not.be.disabled");
			}
		);
	});

	it("Creates new Capabilities", () => {
		cy.visit("/core/new-capability");
		cy.get("input[name=name]").should("be.enabled").and("be.empty");
		cy.get("button").contains("Save").should("not.be.disabled");
		cy.get("input[name=lastUpdated]").should("not.exist");
	});
});
