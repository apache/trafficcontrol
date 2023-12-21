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

import { HttpParams } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { TestBed } from "@angular/core/testing";

import { APITestingService } from "./base-api.service";

/**
 * An implementation of the APITestingService that does nothing more than make
 * it instantiable.
 */
@Injectable()
class TestingAPITestingService extends APITestingService {
}

describe("API abstract base class", () => {
	const request = {
		request: "data",
		sent: "in the request"
	};

	const params = new HttpParams({
		fromObject: {param: "value"}
	});

	const path = "path";

	let service: TestingAPITestingService;

	beforeEach(() => {
		TestBed.configureTestingModule({
			providers: [
				TestingAPITestingService,
			]
		});
		service = TestBed.inject(TestingAPITestingService);
	});

	it("stores requests when params are passed as an object", () => {
		service.get(path, request, {param: "value"});
		const req = service.requestStack.pop();
		if (!req) {
			return fail("expected a GET request pushed onto the stack, but stack was empty");
		}
		expect(service.requestStack.length).toBe(0);
		expect(req.url).toBe(`/api/${service.apiVersion}/${path}`);
		expect(req.method).toBe("GET");
		expect(req.params.keys().length).toEqual(1);
		expect(req.params.get("param")).toBe("value");

		// The Angular HTTP client refuses to send request bodies with GET
		// requests.
		// expect(req.body).toEqual(request);
		expect(req.body).toBeNull();
	});
	it("stores GET requests", () => {
		service.get(path, request, params);
		const req = service.requestStack.pop();
		if (!req) {
			return fail("expected a GET request pushed onto the stack, but stack was empty");
		}
		expect(service.requestStack.length).toBe(0);
		expect(req.url).toBe(`/api/${service.apiVersion}/${path}`);
		expect(req.method).toBe("GET");
		expect(req.params.keys().length).toEqual(1);
		expect(req.params.get("param")).toBe("value");

		// The Angular HTTP client refuses to send request bodies with GET
		// requests.
		// expect(req.body).toEqual(request);
		expect(req.body).toBeNull();
	});
	it("stores DELETE requests", () => {
		service.delete(path, request, params);
		const req = service.requestStack.pop();
		if (!req) {
			return fail("expected a DELETE request pushed onto the stack, but stack was empty");
		}
		expect(service.requestStack.length).toBe(0);
		expect(req.url).toBe(`/api/${service.apiVersion}/${path}`);
		expect(req.method).toBe("DELETE");
		expect(req.params).toEqual(params);

		// The Angular HTTP client refuses to send request bodies with DELETE
		// requests.
		// expect(req.body).toEqual(request);
		expect(req.body).toBeNull();
	});
	it("stores PUT requests", () => {
		service.put(path, request, params);
		const req = service.requestStack.pop();
		if (!req) {
			return fail("expected a PUT request pushed onto the stack, but stack was empty");
		}
		expect(service.requestStack.length).toBe(0);
		expect(req.url).toBe(`/api/${service.apiVersion}/${path}`);
		expect(req.params).toEqual(params);
		expect(req.body).toEqual(request);
		expect(req.method).toBe("PUT");
	});
	it("stores POST requests", () => {
		service.post(path, request, params);
		const req = service.requestStack.pop();
		if (!req) {
			return fail("expected a POST request pushed onto the stack, but stack was empty");
		}
		expect(service.requestStack.length).toBe(0);
		expect(req.url).toBe(`/api/${service.apiVersion}/${path}`);
		expect(req.params.keys().length).toEqual(1);
		expect(req.params.get("param")).toBe("value");
		expect(req.body).toEqual(request);
		expect(req.method).toBe("POST");
	});
	it("stores PATCH requests", () => {
		service.patch(path, request, params);
		const req = service.requestStack.pop();
		if (!req) {
			return fail("expected a PATCH request pushed onto the stack, but stack was empty");
		}
		expect(service.requestStack.length).toBe(0);
		expect(req.url).toBe(`/api/${service.apiVersion}/${path}`);
		expect(req.params).toEqual(params);
		expect(req.body).toEqual(request);
		expect(req.method).toBe("PATCH");
	});
	it("stores HEAD requests", () => {
		service.head(path, request, params);
		const req = service.requestStack.pop();
		if (!req) {
			return fail("expected a HEAD request pushed onto the stack, but stack was empty");
		}
		expect(service.requestStack.length).toBe(0);
		expect(req.url).toBe(`/api/${service.apiVersion}/${path}`);
		expect(req.method).toBe("HEAD");
		expect(req.params).toEqual(params);

		// The Angular HTTP client refuses to send request bodies with HEAD
		// requests.
		// expect(req.body).toEqual(request);
		expect(req.body).toBeNull();
	});
	it("stores OPTIONS requests", () => {
		service.options(path, request, params);
		const req = service.requestStack.pop();
		if (!req) {
			return fail("expected a OPTIONS request pushed onto the stack, but stack was empty");
		}
		expect(service.requestStack.length).toBe(0);
		expect(req.url).toBe(`/api/${service.apiVersion}/${path}`);
		expect(req.method).toBe("OPTIONS");
		expect(req.params).toEqual(params);

		// The Angular HTTP client refuses to send request bodies with OPTIONS
		// requests.
		// expect(req.body).toEqual(request);
		expect(req.body).toBeNull();
	});
});
