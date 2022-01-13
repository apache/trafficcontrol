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
import { waitForAsync, ComponentFixture, TestBed } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { RouterTestingModule } from "@angular/router/testing";
import { MatDialogModule } from "@angular/material/dialog";

import { of } from "rxjs";
import {DeliveryServiceService, InvalidationJobService, UserService} from "src/app/shared/api";


import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import { CustomvalidityDirective } from "../../shared/validation/customvalidity.directive";
import { DeliveryService, GeoLimit, GeoProvider, InvalidationJob, JobType } from "../../models";
import {TpHeaderComponent} from "../../shared/tp-header/tp-header.component";
import { InvalidationJobsComponent } from "./invalidation-jobs.component";

describe("InvalidationJobsComponent", () => {
	let component: InvalidationJobsComponent;
	let fixture: ComponentFixture<InvalidationJobsComponent>;

	beforeEach(waitForAsync(() => {
		// mock the API
		const mockAPIService = jasmine.createSpyObj(["getInvalidationJobs", "getDeliveryServices"]);
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);
		mockAPIService.getInvalidationJobs.and.returnValue(of({
			startTime: new Date(),
		} as InvalidationJob));
		mockAPIService.getDeliveryServices.and.returnValue(of({
			active: true,
			anonymousBlockingEnabled: false,
			cdnId: 0,
			displayName: "test DS",
			dscp: 0,
			geoLimit: GeoLimit.NONE,
			geoProvider: GeoProvider.MAX_MIND,
			ipv6RoutingEnabled: true,
			lastUpdated: new Date(),
			logsEnabled: true,
			longDesc: "A test Delivery Service for API mock-ups",
			missLat: 0,
			missLong: 0,
			multiSiteOrigin: false,
			regionalGeoBlocking: false,
			routingName: "test-DS",
			typeId: 0,
			xmlId: "test-DS"
		} as DeliveryService));

		TestBed.configureTestingModule({
			declarations: [
				InvalidationJobsComponent,
				TpHeaderComponent,
				CustomvalidityDirective
			],
			imports: [
				FormsModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule,
				MatDialogModule
			],
			providers: [
				{ provide: InvalidationJobService, useValue: mockAPIService },
				{ provide: DeliveryServiceService, useValue: mockAPIService },
				{ provide: UserService, useValue: mockAPIService },
				{ provide: CurrentUserService, useValue: mockCurrentUserService }
			]
		});

		TestBed.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(InvalidationJobsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
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
		j.startTime = new Date();
		expect(component.isInProgress(j)).toBeFalse();
	});

	afterAll(() => {
		try{
			TestBed.resetTestingModule();
		} catch (e) {
			console.error("error in InvalidationJobsComponent afterAll:", e);
		}
	});
});
