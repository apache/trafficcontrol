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

import { HttpClient } from "@angular/common/http";
import { HttpClientTestingModule, HttpTestingController } from "@angular/common/http/testing";
import { Injectable } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import type { Observable } from "rxjs";

import { APIService, hasAlerts, type QueryParams } from "./base-api.service";

describe("utilities provided to API services", () => {
	it("correctly determines whether something is an alerts response", () => {
		expect(hasAlerts({})).toBeFalse();
		expect(hasAlerts({alerts: []})).toBeTrue();
		expect(hasAlerts({alerts:[{}]})).toBeFalse();
		expect(hasAlerts({alerts: [{level: "success", text: "succeeded"}]})).toBeTrue();
		expect(hasAlerts({alerts: [{level: "success", text: "succeeded"}, null]})).toBeFalse();
	});
});

/**
 * An implementation of the APIService that just passes method calls through for
 * testing purposes.
 */
@Injectable()
class TestingAPIService extends APIService {

	constructor(httpClient: HttpClient) {
		super(httpClient);
	}

	/**
	 * Calls through to {@link APIService.get} for testing purposes.
	 *
	 * @param path The request path.
	 * @param data Any and all request data.
	 * @param params Any and all query string parameters.
	 * @returns The server's response.
	 */
	public testGet<T>(path: string, data?: object | undefined, params?: QueryParams): Observable<T> {
		return this.get<T>(path, data, params);
	}
	/**
	 * Calls through to {@link APIService.delete} for testing purposes.
	 *
	 * @param path The request path.
	 * @param data Any and all request data.
	 * @param params Any and all query string parameters.
	 * @returns The server's response.
	 */
	public testDelete<T>(path: string, data?: object | undefined, params?: QueryParams): Observable<T> {
		return this.delete<T>(path, data, params);
	}
	/**
	 * Calls through to {@link APIService.put} for testing purposes.
	 *
	 * @param path The request path.
	 * @param data Any and all request data.
	 * @param params Any and all query string parameters.
	 * @returns The server's response.
	 */
	public testPut<T>(path: string, data?: object | undefined, params?: QueryParams): Observable<T> {
		return this.put<T>(path, data, params);
	}
	/**
	 * Calls through to {@link APIService.post} for testing purposes.
	 *
	 * @param path The request path.
	 * @param data Any and all request data.
	 * @param params Any and all query string parameters.
	 * @returns The server's response.
	 */
	public testPost<T>(path: string, data?: object | undefined, params?: QueryParams): Observable<T> {
		return this.post<T>(path, data, params);
	}
	/**
	 * Calls through to {@link APIService.options} for testing purposes.
	 *
	 * @param path The request path.
	 * @param data Any and all request data.
	 * @param params Any and all query string parameters.
	 * @returns The server's response.
	 */
	public testOptions<T>(path: string, data?: object | undefined, params?: QueryParams): Observable<T> {
		return this.options<T>(path, data, params);
	}
	/**
	 * Calls through to {@link APIService.head} for testing purposes.
	 *
	 * @param path The request path.
	 * @param data Any and all request data.
	 * @param params Any and all query string parameters.
	 * @returns The server's response.
	 */
	public testHead<T>(path: string, data?: object | undefined, params?: QueryParams): Observable<T> {
		return this.head<T>(path, data, params);
	}
	/**
	 * Calls through to {@link APIService.patch} for testing purposes.
	 *
	 * @param path The request path.
	 * @param data Any and all request data.
	 * @param params Any and all query string parameters.
	 * @returns The server's response.
	 */
	public testPatch<T>(path: string, data?: object | undefined, params?: QueryParams): Observable<T> {
		return this.patch(path, data, params);
	}
}

describe("API abstract base class", () => {
	const request = {
		request: "data",
		sent: "in the request"
	};

	const params = {
		param: "value"
	};

	const response = {
		response: {
			some: "data"
		}
	};

	let service: TestingAPIService;
	let httpTestingController: HttpTestingController;

	beforeEach(async () => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [
				{provide: APIService, useClass: TestingAPIService},
			]
		});
		service = TestBed.inject(APIService) as TestingAPIService;
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("makes GET requests", async () => {
		const responseP = service.testGet("path", request, params).toPromise();
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/path`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("param")).toBe(params.param);
		expect(req.request.body).toEqual(request);
		req.flush(response);
		await expectAsync(responseP).toBeResolvedTo(response.response);
	});
	it("makes DELETE requests", async () => {
		const responseP = service.testDelete("path", request, params).toPromise();
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/path`);
		expect(req.request.method).toBe("DELETE");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("param")).toBe(params.param);
		expect(req.request.body).toEqual(request);
		req.flush(response);
		await expectAsync(responseP).toBeResolvedTo(response.response);
	});
	it("makes PUT requests", async () => {
		const responseP = service.testPut("path", request, params).toPromise();
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/path`);
		expect(req.request.method).toBe("PUT");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("param")).toBe(params.param);
		expect(req.request.body).toEqual(request);
		req.flush(response);
		await expectAsync(responseP).toBeResolvedTo(response.response);
	});
	it("makes POST requests", async () => {
		const responseP = service.testPost("path", request, params).toPromise();
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/path`);
		expect(req.request.method).toBe("POST");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("param")).toBe(params.param);
		expect(req.request.body).toEqual(request);
		req.flush(response);
		await expectAsync(responseP).toBeResolvedTo(response.response);
	});
	it("makes PATCH requests", async () => {
		const responseP = service.testPatch("path", request, params).toPromise();
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/path`);
		expect(req.request.method).toBe("PATCH");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("param")).toBe(params.param);
		expect(req.request.body).toEqual(request);
		req.flush(response);
		await expectAsync(responseP).toBeResolvedTo(response.response);
	});
	it("makes HEAD requests", async () => {
		const responseP = service.testHead("path", request, params).toPromise();
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/path`);
		expect(req.request.method).toBe("HEAD");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("param")).toBe(params.param);
		expect(req.request.body).toEqual(request);
		req.flush(response);
		await expectAsync(responseP).toBeResolvedTo(response.response);
	});
	it("makes OPTIONS requests", async () => {
		const responseP = service.testOptions("path", request, params).toPromise();
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/path`);
		expect(req.request.method).toBe("OPTIONS");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("param")).toBe(params.param);
		expect(req.request.body).toEqual(request);
		req.flush(response);
		await expectAsync(responseP).toBeResolvedTo(response.response);
	});

	it("throws an error when there's no response body", async () => {
		const responseP = service.testGet("path", request, params).toPromise();
		const req = httpTestingController.expectOne({method: "GET"});
		req.flush(null);

		await expectAsync(responseP).toBeRejected();
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
