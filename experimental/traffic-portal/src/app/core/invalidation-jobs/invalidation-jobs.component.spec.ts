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
import { type ComponentFixture, TestBed } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatDialog, MatDialogModule } from "@angular/material/dialog";
import { Router } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { type Observable, of } from "rxjs";

import { DeliveryServiceService, InvalidationJobService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { defaultDeliveryService, type InvalidationJob, JobType } from "src/app/models";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import { TpHeaderComponent } from "src/app/shared/tp-header/tp-header.component";
import { CustomvalidityDirective } from "src/app/shared/validation/customvalidity.directive";

import { InvalidationJobsComponent } from "./invalidation-jobs.component";

describe("InvalidationJobsComponent", () => {
	let component: InvalidationJobsComponent;
	let fixture: ComponentFixture<InvalidationJobsComponent>;
	let router: Router;
	let job: InvalidationJob;

	beforeEach(async () => {
		// mock the API
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);

		await TestBed.configureTestingModule({
			declarations: [
				InvalidationJobsComponent,
				TpHeaderComponent,
				CustomvalidityDirective
			],
			imports: [
				APITestingModule,
				FormsModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule.withRoutes([
					{component: InvalidationJobsComponent, path: "deliveryservice/:id/invalidation-jobs"}
				]),
				MatDialogModule
			],
			providers: [
				{ provide: CurrentUserService, useValue: mockCurrentUserService },
				{ provide: MatDialog, useValue: {open: (): {afterClosed: () => Observable<unknown>} => ({
					afterClosed: () => of(true)
				})}}
			]
		}).compileComponents();

		const dsService = TestBed.inject(DeliveryServiceService);
		const ds = await dsService.createDeliveryService({...defaultDeliveryService});

		router = TestBed.inject(Router);
		router.initialNavigation();
		const navigated = await router.navigate(["/deliveryservice", ds.id, "invalidation-jobs"]);
		if (!navigated) {
			return fail("navigation failed");
		}
		expect(router.url).toBe(`/deliveryservice/${ds.id}/invalidation-jobs`);

		const jobService = TestBed.inject(InvalidationJobService);
		job = await jobService.createInvalidationJob({
			deliveryService: ds.xmlId,
			regex: "/",
			startTime: new Date(),
			ttl: 178
		});

		fixture = TestBed.createComponent(InvalidationJobsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	afterEach(async ()=> {
		await expectAsync(component.deleteJob(job.id)).toBeResolved();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("determines in-progress state", ()=>{
		const j = {
			assetUrl: "",
			createdBy: "",
			deliveryService: "",
			id: -1,
			keyword: JobType.PURGE,
			parameters: "TTL:1h",
			startTime: new Date(component.now)
		};
		j.startTime.setDate(j.startTime.getDate()-1);
		expect(component.isInProgress(j)).toBeFalse();
		j.startTime = new Date(component.now);
		j.startTime.setMinutes(j.startTime.getMinutes()-30);
		expect(component.isInProgress(j)).toBeTrue();
		j.startTime.setMinutes(j.startTime.getMinutes()+31);
		expect(component.isInProgress(j)).toBeFalse();
	});

	it("calculates end dates", () => {
		const j = {
			assetUrl: "doesn't matter",
			createdBy: "also doesn't matter",
			deliveryService: "doesn't matter either",
			id: -1,
			keyword: JobType.PURGE,
			parameters: "",
			startTime: new Date(0),
		};
		expect(()=>component.endDate(j)).toThrow();

		j.parameters = "TTL";
		expect(()=>component.endDate(j)).toThrow();

		j.parameters = "TTL:not a number";
		expect(()=>component.endDate(j)).toThrow();

		const expected = new Date(0);
		expected.setHours(expected.getHours()+178);
		j.parameters = "TTL:178h";
		expect(component.endDate(j)).toEqual(expected);
	});

	it("opens the create/edit dialog", () => {
		component.editJob(job);
		component.newJob();
	});
});
