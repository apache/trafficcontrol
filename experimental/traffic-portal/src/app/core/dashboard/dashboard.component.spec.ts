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
import { type ComponentFixture, TestBed, fakeAsync, tick } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { Router } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import {ReplaySubject} from "rxjs";

import { DeliveryServiceService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { defaultDeliveryService } from "src/app/models";
import { AlertService } from "src/app/shared/alert/alert.service";
import { LinechartDirective } from "src/app/shared/charts/linechart.directive";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import { LoadingComponent } from "src/app/shared/loading/loading.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";
import { TpHeaderComponent } from "src/app/shared/navigation/tp-header/tp-header.component";

import { DsCardComponent } from "../deliveryservice/ds-card/ds-card.component";

import { DashboardComponent } from "./dashboard.component";

describe("DashboardComponent", () => {
	let component: DashboardComponent;
	let fixture: ComponentFixture<DashboardComponent>;
	let router: Router;

	beforeEach(async () => {
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "hasPermission", "login", "logout"],
			{capabilities: new Set<string>()});
		const navSvc = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});

		await TestBed.configureTestingModule({
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
				RouterTestingModule.withRoutes([
					{component: DashboardComponent, path: ""}
				])
			],
			providers: [
				{ provide: CurrentUserService, useValue: mockCurrentUserService },
				AlertService,
				{ provide: NavigationService, useValue: navSvc}
			]
		}).compileComponents();
		const service = TestBed.inject(DeliveryServiceService);
		const dss = [
			await service.createDeliveryService({
				...defaultDeliveryService,
				displayName: "FIZZbuzz",
				xmlId: "fizz-buzz"
			}),
			await service.createDeliveryService({
				...defaultDeliveryService,
				displayName: "fooBAR",
				xmlId: "foo-bar"
			})
		];

		router = TestBed.inject(Router);
		router.initialNavigation();

		fixture = TestBed.createComponent(DashboardComponent);
		component = fixture.componentInstance;
		component.deliveryServices = dss;
		fixture.detectChanges();
	});

	it("should exist", () => {
		expect(component).toBeTruthy();
	});

	it('sets the "search" query parameter', fakeAsync(() => {
		expect(router.url).toBe("/");
		component.fuzzControl.setValue("query");
		component.updateURL(new Event(""));
		tick();
		expect(router.url).toBe("/?search=query");
		component.fuzzControl.setValue("");
		component.updateURL(new Event(""));
		tick();
		expect(router.url).toBe("/");
	}));

	it("filters Delivery Services", () => {
		expect(component.fuzzControl.value).toBe("");
		expect(component.filteredDSes).toEqual(component.deliveryServices);

		component.fuzzControl.setValue("fz");
		expect(component.filteredDSes.map(d=>d.displayName)).toEqual(["FIZZbuzz"]);

		component.fuzzControl.setValue("aoeu");
		expect(component.filteredDSes).toEqual([]);

		component.fuzzControl.setValue("fb");
		expect(component.filteredDSes.map(d=>d.displayName)).toEqual(["fooBAR", "FIZZbuzz"]);
	});
});
