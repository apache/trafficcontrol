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

describe("CDN edit/creation page", () => {
	beforeEach(() => {
		cy.login();
	});
	it("Edits an existing CDN", () => {
		cy.fixture("test.data").then(
			(data: CreatedData) => {
				const {cdn} = data;
				cy.visit(`/core/cdns/${cdn.id}`);
				cy.get("mat-card").find("input[name=name]").should("be.enabled").should("have.value", cdn.name);
				cy.get("mat-card").find("input[name=id]").should("be.disabled").should("have.value", cdn.id);
				cy.get("mat-card").find("input[name=lastUpdated]").should("be.disabled");
				cy.get("mat-card").find("button").contains("Save").should("not.be.disabled");
			}
		);
	});

	it("Creates new CDNs", () => {
		cy.visit("/core/cdns/new");
		cy.get("mat-card").find("input[name=name]").should("be.enabled").should("have.value", "");
		cy.get("mat-card").find("input[name=id]").should("not.exist");
		cy.get("mat-card").find("input[name=lastUpdated]").should("not.exist");
		cy.get("mat-card").find("button").contains("Save").should("not.be.disabled");
	});
});
