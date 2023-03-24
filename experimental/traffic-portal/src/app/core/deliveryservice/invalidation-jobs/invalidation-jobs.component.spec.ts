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
import { type Observable, of, ReplaySubject } from "rxjs";
import { GeoLimit, GeoProvider, JobType, ResponseInvalidationJob } from "trafficops-types";

import { CDNService, DeliveryServiceService, InvalidationJobService, TypeService, UserService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { InvalidationJobsComponent } from "src/app/core/deliveryservice/invalidation-jobs/invalidation-jobs.component";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";
import { TpHeaderComponent } from "src/app/shared/navigation/tp-header/tp-header.component";
import { CustomvalidityDirective } from "src/app/shared/validation/customvalidity.directive";

describe("InvalidationJobsComponent", () => {
	let component: InvalidationJobsComponent;
	let fixture: ComponentFixture<InvalidationJobsComponent>;
	let router: Router;
	let job: ResponseInvalidationJob;

	beforeEach(async () => {
		// mock the API
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "hasPermission", "login", "logout"]);
		const headerSvc = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});

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
				})}},
				{ provide: NavigationService, useValue: headerSvc}
			]
		}).compileComponents();

		const dsService = TestBed.inject(DeliveryServiceService);
		const cdnService = TestBed.inject(CDNService);
		const cdn = (await cdnService.getCDNs()).find(c => c.name !== "ALL");
		if (!cdn) {
			throw new Error("can't test a DS card component without any CDNs");
		}
		const typeService = TestBed.inject(TypeService);
		const type = (await typeService.getTypesInTable("deliveryservice")).find(t => t.name === "ANY_MAP");
		if (!type) {
			throw new Error("can't test a DS card component without DS types");
		}
		const tenantService = TestBed.inject(UserService);
		const tenant = (await tenantService.getTenants())[0];

		const ds = await dsService.createDeliveryService({
			active: false,
			anonymousBlockingEnabled: false,
			cacheurl: null,
			cdnId: cdn.id,
			displayName: "FIZZbuzz",
			dscp: 0,
			geoLimit: GeoLimit.NONE,
			geoProvider: GeoProvider.MAX_MIND,
			httpBypassFqdn: null,
			infoUrl: null,
			ipv6RoutingEnabled: true,
			logsEnabled: true,
			longDesc: "",
			missLat: 0,
			missLong: 0,
			multiSiteOrigin: false,
			regionalGeoBlocking: false,
			remapText: null,
			tenantId: tenant.id,
			typeId: type.id,
			xmlId: "fizz-buzz",
		});

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
			invalidationType: JobType.REFRESH,
			regex: "/",
			startTime: new Date(),
			ttlHours: 178
		});

		fixture = TestBed.createComponent(InvalidationJobsComponent);
		component = fixture.componentInstance;
		component.deliveryservice = ds;
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
			invalidationType: JobType.REFETCH,
			startTime: new Date(component.now),
			ttlHours: 1
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
			invalidationType: JobType.REFETCH,
			startTime: new Date(0),
			ttlHours: 178,
		};

		const expected = new Date(0);
		expected.setHours(expected.getHours()+178);
		expect(component.endDate(j)).toEqual(expected);
	});

	it("opens the create/edit dialog", () => {
		component.editJob(job);
		component.newJob();
	});
});
