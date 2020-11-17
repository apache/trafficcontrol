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

import { DeliveryserviceComponent } from "./deliveryservice.component";

import { LinechartDirective } from "../../directives/linechart.directive";
import { DeliveryService, GeoLimit, GeoProvider, TPSData } from "../../models";
import { APIService } from "../../services/api.service";
import { TpHeaderComponent } from "../tp-header/tp-header.component";


describe("DeliveryserviceComponent", () => {
	let component: DeliveryserviceComponent;
	let fixture: ComponentFixture<DeliveryserviceComponent>;

	beforeEach(waitForAsync(() => {
		// mock the API
		const mockAPIService = jasmine.createSpyObj(["getDeliveryServices", "getDSKBPS", "getAllDSTPSData"]);
		mockAPIService.getDeliveryServices.and.returnValue(of({
			active: true,
			anonymousBlockingEnabled: false,
			cdnId: 0,
			displayName: "test DS",
			dscp: 0,
			geoLimit: GeoLimit.None,
			geoProvider: GeoProvider.MaxMind,
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
					data: []
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
					data: []
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
					data: []
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
					data: []
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
					data: []
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
			]
		});
		TestBed.overrideProvider(APIService, { useValue: mockAPIService });
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

	it('sets the "to" and "from" values to "so far today"', () => {
		const now = new Date();
		now.setUTCMilliseconds(0);
		const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());

		expect(component.to).toEqual(now);
		expect(component.from).toEqual(today);

		const nowDate = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, "0")}-${String(now.getDate()).padStart(2, "0")}`;
		const nowTime = `${String(now.getHours()).padStart(2, "0")}:${String(now.getMinutes()).padStart(2, "0")}`;
		expect(nowDate).toEqual(component.toDate.value);
		expect(nowTime).toEqual(component.toTime.value);
		expect(nowDate).toEqual(component.fromDate.value);
		expect("00:00").toEqual(component.fromTime.value);
	});

	afterAll(() => {
		TestBed.resetTestingModule();
	});
});
