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

import { HttpClient, HTTP_INTERCEPTORS } from "@angular/common/http";
import { HttpClientTestingModule, HttpTestingController } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";

import { DateReviverInterceptor } from "./date-reviver.interceptor";

/**
 * Holds some Date instance data fields for testing Date revival.
 */
interface TestResponse {
	lastAuthenticated: Date;
	lastUpdated: Date;
	statusLastUpdated: Date;
}

describe("DateReviverInterceptor", () => {
	let httpClient: HttpClient;
	let httpTestingController: HttpTestingController;

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [
				HttpClientTestingModule
			],
			providers: [
				DateReviverInterceptor,
				{ multi: true, provide: HTTP_INTERCEPTORS, useClass: DateReviverInterceptor}
			]
		});
		// Inject the http service and test controller for each test
		httpClient = TestBed.inject(HttpClient);
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("should be created", () => {
		const interceptor: DateReviverInterceptor = TestBed.inject(DateReviverInterceptor);
		expect(interceptor).toBeTruthy();
	});

	it("revives dates sent in responses", () => {
		const testData = {
			lastAuthenticated: "2022-03-31T08:01:02.000+00",
			lastUpdated: "2022-03-31 08:01:02+00",
			statusLastUpdated: "2022-03-31T08:01:02Z"
		};

		httpClient.get<TestResponse>("/api/data").subscribe(
			data => {
				expect(typeof data).toBe("object");

				const dataProperties = Object.keys(data);
				expect(dataProperties.length).toBe(3);
				expect(dataProperties).toContain("lastAuthenticated");
				expect(data.lastAuthenticated).toBeInstanceOf(Date);
				expect(dataProperties).toContain("lastUpdated");
				expect(data.lastUpdated).toBeInstanceOf(Date);
				expect(dataProperties).toContain("statusLastUpdated");
				expect(data.statusLastUpdated).toBeInstanceOf(Date);

				const expectedDate = new Date(Date.UTC(2022, 2, 31, 8, 1, 2, 0));

				expect(data.lastAuthenticated).toEqual(expectedDate);
				expect(data.lastUpdated).toEqual(expectedDate);
				expect(data.statusLastUpdated).toEqual(expectedDate);
			}
		);

		httpTestingController.expectOne("/api/data").flush(testData);
	});

	it("doesn't modify non-date properties of responses", () => {
		const testData = {
			arr: [],
			bool: true,
			num: 1,
			obj: {},
			str: "some string that doesn't look like a date",
			strNum: "12345"
		};

		httpClient.get<typeof testData>("/api/data").subscribe(
			data => {
				expect(data).toEqual(testData);
			}
		);

		httpTestingController.expectOne("/api/data").flush(testData);
	});

	it("doesn't modify responses where JSON parsing wasn't requested", ()=>{
		const testData = JSON.stringify({
			lastAuthenticated: "2022-03-31T08:01:02.000+00",
			lastUpdated: "2022-03-31 08:01:02+00",
			statusLastUpdated: "2022-03-31T08:01:02Z"
		});

		httpClient.get("/api/data", {responseType: "text"}).subscribe(
			data => {
				expect(data).toEqual(testData);
			}
		);

		httpTestingController.expectOne("/api/data").flush(testData);
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
