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
import { TestBed } from "@angular/core/testing";
import { JobType, type ResponseDeliveryService } from "trafficops-types";

import { DeliveryServiceService } from "../delivery-service.service";

import { InvalidationJobService } from "./invalidation-job.service";

import { APITestingModule } from ".";

describe("InvalidationJobService", () => {
	let service: InvalidationJobService;
	let ds: ResponseDeliveryService;

	beforeEach(async () => {
		TestBed.configureTestingModule({
			imports: [APITestingModule],
			providers: [
				InvalidationJobService,
			]
		});
		service = TestBed.inject(InvalidationJobService);
		const xmlId = `DS-with-Profile Creation Test-${(new Date()).valueOf()}`;
		ds = await TestBed.inject(DeliveryServiceService).createDeliveryService({
			active: true,
			cacheurl: null,
			cdnId: 1,
			displayName: xmlId,
			dscp: 2,
			geoLimit: 0,
			geoProvider: 0,
			httpBypassFqdn: null,
			infoUrl: null,
			logsEnabled: true,
			regionalGeoBlocking: false,
			remapText: null,
			tenantId: 1,
			typeId: 1,
			xmlId,
		});

		let job = await service.createInvalidationJob({
			deliveryService: ds.xmlId,
			invalidationType: JobType.REFRESH,
			regex: "",
			startTime: "2123-04-05T06:07:08Z",
			ttlHours: 5
		});
		expect(job.id).toBeTruthy();

		job = await service.createInvalidationJob({
			deliveryService: "N/A",
			invalidationType: JobType.REFETCH,
			regex: "^.+\\.jpg$",
			startTime: new Date(),
			ttlHours: 170
		});
		expect(job.id).toBeTruthy();

		job = await service.createInvalidationJob({
			deliveryService: ds.xmlId,
			invalidationType: JobType.REFETCH,
			regex: /^.+\.png$/.source,
			startTime: new Date(),
			ttlHours: 9001
		});
		expect(job.id).toBeTruthy();
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("gets invalidation jobs", async () => {
		expect((await service.getInvalidationJobs()).length).toBeGreaterThanOrEqual(3);
	});

	it("filters invalidation jobs by DSID", async () => {
		const jobs = await service.getInvalidationJobs({dsID: ds.id});
		expect(jobs).toHaveSize(2);
		for (const job of jobs) {
			expect(job.deliveryService).toBe(ds.xmlId);
		}
	});

	it("filters invalidation jobs by Delivery Service", async () => {
		const jobs = await service.getInvalidationJobs({deliveryService: ds});
		expect(jobs).toHaveSize(2);
		for (const job of jobs) {
			expect(job.deliveryService).toBe(ds.xmlId);
		}
	});

	it("filters invalidation jobs by username", async () => {
		const user = {
			addressLine1: null,
			addressLine2: null,
			changeLogCount: 0,
			city: null,
			company: null,
			country: null,
			email: "test@que.st" as const,
			fullName: "",
			gid: null,
			id: 1,
			lastAuthenticated: new Date(),
			lastUpdated: new Date(),
			newUser: false,
			phoneNumber: null,
			postalCode: null,
			publicSshKey: null,
			registrationSent: null,
			role: "",
			stateOrProvince: null,
			tenant: "",
			tenantId: 1,
			ucdn: "",
			uid: null,
			username: `test-${(new Date()).valueOf()}`
		};
		await expectAsync(service.getInvalidationJobs({user})).toBeResolvedTo([]);
	});

	it("throws an error if using the currently unsuported User ID filtering", async () => {
		await expectAsync(service.getInvalidationJobs({userId: -1})).toBeRejected();
	});

	it("updates an existing content invalidation job", async () => {
		const job = {...(await service.getInvalidationJobs())[0]};
		const initialDate = new Date(job.startTime.valueOf());
		job.startTime.setDate(initialDate.getDate()+1);
		const updated = await service.updateInvalidationJob(job);
		await expectAsync(service.getInvalidationJobs({id: job.id})).toBeResolvedTo([updated]);
	});

	it("throws an error when asked to update a Job that doesn't exist", async () => {
		const job = {
			assetUrl: "",
			createdBy: "",
			deliveryService: "",
			id: -1,
			invalidationType: JobType.REFETCH,
			startTime: new Date(),
			ttlHours: 0
		};
		await expectAsync(service.updateInvalidationJob(job)).toBeRejected();
	});

	it("deletes jobs", async () => {
		const initialJobs = await service.getInvalidationJobs();
		const initialLength = initialJobs.length;
		const job = initialJobs[0];
		const deleted = await service.deleteInvalidationJob(job);
		expect(job).toEqual(deleted);
		const updatedJobs = await service.getInvalidationJobs();
		expect(updatedJobs).not.toContain(job);
		expect(updatedJobs.length).toEqual(initialLength-1);
	});

	it("deletes jobs by ID", async () => {
		const initialJobs = await service.getInvalidationJobs();
		const initialLength = initialJobs.length;
		const job = initialJobs[0];
		const deleted = await service.deleteInvalidationJob(job.id);
		expect(job).toEqual(deleted);
		const updatedJobs = await service.getInvalidationJobs();
		expect(updatedJobs).not.toContain(job);
		expect(updatedJobs.length).toEqual(initialLength-1);
	});

	it("throws an error when asked to delete a Job that doesn't exist", async () => {
		await expectAsync(service.deleteInvalidationJob(-1)).toBeRejected();
	});
});
