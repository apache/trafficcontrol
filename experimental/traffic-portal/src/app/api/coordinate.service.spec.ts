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

import { CoordinateService } from "./coordinate.service";

describe("CoordinateService", () => {
	let service: CoordinateService;
	let httpTestingController: HttpTestingController;
	const coordinate = {
		id: 1,
		lastUpdated: new Date(),
		latitude: 1.0,
		longitude: -1.0,
		name: "test_coordinate",
	};

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [CoordinateService],
		});
		service = TestBed.inject(CoordinateService);
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("gets multiple Coordinates", async () => {
		const responseP = service.getCoordinates();
		const req = httpTestingController.expectOne(
			`/api/${service.apiVersion}/coordinates`
		);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(0);
		req.flush({ response: [coordinate] });
		await expectAsync(responseP).toBeResolvedTo([coordinate]);
	});

	it("gets a single Coordinate by ID", async () => {
		const responseP = service.getCoordinates(coordinate.id);
		const req = httpTestingController.expectOne(
			(r) => r.url === `/api/${service.apiVersion}/coordinates`
		);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("id")).toBe(String(coordinate.id));
		req.flush({ response: [coordinate] });
		await expectAsync(responseP).toBeResolvedTo(coordinate);
	});

	it("gets a single Coordinate by name", async () => {
		const responseP = service.getCoordinates(coordinate.name);
		const req = httpTestingController.expectOne(
			(r) => r.url === `/api/${service.apiVersion}/coordinates`
		);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("name")).toBe(coordinate.name);
		req.flush({ response: [coordinate] });
		await expectAsync(responseP).toBeResolvedTo(coordinate);
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
