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

import type { CreatedData } from "cypress/support/testing.data";

describe("Delivery Service details page", () => {
	let testData: CreatedData;
	before(() => {
		cy.fixture("test.data").then(
			(data: CreatedData) => {
				testData = data;
			}
		);
	});
	beforeEach(() => {
		cy.login();
		cy.visit(`/core/deliveryservice/${testData.ds.id}`);
	});

	it("Contains all the correct charts and controls", (): void => {
		cy.get("#bandwidthData").should("be.visible");
		cy.get("#tpsChartData").should("be.visible");
		cy.get("#invalidate").should("be.visible");
		cy.get("input[name=fromdate]").should("be.enabled");
		cy.get("input[name=fromtime]").should("be.enabled");
		cy.get("input[name=todate]").should("be.enabled");
		cy.get("input[name=totime]").should("be.enabled");
		cy.get("button[name=timespanRefresh]").should("not.be.disabled");
	});

	it("Sets default properties appropriately", (): void => {
		const now = new Date();
		const date = `${now.getFullYear()}-${(now.getMonth()+1).toString().padStart(2, "0")}-${now.getDate().toString().padStart(2, "0")}`;
		const time = `${now.getHours().toString().padStart(2, "0")}:${now.getMinutes().toString().padStart(2, "0")}`;

		cy.get("input[name=fromdate]").should("have.value", date);
		cy.get("input[name=fromtime]").should("have.value", "00:00");
		cy.get("input[name=todate]").should("have.value", date);
		cy.get("input[name=totime]").should("have.value", time);
	});

	it("Displays the steering action based on the DS type", (): void => {
		cy.get("div.actions > mat-icon").should("be.visible");
		cy.visit(`/core/deliveryservice/${testData.steeringDS.id}`);
		cy.get("div.actions > mat-icon").should("not.exist");
	});
});
