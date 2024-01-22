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

describe("Parameter edit/creation page", () => {
	beforeEach(() => {
		cy.login();
	});
	it("Test parameter", () => {
		cy.fixture("test.data").then(
			(data: CreatedData) => {
				const param = data.parameter;
				cy.visit(`/core/parameters/${param.id}`);
				cy.get("input[name=name]").should("be.enabled").and("have.value", param.name);
				cy.get("input[name=configFile]").should("be.enabled").and("have.value", param.configFile);
				cy.get("textarea[name=value]").should("be.enabled").and("have.value", param.value);
				cy.get("input[name=id]").should("not.be.enabled").and("have.value", param.id);
				cy.get("input[name=lastUpdated]").should("not.be.enabled");
				cy.get("button").contains("Save").should("not.be.disabled");
			}
		);
	});

	it("Creates new Parameters", () => {
		cy.visit("/core/parameters/new");
		cy.get("input[name=name]").should("be.enabled").and("have.value", "");
		cy.get("input[name=configFile]").should("be.enabled").and("have.value", "");
		cy.get("textarea[name=value]").should("be.enabled").and("have.value", "");
		cy.get("input[name=id]").should("not.exist");
		cy.get("input[name=lastUpdated]").should("not.exist");
		cy.get("button").contains("Save").should("not.be.disabled");
	});
});
