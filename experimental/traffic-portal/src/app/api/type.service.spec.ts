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

import { TypeService } from "./type.service";

describe("TypeService", () => {
	let service: TypeService;
	let httpTestingController: HttpTestingController;

	const type = {
		description: "description",
		id: 1,
		lastUpdated: new Date(),
		name: "test type",
		useInTable: "deliveryservice",
	};

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [
				TypeService,
			]
		});
		service = TestBed.inject(TypeService);
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("sends requests for multiple types", async () => {
		const responseP = service.getTypes();
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/types`);
		expect(req.request.method).toBe("GET");
		req.flush({response: [type]});
		await expectAsync(responseP).toBeResolvedTo([type]);
	});
	it("sends requests for a single type by ID", async () => {
		const responseP = service.getTypes(type.id);
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/types`);
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("id")).toBe(String(type.id));
		expect(req.request.method).toBe("GET");
		req.flush({response: [type]});
		await expectAsync(responseP).toBeResolvedTo(type);
	});
	it("sends requests for a single type by name", async () => {
		const responseP = service.getTypes(type.name);
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/types`);
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("name")).toBe(type.name);
		expect(req.request.method).toBe("GET");
		req.flush({response: [type]});
		await expectAsync(responseP).toBeResolvedTo(type);
	});
	it("throws an error when fetching a non-existent type", async () => {
		const responseP = service.getTypes(type.id);
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/types`);
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("id")).toBe(String(type.id));
		expect(req.request.method).toBe("GET");
		req.flush({response: []});
		await expectAsync(responseP).toBeRejected();
	});
	it("sends requests to get types for a specific table", async () => {
		const responseP = service.getTypesInTable("deliveryservice");
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/types`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("useInTable")).toBe("deliveryservice");
		req.flush({response: [type]});
		await expectAsync(responseP).toBeResolvedTo([type]);
	});
	it("sends requests to get server types", async () => {
		const responseP = service.getServerTypes();
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/types`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("useInTable")).toBe("server");
		req.flush({response: [type]});
		await expectAsync(responseP).toBeResolvedTo([type]);
	});
	it("sends request to update a type", async () => {
		const responseP = service.updateType(type);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/types/${type.id}`);
		expect(req.request.method).toBe("PUT");
		expect(req.request.body).toEqual(type);
		req.flush({response: type});
		await expectAsync(responseP).toBeResolvedTo(type);
	});
	it("sends request to create a type", async () => {
		const responseP = service.createType(type);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/types`);
		expect(req.request.method).toBe("POST");
		expect(req.request.body).toEqual(type);
		req.flush({response: type});
		await expectAsync(responseP).toBeResolvedTo(type);
	});
	it("sends request to delete a type", async () => {
		const responseP = service.deleteType(type);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/types/${type.id}`);
		expect(req.request.method).toBe("DELETE");
		req.flush({response: type});
		await expectAsync(responseP).toBeResolvedTo(type);
	});
	it("sends request to delete a type by ID", async () => {
		const responseP = service.deleteType(type.id);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/types/${type.id}`);
		expect(req.request.method).toBe("DELETE");
		req.flush({response: type});
		await expectAsync(responseP).toBeResolvedTo(type);
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
