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

import { ChangeLogsService } from "./change-logs.service";

describe("ChangeLogsService", () => {
	let service: ChangeLogsService;
	let httpTestingController: HttpTestingController;

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [
				ChangeLogsService,
			]
		});
		service = TestBed.inject(ChangeLogsService);
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("gets Change Logs", async () => {
		const params = {query: "param"};
		const responseP = service.getChangeLogs(params);
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/logs`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("query")).toEqual(params.query);
		req.flush({response: []});
		await expectAsync(responseP).toBeResolvedTo([]);
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
