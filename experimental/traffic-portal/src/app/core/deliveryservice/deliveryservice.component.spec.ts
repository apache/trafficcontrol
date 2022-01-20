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
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import { LinechartDirective } from "src/app/shared/charts/linechart.directive";
// import { DeliveryService, GeoLimit, GeoProvider, TPSData } from "src/app/models";
import { TpHeaderComponent } from "src/app/shared/tp-header/tp-header.component";
import { UserService } from "src/app/shared/api";
import { AlertService } from "src/app/shared/alert/alert.service";

import { DeliveryserviceComponent } from "./deliveryservice.component";


describe("DeliveryserviceComponent", () => {
	let component: DeliveryserviceComponent;
	let fixture: ComponentFixture<DeliveryserviceComponent>;

	beforeEach(waitForAsync(() => {
		// mock the API
		const mockAPIService = jasmine.createSpyObj(["getUsers"]);
		const mockAlertService = jasmine.createSpyObj(["newAlert"]);
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);

		TestBed.configureTestingModule({
			declarations: [
				DeliveryserviceComponent,
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
				{ provide: AlertService, useValue: mockAlertService },
				{ provide: CurrentUserService, useValue: mockCurrentUserService },
				{ provide: UserService, useValue: mockAPIService }
			]
		});
		TestBed.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(DeliveryserviceComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should exist", () => {
		expect(component).toBeTruthy();
	});

	afterAll(() => {
		try{
			TestBed.resetTestingModule();
		} catch (e) {
			console.error("error in DeliveryServiceComponent afterAll:", e);
		}
	});
});
