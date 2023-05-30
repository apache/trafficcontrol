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
import type { ResponseCDN } from "trafficops-types";

import { CDNService } from "./cdn.service";

describe("CDNService", () => {
	let service: CDNService;
	let httpTestingController: HttpTestingController;
	const cdn = {
		dnssecEnabled: false,
		domainName: "test.quest",
		id: 1,
		lastUpdated: new Date(),
		name: "TestQuest",
	};

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [
				CDNService,
			]
		});
		service = TestBed.inject(CDNService);
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("gets multiple CDNs", async () => {
		const responseP = service.getCDNs();
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cdns`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(0);
		req.flush({response: [cdn]});
		await expectAsync(responseP).toBeResolvedTo([cdn]);
	});

	it("gets a single CDN by ID", async () => {
		const responseP = service.getCDNs(cdn.id);
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/cdns`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("id")).toBe(String(cdn.id));
		req.flush({response: [cdn]});
		await expectAsync(responseP).toBeResolvedTo(cdn);
	});

	it("throws an error when more than one CDN has a given ID", async () => {
		const responseP = service.getCDNs(cdn.id);
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/cdns`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("id")).toBe(String(cdn.id));
		req.flush({response: [cdn, {}]});
		await expectAsync(responseP).toBeRejected();
	});

	it("creates a new CDN", async () => {
		const responseP = service.createCDN(cdn);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cdns`);
		expect(req.request.method).toBe("POST");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toBe(cdn);
		req.flush({response: cdn});
		await expectAsync(responseP).toBeResolved();
	});

	it("updates an existing CDN", async () => {
		const responseP = service.updateCDN(cdn);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cdns/${cdn.id}`);
		expect(req.request.method).toBe("PUT");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toBe(cdn);
		req.flush({response: cdn});
		await expectAsync(responseP).toBeResolved();
	});

	it("updates an existing CDN by ID", async () => {
		const responseP = service.updateCDN(cdn.id, cdn);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cdns/${cdn.id}`);
		expect(req.request.method).toBe("PUT");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toBe(cdn);
		req.flush({response: cdn});
		await expectAsync(responseP).toBeResolved();
	});

	it("throws an error for invalid call signatures to updateCDN", async () => {
		const responseP = (service as unknown as {updateCDN: (id: number) => Promise<ResponseCDN>}).updateCDN(cdn.id);
		httpTestingController.expectNone({method: "PUT"});
		await expectAsync(responseP).toBeRejected();
	});

	it("deletes an existing CDN", async () => {
		const responseP = service.deleteCDN(cdn);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cdns/${cdn.id}`);
		expect(req.request.method).toBe("DELETE");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toBeNull();
		req.flush({response: cdn});
		await expectAsync(responseP).toBeResolved();
	});

	it("deletes an existing CDN by ID", async () => {
		const responseP = service.deleteCDN(cdn.id);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cdns/${cdn.id}`);
		expect(req.request.method).toBe("DELETE");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toBeNull();
		req.flush({response: cdn});
		await expectAsync(responseP).toBeResolved();
	});

	it("Queues Updates by CDN", async () => {
		const resp = service.queueServerUpdates(cdn);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cdns/${cdn.id}/queue_update`);
		expect(req.request.method).toBe("POST");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toEqual({action: "queue"});
		req.flush({response: { action: "queue", cdnId: cdn.id }});
		await expectAsync(resp).toBeResolved();
	});

	it("Queues Updates by CDN ID", async () => {
		const resp = service.queueServerUpdates(cdn.id);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cdns/${cdn.id}/queue_update`);
		expect(req.request.method).toBe("POST");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toEqual({action: "queue"});
		req.flush({response: { action: "queue", cdnId: cdn.id }});
		await expectAsync(resp).toBeResolved();
	});

	it("Dequeues Updates by CDN", async () => {
		const resp = service.dequeueServerUpdates(cdn);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cdns/${cdn.id}/queue_update`);
		expect(req.request.method).toBe("POST");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toEqual({action: "dequeue"});
		req.flush({response: { action: "dequeue", cdnId: cdn.id }});
		await expectAsync(resp).toBeResolved();
	});

	it("Dequeues Updates by CDN ID", async () => {
		const resp = service.dequeueServerUpdates(cdn.id);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cdns/${cdn.id}/queue_update`);
		expect(req.request.method).toBe("POST");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toEqual({action: "dequeue"});
		req.flush({response: { action: "dequeue", cdnId: cdn.id }});
		await expectAsync(resp).toBeResolved();
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
