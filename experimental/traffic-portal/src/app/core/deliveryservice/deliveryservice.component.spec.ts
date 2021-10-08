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

import { of } from "rxjs";


import { LinechartDirective } from "../../shared/charts/linechart.directive";
import { DeliveryService, GeoLimit, GeoProvider, TPSData } from "../../models";
import {TpHeaderComponent} from "../../shared/tp-header/tp-header.component";
import {DeliveryServiceService, UserService} from "../../shared/api";
import {AlertService} from "../../shared/alert/alert.service";
import {AuthenticationService} from "../../shared/authentication/authentication.service";
import { DeliveryserviceComponent } from "./deliveryservice.component";


describe("DeliveryserviceComponent", () => {
	let component: DeliveryserviceComponent;
	let fixture: ComponentFixture<DeliveryserviceComponent>;

	beforeEach(waitForAsync(() => {
		// mock the API
		const mockAPIService = jasmine.createSpyObj(["getDeliveryServices", "getDSKBPS", "getAllDSTPSData"]);
		const mockAlertService = jasmine.createSpyObj(["newAlert"]);
		const mockAuthenticationService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);
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
		} as DeliveryService ));
		mockAPIService.getDSKBPS.and.returnValue(of({series: {values: []}}), of({series: {values: []}}));
		mockAPIService.getAllDSTPSData.and.returnValue(of({
			clientError: {
				dataSet: {
					data: [],
					label: ""
				},
				fifthPercentile: 0,
				max: 0,
				mean: 0,
				min: 0,
				ninetyEighthPercentile: 0,
				ninetyFifthPercentile: 0
			},
			redirection: {
				dataSet: {
					data: [],
					label: ""
				},
				fifthPercentile: 0,
				max: 0,
				mean: 0,
				min: 0,
				ninetyEighthPercentile: 0,
				ninetyFifthPercentile: 0
			},
			serverError: {
				dataSet: {
					data: [],
					label: ""
				},
				fifthPercentile: 0,
				max: 0,
				mean: 0,
				min: 0,
				ninetyEighthPercentile: 0,
				ninetyFifthPercentile: 0
			},
			success: {
				dataSet: {
					data: [],
					label: ""
				},
				fifthPercentile: 0,
				max: 0,
				mean: 0,
				min: 0,
				ninetyEighthPercentile: 0,
				ninetyFifthPercentile: 0,
			},
			total: {
				dataSet: {
					data: [],
					label: ""
				},
				fifthPercentile: 0,
				max: 0,
				mean: 0,
				min: 0,
				ninetyEighthPercentile: 0,
				ninetyFifthPercentile: 0,
			}
		} as TPSData));

		TestBed.configureTestingModule({
			declarations: [
				DeliveryserviceComponent,
				TpHeaderComponent,
				LinechartDirective
			],
			imports: [
				FormsModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule
			],
			providers: [
				{ provide: DeliveryServiceService, useValue: mockAPIService },
				{ provide: AlertService, useValue: mockAlertService },
				{ provide: AuthenticationService, useValue: mockAuthenticationService },
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
