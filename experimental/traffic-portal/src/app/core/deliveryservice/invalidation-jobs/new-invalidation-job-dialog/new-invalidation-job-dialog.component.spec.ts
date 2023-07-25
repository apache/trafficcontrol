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
import { HttpClientModule } from "@angular/common/http";
import { type ComponentFixture, TestBed, tick, fakeAsync } from "@angular/core/testing";
import { MatDialogModule, MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";
import { JobType } from "trafficops-types";

import { DeliveryServiceService, InvalidationJobService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";

import { NewInvalidationJobDialogComponent, sanitizedRegExpString, timeStringFromDate } from "./new-invalidation-job-dialog.component";

describe("NewInvalidationJobDialogComponent", () => {
	let component: NewInvalidationJobDialogComponent;
	let fixture: ComponentFixture<NewInvalidationJobDialogComponent>;
	const dialogRef = {
		close: jasmine.createSpy("dialog 'close' method", (): void => { /* Do nothing */ })
	};
	const dialogData = {
		dsID: -1
	};

	beforeEach(async () => {
		dialogRef.close = jasmine.createSpy("dialog 'close' method", (): void => { /* Do nothing */ });
		await TestBed.configureTestingModule({
			declarations: [ NewInvalidationJobDialogComponent ],
			imports: [
				MatDialogModule,
				HttpClientModule,
				APITestingModule
			],
			providers: [
				{provide: MatDialogRef, useValue: dialogRef},
				{provide: MAT_DIALOG_DATA, useValue: dialogData},
			]
		}).compileComponents();
		const service = TestBed.inject(DeliveryServiceService);
		// TODO: These are never cleaned up (because the DS service doesn't have
		// a method for DS deletion)
		const ds = await service.createDeliveryService({
			active: true,
			anonymousBlockingEnabled: false,
			cacheurl: null,
			cdnId: 2,
			displayName: "Test DS",
			dscp: 1,
			geoLimit: 0,
			geoProvider: 0,
			httpBypassFqdn: null,
			infoUrl: null,
			ipv6RoutingEnabled: true,
			logsEnabled: true,
			longDesc: "A DS for testing",
			missLat: 0,
			missLong: 0,
			multiSiteOrigin: false,
			regionalGeoBlocking: false,
			remapText: null,
			routingName: "test",
			tenantId: 1,
			typeId: 10,
			xmlId: "test-ds",
		});
		if (ds.id === undefined) {
			return fail("created Delivery Service had no ID");
		}
		dialogData.dsID = ds.id;

		fixture = TestBed.createComponent(NewInvalidationJobDialogComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("closes the dialog", () => {
		expect(dialogRef.close).not.toHaveBeenCalled();
		component.cancel();
		expect(dialogRef.close).toHaveBeenCalled();
	});

	it("submits requests to create new Jobs, then closes the dialog", fakeAsync(() => {
		expect(dialogRef.close).not.toHaveBeenCalled();
		component.onSubmit(new SubmitEvent("submit"));
		tick();
		expect(dialogRef.close).toHaveBeenCalled();
		const service = TestBed.inject(InvalidationJobService);
		expectAsync((async (): Promise<true> => {
			for (const j of await service.getInvalidationJobs()) {
				await service.deleteInvalidationJob(j.id);
			}
			return true;
		})()).toBeResolvedTo(true);
	}));

	it("updates the minimum starting time according to a newly selected starting date", () => {
		component.startDate = new Date();
		component.startMin = new Date();
		component.dateChange();
		expect(component.startMinTime).toBe(timeStringFromDate(component.startMin));
		component.startDate.setDate(component.startDate.getDate()+1);
		component.dateChange();
		expect(component.startMinTime).toBe("00:00");
	});

	it("doesn't try to create the job when the regexp isn't valid", fakeAsync(() => {
		component.regexp.setValue("+\\y");
		component.onSubmit(new SubmitEvent("submit"));
		tick();
		expect(dialogRef.close).not.toHaveBeenCalled();
	}));
});

describe("NewInvalidationJobDialogComponent - editing", () => {
	let component: NewInvalidationJobDialogComponent;
	let fixture: ComponentFixture<NewInvalidationJobDialogComponent>;
	const dialogRef = {
		close: jasmine.createSpy("dialog 'close' method", (): void => { /* Do nothing */ })
	};
	const dialogData = {
		dsID: -1,
		job: {
			assetUrl: "https://some-url.test/followed/by/a/p\\.attern\\.\\b",
			id: -1,
			invalidationType: JobType.REFRESH,
			startTime: new Date(0),
			ttlHours: 178,
		}
	};

	beforeEach(async () => {
		dialogRef.close = jasmine.createSpy("dialog 'close' method", (): void => { /* Do nothing */ });
		await TestBed.configureTestingModule({
			declarations: [ NewInvalidationJobDialogComponent ],
			imports: [
				MatDialogModule,
				HttpClientModule,
				APITestingModule
			],
			providers: [
				{provide: MatDialogRef, useValue: dialogRef},
				{provide: MAT_DIALOG_DATA, useValue: dialogData},
			]
		}).compileComponents();

		const service = TestBed.inject(InvalidationJobService);
		const now = new Date();
		const job = await service.createInvalidationJob({
			deliveryService: "test-xmlid",
			invalidationType: JobType.REFRESH,
			regex: "/",
			startTime: new Date(now.setDate(now.getDate()+1)),
			ttlHours: 178
		});
		if (job.id === undefined) {
			return fail("created Content Invalidation Job had no ID");
		}
		dialogData.job = job;

		fixture = TestBed.createComponent(NewInvalidationJobDialogComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	afterEach(async ()=>{
		const service = TestBed.inject(InvalidationJobService);
		await service.deleteInvalidationJob(dialogData.job.id);
		expect((await service.getInvalidationJobs()).length).toBe(0);
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("submits requests to create new Jobs, then closes the dialog", fakeAsync(() => {
		expect(dialogRef.close).not.toHaveBeenCalled();
		component.onSubmit(new SubmitEvent("submit"));
		tick();
		expect(dialogRef.close).toHaveBeenCalled();
	}));
});

describe("NewInvalidationJobDialogComponent utility functions", () => {
	it("gets a time string from a Date", ()=>{
		const d = new Date();
		d.setHours(0);
		d.setMinutes(0);
		expect(timeStringFromDate(d)).toBe("00:00");
		d.setHours(12);
		d.setMinutes(34);
		expect(timeStringFromDate(d)).toBe("12:34");
	});
	it("sanitizes regular expressions", ()=>{
		expect(sanitizedRegExpString(/\/.+\/my\/path\.jpg/)).toBe("/.+/my/path\\.jpg");
		expect(sanitizedRegExpString(new RegExp("\\/path\\/to\\/content\\/.+\\.m3u8"))).toBe("/path/to/content/.+\\.m3u8");
	});
});
