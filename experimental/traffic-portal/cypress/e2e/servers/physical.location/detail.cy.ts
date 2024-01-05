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

describe("Physical Location creation/edit page", () => {
	beforeEach(() => {
		cy.login();
	});
	it("Edits an existing Physical Location", () => {
		cy.fixture("test.data").then(
			(data: CreatedData) => {
				const {physLoc} = data;
				cy.visit(`/core/phys-locs/${physLoc.id}`);
				cy.get("input[name=name]").should("be.enabled").and("have.value", physLoc.name);
				cy.get("input[name=id]").should("be.disabled").and("have.value", physLoc.id);
				cy.get("input[name=lastUpdated]").should("be.disabled");
				cy.get("input[name=address]").should("be.enabled").and("have.value", physLoc.address);
				cy.get("input[name=city]").should("be.enabled").and("have.value", physLoc.city);
				cy.get("textarea[name=comments]").should("be.enabled").and("have.value", physLoc.comments ?? "");
				cy.get("input[name=email]").should("be.enabled").and("have.value", physLoc.email ?? "");
				cy.get("input[name=poc]").should("be.enabled").and("have.value", physLoc.poc ?? "");
				cy.get("input[name=shortName]").should("be.enabled").and("have.value", physLoc.shortName);
				cy.get("input[name=state]").should("be.enabled").and("have.value", physLoc.state);
				cy.get("input[name=postalCode]").should("be.enabled").and("have.value", physLoc.zip);
				cy.get("input[name=phone]").should("be.enabled").and("have.value", physLoc.phone ?? "");
				cy.get("mat-select[name=region]").should("contain.text", physLoc.region ?? "");
				cy.get("button").contains("Save").should("not.be.disabled");
			}
		);
	});

	it("Creates new Physical Locations", () => {
		cy.visit("/core/phys-locs/new");
		cy.get("input[name=name]").should("be.enabled").and("be.empty");
		cy.get("input[name=id]").should("not.exist");
		cy.get("input[name=lastUpdated]").should("not.exist");
		cy.get("button").contains("Save").should("not.be.disabled");
	});
});
