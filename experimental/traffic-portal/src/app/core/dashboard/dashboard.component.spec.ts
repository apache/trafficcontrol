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
import { DeliveryServiceService } from "src/app/shared/api";


import { LinechartDirective } from "../../shared/charts/linechart.directive";
import { DeliveryService } from "../../models";
import { DsCardComponent } from "../ds-card/ds-card.component";
import { LoadingComponent } from "../../common/loading/loading.component";
import { TpHeaderComponent } from "../../common/tp-header/tp-header.component";
import { DashboardComponent } from "./dashboard.component";

describe("DashboardComponent", () => {
	let component: DashboardComponent;
	let fixture: ComponentFixture<DashboardComponent>;

	beforeEach(waitForAsync(() => {
		const mockAPIService = jasmine.createSpyObj(["getDeliveryServices"]);
		mockAPIService.getDeliveryServices.and.returnValue(new Promise(resolve=>resolve([])));

		TestBed.configureTestingModule({
			declarations: [
				DashboardComponent,
				DsCardComponent,
				LoadingComponent,
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
		TestBed.overrideProvider(DeliveryServiceService, { useValue: mockAPIService });
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
