/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

describe("Servers table page", () => {
	beforeEach(() => {
		cy.login();
		cy.visit("/core/servers");
	});
	it("Filters servers by hostname", () => {
		cy.get("input[name=fuzzControl]").focus().type("edge");
		cy.window().its("location.search").should("contain", "search=edge");
	});
	it("Queues and clears revalidations on a server", () => {
		cy.get("input[name=fuzzControl]").focus().type("edge");

		// We need to force re-rendering of the table every time we do
		// something, or cypress moves too fast and undoes things it's doing
		// before the effects can be seen. This could be fixed by splitting
		// these into separate tests, but that wouldn't be faster and would have
		// the added drawback that it depends on the initial state of the data
		// and the order in which the tests are run.
		const reload = (): void => {
			cy.reload();
			cy.get("button[aria-label='column visibility menu']").click();
			cy.get("input[type=checkbox][name='Reval Pending']").check();
			cy.get("body").click(); // closes the menu so you can interact with other things.
		};

		reload();

		cy.get(".ag-row:visible").first().rightclick();
		cy.get("button").contains("Queue Content Revalidation").click();
		reload();

		cy.get(".ag-cell[col-id=revalPending]").first().should("contain.text", "schedule");
		cy.get(".ag-row:visible").first().rightclick();
		cy.get("button").contains("Clear Queued Content Revalidations").click();
		reload();

		cy.get(".ag-cell[col-id=revalPending]").first().should("contain.text", "done");
		cy.get(".ag-row:visible").first().rightclick();
		cy.get("button").contains("Queue Content Revalidation").click();
		reload();

		cy.get(".ag-cell[col-id=revalPending]").first().should("contain.text", "schedule");
	});
});
