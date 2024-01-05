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

describe("Servers edit/creation page", () => {
	beforeEach(() => {
		cy.login();
	});
	it("Loads a new server form correctly", () => {
		cy.visit("core/servers/new");
		cy.get("input[name=hostname]").should("be.enabled").and("have.value", "");
		cy.get("mat-select[name=cdn]").should("exist");
		cy.get("mat-select[name=cachegroup]").should("exist");
		cy.get("mat-select[name=physLocation]").should("exist");
		cy.get("mat-select[name=status]").should("exist");
		cy.get("input[name=offlineReason]").should("not.exist");
		cy.get("mat-select[name=type]").should("exist");
		cy.get("input[name=httpport]").should("be.enabled");
		cy.get("input[name=httpsport]").should("be.enabled");
		cy.get("input[name=rack]").should("be.enabled");
		cy.get("input[name=serverId]").should("not.exist");
		cy.get("input[name=lastUpdated]").should("not.exist");
		cy.get("input[name=statusLastUpdated]").should("not.exist");
		cy.get("mat-select[name=profiles]").should("exist");
		cy.get("button[aria-label='Add An Interface']").should("not.be.disabled");
		cy.get("input[name=iloIP]").should("be.enabled");
		cy.get("input[name=iloGateway]").should("be.enabled");
		cy.get("input[name=iloNetmask]").should("be.enabled");
		cy.get("input[name=iloUsername]").should("be.enabled");
		cy.get("input[name=iloPassword]").should("be.enabled");
		cy.get("input[name=mgmtIP]").should("be.enabled");
		cy.get("input[name=mgmtIpGateway]").should("be.enabled");
		cy.get("input[name=mgmtIpNetmask]").should("be.enabled");
		cy.get("button[aria-label='Delete Server']").should("not.exist");
		cy.get("button[aria-label='Submit Server']").should("not.be.disabled");
	});

	it("Loads an existing server's edit form correctly", () => {
		cy.fixture("test.data").then(
			(data: CreatedData) => {
				const {edgeServer} = data;
				cy.visit(`/core/servers/${edgeServer.id}`);
				cy.get("input[name=hostname]").should("be.enabled").and("have.value", edgeServer.hostName);
				cy.get("mat-select[name=cdn]").should("contain.text", edgeServer.cdnName);
				cy.get("mat-select[name=cachegroup]").should("contain.text", edgeServer.cachegroup);
				cy.get("mat-select[name=physLocation]").should("contain.text", edgeServer.physLocation);
				cy.get("input[name=status]").should("be.disabled").and("have.value", edgeServer.status);
				cy.get("input[name=offlineReason]").should("not.exist");
				cy.get("mat-select[name=type]").should("contain.text", edgeServer.type);
				cy.get("input[name=httpport]").should("be.enabled").and("have.value", edgeServer.tcpPort ?? "");
				cy.get("input[name=httpsport]").should("be.enabled").and("have.value", edgeServer.httpsPort ?? "");
				cy.get("input[name=rack]").should("be.enabled").and("have.value", edgeServer.rack ?? "");
				cy.get("input[name=serverId]").should("be.disabled").and("have.value", edgeServer.id);
				cy.get("input[name=lastUpdated]").should("be.disabled");
				// TODO: verified that this field doesn't show up until the
				// server's status actually changed - so how were the Nightwatch
				// tests passing?
				// cy.get("input[name=statusLastUpdated]").should("be.disabled");
				cy.get("mat-select[name=profiles]").should("contain.text", edgeServer.profileNames[0]);
				cy.get("button[aria-label='Add An Interface']").should("not.be.disabled");
				cy.get("input[name=iloIP]").should("be.enabled").and("have.value", edgeServer.iloIpAddress ?? "");
				cy.get("input[name=iloGateway]").should("be.enabled").and("have.value", edgeServer.iloIpGateway ?? "" );
				cy.get("input[name=iloNetmask]").should("be.enabled").and("have.value", edgeServer.iloIpNetmask ?? "");
				cy.get("input[name=iloUsername]").should("be.enabled").and("have.value", edgeServer.iloUsername ?? "");
				cy.get("input[name=iloPassword]").should("be.enabled").and("have.value", edgeServer.iloPassword ?? "");
				cy.get("input[name=mgmtIP]").should("be.enabled").and("have.value", edgeServer.mgmtIpAddress ?? "");
				cy.get("input[name=mgmtIpGateway]").should("be.enabled").and("have.value", edgeServer.mgmtIpGateway ?? "");
				cy.get("input[name=mgmtIpNetmask]").should("be.enabled").and("have.value", edgeServer.mgmtIpNetmask ?? "");
				cy.get("button[aria-label='Delete Server']").should("not.be.disabled");
				cy.get("button[aria-label='Submit Server']").should("not.be.disabled");
			}
		);
	});
});
