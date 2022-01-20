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

import { APITestingModule } from "src/app/api/testing";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import { CustomvalidityDirective } from "src/app/shared/validation/customvalidity.directive";
import { JobType } from "src/app/models";
import { TpHeaderComponent } from "src/app/shared/tp-header/tp-header.component";

import { InvalidationJobsComponent } from "./invalidation-jobs.component";

describe("InvalidationJobsComponent", () => {
	let component: InvalidationJobsComponent;
	let fixture: ComponentFixture<InvalidationJobsComponent>;

	beforeEach(waitForAsync(() => {
		// mock the API
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);

		TestBed.configureTestingModule({
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
				RouterTestingModule,
				MatDialogModule
			],
			providers: [
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
		j.startTime.setMinutes(j.startTime.getMinutes()+31);
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
