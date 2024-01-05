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

describe("Profile creation/edit page", () => {
	beforeEach(() => {
		cy.login();
	});
	it("Edits an existing Profile", () => {
		cy.fixture("test.data").then(
			(data: CreatedData) => {
				const {profile} = data;
				cy.visit(`/core/profiles/${profile.id}`);
				cy.get("input[name=name]").should("be.enabled").and("have.value", profile.name);
				cy.get("mat-select[name=cdn]").should("not.be.disabled");
				cy.get("mat-select[name=type]").should("not.be.disabled");
				cy.get("mat-select[name=routingDisabled]").should("not.be.disabled");
				cy.get("textarea[name=description]").should("be.enabled").and("have.value", profile.description);
				cy.get("input[name=id]").should("not.be.enabled").and("have.value", profile.id);
				cy.get("input[name=lastUpdated]").should("not.be.enabled");
				cy.get("button").contains("Save").should("not.be.disabled");
			}
		);
	});

	it("Creates new Profiles", () => {
		cy.visit("/core/profiles/new");
		cy.get("input[name=name]").should("be.enabled").and("be.empty");
		cy.get("mat-select[name=cdn]").should("not.be.disabled");
		cy.get("mat-select[name=type]").should("not.be.disabled");
		cy.get("mat-select[name=routingDisabled]").should("not.be.disabled");
		cy.get("textarea[name=description]").should("be.enabled").and("be.empty");
		cy.get("input[name=id]").should("not.exist");
		cy.get("input[name=lastUpdated]").should("not.exist");
		cy.get("button").contains("Save").should("not.be.disabled");
	});
});
