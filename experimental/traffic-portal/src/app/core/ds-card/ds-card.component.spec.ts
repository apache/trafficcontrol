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
import { RouterTestingModule } from "@angular/router/testing";

import { APITestingModule } from "src/app/api/testing";
import type { DeliveryService } from "src/app/models";
import { LinechartDirective } from "src/app/shared/charts/linechart.directive";
import { LoadingComponent } from "src/app/shared/loading/loading.component";

import { DsCardComponent } from "./ds-card.component";

describe("DsCardComponent", () => {
	let component: DsCardComponent;
	let fixture: ComponentFixture<DsCardComponent>;

	beforeEach(waitForAsync(() => {

		TestBed.configureTestingModule({
			declarations: [
				DsCardComponent,
				LoadingComponent,
				LinechartDirective
			],
			imports: [
				APITestingModule,
				HttpClientModule,
				RouterTestingModule,
			],
		});
		TestBed.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(DsCardComponent);
		component = fixture.componentInstance;
		component.deliveryService = {xmlId: "test-ds"} as DeliveryService;
		fixture.detectChanges();
	});

	it("should exist", () => {
		expect(component).toBeTruthy();
	});

	afterAll(() => {
		try{
			TestBed.resetTestingModule();
		} catch (e) {
			console.error("error in DSCardComponent afterAll:", e);
		}
	});
});
