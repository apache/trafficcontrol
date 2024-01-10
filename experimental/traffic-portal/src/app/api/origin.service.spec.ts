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
import {
	HttpClientTestingModule,
	HttpTestingController,
} from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";

import { OriginService } from "./origin.service";

describe("OriginService", () => {
	let service: OriginService;
	let httpTestingController: HttpTestingController;
	const origin = {
		cachegroup: null,
		cachegroupId: null,
		coordinate: null,
		coordinateId: null,
		deliveryService: "test",
		deliveryServiceId: 1,
		fqdn: "origin.infra.ciab.test",
		id: 1,
		ip6Address: null,
		ipAddress: null,
		isPrimary: false,
		lastUpdated: new Date(),
		name: "test",
		port: null,
		profile: null,
		profileId: null,
		protocol: "http" as never,
		tenant: "root",
		tenantId: 1,
	};

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [OriginService],
		});
		service = TestBed.inject(OriginService);
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("gets multiple Origins", async () => {
		const responseP = service.getOrigins();
		const req = httpTestingController.expectOne(
			`/api/${service.apiVersion}/origins`
		);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(0);
		req.flush({ response: [origin] });
		await expectAsync(responseP).toBeResolvedTo([origin]);
	});

	it("gets a single Origin by ID", async () => {
		const responseP = service.getOrigins(origin.id);
		const req = httpTestingController.expectOne(
			(r) => r.url === `/api/${service.apiVersion}/origins`
		);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("id")).toBe(String(origin.id));
		req.flush({ response: [origin] });
		await expectAsync(responseP).toBeResolvedTo(origin);
	});

	it("gets a single Origin by name", async () => {
		const responseP = service.getOrigins(origin.name);
		const req = httpTestingController.expectOne(
			(r) => r.url === `/api/${service.apiVersion}/origins`
		);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("name")).toBe(origin.name);
		req.flush({ response: [origin] });
		await expectAsync(responseP).toBeResolvedTo(origin);
	});

	it("submits requests to create new Origins", async () => {
		const responseP = service.createOrigin({
			...origin,
			tenantID: origin.tenantId,
		});
		const req = httpTestingController.expectOne(
			`/api/${service.apiVersion}/origins`
		);
		expect(req.request.method).toBe("POST");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body.name).toEqual(origin.name);
		req.flush({ response: origin });
		await expectAsync(responseP).toBeResolvedTo(origin);
	});

	it("submits requests to update existing Origins", async () => {
		const responseP = service.updateOrigin(origin);
		const req = httpTestingController.expectOne(
			`/api/${service.apiVersion}/origins?id=${origin.id}`
		);
		expect(req.request.method).toBe("PUT");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toEqual(origin);
		req.flush({ response: origin });
		await expectAsync(responseP).toBeResolvedTo(origin);
	});

	it("submits requests to delete Origins", async () => {
		let responseP = service.deleteOrigin(origin);
		let req = httpTestingController.expectOne(
			`/api/${service.apiVersion}/origins?id=${origin.id}`
		);
		expect(req.request.method).toBe("DELETE");
		expect(req.request.params.keys().length).toBe(0);
		req.flush({ alerts: [] });
		await expectAsync(responseP).toBeResolved();

		responseP = service.deleteOrigin(origin.id);
		req = httpTestingController.expectOne(
			`/api/${service.apiVersion}/origins?id=${origin.id}`
		);
		expect(req.request.method).toBe("DELETE");
		expect(req.request.params.keys().length).toBe(0);
		req.flush({ alerts: [] });
		await expectAsync(responseP).toBeResolved();
	});

	it("sends requests for multiple origins by ID", async () => {
		const responseParams = service.getOrigins(origin.id);
		const req = httpTestingController.expectOne(
			(r) => r.url === `/api/${service.apiVersion}/origins`
		);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("id")).toBe(String(origin.id));
		const data = {
			response: [
				{
					cachegroup: null,
					cachegroupId: null,
					coordinate: null,
					coordinateId: null,
					deliveryService: "test",
					deliveryServiceId: 1,
					fqdn: "origin.infra.ciab.test",
					id: 1,
					ip6Address: null,
					ipAddress: null,
					isPrimary: false,
					lastUpdated: new Date(),
					name: "test",
					port: null,
					profile: null,
					profileId: null,
					protocol: "http" as never,
					tenant: "root",
					tenantId: 1,
				},
				{
					cachegroup: null,
					cachegroupId: null,
					coordinate: null,
					coordinateId: null,
					deliveryService: "test",
					deliveryServiceId: 1,
					fqdn: "origin.infra.ciab.test",
					id: 1,
					ip6Address: null,
					ipAddress: null,
					isPrimary: false,
					lastUpdated: new Date(),
					name: "test2",
					port: null,
					profile: null,
					profileId: null,
					protocol: "http" as never,
					tenant: "root",
					tenantId: 1,
				},
			],
		};
		req.flush(data);
		await expectAsync(responseParams).toBeRejectedWithError(
			`Traffic Ops responded with 2 Origins by identifier ${origin.id}`
		);
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
