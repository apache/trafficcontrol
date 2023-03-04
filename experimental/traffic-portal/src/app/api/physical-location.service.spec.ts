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

import { PhysicalLocationService } from "./physical-location.service";

describe("PhysicalLocationService", () => {
	let service: PhysicalLocationService;
	let httpTestingController: HttpTestingController;
	const physLoc = {
		address: "address",
		city: "city",
		comments: null,
		email: null,
		id: 1,
		lastUpdated: new Date(),
		name: "testquest",
		phone: null,
		poc: null,
		region: null,
		regionId: 2,
		shortName: "testquest",
		state: "state",
		zip: "zip",
	};

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [
				PhysicalLocationService,
			]
		});
		service = TestBed.inject(PhysicalLocationService);
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("gets multiple Physical Locations", async () => {
		const responseP = service.getPhysicalLocations();
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/phys_locations`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(0);
		req.flush({response: [physLoc]});
		await expectAsync(responseP).toBeResolvedTo([physLoc]);
	});

	it("gets a single Physical Location by ID", async () => {
		const responseP = service.getPhysicalLocations(physLoc.id);
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/phys_locations`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("id")).toBe(String(physLoc.id));
		req.flush({response: [physLoc]});
		await expectAsync(responseP).toBeResolvedTo(physLoc);
	});

	it("gets a single Physical Location by name", async () => {
		const responseP = service.getPhysicalLocations(physLoc.name);
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/phys_locations`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("name")).toBe(physLoc.name);
		req.flush({response: [physLoc]});
		await expectAsync(responseP).toBeResolvedTo(physLoc);
	});

	it("submits requests to create new Physical Locations", async () => {
		const responseP = service.createPhysicalLocation(physLoc);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/phys_locations`);
		expect(req.request.method).toBe("POST");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toEqual(physLoc);
		req.flush({response: physLoc});
		await expectAsync(responseP).toBeResolvedTo(physLoc);
	});

	it("submits requests to update existing Physical Locations", async () => {
		const responseP = service.updatePhysicalLocation(physLoc);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/phys_locations/${physLoc.id}`);
		expect(req.request.method).toBe("PUT");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toEqual(physLoc);
		req.flush({response: physLoc});
		await expectAsync(responseP).toBeResolvedTo(physLoc);
	});

	it("submits requests to delete Physical Locations", async () => {
		let responseP = service.deletePhysicalLocation(physLoc);
		let req = httpTestingController.expectOne(`/api/${service.apiVersion}/phys_locations/${physLoc.id}`);
		expect(req.request.method).toBe("DELETE");
		expect(req.request.params.keys().length).toBe(0);
		req.flush({alerts: []});
		await expectAsync(responseP).toBeResolvedTo(undefined);

		responseP = service.deletePhysicalLocation(physLoc.id);
		req = httpTestingController.expectOne(`/api/${service.apiVersion}/phys_locations/${physLoc.id}`);
		expect(req.request.method).toBe("DELETE");
		expect(req.request.params.keys().length).toBe(0);
		req.flush({alerts: []});
		await expectAsync(responseP).toBeResolvedTo(undefined);
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
