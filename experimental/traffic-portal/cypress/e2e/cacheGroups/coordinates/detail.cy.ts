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

describe("Coordinate edit/creation page", () => {
	beforeEach(() => {
		cy.login();
	});
	it("Edits an existing Coordinate", () => {
		cy.fixture("test.data").then(
			(data: CreatedData) => {
				const coordinate = data.coordinate;
				cy.visit(`/core/coordinates/${coordinate.id}`);
				cy.get("mat-card").find("input[name=latitude]").should("be.enabled").should("have.value", coordinate.latitude);
				cy.get("mat-card").find("input[name=longitude]").should("be.enabled").should("have.value", coordinate.longitude);
				cy.get("mat-card").find("input[name=name]").should("be.enabled").should("have.value", coordinate.name);
				cy.get("mat-card").find("input[name=id]").should("not.be.enabled").should("have.value", coordinate.id);
				cy.get("mat-card").find("input[name=lastUpdated]").should("not.be.enabled");
				cy.get("mat-card").find("button").contains("Save").should("not.be.disabled");
			}
		);
	});

	it("Creates new Coordinate", () => {
		cy.visit("/core/coordinates/new");
		cy.get("mat-card").find("input[name=latitude]").should("be.enabled");
		cy.get("mat-card").find("input[name=longitude]").should("be.enabled");
		cy.get("mat-card").find("input[name=name]").should("be.enabled").should("have.value", "");
		cy.get("mat-card").find("input[name=id]").should("not.exist");
		cy.get("mat-card").find("input[name=lastUpdated]").should("not.exist");
		cy.get("mat-card").find("button").contains("Save").should("not.be.disabled");
	});
});
