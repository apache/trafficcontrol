/**
 * @license Apache-2.0
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
import { HttpClientTestingModule, HttpTestingController } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { type ResponseServer } from "trafficops-types";

import { ServerService } from "./server.service";

describe("ServerService", () => {
	let service: ServerService;
	let httpTestingController: HttpTestingController;
	const server: ResponseServer = {
		cachegroup: "cachegroup",
		cachegroupId: 1,
		cdnId: 1,
		cdnName: "cdnName",
		domainName: "quest",
		guid: null,
		hostName: "test",
		httpsPort: null,
		id: 1,
		iloIpAddress: null,
		iloIpGateway: null,
		iloIpNetmask: null,
		iloPassword: null,
		iloUsername: null,
		interfaces: [],
		lastUpdated: new Date(),
		mgmtIpAddress: null,
		mgmtIpGateway: null,
		mgmtIpNetmask: null,
		offlineReason: null,
		physLocation: "physicalLocation",
		physLocationId: 1,
		profileNames: ["profile"],
		rack: null,
		revalPending: false,
		routerHostName: null,
		routerPortName: null,
		status: "ONLINE",
		statusId: 1,
		statusLastUpdated: null,
		tcpPort: null,
		type: "EDGE",
		typeId: 1,
		updPending: false,
		xmppId: "xmppId",
	};

	const status = {
		description: "description",
		id: 1,
		lastUpdated: new Date(),
		name: "testquest",
	};

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [
				ServerService,
			]
		});
		service = TestBed.inject(ServerService);
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	describe("server methods", () => {
		it("sends requests for multiple servers", async () => {
			const responseP = service.getServers();
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/servers`);
			expect(req.request.method).toBe("GET");
			req.flush({response: [server]});
			await expectAsync(responseP).toBeResolvedTo([server]);
		});
		it("sends requests for a single server by ID", async () => {
			const responseP = service.getServers(server.id);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/servers`);
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(server.id));
			expect(req.request.method).toBe("GET");
			req.flush({response: [server]});
			await expectAsync(responseP).toBeResolvedTo(server);
		});
		it("sends requests for a single server by hostName", async () => {
			const responseP = service.getServers(server.hostName);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/servers`);
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("hostName")).toBe(server.hostName);
			expect(req.request.method).toBe("GET");
			req.flush({response: [server]});
			await expectAsync(responseP).toBeResolvedTo(server);
		});
		it("returns a single server when TO presents it with multiple matches by hostName", async () => {
			const responseP = service.getServers(server.hostName);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/servers`);
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("hostName")).toBe(server.hostName);
			expect(req.request.method).toBe("GET");
			req.flush({response: [server, {...server, id: server.id + 1}]});
			await expectAsync(responseP).toBeResolvedTo(server);
		});
		it("throws an error when fetching a non-existent server", async () => {
			const responseP = service.getServers(server.id);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/servers`);
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(server.id));
			expect(req.request.method).toBe("GET");
			req.flush({response: []});
			await expectAsync(responseP).toBeRejected();
		});
		it("sends requests for creating a new server", async () => {
			const responseP = service.createServer(server);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/servers`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual(server);
			req.flush({response: server});
			await expectAsync(responseP).toBeResolvedTo(server);
		});

		it("updates a server by ID", async ()  => {
			const resp = service.updateServer(server.id, server);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/servers/${server.id}`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual(server);

			req.flush({response: server});
			await expectAsync(resp).toBeResolvedTo(server);
		});

		it("updates a server", async ()  => {
			const resp = service.updateServer(server);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/servers/${server.id}`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual(server);

			req.flush({response: server});
			await expectAsync(resp).toBeResolvedTo(server);
		});

		it("delete a server", async ()  => {
			const resp = service.deleteServer(server);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/servers/${server.id}`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.body).toBeNull();

			req.flush({response: server});
			await expectAsync(resp).toBeResolvedTo(server);
		});

		it("delete a server by ID", async ()  => {
			const resp = service.deleteServer(server.id);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/servers/${server.id}`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.body).toBeNull();

			req.flush({response: server});
			await expectAsync(resp).toBeResolvedTo(server);
		});
	});

	describe("Status-related methods", () => {
		it("sends requests for multiple statuses", async () => {
			const responseP = service.getStatuses();
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/statuses`);
			expect(req.request.method).toBe("GET");
			req.flush({response: [status]});
			await expectAsync(responseP).toBeResolvedTo([status]);
		});
		it("sends requests for a single status by ID", async () => {
			const responseP = service.getStatuses(status.id);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/statuses`);
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(status.id));
			expect(req.request.method).toBe("GET");
			req.flush({response: [status]});
			await expectAsync(responseP).toBeResolvedTo(status);
		});
		it("sends requests for a single status by name", async () => {
			const responseP = service.getStatuses(status.name);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/statuses`);
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("name")).toBe(status.name);
			expect(req.request.method).toBe("GET");
			req.flush({response: [status]});
			await expectAsync(responseP).toBeResolvedTo(status);
		});
		it("throws an error when fetching a non-existent status", async () => {
			const responseP = service.getStatuses(status.id);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/statuses`);
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(status.id));
			expect(req.request.method).toBe("GET");
			req.flush({response: []});
			await expectAsync(responseP).toBeRejected();
		});
		it("updates a server's status", async () => {
			const responseP = service.updateStatus(server, status);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/servers/${server.id}/status`);
			expect(req.request.method).toBe("PUT");
			// We specifically don't care whether the offline reason truly isn't
			// specified or has the value `undefined`. Actually, `null` would
			// also work for the API, but typings forbid that.
			expect([{offlineReason: undefined, status: status.name}, {status: status.name}]).toContain(req.request.body);
			req.flush({});
			await expectAsync(responseP).toBeResolved();
		});
		it("updates a server's status with an offlineReason", async () => {
			const offlineReason = "because I told you to";
			const responseP = service.updateStatus(server.id, status.name, offlineReason);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/servers/${server.id}/status`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual({offlineReason, status: status.name});
			req.flush({});
			await expectAsync(responseP).toBeResolved();
		});
		it("creates a new status", async () => {
			const responseP = service.createStatus(status);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/statuses`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual(status);
			req.flush({response: status});
			await expectAsync(responseP).toBeResolvedTo(status);
		});
		it("updates an existing status", async () => {
			const responseP = service.updateStatusDetail(status);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/statuses/${status.id}`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual(status);
			req.flush({response: status});
			await expectAsync(responseP).toBeResolvedTo(status);
		});
		it("deletes statuses", async () => {
			const responseP = service.deleteStatus(status);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/statuses/${status.id}`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.body).toBeNull();
			req.flush({response: status});
			await expectAsync(responseP).toBeResolvedTo(status);
		});
		it("deletes statuses by ID", async () => {
			const responseP = service.deleteStatus(status.id);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/statuses/${status.id}`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.body).toBeNull();
			req.flush({response: status});
			await expectAsync(responseP).toBeResolvedTo(status);
		});
	});

	describe("other methods", () => {
		const serverCheck = {
			adminState: "ONLINE",
			cacheGroup: "cacheGroup",
			hostName: "hostName",
			id: 1,
			profile: "profile",
			revalPending: false,
			type: "EDGE",
			updPending: false,
		};

		it("queues updates on a server", async () => {
			const responseP = service.queueUpdates(server);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/servers/${server.id}/queue_update`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual({action: "queue"});
			const response = {action: "queue" as const, serverId: server.id};
			req.flush({response});
			await expectAsync(responseP).toBeResolvedTo(response);
		});
		it("queues updates on a server by ID", async () => {
			const responseP = service.queueUpdates(server.id);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/servers/${server.id}/queue_update`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual({action: "queue"});
			const response = {action: "queue" as const, serverId: server.id};
			req.flush({response});
			await expectAsync(responseP).toBeResolvedTo(response);
		});
		it("de-queues updates on a server", async () => {
			const responseP = service.clearUpdates(server);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/servers/${server.id}/queue_update`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual({action: "dequeue"});
			const response = {action: "dequeue" as const, serverId: server.id};
			req.flush({response});
			await expectAsync(responseP).toBeResolvedTo(response);
		});
		it("de-queues updates on a server by ID", async () => {
			const responseP = service.clearUpdates(server.id);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/servers/${server.id}/queue_update`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual({action: "dequeue"});
			const response = {action: "dequeue" as const, serverId: server.id};
			req.flush({response});
			await expectAsync(responseP).toBeResolvedTo(response);
		});
		it("sends a request for multiple Serverchecks", async () => {
			const responseP = service.getServerChecks();
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/servercheck`);
			expect(req.request.method).toBe("GET");
			const response = [serverCheck];
			req.flush({response});
			await expectAsync(responseP).toBeResolvedTo(response);
		});
		it("gets Serverchecks for a single server by ID", async () => {
			const responseP = service.getServerChecks(server.id);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/servercheck`);
			expect(req.request.method).toBe("GET");
			const response = [{...serverCheck, id: 2}, serverCheck];
			req.flush({response});
			await expectAsync(responseP).toBeResolvedTo(response[1]);
		});
		it("throws an error if the checks for the requested server cannot be found", async () => {
			const responseP = service.getServerChecks(server.id);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/servercheck`);
			expect(req.request.method).toBe("GET");
			const response = [{...serverCheck, id: server.id+1}];
			req.flush({response});
			await expectAsync(responseP).toBeRejected(response[1]);
		});
	});

	describe("static methods", () => {
		it("finds a service address", () => {
			const infs = [
				{
					ipAddresses: [
						{
							address: "",
							gateway: "",
							serviceAddress: false
						}
					],
					maxBandwidth: null,
					monitor: false,
					mtu: null,
					name: "eth0"
				},
				{
					ipAddresses: [
						{
							address: "",
							gateway: "",
							serviceAddress: false
						},
						{
							address: "",
							gateway: "",
							serviceAddress: false
						}
					],
					maxBandwidth: null,
					monitor: false,
					mtu: null,
					name: "eth1"
				},
				{
					ipAddresses: [
						{
							address: "",
							gateway: "",
							serviceAddress: false
						},
						{
							address: "",
							gateway: "",
							serviceAddress: true
						},
					],
					maxBandwidth: null,
					monitor: false,
					mtu: null,
					name: "eth2"
				}
			];
			const serviceInf = ServerService.getServiceInterface(infs);
			expect(serviceInf).toBe(infs[2]);
		});
		it("throws an error when a server has no service addresses", () => {
			expect(()=>ServerService.getServiceInterface({
				cachegroupId: -1,
				cdnId: -1,
				domainName: "",
				hostName: "",
				interfaces: [{
					ipAddresses: [
						{
							address: "",
							gateway: "",
							serviceAddress: false
						},
						{
							address: "",
							gateway: "",
							serviceAddress: false
						}
					],
					maxBandwidth: null,
					monitor: false,
					mtu: null,
					name: "eth0"
				}],
				physLocationId: -1,
				profileNames: [],
				statusId: -1,
				typeId: -1,
			})).toThrow();
		});
		it("extracts netmasks", () => {
			const [addr, netmask] = ServerService.extractNetmask("192.168.0.1/16");
			expect(addr).toBe("192.168.0.1");
			expect(netmask).toBe("255.255.0.0");
		});
		it("doesn't break when a plain address (no CIDR suffix) is passed", () => {
			const raw = "192.168.0.1";
			const [addr, netmask] = ServerService.extractNetmask(raw);
			expect(addr).toBe(raw);
			expect(netmask).toBeUndefined();
		});
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
