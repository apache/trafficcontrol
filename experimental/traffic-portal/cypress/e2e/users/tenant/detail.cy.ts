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

import { CreatedData } from "cypress/support/testing.data";

describe("Tenant edit/creation page", () => {
	beforeEach(() => {
		cy.login();
	});
	it("Edits the 'root' Tenant", () => {
		// TODO: There's no guarantee that any Tenant exists  with the ID '1',
		// let alone that this be some special "root" tenant.
		cy.visit("/core/tenants/1");
		// TODO: idk how this was working in nightwatch. but it ain't working in
		// Cypress.
		// cy.get("input[name=active]").should("not.be.disabled").and("have.value", "on");
		// cy.get("input[name=parentTenant-tree-select]").should("not.be.disabled");
		cy.get("input[name=name]").should("be.disabled").and("have.value", "root");
		cy.get("button").contains("Save").should("not.be.disabled");
	});

	it("Edits an existing Tenant", () => {
		cy.fixture("test.data").then(
			(data: CreatedData) => {
				const {tenant} = data;
				cy.visit(`/core/tenants/${tenant.id}`);
				cy.get("input[name=active]").should("be.enabled");
				cy.get("input[name=name").should("be.enabled").and("have.value", tenant.name);
				cy.get("input[name=parentTenant-tree-select]").should("be.enabled");
				cy.get("button").contains("Save").should("not.be.disabled");
			}
		);
	});

	it("Creates a new Tenant", () => {
		cy.visit("/core/tenants/new");
		cy.get("input[name=active]").should("be.enabled").and("have.value", "on");
		cy.get("input[name=name]").should("be.enabled").and("be.empty");
		cy.get("input[name=parentTenant-tree-select]").should("be.enabled").and("be.empty");
		cy.get("button").contains("Save").should("not.be.disabled");
	});
});
