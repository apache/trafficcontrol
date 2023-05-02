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
import { JobType } from "trafficops-types";

import { testDS } from "./delivery-service.service.spec";
import { InvalidationJobService } from "./invalidation-job.service";

describe("InvalidationJobService", () => {
	let service: InvalidationJobService;
	let httpTestingController: HttpTestingController;
	const job = {
		assetUrl: "asset URL",
		createdBy: "created by",
		deliveryService: testDS.xmlId,
		id: 1,
		invalidationType: JobType.REFETCH,
		startTime: new Date(),
		ttlHours: 5,
	};

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [
				InvalidationJobService,
			]
		});
		service = TestBed.inject(InvalidationJobService);
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("sends requests for OS Versions", async () => {
		const opts = {
			deliveryService: testDS,
			dsID: testDS.id+1,
			id: testDS.id+2,
			user: {
				addressLine1: null,
				addressLine2: null,
				changeLogCount: 0,
				city: null,
				company: null,
				country: null,
				email: "a@b.c" as const,
				fullName: "",
				gid: null,
				id: testDS.id+3,
				lastAuthenticated: null,
				lastUpdated: new Date(),
				newUser: false,
				phoneNumber: null,
				postalCode: null,
				publicSshKey: null,
				registrationSent: null,
				role: "admin",
				stateOrProvince: null,
				tenant: "root",
				tenantId: 1,
				ucdn: "",
				uid: null,
				username: ""
			},
			userId: testDS.id+4
		};
		const responseP = service.getInvalidationJobs(opts);
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/jobs`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(3);
		expect(req.request.params.get("id")).toBe(String(opts.id));
		expect(req.request.params.get("dsId")).toBe(String(testDS.id));
		expect(req.request.params.get("userId")).toBe(String(opts.user.id));
		req.flush({response: [job]});
		await expectAsync(responseP).toBeResolvedTo([job]);
	});

	it("submits requests to create a new content invalidation job", async () => {
		const requestJob = {
			deliveryService: testDS.xmlId,
			invalidationType: JobType.REFETCH,
			regex: "asset URL",
			startTime: job.startTime,
			ttlHours: 5
		};
		const responseP = service.createInvalidationJob(requestJob);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/jobs`);
		expect(req.request.method).toBe("POST");
		expect(req.request.body).toEqual(requestJob);
		expect(req.request.params.keys().length).toBe(0);
		req.flush({response: job});
		await expectAsync(responseP).toBeResolvedTo(job);
	});

	it("submits requests to update an existing content invalidation job", async () => {
		const responseP = service.updateInvalidationJob(job);
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/jobs`);
		expect(req.request.method).toBe("PUT");
		expect(req.request.body).toEqual(job);
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("id")).toBe(String(job.id));
		req.flush({response: job});
		await expectAsync(responseP).toBeResolvedTo(job);
	});

	it("submits requests to delete a job", async () => {
		const responseP = service.deleteInvalidationJob(job);
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/jobs`);
		expect(req.request.method).toBe("DELETE");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("id")).toBe(String(job.id));
		req.flush({response: job});
		await expectAsync(responseP).toBeResolvedTo(job);
	});

	it("submits requests to delete a job by ID", async () => {
		const responseP = service.deleteInvalidationJob(job.id);
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/jobs`);
		expect(req.request.method).toBe("DELETE");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("id")).toBe(String(job.id));
		req.flush({response: job});
		await expectAsync(responseP).toBeResolvedTo(job);
	});
	afterEach(() => {
		httpTestingController.verify();
	});
});
