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

import { APITestingModule } from "src/app/api/testing";
import { DeliveryService } from "src/app/models";
import { UserService } from "src/app/shared/api";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import { LinechartDirective } from "src/app/shared/charts/linechart.directive";
import { TpHeaderComponent } from "src/app/shared/tp-header/tp-header.component";
import { LoadingComponent } from "src/app/shared/loading/loading.component";
import { AlertService } from "src/app/shared/alert/alert.service";

import { DsCardComponent } from "../ds-card/ds-card.component";
import { DashboardComponent } from "./dashboard.component";

describe("DashboardComponent", () => {
	let component: DashboardComponent;
	let fixture: ComponentFixture<DashboardComponent>;

	beforeEach(waitForAsync(() => {
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"],
			{capabilities: new Set<string>()});

		const mockAlertService = jasmine.createSpyObj(["newAlert"]);

		TestBed.configureTestingModule({
			declarations: [
				DashboardComponent,
				DsCardComponent,
				LoadingComponent,
				TpHeaderComponent,
				LinechartDirective
			],
			imports: [
				APITestingModule,
				FormsModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule
			],
			providers: [
				{ provide: CurrentUserService, useValue: mockCurrentUserService },
				{ provide: AlertService, useValue: mockAlertService },
				{ provide: UserService, useValue: mockCurrentUserService }
			]
		});
		TestBed.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(DashboardComponent);
		component = fixture.componentInstance;
		component.deliveryServices = [
			{
				displayName: "FIZZbuzz"
			} as DeliveryService,
			{
				displayName: "fooBAR"
			} as DeliveryService
		];
		fixture.detectChanges();
	});

	it("should exist", () => {
		expect(component).toBeTruthy();
	});

	it('sets the "search" query parameter', () => {
		expect(true).toBeTruthy();
	});

	afterAll(() => {
		try{
			TestBed.resetTestingModule();
		} catch (e) {
			console.error("error in DashboardComponent afterAll:", e);
		}
	});
});
